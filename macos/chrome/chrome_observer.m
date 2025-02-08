// #import <ApplicationServices/ApplicationServices.h>
#import <Cocoa/Cocoa.h>
// #include <CoreFoundation/CoreFoundation.h>
// #include <Foundation/Foundation.h>
// #include <stdlib.h>
extern void goTabChangeCallback(const char *tabInfo);

NSString *getTabData() {
  // Run AppleScript to get tab details
  NSString *script = @"tell application \"Google Chrome\"\n"
                      "set frontTab to active tab of front window\n"
                      "set tabTitle to title of frontTab\n"
                      "set tabURL to URL of frontTab\n"
                      "return tabURL & \"|\" & tabTitle\n"
                      "end tell";

  NSAppleScript *appleScript = [[NSAppleScript alloc] initWithSource:script];
  NSDictionary *errorInfo = nil;
  NSAppleEventDescriptor *result =
      [appleScript executeAndReturnError:&errorInfo];

  if (errorInfo) {
    NSLog(@"AppleScript Error: %@", errorInfo);
    return @"";
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

  NSString *result = getTabData();
  char *finalInfo = strdup([result UTF8String]);

  // dispatch_async(dispatch_get_main_queue(), ^{
  goTabChangeCallback(finalInfo);
  free(finalInfo);
  // });
}

void registerAllAXEvents(AXObserverRef observer, AXUIElementRef appElement) {
  CFStringRef events[] = {
      kAXTitleChangedNotification, 
      kAXFocusedUIElementChangedNotification,
      kAXFocusedWindowChangedNotification,
  };

  size_t eventCount = sizeof(events) / sizeof(events[0]);

  for (size_t i = 0; i < eventCount; i++) {
    if (AXObserverAddNotification(observer, appElement, events[i], NULL) ==
        kAXErrorSuccess) {
      NSLog(@"✅ Listening for event: %@", events[i]);
    } else {
      NSLog(@"❌ Failed to observe event: %@", events[i]);
    }
  }
}

// Observer chrome
void startTabObserver(int pid) {
  NSLog(@"🚀 Starting Tab Observer for Chrome (PID: %d)", pid);

  // Get Chrome application AXUIElement
  AXUIElementRef appElement = AXUIElementCreateApplication(pid);
  if (!appElement) {
    NSLog(@"❌ Failed to get Chrome AXUIElement");
    return;
  }

  // Set up the observer
  AXObserverRef observer;
  if (AXObserverCreate(pid, tabChangeCallback, &observer) != kAXErrorSuccess) {
    NSLog(@"❌ Failed to create AXObserver");
    CFRelease(appElement);
    return;
  }

  // Add observer to the run loop
  CFRunLoopAddSource(CFRunLoopGetCurrent(),
                     AXObserverGetRunLoopSource(observer),
                     kCFRunLoopDefaultMode);

  registerAllAXEvents(observer, appElement);

  // CFRelease(observer);
  CFRelease(appElement);
}