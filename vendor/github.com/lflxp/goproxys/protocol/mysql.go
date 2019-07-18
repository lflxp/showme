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

	log "github.com/sirupsen/logrus"
)

// mysql protocol handler
// load balance
// 参考资料：https://www.flysnow.org/2016/12/26/golang-socket5-proxy.html
func HandleMysqlRequestTCP(client net.Conn) {
	log.Println("mysql conn", client.LocalAddr().String(), client.RemoteAddr().String())
	if client == nil {
		return
	}
	defer client.Close()

	server, err := net.Dial("tcp", net.JoinHostPort("10.128.142.132", "3306"))
	if err != nil {
		log.Error(err)
		return
	}
	defer server.Close()
	// client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) // 响应客户端连接成功
	// 进行转发
	go io.Copy(server, client)
	io.Copy(client, server)
}
