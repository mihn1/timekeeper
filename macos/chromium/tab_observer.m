// #import <ApplicationServices/ApplicationServices.h>
#import <Cocoa/Cocoa.h>
// #include <CoreFoundation/CoreFoundation.h>
// #include <Foundation/Foundation.h>
// #include <stdlib.h>
extern void goTabChangeCallback(const char* tabInfo, const char* browserName);

// Define a struct to hold the observer context
typedef struct {
  char *name;
} ObserverContext;

NSString *getTabData(const char *browserName) {
  // Run AppleScript to get tab details
  NSString *script = [NSString stringWithFormat:
                      @"tell application \"%s\"\n"
                      "set frontTab to active tab of front window\n"
                      "set tabTitle to title of frontTab\n"
                      "set tabURL to URL of frontTab\n"
                      "return tabURL & \"|\" & tabTitle\n"
                      "end tell", browserName];

  NSAppleScript *appleScript = [[NSAppleScript alloc] initWithSource:script];
  NSDictionary *errorInfo = nil;
  NSAppleEventDescriptor *result =
      [appleScript executeAndReturnError:&errorInfo];

  if (errorInfo) {
    NSLog(@"AppleScript Error for %s: %@", browserName, errorInfo);
    return NULL;
  }

  return [result stringValue];
}

// Recursively find the first element by role (e.g., "AXWebArea")
AXUIElementRef findElementByRole(AXUIElementRef element,
                                 CFStringRef targetRole) {
  CFStringRef role;
  if (AXUIElementCopyAttributeValue(element, kAXRoleAttribute,
                                    (CFTypeRef *)&role) == kAXErrorSuccess) {
    if (CFStringCompare(role, targetRole, 0) == kCFCompareEqualTo) {
      CFRelease(role);
      return (AXUIElementRef)CFRetain(element); // Ensure it remains valid
    }
    CFRelease(role);
  }

  // Search children if not found
  CFArrayRef children;
  if (AXUIElementCopyAttributeValue(element, kAXChildrenAttribute,
                                    (CFTypeRef *)&children) ==
          kAXErrorSuccess &&
      children) {
    for (CFIndex i = 0; i < CFArrayGetCount(children); i++) {
      AXUIElementRef child =
          (AXUIElementRef)CFArrayGetValueAtIndex(children, i);
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

// Callback function triggered when the active tab changes
void tabChangeCallback(AXObserverRef observer, AXUIElementRef element,
                       CFStringRef notification, void *context) {
  if (!context) {
    printf("❌ Context is NULL. Ignoring event.\n");
    return;
  }

   // Cast the context back to ObserverContext
    ObserverContext* ctx = (ObserverContext*)context;
    
    if (!ctx->name) {
        printf("❌ Context->name is NULL. Ignoring event.\n");
        return;
    }

  NSString *role = NULL;
  CFTypeRef roleValue = NULL;
  if (AXUIElementCopyAttributeValue(element, kAXRoleAttribute, &roleValue) ==
          kAXErrorSuccess &&
      roleValue) {
    role = (__bridge NSString *)roleValue;
    CFRelease(roleValue);
  } else {
    return;
  }

  // Check notification type is title
  if (CFStringCompare(notification, kAXTitleChangedNotification, 0) ==
      kCFCompareEqualTo) {
    if (![role isEqual:@"AXWindow"]) {
      // MainElement not found for kAXTitleChangedNotification
      return;
    }
  }

  // Check notification type is UIFocusElementChanged
  if (CFStringCompare(notification, kAXFocusedUIElementChangedNotification,
                      0) == kCFCompareEqualTo) {
    if (![role isEqual:@"AXWebArea"]) {
      // MainElement not found for kAXFocusedUIElementChangedNotification
      return;
    }
  }

  NSString *result = getTabData(ctx->name);
  if (!result) {
    return;
  }

  char *finalInfo = strdup([result UTF8String]);
  char *browserName = strdup(ctx->name);

  // dispatch_async(dispatch_get_main_queue(), ^{
  goTabChangeCallback(finalInfo, browserName);
  free(finalInfo);
  free(browserName);
  // });
}

void registerAllAXEvents(AXObserverRef observer, AXUIElementRef appElement,
                         ObserverContext *context) {
  CFStringRef events[] = {
      kAXTitleChangedNotification,
      kAXFocusedUIElementChangedNotification,
      kAXFocusedWindowChangedNotification,
  };

  size_t eventCount = sizeof(events) / sizeof(events[0]);

  for (size_t i = 0; i < eventCount; i++) {
    if (AXObserverAddNotification(observer, appElement, events[i], context) ==
        kAXErrorSuccess) {
      NSLog(@"✅ Listening for event: %@", events[i]);
    } else {
      free(context->name);
      free(context);
      NSLog(@"❌ Failed to observe event: %@", events[i]);
    }
  }
}

void stopObserver(AXObserverRef observer, void* context) {
    if (!observer) return;

    printf("🛑 Stopping observer and freeing memory...\n");

    // Free context memory
    ObserverContext* ctx = (ObserverContext*)context;
    if (ctx) {
        free(ctx->name);
        free(ctx);
    }

    // Remove observer from the run loop
    CFRunLoopRemoveSource(CFRunLoopGetCurrent(), AXObserverGetRunLoopSource(observer), kCFRunLoopDefaultMode);
    CFRelease(observer);
}

// Observer tab changes in the given browser
void startTabObserver(int pid, const char *browserName) {
  AXUIElementRef appElement = AXUIElementCreateApplication(pid);
  if (!appElement) {
    NSLog(@"❌ Failed to get main app AXUIElement");
    return;
  }

  // Allocate memory for context and store the name
  ObserverContext *context = (ObserverContext *)malloc(sizeof(ObserverContext));
  context->name =
      strdup(browserName); // Copy name string to avoid memory issues

  // Set up the observer
  AXObserverRef observer;
  if (AXObserverCreate(pid, tabChangeCallback, &observer) != kAXErrorSuccess) {
    NSLog(@"❌ Failed to create AXObserver");
    free(context->name);
    free(context);
    CFRelease(appElement);
    return;
  }

  // Add observer to the run loop
  CFRunLoopAddSource(CFRunLoopGetCurrent(),
                     AXObserverGetRunLoopSource(observer),
                     kCFRunLoopDefaultMode);

  registerAllAXEvents(observer, appElement, context);

  // CFRelease(observer);
  CFRelease(appElement);
}