package server

import (
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"time"
)

type Client struct {
	ID                string
	applicationType   op.ApplicationType
	redirectUris      []string
	redirectUrisGlobs []string
	devMode           bool
	responseType      []oidc.ResponseType
}

// client を保存しておくメモリ領域
var clients = make(map[string]*Client)

func init() {
	clients["client_id"] = NewClient(
		"client_id",
		op.ApplicationTypeWeb,
		[]string{"http://localhost:8081/auth/callback"},
		[]string{"http://localhost:8081/**"},
		true,
		[]oidc.ResponseType{oidc.ResponseTypeCode},
	)
	clients["web"] = NewClient(
		"web",
		op.ApplicationTypeWeb,
		[]string{"http://localhost:8081/auth/callback"},
		[]string{"http://localhost:8081/**"},
		true,
		[]oidc.ResponseType{oidc.ResponseTypeCode},
	)
}

var _ op.Client = (*Client)(nil)

func NewClient(id string, applicationType op.ApplicationType, redirectUris []string, redirectUrisGlobs []string, devMode bool, responseType []oidc.ResponseType) *Client {
	return &Client{ID: id, applicationType: applicationType, redirectUris: redirectUris, redirectUrisGlobs: redirectUrisGlobs, devMode: devMode, responseType: responseType}
}

func (c Client) GetID() string {
	return c.ID
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
	return oidc.AuthMethodBasic
}

func (c Client) ResponseTypes() []oidc.ResponseType {
	return c.responseType
}

func (c Client) GrantTypes() []oidc.GrantType {
	return []oidc.GrantType{oidc.GrantTypeCode}

}

func (c Client) LoginURL(s string) string {
	return "http://op.host:8080/login?request_id=" + s
}

func (c Client) AccessTokenType() op.AccessTokenType {
	return op.AccessTokenTypeJWT
}

func (c Client) IDTokenLifetime() time.Duration {
	return 5 * time.Minute
}

func (c Client) DevMode() bool {
	return c.devMode
}

func (c Client) RestrictAdditionalIdTokenScopes() func(scopes []string) []string {
	// 全て通す
	return func(scopes []string) []string {
		return scopes
	}
}

func (c Client) RestrictAdditionalAccessTokenScopes() func(scopes []string) []string {
	// 全て通す
	return func(scopes []string) []string {
		return scopes
	}
}

func (c Client) IsScopeAllowed(scope string) bool {
	//TODO implement me
	panic("implement me")
}

func (c Client) IDTokenUserinfoClaimsAssertion() bool {
	return true
}

func (c Client) ClockSkew() time.Duration {
	return 5 * time.Second
}

func (c Client) RedirectURIGlobs() []string {
	return c.redirectUrisGlobs
}

func (c Client) PostLogoutRedirectURIGlobs() []string {
	//TODO implement me
	panic("implement me")
}
