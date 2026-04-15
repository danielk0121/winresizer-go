package core

import (
	"testing"
)

func TestCalculateWindowPosition_HalfSplits(t *testing.T) {
	sw, sh := 1440.0, 900.0
	gap := 0.0

	tests := []struct {
		mode string
		want Rect
	}{
		{"left_half",   Rect{0, 0, 720, 900}},
		{"right_half",  Rect{720, 0, 720, 900}},
		{"top_half",    Rect{0, 0, 1440, 450}},
		{"bottom_half", Rect{0, 450, 1440, 450}},
	}

	for _, tt := range tests {
		got, err := CalculateWindowPosition(sw, sh, tt.mode, gap)
		if err != nil {
			t.Errorf("[%s] 예상치 못한 에러: %v", tt.mode, err)
			continue
		}
		if !IsSimilar(got, tt.want, 1.0) {
			t.Errorf("[%s] got %+v, want %+v", tt.mode, got, tt.want)
		}
	}
}

func TestCalculateWindowPosition_ThirdSplits(t *testing.T) {
	sw, sh := 1440.0, 900.0
	gap := 0.0
	u := sw / 3 // 480

	tests := []struct {
		mode string
		want Rect
	}{
		{"left_1/3",   Rect{0, 0, u, sh}},
		{"center_1/3", Rect{u, 0, u, sh}},
		{"right_1/3",  Rect{2 * u, 0, u, sh}},
		{"left_2/3",   Rect{0, 0, 2 * u, sh}},
		{"right_2/3",  Rect{u, 0, 2 * u, sh}},
	}

	for _, tt := range tests {
		got, err := CalculateWindowPosition(sw, sh, tt.mode, gap)
		if err != nil {
			t.Errorf("[%s] 예상치 못한 에러: %v", tt.mode, err)
			continue
		}
		if !IsSimilar(got, tt.want, 1.0) {
			t.Errorf("[%s] got %+v, want %+v", tt.mode, got, tt.want)
		}
	}
}

func TestCalculateWindowPosition_QuarterSplits(t *testing.T) {
	sw, sh := 1440.0, 900.0
	gap := 0.0

	tests := []struct {
		mode string
		want Rect
	}{
		{"top_left_1/4",     Rect{0, 0, sw / 2, sh / 2}},
		{"top_right_1/4",    Rect{sw / 2, 0, sw / 2, sh / 2}},
		{"bottom_left_1/4",  Rect{0, sh / 2, sw / 2, sh / 2}},
		{"bottom_right_1/4", Rect{sw / 2, sh / 2, sw / 2, sh / 2}},
	}

	for _, tt := range tests {
		got, err := CalculateWindowPosition(sw, sh, tt.mode, gap)
		if err != nil {
			t.Errorf("[%s] 예상치 못한 에러: %v", tt.mode, err)
			continue
		}
		if !IsSimilar(got, tt.want, 1.0) {
			t.Errorf("[%s] got %+v, want %+v", tt.mode, got, tt.want)
		}
	}
}

func TestCalculateWindowPosition_Custom(t *testing.T) {
	sw, sh := 1440.0, 900.0
	gap := 0.0

	tests := []struct {
		mode string
		wantX float64
		wantW float64
	}{
		{"left_custom:60",  0,          sw * 0.6},
		{"right_custom:60", sw * 0.4,   sw * 0.6},
	}

	for _, tt := range tests {
		got, err := CalculateWindowPosition(sw, sh, tt.mode, gap)
		if err != nil {
			t.Errorf("[%s] 예상치 못한 에러: %v", tt.mode, err)
			continue
		}
		if abs64(got.X-tt.wantX) > 1.0 || abs64(got.W-tt.wantW) > 1.0 {
			t.Errorf("[%s] got X=%.1f W=%.1f, want X=%.1f W=%.1f", tt.mode, got.X, got.W, tt.wantX, tt.wantW)
		}
	}
}

func TestCalculateWindowPosition_WithGap(t *testing.T) {
	sw, sh := 1440.0, 900.0
	gap := 4.0

	got, err := CalculateWindowPosition(sw, sh, "left_half", gap)
	if err != nil {
		t.Fatalf("예상치 못한 에러: %v", err)
	}
	// gap이 적용되면 x > 0, w < sw/2
	if got.X < gap-0.1 {
		t.Errorf("gap 적용 후 X가 너무 작음: %v", got.X)
	}
	if got.W >= sw/2 {
		t.Errorf("gap 적용 후 W가 너무 큼: %v", got.W)
	}
}

// 화면 밖으로 나가지 않는지 검증
func TestCalculateWindowPosition_NoBoundaryOverflow(t *testing.T) {
	sw, sh := 1440.0, 900.0
	gap := 0.0
	modes := []string{
		"left_half", "right_half", "top_half", "bottom_half",
		"left_1/3", "center_1/3", "right_1/3",
		"left_2/3", "right_2/3",
		"top_left_1/4", "top_right_1/4", "bottom_left_1/4", "bottom_right_1/4",
		"maximize",
	}
	for _, mode := range modes {
		got, _ := CalculateWindowPosition(sw, sh, mode, gap)
		if got.X < 0 || got.Y < 0 {
			t.Errorf("[%s] 음수 좌표: %+v", mode, got)
		}
		if got.X+got.W > sw+1 {
			t.Errorf("[%s] 화면 우측 초과: %+v", mode, got)
		}
		if got.Y+got.H > sh+1 {
			t.Errorf("[%s] 화면 하단 초과: %+v", mode, got)
		}
	}
}
