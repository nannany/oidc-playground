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

// op.AuthRequestを実装していることを確認
var _ op.AuthRequest = (*AuthRequest)(nil)

func (a AuthRequest) GetID() string {
	//TODO implement me
	panic("implement me")
}

func (a AuthRequest) GetACR() string {
	//TODO implement me
	panic("implement me")
}

func (a AuthRequest) GetAMR() []string {
	//TODO implement me
	panic("implement me")
}

func (a AuthRequest) GetAudience() []string {
	//TODO implement me
	panic("implement me")
}

func (a AuthRequest) GetAuthTime() time.Time {
	//TODO implement me
	panic("implement me")
}

func (a AuthRequest) GetClientID() string {
	//TODO implement me
	panic("implement me")
}

func (a AuthRequest) GetCodeChallenge() *oidc.CodeChallenge {
	//TODO implement me
	panic("implement me")
}

func (a AuthRequest) GetNonce() string {
	//TODO implement me
	panic("implement me")
}

func (a AuthRequest) GetRedirectURI() string {
	//TODO implement me
	panic("implement me")
}

func (a AuthRequest) GetResponseType() oidc.ResponseType {
	//TODO implement me
	panic("implement me")
}

func (a AuthRequest) GetResponseMode() oidc.ResponseMode {
	//TODO implement me
	panic("implement me")
}

func (a AuthRequest) GetScopes() []string {
	//TODO implement me
	panic("implement me")
}

func (a AuthRequest) GetState() string {
	//TODO implement me
	panic("implement me")
}

func (a AuthRequest) GetSubject() string {
	//TODO implement me
	panic("implement me")
}

func (a AuthRequest) Done() bool {
	//TODO implement me
	panic("implement me")
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