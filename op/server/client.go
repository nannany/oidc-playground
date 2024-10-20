package server

import (
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"time"
)

type Client struct {
	applicationType   op.ApplicationType
	redirectUris      []string
	redirectUrisGlobs []string
	devMode           bool
	responseType      []oidc.ResponseType
}

func NewClient(applicationType op.ApplicationType, redirectUris []string, redirectUrisGlobs []string, devMode bool, responseType []oidc.ResponseType) *Client {
	return &Client{applicationType: applicationType, redirectUris: redirectUris, redirectUrisGlobs: redirectUrisGlobs, devMode: devMode, responseType: responseType}
}

func (c Client) GetID() string {
	//TODO implement me
	panic("implement me")
}

func (c Client) RedirectURIs() []string {
	return c.redirectUris
}

func (c Client) PostLogoutRedirectURIs() []string {
	//TODO implement me
	panic("implement me")
}

func (c Client) ApplicationType() op.ApplicationType {
	return c.applicationType
}

func (c Client) AuthMethod() oidc.AuthMethod {
	//TODO implement me
	panic("implement me")
}

func (c Client) ResponseTypes() []oidc.ResponseType {
	return c.responseType
}

func (c Client) GrantTypes() []oidc.GrantType {
	//TODO implement me
	panic("implement me")
}

func (c Client) LoginURL(s string) string {
	//TODO implement me
	panic("implement me")
}

func (c Client) AccessTokenType() op.AccessTokenType {
	//TODO implement me
	panic("implement me")
}

func (c Client) IDTokenLifetime() time.Duration {
	//TODO implement me
	panic("implement me")
}

func (c Client) DevMode() bool {
	return c.devMode
}

func (c Client) RestrictAdditionalIdTokenScopes() func(scopes []string) []string {
	//TODO implement me
	panic("implement me")
}

func (c Client) RestrictAdditionalAccessTokenScopes() func(scopes []string) []string {
	//TODO implement me
	panic("implement me")
}

func (c Client) IsScopeAllowed(scope string) bool {
	//TODO implement me
	panic("implement me")
}

func (c Client) IDTokenUserinfoClaimsAssertion() bool {
	//TODO implement me
	panic("implement me")
}

func (c Client) ClockSkew() time.Duration {
	//TODO implement me
	panic("implement me")
}

func (c Client) RedirectURIGlobs() []string {
	return c.redirectUrisGlobs
}

func (c Client) PostLogoutRedirectURIGlobs() []string {
	//TODO implement me
	panic("implement me")
}
