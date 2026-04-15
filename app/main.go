package main

import (
	"winresizer/core"
	"winresizer/server"
	"winresizer/ui"
	"winresizer/utils"

	"github.com/getlantern/systray"
)

func main() {
	utils.Log.Infof("WinResizer 시작")

	// 1. 웹서버 초기화 및 백그라운드 시작
	srv := server.New()
	go srv.Start()

	// 런타임 정보(port, pid) 설정 파일에 기록
	if err := core.SaveRuntimeInfo(srv.Port); err != nil {
		utils.Log.Warnf("런타임 정보 기록 실패: %v", err)
	}

	// 2. 단축키 리스너 시작 (별도 goroutine — Carbon 이벤트 루프)
	go core.StartHotkeyManager()

	// 3. systray 시작 (메인 스레드 점유 — macOS 요구사항)
	systray.Run(
		ui.OnReady(srv.Port, func() {
			core.StopHotkeyListener()
			utils.Log.Close()
		}),
		ui.OnExit(),
	)
}
