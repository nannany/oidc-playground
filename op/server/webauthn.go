package server

import "github.com/go-webauthn/webauthn/webauthn"

var WebAuthn *webauthn.WebAuthn

var sessionData = make(map[string]*webauthn.SessionData)

func init() {
	webAuthnConfig := &webauthn.Config{
		RPDisplayName: "Go WebAuthn",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:8080"},
	}
	if WebAuthn, _ = webauthn.New(webAuthnConfig); WebAuthn == nil {
		panic("webauthn new error")
	}
}
