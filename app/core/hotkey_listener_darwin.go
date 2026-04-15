package core

/*
#cgo LDFLAGS: -framework Carbon
#include "hotkey_listener_darwin.h"

// Go 콜백을 C 함수 포인터로 전달하기 위한 브릿지
extern void goHotkeyCallback(int hotkeyID);
*/
import "C"
import (
	"sync"
	"winresizer/utils"
)

// hotkeyCallback은 단축키 감지 시 실행할 함수입니다.
var (
	hotkeyCallbackFn func(id int)
	hotkeyCallbackMu sync.Mutex
)

// goHotkeyCallback은 C에서 호출되는 Go 콜백 브릿지입니다.
//
//export goHotkeyCallback
func goHotkeyCallback(hotkeyID C.int) {
	hotkeyCallbackMu.Lock()
	fn := hotkeyCallbackFn
	hotkeyCallbackMu.Unlock()
	if fn != nil {
		fn(int(hotkeyID))
	}
}

// StartHotkeyListener는 Carbon 이벤트 루프를 시작합니다.
// 단축키 감지 시 callback(hotkeyID)을 호출합니다.
// 이 함수는 블로킹이므로 별도 goroutine에서 호출해야 합니다.
func StartHotkeyListener(callback func(id int)) {
	hotkeyCallbackMu.Lock()
	hotkeyCallbackFn = callback
	hotkeyCallbackMu.Unlock()

	utils.Log.Infof("Carbon 단축키 리스너 시작")
	C.startHotkeyListener(C.HotkeyCallback(C.goHotkeyCallback))
}

// RegisterHotkey는 단축키를 등록합니다.
// id: 콜백에서 식별할 고유 번호, keyCode/modifiers: Carbon 키코드 및 수식키
func RegisterHotkey(id int, keyCode uint32, modifiers uint32) bool {
	ok := C.registerHotkey(C.int(id), C.UInt32(keyCode), C.UInt32(modifiers))
	if !bool(ok) {
		utils.Log.Warnf("단축키 등록 실패: id=%d keyCode=%d modifiers=%d", id, keyCode, modifiers)
	}
	return bool(ok)
}

// UnregisterHotkey는 특정 ID의 단축키를 해제합니다.
func UnregisterHotkey(id int) {
	C.unregisterHotkey(C.int(id))
}

// UnregisterAllHotkeys는 등록된 모든 단축키를 해제합니다.
func UnregisterAllHotkeys() {
	C.unregisterAllHotkeys()
}

// StopHotkeyListener는 Carbon 이벤트 루프를 종료합니다.
func StopHotkeyListener() {
	utils.Log.Infof("Carbon 단축키 리스너 종료")
	C.stopHotkeyListener()
}
