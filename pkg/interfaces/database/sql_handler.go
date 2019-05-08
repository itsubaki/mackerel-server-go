package database

import "database/sql"

type DB interface {
	New() DB
	Close() error
	CommonDB() interface{}
	Dialect() interface{}
	Callback() interface{}
	SetLogger(log interface{})
	LogMode(enable bool) DB
	BlockGlobalUpdate(enable bool) DB
	HasBlockGlobalUpdate() bool
	SingularTable(enable bool)
	NewScope(value interface{}) interface{}
	QueryExpr() interface{}
	SubQuery() interface{}
	Where(query interface{}, args ...interface{}) DB
	Or(query interface{}, args ...interface{}) DB
	Not(query interface{}, args ...interface{}) DB
	Limit(limit interface{}) DB
	Offset(offset interface{}) DB
	Order(value interface{}, reorder ...bool) DB
	Select(query interface{}, args ...interface{}) DB
	Omit(columns ...string) DB
	Group(query string) DB
	Having(query interface{}, values ...interface{}) DB
	Joins(query string, args ...interface{}) DB
	Scopes(funcs ...func(DB) DB) DB
	Unscoped() DB
	Attrs(attrs ...interface{}) DB
	Assign(attrs ...interface{}) DB
	First(out interface{}, where ...interface{}) DB
	Take(out interface{}, where ...interface{}) DB
	Last(out interface{}, where ...interface{}) DB
	Find(out interface{}, where ...interface{}) DB
	Preloads(out interface{}) DB
	Scan(dest interface{}) DB
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	ScanRows(rows *sql.Rows, result interface{}) error
	Pluck(column string, value interface{}) DB
	Count(value interface{}) DB
	Related(value interface{}, foreignKeys ...string) DB
	FirstOrInit(out interface{}, where ...interface{}) DB
	FirstOrCreate(out interface{}, where ...interface{}) DB
	Update(attrs ...interface{}) DB
	Updates(values interface{}, ignoreProtectedAttrs ...bool) DB
	UpdateColumn(attrs ...interface{}) DB
	UpdateColumns(values interface{}) DB
	Save(value interface{}) DB
	Create(value interface{}) DB
	Delete(value interface{}, where ...interface{}) DB
	Raw(sql string, values ...interface{}) DB
	Exec(sql string, values ...interface{}) DB
	Model(value interface{}) DB
	Table(name string) DB
	Debug() DB
	Begin() DB
	Commit() DB
	Rollback() DB
	NewRecord(value interface{}) bool
	RecordNotFound() bool
	CreateTable(models ...interface{}) DB
	DropTable(values ...interface{}) DB
	DropTableIfExists(values ...interface{}) DB
	HasTable(value interface{}) bool
	AutoMigrate(values ...interface{}) DB
	ModifyColumn(column string, typ string) DB
	DropColumn(column string) DB
	AddIndex(indexName string, columns ...string) DB
	AddUniqueIndex(indexName string, columns ...string) DB
	RemoveIndex(indexName string) DB
	AddForeignKey(field string, dest string, onDelete string, onUpdate string) DB
	RemoveForeignKey(field string, dest string) DB
	Association(column string) interface{}
	Preload(column string, conditions ...interface{}) DB
	Set(name string, value interface{}) DB
	InstantSet(name string, value interface{}) DB
	Get(name string) (value interface{}, ok bool)
	SetJoinTableHandler(source interface{}, column string, handler interface{})
	AddError(err error) error
	GetErrors() []error
}
