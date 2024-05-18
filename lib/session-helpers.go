package lib

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func SetCookieSession(ctx echo.Context, name string, value string) error {
	s, err := session.Get(name, ctx)
	if err != nil {
		return err
	}

	s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   36000,
		HttpOnly: true,
		Secure:   true,
	}
	s.Values[name] = value
	if err := s.Save(ctx.Request(), ctx.Response()); err != nil {
		return err
	}

	return nil
}

func DeleteCookieSession(ctx echo.Context, name string) error {
	s, err := session.Get(name, ctx)
	if err != nil {
		return err
	}

	s.Options.MaxAge = -1
	if err := s.Save(ctx.Request(), ctx.Response()); err != nil {
		return err
	}

	return nil
}

func GetProfile(ctx echo.Context) (map[string]interface{}, error) {
	profileSession, err := session.Get("profile", ctx)
	if err != nil {
		return nil, err
	}

	profile := profileSession.Values["profile"].(map[string]interface{})
	return profile, nil
}
