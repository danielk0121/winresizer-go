package core

import (
	"fmt"
	"os/exec"
	"sync"
	"winresizer/utils"
)

// windowStateStore는 창 원래 상태를 저장하는 캐시입니다 (복구용).
var (
	windowStateStore = map[int]Rect{} // key: PID
	windowStateMu    sync.Mutex
)

// HotkeyManager는 설정 기반으로 단축키를 등록/관리합니다.
type HotkeyManager struct {
	mu      sync.Mutex
	idMap   map[int]string // hotkeyID → mode
	running bool
	stopCh  chan struct{}
}

var globalHotkeyManager = &HotkeyManager{
	idMap:  map[int]string{},
	stopCh: make(chan struct{}),
}

// StartHotkeyManager는 설정을 읽어 단축키를 등록하고 리스너를 시작합니다.
func StartHotkeyManager() {
	globalHotkeyManager.start()
}

// RestartHotkeyManager는 단축키 리스너를 재시작합니다 (설정 변경 시 호출).
func RestartHotkeyManager() {
	globalHotkeyManager.restart()
}

func (hm *HotkeyManager) start() {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	cfg, err := GetConfig()
	if err != nil {
		utils.Log.Errorf("설정 로드 실패: %v", err)
		return
	}

	// 단축키 등록
	id := 1
	for name, sc := range cfg.Shortcuts {
		if sc.Keycode == 0 {
			continue // 미설정 단축키 스킵
		}
		if RegisterHotkey(id, sc.Keycode, sc.Modifiers) {
			hm.idMap[id] = sc.Mode
			utils.Log.Debugf("단축키 등록: [%s] id=%d keycode=%d modifiers=%d mode=%s", name, id, sc.Keycode, sc.Modifiers, sc.Mode)
			id++
		}
	}

	hm.running = true

	// Carbon 이벤트 루프는 블로킹이므로 goroutine에서 실행
	go StartHotkeyListener(func(hotkeyID int) {
		hm.mu.Lock()
		mode, ok := hm.idMap[hotkeyID]
		hm.mu.Unlock()
		if ok {
			utils.Log.Debugf("단축키 감지: id=%d mode=%s", hotkeyID, mode)
			if err := ExecuteWindowCommand(mode); err != nil {
				utils.Log.Errorf("창 조절 실패: %v", err)
			}
		}
	})
}

func (hm *HotkeyManager) restart() {
	utils.Log.Infof("단축키 리스너 재시작")
	UnregisterAllHotkeys()
	StopHotkeyListener()

	hm.mu.Lock()
	hm.idMap = map[int]string{}
	hm.running = false
	hm.mu.Unlock()

	InvalidateCache()
	hm.start()
}

// ExecuteWindowCommand는 지정된 모드로 활성 창을 조절합니다.
func ExecuteWindowCommand(mode string) error {
	// 특수 명령: 권한 설정 화면 열기
	switch mode {
	case "open_accessibility":
		return openSystemPrefs("x-apple.systempreferences:com.apple.preference.security?Privacy_Accessibility")
	case "open_input_monitoring":
		return openSystemPrefs("x-apple.systempreferences:com.apple.preference.security?Privacy_ListenEvent")
	}

	// 권한 확인
	if !CheckAccessibilityPermission() {
		utils.Log.Warnf("손쉬운 사용 권한 없음")
		return fmt.Errorf("손쉬운 사용 권한이 필요합니다")
	}

	// 활성 앱 PID
	pid := GetActiveAppPID()
	if pid <= 0 {
		return fmt.Errorf("활성 앱을 찾을 수 없습니다")
	}

	// 현재 창 상태 조회
	currentFrame := GetWindowFrame(pid)
	if currentFrame.W == 0 && currentFrame.H == 0 {
		return fmt.Errorf("활성 창을 찾을 수 없습니다")
	}

	// 복구 명령
	if mode == "restore" {
		windowStateMu.Lock()
		saved, ok := windowStateStore[pid]
		windowStateMu.Unlock()
		if ok {
			SetWindowFrame(pid, saved)
			windowStateMu.Lock()
			delete(windowStateStore, pid)
			windowStateMu.Unlock()
			utils.Log.Infof("창 복구 완료: pid=%d", pid)
		} else {
			utils.Log.Infof("저장된 창 상태 없음: pid=%d", pid)
		}
		ActivateApp(pid)
		return nil
	}

	// 원래 상태 저장 (최초 1회)
	windowStateMu.Lock()
	if _, exists := windowStateStore[pid]; !exists {
		windowStateStore[pid] = currentFrame
	}
	windowStateMu.Unlock()

	// 모니터 판별
	monitors := GetAllMonitors()
	if len(monitors) == 0 {
		return fmt.Errorf("모니터 정보를 가져올 수 없습니다")
	}
	targetMonitor := FindActiveMonitor(monitors, currentFrame)

	// 다음 디스플레이 이동
	if mode == "next_display" {
		return moveToNextDisplay(pid, currentFrame, monitors, targetMonitor)
	}

	// 설정에서 gap 조회
	cfg, _ := GetConfig()
	gap := 0.0
	if cfg != nil {
		gap = float64(cfg.Settings.Gap)
	}

	// 목표 좌표 계산
	screenW := float64(targetMonitor.Width)
	screenH := float64(targetMonitor.Height)
	targetRelative, err := CalculateWindowPosition(screenW, screenH, mode, gap)
	if err != nil {
		utils.Log.Warnf("좌표 계산 실패 (%s): %v", mode, err)
	}

	// 모니터 절대 좌표로 변환
	targetAbs := Rect{
		X: targetRelative.X + float64(targetMonitor.X),
		Y: targetRelative.Y + float64(targetMonitor.Y),
		W: targetRelative.W,
		H: targetRelative.H,
	}

	// 이미 해당 위치에 있으면 스마트 사이클 시도, 사이클 없으면 다음 모니터로 이동
	if isAlreadyAligned(currentFrame, targetAbs, mode, targetMonitor) {
		if cycled := nextCycleMode(mode); cycled != "" {
			utils.Log.Debugf("스마트 사이클: %s → %s", mode, cycled)
			return ExecuteWindowCommand(cycled)
		}
		return moveToNextDisplay(pid, currentFrame, monitors, targetMonitor)
	}

	// 창 이동
	SetWindowFrame(pid, targetAbs)

	// Re-anchoring: 앱 최소 크기 제한으로 인한 화면 밖 짤림 보정
	actual := GetWindowFrame(pid)
	corrected := reanchor(actual, targetMonitor, mode, gap)
	if corrected != actual {
		utils.Log.Infof("Re-anchoring 적용: pid=%d mode=%s", pid, mode)
		SetWindowFrame(pid, corrected)
	}

	// 포커스 재부여 (포커스 유실 방지)
	ActivateApp(pid)
	utils.Log.Infof("창 조절 완료: pid=%d mode=%s", pid, mode)
	return nil
}

// isAlreadyAligned는 창이 이미 목표 위치에 정렬되어 있는지 판별합니다.
// 크기까지 유사한 경우를 기본으로 하고, 앱 최소 크기 제한으로 크기가 다를 때만 엣지 보조 판정을 적용합니다.
func isAlreadyAligned(current, target Rect, mode string, monitor Monitor) bool {
	const posTolerance = 20.0  // 위치 오차 허용
	const sizeTolerance = 20.0 // 크기 오차 허용

	// 1차: 위치 + 크기 모두 유사한지 확인
	if IsSimilar(current, target, posTolerance) {
		return true
	}

	// 2차: 크기가 다를 때(앱 최소 크기 제한)는 엣지만 확인
	// 단, 크기가 목표보다 작은 경우는 정렬 안 된 것으로 봄 (목표 크기에 못 미치면 아직 이동 필요)
	posAligned := false
	switch {
	case contains(mode, "right"):
		currRight := current.X + current.W
		expRight := target.X + target.W
		// 우측 엣지가 맞거나, 앱 최소 크기로 인해 화면 밖으로 나간 경우
		posAligned = (abs64(currRight-expRight) <= posTolerance || currRight > expRight+5) &&
			abs64(current.Y-target.Y) <= posTolerance
	case contains(mode, "bottom"):
		currBottom := current.Y + current.H
		expBottom := target.Y + target.H
		posAligned = (abs64(currBottom-expBottom) <= posTolerance || currBottom > expBottom+5) &&
			abs64(current.X-target.X) <= posTolerance
	case contains(mode, "left"), contains(mode, "top"), mode == "maximize":
		posAligned = abs64(current.X-target.X) <= posTolerance &&
			abs64(current.Y-target.Y) <= posTolerance
	}

	if !posAligned {
		return false
	}

	// 위치는 맞는데 크기가 다른 경우: 앱 최소 크기 제한으로 인한 것인지 확인
	// 목표보다 창이 더 큰 경우만 "이미 정렬됨"으로 간주 (최소 크기 제한 때문에 더 커진 것)
	return current.W >= target.W-sizeTolerance && current.H >= target.H-sizeTolerance
}

// reanchor는 우측/하단 정렬 시 창이 화면 밖으로 나가지 않도록 좌표를 보정합니다.
func reanchor(actual Rect, monitor Monitor, mode string, gap float64) Rect {
	corrected := actual

	if contains(mode, "right") {
		expectedRight := float64(monitor.X+monitor.Width) - gap
		actualRight := actual.X + actual.W
		if actualRight > expectedRight+5 {
			corrected.X = expectedRight - actual.W
		}
	}

	if contains(mode, "bottom") {
		expectedBottom := float64(monitor.Y+monitor.Height) - gap
		actualBottom := actual.Y + actual.H
		if actualBottom > expectedBottom+5 {
			corrected.Y = expectedBottom - actual.H
		}
	}

	return corrected
}

// moveToNextDisplay는 창을 다음 모니터로 이동합니다.
func moveToNextDisplay(pid int, current Rect, monitors []Monitor, currentMonitor Monitor) error {
	if len(monitors) <= 1 {
		return nil
	}

	idx := 0
	for i, m := range monitors {
		if m == currentMonitor {
			idx = i
			break
		}
	}
	next := monitors[(idx+1)%len(monitors)]

	// 저장된 원래 상태 기준으로 상대 좌표 유지
	windowStateMu.Lock()
	saved, ok := windowStateStore[pid]
	windowStateMu.Unlock()
	if !ok {
		saved = current
	}

	relX := saved.X - float64(currentMonitor.X)
	relY := saved.Y - float64(currentMonitor.Y)
	newW := min64(saved.W, float64(next.Width))
	newH := min64(saved.H, float64(next.Height))
	newX := clamp(float64(next.X)+relX, float64(next.X), float64(next.X+next.Width)-newW)
	newY := clamp(float64(next.Y)+relY, float64(next.Y), float64(next.Y+next.Height)-newH)

	SetWindowFrame(pid, Rect{newX, newY, newW, newH})
	ActivateApp(pid)
	utils.Log.Infof("다음 모니터로 이동: pid=%d", pid)
	return nil
}

func openSystemPrefs(url string) error {
	return exec.Command("open", url).Run()
}

// cycleMap은 스마트 사이클 순서를 정의합니다.
// left/right: 1/2 → 1/3 → 2/3 → 1/2 순환
// 상하/쿼터/maximize 등: 사이클 없음 (다음 모니터로 이동)
var cycleMap = map[string]string{
	"left_half": "left_1/3",
	"left_1/3":  "left_2/3",
	"left_2/3":  "left_half",
	"right_half": "right_1/3",
	"right_1/3":  "right_2/3",
	"right_2/3":  "right_half",
}

// nextCycleMode는 현재 모드의 다음 사이클 모드를 반환합니다.
// 사이클이 없으면 빈 문자열을 반환합니다.
func nextCycleMode(mode string) string {
	next, ok := cycleMap[mode]
	if !ok {
		return ""
	}
	// 자기 자신을 가리키면 사이클 없음
	if next == mode {
		return ""
	}
	return next
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func min64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
