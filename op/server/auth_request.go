package server

import (
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"golang.org/x/text/language"
	"time"
)

type AuthRequest struct {
	ID            string
	CreationDate  time.Time
	ApplicationID string
	CallbackURI   string
	TransferState string
	Prompt        []string
	UiLocales     []language.Tag
	LoginHint     string
	MaxAuthAge    *time.Duration
	UserID        string
	Scopes        []string
	ResponseType  oidc.ResponseType
	ResponseMode  oidc.ResponseMode
	Nonce         string

	done     bool
	authTime time.Time
}

// authorization reqを保存しておくメモリ領域
var AuthRequests = make(map[string]*AuthRequest)

// codeとauth reqを保存しておくメモリ領域
var codeAndAuthRequest = make(map[string]*AuthRequest)

// op.AuthRequestを実装していることを確認
var _ op.AuthRequest = (*AuthRequest)(nil)

func (a AuthRequest) GetID() string {
	return a.ID
}

func (a AuthRequest) GetACR() string {
	return "urn:mace:incommon:iap:bronze"
}

func (a AuthRequest) GetAMR() []string {
	return []string{"pwd"}
}

func (a AuthRequest) GetAudience() []string {
	return []string{a.ApplicationID}
}

func (a AuthRequest) GetAuthTime() time.Time {
	return a.authTime
}

func (a AuthRequest) GetClientID() string {
	return a.ApplicationID
}

func (a AuthRequest) GetCodeChallenge() *oidc.CodeChallenge {
	//TODO implement me
	panic("implement me")
}

func (a AuthRequest) GetNonce() string {
	return a.Nonce
}

func (a AuthRequest) GetRedirectURI() string {
	return a.CallbackURI
}

func (a AuthRequest) GetResponseType() oidc.ResponseType {
	return a.ResponseType
}

func (a AuthRequest) GetResponseMode() oidc.ResponseMode {
	return a.ResponseMode
}

func (a AuthRequest) GetScopes() []string {
	return a.Scopes
}

func (a AuthRequest) GetState() string {
	return a.TransferState
}

func (a AuthRequest) GetSubject() string {
	return a.UserID
}

func (a AuthRequest) Done() bool {
	return a.done
}

func PromptToInternal(oidcPrompt oidc.SpaceDelimitedArray) []string {
	prompts := make([]string, 0, len(oidcPrompt))
	for _, oidcPrompt := range oidcPrompt {
		switch oidcPrompt {
		case oidc.PromptNone,
			oidc.PromptLogin,
			oidc.PromptConsent,
			oidc.PromptSelectAccount:
			prompts = append(prompts, oidcPrompt)
		}
	}
	return prompts
}

func MaxAgeToInternal(maxAge *uint) *time.Duration {
	if maxAge == nil {
		return nil
	}
	dur := time.Duration(*maxAge) * time.Second
	return &dur
}

func (a AuthRequest) DeepCopy() AuthRequest {
	var cp AuthRequest
	cp = a
	cp.Prompt = make([]string, len(a.Prompt))
	copy(cp.Prompt, a.Prompt)
	cp.Scopes = make([]string, len(a.Scopes))
	copy(cp.Scopes, a.Scopes)
	cp.UiLocales = make([]language.Tag, len(a.UiLocales))
	copy(cp.UiLocales, a.UiLocales)

	return cp
}
