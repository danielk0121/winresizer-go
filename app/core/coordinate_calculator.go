package core

import (
	"fmt"
	"strings"
	"strconv"
)

// Rect는 창의 위치와 크기를 나타냅니다.
type Rect struct {
	X, Y, W, H float64
}

// CalculateWindowPosition은 스크린 크기와 모드에 따라 창의 목표 좌표를 계산합니다.
// gap은 창과 화면 경계 사이의 여백(px)입니다.
func CalculateWindowPosition(screenW, screenH float64, mode string, gap float64) (Rect, error) {
	g := gap

	// 1/2 분할
	switch mode {
	case "left_half":
		return Rect{g, g, screenW/2 - g*1.5, screenH - g*2}, nil
	case "right_half":
		return Rect{screenW/2 + g*0.5, g, screenW/2 - g*1.5, screenH - g*2}, nil
	case "top_half":
		return Rect{g, g, screenW - g*2, screenH/2 - g*1.5}, nil
	case "bottom_half":
		return Rect{g, screenH/2 + g*0.5, screenW - g*2, screenH/2 - g*1.5}, nil

	// 1/4 분할
	case "top_left_1/4":
		return Rect{g, g, screenW/2 - g*1.5, screenH/2 - g*1.5}, nil
	case "top_right_1/4":
		return Rect{screenW/2 + g*0.5, g, screenW/2 - g*1.5, screenH/2 - g*1.5}, nil
	case "bottom_left_1/4":
		return Rect{g, screenH/2 + g*0.5, screenW/2 - g*1.5, screenH/2 - g*1.5}, nil
	case "bottom_right_1/4":
		return Rect{screenW/2 + g*0.5, screenH/2 + g*0.5, screenW/2 - g*1.5, screenH/2 - g*1.5}, nil

	// 1/3 분할
	case "left_1/3":
		u := screenW / 3
		return Rect{g, g, u - g*1.5, screenH - g*2}, nil
	case "center_1/3":
		u := screenW / 3
		return Rect{u + g*0.5, g, u - g*1.0, screenH - g*2}, nil
	case "right_1/3":
		u := screenW / 3
		return Rect{2*u + g*0.5, g, u - g*1.5, screenH - g*2}, nil

	// 2/3 분할
	case "left_2/3":
		u := screenW / 3
		return Rect{g, g, 2*u - g*1.0, screenH - g*2}, nil
	case "right_2/3":
		u := screenW / 3
		return Rect{u + g*0.5, g, 2*u - g*1.0, screenH - g*2}, nil

	// 최대화
	case "maximize":
		return Rect{g, g, screenW - g*2, screenH - g*2}, nil
	}

	// 커스텀 비율: left_custom:75, right_custom:60 등
	if strings.Contains(mode, "_custom:") {
		parts := strings.SplitN(mode, "_custom:", 2)
		if len(parts) == 2 {
			direction := parts[0]
			pct, err := strconv.Atoi(parts[1])
			if err != nil || pct < 1 || pct > 100 {
				return Rect{}, fmt.Errorf("잘못된 커스텀 비율: %s", mode)
			}
			ratio := float64(pct) / 100.0
			switch direction {
			case "left":
				return Rect{g, g, screenW*ratio - g*2, screenH - g*2}, nil
			case "right":
				startX := screenW * (1 - ratio)
				return Rect{startX + g, g, screenW*ratio - g*2, screenH - g*2}, nil
			case "top":
				return Rect{g, g, screenW - g*2, screenH*ratio - g*2}, nil
			case "bottom":
				startY := screenH * (1 - ratio)
				return Rect{g, startY + g, screenW - g*2, screenH*ratio - g*2}, nil
			}
		}
	}

	// 알 수 없는 모드: 최대화로 폴백
	return Rect{g, g, screenW - g*2, screenH - g*2}, fmt.Errorf("알 수 없는 모드: %s", mode)
}

// IsSimilar는 두 Rect가 tolerance 범위 내에서 같은지 비교합니다.
func IsSimilar(a, b Rect, tolerance float64) bool {
	return abs64(a.X-b.X) <= tolerance &&
		abs64(a.Y-b.Y) <= tolerance &&
		abs64(a.W-b.W) <= tolerance &&
		abs64(a.H-b.H) <= tolerance
}

func abs64(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
