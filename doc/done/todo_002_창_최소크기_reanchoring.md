# TODO-002: 창 최소 크기 Re-anchoring 구현

## 우선순위
1순위 (핵심 버그)

## 문제
커스텀 비율(예: 35%)이 앱(크롬 등)의 최소 폭 제한보다 작을 때:
- 우측 정렬 시 창이 오른쪽 화면 밖으로 짤림
- 스마트 사이클이 짤린 상태를 "미정렬"로 판단해 무한 반복

## 원인
`CoordinateCalculator`는 앱 최소 크기를 모르기 때문에 `x = 모니터너비 × (1 - 비율)`로 계산.
실제 창이 더 넓게 적용되면 오른쪽 경계 초과 발생.

## 구현 위치
`app/core/window_controller_darwin.m` — `setWindowFrame()` 호출 후

## 구현 방법
```
// 우측 정렬 계열 모드에서만 적용
actualBounds = getWindowFrame(pid)
if actualBounds.width > targetWidth {
    correctedX = monitorRight - actualBounds.width
    setWindowPosition(pid, correctedX, actualBounds.y)
}
```

우측 정렬 모드 목록: `right_half`, `right_1/3`, `right_2/3`, `right_custom:*`,
`top_right_1/4`, `bottom_right_1/4`

## 검증 방법 (e2e)
- 크롬에 커스텀 35% 우측 정렬 적용 → 화면 짤림 없는지 확인
- 짤린 상태에서 스마트 사이클 → 다음 단계로 정상 진행 확인

## 참고 문서
- `doc/ing/known_issue_창_최소크기_화면짤림.md`

## 작업 결과

**상태**: 완료

**수정 파일**: `app/core/window_controller.go`, `app/core/window_controller_test.go`

**변경 내용**:
`ExecuteWindowCommand()` 내에 `reanchor()` 함수 호출 추가. `setWindowFrame` 직후 실제 창 위치를 다시 읽어 화면 밖 초과 여부를 판별하고 X/Y 보정.

```go
// 창 이동
SetWindowFrame(pid, targetAbs)

// Re-anchoring: 앱 최소 크기 제한으로 인한 화면 밖 짤림 보정
actual := GetWindowFrame(pid)
corrected := reanchor(actual, targetMonitor, mode, gap)
if corrected != actual {
    SetWindowFrame(pid, corrected)
}
```

`reanchor()` 함수: `right` 계열은 우측 엣지 초과 시 X 보정, `bottom` 계열은 하단 엣지 초과 시 Y 보정.

**추가된 테스트**:
- `TestReanchor_RightOverflow` / `TestReanchor_RightNoOverflow`
- `TestReanchor_BottomOverflow` / `TestReanchor_WithGap` / `TestReanchor_LeftMode_NoChange`
- `TestReanchor_CustomRightMode` / `TestReanchor_BottomRightQuarter`
- `TestIsAlreadyAligned_RightEdge_MinSizeExpanded`

**커밋**: `d1774fa`
