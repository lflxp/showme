package http

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/juju/errors"
	"github.com/lflxp/showme/utils/decoder"
)

var SKIP = errors.New("Skip msg")

type Decoder struct {
	Buf    *bufio.Reader
	filter *Filter
}

func (d *Decoder) Decode(reader io.Reader, writer io.Writer, opts *decoder.Options) error {
	d.Buf = bufio.NewReader(reader)
	for {
		msg, err := d.DecodeHttp()
		if err != nil {
			if err == SKIP {
				continue
			}
			log.Println(err)
			continue
		}
		writer.Write([]byte(msg.StringHeader()))
		if opts.DeepDecode {
			_msg, err := msg.DecodeBody()
			if err != nil {
				log.Println(err)
				continue
			}
			writer.Write([]byte(_msg))
		} else {
			writer.Write([]byte(msg.RawBody()))
		}
		writer.Write([]byte("\n"))
	}
	return nil
}

func (d *Decoder) DecodeHttp() (Http, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	isReq := true
	firstLine, err := d.Buf.ReadString('\n')
	if err != nil {
		return nil, err
	}
	firstLine = firstLine[:len(firstLine)-2]
	f := strings.SplitN(firstLine, " ", 3)
	fLen := len(f)
	if fLen < 3 {
		return nil, errors.New("bad http msg: " + firstLine)
	} else if fLen == 3 {
		if f[0][:2] == "HT" { // eg: HTTP/1.1 200 OK
			isReq = false
		}
	} else {
		isReq = false
	}
	if isReq {
		req := new(HttpReq)
		req.Method = f[0]
		req.Url = f[1]
		req.Version = f[2]
		req.Headers, err = parseHeaders(d.Buf)
		if err != nil {
			return nil, err
		}
		req.body = parseBody(req.Headers, d.Buf)
		if !d.filter.IsEmpty() && !req.Match(d.filter) {
			return nil, SKIP
		}
		return req, nil
	} else {
		// it's http response
		resp := new(HttpResp)
		resp.Version = f[0]
		resp.StatusCode, err = strconv.Atoi(f[1])
		if err != nil {
			return nil, errors.New("Invalid http resp: " + firstLine)
		}
		resp.StatusMsg = f[2]
		resp.Headers, err = parseHeaders(d.Buf)
		if err != nil {
			return nil, err
		}
		resp.body = parseBody(resp.Headers, d.Buf)
		if !d.filter.IsEmpty() && !resp.Match(d.filter) {
			return nil, SKIP
		}
		return resp, nil
	}
}

func parseHeaders(buf *bufio.Reader) (map[string]string, error) {
	var err error
	var line string
	headers := make(map[string]string)
	for {
		// parse headers
		if line, err = buf.ReadString('\n'); err != nil {
			return nil, err
		}
		if line == "\r\n" {
			// end of header
			break
		}
		line = line[:len(line)-2] // remove \r\n
		result := strings.SplitN(line, ":", 2)
		headers[strings.ToLower(strings.TrimSpace(result[0]))] = strings.TrimSpace(result[1])
	}
	return headers, nil
}

func parseBody(headers map[string]string, reader *bufio.Reader) []byte {
	length, ok := headers["content-length"]
	if ok {
		bodyLen, _ := strconv.Atoi(length)
		body := make([]byte, bodyLen)
		reader.Read(body)
		return body
	}
	return nil
}

func (d *Decoder) SetFilter(filter string) {
	d.filter = NewFilter(filter)
}

func init() {
	decoder.Register("http", new(Decoder))
}
