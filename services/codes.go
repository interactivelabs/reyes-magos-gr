package services

import (
	"math/rand"
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/models"
	"strings"
	"time"
)

type CodesServiceApp struct {
	CodesStore store.CodesStore
}

func NewCodesService(codesStore store.CodesStore) *CodesServiceApp {
	return &CodesServiceApp{
		CodesStore: codesStore,
	}
}

type CodesService interface {
	CreateCode() (code models.Code, err error)
	CreateCodeBatch(Count int64) (codes []models.Code, err error)
}

func newCode() models.Code {
	return models.Code{
		Code:       generateRandomString(6),
		Expiration: time.Now().AddDate(0, 0, 10).Format(time.RFC3339),
	}
}

func (s *CodesServiceApp) CreateCode() (code models.Code, err error) {
	code = newCode()
	_, codeRow, err := s.CodesStore.CreateCode(code)
	if err != nil {
		return models.Code{}, err
	}

	return codeRow, nil
}

func (s *CodesServiceApp) CreateCodeBatch(Count int64) (codes []models.Code, err error) {
	for i := 0; i < int(Count); i++ {
		code := newCode()
		_, codeRow, err := s.CodesStore.CreateCode(code)
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
	for range length {
		sb.WriteByte(chars[rand.Intn(len(chars))])
	}
	return sb.String()
}
