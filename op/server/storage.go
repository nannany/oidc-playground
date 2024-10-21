package server

import (
	"context"
	"fmt"
	"github.com/go-jose/go-jose/v4"
	"github.com/google/uuid"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"time"
)

// Storage is a minimal implementation of op.Storage
// zitadel に合わせるために作った
type Storage struct {
}

func (s Storage) SetUserinfoFromRequest(ctx context.Context, userInfo *oidc.UserInfo, request op.IDTokenRequest, scopes []string) error {
	user := Users[request.GetSubject()]
	if user == nil {
		return fmt.Errorf("user not found")
	}
	for _, scope := range scopes {
		switch scope {
		case oidc.ScopeOpenID:
			userInfo.Subject = user.ID
		case oidc.ScopeEmail:
			userInfo.Email = user.Email
			userInfo.EmailVerified = oidc.Bool(user.EmailVerified)
		case oidc.ScopeProfile:
			userInfo.PreferredUsername = user.Username
			userInfo.Name = user.FirstName + " " + user.LastName
			userInfo.FamilyName = user.LastName
			userInfo.GivenName = user.FirstName
			userInfo.Locale = oidc.NewLocale(user.PreferredLanguage)
		case oidc.ScopePhone:
			userInfo.PhoneNumber = user.Phone
			userInfo.PhoneNumberVerified = user.PhoneVerified
		}
	}
	return nil
}

func (s Storage) ValidateTokenExchangeRequest(ctx context.Context, request op.TokenExchangeRequest) error {
	//TODO implement me
	panic("implement me")
}

func (s Storage) CreateTokenExchangeRequest(ctx context.Context, request op.TokenExchangeRequest) error {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetPrivateClaimsFromTokenExchangeRequest(ctx context.Context, request op.TokenExchangeRequest) (claims map[string]any, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) SetUserinfoFromTokenExchangeRequest(ctx context.Context, userinfo *oidc.UserInfo, request op.TokenExchangeRequest) error {
	//TODO implement me
	panic("implement me")
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

func (s Storage) DeleteAuthRequest(ctx context.Context, authReqID string) error {
	AuthRequests[authReqID] = nil

	return nil
}

func (s Storage) CreateAccessToken(ctx context.Context, request op.TokenRequest) (accessTokenID string, expiration time.Time, err error) {
	var applicationID string
	switch req := request.(type) {
	case *AuthRequest:
		// if authenticated for an app (auth code / implicit flow) we must save the client_id to the token
		applicationID = req.ApplicationID
	case op.TokenExchangeRequest:
		applicationID = req.GetClientID()
	}

	token, err := s.accessToken(applicationID, "", request.GetSubject(), request.GetAudience(), request.GetScopes())
	if err != nil {
		return "", time.Time{}, err
	}
	return token.ID, token.Expiration, nil
}

func (s *Storage) accessToken(applicationID, refreshTokenID, subject string, audience, scopes []string) (*Token, error) {
	token := &Token{
		ID:             uuid.NewString(),
		ApplicationID:  applicationID,
		RefreshTokenID: refreshTokenID,
		Subject:        subject,
		Audience:       audience,
		Expiration:     time.Now().Add(5 * time.Minute),
		Scopes:         scopes,
	}
	Tokens[token.ID] = token
	return token, nil
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
	return MySigningKey, nil
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
	return nil
}

func (s Storage) SetUserinfoFromScopes(ctx context.Context, userinfo *oidc.UserInfo, userID, clientID string, scopes []string) error {
	return nil
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
	// todo: ここの処理よくわかってない
	return map[string]any{}, nil
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

var _ op.TokenExchangeStorage = (*Storage)(nil)

var _ op.CanSetUserinfoFromRequest = (*Storage)(nil)
