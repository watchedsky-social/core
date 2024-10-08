// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

var (
	Q         = new(Query)
	Alert     *alert
	Didwid    *didwid
	Mapsearch *mapsearch
	SavedArea *savedArea
	Zone      *zone
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	Alert = &Q.Alert
	Didwid = &Q.Didwid
	Mapsearch = &Q.Mapsearch
	SavedArea = &Q.SavedArea
	Zone = &Q.Zone
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:        db,
		Alert:     newAlert(db, opts...),
		Didwid:    newDidwid(db, opts...),
		Mapsearch: newMapsearch(db, opts...),
		SavedArea: newSavedArea(db, opts...),
		Zone:      newZone(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	Alert     alert
	Didwid    didwid
	Mapsearch mapsearch
	SavedArea savedArea
	Zone      zone
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:        db,
		Alert:     q.Alert.clone(db),
		Didwid:    q.Didwid.clone(db),
		Mapsearch: q.Mapsearch.clone(db),
		SavedArea: q.SavedArea.clone(db),
		Zone:      q.Zone.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:        db,
		Alert:     q.Alert.replaceDB(db),
		Didwid:    q.Didwid.replaceDB(db),
		Mapsearch: q.Mapsearch.replaceDB(db),
		SavedArea: q.SavedArea.replaceDB(db),
		Zone:      q.Zone.replaceDB(db),
	}
}

type queryCtx struct {
	Alert     IAlertDo
	Didwid    IDidwidDo
	Mapsearch IMapsearchDo
	SavedArea ISavedAreaDo
	Zone      IZoneDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		Alert:     q.Alert.WithContext(ctx),
		Didwid:    q.Didwid.WithContext(ctx),
		Mapsearch: q.Mapsearch.WithContext(ctx),
		SavedArea: q.SavedArea.WithContext(ctx),
		Zone:      q.Zone.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
