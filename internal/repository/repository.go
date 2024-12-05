package repository

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	SaveLongUrl(string, string) (error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) SaveLongUrl(longUrl, shortUrl string) (error) {
	stmt,err := r.db.Prepare("INSERT INTO ShortUrl (LongUrl, ShortUrl) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	_, err = stmt.Exec(longUrl, shortUrl)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}
	return nil
}

