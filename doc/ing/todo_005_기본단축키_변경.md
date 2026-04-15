# TODO-005: 기본 단축키 변경 (크롬 충돌 방지)

## 우선순위
2순위 (미구현 기능)

## 문제
현재 `default-config.json`의 기본 단축키가 `Cmd+Opt+방향키`로 설정되어 있어
크롬 탭 이동(Cmd+Opt+Left/Right) 단축키와 충돌 발생.

## 변경 내용
- **현재**: `alt + cmd + shift + left/right/up/down`
- **변경 후**: `Ctrl + Opt + Cmd + 방향키`

## 구현 위치
`app/config/default-config.json` — `shortcuts` 섹션의 `display` 및 `keycode`/`modifiers` 값

## 구현 방법
1. `Ctrl + Opt + Cmd + Left` 조합의 keycode, modifiers 값 확인
2. `default-config.json` 4방향 기본 단축키 값 업데이트
3. 기존 사용자 설정(`~/Library/Application Support/WinResizer/config.json`)에는 영향 없음

## 검증 방법
- 앱 최초 실행(config.json 없는 상태) → 기본 단축키 `Ctrl+Opt+Cmd+방향키`로 등록 확인
- 크롬에서 기존 Cmd+Opt+Left/Right → 탭 이동 정상 동작 확인

## 참고 문서
- `doc/ing/known_issue_단축키_충돌.md`

## 작업 결과

**상태**: 완료

**수정 파일**: `app/config/default-config.json`

**변경 내용**:
4방향 기본 단축키를 `Ctrl+Opt+Cmd+방향키`로 설정.

| 단축키 | keycode | modifiers |
|--------|---------|-----------|
| Ctrl+Opt+Cmd+Left  | 123 | 6400 |
| Ctrl+Opt+Cmd+Right | 124 | 6400 |
| Ctrl+Opt+Cmd+Up    | 126 | 6400 |
| Ctrl+Opt+Cmd+Down  | 125 | 6400 |

modifiers 6400 = cmdKey(256) + optionKey(2048) + controlKey(4096)

기존 `~/Library/Application Support/WinResizer/config.json` 사용자 설정에는 영향 없음.

**커밋**: `9c5e87a`
