package services

import (
	"reyes-magos-gr/db/model"
	"reyes-magos-gr/db/repository"

	"github.com/dranikpg/dto-mapper"
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

type CreateVolunteerRequest struct {
	Name     string `form:"name" validate:"required"`
	Email    string `form:"email" validate:"required"`
	Phone    string `form:"phone"`
	Address  string `form:"address" validate:"required"`
	Address2 string `form:"address2"`
	Country  string `form:"country" validate:"required"`
	State    string `form:"state" validate:"required"`
	City     string `form:"city" validate:"required"`
	Province string `form:"province"`
	ZipCode  string `form:"zip_code" validate:"required"`
}

func (s VolunteersService) CreateVolunteer(tr CreateVolunteerRequest) (volunteerID int64, err error) {
	var volunteer model.Volunteer
	err = dto.Map(&volunteer, tr)
	if err != nil {
		return 0, err
	}

	volunteerID, err = s.VolunteersRepository.CreateVolunteer(volunteer)
	if err != nil {
		return 0, err
	}

	return volunteerID, nil
}

func (s VolunteersService) CreateAndGetVolunteer(tr CreateVolunteerRequest) (volunteer model.Volunteer, err error) {
	volunteerID, err := s.CreateVolunteer(tr)
	if err != nil {
		return model.Volunteer{}, err
	}

	volunteer, err = s.VolunteersRepository.GetVolunteerByID(volunteerID)
	if err != nil {
		return model.Volunteer{}, err
	}

	return volunteer, nil
}

type UpdateVolunteerRequest struct {
	VolunteerID int64  `form:"volunteer_id" validate:"required"`
	Name        string `form:"name" validate:"required"`
	Email       string `form:"email" validate:"required"`
	Phone       string `form:"phone"`
	Address     string `form:"address" validate:"required"`
	Address2    string `form:"address2"`
	Country     string `form:"country" validate:"required"`
	State       string `form:"state" validate:"required"`
	City        string `form:"city" validate:"required"`
	Province    string `form:"province"`
	ZipCode     string `form:"zip_code" validate:"required"`
}

func (h VolunteersService) UpdateVolunteer(tr UpdateVolunteerRequest) (volunteer model.Volunteer, err error) {
	err = dto.Map(&volunteer, tr)
	if err != nil {
		return model.Volunteer{}, err
	}

	err = h.VolunteersRepository.UpdateVolunteer(volunteer)
	if err != nil {
		return model.Volunteer{}, err
	}

	return volunteer, nil
}
