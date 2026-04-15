// gen_icon: WR 텍스트 아이콘을 각 해상도로 생성합니다.
// - icon.iconset/ : .app 번들용 icns 소스
// - tray_icon.png : 트레이/독바 아이콘 (22px)
package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"
)

func main() {
	outDir := "icon.iconset"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		panic(err)
	}

	type iconSpec struct {
		filename string
		size     int
	}
	specs := []iconSpec{
		{"icon_16x16.png", 16},
		{"icon_16x16@2x.png", 32},
		{"icon_32x32.png", 32},
		{"icon_32x32@2x.png", 64},
		{"icon_128x128.png", 128},
		{"icon_128x128@2x.png", 256},
		{"icon_256x256.png", 256},
		{"icon_256x256@2x.png", 512},
		{"icon_512x512.png", 512},
		{"icon_512x512@2x.png", 1024},
	}

	for _, spec := range specs {
		img := renderIcon(spec.size)
		f, err := os.Create(filepath.Join(outDir, spec.filename))
		if err != nil {
			panic(err)
		}
		if err := png.Encode(f, img); err != nil {
			panic(err)
		}
		f.Close()
	}

	// 트레이/독바 아이콘 (22px)
	tray := renderIcon(22)
	f, err := os.Create("tray_icon.png")
	if err != nil {
		panic(err)
	}
	if err := png.Encode(f, tray); err != nil {
		panic(err)
	}
	f.Close()
}

func renderIcon(size int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(img, img.Bounds(), image.Transparent, image.Point{}, draw.Src)

	s := float64(size)
	radius := s * 0.22
	white := color.RGBA{255, 255, 255, 255}
	drawRoundedRect(img, 0, 0, size, size, radius, white)

	black := color.RGBA{0, 0, 0, 255}
	drawWR(img, size, black)

	return img
}

func drawRoundedRect(img *image.RGBA, x0, y0, x1, y1 int, r float64, c color.RGBA) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			if inRoundedRect(float64(x)+0.5, float64(y)+0.5, float64(x0), float64(y0), float64(x1), float64(y1), r) {
				img.SetRGBA(x, y, c)
			}
		}
	}
}

func inRoundedRect(px, py, x0, y0, x1, y1, r float64) bool {
	if px >= x0+r && px < x1-r {
		return py >= y0 && py < y1
	}
	if py >= y0+r && py < y1-r {
		return px >= x0 && px < x1
	}
	corners := [][2]float64{
		{x0 + r, y0 + r},
		{x1 - r, y0 + r},
		{x0 + r, y1 - r},
		{x1 - r, y1 - r},
	}
	for _, c := range corners {
		dx := px - c[0]
		dy := py - c[1]
		if math.Sqrt(dx*dx+dy*dy) < r {
			return true
		}
	}
	return false
}

// drawWR은 "WR" 두 글자를 아이콘 중앙에 그립니다.
// 각 글자는 세그먼트(선분) 집합으로 정의되고 크기에 맞게 스케일됩니다.
func drawWR(img *image.RGBA, size int, c color.RGBA) {
	s := float64(size)

	// 전체 텍스트 영역: 아이콘의 64% 너비, 45% 높이, 중앙 정렬
	textW := s * 0.62
	textH := s * 0.44
	textX := (s - textW) / 2
	textY := (s - textH) / 2

	strokeW := s * 0.085 // 획 두께

	// 글자 너비/간격 비율 (W:gap:R = 1:0.12:0.85)
	wWidth := textW * 0.52
	gap := textW * 0.06
	rWidth := textW - wWidth - gap

	drawLetterW(img, textX, textY, wWidth, textH, strokeW, c)
	drawLetterR(img, textX+wWidth+gap, textY, rWidth, textH, strokeW, c)
}

// 선분을 두께 있게 그립니다 (픽셀 단위)
func drawThickLine(img *image.RGBA, x0, y0, x1, y1, thickness float64, c color.RGBA) {
	// 선분 방향 벡터
	dx := x1 - x0
	dy := y1 - y0
	length := math.Sqrt(dx*dx + dy*dy)
	if length < 0.001 {
		return
	}
	// 법선 벡터 (두께 방향)
	nx := -dy / length * thickness / 2
	ny := dx / length * thickness / 2

	// 바운딩 박스
	minX := int(math.Min(math.Min(x0+nx, x0-nx), math.Min(x1+nx, x1-nx))) - 1
	maxX := int(math.Max(math.Max(x0+nx, x0-nx), math.Max(x1+nx, x1-nx))) + 2
	minY := int(math.Min(math.Min(y0+ny, y0-ny), math.Min(y1+ny, y1-ny))) - 1
	maxY := int(math.Max(math.Max(y0+ny, y0-ny), math.Max(y1+ny, y1-ny))) + 2

	bounds := img.Bounds()
	for py := minY; py <= maxY; py++ {
		for px := minX; px <= maxX; px++ {
			if px < bounds.Min.X || px >= bounds.Max.X || py < bounds.Min.Y || py >= bounds.Max.Y {
				continue
			}
			fpx := float64(px) + 0.5
			fpy := float64(py) + 0.5
			// 점에서 선분까지의 거리
			dist := distToSegment(fpx, fpy, x0, y0, x1, y1)
			if dist <= thickness/2+0.5 {
				alpha := math.Min(1.0, thickness/2+0.5-dist)
				existing := img.RGBAAt(px, py)
				a := uint8(float64(c.A) * alpha)
				if a > existing.A {
					img.SetRGBA(px, py, color.RGBA{c.R, c.G, c.B, a})
				}
			}
		}
	}
}

func distToSegment(px, py, x0, y0, x1, y1 float64) float64 {
	dx := x1 - x0
	dy := y1 - y0
	lenSq := dx*dx + dy*dy
	if lenSq < 0.0001 {
		return math.Sqrt((px-x0)*(px-x0) + (py-y0)*(py-y0))
	}
	t := ((px-x0)*dx + (py-y0)*dy) / lenSq
	t = math.Max(0, math.Min(1, t))
	projX := x0 + t*dx
	projY := y0 + t*dy
	return math.Sqrt((px-projX)*(px-projX) + (py-projY)*(py-projY))
}

// drawLetterW는 W 글자를 그립니다.
// W = 4개의 대각선 획 (∧∧ 형태를 뒤집은 모양)
func drawLetterW(img *image.RGBA, x, y, w, h, sw float64, c color.RGBA) {
	// 5개 꼭짓점: 좌상, 좌중하, 중앙상, 우중하, 우상
	// 실제 W: top-left → bottom-center-left → middle-top → bottom-center-right → top-right
	x0 := x
	x1 := x + w*0.25
	x2 := x + w*0.50
	x3 := x + w*0.75
	x4 := x + w
	yTop := y
	yBot := y + h
	yMid := y + h*0.55 // 중간 꼭짓점 높이

	drawThickLine(img, x0, yTop, x1, yBot, sw, c)
	drawThickLine(img, x1, yBot, x2, yMid, sw, c)
	drawThickLine(img, x2, yMid, x3, yBot, sw, c)
	drawThickLine(img, x3, yBot, x4, yTop, sw, c)
}

// drawLetterR는 R 글자를 그립니다.
// R = 수직 획 + 위쪽 반원 + 오른쪽 하단 대각선
func drawLetterR(img *image.RGBA, x, y, w, h, sw float64, c color.RGBA) {
	// 수직 획 (왼쪽)
	drawThickLine(img, x+sw/2, y, x+sw/2, y+h, sw, c)

	// 위쪽 범프: 반원 (수직 획 상단 ~ 중간)
	// 반원 중심 및 반지름
	bumpH := h * 0.48          // 범프가 차지하는 높이
	rOuter := bumpH * 0.52     // 외부 반지름
	rInner := rOuter - sw*0.95 // 내부 반지름
	cx := x + sw/2
	cy := y + bumpH/2

	for py := int(y); py <= int(y+bumpH+1); py++ {
		for px := int(x); px <= int(x+w+1); px++ {
			fpx := float64(px) + 0.5
			fpy := float64(py) + 0.5
			dx := fpx - cx
			dy := fpy - cy
			dist := math.Sqrt(dx*dx + dy*dy)
			// 오른쪽 반원만 (dx >= 0)
			if dx >= -sw*0.1 && dist >= rInner && dist <= rOuter {
				bounds := img.Bounds()
				if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
					img.SetRGBA(px, py, c)
				}
			}
		}
	}

	// 범프 연결: 수직 획에서 범프 끝점까지 수평선 (위/아래)
	drawThickLine(img, x+sw/2, y+sw*0.1, x+cx-x, y+sw*0.1, sw*0.0, c) // 위 연결 (두께 0 = 점)
	// 범프 하단 연결 수평선
	drawThickLine(img, x+sw/2, y+bumpH, x+sw/2+rInner*0.7, y+bumpH, sw*0.85, c)

	// 오른쪽 하단 대각선 (leg)
	legStartX := x + sw/2 + rInner*0.5
	legStartY := y + bumpH - sw*0.3
	legEndX := x + w
	legEndY := y + h
	drawThickLine(img, legStartX, legStartY, legEndX, legEndY, sw, c)
}
