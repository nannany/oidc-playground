package server

import (
	"context"
	"github.com/go-jose/go-jose/v4"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"time"
)

// Storage is a minimal implementation of op.Storage
// zitadel に合わせるために作った
type Storage struct {
}

func (s Storage) CreateAuthRequest(ctx context.Context, request *oidc.AuthRequest, s2 string) (op.AuthRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) AuthRequestByID(ctx context.Context, s2 string) (op.AuthRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) AuthRequestByCode(ctx context.Context, code string) (op.AuthRequest, error) {
	return codeAndAuthRequest[code], nil
}

func (s Storage) SaveAuthCode(ctx context.Context, authRequestID string, code string) error {
	codeAndAuthRequest[code] = AuthRequests[authRequestID]

	return nil
}

func (s Storage) DeleteAuthRequest(ctx context.Context, s2 string) error {
	//TODO implement me
	panic("implement me")
}

func (s Storage) CreateAccessToken(ctx context.Context, request op.TokenRequest) (accessTokenID string, expiration time.Time, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) CreateAccessAndRefreshTokens(ctx context.Context, request op.TokenRequest, currentRefreshToken string) (accessTokenID string, newRefreshTokenID string, expiration time.Time, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) TokenRequestByRefreshToken(ctx context.Context, refreshTokenID string) (op.RefreshTokenRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) TerminateSession(ctx context.Context, userID string, clientID string) error {
	//TODO implement me
	panic("implement me")
}

func (s Storage) RevokeToken(ctx context.Context, tokenOrTokenID string, userID string, clientID string) *oidc.Error {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetRefreshTokenInfo(ctx context.Context, clientID string, token string) (userID string, tokenID string, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) SigningKey(ctx context.Context) (op.SigningKey, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) SignatureAlgorithms(ctx context.Context) ([]jose.SignatureAlgorithm, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) KeySet(ctx context.Context) ([]op.Key, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	return clients[clientID], nil
}

func (s Storage) AuthorizeClientIDSecret(ctx context.Context, clientID, clientSecret string) error {
	//TODO implement me
	panic("implement me")
}

func (s Storage) SetUserinfoFromScopes(ctx context.Context, userinfo *oidc.UserInfo, userID, clientID string, scopes []string) error {
	//TODO implement me
	panic("implement me")
}

func (s Storage) SetUserinfoFromToken(ctx context.Context, userinfo *oidc.UserInfo, tokenID, subject, origin string) error {
	//TODO implement me
	panic("implement me")
}

func (s Storage) SetIntrospectionFromToken(ctx context.Context, userinfo *oidc.IntrospectionResponse, tokenID, subject, clientID string) error {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetPrivateClaimsFromScopes(ctx context.Context, userID, clientID string, scopes []string) (map[string]any, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetKeyByIDAndClientID(ctx context.Context, keyID, clientID string) (*jose.JSONWebKey, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) ValidateJWTProfileScopes(ctx context.Context, userID string, scopes []string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) Health(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

var _ op.Storage = (*Storage)(nil)
