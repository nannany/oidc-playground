package server

import (
	"github.com/go-webauthn/webauthn/webauthn"
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

// User は webauthn.User インターフェースを実装することを示す
var _ webauthn.User = (*User)(nil)

func (u User) WebAuthnID() []byte {
	return []byte(u.ID)
}

func (u User) WebAuthnName() string {
	return u.Username
}

func (u User) WebAuthnDisplayName() string {
	return u.Username
}

func (u User) WebAuthnCredentials() []webauthn.Credential {
	return nil
}

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
