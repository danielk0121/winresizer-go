# TODO-016: 아이콘 디자인 개선 (트레이 + 독바 통일)

## 우선순위
2순위 (UI 일관성)

## 문제 현황
1. **G 글자가 너무 크고 획이 두꺼움** — 아이콘 전체를 꽉 채워 촌스러움
2. **트레이 아이콘(독바)과 앱 아이콘(Finder/독) 디자인 불일치**
   - 트레이(`tray_icon.png`): 22×22 소형, 흰 배경 + G
   - 앱 번들(`icon.icns`): 별도 생성, 스타일 다름
3. **독바 아이콘 크기** — 카카오톡/텔레그램 대비 너무 큼

## 목표
- 흰색 둥근 사각형 배경 + 검은색 **"WR"** 텍스트
- 트레이 아이콘, 웹서버 파비콘, 독바 아이콘, 앱 번들 아이콘 모두 동일 디자인으로 교체
- G 글자 방식 폐기 → WR 텍스트 벡터 렌더링으로 전환

## 구현 방법
`build-work/gen_icon/main.go` 수정:
- `drawGShape()` 제거 → `drawWRText()` 구현
  - W, R 두 글자를 픽셀 단위 벡터로 렌더링
  - 전체 크기의 약 60% 영역 안에 두 글자 배치
  - 획 두께: 전체의 약 8~9%
- 라운드 반경: size×0.22
- `tray_icon.png` 22px도 동일 함수로 생성
- `iconutil`로 `icon.icns` 재생성
- `app/ui/tray_icon.png` 교체
- `app/ui/static/favicon.png` 동일 파일로 교체
- `build.sh` 빌드 전 gen_icon 자동 실행 추가

## 검증 방법
- `build-work/gen_icon/icon.iconset/icon_256x256.png` 미리보기로 시각 확인
- 카카오톡/텔레그램과 독바에서 나란히 비교
- 빌드 후 Finder 아이콘, 독 아이콘, 트레이 아이콘 모두 동일한지 확인

## 작업 결과

**상태**: 완료

**수정/생성 파일**:
- `build-work/gen_icon/wr_icon.svg` — WR 텍스트 SVG 아이콘 소스
- `build-work/gen_icon/gen_icon.js` — Node.js sharp 라이브러리로 SVG → PNG 변환
- `app/ui/icon.icns` — 앱 번들 아이콘 (iconutil 생성)
- `app/ui/tray_icon.png` — 트레이 아이콘 (22px, 88px 렌더 후 다운스케일)
- `app/ui/static/favicon.png` — 웹서버 파비콘

**변경 내용**:
- G 글자 방식 폐기 → 흰색 둥근 사각형 배경 + 검은색 "WR" 텍스트로 통일
- SVG 렌더링: `rsvg-convert` 대신 Node.js `sharp` 라이브러리 사용 (macOS 12 호환성 문제로)
- 트레이 아이콘(22px)은 88px로 렌더링 후 다운스케일하여 글자 깨짐 방지
- 트레이 아이콘, 파비콘, 앱 번들 아이콘 모두 동일 SVG 소스에서 생성
