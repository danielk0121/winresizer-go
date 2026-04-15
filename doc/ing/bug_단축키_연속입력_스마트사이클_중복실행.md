# 버그: 단축키 연속 입력 시 스마트 사이클 중복 실행

## 현상
단축키를 빠르게 반복 입력하면, 한 번 입력에 스마트 사이클이 여러 단계씩 건너뛰며 연속 실행된다.

```
2026-04-16 07:32:12 [DEBUG] 스마트 사이클: left_half → left_1/3
2026-04-16 07:32:12 [DEBUG] 스마트 사이클: left_1/3 → left_2/3
2026-04-16 07:32:12 [DEBUG] 스마트 사이클: left_2/3 → left_half
...
```
로그 타임스탬프가 모두 동일 초(07:32:12)에 몰려 있어, 단 1회 단축키 입력에 여러 사이클이 실행됨을 확인.

## 원인
`ExecuteWindowCommand` 실행이 완료(창 이동 + AX API 반영)되기 전에 Carbon 이벤트 루프가 다음 단축키 이벤트를 콜백으로 전달한다.
콜백은 별도 goroutine에서 처리되므로, 이전 명령이 끝나지 않은 상태에서 `GetWindowFrame`을 읽으면 창이 아직 이전 위치에 있는 것으로 보여 `isAlreadyAligned`가 "이미 정렬됨"으로 판단, 즉시 스마트 사이클을 발동한다.
이 과정이 연쇄적으로 반복된다.

## 해결 방안
`ExecuteWindowCommand` 실행 중에는 추가 호출을 무시하는 **처리 중 플래그(debounce)**를 적용한다.
- `atomic.Bool` 타입의 `commandRunning` 플래그 사용
- 진입 시 `CompareAndSwap(false, true)`로 선점, 종료 시 `Store(false)` 해제
- 이미 실행 중이면 즉시 `nil` 반환 (이벤트 드롭)

## 수정 내용
- `app/core/window_controller.go`
  - `commandRunning atomic.Bool` 플래그 추가
  - `ExecuteWindowCommand` 진입 시 `CompareAndSwap(false, true)`로 선점, `defer`로 해제
  - 이미 실행 중이면 즉시 드롭하고 `nil` 반환
