package controller_test

type ContextMock struct {
	header map[string]string
}

func Context() *ContextMock {
	return &ContextMock{
		header: make(map[string]string),
	}
}

func (c *ContextMock) Param(param string) string {
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

func (c *ContextMock) Set(key string, val interface{}) {

}

func (c *ContextMock) GetString(key string) string {
	return ""
}

func (c *ContextMock) Bind(interface{}) error {
	return nil
}

func (c *ContextMock) BindJSON(interface{}) error {
	return nil
}

func (c *ContextMock) Status(int) {

}

func (c *ContextMock) JSON(int, interface{}) {

}

func (c *ContextMock) Next() {

}

func (c *ContextMock) Abort() {

}
