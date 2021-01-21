package roles

/*
roles功能文件夹主要用来实现Shell Filter里面具体的Config Pipeline功能的具体实现，比如：copy（上传文件）、template（模板）、shell等功能。

接口定义：

/ 执行过程生命周期接口
type RoleLifeCycle interface {
	// 环境准备
	Pre()

	// 执行前
	Before()

	// 执行中
	// 返回是否执行信号
	Run() error

	// 执行后
	After()

	// 是否执行hook
	IsHook() (string, string, bool)

	// 钩子函数
	Hooks(string, string, func(string, string) error) error
}

调用YAML文件格式：

config: # 详细配置信息
  - stage: build
    name: template 模板测试
    template:
      src: template.service.j2
      dest: /tmp/docker.service
  - stage: test
    name: 上传文件
    copy:
      src: "{{ .item }}"
      dest: /tmp/{{ .item }}
    with_items:
      - LICENSE
    tags: # 指定主机执行
      - 192.168.0.10
  - stage: what
    name: 非stage测试
    shell: whoami
  - stage: build
    name: 获取go version
    shell: lsb_release -a
  - stage: test
    name: 获取主机名
    shell: "{{.item}}"
    with_items:
    - hostname
    - ip a|grep eth0
    - pwd
    - uname -a
    - docker ps &&
      docker images
    tags:
      - 192.168.0.250
  - stage: test
    name: 查看docker信息
    shell: systemctl status sshd
*/
