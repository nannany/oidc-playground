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
}

var Users = make(map[string]*User)
