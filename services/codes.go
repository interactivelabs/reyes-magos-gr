package services

import (
	"math/rand"
	"reyes-magos-gr/db/model"
	"reyes-magos-gr/db/repository"
	"strings"
	"time"
)

type CodesService struct {
	CodesRepository repository.CodesRepository
}

func NewCode() model.Code {
	return model.Code{
		Code:       generateRandomString(6),
		Expiration: time.Now().AddDate(0, 0, 10).Format(time.RFC3339),
	}
}

func (s CodesService) CreateCode() (code model.Code, err error) {
	code = NewCode()
	_, codeRow, err := s.CodesRepository.CreateCode(code)
	if err != nil {
		return model.Code{}, err
	}

	return codeRow, nil
}

func (s CodesService) CreateCodeBatch(Count int64) (codes []model.Code, err error) {
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
