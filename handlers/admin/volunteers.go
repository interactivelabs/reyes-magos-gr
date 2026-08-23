package admin

import (
	"net/http"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/services"
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/dtos"
	"reyes-magos-gr/store/models"
	views "reyes-magos-gr/views/admin/volunteers"
	"strconv"
	"strings"

	"github.com/dranikpg/dto-mapper"
	"github.com/labstack/echo/v4"
)

type VolunteersHandler struct {
	VolunteersStore   store.VolunteersStore
	VolunteersService services.VolunteersService
}

func NewVolunteersHandler(
	volunteersStore store.VolunteersStore,
	volunteersService services.VolunteersService,
) *VolunteersHandler {
	return &VolunteersHandler{
		VolunteersStore:   volunteersStore,
		VolunteersService: volunteersService,
	}
}

func (h *VolunteersHandler) VolunteersViewHandler(ctx echo.Context) error {
	q := strings.TrimSpace(ctx.QueryParam("q"))

	groups, resultLabel, err := h.fetchVolunteerGroups(q)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if ctx.Request().Header.Get("HX-Request") != "" {
		return lib.Render(ctx, views.VolunteerSearchResults(groups, resultLabel))
	}

	return lib.Render(ctx, views.AdminVolunteers(groups, resultLabel, q))
}

// fetchVolunteerGroups fetches the active volunteers with their given/held
// counts and last-activity date, filtering by name/email when q is set.
func (h *VolunteersHandler) fetchVolunteerGroups(
	q string,
) (groups []dtos.VolunteerLocationGroup, resultLabel string, err error) {
	allGroups, err := h.VolunteersService.GetActiveVolunteersWithStatsGroupedByLocation()
	if err != nil {
		return nil, "", err
	}

	filtered, total := filterVolunteerGroups(allGroups, q)

	return filtered, views.VolunteerCountLabel(total), nil
}

func filterVolunteerGroups(
	groups []dtos.VolunteerLocationGroup,
	q string,
) (filtered []dtos.VolunteerLocationGroup, total int) {
	if q == "" {
		for _, group := range groups {
			total += len(group.Volunteers)
		}
		return groups, total
	}

	needle := strings.ToLower(q)
	for _, group := range groups {
		var matched []dtos.VolunteerListItem
		for _, volunteer := range group.Volunteers {
			if strings.Contains(strings.ToLower(volunteer.Name), needle) ||
				strings.Contains(strings.ToLower(volunteer.Email), needle) {
				matched = append(matched, volunteer)
			}
		}
		if len(matched) > 0 {
			filtered = append(filtered, dtos.VolunteerLocationGroup{
				Location:   group.Location,
				Volunteers: matched,
			})
			total += len(matched)
		}
	}

	return filtered, total
}

func (h *VolunteersHandler) VolunteersCreateHandler(ctx echo.Context) error {
	return lib.Render(ctx, views.CreateVolunteerForm())
}

type CreateVolunteerRequest struct {
	Name     string `form:"name"     validate:"required"`
	Email    string `form:"email"    validate:"required"`
	Phone    string `form:"phone"`
	Address  string `form:"address"  validate:"required"`
	Address2 string `form:"address2"`
	Country  string `form:"country"  validate:"required"`
	State    string `form:"state"    validate:"required"`
	City     string `form:"city"     validate:"required"`
	Province string `form:"province"`
	ZipCode  string `form:"zip_code" validate:"required"`
}

func (h *VolunteersHandler) VolunteersCreatePostHandler(ctx echo.Context) error {
	tr := new(CreateVolunteerRequest)
	if err := ctx.Bind(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var volunteer models.Volunteer
	err := dto.Map(&volunteer, tr)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err = h.VolunteersService.CreateAndGetVolunteer(volunteer)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Redirect(303, "/admin/volunteers")
}

func (h *VolunteersHandler) VolunteersUpdateViewHandler(ctx echo.Context) error {
	volunteerID, err := strconv.ParseInt(ctx.Param("volunteer_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	volunteer, err := h.VolunteersStore.GetVolunteerByID(volunteerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, views.UpdateVolunteerForm(volunteer))
}

func (h *VolunteersHandler) VolunteersUpdatePutHandler(ctx echo.Context) error {
	volunteerIDStr := ctx.Param("volunteer_id")
	volunteerID, err := strconv.ParseInt(volunteerIDStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid volunteer ID")
	}

	tr := new(CreateVolunteerRequest)
	if err := ctx.Bind(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ctx.Validate(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var volunteer models.Volunteer
	err = dto.Map(&volunteer, tr)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	volunteer, err = h.VolunteersService.UpdateVolunteer(volunteer, volunteerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Redirect(303, "/admin/volunteers")
}

func (h *VolunteersHandler) VolunteersDeleteHandler(ctx echo.Context) error {
	volunteerIDStr := ctx.Param("volunteer_id")
	volunteerID, err := strconv.ParseInt(volunteerIDStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid volunteer ID")
	}

	err = h.VolunteersStore.DeleteVolunteer(volunteerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}
