package proxy

import (
	"log"
	"net"

	"github.com/lflxp/goproxys/protocol"
)

func RunProxy(types string) {
	if types == "httprp" {
		protocol.RunHttpProxy()
	} else {
		// cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// config := &tls.Config{Certificates: []tls.Certificate{cer}}

		l, err := net.Listen("tcp", ":8081")
		// l, err := tls.Listen("tcp", ":8081", config)
		if err != nil {
			log.Panic(err)
		}
		log.Println("Started Proxy")

		if types == "http" {
			log.Println("Http Proxy Listening port: 8081")
		} else if types == "socket5" {
			log.Println("Socket5 Proxy Listening port: 8081")
		} else if types == "mysql" {
			log.Println("Mysql Proxy Listening port: 8081")
		} else if types == "ss" {
			log.Println("Socket5 Cipher Proxy Listening port: 8081")
		}

		for {
			client, err := l.Accept()
			if err != nil {
				log.Panic(err)
			}
			log.Println(client.RemoteAddr().String())

			if types == "http" {
				go protocol.HandleHttpRequestTCP(client)
			} else if types == "socket5" {
				go protocol.HandleSocket5RequestTCP(client)
			} else if types == "mysql" {
				go protocol.HandleMysqlRequestTCP(client)
			} else if types == "ss" {
				go protocol.HandleSocket5CipherRequestTCP(client)
			}
		}
	}

}
