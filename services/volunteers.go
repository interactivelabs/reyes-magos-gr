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

func (s VolunteersService) GetVolunteerCodesByEmail(email string) ([]model.Code, error) {
	volunteer, err := s.VolunteersRepository.GetVolunteerByEmail(email)
	if err != nil {
		return nil, err
	}

	volunteerCodes, err := s.VolunteerCodesRepository.GetAllVolunteerCodesByVolunteerID(volunteer.VolunteerID)
	if err != nil {
		return nil, err
	}

	return volunteerCodes, nil
}
