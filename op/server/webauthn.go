package server

import "github.com/go-webauthn/webauthn/webauthn"

var WebAuthn *webauthn.WebAuthn

var SessionData = make(map[string]*webauthn.SessionData)

func init() {
	webAuthnConfig := &webauthn.Config{
		RPDisplayName: "Go WebAuthn",
		RPID:          "satyr-ample-supposedly.ngrok-free.app",
		RPOrigins:     []string{"https://satyr-ample-supposedly.ngrok-free.app"},
	}
	if WebAuthn, _ = webauthn.New(webAuthnConfig); WebAuthn == nil {
		panic("webauthn new error")
	}
}
