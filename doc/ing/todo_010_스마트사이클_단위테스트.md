# TODO-010: 스마트 사이클 단위 테스트

## 우선순위
4순위 (테스트 보완)

## 목표
TODO-003 스마트 사이클 구현 후 `IsNearlyEqual` 판별 로직 및 사이클 전환 로직 테스트.

## 테스트 파일 위치
`app/core/window_controller_test.go` (또는 `coordinate_calculator_test.go`)

## 테스트 케이스 목록
| # | 테스트명 | 시나리오 |
|---|---------|---------|
| 1 | TestIsNearlyEqual_ExactMatch | 정확히 일치하는 크기 → true |
| 2 | TestIsNearlyEqual_WithinTolerance | 허용 오차(5px) 내 → true |
| 3 | TestIsNearlyEqual_OutOfTolerance | 허용 오차 초과 → false |
| 4 | TestNextCycleMode_LeftHalf | left_half → left_1/3 전환 |
| 5 | TestNextCycleMode_Left13 | left_1/3 → left_2/3 전환 |
| 6 | TestNextCycleMode_Left23 | left_2/3 → left_half 복귀 |
| 7 | TestNextCycleMode_RightHalf | right_half 사이클 전환 |
| 8 | TestNextCycleMode_NoMatch | 현재 크기가 어느 단계도 아닐 때 첫 단계 시작 |

## 의존성
- TODO-003 (스마트 사이클 구현) 완료 후 진행

## 검증 방법
```bash
go test ./app/core/ -run TestIsNearlyEqual -v
go test ./app/core/ -run TestNextCycleMode -v
```

## 작업 결과

**상태**: 완료

**수정 파일**: `app/core/window_controller_test.go`

**변경 내용**:
`TestNextCycleMode_*` 3개는 todo-003에서 이미 추가됨. 이번에 `IsSimilar` 직접 단위테스트 3개 추가.

| 테스트 | 결과 |
|--------|------|
| TestIsSimilar_ExactMatch | PASS |
| TestIsSimilar_WithinTolerance | PASS |
| TestIsSimilar_OutOfTolerance | PASS |
| TestNextCycleMode_LeftCycle | PASS (기존) |
| TestNextCycleMode_RightCycle | PASS (기존) |
| TestNextCycleMode_NoCycle | PASS (기존) |

core 패키지 전체 테스트 23개 모두 통과.

**커밋**: `bcc24fc`
