package controller_test

import (
	"encoding/json"
	"fmt"
)

type ContextMock struct {
	header  map[string]string
	param   map[string]string
	ctx     map[string]any
	reqbody []byte
	resbody any
	status  int
}

func Context() *ContextMock {
	return &ContextMock{
		header: make(map[string]string),
		param:  make(map[string]string),
		ctx:    make(map[string]any),
		status: -1,
	}
}

func (c *ContextMock) SetRequestBody(b []byte) {
	c.reqbody = b
}

func (c *ContextMock) ResponseBody() any {
	return c.resbody
}

func (c *ContextMock) SetParam(key, value string) {
	c.param[key] = value
}

func (c *ContextMock) Param(key string) string {
	v, ok := c.param[key]
	if ok {
		return v
	}

	return ""
}

func (c *ContextMock) SetHeader(key, value string) {
	c.header[key] = value
}

func (c *ContextMock) GetHeader(key string) string {
	v, ok := c.header[key]
	if ok {
		return v
	}

	return ""
}

func (c *ContextMock) Query(key string) string {
	return ""
}

func (c *ContextMock) DefaultQuery(key, val string) string {
	return ""
}

func (c *ContextMock) QueryArray(key string) []string {
	return []string{}
}

func (c *ContextMock) QueryMap(key string) map[string]string {
	return make(map[string]string)
}

func (c *ContextMock) Set(key string, val any) {
	c.ctx[key] = val
}

func (c *ContextMock) GetString(key string) string {
	v, ok := c.ctx[key]
	if ok {
		return v.(string)
	}

	return ""
}

func (c *ContextMock) Bind(any) error {
	return nil
}

func (c *ContextMock) BindJSON(b any) error {
	if err := json.Unmarshal(c.reqbody, &b); err != nil {
		return fmt.Errorf("unmarshal: %v", err)
	}

	return nil
}

func (c *ContextMock) Status(status int) {
	c.status = status
}

func (c *ContextMock) GetStatus() int {
	return c.status
}

func (c *ContextMock) JSON(status int, body any) {
	c.Status(status)
	c.resbody = body
}

func (c *ContextMock) Next() {

}

func (c *ContextMock) Abort() {

}
