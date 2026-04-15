# TODO-011: 트레이 아이콘 확인 + 웹 파비콘 통일

## 우선순위
2순위 (UI 일관성)

## 목표
1. 빌드 결과물(`dist/WinResizer.app`)에 변경된 트레이 아이콘(흰색 배경 + G 글자)이 적용됐는지 확인
2. 웹 파비콘(`app/ui/static/favicon.png`)을 트레이 아이콘(`app/ui/tray_icon.png`)과 동일한 이미지로 교체

## 현재 상태
- `app/ui/tray_icon.png`: 22×22 RGBA PNG, 흰색 둥근 사각형 배경(radius=3) + 검정 G 글자
- `app/ui/static/favicon.png`: 다른 이미지 (트레이 아이콘과 불일치)
- `app/ui/templates/index.html` line 7: `<link rel="icon" type="image/png" href="/static/favicon.png">`

## 구현 방법
1. 빌드 결과물 트레이 아이콘 확인 (시각적 검증 또는 파일 비교)
2. `app/ui/static/favicon.png`를 `app/ui/tray_icon.png`와 동일한 파일로 교체

```bash
cp app/ui/tray_icon.png app/ui/static/favicon.png
```

## 검증 방법
- 웹앱 접속 후 브라우저 탭 아이콘 확인 → 트레이 아이콘과 동일한 G 마크 표시
- `md5 app/ui/tray_icon.png app/ui/static/favicon.png` → 해시 일치 확인

## 작업 결과

**상태**: 완료

**수정 파일**: `app/ui/static/favicon.png`

**변경 내용**:
`tray_icon.png`를 `favicon.png`로 복사하여 동일한 이미지로 통일.

```bash
cp app/ui/tray_icon.png app/ui/static/favicon.png
```

**검증 결과**:
```
MD5 (tray_icon.png)      = a9694723ae05cce69501e5a24e16ea76
MD5 (favicon.png)        = a9694723ae05cce69501e5a24e16ea76
```
