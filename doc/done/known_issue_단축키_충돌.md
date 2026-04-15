# 알려진 이슈: 단축키 조합 충돌

## 1. 현상
크롬 브라우저가 활성화된 상태에서 `Option+Command+Left/Right` 단축키 입력 시,
창 크기 조절 대신 크롬의 **탭 이동** 기능이 실행됨.

## 2. 원인
macOS 애플리케이션이 시스템 전역 단축키보다 우선순위가 높은 자체 단축키를 가지고 있어
Carbon `RegisterEventHotKey`로 등록된 이벤트보다 앱 단축키가 먼저 처리됨.

## 3. 조치 사항
- **기본 단축키 변경**: 충돌 가능성이 낮은 `Control+Option+Command` 조합을 기본값으로 설정
  - `Control+Option+Command+Left`: 좌측 절반
  - `Control+Option+Command+Right`: 우측 절반
  - `Control+Option+Command+Up`: 위쪽 절반
  - `Control+Option+Command+Down`: 아래쪽 절반
- **UI 경고**: 사용자가 단축키 설정 시 시스템 예약 단축키 입력 시 경고 노출

## 4. 구현 시 고려사항
- `default-config.json`에 `Control+Option+Command` 조합을 기본값으로 반영할 것
- 웹 UI에서 단축키 녹화 시 충돌 가능 조합(`Cmd+Alt+방향키` 등)에 대한 경고 메시지 제공
