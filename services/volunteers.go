package services

import (
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/dtos"
	"reyes-magos-gr/store/models"

	"github.com/dranikpg/dto-mapper"
)

type VolunteersService struct {
	CartsRepository          store.CartsRepository
	CodesRepository          store.CodesRepository
	OrdersRepository         store.OrdersRepository
	VolunteersRepository     store.VolunteersRepository
	VolunteerCodesRepository store.VolunteerCodesRepository
}

func (s VolunteersService) GetVolunteerByEmail(email string) (models.Volunteer, error) {
	return s.VolunteersRepository.GetVolunteerByEmail(email)
}

func (s VolunteersService) GetVolunteerCodesByEmail(email string) (codes []models.Code, givenCodes []models.Code, err error) {
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

func (s VolunteersService) GetVolunteerOrdersByEmail(email string) (orders []models.Order, err error) {
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

func (s VolunteersService) GetVolunteerCartByEmail(email string) (cartItems []dtos.CartItem, err error) {
	volunteer, err := s.VolunteersRepository.GetVolunteerByEmail(email)
	if err != nil {
		return nil, err
	}

	cartItems, err = s.CartsRepository.GetCartToys(volunteer.VolunteerID)
	if err != nil {
		return nil, err
	}

	return cartItems, nil
}

func GroupVolunteersByLocation(volunteers []models.Volunteer) (groupedVolunteers map[string][]models.Volunteer) {
	groupedVolunteers = make(map[string][]models.Volunteer)
	for _, volunteer := range volunteers {
		location := volunteer.State + ", " + volunteer.City
		groupedVolunteers[location] = append(groupedVolunteers[location], volunteer)
	}
	return groupedVolunteers
}

func (s VolunteersService) GetActiveVolunteersGrupedByLocation() (groupedVolunteers map[string][]models.Volunteer, err error) {

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
	var volunteer models.Volunteer
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

func (s VolunteersService) CreateAndGetVolunteer(tr CreateVolunteerRequest) (volunteer models.Volunteer, err error) {
	volunteerID, err := s.CreateVolunteer(tr)
	if err != nil {
		return models.Volunteer{}, err
	}

	volunteer, err = s.VolunteersRepository.GetVolunteerByID(volunteerID)
	if err != nil {
		return models.Volunteer{}, err
	}

	return volunteer, nil
}

func (h VolunteersService) UpdateVolunteer(tr CreateVolunteerRequest, volunteerID int64) (volunteer models.Volunteer, err error) {
	err = dto.Map(&volunteer, tr)
	if err != nil {
		return models.Volunteer{}, err
	}

	volunteer.VolunteerID = volunteerID

	err = h.VolunteersRepository.UpdateVolunteer(volunteer)
	if err != nil {
		return models.Volunteer{}, err
	}

	return volunteer, nil
}
