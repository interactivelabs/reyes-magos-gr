package services

import (
	"database/sql"
	"errors"
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/dtos"
	"reyes-magos-gr/store/models"
	"time"
)

type VolunteersServiceApp struct {
	CartsStore          store.CartsStore
	CodesStore          store.CodesStore
	OrdersStore         store.OrdersStore
	ToysStore           store.ToysStore
	VolunteersStore     store.VolunteersStore
	VolunteerCodesStore store.VolunteerCodesStore
}

func NewVolunteersService(
	cartsStore store.CartsStore,
	codesStore store.CodesStore,
	ordersStore store.OrdersStore,
	toysStore store.ToysStore,
	volunteersStore store.VolunteersStore,
	volunteerCodesStore store.VolunteerCodesStore,
) *VolunteersServiceApp {
	return &VolunteersServiceApp{
		CartsStore:          cartsStore,
		CodesStore:          codesStore,
		OrdersStore:         ordersStore,
		ToysStore:           toysStore,
		VolunteersStore:     volunteersStore,
		VolunteerCodesStore: volunteerCodesStore,
	}
}

type VolunteersService interface {
	GetVolunteerByEmail(email string) (models.Volunteer, error)
	GetVolunteerCodesByEmail(
		email string,
	) (codes []models.Code, givenCodes []models.Code, err error)
	GetVolunteerOrdersByEmail(email string) (orders []models.OrderDetails, err error)
	GetRecentlyCompletedVolunteerOrdersByEmail(
		email string,
		since time.Time,
		limit int,
	) (orders []models.OrderDetails, err error)
	EnrichOrder(order models.Order) (models.OrderDetails, error)
	GetVolunteerCartByEmail(email string) (cartItems []dtos.CartItem, err error)
	CreateVolunteerCartItem(email string, toyID int64) (CartID int64, err error)
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

func (s *VolunteersServiceApp) EnrichOrder(order models.Order) (models.OrderDetails, error) {
	toy, err := s.ToysStore.GetToyByID(order.ToyID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.OrderDetails{}, err
	}

	code, err := s.CodesStore.GetCodeByID(order.CodeID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.OrderDetails{}, err
	}

	return models.OrderDetails{
		Order:        order,
		ToyName:      toy.ToyName,
		ToySourceURL: toy.SourceURL,
		Code:         code.Code,
	}, nil
}

func (s *VolunteersServiceApp) enrichOrders(orders []models.Order) ([]models.OrderDetails, error) {
	details := make([]models.OrderDetails, 0, len(orders))
	for _, order := range orders {
		detail, err := s.EnrichOrder(order)
		if err != nil {
			return nil, err
		}
		details = append(details, detail)
	}
	return details, nil
}

func (s *VolunteersServiceApp) GetVolunteerOrdersByEmail(
	email string,
) (orders []models.OrderDetails, err error) {
	volunteer, err := s.VolunteersStore.GetVolunteerByEmail(email)
	if err != nil {
		return nil, err
	}

	pendingOrders, err := s.OrdersStore.GetPendingOrdersByVolunteerID(volunteer.VolunteerID)
	if err != nil {
		return nil, err
	}

	return s.enrichOrders(pendingOrders)
}

func (s *VolunteersServiceApp) GetRecentlyCompletedVolunteerOrdersByEmail(
	email string,
	since time.Time,
	limit int,
) (orders []models.OrderDetails, err error) {
	volunteer, err := s.VolunteersStore.GetVolunteerByEmail(email)
	if err != nil {
		return nil, err
	}

	completedOrders, err := s.OrdersStore.GetRecentlyCompletedOrdersByVolunteerID(volunteer.VolunteerID, since, limit)
	if err != nil {
		return nil, err
	}

	return s.enrichOrders(completedOrders)
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

func (s *VolunteersServiceApp) CreateVolunteerCartItem(
	email string,
	toyID int64,
) (CartID int64, err error) {
	volunteer, err := s.VolunteersStore.GetVolunteerByEmail(email)
	if err != nil {
		return 0, err
	}

	item := models.CartItem{
		ToyID:       toyID,
		VolunteerID: volunteer.VolunteerID,
		Used:        0,
		Deleted:     0,
	}

	return s.CartsStore.CreateCartItem(item)
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
