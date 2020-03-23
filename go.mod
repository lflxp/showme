module github.com/lflxp/showme

go 1.12

require (
	github.com/astaxie/beego v1.12.1
	github.com/c-bata/go-prompt v0.2.3
	github.com/coreos/bbolt v1.3.2
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/xorm v0.7.9
	github.com/google/gopacket v1.1.17
	github.com/jroimartin/gocui v0.4.0
	github.com/juju/errors v0.0.0-20190207033735-e65537c515d7
	github.com/juju/loggo v0.0.0-20190526231331-6e530bcce5d8 // indirect
	github.com/juju/testing v0.0.0-20191001232224-ce9dec17d28b // indirect
	github.com/lflxp/goproxys v0.0.0-20200308164541-294d52d6ffa9
	github.com/lflxp/lflxp-api v0.0.0-20200323063154-619ec5845ffa
	github.com/lflxp/lflxp-monitor v0.0.0-20200323111401-5aff81ef4fa3
	github.com/lflxp/lflxp-orzdba v0.0.0-20200323151836-eee154e8db7a
	github.com/lflxp/lflxp-scan v0.0.0-20200323114511-9ac561b61f89
	github.com/lflxp/lflxp-sflowtool v0.0.0-20200323103145-8e12626667ee
	github.com/lflxp/lflxp-static v0.0.0-20200323072822-0e507513cc6f
	github.com/lflxp/lflxp-tty v0.0.0-20200323112110-fd85c0eb6b1d
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/mitchellh/go-homedir v1.1.0
	github.com/shadowsocks/shadowsocks-go v0.0.0-20190614083952-6a03846ca9c0
	github.com/shirou/gopsutil v2.20.2+incompatible
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.4.0
	github.com/tatsushid/go-fastping v0.0.0-20160109021039-d7bb493dee3e
	github.com/ugorji/go/codec v1.1.7
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
)

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
replace github.com/shirou/gopsutil v2.20.2+incompatible => github.com/shirou/gopsutil v2.18.12+incompatible