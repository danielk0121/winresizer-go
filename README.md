# winresizer-go

macOS 창 크기/위치를 글로벌 단축키로 조절하는 트레이 앱. Python 구현체를 Go로 포팅하여 배포 단순화 및 성능 개선.

## 주요 특징

- 글로벌 단축키로 창을 화면의 절반, 1/3, 2/3, 사분면 등으로 즉시 배치
- 동일 단축키 반복 시 비율 순환 (스마트 사이클링)
- 멀티모니터 지원 — 창이 위치한 모니터 기준으로 동작
- 브라우저 기반 설정 UI (메뉴바 아이콘 클릭으로 접근)
- 단일 바이너리 배포 (~15–25 MB), 별도 런타임 불필요

## 요구사항

| 항목 | 내용 |
|---|---|
| OS | macOS 12.0 이상 (Intel / Apple Silicon) |
| Go | 1.22 이상 |
| 기타 | Xcode Command Line Tools (CGo 빌드용) |

## 빌드

```bash
# 단순 빌드
cd app
go build -o winresizer ./...

# .app 번들 + DMG 생성
./build/build.sh

# Universal Binary (Intel + Apple Silicon)
GOARCH=amd64 go build -o winresizer-amd64 ./...
GOARCH=arm64 go build -o winresizer-arm64 ./...
lipo -create winresizer-amd64 winresizer-arm64 -output winresizer
```

## 프로젝트 구조

```
winresizer-go/
├── app/            # Go 소스코드 루트
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── config/     # 기본 설정 파일 (default-config.json)
│   ├── core/       # 창 제어, 단축키, 멀티모니터, 설정 관리 (CGo 포함)
│   ├── server/     # 설정 웹서버 및 API 핸들러
│   ├── ui/         # 메뉴바 트레이 UI 및 웹 UI 정적 파일
│   ├── utils/      # 로거 등 공통 유틸리티
│   └── build/      # .app 번들 및 DMG 빌드 스크립트, plist
├── doc/            # 설계 문서
└── ref/            # 참고 자료
```

## Go 개발 참고

### 테스트 파일 위치
Go는 별도 `tests/` 폴더를 만들지 않는다. 테스트 파일은 소스 파일과 **같은 폴더**에 `_test.go` 접미사로 작성한다.

```
core/
├── window_controller.go
├── window_controller_test.go   ← 테스트 파일
├── hotkey_listener.go
└── hotkey_listener_test.go     ← 테스트 파일
```

```bash
# 특정 패키지 테스트
go test ./core/...

# 전체 테스트
go test ./...
```

## 설정 파일 경로

| 항목 | 경로 |
|---|---|
| 설정 | `~/Library/Application Support/WinResizer/config.json` |
| 로그 | `~/Library/Application Support/WinResizer/log/` |

## 설정 파일 구조 (default-config.json vs config.json)

파이썬 구현체와 동일하게 **2개 파일** 구조를 사용한다.

| 파일 | 위치 | 역할 |
|---|---|---|
| `default-config.json` | `app/config/default-config.json` (소스 내장) | 앱 기본값 정의. 소스코드에 포함되어 배포됨 |
| `config.json` | `~/Library/Application Support/WinResizer/config.json` | 사용자 실제 설정. 런타임에 생성/수정됨 |

**동작 흐름:**

```
앱 시작
  └─ config.json 존재?
       ├─ YES → config.json 로드 (사용자 설정 사용)
       └─ NO  → default-config.json 로드 → config.json으로 복사 생성
```

- `default-config.json`은 읽기 전용 기준값. 앱이 직접 수정하지 않는다.
- 사용자가 설정 UI에서 변경하면 `config.json`에만 저장된다.
- `POST /api/config/reset` 호출 시 `default-config.json` 값을 응답으로 반환하지만, `config.json` 파일은 덮어쓰지 않는다 (설계 의도).
- 관련 구현: `app/core/config_manager.go` — `LoadConfig()`, `LoadDefaultConfig()`
