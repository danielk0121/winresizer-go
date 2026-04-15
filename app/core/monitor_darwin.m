#include "monitor_darwin.h"
#include <stdlib.h>

// Cocoa NSScreen 을 통해 모니터 정보를 조회합니다.
// AppKit 좌표(하단 원점) → Quartz 좌표(상단 원점)로 변환합니다.
#import <AppKit/AppKit.h>

MonitorInfo* getAllMonitors(int* count) {
    NSArray<NSScreen*>* screens = [NSScreen screens];
    *count = (int)[screens count];
    if (*count == 0) return NULL;

    MonitorInfo* monitors = (MonitorInfo*)malloc(sizeof(MonitorInfo) * (*count));

    // 메인 모니터 전체 높이 (좌표 변환 기준)
    CGFloat mainHeight = [[screens objectAtIndex:0] frame].size.height;

    for (int i = 0; i < *count; i++) {
        NSScreen* screen = [screens objectAtIndex:i];

        // visibleFrame: 메뉴바 및 Dock을 제외한 사용 가능한 영역
        NSRect vf = [screen visibleFrame];

        // AppKit(하단 원점) → Quartz(상단 원점) 변환
        // Quartz_Y = mainHeight - (AppKit_Y + height)
        int quartzY = (int)(mainHeight - (vf.origin.y + vf.size.height));

        monitors[i].x      = (int)vf.origin.x;
        monitors[i].y      = quartzY;
        monitors[i].width  = (int)vf.size.width;
        monitors[i].height = (int)vf.size.height;
    }

    return monitors;
}
