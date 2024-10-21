package server

import (
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/text/language"
)

type User struct {
	ID                string
	Email             string
	EmailVerified     oidc.Bool
	Username          string
	FirstName         string
	LastName          string
	PreferredLanguage language.Tag
	Phone             string
	PhoneVerified     bool
	Password          string
}

var Users = make(map[string]*User)

func init() {
	Users["21e204ab-b1f4-4a37-b4cf-28cffabdfe49"] = &User{
		ID:                "21e204ab-b1f4-4a37-b4cf-28cffabdfe49",
		Email:             "test@example.com",
		EmailVerified:     true,
		Username:          "test",
		FirstName:         "Test",
		LastName:          "User",
		PreferredLanguage: language.Japanese,
		Phone:             "+1234567890",
		PhoneVerified:     true,
		Password:          "password",
	}
}
