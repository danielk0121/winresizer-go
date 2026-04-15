# WinResizer Go 포팅 개발 기획서

> 기존 Python 구현체를 Go 언어로 재구현하는 프로젝트

---

## 1. 프로젝트 개요

| 항목 | 내용 |
|---|---|
| 프로젝트명 | winresizer-go |
| 목적 | Python 구현체를 Go로 포팅하여 배포 단순화 및 성능 개선 |
| 대상 OS | macOS (Intel / Apple Silicon) |
| 언어 | Go 1.22+ (CGo 포함) |
| 빌드 결과물 | 단일 바이너리 → `.app` 번들 → `.dmg` |

---

## 2. 포팅 동기

### 기존 Python 구현의 한계

| 문제 | 설명 |
|---|---|
| 무거운 배포 | PyInstaller가 Python 런타임 전체를 번들에 포함 → ~100~200MB |
| 복잡한 빌드 | `WinResizer.spec` + `build.sh` + PyInstaller 설치 필요 |
| 사용자 머신 의존 | Python 없이도 실행 가능하지만 번들 크기가 큼 |
| 런타임 오버헤드 | 트레이 앱임에도 Python 인터프리터 상시 기동 |

### Go 포팅 기대 효과

| 항목 | Python | Go |
|---|---|---|
| 배포 크기 | ~100~200 MB | **~15~25 MB** |
| 빌드 방법 | PyInstaller (복잡) | `go build` 한 줄 |
| 사용자 머신 요구사항 | 없음 (번들) | **없음 (바이너리)** |
| 메모리 사용 | ~50~100 MB | **~10~20 MB** |
| 시작 속도 | 느림 (인터프리터 기동) | 빠름 (네이티브) |

> **핵심**: Python도 내부적으로 pyobjc(C 래퍼)를 통해 macOS API를 호출하므로,  
> Go의 CGo 사용은 Python 대비 복잡도 증가가 없음.

---

## 3. 기술 스택

### Go 라이브러리

| 역할 | Python (현재) | Go 대응 |
|---|---|---|
| 메뉴바 트레이 | `rumps` | `github.com/getlantern/systray` |
| 웹 서버 (설정 UI) | `Flask` | `net/http` (표준 라이브러리) |
| 설정 파일 (JSON) | `json` 표준 | `encoding/json` (표준 라이브러리) |
| 글로벌 단축키 | `pynput` | CGo + Carbon `RegisterEventHotKey` |
| 창 제어 (AX API) | `pyobjc` | CGo + `ApplicationServices` 프레임워크 |
| 멀티모니터 감지 | `AppKit.NSScreen` | CGo + `Cocoa NSScreen` |

### CGo 연동 macOS 프레임워크

```
ApplicationServices  → AX API (창 위치/크기 읽기, 쓰기)
Carbon               → 글로벌 단축키 등록/해제
Cocoa (AppKit)       → NSScreen 멀티모니터 정보
```

---

## 4. 프로젝트 구조

```
winresizer-go/
├── main.go                        # 진입점 — systray 초기화, 컴포넌트 연결
├── go.mod
├── go.sum
│
├── core/
│   ├── window_controller.go       # 창 조절 로직 및 스마트 사이클링
│   ├── window_controller_darwin.c # CGo — AX API C 구현체
│   ├── window_controller_darwin.h # CGo — 헤더
│   ├── hotkey_listener.go         # 글로벌 단축키 감지 및 콜백 등록
│   ├── hotkey_listener_darwin.c   # CGo — Carbon RegisterEventHotKey
│   ├── hotkey_listener_darwin.h
│   ├── monitor.go                 # 멀티모니터 감지 및 활성 모니터 계산
│   ├── monitor_darwin.c           # CGo — NSScreen
│   ├── monitor_darwin.h
│   └── config_manager.go         # config.json 읽기/쓰기
│
├── server/
│   ├── web_server.go              # net/http 설정 웹서버
│   └── handlers.go                # API 핸들러 (/api/config, /api/execute 등)
│
├── ui/
│   ├── tray.go                    # systray 메뉴바 UI
│   ├── static/                    # 웹 UI 정적 파일 (JS, CSS) — 기존 재사용
│   └── templates/                 # 웹 UI HTML 템플릿 — 기존 재사용
│
├── build/
│   ├── build.sh                   # .app 번들 및 DMG 빌드 스크립트
│   ├── Info.plist                 # macOS 앱 메타데이터
│   └── entitlements.plist         # 권한 선언 (Accessibility, Input Monitoring)
│
└── doc/
    └── (기존 설계 문서 이관)
```

---

## 5. 핵심 컴포넌트 설계

### 5.1 실행 중 프로세스 구성

| # | 구성 요소 | 역할 | 비고 |
|---|---|---|---|
| 1 | **메인 고루틴** (systray) | macOS 메뉴바 UI 및 앱 생명주기 | CGo AX API 직접 호출 |
| 2 | **HotkeyListener 고루틴** | Carbon 글로벌 단축키 감지 | 백그라운드 상시 대기 |
| 3 | **웹서버 고루틴** | `net/http` 설정 UI 및 API | 랜덤 포트 할당 |

> Python 버전과 스레드 구성 동일. `goroutine` + `channel`로 안전한 통신.

### 5.2 웹서버 API (기존 동일 유지)

| 메서드 | 경로 | 설명 |
|---|---|---|
| GET | `/` | 설정 페이지 (HTML) |
| GET | `/api/status` | 앱 상태 조회 |
| GET | `/api/config` | 현재 설정값 조회 |
| POST | `/api/config` | 설정 저장 및 단축키 리스너 재시작 |
| POST | `/api/config/reset` | 기본값 반환 |
| POST | `/api/execute` | 창 조절 명령 실행 |
| GET | `/api/execute` | 창 조절 명령 실행 (`?mode=left_half`) |

### 5.3 WindowController 모드 (기존 동일 유지)

```
half:   left_half, right_half, top_half, bottom_half
thirds: left_third, center_third, right_third
two_thirds: left_two_thirds, right_two_thirds
quarter: top_left, top_right, bottom_left, bottom_right
full:   maximize, restore
```

### 5.4 스마트 사이클링

동일 명령 반복 시 단계적 순환 (Python 로직 그대로 이식):

```
left_half 1회: 1/2 배치
left_half 2회: 2/3 배치
left_half 3회: 1/3 배치
left_half 4회: 1/2 배치 (반복)
```

---

## 6. CGo 핵심 구현 계획

### 6.1 창 제어 (AX API)

```c
// window_controller_darwin.c
#cgo LDFLAGS: -framework ApplicationServices

// 현재 활성 앱의 PID 조회
pid_t getActiveAppPID();

// 창의 현재 위치/크기 조회
CGRect getWindowFrame(pid_t pid);

// 창 위치/크기 변경
void setWindowFrame(pid_t pid, float x, float y, float w, float h);

// Accessibility 권한 확인
bool checkAccessibilityPermission();
```

### 6.2 글로벌 단축키 (Carbon)

```c
// hotkey_listener_darwin.c
#cgo LDFLAGS: -framework Carbon

// 단축키 등록
EventHotKeyRef registerHotKey(UInt32 keyCode, UInt32 modifiers, int hotkeyID);

// 단축키 해제
void unregisterHotKey(EventHotKeyRef ref);

// 이벤트 루프 실행 (별도 스레드)
void runEventLoop();
```

### 6.3 멀티모니터 (NSScreen)

```c
// monitor_darwin.c
#cgo LDFLAGS: -framework Cocoa

// 전체 모니터 목록 및 영역 조회
MonitorInfo* getAllMonitors(int* count);

// 특정 창이 위치한 모니터 자동 감지
MonitorInfo getActiveMonitor(CGRect windowFrame);
```

---

## 7. 데이터 플로우

### 단축키 → 창 조절

```
사용자 키 입력
  └─ HotkeyListener 고루틴 (Carbon 이벤트)
       └─ channel → WindowController
            └─ getActiveAppPID() [CGo]
                 └─ getActiveMonitor() [CGo]
                      └─ 사이클링 계산 (순수 Go)
                           └─ setWindowFrame() [CGo]
```

### 설정 변경

```
브라우저 설정 변경
  └─ POST /api/config
       └─ config_manager.SaveConfig() [순수 Go]
            └─ channel → HotkeyListener 재시작
                 └─ 기존 단축키 unregisterHotKey() [CGo]
                      └─ 새 설정으로 registerHotKey() [CGo]
```

---

## 8. 개발 단계 (로드맵)

### Phase 1 — CGo 프로토타입 (핵심 검증)
- [ ] CGo + AX API: 활성 창 감지 및 위치/크기 변경
- [ ] CGo + Carbon: 글로벌 단축키 등록/감지
- [ ] CGo + NSScreen: 멀티모니터 정보 조회
- [ ] **검증**: 기존 Python과 동작 동일성 확인

### Phase 2 — 핵심 로직 구현
- [ ] `config_manager.go`: config.json 읽기/쓰기
- [ ] `window_controller.go`: 스마트 사이클링 전체 모드 구현
- [ ] `hotkey_listener.go`: 설정 기반 단축키 동적 등록

### Phase 3 — UI 및 서버
- [ ] `web_server.go`: API 전체 구현
- [ ] `tray.go`: systray 메뉴바 UI
- [ ] 기존 웹 UI (HTML/JS/CSS) 재사용 및 Go `embed` 패키지로 번들링

### Phase 4 — 빌드 및 배포
- [ ] `.app` 번들 구조 생성 (`Info.plist`, 아이콘 포함)
- [ ] `entitlements.plist` 권한 설정 (Accessibility, Input Monitoring)
- [ ] `build.sh` 작성: `go build` → `.app` → `.dmg`
- [ ] Apple Silicon / Intel 크로스 빌드 (`GOARCH=arm64`, `amd64`)
- [ ] Universal Binary 생성 (`lipo` 사용)

### Phase 5 — 테스트 및 안정화
- [ ] 순수 Go 로직 단위 테스트 (사이클링, config 파싱 등)
- [ ] 기존 Python 버전과 기능 동일성 검증
- [ ] macOS 버전별 호환성 확인 (Ventura, Sonoma, Sequoia)

---

## 9. 빌드 시스템

### 개발 환경 요구사항

| 항목 | 요구사항 |
|---|---|
| Go | 1.22 이상 |
| Xcode Command Line Tools | CGo 컴파일에 필요 |
| macOS SDK | 13.0 이상 권장 |

### 빌드 명령

```bash
# 단순 빌드
go build -o winresizer ./...

# .app 번들 생성
./build/build.sh

# Universal Binary (Intel + Apple Silicon)
GOARCH=amd64 go build -o winresizer-amd64 ./...
GOARCH=arm64 go build -o winresizer-arm64 ./...
lipo -create winresizer-amd64 winresizer-arm64 -output winresizer
```

### 사용자 배포 파일 구조

```
WinResizer.app/
├── Contents/
│   ├── MacOS/
│   │   └── winresizer          ← 단일 바이너리 (~15~25 MB)
│   ├── Resources/
│   │   └── AppIcon.icns
│   └── Info.plist
```

---

## 10. 로그 및 설정 파일 경로 (기존 동일 유지)

| 항목 | 경로 |
|---|---|
| 설정 파일 | `~/Library/Application Support/WinResizer/config.json` |
| 로그 파일 | `~/Library/Application Support/WinResizer/log/winresizer_YYYYMMDD_HHMMSS_KST.log` |

---

## 11. 재사용 가능한 기존 자산

| 자산 | 재사용 여부 | 비고 |
|---|---|---|
| 웹 UI (HTML/JS/CSS) | ✅ 100% 재사용 | Go `embed` 패키지로 바이너리에 내장 |
| API 스펙 (엔드포인트) | ✅ 동일 유지 | 하위 호환성 보장 |
| config.json 스키마 | ✅ 동일 유지 | 기존 설정 파일 그대로 사용 가능 |
| 설계 문서 (`doc/`) | ✅ 이관 | Go 버전에 맞게 업데이트 |
| 빌드/배포 스크립트 | ⚠️ 수정 필요 | PyInstaller → `go build` 으로 교체 |
| 테스트 코드 | ⚠️ 재작성 | Python → Go 테스트로 변환 |

---

*본 기획서는 Python 구현체(danielk0121/winresizer)를 기반으로 작성되었습니다.*
