# TODO-014: 앱 번들 아이콘(.icns) G 글자 아이콘으로 교체

## 우선순위
2순위 (UI 일관성)

## 원인 분석
- `.app` 번들 아이콘(Finder에 표시되는 아이콘)은 `app/ui/icon.icns`에서 온다
- `build.sh`에서 `icon.icns` → `Contents/Resources/icon.icns`로 복사
- 현재 `icon.icns`(34KB, 512×512)는 기존 사각형 2개 겹친 아이콘
- `tray_icon.png`(376B, 22×22, G 글자 흰색 배경)는 런타임 트레이 아이콘용이며 번들 아이콘과 별개
- 즉, 트레이 아이콘은 G 글자로 바뀌었지만 Finder용 `.app` 아이콘은 아직 갱신 안 됨

## 목표
`app/ui/icon.icns`를 G 글자 디자인(흰색 배경 + 검정 G)의 icns 파일로 교체하여 Finder와 Dock에서도 동일한 아이콘이 표시되도록 한다.

## 구현 방법

### 1. tray_icon.png → 고해상도 PNG 생성
macOS `.icns`는 여러 해상도(16, 32, 64, 128, 256, 512, 1024px)를 포함해야 한다.
`tray_icon.png`는 22×22로 작으므로, 동일한 디자인으로 큰 해상도 PNG를 생성해야 한다.

Go 코드 또는 `sips` / `iconutil` 로 생성:

```bash
# 1. iconset 폴더 구성
mkdir -p build-work/icon.iconset
# 각 해상도별 PNG 생성 (Go 스크립트 또는 sips로 리사이즈)
# icon_16x16.png, icon_32x32.png, ..., icon_512x512@2x.png

# 2. icns 변환
iconutil -c icns build-work/icon.iconset -o app/ui/icon.icns
```

### 방법 A: Go 스크립트로 PNG 생성 후 iconutil
- G 글자 아이콘을 생성하는 Go 스크립트 작성 (`build-work/gen_icon.go`)
- 16, 32, 64, 128, 256, 512, 1024px 각각 렌더링
- `iconutil`로 `.icns` 생성

### 방법 B: tray_icon.png를 고해상도로 리사이즈
- `sips -z {size} {size} tray_icon.png` 로 각 해상도 생성 (품질 저하 가능)

## 권장 방법
방법 A (Go 스크립트) — `tray_icon.png` 생성 코드(`tray_icon_gen.go` 등)가 이미 존재한다면 해상도 파라미터만 바꿔 재사용 가능.

## 검증 방법
- `build.sh` 실행 후 Finder에서 `dist/WinResizer.app` 아이콘 확인
- G 글자 흰색 배경 아이콘 표시 확인

## 작업 결과

**상태**: 미완료
