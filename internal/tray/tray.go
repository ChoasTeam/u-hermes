//go:build windows

package tray

import (
	"fmt"
	"os"

	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

type Tray struct {
	onOpenChat    func()
	onCheckUpdate func()
	onQuit        func()
	port          int
}

func New(port int, onOpenChat, onCheckUpdate, onQuit func()) *Tray {
	return &Tray{
		port:          port,
		onOpenChat:    onOpenChat,
		onCheckUpdate: onCheckUpdate,
		onQuit:        onQuit,
	}
}

func (t *Tray) Run() {
	systray.Run(t.onReady, t.onExit)
}

func (t *Tray) onReady() {
	systray.SetIcon(getIcon())
	systray.SetTitle("U-Hermes")
	systray.SetTooltip("U-Hermes 运行中")

	mOpen := systray.AddMenuItem("打开聊天界面", "打开浏览器")
	systray.AddMenuItem("查看状态", "查看运行状态")
	systray.AddSeparator()
	mUpdate := systray.AddMenuItem("检查更新", "检查新版本")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出 U-Hermes", "安全退出")

	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				open.Run(fmt.Sprintf("http://127.0.0.1:%d/chat", t.port))
			case <-mUpdate.ClickedCh:
				if t.onCheckUpdate != nil {
					t.onCheckUpdate()
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func (t *Tray) onExit() {
	if t.onQuit != nil {
		t.onQuit()
	}
	os.Exit(0)
}

func (t *Tray) Quit() {
	systray.Quit()
}

func getIcon() []byte {
	return []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
}
