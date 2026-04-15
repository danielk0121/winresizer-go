package core

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "monitor_darwin.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Monitor는 하나의 모니터 정보를 나타냅니다 (Quartz 좌표계 기준).
type Monitor struct {
	X, Y, Width, Height int
}

// GetAllMonitors는 연결된 모든 모니터의 사용 가능한 영역을 반환합니다.
// 메뉴바와 Dock을 제외한 영역이며, Quartz 좌표계(상단 왼쪽 원점) 기준입니다.
func GetAllMonitors() []Monitor {
	var count C.int
	raw := C.getAllMonitors(&count)
	if raw == nil || count == 0 {
		return nil
	}
	defer C.free(unsafe.Pointer(raw))

	n := int(count)
	// C 배열을 Go 슬라이스로 변환
	cSlice := (*[1 << 10]C.MonitorInfo)(unsafe.Pointer(raw))[:n:n]

	monitors := make([]Monitor, n)
	for i, m := range cSlice {
		monitors[i] = Monitor{
			X:      int(m.x),
			Y:      int(m.y),
			Width:  int(m.width),
			Height: int(m.height),
		}
	}
	return monitors
}

// FindActiveMonitor는 창의 중심점이 속한 모니터를 반환합니다.
// 일치하는 모니터가 없으면 첫 번째 모니터를 반환합니다.
func FindActiveMonitor(monitors []Monitor, windowRect Rect) Monitor {
	if len(monitors) == 0 {
		return Monitor{}
	}
	centerX := windowRect.X + windowRect.W/2
	centerY := windowRect.Y + windowRect.H/2

	for _, m := range monitors {
		if float64(m.X) <= centerX && centerX < float64(m.X+m.Width) &&
			float64(m.Y) <= centerY && centerY < float64(m.Y+m.Height) {
			return m
		}
	}
	return monitors[0]
}
