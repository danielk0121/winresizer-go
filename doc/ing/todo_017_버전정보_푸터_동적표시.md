# TODO-017: 웹앱 푸터 버전 정보 동적 표시

## 우선순위
3순위 (UI 개선)

## 문제 현황
- 푸터의 버전 문자열(`ver-260414-2328`)이 HTML에 하드코딩되어 있음
- `app/ui/version-time.txt` 파일에 실제 버전 정보가 있으나 미사용

## 목표
- `version-time.txt` 내용을 서버 기동 시 읽어 템플릿에 주입
- 푸터에 동적으로 버전 표시

## 구현 방법
1. `assets.go` embed 지시자에 `version-time.txt` 추가
2. `handlers.go` `handleIndex`에서 `Assets.ReadFile("version-time.txt")` 로 읽어 템플릿 데이터로 전달
3. `index.html` 푸터의 하드코딩 버전 문자열을 `{{.Version}}`으로 교체

## 작업 결과

**상태**: 완료

**수정 파일**:
- `app/ui/assets.go` — `version-time.txt` embed 추가
- `app/server/handlers.go` — `handleIndex`에서 버전 파일 읽어 템플릿에 주입
- `app/ui/templates/index.html` — 푸터 버전 하드코딩 → `{{.Version}}` 치환
