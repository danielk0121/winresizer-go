package core

import (
	"testing"
)

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
