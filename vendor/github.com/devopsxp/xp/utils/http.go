/*
有关Http协议GET和POST请求的封装
https://studygolang.com/articles/17373?fr=sidebar
*/
package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// 发送GET请求
// url:请求地址
// response:请求返回的内容
func Get(url string) (response []byte, err error) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return response, err
	}

	resp.Header.Set("Content-Type", "application/json")

	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return response, err
		}
	}

	response = result.Bytes()
	return
}

// 发送POST请求
// url:请求地址，data:POST请求提交的数据,contentType:请求体格式，如：application/json
// content:请求放回的内容
func Post(url string, data interface{}, contentType string) (result []byte, err error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("content-type", contentType)
	if err != nil {
		return
	}
	defer req.Body.Close()

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	result, _ = ioutil.ReadAll(resp.Body)
	return
}
