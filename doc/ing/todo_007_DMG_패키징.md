# TODO-007: DMG 패키징

## 우선순위
3순위 (배포 준비)

## 목표
사용자 배포용 `.dmg` 설치 이미지 생성.

## 구현 방법
`hdiutil`을 이용한 DMG 생성:

```bash
# 1. 임시 디렉토리 구성
mkdir -p dist/dmg
cp -r WinResizer.app dist/dmg/
ln -s /Applications dist/dmg/Applications  # 드래그 설치용 심볼릭 링크

# 2. DMG 생성
hdiutil create -volname "WinResizer" \
  -srcfolder dist/dmg \
  -ov -format UDZO \
  dist/WinResizer.dmg
```

### 선택 사항
- DMG 배경 이미지 적용 (사용자 안내용)
- 아이콘 배치 위치 설정 (`AppleScript`로 Finder 뷰 커스텀)

## 의존성
- TODO-006 (앱 번들 빌드) 완료 후 진행

## 검증 방법
- DMG 마운트 후 Applications 폴더로 드래그 설치
- 설치된 앱 실행 → 정상 동작 확인

## 작업 결과

**상태**: 완료

**수정 파일**: `app/build/build.sh`

**변경 내용**:
`build.sh`의 DMG 생성 단계에 `/Applications` 심볼릭 링크 추가. 드래그 설치 UX 개선.

```bash
mkdir -p dist/dmg
cp -r "${BUNDLE_DIR}" dist/dmg/
ln -sf /Applications dist/dmg/Applications   # 드래그 설치용

hdiutil create -volname "${APP_NAME}" \
    -srcfolder dist/dmg \
    -ov -format UDZO \
    "dist/${APP_NAME}.dmg"

rm -rf dist/dmg
```

빌드 완료 후 `dist/WinResizer.dmg` (~10MB) 생성 확인.

**커밋**: `2a0367d`
