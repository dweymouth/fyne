//go:build darwin

#import <AppKit/AppKit.h>
#import <AppKit/NSEvent.h>
#import <IOKit/pwr_mgt/IOPMLib.h>

extern void fyneAppReopened(void);

@interface FyneReopenDelegate : NSObject<NSApplicationDelegate>
@property (nonatomic, strong) id<NSApplicationDelegate> wrapped;
@end

@implementation FyneReopenDelegate

- (BOOL)applicationShouldHandleReopen:(NSApplication *)app
                    hasVisibleWindows:(BOOL)hasVisible {
    if (!hasVisible) {
        fyneAppReopened();
    }
    if ([self.wrapped respondsToSelector:_cmd]) {
        return [self.wrapped applicationShouldHandleReopen:app
                                        hasVisibleWindows:hasVisible];
    }
    return YES;
}

- (BOOL)respondsToSelector:(SEL)sel {
    return [super respondsToSelector:sel] ||
           [self.wrapped respondsToSelector:sel];
}

- (id)forwardingTargetForSelector:(SEL)sel {
    return [self.wrapped respondsToSelector:sel] ? self.wrapped : nil;
}

@end

void installFyneReopenDelegate(void) {
    FyneReopenDelegate *d = [FyneReopenDelegate new];
    d.wrapped = [NSApp delegate];
    [NSApp setDelegate:d];
}

IOPMAssertionID currentDisableID;

void setDisableDisplaySleep(BOOL disable) {
    if (!disable) {
        if (currentDisableID == 0) {
            return;
        }

        IOPMAssertionRelease(currentDisableID);
        currentDisableID = 0;
        return;
    }

    if (currentDisableID != 0) {
        return;
    }
    IOPMAssertionID assertionID;
    IOReturn success = IOPMAssertionCreateWithName(kIOPMAssertionTypeNoDisplaySleep,
        kIOPMAssertionLevelOn, (CFStringRef)@"App disabled screensaver", &assertionID);

    if (success == kIOReturnSuccess) {
        currentDisableID = assertionID;
    }
}

// https://developer.apple.com/documentation/appkit/nsevent/1528384-doubleclickinterval?language=objc
double doubleClickInterval() {
    return [NSEvent doubleClickInterval];
}
