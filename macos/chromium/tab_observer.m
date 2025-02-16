#import <Cocoa/Cocoa.h>
#include <Foundation/Foundation.h>
#include <stdlib.h>
extern void goTabChangeCallback(const char *tabInfo, const char *browserName);

// Define a struct to hold the observer context
typedef struct {
  char *name;
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
    free(data->context);
  }
  
  free(data);
}

// Cleanup all observers.
void cleanupAllObservers(void) {
  for (int i = 0; i < observerCount; i++) {
    if (observerList[i]) {
      cleanupObserver(observerList[i]);
      observerList[i] = NULL;
    }
  }
  observerCount = 0;
}

// Returns tab data via AppleScript.
NSString *getTabData(const char *browserName) {
  NSString *script = [NSString stringWithFormat:@"tell application \"%s\"\n"
                        "set frontTab to active tab of front window\n"
                        "set tabTitle to title of frontTab\n"
                        "set tabURL to URL of frontTab\n"
                        "return tabURL & \"|\" & tabTitle\n"
                        "end tell", browserName];
  NSAppleScript *appleScript = [[NSAppleScript alloc] initWithSource:script];
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
  
  NSString *result = getTabData(ctx->name);
  if (!result)
    return;
  
  char *finalInfo = strdup([result UTF8String]);
  char *browserNameDup = strdup(ctx->name);
  
  goTabChangeCallback(finalInfo, browserNameDup);
  
  free(finalInfo);
  free(browserNameDup);
}

// Registers all accessibility notifications.
void registerAllAXEvents(AXObserverRef observer, AXUIElementRef appElement,
                         ObserverContext *context) {
  CFStringRef events[] = {
    kAXTitleChangedNotification,
    kAXFocusedUIElementChangedNotification,
    kAXFocusedWindowChangedNotification,
  };
  size_t eventCount = sizeof(events) / sizeof(events[0]);
  for (size_t i = 0; i < eventCount; i++) {
    if (AXObserverAddNotification(observer, appElement, events[i], context) == kAXErrorSuccess) {
      NSLog(@"✅ Listening for event: %@", events[i]);
    } else {
      NSLog(@"❌ Failed to observe event: %@", events[i]);
    }
  }
}

// Starts the observer for a given application.
void startTabObserver(int pid, const char *browserName) {
  AXUIElementRef appElement = AXUIElementCreateApplication(pid);
  if (!appElement) {
    NSLog(@"❌ Failed to get main app AXUIElement");
    return;
  }
  
  ObserverData *data = (ObserverData *)malloc(sizeof(ObserverData));
  data->context = (ObserverContext *)malloc(sizeof(ObserverContext));
  if (!data->context) {
    NSLog(@"❌ Failed to allocate context");
    free(data);
    CFRelease(appElement);
    return;
  }
  data->context->name = strdup(browserName);
  
  if (AXObserverCreate(pid, tabChangeCallback, &data->observer) != kAXErrorSuccess) {
    NSLog(@"❌ Failed to create AXObserver");
    cleanupObserver(data);
    CFRelease(appElement);
    return;
  }
  
  CFRunLoopAddSource(CFRunLoopGetCurrent(),
                     AXObserverGetRunLoopSource(data->observer),
                     kCFRunLoopDefaultMode);
  
  registerAllAXEvents(data->observer, appElement, data->context);
  
  data->pid = pid;
  addObserver(data);
  
  atexit(cleanupAllObservers);
  
  CFRelease(appElement);
}