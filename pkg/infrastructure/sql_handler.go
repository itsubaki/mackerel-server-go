package infrastructure

import (
	"database/sql"

	"github.com/itsubaki/mackerel-api/pkg/interfaces/database"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type SQLHandler struct {
	DB *gorm.DB
}

func NewGorm() database.DB {
	db, err := gorm.Open("mysql", "root:secret@tcp(127.0.0.1:3307)/mackerel")
	if err != nil {
		panic(err)
	}

	return db
}

func (h *SQLHandler) New() database.DB {
	return h.DB.New()
}

func (h *SQLHandler) Close() error {
	return h.DB.Close()
}

func (h *SQLHandler) CommonDB() interface{} {
	return h.DB.CommonDB()
}

func (h *SQLHandler) Dialect() interface{} {
	return h.DB.Dialect()
}

func (h *SQLHandler) Callback() interface{} {
	return h.DB.Callback()
}

func (h *SQLHandler) SetLogger(log interface{}) {
	h.DB.SetLogger(log.(gorm.Logger))
}

func (h *SQLHandler) LogMode(enable bool) database.DB {
	return h.DB.LogMode(enable)
}

func (h *SQLHandler) BlockGlobalUpdate(enable bool) database.DB {
	return h.DB.BlockGlobalUpdate(enable)
}

func (h *SQLHandler) HasBlockGlobalUpdate() bool {
	return h.DB.HasBlockGlobalUpdate()
}

func (h *SQLHandler) SingularTable(enable bool) {
	h.DB.SingularTable(enable)
}

func (h *SQLHandler) NewScope(value interface{}) interface{} {
	return h.DB.NewScope(value)
}

func (h *SQLHandler) QueryExpr() interface{} {
	return h.DB.QueryExpr()
}
func (h *SQLHandler) SubQuery() interface{} {
	return h.DB.SubQuery()
}

func (h *SQLHandler) Where(query interface{}, args ...interface{}) database.DB {
	return h.DB.Where(query, args...)
}

func (h *SQLHandler) Or(query interface{}, args ...interface{}) database.DB {
	return h.DB.Or(query, args...)
}

func (h *SQLHandler) Not(query interface{}, args ...interface{}) database.DB {
	return h.DB.Not(query, args...)
}

func (h *SQLHandler) Limit(limit interface{}) database.DB {
	return h.DB.Limit(limit)
}

func (h *SQLHandler) Offset(offset interface{}) database.DB {
	return h.DB.Offset(offset)
}

func (h *SQLHandler) Order(value interface{}, reorder ...bool) database.DB {
	return h.DB.Order(value, reorder)
}

func (h *SQLHandler) Select(query interface{}, args ...interface{}) database.DB {
	return h.DB.Select(query, args...)
}

func (h *SQLHandler) Omit(columns ...string) database.DB {
	return h.DB.Omit(columns...)
}

func (h *SQLHandler) Group(query string) database.DB {
	return h.DB.Group(query)
}

func (h *SQLHandler) Having(query interface{}, values ...interface{}) database.DB {
	return h.DB.Having(query, values...)
}

func (h *SQLHandler) Joins(query string, args ...interface{}) database.DB {
	return h.DB.Joins(query, args...)
}

func (h *SQLHandler) Scopes(funcs ...func(database.DB) database.DB) database.DB {
	return h.DB.Scopes(funcs...)
}

func (h *SQLHandler) Unscoped() database.DB {
	return h.DB.Unscoped()
}

func (h *SQLHandler) Attrs(attrs ...interface{}) database.DB {
	return h.DB.Attrs(attrs...)
}

func (h *SQLHandler) Assign(attrs ...interface{}) database.DB {
	return h.DB.Assign(attrs...)
}

func (h *SQLHandler) First(out interface{}, where ...interface{}) database.DB {
	return h.DB.First(out, where...)
}

func (h *SQLHandler) Take(out interface{}, where ...interface{}) database.DB {
	return h.DB.Take(out, where...)
}

func (h *SQLHandler) Last(out interface{}, where ...interface{}) database.DB {
	return h.DB.Last(out, where...)
}

func (h *SQLHandler) Find(out interface{}, where ...interface{}) database.DB {
	return h.DB.Find(out, where)
}
func (h *SQLHandler) Preloads(out interface{}) database.DB {
	return h.DB.Preloads(out)
}
func (h *SQLHandler) Scan(dest interface{}) database.DB {
	return h.DB.Scan(dest)
}

func (h *SQLHandler) Row() *sql.Row {
	return h.DB.Row()
}
func (h *SQLHandler) Rows() (*sql.Rows, error) {
	return h.DB.Rows()
}
func (h *SQLHandler) ScanRows(rows *sql.Rows, result interface{}) error {
	return h.DB.ScanRows(rows, result)
}
func (h *SQLHandler) Pluck(column string, value interface{}) database.DB {
	return h.DB.Pluck(column, value)
}
func (h *SQLHandler) Count(value interface{}) database.DB {
	return h.DB.Count(value)
}
func (h *SQLHandler) Related(value interface{}, foreignKeys ...string) database.DB {
	return h.DB.Related(value, foreignKeys...)
}
func (h *SQLHandler) FirstOrInit(out interface{}, where ...interface{}) database.DB {
	return h.DB.FirstOrInit(out, where...)
}

func (h *SQLHandler) FirstOrCreate(out interface{}, where ...interface{}) database.DB {
	return h.DB.FirstOrCreate(out, where...)
}

func (h *SQLHandler) Update(attrs ...interface{}) database.DB {
	return h.DB.Update(attrs...)
}

func (h *SQLHandler) Updates(values interface{}, ignoreProtectedAttrs ...bool) database.DB {
	return h.DB.Updates(values, ignoreProtectedAttrs...)
}

func (h *SQLHandler) UpdateColumn(attrs ...interface{}) database.DB {
	return h.DB.UpdateColumn(attrs...)
}

func (h *SQLHandler) UpdateColumns(values interface{}) database.DB {
	return h.DB.UpdateColumns(values)
}

func (h *SQLHandler) Save(value interface{}) database.DB {
	return h.DB.Save(value)
}

func (h *SQLHandler) Create(value interface{}) database.DB {
	return h.DB.Create(value)
}
func (h *SQLHandler) Delete(value interface{}, where ...interface{}) database.DB {
	return h.DB.Delete(value, where...)
}

func (h *SQLHandler) Raw(sql string, values ...interface{}) database.DB {
	return h.DB.Raw(sql, values...)
}

func (h *SQLHandler) Exec(sql string, values ...interface{}) database.DB {
	return h.DB.Exec(sql, values...)
}

func (h *SQLHandler) Model(value interface{}) database.DB {
	return h.DB.Model(value)
}

func (h *SQLHandler) Table(name string) database.DB {
	return h.DB.Table(name)
}

func (h *SQLHandler) Debug() database.DB {
	return h.DB.Debug()
}

func (h *SQLHandler) Begin() database.DB {
	return h.DB.Begin()
}

func (h *SQLHandler) Commit() database.DB {
	return h.DB.Commit()
}

func (h *SQLHandler) Rollback() database.DB {
	return h.DB.Rollback()
}

func (h *SQLHandler) NewRecord(value interface{}) bool {
	return h.DB.NewRecord(value)
}

func (h *SQLHandler) RecordNotFound() bool {
	return h.DB.RecordNotFound()
}

func (h *SQLHandler) CreateTable(models ...interface{}) database.DB {
	return h.DB.CreateTable(models...)
}

func (h *SQLHandler) DropTable(values ...interface{}) database.DB {
	return h.DB.DropTable(values...)
}

func (h *SQLHandler) DropTableIfExists(values ...interface{}) database.DB {
	return h.DB.DropTableIfExists(values...)
}

func (h *SQLHandler) HasTable(value interface{}) bool {
	return h.DB.HasTable(value)
}

func (h *SQLHandler) AutoMigrate(values ...interface{}) database.DB {
	return h.DB.AutoMigrate(values...)
}

func (h *SQLHandler) ModifyColumn(column string, typ string) database.DB {
	return h.DB.ModifyColumn(column, typ)
}

func (h *SQLHandler) DropColumn(column string) database.DB {
	return h.DB.DropColumn(column)
}

func (h *SQLHandler) AddIndex(indexName string, columns ...string) database.DB {
	return h.DB.AddIndex(indexName, columns...)
}

func (h *SQLHandler) AddUniqueIndex(indexName string, columns ...string) database.DB {
	return h.DB.AddUniqueIndex(indexName, columns...)
}

func (h *SQLHandler) RemoveIndex(indexName string) database.DB {
	return h.DB.RemoveIndex(indexName)
}

func (h *SQLHandler) AddForeignKey(field string, dest string, onDelete string, onUpdate string) database.DB {
	return h.DB.AddForeignKey(field, dest, onDelete, onUpdate)
}

func (h *SQLHandler) RemoveForeignKey(field string, dest string) database.DB {
	return h.DB.RemoveForeignKey(field, dest)
}

func (h *SQLHandler) Association(column string) interface{} {
	return h.DB.Association(column)
}

func (h *SQLHandler) Preload(column string, conditions ...interface{}) database.DB {
	return h.DB.Preload(column, conditions...)
}

func (h *SQLHandler) Set(name string, value interface{}) database.DB {
	return h.DB.Set(name, value)
}

func (h *SQLHandler) InstantSet(name string, value interface{}) database.DB {
	return h.DB.InstantSet(name, value)
}

func (h *SQLHandler) Get(name string) (value interface{}, ok bool) {
	return h.DB.Get(name)
}

func (h *SQLHandler) SetJoinTableHandler(source interface{}, column string, handler interface{}) {
	return h.DB.SetJoinTableHandler(source, column, handler)
}

func (h *SQLHandler) AddError(err error) error {
	return h.DB.AddError(err)
}

func (h *SQLHandler) GetErrors() []error {
	return h.DB.GetErrors()
}
