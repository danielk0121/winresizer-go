# todo: 기본 설정값 스크린샷 기반 업데이트

## 개요
사용자가 실제 사용 중인 설정 화면 스크린샷을 기준으로 `default-config.json`을 업데이트한다.
앱 최초 설치 또는 "기본 세팅값으로 초기화" 시 이 값이 적용된다.

## 변경 내용 (2026-04-16)

### 커스텀 비율 항목 추가
기존 default-config.json에 누락되어 있던 커스텀 비율 4개 항목을 추가.

| 키 | mode | 비율 |
|----|------|------|
| Left Custom | left_custom:60 | 60% |
| Right Custom | right_custom:60 | 60% |
| Top Custom | top_custom:70 | 70% |
| Bottom Custom | bottom_custom:70 | 70% |

### 단축키 변경

| 항목 | 변경 전 | 변경 후 |
|------|---------|---------|
| Left 1/2 | ctrl+opt+cmd+ArrowLeft (6400) | 없음 (0) |
| Right 1/2 | ctrl+opt+cmd+ArrowRight (6400) | 없음 (0) |
| 상단 1/2 | ctrl+opt+cmd+ArrowUp (6400) | ctrl+opt+ArrowUp (6144) |
| 하단 1/2 | ctrl+opt+cmd+ArrowDown (6400) | ctrl+opt+ArrowDown (6144) |

### 사이즈 리사이클 단축키 추가

| 항목 | 단축키 | keycode | modifiers |
|------|--------|---------|-----------|
| 좌측 확장 (+10%) | ctrl+opt+ArrowLeft | 123 | 6144 |
| 좌측 축소 (-10%) | ctrl+opt+shift+ArrowLeft | 123 | 6656 |
| 우측 확장 (+10%) | ctrl+opt+ArrowRight | 124 | 6144 |
| 우측 축소 (-10%) | ctrl+opt+shift+ArrowRight | 124 | 6656 |

### modifiers 비트 플래그 참고
- ctrl = 1<<12 = 4096
- opt = 1<<11 = 2048
- shift = 1<<9 = 512
- cmd = 1<<8 = 256
- ctrl+opt = 6144
- ctrl+opt+shift = 6656
- ctrl+opt+cmd = 6400

## 수정 파일
- `app/config/default-config.json`
