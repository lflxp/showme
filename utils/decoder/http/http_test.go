package http

import (
	"bufio"
	"bytes"
	"testing"
)

func assertEqual(t *testing.T, result interface{}, expected interface{}) {
	if result != expected {
		t.Errorf("result: %v \n no match expected: %v", result, expected)
	}
}

func TestDecodeHttpReq(t *testing.T) {
	decoder := Decoder{}
	decoder.SetFilter("")
	data := []byte("POST /test HTTP/1.1\r\nHost: google.com\r\nUser-Agent:curl\r\nContent-Length: 5\r\n\r\nHello")
	decoder.buf = bufio.NewReader(bytes.NewReader(data))
	_data, err := decoder.decodeHttp()
	if err != nil {
		t.Fatal(err)
	}
	req := _data.(*HttpReq)
	assertEqual(t, req.method, "POST")
	assertEqual(t, req.url, "/test")
	assertEqual(t, req.headers["host"], "google.com")
	assertEqual(t, req.headers["user-agent"], "curl")
	assertEqual(t, string(req.body), "Hello")
}

func TestDecodeHttpResp(t *testing.T) {
	decoder := Decoder{}
	decoder.SetFilter("")
	data := []byte("HTTP/1.1 200 OK\r\nContent-Length:11\r\nHost: google.com\r\n\r\nHello World")
	decoder.buf = bufio.NewReader(bytes.NewReader(data))
	_data, err := decoder.decodeHttp()
	if err != nil {
		t.Fatal(err)
	}
	resp := _data.(*HttpResp)
	assertEqual(t, resp.statusCode, 200)
	assertEqual(t, resp.statusMsg, "OK")
	assertEqual(t, resp.headers["host"], "google.com")
	assertEqual(t, string(resp.body), "Hello World")
}

func TestHttpReqFilter(t *testing.T) {
	decoder := Decoder{}
	decoder.SetFilter("url: /test & method: POST")
	data := []byte("POST /tes/haha HTTP/1.1\r\nHost: google.com\r\nUser-Agent:curl\r\n\r\nHello\r\n")
	decoder.buf = bufio.NewReader(bytes.NewReader(data))
	_data, err := decoder.decodeHttp()
	assertEqual(t, _data, nil)
	assertEqual(t, err, SKIP)

	// match
	data = []byte("POST /test/hahax HTTP/1.1\r\nHost: google.com\r\nUser-Agent:curl\r\n\r\nHello\r\n")
	decoder.buf = bufio.NewReader(bytes.NewReader(data))
	_data, err = decoder.decodeHttp()
	if err != nil {
		t.Fatal(err)
	}
	req := _data.(*HttpReq)
	assertEqual(t, req.url, "/test/hahax")
}

func TestHttpFilterPlainString(t *testing.T) {
	filter := NewFilter("url: /home & method: get")
	pattern, ok := filter.filters["url"]
	assertEqual(t, ok, true)
	assertEqual(t, pattern.MatchString("/home/page"), true)
	assertEqual(t, pattern.MatchString("xxhome"), false)
	pattern, ok = filter.filters["method"]
	assertEqual(t, ok, true)
	assertEqual(t, pattern.MatchString("xxget"), true)
	assertEqual(t, pattern.MatchString("et"), false)
}

func TestHttpFilterRegexp(t *testing.T) {
	filter := NewFilter("url : ^/home$ & method: PUT")
	pattern, ok := filter.filters["url"]
	assertEqual(t, ok, true)
	assertEqual(t, pattern.MatchString("/home/page"), false)
	assertEqual(t, pattern.MatchString("/home"), true)
	pattern, ok = filter.filters["method"]
	assertEqual(t, ok, true)
	assertEqual(t, pattern.MatchString("put"), true)
}
