package postgrestore

import (
	"context"

	"gorm.io/gorm"
)

type DBTransaction struct {
	db *gorm.DB
}

func NewDBTransaction(db *gorm.DB) *DBTransaction {
	return &DBTransaction{db: db}
}

func (r *DBTransaction) Begin(ctx context.Context) (_interface.IDBTransaction, error) {
	tx := r.db.WithContext(ctx).Begin()
	return &DBTransaction{db: tx}, tx.Error
}

func (r *DBTransaction) Commit(ctx context.Context) error {
	return r.db.WithContext(ctx).Commit().Error
}

func (r *DBTransaction) Rollback(ctx context.Context) {
	r.db.WithContext(ctx).Rollback()
}
