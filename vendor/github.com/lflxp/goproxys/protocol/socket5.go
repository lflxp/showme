package protocol

/**
Socket5的客户端和服务端进行双方授权验证通过之后，就开始建立连接了。连接由客户端发起，告诉Sokcet服务端客户端需要访问哪个远程服务器，其中包含，远程服务器的地址和端口，地址可以是IP4，IP6，也可以是域名。

VER	CMD	RSV	ATYP	DST.ADDR	DST.PORT
1	1	X’00’	1	Variable	2

VER代表Socket协议的版本，Soket5默认为0x05，其值长度为1个字节
CMD代表客户端请求的类型，值长度也是1个字节，有三种类型
CONNECT X’01’
BIND X’02’
UDP ASSOCIATE X’03’
RSV保留字，值长度为1个字节
ATYP代表请求的远程服务器地址类型，值长度1个字节，有三种类型
IP V4 address: X’01’
DOMAINNAME: X’03’
IP V6 address: X’04’
DST.ADDR代表远程服务器的地址，根据ATYP进行解析，值长度不定。
DST.PORT代表远程服务器的端口，要访问哪个端口的意思，值长度2个字节
*/

import (
	"io"
	"net"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// socket5 protocol handler
// 参考资料：https://www.flysnow.org/2016/12/26/golang-socket5-proxy.html
func HandleSocket5RequestTCP(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("Byte Origin: %v", b[:])
	if b[0] == 0x05 { // 只处理Socket5协议
		//客户端回应：Socket服务端不需要验证方式
		client.Write([]byte{0x05, 0x00})
		n, err = client.Read(b[:])
		if err != nil {
			log.Error(err)
			return
		}
		var host, port string
		switch b[3] {
		case 0x01: // IP V4
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
		case 0x03: //域名
			host = string(b[5 : n-2]) //b[4]表示域名的长度
		case 0x04: //IP V6
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

		log.Debugf("Socket5: %s:%s", host, port)
		// 获得了请求的host和port，就开始拨号吧
		server, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			log.Error(err)
			return
		}
		defer server.Close()
		client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) // 响应客户端连接成功
		// 进行转发
		go io.Copy(server, client)
		io.Copy(client, server)
	}
}
