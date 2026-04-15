package core

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework ApplicationServices -framework AppKit
#include "window_controller_darwin.h"
*/
import "C"

// CheckAccessibilityPermission은 손쉬운 사용 권한 여부를 반환합니다.
func CheckAccessibilityPermission() bool {
	return bool(C.checkAccessibilityPermission())
}

// CheckInputMonitoringPermission은 입력 모니터링 권한 여부를 반환합니다.
func CheckInputMonitoringPermission() bool {
	return bool(C.checkInputMonitoringPermission())
}

// GetActiveAppPID는 현재 포커스된 앱의 PID를 반환합니다.
func GetActiveAppPID() int {
	return int(C.getActiveAppPID())
}

// GetWindowFrame은 PID에 해당하는 앱의 활성 창 위치/크기를 반환합니다.
func GetWindowFrame(pid int) Rect {
	f := C.getWindowFrame(C.pid_t(pid))
	return Rect{
		X: float64(f.x),
		Y: float64(f.y),
		W: float64(f.width),
		H: float64(f.height),
	}
}

// SetWindowFrame은 PID에 해당하는 앱의 활성 창 위치/크기를 변경합니다.
func SetWindowFrame(pid int, r Rect) {
	C.setWindowFrame(C.pid_t(pid), C.float(r.X), C.float(r.Y), C.float(r.W), C.float(r.H))
}

// ActivateApp은 앱을 최상위로 올려 포커스를 재부여합니다.
func ActivateApp(pid int) {
	C.activateApp(C.pid_t(pid))
}
