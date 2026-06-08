//go:build !windows

package tray

type Tray struct{}

func New(port int, onOpenChat, onCheckUpdate, onQuit func()) *Tray {
	return &Tray{}
}

func (t *Tray) Run() {
	select {} // Block forever on non-Windows
}

func (t *Tray) Quit() {}
