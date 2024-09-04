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

	"time"
)

func newAlert(db *gorm.DB, opts ...gen.DOOption) alert {
	_alert := alert{}

	_alert.alertDo.UseDB(db, opts...)
	_alert.alertDo.UseModel(&models.Alert{})

	tableName := _alert.alertDo.TableName()
	_alert.ALL = field.NewAsterisk(tableName)
	_alert.ID = field.NewString(tableName, "id")
	_alert.AreaDesc = field.NewString(tableName, "area_desc")
	_alert.Headline = field.NewString(tableName, "headline")
	_alert.Description = field.NewString(tableName, "description")
	_alert.Severity = field.NewString(tableName, "severity")
	_alert.Certainty = field.NewString(tableName, "certainty")
	_alert.Urgency = field.NewString(tableName, "urgency")
	_alert.Event = field.NewString(tableName, "event")
	_alert.Sent = field.NewTime(tableName, "sent")
	_alert.Effective = field.NewTime(tableName, "effective")
	_alert.Onset = field.NewTime(tableName, "onset")
	_alert.Expires = field.NewTime(tableName, "expires")
	_alert.Ends = field.NewTime(tableName, "ends")
	_alert.ReferenceIds = field.NewField(tableName, "reference_ids")
	_alert.Border = field.NewField(tableName, "border")
	_alert.MessageType = field.NewString(tableName, "message_type")
	_alert.SkeetInfo = field.NewField(tableName, "skeet_info")

	_alert.fillFieldMap()

	return _alert
}

type alert struct {
	alertDo alertDo

	ALL          field.Asterisk
	ID           field.String
	AreaDesc     field.String
	Headline     field.String
	Description  field.String
	Severity     field.String
	Certainty    field.String
	Urgency      field.String
	Event        field.String
	Sent         field.Time
	Effective    field.Time
	Onset        field.Time
	Expires      field.Time
	Ends         field.Time
	ReferenceIds field.Field
	Border       field.Field
	MessageType  field.String
	SkeetInfo    field.Field

	fieldMap map[string]field.Expr
}

func (a alert) Table(newTableName string) *alert {
	a.alertDo.UseTable(newTableName)
	return a.updateTableName(newTableName)
}

func (a alert) As(alias string) *alert {
	a.alertDo.DO = *(a.alertDo.As(alias).(*gen.DO))
	return a.updateTableName(alias)
}

func (a *alert) updateTableName(table string) *alert {
	a.ALL = field.NewAsterisk(table)
	a.ID = field.NewString(table, "id")
	a.AreaDesc = field.NewString(table, "area_desc")
	a.Headline = field.NewString(table, "headline")
	a.Description = field.NewString(table, "description")
	a.Severity = field.NewString(table, "severity")
	a.Certainty = field.NewString(table, "certainty")
	a.Urgency = field.NewString(table, "urgency")
	a.Event = field.NewString(table, "event")
	a.Sent = field.NewTime(table, "sent")
	a.Effective = field.NewTime(table, "effective")
	a.Onset = field.NewTime(table, "onset")
	a.Expires = field.NewTime(table, "expires")
	a.Ends = field.NewTime(table, "ends")
	a.ReferenceIds = field.NewField(table, "reference_ids")
	a.Border = field.NewField(table, "border")
	a.MessageType = field.NewString(table, "message_type")
	a.SkeetInfo = field.NewField(table, "skeet_info")

	a.fillFieldMap()

	return a
}

func (a *alert) WithContext(ctx context.Context) IAlertDo { return a.alertDo.WithContext(ctx) }

func (a alert) TableName() string { return a.alertDo.TableName() }

func (a alert) Alias() string { return a.alertDo.Alias() }

func (a alert) Columns(cols ...field.Expr) gen.Columns { return a.alertDo.Columns(cols...) }

func (a *alert) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (a *alert) fillFieldMap() {
	a.fieldMap = make(map[string]field.Expr, 17)
	a.fieldMap["id"] = a.ID
	a.fieldMap["area_desc"] = a.AreaDesc
	a.fieldMap["headline"] = a.Headline
	a.fieldMap["description"] = a.Description
	a.fieldMap["severity"] = a.Severity
	a.fieldMap["certainty"] = a.Certainty
	a.fieldMap["urgency"] = a.Urgency
	a.fieldMap["event"] = a.Event
	a.fieldMap["sent"] = a.Sent
	a.fieldMap["effective"] = a.Effective
	a.fieldMap["onset"] = a.Onset
	a.fieldMap["expires"] = a.Expires
	a.fieldMap["ends"] = a.Ends
	a.fieldMap["reference_ids"] = a.ReferenceIds
	a.fieldMap["border"] = a.Border
	a.fieldMap["message_type"] = a.MessageType
	a.fieldMap["skeet_info"] = a.SkeetInfo
}

func (a alert) clone(db *gorm.DB) alert {
	a.alertDo.ReplaceConnPool(db.Statement.ConnPool)
	return a
}

func (a alert) replaceDB(db *gorm.DB) alert {
	a.alertDo.ReplaceDB(db)
	return a
}

type alertDo struct{ gen.DO }

type IAlertDo interface {
	gen.SubQuery
	Debug() IAlertDo
	WithContext(ctx context.Context) IAlertDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IAlertDo
	WriteDB() IAlertDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IAlertDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IAlertDo
	Not(conds ...gen.Condition) IAlertDo
	Or(conds ...gen.Condition) IAlertDo
	Select(conds ...field.Expr) IAlertDo
	Where(conds ...gen.Condition) IAlertDo
	Order(conds ...field.Expr) IAlertDo
	Distinct(cols ...field.Expr) IAlertDo
	Omit(cols ...field.Expr) IAlertDo
	Join(table schema.Tabler, on ...field.Expr) IAlertDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IAlertDo
	RightJoin(table schema.Tabler, on ...field.Expr) IAlertDo
	Group(cols ...field.Expr) IAlertDo
	Having(conds ...gen.Condition) IAlertDo
	Limit(limit int) IAlertDo
	Offset(offset int) IAlertDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IAlertDo
	Unscoped() IAlertDo
	Create(values ...*models.Alert) error
	CreateInBatches(values []*models.Alert, batchSize int) error
	Save(values ...*models.Alert) error
	First() (*models.Alert, error)
	Take() (*models.Alert, error)
	Last() (*models.Alert, error)
	Find() ([]*models.Alert, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.Alert, err error)
	FindInBatches(result *[]*models.Alert, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*models.Alert) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IAlertDo
	Assign(attrs ...field.AssignExpr) IAlertDo
	Joins(fields ...field.RelationField) IAlertDo
	Preload(fields ...field.RelationField) IAlertDo
	FirstOrInit() (*models.Alert, error)
	FirstOrCreate() (*models.Alert, error)
	FindByPage(offset int, limit int) (result []*models.Alert, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IAlertDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	InsertOptimizedAlert(id string, areaDesc string, headline string, description string, severity *string, certainty *string, urgency *string, event *string, sent time.Time, effective time.Time, onset *time.Time, expires *time.Time, ends *time.Time, referenceIDs *models.StringSlice, border *models.Geometry, messageType *string) (err error)
	GetCustomAlertURIs(watchID string, limit uint) (result []*models.Alert, err error)
	GetCustomAlertURIsWithCursor(watchID string, limit uint, cursor uint) (result []*models.Alert, err error)
}

// INSERT INTO alerts (
//
//	id, area_desc, headline, description, severity, certainty, urgency, event,
//	sent, effective, onset, expires, ends, reference_ids, border, message_type
//	) VALUES (
//	@id, @areaDesc, @headline, @description, @severity, @certainty,
//	@urgency, @event, @sent, @effective, @onset, @expires, @ends,
//	@referenceIDs, ST_UnaryUnion(@border),
//	@messageType
//	) ON CONFLICT(id) DO UPDATE SET
//	area_desc = EXCLUDED.area_desc,
//	headline = EXCLUDED.headline,
//	description = EXCLUDED.description,
//	severity = EXCLUDED.severity,
//	certainty = EXCLUDED.certainty,
//	urgency = EXCLUDED.urgency,
//	event = EXCLUDED.event,
//	sent = EXCLUDED.sent,
//	effective = EXCLUDED.effective,
//	onset = EXCLUDED.onset,
//	expires = EXCLUDED.expires,
//	ends = EXCLUDED.ends,
//	reference_ids = EXCLUDED.reference_ids,
//	border = EXCLUDED.border,
//	message_type = EXCLUDED.message_type;
func (a alertDo) InsertOptimizedAlert(id string, areaDesc string, headline string, description string, severity *string, certainty *string, urgency *string, event *string, sent time.Time, effective time.Time, onset *time.Time, expires *time.Time, ends *time.Time, referenceIDs *models.StringSlice, border *models.Geometry, messageType *string) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	params = append(params, areaDesc)
	params = append(params, headline)
	params = append(params, description)
	params = append(params, severity)
	params = append(params, certainty)
	params = append(params, urgency)
	params = append(params, event)
	params = append(params, sent)
	params = append(params, effective)
	params = append(params, onset)
	params = append(params, expires)
	params = append(params, ends)
	params = append(params, referenceIDs)
	params = append(params, border)
	params = append(params, messageType)
	generateSQL.WriteString("INSERT INTO alerts ( id, area_desc, headline, description, severity, certainty, urgency, event, sent, effective, onset, expires, ends, reference_ids, border, message_type ) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ST_UnaryUnion(?), ? ) ON CONFLICT(id) DO UPDATE SET area_desc = EXCLUDED.area_desc, headline = EXCLUDED.headline, description = EXCLUDED.description, severity = EXCLUDED.severity, certainty = EXCLUDED.certainty, urgency = EXCLUDED.urgency, event = EXCLUDED.event, sent = EXCLUDED.sent, effective = EXCLUDED.effective, onset = EXCLUDED.onset, expires = EXCLUDED.expires, ends = EXCLUDED.ends, reference_ids = EXCLUDED.reference_ids, border = EXCLUDED.border, message_type = EXCLUDED.message_type; ")

	var executeSQL *gorm.DB
	executeSQL = a.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

// WITH target_area AS (SELECT border FROM saved_areas WHERE id = @watchID LIMIT 1)
//
//	SELECT a.skeet_info AS skeet_info, EXTRACT(EPOCH FROM a.sent) * 1000 as sent
//	FROM alerts a, target_area t
//	WHERE a.skeet_info IS NOT NULL
//	AND (a.border && t.border)
//	AND ST_Intersects(a.border, t.border)
//	LIMIT @limit
//	ORDER BY a.sent DESC;
func (a alertDo) GetCustomAlertURIs(watchID string, limit uint) (result []*models.Alert, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, watchID)
	params = append(params, limit)
	generateSQL.WriteString("WITH target_area AS (SELECT border FROM saved_areas WHERE id = ? LIMIT 1) SELECT a.skeet_info AS skeet_info, EXTRACT(EPOCH FROM a.sent) * 1000 as sent FROM alerts a, target_area t WHERE a.skeet_info IS NOT NULL AND (a.border && t.border) AND ST_Intersects(a.border, t.border) LIMIT ? ORDER BY a.sent DESC; ")

	var executeSQL *gorm.DB
	executeSQL = a.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// WITH target_area AS (SELECT border FROM saved_areas WHERE ID = @watchID LIMIT 1)
//
//	SELECT a.skeet_info AS skeet_info, EXTRACT(EPOCH FROM a.sent) * 1000 as sent
//	FROM alerts a, target_area t
//	WHERE a.skeet_info IS NOT NULL
//	AND a.border && t.border AND ST_Intersects(a.border, t.border)
//	AND sent < @cursor
//	ORDER BY a.sent LIMIT @limit DESC;
func (a alertDo) GetCustomAlertURIsWithCursor(watchID string, limit uint, cursor uint) (result []*models.Alert, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, watchID)
	params = append(params, cursor)
	params = append(params, limit)
	generateSQL.WriteString("WITH target_area AS (SELECT border FROM saved_areas WHERE ID = ? LIMIT 1) SELECT a.skeet_info AS skeet_info, EXTRACT(EPOCH FROM a.sent) * 1000 as sent FROM alerts a, target_area t WHERE a.skeet_info IS NOT NULL AND a.border && t.border AND ST_Intersects(a.border, t.border) AND sent < ? ORDER BY a.sent LIMIT ? DESC; ")

	var executeSQL *gorm.DB
	executeSQL = a.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (a alertDo) Debug() IAlertDo {
	return a.withDO(a.DO.Debug())
}

func (a alertDo) WithContext(ctx context.Context) IAlertDo {
	return a.withDO(a.DO.WithContext(ctx))
}

func (a alertDo) ReadDB() IAlertDo {
	return a.Clauses(dbresolver.Read)
}

func (a alertDo) WriteDB() IAlertDo {
	return a.Clauses(dbresolver.Write)
}

func (a alertDo) Session(config *gorm.Session) IAlertDo {
	return a.withDO(a.DO.Session(config))
}

func (a alertDo) Clauses(conds ...clause.Expression) IAlertDo {
	return a.withDO(a.DO.Clauses(conds...))
}

func (a alertDo) Returning(value interface{}, columns ...string) IAlertDo {
	return a.withDO(a.DO.Returning(value, columns...))
}

func (a alertDo) Not(conds ...gen.Condition) IAlertDo {
	return a.withDO(a.DO.Not(conds...))
}

func (a alertDo) Or(conds ...gen.Condition) IAlertDo {
	return a.withDO(a.DO.Or(conds...))
}

func (a alertDo) Select(conds ...field.Expr) IAlertDo {
	return a.withDO(a.DO.Select(conds...))
}

func (a alertDo) Where(conds ...gen.Condition) IAlertDo {
	return a.withDO(a.DO.Where(conds...))
}

func (a alertDo) Order(conds ...field.Expr) IAlertDo {
	return a.withDO(a.DO.Order(conds...))
}

func (a alertDo) Distinct(cols ...field.Expr) IAlertDo {
	return a.withDO(a.DO.Distinct(cols...))
}

func (a alertDo) Omit(cols ...field.Expr) IAlertDo {
	return a.withDO(a.DO.Omit(cols...))
}

func (a alertDo) Join(table schema.Tabler, on ...field.Expr) IAlertDo {
	return a.withDO(a.DO.Join(table, on...))
}

func (a alertDo) LeftJoin(table schema.Tabler, on ...field.Expr) IAlertDo {
	return a.withDO(a.DO.LeftJoin(table, on...))
}

func (a alertDo) RightJoin(table schema.Tabler, on ...field.Expr) IAlertDo {
	return a.withDO(a.DO.RightJoin(table, on...))
}

func (a alertDo) Group(cols ...field.Expr) IAlertDo {
	return a.withDO(a.DO.Group(cols...))
}

func (a alertDo) Having(conds ...gen.Condition) IAlertDo {
	return a.withDO(a.DO.Having(conds...))
}

func (a alertDo) Limit(limit int) IAlertDo {
	return a.withDO(a.DO.Limit(limit))
}

func (a alertDo) Offset(offset int) IAlertDo {
	return a.withDO(a.DO.Offset(offset))
}

func (a alertDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IAlertDo {
	return a.withDO(a.DO.Scopes(funcs...))
}

func (a alertDo) Unscoped() IAlertDo {
	return a.withDO(a.DO.Unscoped())
}

func (a alertDo) Create(values ...*models.Alert) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Create(values)
}

func (a alertDo) CreateInBatches(values []*models.Alert, batchSize int) error {
	return a.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (a alertDo) Save(values ...*models.Alert) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Save(values)
}

func (a alertDo) First() (*models.Alert, error) {
	if result, err := a.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*models.Alert), nil
	}
}

func (a alertDo) Take() (*models.Alert, error) {
	if result, err := a.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*models.Alert), nil
	}
}

func (a alertDo) Last() (*models.Alert, error) {
	if result, err := a.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*models.Alert), nil
	}
}

func (a alertDo) Find() ([]*models.Alert, error) {
	result, err := a.DO.Find()
	return result.([]*models.Alert), err
}

func (a alertDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*models.Alert, err error) {
	buf := make([]*models.Alert, 0, batchSize)
	err = a.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (a alertDo) FindInBatches(result *[]*models.Alert, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return a.DO.FindInBatches(result, batchSize, fc)
}

func (a alertDo) Attrs(attrs ...field.AssignExpr) IAlertDo {
	return a.withDO(a.DO.Attrs(attrs...))
}

func (a alertDo) Assign(attrs ...field.AssignExpr) IAlertDo {
	return a.withDO(a.DO.Assign(attrs...))
}

func (a alertDo) Joins(fields ...field.RelationField) IAlertDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Joins(_f))
	}
	return &a
}

func (a alertDo) Preload(fields ...field.RelationField) IAlertDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Preload(_f))
	}
	return &a
}

func (a alertDo) FirstOrInit() (*models.Alert, error) {
	if result, err := a.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*models.Alert), nil
	}
}

func (a alertDo) FirstOrCreate() (*models.Alert, error) {
	if result, err := a.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*models.Alert), nil
	}
}

func (a alertDo) FindByPage(offset int, limit int) (result []*models.Alert, count int64, err error) {
	result, err = a.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = a.Offset(-1).Limit(-1).Count()
	return
}

func (a alertDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = a.Count()
	if err != nil {
		return
	}

	err = a.Offset(offset).Limit(limit).Scan(result)
	return
}

func (a alertDo) Scan(result interface{}) (err error) {
	return a.DO.Scan(result)
}

func (a alertDo) Delete(models ...*models.Alert) (result gen.ResultInfo, err error) {
	return a.DO.Delete(models)
}

func (a *alertDo) withDO(do gen.Dao) *alertDo {
	a.DO = *do.(*gen.DO)
	return a
}
