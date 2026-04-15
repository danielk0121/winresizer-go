# TODO-001: 창 이동 후 포커스 유실 수정

## 우선순위
1순위 (핵심 버그)

## 문제
단축키로 창을 이동/크기 조절한 직후, 연달아 단축키 입력 시 반응 없음.
스마트 사이클 연속 입력이 불가능해지는 증상.

## 원인
- AX API로 창 위치/크기 변경 후 macOS가 해당 창의 포커스 상태를 일시 해제
- Carbon 이벤트 핸들러에서 수식 키(Cmd, Ctrl 등) 상태가 함께 초기화되어 후속 입력 무시

## 구현 위치
`app/core/window_controller_darwin.m` — `setWindowFrame()` 함수 직후

## 구현 방법
```c
// setWindowFrame() 호출 후 포커스 재부여
AXUIElementSetAttributeValue(appElement, kAXFrontmostAttribute, kCFBooleanTrue);
```

## 검증 방법 (e2e)
1. Chrome 활성화
2. 좌측 절반 단축키 입력 → 창 이동 확인
3. 즉시 동일 단축키 재입력 → 1/3으로 스마트 사이클 동작 확인 (포커스 유실 없어야 함)

## 참고 문서
- `doc/ing/known_issue_포커스_유실.md`

## 작업 결과

**상태**: 완료

**수정 파일**: `app/core/window_controller_darwin.m`

**변경 내용**:
`activateApp()` 함수에 `kAXFrontmostAttribute` AX API 호출을 `NSRunningApplication activate` 이전에 추가.
```c
void activateApp(pid_t pid) {
    // 1차: AX API로 kAXFrontmostAttribute 직접 설정 (창 단위 포커스 재부여)
    AXUIElementRef appRef = AXUIElementCreateApplication(pid);
    if (appRef) {
        AXUIElementSetAttributeValue(appRef, kAXFrontmostAttribute, kCFBooleanTrue);
        CFRelease(appRef);
    }
    // 2차: NSRunningApplication activate (앱 레벨 포커스 재부여)
    NSRunningApplication* app = [NSRunningApplication runningApplicationWithProcessIdentifier:pid];
    if (app) {
        [app activateWithOptions:NSApplicationActivateIgnoringOtherApps];
    }
}
```

**커밋**: `d1774fa` (embed 정적 파일, 트레이 아이콘 내장; 단축키 keycode/modifiers 방식으로 교체; deleteHotkey 버그 수정)
