.PHONY: push pull install run clean tty build gopacket bindata swag

# 默认位置 以后都保持不变
push: asset pull
	git add .
	git commit -m "${m}"
	git push origin $(shell git branch|grep '*'|awk '{print $$2}')

pull:
	git pull origin $(shell git branch|grep '*'|awk '{print $$2}')

build: Makefile main.go
	GOOS=linux GOARCH=amd64 go build
	chmod +x showme 
	./showme -h

install: Makefile main.go 
	go install -v
	showme -h

gopacket: Makefile main.go  
	@echo please install pcap first
	@echo yum install libpcap-devel
	go install -tags=gopacket
	showme -h

# 静态文件转go二进制文件
asset: bindata
	# cd tty/static && go-bindata -o=../asset.go -pkg=tty ./
	# cd executors/httpstatic/static && go-bindata -o=../asset.go -pkg=httpstatic ./
	# cd boltapi/static && go-bindata -o=../asset.go -pkg=boltapi ./

swag:
	# cd boltapi && swag init

run: main.go
	go run main.go static ${n}

released: 
    # @echo goreleaser release

# tty功能测试
tty: 
	go run main.go tty -w -m 1 -d -a -u admin -p admin bash 
	# go run main.go tty -w -m 10 -r -d showme proxy http

bindata:
	@echo 安装预制环境
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/elazarl/go-bindata-assetfs/...

clean:
	rm -f 123.mp4
	rm -f 1.db
	rm -f tty/asset.go
	rm -f executors/httpstatic/asset.go
	rm -f showme
	rm -f *.crt
	rm -f *.key
	rm -f *.csr
	rm -f tls/*
	rm -f *.tar.gz

.PHONY: windows
windows: 
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build

.PHONY: crt csr key
crt: csr
	openssl x509 -req -sha256 -days 3650 -in tls/server.csr -signkey tls/server.key -out tls/server.crt

csr: key
	openssl req -nodes -new -key tls/server.key -subj "/CN=www.lflxp.cn" -out tls/server.csr

key: clean
	openssl genrsa -out tls/server.key 2048

tags:
	@echo 只需要在go build指令后用-tags指定编译条件即可
	go build -tags govcl

test:
	go test -v -cover ./utils/...