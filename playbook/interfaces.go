package playbook

// 代理模式¶ https://design-patterns.readthedocs.io/zh_CN/latest/structural_patterns/proxy.html
type Pipeline interface {
	Run() error
}

// https://design-patterns.readthedocs.io/zh_CN/latest/behavioral_patterns/state.html
type Proxy interface {
	Parse() error // 解析
	Exec()
}
