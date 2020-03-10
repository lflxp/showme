.PHONY: push pull install run clean asset tty build gopacket bindata

# 默认位置 以后都保持不变
push: asset pull
	git add .
	git commit -m "${m}"
	git push origin $(shell git branch|grep '*'|awk '{print $$2}')

pull:
	git pull origin $(shell git branch|grep '*'|awk '{print $$2}')

build: Makefile main.go asset
	go build
	chmod +x showme 
	./showme -h

install: Makefile main.go asset
	go install
	showme -h

gopacket: Makefile main.go asset
	go build -tags=gopacket
	chmod +x showme 
	./showme -h

# 静态文件转go二进制文件
asset: bindata
	cd tty/static && go-bindata -o=../asset.go -pkg=tty ./

run: main.go
	go run main.go static

# tty功能测试
tty: asset
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
	rm -f showme

