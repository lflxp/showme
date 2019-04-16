package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/ugorji/go/codec"
)

var mh codec.MsgpackHandle

type Http interface {
	Match(*Filter) bool
	StringHeader() string
	DecodeBody() (string, error)
	RawBody() []byte
}

func prettyPrint(v map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range v {
		switch reflect.ValueOf(value).Kind() {
		case reflect.Map:
			result[key] = _prettyMap(value.(map[interface{}]interface{}))
		case reflect.Slice:
			result[key] = _prettySlice(value)
		default:
			result[key] = value
		}
	}
	return result
}

func _prettyMap(m map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		switch reflect.ValueOf(v).Kind() {
		case reflect.Map:
			result[k.(string)] = _prettyMap(v.(map[interface{}]interface{}))
		case reflect.Slice:
			if data, ok := v.([]byte); ok {
				result[k.(string)] = string(data)
			} else {
				result[k.(string)] = _prettySlice(v)
			}
		default:
			switch k.(type) {
			case int:
				result[strconv.Itoa(k.(int))] = v
			case int8:
				result[strconv.Itoa(int(k.(int8)))] = v
			case int16:
				result[strconv.Itoa(int(k.(int16)))] = v
			case int32:
				result[strconv.Itoa(int(k.(int32)))] = v
			case int64:
				result[strconv.Itoa(int(k.(int64)))] = v
			default:
				result[k.(string)] = v
			}
		}
	}
	return result
}

func _prettySlice(s interface{}) []interface{} {
	result := make([]interface{}, 0)
	rs := reflect.ValueOf(s)
	for i := 0; i < rs.Len(); i++ {
		v := rs.Index(i).Interface()
		switch reflect.ValueOf(v).Kind() {
		case reflect.Map:
			result = append(result, _prettyMap(v.(map[interface{}]interface{})))
		case reflect.Slice:
			if data, ok := v.([]byte); ok {
				result = append(result, string(data))
			} else {
				result = append(result, _prettySlice(v))
			}
		default:
			result = append(result, v)
		}
	}
	return result
}

func decodeToString(contentType string, data []byte) (string, error) {
	var v map[string]interface{}
	switch contentType {
	case "application/msgpack":
		dec := codec.NewDecoder(bytes.NewReader(data), &mh)
		err := dec.Decode(&v)
		if err != nil {
			return "", err
		}
		pv, err := json.MarshalIndent(prettyPrint(v), "", " ")
		if err != nil {
			return "", err
		}
		return fmt.Sprint(string(pv)), nil
	}
	return string(data), nil
}

type HttpReq struct {
	Method  string
	Url     string
	Version string
	Headers map[string]string
	body    []byte
}

func (m *HttpReq) RawBody() []byte {
	return m.body
}

func (m *HttpReq) DecodeBody() (string, error) {
	return decodeToString(m.Headers["content-type"], m.body)
}

func (m *HttpReq) StringHeader() string {
	headStr := ""
	for k, v := range m.Headers {
		headStr += k + ": " + v + "\r\n"
	}
	return fmt.Sprintf("%s %s %s\r\n%s\r\n", m.Method, m.Url, m.Version, headStr)
}

func (m *HttpReq) Match(filter *Filter) bool {
	filters := make(map[string]*regexp.Regexp)
	mapCopy(filters, filter.filters)
	if _, ok := filters["method"]; ok && !filters["method"].MatchString(m.Method) {
		return false
	}
	delete(filters, "method")
	if _, ok := filters["url"]; ok && !filters["url"].MatchString(m.Url) {
		return false
	}
	delete(filters, "url")
	if _, ok := filters["version"]; ok && !filters["version"].MatchString(m.Version) {
		return false
	}
	delete(filters, "version")
	if _, ok := filters["body"]; ok && !matchBody(filters["body"], m.body) {
		return false
	}
	delete(filters, "body")
	if len(filters) > 0 && !matchHeaders(filters, m.Headers) {
		return false
	}
	return true
}

type HttpResp struct {
	Version    string
	StatusCode int
	StatusMsg  string
	Headers    map[string]string
	body       []byte
}

func (m *HttpResp) RawBody() []byte {
	return m.body
}

func (m *HttpResp) DecodeBody() (string, error) {
	return decodeToString(m.Headers["content-type"], m.body)
}

func (m *HttpResp) StringHeader() string {
	headStr := ""
	for k, v := range m.Headers {
		headStr += k + ": " + v + "\r\n"
	}
	return fmt.Sprintf("%s %s %s\r\n%s\r\n", m.Version, m.StatusCode, m.StatusMsg, headStr)
}

func (m *HttpResp) Match(filter *Filter) bool {
	filters := make(map[string]*regexp.Regexp)
	mapCopy(filters, filter.filters)
	if _, ok := filters["version"]; ok && !filters["version"].MatchString(m.Version) {
		return false
	}
	delete(filters, "version")
	if _, ok := filters["statusCode"]; ok && !filters["statusCode"].MatchString(strconv.Itoa(m.StatusCode)) {
		return false
	}
	delete(filters, "statusCode")
	if _, ok := filters["statusMsg"]; ok && !filters["statusMsg"].MatchString(m.StatusMsg) {
		return false
	}
	delete(filters, "statusMsg")
	if _, ok := filters["body"]; ok && !matchBody(filters["body"], m.body) {
		return false
	}
	delete(filters, "body")
	if len(filters) > 0 && !matchHeaders(filters, m.Headers) {
		return false
	}
	return true
}

func mapCopy(dst, src interface{}) {
	dv, sv := reflect.ValueOf(dst), reflect.ValueOf(src)

	for _, k := range sv.MapKeys() {
		dv.SetMapIndex(k, sv.MapIndex(k))
	}
}

func matchHeaders(rules map[string]*regexp.Regexp, Headers map[string]string) bool {
	for h, pattern := range rules {
		value, ok := Headers[h]
		if !ok {
			// header not exist, so not match
			return false
		}
		if !pattern.MatchString(value) {
			return false
		}
	}
	return true
}

func matchBody(pattern *regexp.Regexp, body []byte) bool {
	if !pattern.Match(body) {
		return false
	}
	return true
}
