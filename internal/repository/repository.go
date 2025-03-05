package repository

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
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

func initDB(dbDriver, dbConnectionString, sqlFilePath string) (*sql.DB, error) {
	// 1. Подключение к базе данных.
	db, err := sql.Open(dbDriver, dbConnectionString)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Проверка подключения
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка ping к БД: %w", err)
	}

	// 2. Проверка и чтение SQL файла.
	if _, err := os.Stat(sqlFilePath); err != nil {
		return nil, fmt.Errorf("SQL файл не найден: %w", err)
	}

	sqlBytes, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения SQL файла: %w", err)
	}
	sqlScript := string(sqlBytes)

	// 3. Выполнение SQL скрипта.
	_, err = db.Exec(sqlScript)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения SQL скрипта: %w", err)
	}

	fmt.Println("База данных инициализирована успешно.")
	return db, nil
}
