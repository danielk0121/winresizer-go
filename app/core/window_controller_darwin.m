#include "window_controller_darwin.h"
#include <ApplicationServices/ApplicationServices.h>
#import <AppKit/AppKit.h>

bool checkAccessibilityPermission(void) {
    return AXIsProcessTrusted();
}

bool checkInputMonitoringPermission(void) {
    // CGPreflightListenEventAccess: 입력 모니터링 권한 확인 (macOS 10.15+)
    return CGPreflightListenEventAccess();
}

pid_t getActiveAppPID(void) {
    NSRunningApplication* app = [[NSWorkspace sharedWorkspace] frontmostApplication];
    if (!app) return -1;
    return [app processIdentifier];
}

// PID에 해당하는 앱의 첫 번째 일반 창(AXWindow)을 가져옵니다.
static AXUIElementRef _getMainWindow(pid_t pid) {
    AXUIElementRef appRef = AXUIElementCreateApplication(pid);
    if (!appRef) return NULL;

    CFArrayRef windows = NULL;
    AXUIElementCopyAttributeValue(appRef, kAXWindowsAttribute, (CFTypeRef*)&windows);
    CFRelease(appRef);

    if (!windows || CFArrayGetCount(windows) == 0) {
        if (windows) CFRelease(windows);
        return NULL;
    }

    // 첫 번째 창 반환 (retain 하여 호출자가 Release 책임)
    AXUIElementRef win = (AXUIElementRef)CFArrayGetValueAtIndex(windows, 0);
    CFRetain(win);
    CFRelease(windows);
    return win;
}

WindowFrame getWindowFrame(pid_t pid) {
    WindowFrame result = {0, 0, 0, 0};

    AXUIElementRef win = _getMainWindow(pid);
    if (!win) return result;

    // 위치
    AXValueRef posVal = NULL;
    CGPoint pos = {0, 0};
    if (AXUIElementCopyAttributeValue(win, kAXPositionAttribute, (CFTypeRef*)&posVal) == kAXErrorSuccess) {
        AXValueGetValue(posVal, kAXValueCGPointType, &pos);
        CFRelease(posVal);
    }

    // 크기
    AXValueRef sizeVal = NULL;
    CGSize size = {0, 0};
    if (AXUIElementCopyAttributeValue(win, kAXSizeAttribute, (CFTypeRef*)&sizeVal) == kAXErrorSuccess) {
        AXValueGetValue(sizeVal, kAXValueCGSizeType, &size);
        CFRelease(sizeVal);
    }

    CFRelease(win);

    result.x = (float)pos.x;
    result.y = (float)pos.y;
    result.width = (float)size.width;
    result.height = (float)size.height;
    return result;
}

void setWindowFrame(pid_t pid, float x, float y, float width, float height) {
    AXUIElementRef win = _getMainWindow(pid);
    if (!win) return;

    // 위치 설정
    CGPoint pos = {x, y};
    AXValueRef posVal = AXValueCreate(kAXValueCGPointType, &pos);
    AXUIElementSetAttributeValue(win, kAXPositionAttribute, posVal);
    CFRelease(posVal);

    // 크기 설정
    CGSize size = {width, height};
    AXValueRef sizeVal = AXValueCreate(kAXValueCGSizeType, &size);
    AXUIElementSetAttributeValue(win, kAXSizeAttribute, sizeVal);
    CFRelease(sizeVal);

    CFRelease(win);
}

void activateApp(pid_t pid) {
    // 1차: AX API로 kAXFrontmostAttribute 직접 설정 (창 단위 포커스 재부여)
    AXUIElementRef appRef = AXUIElementCreateApplication(pid);
    if (appRef) {
        AXUIElementSetAttributeValue(appRef, kAXFrontmostAttribute, kCFBooleanTrue);
        CFRelease(appRef);
    }

    // 2차: NSRunningApplication activate (앱 레벨 포커스 재부여)
    NSRunningApplication* app = [NSRunningApplication runningApplicationWithProcessIdentifier:pid];
    if (app) {
        [app activateWithOptions:NSApplicationActivateIgnoringOtherApps];
    }
}
