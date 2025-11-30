package controller

type Context interface {
	Param(param string) string
	GetHeader(key string) string
	Query(key string) string
	DefaultQuery(key, val string) string
	QueryArray(key string) []string
	QueryMap(key string) map[string]string
	Set(key string, val any)
	GetString(key string) string
	Bind(any) error
	BindJSON(any) error
	Status(int)
	JSON(int, any)
	Next()
	Abort()
}
