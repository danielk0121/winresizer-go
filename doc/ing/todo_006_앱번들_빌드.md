# TODO-006: 앱 번들(.app) 빌드 구성

## 우선순위
3순위 (배포 준비)

## 목표
macOS 표준 `.app` 번들 구조로 패키징하여 더블클릭 실행 가능하도록 구성.

## 구현 방법
```
WinResizer.app/
  Contents/
    Info.plist        ← 앱 메타데이터 (Bundle ID, 버전, 권한 등)
    MacOS/
      winresizer      ← 빌드된 바이너리
    Resources/
      icon.icns       ← 앱 아이콘 (icns 형식)
    config/
      default-config.json
```

### Info.plist 필수 항목
- `CFBundleIdentifier`: `com.winresizer.app`
- `CFBundleName`: `WinResizer`
- `NSAccessibilityUsageDescription`: 접근성 권한 설명
- `com.apple.security.automation.apple-events`: AX API 사용

### 빌드 스크립트 (`scripts/build_app.sh`)
```bash
go build -o WinResizer.app/Contents/MacOS/winresizer ./app
cp -r app/config WinResizer.app/Contents/
# iconutil로 icon.icns 생성
```

## 검증 방법
- `.app` 더블클릭 → 트레이 아이콘 표시 확인
- `~/Library/Application Support/WinResizer/config.json` 생성 확인

## 작업 결과

**상태**: 완료

**수정 파일**: `app/build/build.sh`

**변경 내용**:
기존 스크립트는 `GOARCH=amd64/arm64` 크로스컴파일을 시도했으나, CGo 포함 코드는 Apple Silicon에서 amd64 크로스컴파일이 실패. 이후 todo-008에서 Xcode clang이 크로스컴파일을 지원함을 확인하여 Universal Binary로 전환.

빌드 결과:
```
dist/
├── WinResizer.app/
│   └── Contents/
│       ├── Info.plist
│       ├── MacOS/WinResizer   ← Universal Binary (arm64 + x86_64)
│       └── Resources/
└── WinResizer.dmg
```

`dist/`는 `.gitignore`에 포함되어 버전 관리 대상 제외.

**커밋**: `0bc7865` (단독 빌드), `dd2c211` (Universal Binary로 교체)
