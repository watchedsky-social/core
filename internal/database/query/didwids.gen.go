// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/watchedsky-social/core/internal/database/models"
)

func newDidwid(db *gorm.DB, opts ...gen.DOOption) didwid {
	_didwid := didwid{}

	_didwid.didwidDo.UseDB(db, opts...)
	_didwid.didwidDo.UseModel(&models.Didwid{})

	tableName := _didwid.didwidDo.TableName()
	_didwid.ALL = field.NewAsterisk(tableName)
	_didwid.Did = field.NewString(tableName, "did")
	_didwid.Wid = field.NewString(tableName, "wid")
	_didwid.CreatedAt = field.NewTime(tableName, "created_at")

	_didwid.fillFieldMap()

	return _didwid
}

type didwid struct {
	didwidDo didwidDo

	ALL       field.Asterisk
	Did       field.String
	Wid       field.String
	CreatedAt field.Time

	fieldMap map[string]field.Expr
}

func (d didwid) Table(newTableName string) *didwid {
	d.didwidDo.UseTable(newTableName)
	return d.updateTableName(newTableName)
}

func (d didwid) As(alias string) *didwid {
	d.didwidDo.DO = *(d.didwidDo.As(alias).(*gen.DO))
	return d.updateTableName(alias)
}

func (d *didwid) updateTableName(table string) *didwid {
	d.ALL = field.NewAsterisk(table)
	d.Did = field.NewString(table, "did")
	d.Wid = field.NewString(table, "wid")
	d.CreatedAt = field.NewTime(table, "created_at")

	d.fillFieldMap()

	return d
}

func (d *didwid) WithContext(ctx context.Context) IDidwidDo { return d.didwidDo.WithContext(ctx) }

func (d didwid) TableName() string { return d.didwidDo.TableName() }

func (d didwid) Alias() string { return d.didwidDo.Alias() }

func (d didwid) Columns(cols ...field.Expr) gen.Columns { return d.didwidDo.Columns(cols...) }

func (d *didwid) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := d.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (d *didwid) fillFieldMap() {
	d.fieldMap = make(map[string]field.Expr, 3)
	d.fieldMap["did"] = d.Did
	d.fieldMap["wid"] = d.Wid
	d.fieldMap["created_at"] = d.CreatedAt
}

func (d didwid) clone(db *gorm.DB) didwid {
	d.didwidDo.ReplaceConnPool(db.Statement.ConnPool)
	return d
}

func (d didwid) replaceDB(db *gorm.DB) didwid {
	d.didwidDo.ReplaceDB(db)
	return d
}

type didwidDo struct{ gen.DO }

type IDidwidDo interface {
	gen.SubQuery
	Debug() IDidwidDo
	WithContext(ctx context.Context) IDidwidDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IDidwidDo
	WriteDB() IDidwidDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IDidwidDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IDidwidDo
	Not(conds ...gen.Condition) IDidwidDo
	Or(conds ...gen.Condition) IDidwidDo
	Select(conds ...field.Expr) IDidwidDo
	Where(conds ...gen.Condition) IDidwidDo
	Order(conds ...field.Expr) IDidwidDo
	Distinct(cols ...field.Expr) IDidwidDo
	Omit(cols ...field.Expr) IDidwidDo
	Join(table schema.Tabler, on ...field.Expr) IDidwidDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IDidwidDo
	RightJoin(table schema.Tabler, on ...field.Expr) IDidwidDo
	Group(cols ...field.Expr) IDidwidDo
	Having(conds ...gen.Condition) IDidwidDo
	Limit(limit int) IDidwidDo
	Offset(offset int) IDidwidDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IDidwidDo
	Unscoped() IDidwidDo
	Create(values ...*models.Didwid) error
	CreateInBatches(values []*models.Didwid, batchSize int) error
	Save(values ...*models.Didwid) error
	First() (*models.Didwid, error)
	Take() (*models.Didwid, error)
	Last() (*models.Didwid, error)
	Find() ([]*models.Didwid, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.Didwid, err error)
	FindInBatches(result *[]*models.Didwid, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*models.Didwid) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IDidwidDo
	Assign(attrs ...field.AssignExpr) IDidwidDo
	Joins(fields ...field.RelationField) IDidwidDo
	Preload(fields ...field.RelationField) IDidwidDo
	FirstOrInit() (*models.Didwid, error)
	FirstOrCreate() (*models.Didwid, error)
	FindByPage(offset int, limit int) (result []*models.Didwid, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IDidwidDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (d didwidDo) Debug() IDidwidDo {
	return d.withDO(d.DO.Debug())
}

func (d didwidDo) WithContext(ctx context.Context) IDidwidDo {
	return d.withDO(d.DO.WithContext(ctx))
}

func (d didwidDo) ReadDB() IDidwidDo {
	return d.Clauses(dbresolver.Read)
}

func (d didwidDo) WriteDB() IDidwidDo {
	return d.Clauses(dbresolver.Write)
}

func (d didwidDo) Session(config *gorm.Session) IDidwidDo {
	return d.withDO(d.DO.Session(config))
}

func (d didwidDo) Clauses(conds ...clause.Expression) IDidwidDo {
	return d.withDO(d.DO.Clauses(conds...))
}

func (d didwidDo) Returning(value interface{}, columns ...string) IDidwidDo {
	return d.withDO(d.DO.Returning(value, columns...))
}

func (d didwidDo) Not(conds ...gen.Condition) IDidwidDo {
	return d.withDO(d.DO.Not(conds...))
}

func (d didwidDo) Or(conds ...gen.Condition) IDidwidDo {
	return d.withDO(d.DO.Or(conds...))
}

func (d didwidDo) Select(conds ...field.Expr) IDidwidDo {
	return d.withDO(d.DO.Select(conds...))
}

func (d didwidDo) Where(conds ...gen.Condition) IDidwidDo {
	return d.withDO(d.DO.Where(conds...))
}

func (d didwidDo) Order(conds ...field.Expr) IDidwidDo {
	return d.withDO(d.DO.Order(conds...))
}

func (d didwidDo) Distinct(cols ...field.Expr) IDidwidDo {
	return d.withDO(d.DO.Distinct(cols...))
}

func (d didwidDo) Omit(cols ...field.Expr) IDidwidDo {
	return d.withDO(d.DO.Omit(cols...))
}

func (d didwidDo) Join(table schema.Tabler, on ...field.Expr) IDidwidDo {
	return d.withDO(d.DO.Join(table, on...))
}

func (d didwidDo) LeftJoin(table schema.Tabler, on ...field.Expr) IDidwidDo {
	return d.withDO(d.DO.LeftJoin(table, on...))
}

func (d didwidDo) RightJoin(table schema.Tabler, on ...field.Expr) IDidwidDo {
	return d.withDO(d.DO.RightJoin(table, on...))
}

func (d didwidDo) Group(cols ...field.Expr) IDidwidDo {
	return d.withDO(d.DO.Group(cols...))
}

func (d didwidDo) Having(conds ...gen.Condition) IDidwidDo {
	return d.withDO(d.DO.Having(conds...))
}

func (d didwidDo) Limit(limit int) IDidwidDo {
	return d.withDO(d.DO.Limit(limit))
}

func (d didwidDo) Offset(offset int) IDidwidDo {
	return d.withDO(d.DO.Offset(offset))
}

func (d didwidDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IDidwidDo {
	return d.withDO(d.DO.Scopes(funcs...))
}

func (d didwidDo) Unscoped() IDidwidDo {
	return d.withDO(d.DO.Unscoped())
}

func (d didwidDo) Create(values ...*models.Didwid) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Create(values)
}

func (d didwidDo) CreateInBatches(values []*models.Didwid, batchSize int) error {
	return d.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (d didwidDo) Save(values ...*models.Didwid) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Save(values)
}

func (d didwidDo) First() (*models.Didwid, error) {
	if result, err := d.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*models.Didwid), nil
	}
}

func (d didwidDo) Take() (*models.Didwid, error) {
	if result, err := d.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*models.Didwid), nil
	}
}

func (d didwidDo) Last() (*models.Didwid, error) {
	if result, err := d.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*models.Didwid), nil
	}
}

func (d didwidDo) Find() ([]*models.Didwid, error) {
	result, err := d.DO.Find()
	return result.([]*models.Didwid), err
}

func (d didwidDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.Didwid, err error) {
	buf := make([]*models.Didwid, 0, batchSize)
	err = d.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (d didwidDo) FindInBatches(result *[]*models.Didwid, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return d.DO.FindInBatches(result, batchSize, fc)
}

func (d didwidDo) Attrs(attrs ...field.AssignExpr) IDidwidDo {
	return d.withDO(d.DO.Attrs(attrs...))
}

func (d didwidDo) Assign(attrs ...field.AssignExpr) IDidwidDo {
	return d.withDO(d.DO.Assign(attrs...))
}

func (d didwidDo) Joins(fields ...field.RelationField) IDidwidDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Joins(_f))
	}
	return &d
}

func (d didwidDo) Preload(fields ...field.RelationField) IDidwidDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Preload(_f))
	}
	return &d
}

func (d didwidDo) FirstOrInit() (*models.Didwid, error) {
	if result, err := d.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*models.Didwid), nil
	}
}

func (d didwidDo) FirstOrCreate() (*models.Didwid, error) {
	if result, err := d.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*models.Didwid), nil
	}
}

func (d didwidDo) FindByPage(offset int, limit int) (result []*models.Didwid, count int64, err error) {
	result, err = d.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = d.Offset(-1).Limit(-1).Count()
	return
}

func (d didwidDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = d.Count()
	if err != nil {
		return
	}

	err = d.Offset(offset).Limit(limit).Scan(result)
	return
}

func (d didwidDo) Scan(result interface{}) (err error) {
	return d.DO.Scan(result)
}

func (d didwidDo) Delete(models ...*models.Didwid) (result gen.ResultInfo, err error) {
	return d.DO.Delete(models)
}

func (d *didwidDo) withDO(do gen.Dao) *didwidDo {
	d.DO = *do.(*gen.DO)
	return d
}
