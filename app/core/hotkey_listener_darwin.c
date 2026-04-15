#include "hotkey_listener_darwin.h"
#include <stdlib.h>
#include <stdio.h>

// 최대 등록 단축키 수
#define MAX_HOTKEYS 64

static EventHotKeyRef  s_hotkeyRefs[MAX_HOTKEYS];
static int             s_hotkeyIDs[MAX_HOTKEYS];
static int             s_hotkeyCount = 0;
static HotkeyCallback  s_callback    = NULL;
static EventHandlerRef s_handlerRef  = NULL;
static bool            s_running     = false;
static CFRunLoopRef    s_runLoop     = NULL;

// Carbon 이벤트 핸들러: 단축키 감지 시 Go 콜백 호출
static OSStatus hotkeyHandler(EventHandlerCallRef handlerRef, EventRef event, void* userData) {
    EventHotKeyID hotkeyID;
    GetEventParameter(event, kEventParamDirectObject, typeEventHotKeyID,
                      NULL, sizeof(hotkeyID), NULL, &hotkeyID);
    if (s_callback) {
        s_callback((int)hotkeyID.id);
    }
    return noErr;
}

void startHotkeyListener(HotkeyCallback callback) {
    s_callback = callback;

    // 이벤트 핸들러 등록
    EventTypeSpec eventType = {kEventClassKeyboard, kEventHotKeyPressed};
    InstallApplicationEventHandler(NewEventHandlerUPP(hotkeyHandler),
                                   1, &eventType, NULL, &s_handlerRef);
    s_running = true;
    s_runLoop = CFRunLoopGetCurrent();
    // CFRunLoop 기반 이벤트 루프 실행 (goroutine에서 호출)
    CFRunLoopRun();
}

bool registerHotkey(int hotkeyID, UInt32 keyCode, UInt32 modifiers) {
    if (s_hotkeyCount >= MAX_HOTKEYS) return false;

    EventHotKeyID hkID = {(OSType)'WRsz', (UInt32)hotkeyID};
    EventHotKeyRef ref = NULL;

    OSStatus status = RegisterEventHotKey(keyCode, modifiers, hkID,
                                          GetApplicationEventTarget(), 0, &ref);
    if (status != noErr) {
        fprintf(stderr, "[hotkey] 등록 실패: hotkeyID=%d status=%d\n", hotkeyID, status);
        return false;
    }

    s_hotkeyRefs[s_hotkeyCount] = ref;
    s_hotkeyIDs[s_hotkeyCount]  = hotkeyID;
    s_hotkeyCount++;
    return true;
}

void unregisterHotkey(int hotkeyID) {
    for (int i = 0; i < s_hotkeyCount; i++) {
        if (s_hotkeyIDs[i] == hotkeyID) {
            UnregisterEventHotKey(s_hotkeyRefs[i]);
            // 배열에서 제거 (마지막 항목과 교체)
            s_hotkeyRefs[i] = s_hotkeyRefs[s_hotkeyCount - 1];
            s_hotkeyIDs[i]  = s_hotkeyIDs[s_hotkeyCount - 1];
            s_hotkeyCount--;
            return;
        }
    }
}

void unregisterAllHotkeys(void) {
    for (int i = 0; i < s_hotkeyCount; i++) {
        UnregisterEventHotKey(s_hotkeyRefs[i]);
    }
    s_hotkeyCount = 0;
}

void stopHotkeyListener(void) {
    unregisterAllHotkeys();
    if (s_running) {
        if (s_runLoop) {
            CFRunLoopStop(s_runLoop);
            s_runLoop = NULL;
        }
        s_running = false;
    }
    if (s_handlerRef) {
        RemoveEventHandler(s_handlerRef);
        s_handlerRef = NULL;
    }
}
