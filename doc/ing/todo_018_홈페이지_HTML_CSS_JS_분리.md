# 홈페이지 HTML·CSS·JS 파일 분리

## 목표

`homepage/index.html`에 포함된 인라인 CSS와 JavaScript를 각각 외부 파일로 분리하여 유지보수성을 높인다.

## 작업 범위

- CSS를 `homepage/style.css`로 이동하고 `index.html`에서 링크한다.
- JavaScript를 `homepage/app.js`로 이동하고 `index.html`에서 불러온다.
- 기존 한국어·영어 전환 및 화면 스타일이 동일하게 동작하는지 확인한다.
- 이어서 가로 340px 모바일 환경에서 버튼 문구와 레이아웃이 깨지지 않도록 CSS를 보완한다.

## 검증 기준

- 340px 뷰포트에서 가로 스크롤 없이 주요 영역이 표시된다.
- 다운로드·GitHub 버튼의 문구가 글자 단위로 세로 줄바꿈되지 않는다.
- 언어 전환 버튼이 정상 동작한다.
