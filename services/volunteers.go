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

func (s VolunteersService) GetVolunteerCodesByEmail(email string) (codes []model.Code, givenCodes []model.Code, err error) {
	volunteer, err := s.VolunteersRepository.GetVolunteerByEmail(email)
	if err != nil {
		return nil, nil, err
	}

	codes, err = s.VolunteerCodesRepository.GetActiveVolunteerCodesByVolunteerID(volunteer.VolunteerID)
	if err != nil {
		return nil, nil, err
	}

	givenCodes, err = s.VolunteerCodesRepository.GetGivenVolunteerCodesByVolunteerID(volunteer.VolunteerID)
	if err != nil {
		return nil, nil, err
	}

	return codes, givenCodes, nil
}

func (s VolunteersService) GetVolunteerOrdersByEmail(email string) (orders []model.Order, err error) {
	volunteer, err := s.VolunteersRepository.GetVolunteerByEmail(email)
	if err != nil {
		return nil, err
	}

	orders, err = s.OrdersRepository.GetPendingOrdersByVolunteerID(volunteer.VolunteerID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func GroupVolunteersByLocation(volunteers []model.Volunteer) (groupedVolunteers map[string][]model.Volunteer) {
	groupedVolunteers = make(map[string][]model.Volunteer)
	for _, volunteer := range volunteers {
		location := volunteer.State + ", " + volunteer.City
		groupedVolunteers[location] = append(groupedVolunteers[location], volunteer)
	}
	return groupedVolunteers
}

func (s VolunteersService) GetActiveVolunteersGrupedByLocation() (groupedVolunteers map[string][]model.Volunteer, err error) {

	allVolunteers, err := s.VolunteersRepository.GetActiveVolunteers()
	if err != nil {
		return nil, err
	}

	return GroupVolunteersByLocation(allVolunteers), nil
}
