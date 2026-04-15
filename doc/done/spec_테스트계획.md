# 테스트 계획 및 전략

## 1. 개요
단축키 감지, 윈도우 조절 로직, 좌표 계산 등 핵심 모듈의 신뢰성을 확보하기 위해
순수 Go 로직은 단위 테스트로, CGo 연동 기능은 수동 검증으로 진행한다.

## 2. 테스트 환경
- **프레임워크:** Go 표준 `testing` 패키지
- **실행 명령:**
  ```bash
  # 특정 패키지
  go test ./core/...

  # 전체
  go test ./...
  ```
- **파일 위치:** 각 패키지 폴더 내 `_test.go` 접미사 파일 (별도 `tests/` 폴더 없음)

## 3. 테스트 대상

### 3.1 단위 테스트 (자동화 가능 — 순수 Go 로직)
- **CoordinateCalculator** (`core/coordinate_calculator_test.go`)
  - 1/2, 1/3, 2/3, 1/4 등 모든 모드에서 좌표가 모니터 영역을 벗어나지 않는지 검증
  - 우측/하단 정렬 시 창이 화면 밖으로 짤리지 않도록 Clamp 로직 검증
  - 정수 연산 오차 방지 확인
- **ConfigManager** (`core/config_manager_test.go`)
  - `config.json` 파일 생성, 읽기, 쓰기 후 데이터 일치 여부 확인
  - `default-config.json` 기반 초기화 동작 확인
- **SmartCycle** (`core/window_controller_test.go`)
  - 동일 단축키 반복 시 1/2 → 1/3 → 2/3 → 1/2 순환 검증
  - `is_nearly_equal` 비교 로직 정확성 검증

### 3.2 수동 검증 (CGo 의존 — 자동화 어려움)
- CGo + AX API: 활성 창 감지 및 위치/크기 변경
- CGo + Carbon: 글로벌 단축키 등록/감지
- CGo + NSScreen: 멀티모니터 정보 조회
- 단축키 변경 후 리스너 재시작 및 즉시 반영 확인
- 접근성 권한 없을 때 경고 안내 동작 확인

## 4. 알려진 엣지 케이스 (Python 버전에서 발견된 버그 기반)
- **창 최소 크기 제한 충돌**: 앱의 최소 폭 제한으로 인해 우측 정렬 시 화면 밖 짤림 → Clamp 로직 필수
- **포커스 유실**: AX API로 창 이동 후 포커스가 해제될 수 있음 → 창 이동 직후 명시적 포커스 재부여 필요
- **단축키 충돌**: `Option+Command+Left/Right`는 크롬 탭 이동과 충돌 → 기본값을 `Control+Option+Command` 조합으로 설정

## 5. 배포 조건
- 모든 단위 테스트 통과 필수
- macOS 12 환경에서 실행 확인
- `config.json` 유효성 검사 통과
