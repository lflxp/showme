package module

import (
	. "github.com/devopsxp/xp/plugin"
	log "github.com/sirupsen/logrus"
)

// 功能：设置默认Start|Stop|Status 实现Filter Interface
type LifeCycle struct {
	name   string
	status StatusPlugin
}

func (l *LifeCycle) Start() {
	l.status = Started
	log.Debugf("%s plugin started.\n", l.name)
}

func (l *LifeCycle) Stop() {
	l.status = Stopped
	log.Debugf("%s plugin stopped.\n", l.name)
}

func (l *LifeCycle) Status() StatusPlugin {
	return l.status
}
