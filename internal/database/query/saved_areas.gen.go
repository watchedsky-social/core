// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/watchedsky-social/core/internal/database/models"
)

func newSavedArea(db *gorm.DB, opts ...gen.DOOption) savedArea {
	_savedArea := savedArea{}

	_savedArea.savedAreaDo.UseDB(db, opts...)
	_savedArea.savedAreaDo.UseModel(&models.SavedArea{})

	tableName := _savedArea.savedAreaDo.TableName()
	_savedArea.ALL = field.NewAsterisk(tableName)
	_savedArea.ID = field.NewString(tableName, "id")
	_savedArea.PassedZones = field.NewString(tableName, "passed_zones")
	_savedArea.Border = field.NewField(tableName, "border")
	_savedArea.CalculatedZones = field.NewString(tableName, "calculated_zones")

	_savedArea.fillFieldMap()

	return _savedArea
}

type savedArea struct {
	savedAreaDo savedAreaDo

	ALL             field.Asterisk
	ID              field.String
	PassedZones     field.String
	Border          field.Field
	CalculatedZones field.String

	fieldMap map[string]field.Expr
}

func (s savedArea) Table(newTableName string) *savedArea {
	s.savedAreaDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s savedArea) As(alias string) *savedArea {
	s.savedAreaDo.DO = *(s.savedAreaDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *savedArea) updateTableName(table string) *savedArea {
	s.ALL = field.NewAsterisk(table)
	s.ID = field.NewString(table, "id")
	s.PassedZones = field.NewString(table, "passed_zones")
	s.Border = field.NewField(table, "border")
	s.CalculatedZones = field.NewString(table, "calculated_zones")

	s.fillFieldMap()

	return s
}

func (s *savedArea) WithContext(ctx context.Context) ISavedAreaDo {
	return s.savedAreaDo.WithContext(ctx)
}

func (s savedArea) TableName() string { return s.savedAreaDo.TableName() }

func (s savedArea) Alias() string { return s.savedAreaDo.Alias() }

func (s savedArea) Columns(cols ...field.Expr) gen.Columns { return s.savedAreaDo.Columns(cols...) }

func (s *savedArea) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *savedArea) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 4)
	s.fieldMap["id"] = s.ID
	s.fieldMap["passed_zones"] = s.PassedZones
	s.fieldMap["border"] = s.Border
	s.fieldMap["calculated_zones"] = s.CalculatedZones
}

func (s savedArea) clone(db *gorm.DB) savedArea {
	s.savedAreaDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s savedArea) replaceDB(db *gorm.DB) savedArea {
	s.savedAreaDo.ReplaceDB(db)
	return s
}

type savedAreaDo struct{ gen.DO }

type ISavedAreaDo interface {
	gen.SubQuery
	Debug() ISavedAreaDo
	WithContext(ctx context.Context) ISavedAreaDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ISavedAreaDo
	WriteDB() ISavedAreaDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ISavedAreaDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ISavedAreaDo
	Not(conds ...gen.Condition) ISavedAreaDo
	Or(conds ...gen.Condition) ISavedAreaDo
	Select(conds ...field.Expr) ISavedAreaDo
	Where(conds ...gen.Condition) ISavedAreaDo
	Order(conds ...field.Expr) ISavedAreaDo
	Distinct(cols ...field.Expr) ISavedAreaDo
	Omit(cols ...field.Expr) ISavedAreaDo
	Join(table schema.Tabler, on ...field.Expr) ISavedAreaDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ISavedAreaDo
	RightJoin(table schema.Tabler, on ...field.Expr) ISavedAreaDo
	Group(cols ...field.Expr) ISavedAreaDo
	Having(conds ...gen.Condition) ISavedAreaDo
	Limit(limit int) ISavedAreaDo
	Offset(offset int) ISavedAreaDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ISavedAreaDo
	Unscoped() ISavedAreaDo
	Create(values ...*models.SavedArea) error
	CreateInBatches(values []*models.SavedArea, batchSize int) error
	Save(values ...*models.SavedArea) error
	First() (*models.SavedArea, error)
	Take() (*models.SavedArea, error)
	Last() (*models.SavedArea, error)
	Find() ([]*models.SavedArea, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.SavedArea, err error)
	FindInBatches(result *[]*models.SavedArea, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*models.SavedArea) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ISavedAreaDo
	Assign(attrs ...field.AssignExpr) ISavedAreaDo
	Joins(fields ...field.RelationField) ISavedAreaDo
	Preload(fields ...field.RelationField) ISavedAreaDo
	FirstOrInit() (*models.SavedArea, error)
	FirstOrCreate() (*models.SavedArea, error)
	FindByPage(offset int, limit int) (result []*models.SavedArea, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ISavedAreaDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	InsertOptimizedSavedArea(id string, passedZones string, calculatedZones string, border *models.Geometry) (err error)
}

// INSERT INTO saved_areas (
//
//	  id, passed_zones, calculated_zones, border
//	) VALUES (
//	  @id, @passedZones, @calculatedZones,
//	  ST_UnaryUnion(@border)
//	) ON CONFLICT(id) DO UPDATE SET
//	  passed_zones = EXCLUDED.passed_zones,
//	  calculated_zones = EXCLUDED.calculated_zones,
//	  border = EXCLUDED.border;
func (s savedAreaDo) InsertOptimizedSavedArea(id string, passedZones string, calculatedZones string, border *models.Geometry) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	params = append(params, passedZones)
	params = append(params, calculatedZones)
	params = append(params, border)
	generateSQL.WriteString("INSERT INTO saved_areas ( id, passed_zones, calculated_zones, border ) VALUES ( ?, ?, ?, ST_UnaryUnion(?) ) ON CONFLICT(id) DO UPDATE SET passed_zones = EXCLUDED.passed_zones, calculated_zones = EXCLUDED.calculated_zones, border = EXCLUDED.border; ")

	var executeSQL *gorm.DB
	executeSQL = s.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (s savedAreaDo) Debug() ISavedAreaDo {
	return s.withDO(s.DO.Debug())
}

func (s savedAreaDo) WithContext(ctx context.Context) ISavedAreaDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s savedAreaDo) ReadDB() ISavedAreaDo {
	return s.Clauses(dbresolver.Read)
}

func (s savedAreaDo) WriteDB() ISavedAreaDo {
	return s.Clauses(dbresolver.Write)
}

func (s savedAreaDo) Session(config *gorm.Session) ISavedAreaDo {
	return s.withDO(s.DO.Session(config))
}

func (s savedAreaDo) Clauses(conds ...clause.Expression) ISavedAreaDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s savedAreaDo) Returning(value interface{}, columns ...string) ISavedAreaDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s savedAreaDo) Not(conds ...gen.Condition) ISavedAreaDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s savedAreaDo) Or(conds ...gen.Condition) ISavedAreaDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s savedAreaDo) Select(conds ...field.Expr) ISavedAreaDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s savedAreaDo) Where(conds ...gen.Condition) ISavedAreaDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s savedAreaDo) Order(conds ...field.Expr) ISavedAreaDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s savedAreaDo) Distinct(cols ...field.Expr) ISavedAreaDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s savedAreaDo) Omit(cols ...field.Expr) ISavedAreaDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s savedAreaDo) Join(table schema.Tabler, on ...field.Expr) ISavedAreaDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s savedAreaDo) LeftJoin(table schema.Tabler, on ...field.Expr) ISavedAreaDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s savedAreaDo) RightJoin(table schema.Tabler, on ...field.Expr) ISavedAreaDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s savedAreaDo) Group(cols ...field.Expr) ISavedAreaDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s savedAreaDo) Having(conds ...gen.Condition) ISavedAreaDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s savedAreaDo) Limit(limit int) ISavedAreaDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s savedAreaDo) Offset(offset int) ISavedAreaDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s savedAreaDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ISavedAreaDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s savedAreaDo) Unscoped() ISavedAreaDo {
	return s.withDO(s.DO.Unscoped())
}

func (s savedAreaDo) Create(values ...*models.SavedArea) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s savedAreaDo) CreateInBatches(values []*models.SavedArea, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s savedAreaDo) Save(values ...*models.SavedArea) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s savedAreaDo) First() (*models.SavedArea, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*models.SavedArea), nil
	}
}

func (s savedAreaDo) Take() (*models.SavedArea, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*models.SavedArea), nil
	}
}

func (s savedAreaDo) Last() (*models.SavedArea, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*models.SavedArea), nil
	}
}

func (s savedAreaDo) Find() ([]*models.SavedArea, error) {
	result, err := s.DO.Find()
	return result.([]*models.SavedArea), err
}

func (s savedAreaDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.SavedArea, err error) {
	buf := make([]*models.SavedArea, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s savedAreaDo) FindInBatches(result *[]*models.SavedArea, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s savedAreaDo) Attrs(attrs ...field.AssignExpr) ISavedAreaDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s savedAreaDo) Assign(attrs ...field.AssignExpr) ISavedAreaDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s savedAreaDo) Joins(fields ...field.RelationField) ISavedAreaDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s savedAreaDo) Preload(fields ...field.RelationField) ISavedAreaDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s savedAreaDo) FirstOrInit() (*models.SavedArea, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*models.SavedArea), nil
	}
}

func (s savedAreaDo) FirstOrCreate() (*models.SavedArea, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*models.SavedArea), nil
	}
}

func (s savedAreaDo) FindByPage(offset int, limit int) (result []*models.SavedArea, count int64, err error) {
	result, err = s.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = s.Offset(-1).Limit(-1).Count()
	return
}

func (s savedAreaDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s savedAreaDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s savedAreaDo) Delete(models ...*models.SavedArea) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *savedAreaDo) withDO(do gen.Dao) *savedAreaDo {
	s.DO = *do.(*gen.DO)
	return s
}