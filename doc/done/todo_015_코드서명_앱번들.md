# TODO-015: 앱 번들 코드 서명 (Ad-hoc) 적용

## 우선순위
1순위 (버그 수정 — 손쉬운 사용 권한 불인식)

## 원인 분석
```
codesign -dv /Applications/WinResizer.app
→ Signature=adhoc
→ Info.plist=not bound   ← 문제
→ Sealed Resources=none  ← 문제
```

`go build`로 생성된 바이너리는 링커가 자동으로 adhoc 서명을 붙이지만,
`Info.plist`와 `Resources/`는 번들 서명에 포함되지 않는다.

macOS TCC(투명성·동의·제어)는 앱 번들이 올바르게 서명돼야 손쉬운 사용 권한을 부여한다.
`Info.plist=not bound` 상태에서는 목록에 추가해도 `AXIsProcessTrusted()`가 `false`를 반환한다.

## 목표
`build.sh`에서 앱 번들 빌드 후 `codesign`으로 전체 번들을 서명하여
손쉬운 사용 권한이 정상 동작하도록 한다.

## 구현 방법
Apple Developer 계정 없이도 ad-hoc 서명으로 동일 Mac에서 동작 가능.

```bash
# Info.plist를 번들에 바인딩하고 전체 번들 서명
codesign --force --deep --sign - \
  --entitlements entitlements.plist \
  WinResizer.app
```

`--sign -`: ad-hoc 서명 (개발자 인증서 불필요)
`--deep`: 번들 내 모든 바이너리 재귀 서명
`--force`: 기존 서명 덮어쓰기
`--entitlements`: Info.plist의 권한 항목을 서명에 포함

## 검증 방법
```bash
codesign -dv /Applications/WinResizer.app
# Info.plist=found  ← 확인
# Sealed Resources=v2 rules  ← 확인

# 권한 확인
/Applications/WinResizer.app/Contents/MacOS/WinResizer &
# 단축키 실행 → 창 조절 동작 확인
```

## 작업 결과

**상태**: 완료

**수정 파일**: `build-work/build.sh`

**변경 내용**:
DMG 생성 전에 `codesign` 단계 추가:

```bash
codesign --force --deep --sign - \
    --entitlements "${SCRIPT_DIR}/entitlements.plist" \
    "${BUNDLE_DIR}"
```

**검증 결과**:
```
codesign -dv dist/WinResizer.app
→ Identifier=com.winresizer.app
→ Info.plist entries=12       ← 번들에 바인딩됨
→ Sealed Resources version=2  ← 리소스 서명됨
→ Signature=adhoc
```
