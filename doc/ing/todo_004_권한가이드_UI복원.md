# TODO-004: 권한 가이드 UI 복원

## 우선순위
2순위 (미구현 기능)

## 문제
`app/ui/static/app.js`에서 개발 편의상 권한 미승인 시 가이드 오버레이 표시 로직을 비활성화한 상태.
권한이 없는 사용자가 앱을 처음 실행할 때 안내가 없음.

## 구현 위치
`app/ui/static/app.js` — 권한 상태 체크 및 오버레이 표시 부분

## 구현 방법
1. `GET /api/status` 응답에서 `accessibility_granted`, `input_monitoring_granted` 확인
2. 하나라도 `false`이면 가이드 오버레이 표시
3. 오버레이에 `open_accessibility` / `open_input_monitoring` 버튼 제공
4. 권한 승인 후 오버레이 자동 해제 (폴링 또는 사용자 수동 새로고침)

## 검증 방법 (e2e)
- 접근성 권한 미부여 상태에서 앱 실행 → 가이드 오버레이 표시 확인
- 권한 부여 후 → 오버레이 사라지고 정상 UI 표시 확인

## 참고 문서
- `doc/ing/spec_권한_가이드_UX.md`

## 작업 결과

**상태**: 완료

**수정 파일**: `app/ui/static/app.js`

**변경 내용**:
개발용 강제 숨김 코드 제거, 권한 상태 기반 오버레이 표시/숨김 함수 추가.

- `showGuideOverlay(acc, inp)`: 미승인 권한에 따라 단계 표시, 완료 단계는 `step-done` 클래스 적용, 주 버튼 onclick을 미승인 권한 오픈 커맨드로 설정
- `hideGuideOverlay()`: 오버레이 숨김
- `checkStatus()` 내 로직: 모든 권한 승인 시 `hideGuideOverlay()`, 하나라도 미승인 시 `showGuideOverlay(acc, inp)` 호출

**커밋**: `ebc7600`
