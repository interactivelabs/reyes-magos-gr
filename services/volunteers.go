package services

import (
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/dtos"
	"reyes-magos-gr/store/models"
)

type VolunteersServiceApp struct {
	CartsStore               store.CartsStore
	CodesRepository          store.CodesRepository
	OrdersRepository         store.OrdersRepository
	VolunteersRepository     store.VolunteersRepository
	VolunteerCodesRepository store.VolunteerCodesRepository
}

func NewVolunteersService(cartsStore store.CartsStore, codesRepository store.CodesRepository, ordersRepository store.OrdersRepository, volunteersRepository store.VolunteersRepository, volunteerCodesRepository store.VolunteerCodesRepository) *VolunteersServiceApp {
	return &VolunteersServiceApp{
		CartsStore:               cartsStore,
		CodesRepository:          codesRepository,
		OrdersRepository:         ordersRepository,
		VolunteersRepository:     volunteersRepository,
		VolunteerCodesRepository: volunteerCodesRepository,
	}
}

type VolunteersService interface {
	GetVolunteerByEmail(email string) (models.Volunteer, error)
	GetVolunteerCodesByEmail(email string) (codes []models.Code, givenCodes []models.Code, err error)
	GetVolunteerOrdersByEmail(email string) (orders []models.Order, err error)
	GetVolunteerCartByEmail(email string) (cartItems []dtos.CartItem, err error)
	GetActiveVolunteersGrupedByLocation() (groupedVolunteers map[string][]models.Volunteer, err error)
	CreateAndGetVolunteer(volunteer models.Volunteer) (models.Volunteer, error)
	UpdateVolunteer(volunteer models.Volunteer, volunteerID int64) (models.Volunteer, error)
}

func (s VolunteersServiceApp) GetVolunteerByEmail(email string) (models.Volunteer, error) {
	return s.VolunteersRepository.GetVolunteerByEmail(email)
}

func (s VolunteersServiceApp) GetVolunteerCodesByEmail(email string) (codes []models.Code, givenCodes []models.Code, err error) {
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

func (s VolunteersServiceApp) GetVolunteerOrdersByEmail(email string) (orders []models.Order, err error) {
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

func (s VolunteersServiceApp) GetVolunteerCartByEmail(email string) (cartItems []dtos.CartItem, err error) {
	volunteer, err := s.VolunteersRepository.GetVolunteerByEmail(email)
	if err != nil {
		return nil, err
	}

	cartItems, err = s.CartsStore.GetCartToys(volunteer.VolunteerID)
	if err != nil {
		return nil, err
	}

	return cartItems, nil
}

func (s VolunteersServiceApp) GetActiveVolunteersGrupedByLocation() (groupedVolunteers map[string][]models.Volunteer, err error) {

	allVolunteers, err := s.VolunteersRepository.GetActiveVolunteers()
	if err != nil {
		return nil, err
	}

	return GroupVolunteersByLocation(allVolunteers), nil
}

func (s VolunteersServiceApp) CreateAndGetVolunteer(volunteer models.Volunteer) (models.Volunteer, error) {
	volunteerID, err := s.VolunteersRepository.CreateVolunteer(volunteer)
	if err != nil {
		return models.Volunteer{}, err
	}

	volunteer, err = s.VolunteersRepository.GetVolunteerByID(volunteerID)
	if err != nil {
		return models.Volunteer{}, err
	}

	return volunteer, nil
}

func (h VolunteersServiceApp) UpdateVolunteer(volunteer models.Volunteer, volunteerID int64) (models.Volunteer, error) {

	volunteer.VolunteerID = volunteerID

	err := h.VolunteersRepository.UpdateVolunteer(volunteer)
	if err != nil {
		return models.Volunteer{}, err
	}

	return volunteer, nil
}

func GroupVolunteersByLocation(volunteers []models.Volunteer) (groupedVolunteers map[string][]models.Volunteer) {
	groupedVolunteers = make(map[string][]models.Volunteer)
	for _, volunteer := range volunteers {
		location := volunteer.State + ", " + volunteer.City
		groupedVolunteers[location] = append(groupedVolunteers[location], volunteer)
	}
	return groupedVolunteers
}
