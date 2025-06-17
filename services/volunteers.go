package services

import (
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/dtos"
	"reyes-magos-gr/store/models"
)

type VolunteersServiceApp struct {
	CartsStore          store.CartsStore
	CodesStore     store.CodesStore
	OrdersStore    store.OrdersStore
	VolunteersStore     store.VolunteersStore
	VolunteerCodesStore store.VolunteerCodesStore
}

func NewVolunteersService(
	cartsStore store.CartsStore,
	codesStore store.CodesStore,
	ordersStore store.OrdersStore,
	volunteersStore store.VolunteersStore,
	volunteerCodesStore store.VolunteerCodesStore,
) *VolunteersServiceApp {
	return &VolunteersServiceApp{
		CartsStore:          cartsStore,
		CodesStore:     codesStore,
		OrdersStore:    ordersStore,
		VolunteersStore:     volunteersStore,
		VolunteerCodesStore: volunteerCodesStore,
	}
}

type VolunteersService interface {
	GetVolunteerByEmail(email string) (models.Volunteer, error)
	GetVolunteerCodesByEmail(
		email string,
	) (codes []models.Code, givenCodes []models.Code, err error)
	GetVolunteerOrdersByEmail(email string) (orders []models.Order, err error)
	GetVolunteerCartByEmail(email string) (cartItems []dtos.CartItem, err error)
	GetActiveVolunteersGrupedByLocation() (groupedVolunteers map[string][]models.Volunteer, err error)
	CreateAndGetVolunteer(volunteer models.Volunteer) (models.Volunteer, error)
	UpdateVolunteer(volunteer models.Volunteer, volunteerID int64) (models.Volunteer, error)
}

func (s *VolunteersServiceApp) GetVolunteerByEmail(email string) (models.Volunteer, error) {
	return s.VolunteersStore.GetVolunteerByEmail(email)
}

func (s *VolunteersServiceApp) GetVolunteerCodesByEmail(
	email string,
) (codes []models.Code, givenCodes []models.Code, err error) {
	volunteer, err := s.VolunteersStore.GetVolunteerByEmail(email)
	if err != nil {
		return nil, nil, err
	}

	codes, err = s.VolunteerCodesStore.GetActiveVolunteerCodesByVolunteerID(
		volunteer.VolunteerID,
	)
	if err != nil {
		return nil, nil, err
	}

	givenCodes, err = s.VolunteerCodesStore.GetGivenVolunteerCodesByVolunteerID(
		volunteer.VolunteerID,
	)
	if err != nil {
		return nil, nil, err
	}

	return codes, givenCodes, nil
}

func (s *VolunteersServiceApp) GetVolunteerOrdersByEmail(
	email string,
) (orders []models.Order, err error) {
	volunteer, err := s.VolunteersStore.GetVolunteerByEmail(email)
	if err != nil {
		return nil, err
	}

	orders, err = s.OrdersStore.GetPendingOrdersByVolunteerID(volunteer.VolunteerID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *VolunteersServiceApp) GetVolunteerCartByEmail(
	email string,
) (cartItems []dtos.CartItem, err error) {
	volunteer, err := s.VolunteersStore.GetVolunteerByEmail(email)
	if err != nil {
		return nil, err
	}

	cartItems, err = s.CartsStore.GetCartToys(volunteer.VolunteerID)
	if err != nil {
		return nil, err
	}

	return cartItems, nil
}

func (s *VolunteersServiceApp) GetActiveVolunteersGrupedByLocation() (groupedVolunteers map[string][]models.Volunteer, err error) {

	allVolunteers, err := s.VolunteersStore.GetActiveVolunteers()
	if err != nil {
		return nil, err
	}

	return groupVolunteersByLocation(allVolunteers), nil
}

func (s *VolunteersServiceApp) CreateAndGetVolunteer(
	volunteer models.Volunteer,
) (models.Volunteer, error) {
	volunteerID, err := s.VolunteersStore.CreateVolunteer(volunteer)
	if err != nil {
		return models.Volunteer{}, err
	}

	volunteer, err = s.VolunteersStore.GetVolunteerByID(volunteerID)
	if err != nil {
		return models.Volunteer{}, err
	}

	return volunteer, nil
}

func (s *VolunteersServiceApp) UpdateVolunteer(
	volunteer models.Volunteer,
	volunteerID int64,
) (models.Volunteer, error) {

	volunteer.VolunteerID = volunteerID

	err := s.VolunteersStore.UpdateVolunteer(volunteer)
	if err != nil {
		return models.Volunteer{}, err
	}

	return volunteer, nil
}

func groupVolunteersByLocation(
	volunteers []models.Volunteer,
) (groupedVolunteers map[string][]models.Volunteer) {
	groupedVolunteers = make(map[string][]models.Volunteer)
	for _, volunteer := range volunteers {
		location := volunteer.State + ", " + volunteer.City
		groupedVolunteers[location] = append(groupedVolunteers[location], volunteer)
	}
	return groupedVolunteers
}
