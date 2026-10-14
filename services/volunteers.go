package services

import (
	"database/sql"
	"errors"
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/dtos"
	"reyes-magos-gr/store/models"
	"sort"
	"time"
)

type VolunteersServiceApp struct {
	CartsStore            store.CartsStore
	CodesStore            store.CodesStore
	OrdersStore           store.OrdersStore
	ToysStore             store.ToysStore
	VolunteersStore       store.VolunteersStore
	VolunteerCodesStore   store.VolunteerCodesStore
	VolunteersOrdersStore store.VolunteersOrdersStore
}

func NewVolunteersService(
	cartsStore store.CartsStore,
	codesStore store.CodesStore,
	ordersStore store.OrdersStore,
	toysStore store.ToysStore,
	volunteersStore store.VolunteersStore,
	volunteerCodesStore store.VolunteerCodesStore,
	volunteersOrdersStore store.VolunteersOrdersStore,
) *VolunteersServiceApp {
	return &VolunteersServiceApp{
		CartsStore:            cartsStore,
		CodesStore:            codesStore,
		OrdersStore:           ordersStore,
		ToysStore:             toysStore,
		VolunteersStore:       volunteersStore,
		VolunteerCodesStore:   volunteerCodesStore,
		VolunteersOrdersStore: volunteersOrdersStore,
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
	DuplicateVolunteerCartItem(email string, cartID int64) error
	DeleteVolunteerCartItem(email string, cartID int64) error
	GetActiveVolunteersGrupedByLocation() (groupedVolunteers map[string][]models.Volunteer, err error)
	GetActiveVolunteersWithStatsGroupedByLocation() (groups []dtos.VolunteerLocationGroup, err error)
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
	return enrichOrder(s.ToysStore, s.CodesStore, order)
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

	return enrichOrders(s.ToysStore, s.CodesStore, pendingOrders)
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

	completedOrders, err := s.OrdersStore.GetRecentlyCompletedOrdersByVolunteerID(
		volunteer.VolunteerID,
		since,
		limit,
	)
	if err != nil {
		return nil, err
	}

	return enrichOrders(s.ToysStore, s.CodesStore, completedOrders)
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

var (
	ErrCartItemNotFound = errors.New("cart item not found")
	ErrNotEnoughCodes   = errors.New("not enough codes available")
)

// getOwnedCartItem returns the cart item only when it belongs to the volunteer,
// so a volunteer can't act on someone else's cart by guessing IDs.
func (s *VolunteersServiceApp) getOwnedCartItem(
	volunteerID int64,
	cartID int64,
) (models.CartItem, error) {
	item, err := s.CartsStore.GetCartItemByID(cartID)
	if errors.Is(err, sql.ErrNoRows) {
		return models.CartItem{}, ErrCartItemNotFound
	}
	if err != nil {
		return models.CartItem{}, err
	}
	if item.VolunteerID != volunteerID {
		return models.CartItem{}, ErrCartItemNotFound
	}
	return item, nil
}

// DuplicateVolunteerCartItem adds another entry for the same toy, as long as
// the volunteer still has an available code to cover it.
func (s *VolunteersServiceApp) DuplicateVolunteerCartItem(email string, cartID int64) error {
	volunteer, err := s.VolunteersStore.GetVolunteerByEmail(email)
	if err != nil {
		return err
	}

	item, err := s.getOwnedCartItem(volunteer.VolunteerID, cartID)
	if err != nil {
		return err
	}

	cartItems, err := s.CartsStore.GetCartToys(volunteer.VolunteerID)
	if err != nil {
		return err
	}

	codes, err := s.VolunteerCodesStore.GetActiveVolunteerCodesByVolunteerID(volunteer.VolunteerID)
	if err != nil {
		return err
	}

	if len(cartItems) >= len(codes) {
		return ErrNotEnoughCodes
	}

	_, err = s.CartsStore.CreateCartItem(models.CartItem{
		ToyID:       item.ToyID,
		VolunteerID: volunteer.VolunteerID,
	})
	return err
}

func (s *VolunteersServiceApp) DeleteVolunteerCartItem(email string, cartID int64) error {
	volunteer, err := s.VolunteersStore.GetVolunteerByEmail(email)
	if err != nil {
		return err
	}

	if _, err := s.getOwnedCartItem(volunteer.VolunteerID, cartID); err != nil {
		return err
	}

	return s.CartsStore.DeleteCartItem(cartID)
}

func (s *VolunteersServiceApp) GetActiveVolunteersGrupedByLocation() (groupedVolunteers map[string][]models.Volunteer, err error) {

	allVolunteers, err := s.VolunteersStore.GetActiveVolunteers()
	if err != nil {
		return nil, err
	}

	return groupVolunteersByLocation(allVolunteers), nil
}

func (s *VolunteersServiceApp) GetActiveVolunteersWithStatsGroupedByLocation() (
	groups []dtos.VolunteerLocationGroup,
	err error,
) {
	items, err := s.VolunteersOrdersStore.GetActiveVolunteersWithStats()
	if err != nil {
		return nil, err
	}

	return groupVolunteerListItemsByLocation(items), nil
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

func groupVolunteerListItemsByLocation(
	items []dtos.VolunteerListItem,
) []dtos.VolunteerLocationGroup {
	byLocation := make(map[string][]dtos.VolunteerListItem)
	for _, item := range items {
		location := item.State + ", " + item.City
		byLocation[location] = append(byLocation[location], item)
	}

	locations := make([]string, 0, len(byLocation))
	for location := range byLocation {
		locations = append(locations, location)
	}
	sort.Strings(locations)

	groups := make([]dtos.VolunteerLocationGroup, 0, len(locations))
	for _, location := range locations {
		groups = append(groups, dtos.VolunteerLocationGroup{
			Location:   location,
			Volunteers: byLocation[location],
		})
	}

	return groups
}
