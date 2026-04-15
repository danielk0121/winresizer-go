# TODO-003: 스마트 사이클 구현

## 우선순위
2순위 (미구현 기능)

## 기능 설명
동일 단축키를 반복 입력할 때 창 크기를 순환 변경:
- `1/2 → 1/3 → 2/3 → 1/2` 순서로 사이클

## 구현 위치
`app/core/window_controller.go` — `ExecuteWindowCommand()` 내부

## 구현 방법
```go
// 현재 창 상태와 요청 모드를 비교하여 다음 사이클 모드 결정
// IsNearlyEqual(actual, expected, tolerance) 로 현재 상태 판별
func nextCycleMode(currentMode string, currentBounds, monitorBounds Rect) string {
    // left_half → left_1/3 → left_2/3 → left_half
    // right_half → right_1/3 → right_2/3 → right_half
}
```

- `IsNearlyEqual` 함수: 실제 창 크기와 기대 비율의 오차 허용 범위(예: 5px) 내 판별
- 사이클 맵을 정의하여 현재 모드 → 다음 모드 전환

## 검증 방법 (e2e)
1. Chrome 좌측 절반 단축키 입력 → `left_half` 적용
2. 동일 단축키 재입력 → `left_1/3` 전환 확인
3. 동일 단축키 재입력 → `left_2/3` 전환 확인
4. 동일 단축키 재입력 → `left_half` 복귀 확인

## 참고 문서
- `doc/ing/spec_스마트사이클.md`

## 작업 결과

**상태**: 완료

**수정 파일**: `app/core/window_controller.go`, `app/core/window_controller_test.go`

**변경 내용**:
`cycleMap`과 `nextCycleMode()` 함수 추가. `ExecuteWindowCommand()` 내에서 `isAlreadyAligned` 판정 시 사이클 여부 먼저 확인 후 `moveToNextDisplay` 폴백.

```go
var cycleMap = map[string]string{
    "left_half": "left_1/3",  "left_1/3": "left_2/3",  "left_2/3": "left_half",
    "right_half": "right_1/3", "right_1/3": "right_2/3", "right_2/3": "right_half",
}

// 이미 해당 위치에 있으면 스마트 사이클 시도, 사이클 없으면 다음 모니터로 이동
if isAlreadyAligned(currentFrame, targetAbs, mode, targetMonitor) {
    if cycled := nextCycleMode(mode); cycled != "" {
        return ExecuteWindowCommand(cycled)
    }
    return moveToNextDisplay(pid, currentFrame, monitors, targetMonitor)
}
```

상하/쿼터/maximize 등 사이클 없는 모드는 `nextCycleMode`가 `""` 반환 → 다음 모니터로 이동.

**추가된 테스트**:
- `TestNextCycleMode_LeftCycle` / `TestNextCycleMode_RightCycle` / `TestNextCycleMode_NoCycle`

**커밋**: `d1774fa`
