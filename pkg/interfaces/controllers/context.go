package controllers

type Context interface {
	Param(string) string
	GetHeader(key string) string
	Query(key string) string
	DefaultQuery(key, val string) string
	QueryArray(key string) []string
	QueryMap(key string) map[string]string
	Set(key string, val interface{})
	GetString(key string) string
	Bind(interface{}) error
	BindJSON(interface{}) error
	Status(int)
	JSON(int, interface{})
	Next()
	Abort()
}
