package database

type SQLHandler interface {
	ShutdownHook()
}
