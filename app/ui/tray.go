package ui

import (
	"fmt"
	"os/exec"
	"winresizer/utils"

	"github.com/getlantern/systray"
)

// OnReady는 systray 초기화 완료 시 호출됩니다.
// webPort: gin 웹서버 포트, onQuit: 앱 종료 콜백
func OnReady(webPort int, onQuit func()) func() {
	return func() {
		// 트레이 아이콘 설정
		systray.SetTitle("WinResizer")
		systray.SetTooltip("WinResizer — 창 크기 조절기")

		// 메뉴 구성
		mSettings := systray.AddMenuItem("설정 (Preferences...)", "설정 페이지 열기")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("종료 (Quit)", "WinResizer 종료")

		utils.Log.Infof("트레이 메뉴 준비 완료 (포트: %d)", webPort)

		// 메뉴 이벤트 처리
		go func() {
			for {
				select {
				case <-mSettings.ClickedCh:
					openBrowser(webPort)
				case <-mQuit.ClickedCh:
					utils.Log.Infof("WinResizer 종료")
					systray.Quit()
					if onQuit != nil {
						onQuit()
					}
				}
			}
		}()
	}
}

// OnExit는 systray 종료 시 호출됩니다.
func OnExit() func() {
	return func() {
		utils.Log.Infof("트레이 앱 종료 처리 완료")
	}
}

// openBrowser는 지정된 포트의 설정 페이지를 기본 브라우저로 엽니다.
func openBrowser(port int) {
	url := fmt.Sprintf("http://127.0.0.1:%d", port)
	if err := exec.Command("open", url).Start(); err != nil {
		utils.Log.Errorf("브라우저 열기 실패: %v", err)
	}
}
