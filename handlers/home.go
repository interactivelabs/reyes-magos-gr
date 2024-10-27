package handlers

import (
	"reyes-magos-gr/lib"
	home "reyes-magos-gr/views/home"
	support "reyes-magos-gr/views/support"
	volunteer "reyes-magos-gr/views/volunteer"

	"github.com/labstack/echo/v4"
)

type HomeHandler struct{}

func (h HomeHandler) HomeViewHandler(ctx echo.Context) error {
	return lib.Render(ctx, home.Home())
}

func (h HomeHandler) SupportViewHandler(ctx echo.Context) error {
	return lib.Render(ctx, support.Support())
}

func (h HomeHandler) NotVolunteerHandler(ctx echo.Context) error {
	return lib.Render(ctx, volunteer.NotVolunteer())
}

func (h HomeHandler) VerifyEmailHandler(ctx echo.Context) error {
	return lib.Render(ctx, volunteer.VerifyEmail())
}
