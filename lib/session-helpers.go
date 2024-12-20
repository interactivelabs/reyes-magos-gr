package lib

import (
	"context"
	"encoding/json"

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

func GetSessionProfile(ctx echo.Context) (map[string]interface{}, error) {
	profileSession, err := session.Get("profile", ctx)
	if err != nil {
		return nil, err
	}

	profile := profileSession.Values["profile"]
	if profile == nil {
		return nil, nil
	}

	var profileMap map[string]interface{}
	if err := json.Unmarshal([]byte(profile.(string)), &profileMap); err != nil {
		return nil, err
	}

	return profileMap, nil
}

type contextKey string

const (
	profileKey contextKey = "profile"
)

type ProfileView struct {
	IsAdmin  bool
	Nickname string
	Email    string
	Picture  string
}

func GetProfileView(ctx echo.Context) (profile ProfileView) {
	sessionProfile, err := GetSessionProfile(ctx)
	if err != nil || sessionProfile == nil {
		return profile
	}

	if sessionProfile["dl_admin"] == "true" {
		profile.IsAdmin = true
	}

	if nickname, ok := sessionProfile["nickname"].(string); ok {
		profile.Nickname = nickname
	}

	if email, ok := sessionProfile["name"].(string); ok {
		profile.Email = email
	}

	if picture, ok := sessionProfile["picture"].(string); ok {
		profile.Picture = picture
	}

	return profile
}

func GetProfile(ctx context.Context) (profile ProfileView) {
	if profile, ok := ctx.Value(profileKey).(ProfileView); ok {
		return profile
	}
	return ProfileView{}
}
