package utils

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type LogFileWriter struct {
	File *os.File
	//write count
	Size int64
}

func (p LogFileWriter) Write(data []byte) (n int, err error) {
	if p.File == nil {
		return 0, errors.New("logFileWriter is nil")
	}
	if p.File == nil {
		return 0, errors.New("file not opened")
	}
	n, e := p.File.Write(data)
	p.Size += int64(n)
	//文件最大 64M byte
	if p.Size > 1024*1024*64 {
		p.File.Close()
		p.File, _ = os.OpenFile(fmt.Sprintf("%s.log", time.Now().Format("20060102150405")), os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
		p.Size = 0
	}
	return n, e
}
