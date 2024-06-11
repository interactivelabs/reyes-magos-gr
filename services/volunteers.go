package services

import (
	"reyes-magos-gr/db/model"
	"reyes-magos-gr/db/repository"
)

type VolunteersService struct {
	CodesRepository          repository.CodesRepository
	OrdersRepository         repository.OrdersRepository
	VolunteersRepository     repository.VolunteersRepository
	VolunteerCodesRepository repository.VolunteerCodesRepository
}

func (s VolunteersService) GetVolunteerCodesByEmail(email string) (codes []model.Code, err error) {
	volunteer, err := s.VolunteersRepository.GetVolunteerByEmail(email)
	if err != nil {
		return nil, err
	}

	codes, err = s.VolunteerCodesRepository.GetAllVolunteerCodesByVolunteerID(volunteer.VolunteerID)
	if err != nil {
		return nil, err
	}

	return codes, nil
}

func (s VolunteersService) GetVolunteerOrdersByEmail(email string) (orders []model.Order, err error) {
	volunteer, err := s.VolunteersRepository.GetVolunteerByEmail(email)
	if err != nil {
		return nil, err
	}

	orders, err = s.OrdersRepository.GetOrdersByVolunteerID(volunteer.VolunteerID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
