package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSaveLongUrl(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	longUrl := "https://example.com"
	shortUrl := "short1"

	mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO urls (LongUrl, ShortUrl) VALUES ($1, $2)")).
	ExpectExec(). 
	WithArgs(longUrl,shortUrl). 
	WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SaveLongUrl(longUrl,shortUrl)
	assert.NoError(t, err)

}


func TestGetLongUrl(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	shortUrl := "short1"

	rows := sqlmock.NewRows([]string{"LongUrl"}).AddRow("https://example.com")

	mock.ExpectPrepare(regexp.QuoteMeta("SELECT LongUrl FROM urls WHERE ShortUrl = $1")).
		ExpectQuery().WithArgs(shortUrl).WillReturnRows(rows)

	longUrl, err := repo.GetLongUrl(shortUrl)
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", longUrl)

	mock.ExpectPrepare(regexp.QuoteMeta("SELECT LongUrl FROM urls WHERE ShortUrl = $1")).
		ExpectQuery().WithArgs(shortUrl).WillReturnError(sql.ErrNoRows)

	_, err = repo.GetLongUrl(shortUrl)
	assert.Error(t, err)
}

func TestFindShortUrl(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	shortUrl := "short1"

	rows := sqlmock.NewRows([]string{"ShortUrl"}).AddRow("short1")

	mock.ExpectPrepare(regexp.QuoteMeta("SELECT ShortUrl FROM urls WHERE ShortUrl = $1")).
		ExpectQuery().WithArgs(shortUrl).WillReturnRows(rows)

	ok, err := repo.FindShortUrl(shortUrl)
	assert.NoError(t, err)
	assert.True(t, ok)

	mock.ExpectPrepare(regexp.QuoteMeta("SELECT ShortUrl FROM urls WHERE ShortUrl = $1")).
		ExpectQuery().
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	ok, err = repo.FindShortUrl("nonexistent")
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.NoError(t, mock.ExpectationsWereMet())
}