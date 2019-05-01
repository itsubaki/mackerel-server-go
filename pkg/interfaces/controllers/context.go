package controllers

type Context interface {
	Param(string) string
	Query(key string) string
	QueryArray(key string) []string
	QueryMap(key string) map[string]string
	Bind(interface{}) error
	BindJSON(interface{}) error
	Status(int)
	JSON(int, interface{})
}
