package service

import (
	"auto-backend-trainee-assignment/internal/models"
	"auto-backend-trainee-assignment/internal/repository"
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

type Service interface {
	GenerateShortUrl(models.ResponseUrl) (string,error)
	GetLongUrl(string) (string,error)
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GenerateShortUrl(RespUrl models.ResponseUrl) (string,error) {
	var ShortUrl string
	switch RespUrl.CustomUrl {
	case "":
		ShortUrl = GenShortUrl()
	default:
		if !isValidCustomUrl(RespUrl.CustomUrl) {
            return "", fmt.Errorf("invalid custom URL: %s", RespUrl.CustomUrl)
        }

        fmt.Println(RespUrl.CustomUrl)
        ok, err := s.repo.FindShortUrl(RespUrl.CustomUrl)
        if err != nil {
            return "", fmt.Errorf("failed to scan row: %w", err)
        }
        if ok {
            return "", fmt.Errorf("short URL already exists")
        }
        ShortUrl = RespUrl.CustomUrl
	}
	
	
	err := s.repo.SaveLongUrl(RespUrl.Url,ShortUrl)
	if err != nil {
		return "", fmt.Errorf("failed to save long URL: %w", err)
	}
	return ShortUrl,nil

}

func (s *service)GetLongUrl(ShortUrl string) (string,error) {
	LongUrl,err := s.repo.GetLongUrl(ShortUrl)
	if err != nil {
		return "",err
	}
	return LongUrl,nil
}

func GenShortUrl() (string) {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8
	short := make([]byte, length)
	for i := range short {
		short[i] = charset[rand.Intn(len(charset))]
	}
	return string(short)
}


func isValidCustomUrl(customUrl string) bool {
    // Регулярное выражение: допускаются только буквы, цифры, дефисы и подчеркивания
    validUrlPattern := `^[a-zA-Z0-9_-]+$`
    match, _ := regexp.MatchString(validUrlPattern, customUrl)
    return match
}