.PHONY: push pull install run clean asset tty build gopacket bindata

build: Makefile main.go asset
	go build
	chmod +x showme 
	./showme -h

install: Makefile main.go asset
	go install
	showme -h

push:
	git add .
	git commit -m "auto `date`"
	git push origin $(shell git branch|grep '*'|awk '{print $$2}')

pull:
	git pull origin $(shell git branch|grep '*'|awk '{print $$2}')

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
	go run main.go tty

bindata:
	@echo 安装预制环境
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/elazarl/go-bindata-assetfs/...

clean:
	rm -f 123.mp4
	rm -f 1.db
	rm -f tty/asset.go

