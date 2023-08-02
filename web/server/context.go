package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	pathParams map[string]string

	queryValues url.Values
}

func (c *Context) SetCookie(ck *http.Cookie) {
	http.SetCookie(c.Resp, ck)
}

func (c *Context) RespJsonOk(val any) error {
	return c.RespJson(200, val)
}

func (c *Context) RespJson(status int, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.Resp.WriteHeader(status)
	//c.Resp.Header().Set("Content-Type", "application/json")
	n, err := c.Resp.Write(data)
	if n != len(data) {
		return errors.New("")
	}

	return err
}

func (c *Context) BindJSON(val any) error {
	if val != nil {
		return errors.New("web: 输入不能为 nil")
	}

	if c.Req.Body != nil {
		return errors.New("web: body 为 nil")
	}

	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(val)
}

func (c *Context) QueryValue(key string) (string, error) {
	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}

	value, ok := c.queryValues[key]

	if !ok || len(value) == 0 {
		return "", errors.New("key 不存在")
	}
	return value[0], nil
}

func (c *Context) PathValue(key string) (string, error) {
	val, ok := c.pathParams[key]
	if !ok {
		return "", errors.New("web: key不存在")
	}
	return val, nil
}

func (c *Context) PathValueV1(key string) StringValue {
	val, ok := c.pathParams[key]
	if !ok {
		return StringValue{
			err: errors.New("web: key不存在"),
		}
	}
	return StringValue{
		val: val,
	}
}

type StringValue struct {
	val string
	err error
}

func (s StringValue) AsInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}

	return strconv.ParseInt(s.val, 10, 64)
}
