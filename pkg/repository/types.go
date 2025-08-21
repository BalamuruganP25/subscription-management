package repository

import "database/sql"

type CrudRepo interface {
}

type CurdRepository struct {
	db *sql.DB
}

func NewCurdRepo(db *sql.DB) *CurdRepository {
	return &CurdRepository{
		db: db,
	}
}
