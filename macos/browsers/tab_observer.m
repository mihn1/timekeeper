#import <Cocoa/Cocoa.h>
#include <Foundation/Foundation.h>
#include <stdlib.h>
extern void goTabChangeCallback(const char *tabInfo, const char *browserName);

// Define a struct to hold the observer context
typedef struct {
  char *name;
  char *script;
} ObserverContext;

typedef struct {
  AXObserverRef observer;
  ObserverContext *context;
  int pid;
} ObserverData;

#define MAX_OBSERVERS 32
static ObserverData *observerList[MAX_OBSERVERS];
static int observerCount = 0;

// Add an observer to the array.
static void addObserver(ObserverData *data) {
  if (observerCount < MAX_OBSERVERS) {
    observerList[observerCount++] = data;
  } else {
    NSLog(@"❌ Maximum observers reached");
  }
}

// Remove an observer from the array.
static void removeObserver(ObserverData *data) {
  int index = -1;
  for (int i = 0; i < observerCount; i++) {
    if (observerList[i] == data) {
      index = i;
      break;
    }
  }
  if (index >= 0) {
    for (int j = index; j < observerCount - 1; j++) {
      observerList[j] = observerList[j + 1];
    }
    observerCount--;
  }
}

// Cleanup for a single observer.
void cleanupObserver(ObserverData *data) {
  if (!data)
    return;
  
  if (data->context) {
    NSLog(@"🧹 Cleaning up context for %s", data->context->name);
  } else {
    NSLog(@"🧹 Cleaning up observer");
  }
  
  if (data->observer) {
    CFRunLoopSourceRef source = AXObserverGetRunLoopSource(data->observer);
    if (source) {
      CFRunLoopRemoveSource(CFRunLoopGetCurrent(), source, kCFRunLoopDefaultMode);
    }
    CFRelease(data->observer);
  }
  
  if (data->context) {
    free(data->context->name);
    free(data->context->script);
    free(data->context);
  }
  
  free(data);
}

// Cleanup an observer by name
void cleanupObserverByName(const char *name) {
  for (int i = 0; i < observerCount; i++) {
    if (observerList[i] && observerList[i]->context &&
        observerList[i]->context->name &&
        strcmp(observerList[i]->context->name, name) == 0) {
      
      ObserverData *data = observerList[i];
      removeObserver(data);
      cleanupObserver(data);
    }
  }
}

// Cleanup all observers.
void cleanupAllObservers(void) {
  for (int i = 0; i < MAX_OBSERVERS; i++) {
    if (observerList[i] != NULL) {
      cleanupObserver(observerList[i]);
      observerList[i] = NULL; // No need to call remove observer as we're cleaning up all
    }
  }
  observerCount = 0;
}

// Returns tab data via AppleScript.
NSString *getTabData(const char *browserName, const char *script) {
  NSString *nsScript = [NSString stringWithUTF8String:script];
  NSAppleScript *appleScript = [[NSAppleScript alloc] initWithSource:nsScript];
  NSDictionary *errorInfo = nil;
  NSAppleEventDescriptor *result = [appleScript executeAndReturnError:&errorInfo];
  if (errorInfo) {
    NSLog(@"AppleScript Error for %s: %@", browserName, errorInfo);
    return NULL;
  }
  return [result stringValue];
}

// Recursively finds the first element by role.
AXUIElementRef findElementByRole(AXUIElementRef element, CFStringRef targetRole) {
  CFStringRef role;
  if (AXUIElementCopyAttributeValue(element, kAXRoleAttribute, (CFTypeRef *)&role) == kAXErrorSuccess) {
    if (CFStringCompare(role, targetRole, 0) == kCFCompareEqualTo) {
      CFRelease(role);
      return (AXUIElementRef)CFRetain(element); // Ensure it stays valid
    }
    CFRelease(role);
  }
  
  CFArrayRef children;
  if (AXUIElementCopyAttributeValue(element, kAXChildrenAttribute, (CFTypeRef *)&children) ==
      kAXErrorSuccess && children) {
    for (CFIndex i = 0; i < CFArrayGetCount(children); i++) {
      AXUIElementRef child = (AXUIElementRef)CFArrayGetValueAtIndex(children, i);
      AXUIElementRef foundElement = findElementByRole(child, targetRole);
      if (foundElement) {
        CFRelease(children);
        return foundElement;
      }
    }
    CFRelease(children);
  }
  return NULL;
}

// Callback for tab changes.
void tabChangeCallback(AXObserverRef observer, AXUIElementRef element,
                       CFStringRef notification, void *context) {
  if (!context) {
    printf("❌ Context is NULL. Ignoring event.\n");
    return;
  }
  ObserverContext *ctx = (ObserverContext *)context;
  if (!ctx->name) {
    printf("❌ Context->name is NULL. Ignoring event.\n");
    return;
  }
  
  NSString *role = NULL;
  CFTypeRef roleValue = NULL;
  if (AXUIElementCopyAttributeValue(element, kAXRoleAttribute, &roleValue) ==
      kAXErrorSuccess && roleValue) {
    role = (__bridge NSString *)roleValue;
    CFRelease(roleValue);
  } else {
    return;
  }
  
  if (CFStringCompare(notification, kAXTitleChangedNotification, 0) == kCFCompareEqualTo) {
    if (![role isEqual:@"AXWindow"]) {
      return;
    }
  }
  
  if (CFStringCompare(notification, kAXFocusedUIElementChangedNotification, 0) == kCFCompareEqualTo) {
    if (![role isEqual:@"AXWebArea"]) {
      return;
    }
  }
  
  NSString *result = getTabData(ctx->name, ctx->script);
  if (!result)
    return;
  
  char *finalInfo = strdup([result UTF8String]);
  char *browserNameDup = strdup(ctx->name);
  
  goTabChangeCallback(finalInfo, browserNameDup);
  
  free(finalInfo);
  free(browserNameDup);
}

// Registers all accessibility notifications.
int registerAllAXEvents(AXObserverRef observer, AXUIElementRef appElement,
                         ObserverContext *context) {
  CFStringRef events[] = {
    kAXTitleChangedNotification,
    kAXFocusedUIElementChangedNotification,
    kAXFocusedWindowChangedNotification,
  };
  int successCount = 0;
  size_t eventCount = sizeof(events) / sizeof(events[0]);
  for (size_t i = 0; i < eventCount; i++) {
    AXError error = AXObserverAddNotification(observer, appElement, events[i], context);
    if (error == kAXErrorSuccess) {
      NSLog(@"✅ Listening for event: %@", events[i]);
      successCount++;
    } else {
        NSLog(@"❌ Failed to observe event: %@ (code %d)", 
            events[i], 
            error);    
    }
  }
  return successCount;
}

// Starts the observer for a given application.
int startTabObserver(int pid, const char *browserName, const char *script) {
  AXUIElementRef appElement = AXUIElementCreateApplication(pid);
  if (!appElement) {
    NSLog(@"❌ Failed to get main app AXUIElement");
    return 0;
  }
  
  ObserverData *data = (ObserverData *)malloc(sizeof(ObserverData));
  data->context = (ObserverContext *)malloc(sizeof(ObserverContext));
  if (!data->context) {
    NSLog(@"❌ Failed to allocate context");
    free(data);
    CFRelease(appElement);
    return 0;
  }
  data->context->name = strdup(browserName);
  data->context->script = strdup(script);
  
  if (AXObserverCreate(pid, tabChangeCallback, &data->observer) != kAXErrorSuccess) {
    NSLog(@"❌ Failed to create AXObserver");
    cleanupObserver(data);
    CFRelease(appElement);
    return 0;
  }
  
  CFRunLoopAddSource(CFRunLoopGetCurrent(),
                     AXObserverGetRunLoopSource(data->observer),
                     kCFRunLoopDefaultMode);
  
  int eventCount = registerAllAXEvents(data->observer, appElement, data->context);
  if (eventCount == 0) {
    NSLog(@"❌ Failed to register any events");
    cleanupObserver(data);
    CFRelease(appElement);
    return 0;
  }
  
  data->pid = pid;
  addObserver(data);
  
  // atexit(cleanupAllObservers);
  
  CFRelease(appElement);

  return 1;
}