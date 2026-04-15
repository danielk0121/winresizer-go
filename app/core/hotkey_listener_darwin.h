#pragma once
#include <Carbon/Carbon.h>
#include <CoreFoundation/CoreFoundation.h>
#include <stdbool.h>

// Go 콜백 함수 타입 (단축키 감지 시 호출)
typedef void (*HotkeyCallback)(int hotkeyID);

// 단축키 핸들러 등록 및 Carbon 이벤트 루프 시작
// callback: 단축키 감지 시 호출될 Go 함수 포인터
void startHotkeyListener(HotkeyCallback callback);

// 단축키를 등록합니다. 성공 시 true 반환.
bool registerHotkey(int hotkeyID, UInt32 keyCode, UInt32 modifiers);

// 특정 ID의 단축키를 해제합니다.
void unregisterHotkey(int hotkeyID);

// 등록된 모든 단축키를 해제합니다.
void unregisterAllHotkeys(void);

// Carbon 이벤트 루프를 종료합니다.
void stopHotkeyListener(void);
