.PHONY: push pull install run clean

push:
	git add .
	git commit -m "auto `date`"
	git push origin $(shell git branch|grep '*'|awk '{print $$2}')

pull:
	git pull origin $(shell git branch|grep '*'|awk '{print $$2}')

install: Makefile main.go
	go build
	chmod +x showme 
	./showme

run: main.go
	go run main.go static

clean:
	rm -f 123.mp4
	rm -f 1.db

