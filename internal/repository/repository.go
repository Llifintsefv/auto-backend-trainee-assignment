package repository

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	SaveLongUrl(string, string) (error)
	GetLongUrl(string)(string,error)
	FindShortUrl(string)(bool,error)
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
	stmt,err := r.db.Prepare("INSERT INTO urls (LongUrl, ShortUrl) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(longUrl, shortUrl)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}
	return nil
}


func (r *repository) GetLongUrl(shortUrl string) (string, error) {
    stmt, err := r.db.Prepare("SELECT LongUrl FROM urls WHERE ShortUrl = $1")
    if err != nil {
        return "", fmt.Errorf("failed to prepare statement: %w", err)
    }
    defer stmt.Close() 

    var longUrl string
    err = stmt.QueryRow(shortUrl).Scan(&longUrl) 
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("shortUrl not found: %w", err)
        }
        return "", fmt.Errorf("failed to scan row: %w", err)
    }

    return longUrl, nil
}

func (r *repository) FindShortUrl(shortUrl string) (bool,error) {
	stmt, err := r.db.Prepare("SELECT ShortUrl FROM urls WHERE ShortUrl = $1")
	if err != nil {
		return false, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var ShortUrl string
	err = stmt.QueryRow(shortUrl).Scan(&ShortUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to scan row: %w", err)
	}

	return true, nil
}