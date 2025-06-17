package services

import (
	"math/rand"
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/models"
	"strings"
	"time"
)

type CodesService struct {
	CodesRepository store.CodesRepository
}

func NewCode() models.Code {
	return models.Code{
		Code:       generateRandomString(6),
		Expiration: time.Now().AddDate(0, 0, 10).Format(time.RFC3339),
	}
}

func (s CodesService) CreateCode() (code models.Code, err error) {
	code = NewCode()
	_, codeRow, err := s.CodesRepository.CreateCode(code)
	if err != nil {
		return models.Code{}, err
	}

	return codeRow, nil
}

func (s CodesService) CreateCodeBatch(Count int64) (codes []models.Code, err error) {
	for i := 0; i < int(Count); i++ {
		code := NewCode()
		_, codeRow, err := s.CodesRepository.CreateCode(code)
		if err != nil {
			return nil, err
		}
		codes = append(codes, codeRow)
	}
	return codes, nil
}

func generateRandomString(length int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(chars[rand.Intn(len(chars))])
	}
	return sb.String()
}
