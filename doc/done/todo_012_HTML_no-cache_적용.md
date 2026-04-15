# TODO-012: HTML no-cache 메타 태그 적용

## 우선순위
2순위 (캐시 제어 강화)

## 목표
`app/ui/templates/index.html`의 `<head>` 섹션에 HTTP 수준 no-cache에 더해 클라이언트(브라우저) 수준 no-cache 메타 태그를 추가한다.

## 현재 상태
- `app/server/web_server.go`: `noCacheMiddleware()`가 모든 응답에 HTTP 헤더 수준 no-cache 적용 중
  ```
  Cache-Control: no-store, no-cache, must-revalidate, max-age=0
  Expires: 0
  ```
- `app/ui/templates/index.html`: `<meta>` no-cache 태그 없음

## 구현 방법
`index.html` `<head>` 섹션에 아래 태그 추가:

```html
<meta http-equiv="Cache-Control" content="no-store, no-cache, must-revalidate, max-age=0">
<meta http-equiv="Pragma" content="no-cache">
<meta http-equiv="Expires" content="0">
```

## 검증 방법
- 브라우저 DevTools → Network 탭 → `index.html` 응답 헤더에 `Cache-Control: no-store` 확인
- 페이지 소스에서 메타 태그 존재 확인

## 작업 결과

**상태**: 완료

**수정 파일**: `app/ui/templates/index.html`

**변경 내용**:
`<head>` 섹션에 no-cache 메타 태그 3개 추가 (line 6~8):

```html
<meta http-equiv="Cache-Control" content="no-store, no-cache, must-revalidate, max-age=0">
<meta http-equiv="Pragma" content="no-cache">
<meta http-equiv="Expires" content="0">
```
