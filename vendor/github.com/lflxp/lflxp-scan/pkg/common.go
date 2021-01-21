package pkg

import "net"

func GetIps() []string {
	rs := []string{}
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		rs = append(rs, err.Error())
		return rs
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				rs = append(rs, ipnet.IP.String())
			}
		}
	}
	return rs
}
