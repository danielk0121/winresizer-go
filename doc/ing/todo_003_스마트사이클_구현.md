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
