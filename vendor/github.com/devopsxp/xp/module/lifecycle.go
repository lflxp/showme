package module

import (
	"github.com/devopsxp/xp/plugin"
	log "github.com/sirupsen/logrus"
)

// 功能：设置默认Start|Stop|Status 实现Filter Interface
type LifeCycle struct {
	name   string
	status plugin.StatusPlugin
}

func (l *LifeCycle) Start() {
	l.status = plugin.Started
	log.Debugf("%s plugin started.\n", l.name)
}

func (l *LifeCycle) Stop() {
	l.status = plugin.Stopped
	log.Debugf("%s plugin stopped.\n", l.name)
}

func (l *LifeCycle) Status() plugin.StatusPlugin {
	return l.status
}
