package database

type SQLHandler interface {
	ShutdownHook()
	Close() error
}
