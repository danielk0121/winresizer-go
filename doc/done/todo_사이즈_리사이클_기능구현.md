## 개요
- 단축키를 반복적으로 입력하면 창 크기 조절이 리사이클 되는 기능 추가
- 시나리오
  - 창이 화면 가운데 100*100 크기로 있음
  - 모니터는 2000*2000 사이즈라고 가정
  - 1/2 우측 이동 단축키를 입력함: 이 단축키를 "우측 정렬" 이라고 하자
  - 우측 정렬 1번 입력 > 창 위치 이동 및 창 크기 변경 (가로: 1000, 세로: 2000)
  - 우측 정렬 2번 입력 > 창 위치 그대로, 창 크기 변경 (가로:  800, 세로: 2000)
  - 우측 정렬 3번 입력 > 창 위치 그대로, 창 크기 변경 (가로:  600, 세로: 2000)
  - 우측 정렬 3번 입력 > 창 위치 그대로, 창 크기 변경 (가로:  400, 세로: 2000)
  - 우측 정렬 3번 입력 > 창 위치 그대로, 창 크기 변경 (가로: 1000, 세로: 2000)
  - 반복
- 옵션 선택
  - 창 크기가 단축키 반복할 때마다 점점 줄어드는 기능
  - 창 크기가 단축키 반복할 때마다 점점 커지는 기능
  - 아마 단축키를 각각 별도 지정 할 수 있도록 하는 ux 가 더 편리할 듯
- 개발 전략
  - 일단 좌/우 정렬, 창 크기 증가/감소 4가지 경우로 4개 단축키 설정 기능을 신규로 추가
  - 웹 ui 위치 : 비율 조절 바로 아래, 새로운 색션 인 것 처럼 배치

## 작업 진행
- 모니터 여러개를 건너며 창이 이동하는 리사이클 기능 삭제
- 단축키를 반복적으로 입력하면 사이즈 조절이 리사이클 되는 기능 추가

## 검증
- `window_controller_test.go`에 `calcSizeRecycle` 헬퍼를 통한 단위 테스트 6개 작성
  - `TestSizeRecycle_GrowLeft`: 좌측 확장 시 폭 +10%, X 고정 확인
  - `TestSizeRecycle_ShrinkLeft`: 좌측 축소 시 폭 -10%, X 고정 확인
  - `TestSizeRecycle_GrowRight`: 우측 확장 시 폭 +10%, 우측 엣지 고정 확인
  - `TestSizeRecycle_ShrinkRight`: 우측 축소 시 폭 -10%, 우측 엣지 고정 확인
  - `TestSizeRecycle_MinClamp`: 최소 10% 이하 클램핑 확인
  - `TestSizeRecycle_MaxClamp`: 최대 100% 초과 클램핑 확인
- `go test ./...` 전체 통과

## 결과
- `app/core/window_controller.go`: `isSizeRecycleMode`, `executeSizeRecycle` 함수 추가
- `app/config/default-config.json`: `Size Grow Left`, `Size Shrink Left`, `Size Grow Right`, `Size Shrink Right` 4개 단축키 항목 추가
- `app/ui/templates/index.html`: 사이즈 리사이클 섹션 UI 추가 (커스텀 비율 섹션과 단축키 섹션 사이)
- `app/ui/static/app.js`: 사이즈 리사이클 렌더링 로직 및 i18n 문자열 추가

## 추가 변경 (2026-04-16)
- 최솟값 20%, 최댓값 90%로 범위 조정 (`sizeRecycleMinRatio = 0.20`, `sizeRecycleMaxRatio = 0.90`)
- 경계 도달 시 랩어라운드(wrap-around) 동작:
  - 확장(grow): 현재 폭이 최대(90%) 초과 시 → 최소(20%)로 점프
  - 축소(shrink): 현재 폭이 최소(20%) 미만 시 → 최대(90%)로 점프
- 단위 테스트 갱신: `TestSizeRecycle_GrowWrapsToMin`, `TestSizeRecycle_ShrinkWrapsToMax` 추가
