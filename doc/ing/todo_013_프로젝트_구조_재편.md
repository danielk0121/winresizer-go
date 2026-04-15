# TODO-013: 프로젝트 구조 재편 (/app/ + /build-work/ 분리)

## 우선순위
2순위 (코드베이스 정리)

## 목표
현재 `/app/` 내에 혼재된 소스코드와 빌드 스크립트를 분리한다.

- `/app/`: Go 소스코드, config, HTML/JS/CSS 리소스만 존재
- `/build-work/`: 빌드 스크립트(`build.sh`, `Info.plist`, `entitlements.plist`)와 `dist/` 출력 폴더

## 현재 구조
```
/app/
  build/          ← build.sh, Info.plist, entitlements.plist
  config/         ← default-config.json
  core/           ← Go 소스
  server/         ← Go 소스
  ui/             ← tray_icon.png, static/, templates/
  main.go
```

## 목표 구조
```
/app/             ← Go 소스 + 리소스만
  config/
  core/
  server/
  ui/
  main.go

/build-work/      ← 빌드 관련만
  build.sh
  Info.plist
  entitlements.plist
  dist/           ← 빌드 출력 (.gitignore)
```

## 구현 방법
1. `build-work/` 디렉토리 생성
2. `app/build/` 내 파일들을 `build-work/`로 이동
3. `build.sh` 내 경로 수정:
   - `go build` 호출 시 소스 경로를 `../app` 또는 `../app/...`로 업데이트
   - `dist/` 경로를 `build-work/dist/`로 유지
4. `.gitignore`의 `dist/` 경로 확인 및 업데이트
5. `app/build/` 빈 디렉토리 제거

## 검증 방법
```bash
cd build-work
bash build.sh
# dist/WinResizer.dmg 생성 확인
lipo -info dist/WinResizer.app/Contents/MacOS/WinResizer
# arm64 x86_64 출력 확인
```

## 작업 결과

**상태**: 완료

**수정/이동 파일**:
- `app/build/build.sh` → `build-work/build.sh`
- `app/build/Info.plist` → `build-work/Info.plist`
- `app/build/entitlements.plist` → `build-work/entitlements.plist`
- `app/build/` 디렉토리 제거
- `.gitignore`: `build-work/dist/` 추가

**변경 내용**:
`build.sh`를 `build-work/`에서 실행 가능하도록 경로 수정:
- `SCRIPT_DIR` / `APP_DIR` 변수로 절대 경로 참조
- `cd "${APP_DIR}"` 후 `go build` 실행 (go.mod가 app/ 하위에 위치)
- `dist/`, `Info.plist`, `icon.icns` 경로 모두 `SCRIPT_DIR` 기준 절대 경로로 변경

**최종 구조**:
```
/app/           ← Go 소스 + 리소스 (go.mod, main.go, core/, server/, ui/, config/)
/build-work/    ← 빌드 관련 (build.sh, Info.plist, entitlements.plist, dist/)
```
