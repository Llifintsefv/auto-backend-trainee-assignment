package service

import (
	"auto-backend-trainee-assignment/internal/models"
	"auto-backend-trainee-assignment/internal/repository"
	"math/rand"
	"time"
)

type Service interface {
	GenerateShortUrl(models.ResponseUrl) (string,error)
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GenerateShortUrl(RespUrl models.ResponseUrl) (string,error) {
	rand.Seed(time.Now().UnixNano())
	ShortUrl := GenShortUrl()
	err := s.repo.SaveLongUrl(RespUrl.Url,ShortUrl)
	if err != nil {

	}
	return ShortUrl,nil

}


func GenShortUrl() (string) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8
	short := make([]byte, length)
	for i := range short {
		short[i] = charset[rand.Intn(len(charset))]
	}
	return string(short)
}