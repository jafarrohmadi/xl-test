package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type NewRepositoryOptions struct {
	Dsn string
}

func NewRepository(options NewRepositoryOptions) *Repository {
	db, err := gorm.Open(postgres.Open(options.Dsn), &gorm.Config{})
	if err != nil {
		return nil
	}

	return &Repository{db: db}
}
