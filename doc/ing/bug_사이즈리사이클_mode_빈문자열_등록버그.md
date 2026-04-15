# 버그: 사이즈 리사이클 단축키 mode 빈 문자열로 등록

## 현상
사이즈 리사이클 단축키(ctrl+opt+ArrowLeft 등)를 입력해도 창이 움직이지 않고 로그에 오류 출력.

```
2026-04-16 07:41:55 [DEBUG] 단축키 감지: id=1 mode=
2026-04-16 07:41:55 [WARN] 좌표 계산 실패 (): 알 수 없는 모드: 
2026-04-16 07:41:55 [INFO] 다음 모니터로 이동: pid=789
```

`mode=` — mode가 빈 문자열로 등록되어 있어 좌표 계산이 실패함.

## 원인

### 1. UI 저장 시 mode 미설정
`app/ui/static/app.js`의 `saveConfig()`에서 커스텀 비율 키(`Left Custom` 등)는 `dir_custom:pct` 형식으로 mode를 재구성하지만, 사이즈 리사이클 키(`Size Grow Left` 등)는 mode를 재구성하지 않는다.

`loadConfigUI()`의 SIZE_RECYCLE_KEYS 초기화 코드:
```js
if (!config.shortcuts[name]) {
    config.shortcuts[name] = { display: '', mode: '', keycode: 0, modifiers: 0 };
}
```
키가 이미 존재하면 건너뛰는데, `mode: ""` 상태로 저장된 경우엔 빈 문자열을 그대로 유지한다.

### 2. 단축키 등록 시 mode 검증 없음
`StartHotkeyManager`에서 `Keycode != 0` 만 확인하고 `mode == ""` 여부는 확인하지 않아, mode가 빈 문자열인 단축키도 Carbon에 등록된다.

### 3. 기존 config.json 마이그레이션 미비
새 단축키 항목(`Size Grow Left` 등)을 default-config.json에 추가했지만, 이미 사용자 config.json에 해당 키가 `mode: ""`로 저장되어 있어 자동 보정되지 않음.

## 해결 방안

### 수정 1: `window_controller.go` — mode 빈 문자열 단축키 등록 스킵
```go
if sc.Keycode == 0 || sc.Mode == "" {
    continue
}
```

### 수정 2: `app.js` — loadConfigUI에서 SIZE_RECYCLE_KEYS mode 항상 보정
```js
const SIZE_RECYCLE_MODE_MAP = {
    'Size Grow Left':   'size_grow_left',
    'Size Shrink Left': 'size_shrink_left',
    'Size Grow Right':  'size_grow_right',
    'Size Shrink Right':'size_shrink_right',
};
for (const name of SIZE_RECYCLE_KEYS) {
    if (!config.shortcuts[name]) {
        config.shortcuts[name] = { display: '', mode: SIZE_RECYCLE_MODE_MAP[name], keycode: 0, modifiers: 0 };
    } else {
        // mode가 잘못된 경우 보정
        config.shortcuts[name].mode = SIZE_RECYCLE_MODE_MAP[name];
    }
}
```

## 수정 파일
- `app/core/window_controller.go`: mode 빈 문자열 단축키 등록 스킵
- `app/ui/static/app.js`: SIZE_RECYCLE_KEYS mode 항상 올바른 값으로 보정
