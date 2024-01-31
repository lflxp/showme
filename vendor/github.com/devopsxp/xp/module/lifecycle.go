package module

import (
	"fmt"
	"log/slog"

	"github.com/devopsxp/xp/plugin"
)

// 功能：设置默认Start|Stop|Status 实现Filter Interface
type LifeCycle struct {
	name   string
	status plugin.StatusPlugin
}

func (l *LifeCycle) Start() {
	l.status = plugin.Started
	slog.Debug(fmt.Sprintf("%s plugin started.\n", l.name))
}

func (l *LifeCycle) Stop() {
	l.status = plugin.Stopped
	slog.Debug(fmt.Sprintf("%s plugin stopped.\n", l.name))
}

func (l *LifeCycle) Status() plugin.StatusPlugin {
	return l.status
}
