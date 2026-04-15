#pragma once
#include <stdint.h>

// 모니터 정보 구조체 (Quartz 좌표계 기준 — 상단 왼쪽이 원점)
typedef struct {
    int x;
    int y;
    int width;
    int height;
} MonitorInfo;

// 연결된 모든 모니터의 사용 가능한 영역(메뉴바/Dock 제외)을 반환합니다.
// count: 반환된 모니터 수가 저장됩니다.
// 반환된 포인터는 호출자가 free() 해야 합니다.
MonitorInfo* getAllMonitors(int* count);
