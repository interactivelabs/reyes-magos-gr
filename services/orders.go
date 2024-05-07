package services

import (
	"errors"
	"reyes-magos-gr/db/model"
	"reyes-magos-gr/db/repository"
	"time"
)

type OrdersService struct {
	CodesRepository          repository.CodesRepository
	OrdersRepository         repository.OrdersRepository
	VolunteerCodesRepository repository.VolunteerCodesRepository
}

func (s OrdersService) CreateOrder(toyID int64, code string) (model.Order, error) {
	var order model.Order

	codeResult, err := s.CodesRepository.GetCode(code)
	if err != nil {
		return order, err
	}

	if codeResult.Used == 1 {
		return order, errors.New("Code already used")
	}

	volunteerID, err := s.VolunteerCodesRepository.GetVolunteerIdByCodeId(codeResult.CodeID)
	if err != nil {
		return order, err
	}

	order = model.Order{
		ToyID:       toyID,
		CodeID:      codeResult.CodeID,
		VolunteerID: volunteerID,
		OrderDate:   time.Now().Format("2006-01-02"),
	}

	orderID, err := s.OrdersRepository.CreateOrder(order)
	if err != nil {
		return order, err
	}

	codeResult.Used = 1
	err = s.CodesRepository.UpdateCode(codeResult)
	if err != nil {
		return order, err
	}

	order, err = s.OrdersRepository.GetOrderByID(orderID)
	if err != nil {
		return order, err
	}

	return order, nil
}
