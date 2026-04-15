package core

import (
	"testing"
)

// --- IsSimilar 테스트 ---

func TestIsSimilar_ExactMatch(t *testing.T) {
	a := Rect{X: 0, Y: 0, W: 720, H: 900}
	b := Rect{X: 0, Y: 0, W: 720, H: 900}
	if !IsSimilar(a, b, 5.0) {
		t.Error("정확히 일치하는 경우 IsSimilar가 false를 반환")
	}
}

func TestIsSimilar_WithinTolerance(t *testing.T) {
	a := Rect{X: 2, Y: 1, W: 718, H: 898}
	b := Rect{X: 0, Y: 0, W: 720, H: 900}
	if !IsSimilar(a, b, 5.0) {
		t.Error("허용 오차(5px) 내 차이인데 IsSimilar가 false를 반환")
	}
}

func TestIsSimilar_OutOfTolerance(t *testing.T) {
	a := Rect{X: 10, Y: 0, W: 720, H: 900}
	b := Rect{X: 0, Y: 0, W: 720, H: 900}
	if IsSimilar(a, b, 5.0) {
		t.Error("허용 오차(5px) 초과인데 IsSimilar가 true를 반환")
	}
}

// --- reanchor 테스트 ---

func TestReanchor_RightOverflow(t *testing.T) {
	// 우측 정렬 시 앱 최소 크기로 창이 화면 밖으로 나간 경우 x 좌표 보정
	monitor := Monitor{X: 0, Y: 0, Width: 1440, Height: 900}
	actual := Rect{X: 900, Y: 0, W: 600, H: 900} // 900+600=1500 > 1440 초과
	corrected := reanchor(actual, monitor, "right_half", 0)

	expectedX := float64(1440) - 600.0 // 840
	if abs64(corrected.X-expectedX) > 1.0 {
		t.Errorf("reanchor 우측 보정 실패: got X=%.1f, want X=%.1f", corrected.X, expectedX)
	}
}

func TestReanchor_RightNoOverflow(t *testing.T) {
	// 화면 안에 있으면 보정 없음
	monitor := Monitor{X: 0, Y: 0, Width: 1440, Height: 900}
	actual := Rect{X: 720, Y: 0, W: 720, H: 900}
	corrected := reanchor(actual, monitor, "right_half", 0)

	if corrected != actual {
		t.Errorf("reanchor 불필요한 보정 발생: got %+v, want %+v", corrected, actual)
	}
}

func TestReanchor_BottomOverflow(t *testing.T) {
	// 하단 정렬 시 화면 밖 보정
	monitor := Monitor{X: 0, Y: 0, Width: 1440, Height: 900}
	actual := Rect{X: 0, Y: 500, W: 1440, H: 500} // 500+500=1000 > 900 초과
	corrected := reanchor(actual, monitor, "bottom_half", 0)

	expectedY := float64(900) - 500.0 // 400
	if abs64(corrected.Y-expectedY) > 1.0 {
		t.Errorf("reanchor 하단 보정 실패: got Y=%.1f, want Y=%.1f", corrected.Y, expectedY)
	}
}

func TestReanchor_WithGap(t *testing.T) {
	// gap 적용 시 보정
	monitor := Monitor{X: 0, Y: 0, Width: 1440, Height: 900}
	gap := 4.0
	actual := Rect{X: 736, Y: 0, W: 710, H: 900} // 736+710=1446 > 1440-4=1436 초과
	corrected := reanchor(actual, monitor, "right_half", gap)

	expectedX := float64(1440) - gap - 710.0 // 726
	if abs64(corrected.X-expectedX) > 1.0 {
		t.Errorf("reanchor gap 보정 실패: got X=%.1f, want X=%.1f", corrected.X, expectedX)
	}
}

func TestReanchor_LeftMode_NoChange(t *testing.T) {
	// 좌측 정렬 모드는 reanchor 대상 아님
	monitor := Monitor{X: 0, Y: 0, Width: 1440, Height: 900}
	actual := Rect{X: 0, Y: 0, W: 720, H: 900}
	corrected := reanchor(actual, monitor, "left_half", 0)

	if corrected != actual {
		t.Errorf("left_half reanchor 불필요한 보정 발생: got %+v", corrected)
	}
}

// --- nextCycleMode 테스트 ---

func TestNextCycleMode_LeftCycle(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"left_half", "left_1/3"},
		{"left_1/3", "left_2/3"},
		{"left_2/3", "left_half"},
	}
	for _, tt := range tests {
		got := nextCycleMode(tt.input)
		if got != tt.want {
			t.Errorf("nextCycleMode(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNextCycleMode_RightCycle(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"right_half", "right_1/3"},
		{"right_1/3", "right_2/3"},
		{"right_2/3", "right_half"},
	}
	for _, tt := range tests {
		got := nextCycleMode(tt.input)
		if got != tt.want {
			t.Errorf("nextCycleMode(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNextCycleMode_NoCycle(t *testing.T) {
	// 사이클 없는 모드는 빈 문자열 반환
	modes := []string{"top_half", "bottom_half", "maximize", "restore", "top_left_1/4", "left_custom:60"}
	for _, mode := range modes {
		got := nextCycleMode(mode)
		if got != "" {
			t.Errorf("nextCycleMode(%q) = %q, want \"\"", mode, got)
		}
	}
}

// --- isAlreadyAligned 테스트 ---

func TestIsAlreadyAligned_ExactMatch(t *testing.T) {
	monitor := Monitor{X: 0, Y: 0, Width: 1440, Height: 900}
	frame := Rect{X: 0, Y: 0, W: 720, H: 900}
	target := Rect{X: 0, Y: 0, W: 720, H: 900}

	if !isAlreadyAligned(frame, target, "left_half", monitor) {
		t.Error("정확히 일치하는 경우 isAlreadyAligned가 false를 반환")
	}
}

func TestIsAlreadyAligned_WithinTolerance(t *testing.T) {
	// 허용 오차 내 위치/크기 차이
	monitor := Monitor{X: 0, Y: 0, Width: 1440, Height: 900}
	frame := Rect{X: 2, Y: 1, W: 718, H: 899}
	target := Rect{X: 0, Y: 0, W: 720, H: 900}

	if !isAlreadyAligned(frame, target, "left_half", monitor) {
		t.Error("허용 오차 내 차이인데 isAlreadyAligned가 false를 반환")
	}
}

func TestIsAlreadyAligned_NotAligned(t *testing.T) {
	// 완전히 다른 위치
	monitor := Monitor{X: 0, Y: 0, Width: 1440, Height: 900}
	frame := Rect{X: 720, Y: 0, W: 720, H: 900} // 우측 절반
	target := Rect{X: 0, Y: 0, W: 720, H: 900}  // 좌측 절반 목표

	if isAlreadyAligned(frame, target, "left_half", monitor) {
		t.Error("다른 위치인데 isAlreadyAligned가 true를 반환")
	}
}

func TestReanchor_CustomRightMode(t *testing.T) {
	// right_custom:35 — 앱 최소 크기로 확장된 경우 보정
	monitor := Monitor{X: 0, Y: 0, Width: 1440, Height: 900}
	// 목표: 35% = 504px, x = 1440-504 = 936
	// 실제: 크롬 최소 폭 600px → x=936, w=600 → 936+600=1536 > 1440 초과
	actual := Rect{X: 936, Y: 0, W: 600, H: 900}
	corrected := reanchor(actual, monitor, "right_custom:35", 0)

	expectedX := float64(1440) - 600.0 // 840
	if abs64(corrected.X-expectedX) > 1.0 {
		t.Errorf("right_custom reanchor 보정 실패: got X=%.1f, want X=%.1f", corrected.X, expectedX)
	}
}

func TestReanchor_BottomRightQuarter(t *testing.T) {
	// bottom_right_1/4 — 우측+하단 모두 보정
	monitor := Monitor{X: 0, Y: 0, Width: 1440, Height: 900}
	actual := Rect{X: 800, Y: 550, W: 700, H: 400} // 우측 1540>1440, 하단 950>900 초과
	corrected := reanchor(actual, monitor, "bottom_right_1/4", 0)

	expectedX := float64(1440) - 700.0 // 740
	expectedY := float64(900) - 400.0  // 500
	if abs64(corrected.X-expectedX) > 1.0 {
		t.Errorf("bottom_right 우측 보정 실패: got X=%.1f, want X=%.1f", corrected.X, expectedX)
	}
	if abs64(corrected.Y-expectedY) > 1.0 {
		t.Errorf("bottom_right 하단 보정 실패: got Y=%.1f, want Y=%.1f", corrected.Y, expectedY)
	}
}

// --- executeSizeRecycle 계산 검증 ---
// executeSizeRecycle은 CGo 호출(SetWindowFrame 등)을 포함하므로 직접 호출 불가.
// 핵심 계산 로직을 순수 함수로 추출하여 검증합니다.

func calcSizeRecycle(current Rect, monitor Monitor, mode string, gap float64) Rect {
	screenW := float64(monitor.Width)
	screenH := float64(monitor.Height)
	step := screenW * 0.10

	minW := screenW*sizeRecycleMinRatio - gap*2
	maxW := screenW*sizeRecycleMaxRatio - gap*2

	newW := current.W
	switch mode {
	case "size_grow_left", "size_grow_right":
		newW = current.W + step
		if newW > maxW+step*0.5 {
			newW = minW
		}
	case "size_shrink_left", "size_shrink_right":
		newW = current.W - step
		if newW < minW-step*0.5 {
			newW = maxW
		}
	}

	newH := screenH - gap*2
	monAbsX := float64(monitor.X)
	monAbsY := float64(monitor.Y)

	var newX float64
	switch mode {
	case "size_grow_left", "size_shrink_left":
		newX = monAbsX + gap
	case "size_grow_right", "size_shrink_right":
		newX = monAbsX + screenW - gap - newW
	default:
		newX = current.X
	}

	return Rect{X: newX, Y: monAbsY + gap, W: newW, H: newH}
}

func TestSizeRecycle_GrowLeft(t *testing.T) {
	// 중간 값(50%)에서 확장 → 60%
	monitor := Monitor{X: 0, Y: 0, Width: 2000, Height: 2000}
	current := Rect{X: 0, Y: 0, W: 1000, H: 2000} // 50%
	got := calcSizeRecycle(current, monitor, "size_grow_left", 0)

	if abs64(got.W-1200) > 1 {
		t.Errorf("size_grow_left: W=%.0f, want 1200", got.W)
	}
	if abs64(got.X-0) > 1 {
		t.Errorf("size_grow_left: X=%.0f, want 0 (좌측 엣지 고정)", got.X)
	}
}

func TestSizeRecycle_ShrinkLeft(t *testing.T) {
	// 중간 값(50%)에서 축소 → 40%
	monitor := Monitor{X: 0, Y: 0, Width: 2000, Height: 2000}
	current := Rect{X: 0, Y: 0, W: 1000, H: 2000}
	got := calcSizeRecycle(current, monitor, "size_shrink_left", 0)

	if abs64(got.W-800) > 1 {
		t.Errorf("size_shrink_left: W=%.0f, want 800", got.W)
	}
	if abs64(got.X-0) > 1 {
		t.Errorf("size_shrink_left: X=%.0f, want 0 (좌측 엣지 고정)", got.X)
	}
}

func TestSizeRecycle_GrowRight(t *testing.T) {
	// 중간 값(50%)에서 확장 → 60%, 우측 엣지 고정
	monitor := Monitor{X: 0, Y: 0, Width: 2000, Height: 2000}
	current := Rect{X: 1000, Y: 0, W: 1000, H: 2000}
	got := calcSizeRecycle(current, monitor, "size_grow_right", 0)

	// 새 폭 = 1000+200 = 1200, 우측 엣지 고정 → x = 2000-1200 = 800
	if abs64(got.W-1200) > 1 {
		t.Errorf("size_grow_right: W=%.0f, want 1200", got.W)
	}
	if abs64(got.X-800) > 1 {
		t.Errorf("size_grow_right: X=%.0f, want 800 (우측 엣지 고정)", got.X)
	}
}

func TestSizeRecycle_ShrinkRight(t *testing.T) {
	// 중간 값(50%)에서 축소 → 40%, 우측 엣지 고정
	monitor := Monitor{X: 0, Y: 0, Width: 2000, Height: 2000}
	current := Rect{X: 1000, Y: 0, W: 1000, H: 2000}
	got := calcSizeRecycle(current, monitor, "size_shrink_right", 0)

	// 새 폭 = 1000-200 = 800, 우측 엣지 고정 → x = 2000-800 = 1200
	if abs64(got.W-800) > 1 {
		t.Errorf("size_shrink_right: W=%.0f, want 800", got.W)
	}
	if abs64(got.X-1200) > 1 {
		t.Errorf("size_shrink_right: X=%.0f, want 1200 (우측 엣지 고정)", got.X)
	}
}

func TestSizeRecycle_GrowWrapsToMin(t *testing.T) {
	// 최댓값(90%)에서 확장 → 최솟값(20%)으로 순환
	monitor := Monitor{X: 0, Y: 0, Width: 2000, Height: 2000}
	maxW := 2000.0 * sizeRecycleMaxRatio // 1800
	minW := 2000.0 * sizeRecycleMinRatio // 400
	current := Rect{X: 0, Y: 0, W: maxW, H: 2000}
	got := calcSizeRecycle(current, monitor, "size_grow_left", 0)

	if abs64(got.W-minW) > 1 {
		t.Errorf("size_grow_left 최대→최소 순환 실패: W=%.0f, want %.0f", got.W, minW)
	}
}

func TestSizeRecycle_ShrinkWrapsToMax(t *testing.T) {
	// 최솟값(20%)에서 축소 → 최댓값(90%)으로 순환
	monitor := Monitor{X: 0, Y: 0, Width: 2000, Height: 2000}
	minW := 2000.0 * sizeRecycleMinRatio // 400
	maxW := 2000.0 * sizeRecycleMaxRatio // 1800
	current := Rect{X: 0, Y: 0, W: minW, H: 2000}
	got := calcSizeRecycle(current, monitor, "size_shrink_left", 0)

	if abs64(got.W-maxW) > 1 {
		t.Errorf("size_shrink_left 최소→최대 순환 실패: W=%.0f, want %.0f", got.W, maxW)
	}
}

func TestIsAlreadyAligned_RightEdge_MinSizeExpanded(t *testing.T) {
	// 우측 정렬 — 앱 최소 크기로 창이 목표보다 더 넓어진 경우에도 "정렬됨"으로 판단
	monitor := Monitor{X: 0, Y: 0, Width: 1440, Height: 900}
	// 목표: right_1/3 = x=960, w=480
	// 실제: 앱 최소 크기 600px → x 보정 후 x=840, w=600
	frame := Rect{X: 840, Y: 0, W: 600, H: 900}
	target := Rect{X: 960, Y: 0, W: 480, H: 900}

	if !isAlreadyAligned(frame, target, "right_1/3", monitor) {
		t.Error("우측 정렬 최소 크기 확장 케이스에서 isAlreadyAligned가 false를 반환")
	}
}
