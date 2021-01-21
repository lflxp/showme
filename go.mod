module github.com/lflxp/showme

go 1.15

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/astaxie/beego v1.12.1
	github.com/c-bata/go-prompt v0.2.3
	github.com/coreos/bbolt v1.3.2
	github.com/devopsxp/xp v0.0.0-20210120084237-346e7693b0f3
	github.com/go-xorm/xorm v0.7.9
	github.com/google/gopacket v1.1.17
	github.com/google/gops v0.3.13 // indirect
	github.com/jroimartin/gocui v0.4.0
	github.com/juju/errors v0.0.0-20200330140219-3fe23663418f
	github.com/juju/loggo v0.0.0-20190526231331-6e530bcce5d8 // indirect
	github.com/juju/testing v0.0.0-20191001232224-ce9dec17d28b // indirect
	github.com/lflxp/goproxys v0.0.0-20200308164541-294d52d6ffa9
	github.com/lflxp/lflxp-api v0.0.0-20200323063154-619ec5845ffa
	github.com/lflxp/lflxp-monitor v0.0.0-20200324040342-939febaf252c
	github.com/lflxp/lflxp-orzdba v0.0.0-20201102102701-b8902cead805
	github.com/lflxp/lflxp-scan v0.0.0-20200323114511-9ac561b61f89
	github.com/lflxp/lflxp-sflowtool v0.0.0-20200323103145-8e12626667ee
	github.com/lflxp/lflxp-static v0.0.0-20200323072822-0e507513cc6f
	github.com/lflxp/lflxp-tty v0.0.0-20201030085208-96cf85de2130
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/shadowsocks/shadowsocks-go v0.0.0-20190614083952-6a03846ca9c0
	github.com/shirou/gopsutil v3.20.10+incompatible
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.1
	github.com/swaggo/swag v1.6.5
	github.com/ugorji/go/codec v1.1.7
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
)

replace (
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.1.0
	github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
	k8s.io/api => k8s.io/api v0.18.2
	k8s.io/client-go => k8s.io/client-go v0.18.2
)
