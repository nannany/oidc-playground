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
	PasskeyCredential []*webauthn.Credential
}

var Users = make(map[string]*User)

var WebAuthnIDUserMap = make(map[byte]*User)

// User は webauthn.User インターフェースを実装することを示す
var _ webauthn.User = (*User)(nil)

func (u *User) WebAuthnID() []byte {
	return []byte(u.ID)
}

func (u *User) WebAuthnName() string {
	return u.Username
}

func (u *User) WebAuthnDisplayName() string {
	return u.Username
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return nil
}

func (u *User) AddCredential(credential *webauthn.Credential) {
	u.PasskeyCredential = append(u.PasskeyCredential, credential)
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
