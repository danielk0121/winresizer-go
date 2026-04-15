#pragma once
#include <stdbool.h>
#include <sys/types.h>

// 창의 위치와 크기를 나타내는 구조체
typedef struct {
    float x;
    float y;
    float width;
    float height;
} WindowFrame;

// Accessibility 권한 확인
bool checkAccessibilityPermission(void);

// 입력 모니터링 권한 확인
bool checkInputMonitoringPermission(void);

// 현재 활성 앱의 PID 반환 (없으면 -1)
pid_t getActiveAppPID(void);

// PID로 활성 창의 위치/크기 반환
WindowFrame getWindowFrame(pid_t pid);

// PID로 창의 위치/크기 변경
void setWindowFrame(pid_t pid, float x, float y, float width, float height);

// 앱을 최상위로 올려 포커스를 재부여합니다 (창 이동 후 포커스 유실 방지)
void activateApp(pid_t pid);
