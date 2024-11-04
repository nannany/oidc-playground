package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"log/slog"
	"time"
)

type MyServer struct {
	op.UnimplementedServer
}

func NewMyServer() *MyServer {
	return &MyServer{UnimplementedServer: op.UnimplementedServer{}}
}

func (m *MyServer) Health(ctx context.Context, r *op.Request[struct{}]) (*op.Response, error) {
	return op.NewResponse("Health check!"), nil
}

func (m *MyServer) VerifyAuthRequest(ctx context.Context, r *op.Request[oidc.AuthRequest]) (*op.ClientRequest[oidc.AuthRequest], error) {

	if r.Data.ClientID == "" {
		slog.Info("client_id is missing")
	}

	return &op.ClientRequest[oidc.AuthRequest]{
		Request: r,
		Client: NewClient(
			"client_id",
			op.ApplicationTypeWeb,
			[]string{"http://localhost:8081/auth/callback"},
			[]string{"http://localhost:8081/**"},
			true,
			[]oidc.ResponseType{oidc.ResponseTypeCode}),
	}, nil

}

func (m *MyServer) Authorize(ctx context.Context, r *op.ClientRequest[oidc.AuthRequest]) (*op.Redirect, error) {
	authReq := r.Data

	// todo: validate request

	// userIDはtoken hint があれば取得できる
	request := authRequestToInternal(authReq, "")

	request.ID = uuid.NewString()

	// メモリにauthorization epに飛んできたリクエストを保存する
	AuthRequests[request.ID] = request

	// promptがnoneまたはpromptが提供されていない場合は、コンテキストにuserIDがあればcallbackさせる
	if authReq.Prompt == nil || len(authReq.Prompt) == 0 || authReq.Prompt[0] == oidc.PromptNone {
		userID := ctx.Value("userID")
		if userID != nil && userID != "" {
			request.UserID = userID.(string)
			// https://github.com/zitadel/oidc/discussions/669
			return op.NewRedirect("http://op.host:8080/auto-login?auth_req_id=" + request.ID), nil
		}
	}

	return op.NewRedirect(r.Client.LoginURL(request.ID)), nil
}

// リクエストを保存形式に合わせるべく変換
func authRequestToInternal(authReq *oidc.AuthRequest, userID string) *AuthRequest {
	return &AuthRequest{
		CreationDate:  time.Now(),
		ApplicationID: authReq.ClientID,
		CallbackURI:   authReq.RedirectURI,
		TransferState: authReq.State,
		Prompt:        PromptToInternal(authReq.Prompt),
		UiLocales:     authReq.UILocales,
		LoginHint:     authReq.LoginHint,
		MaxAuthAge:    MaxAgeToInternal(authReq.MaxAge),
		UserID:        userID,
		Scopes:        authReq.Scopes,
		ResponseType:  authReq.ResponseType,
		ResponseMode:  authReq.ResponseMode,
		Nonce:         authReq.Nonce,
		done:          false,
	}
}

func (m *MyServer) Discovery(ctx context.Context, r *op.Request[struct{}]) (*op.Response, error) {
	return op.NewResponse(&oidc.DiscoveryConfiguration{
		Issuer:                "http://op.host:8080",
		AuthorizationEndpoint: "http://op.host:8080/authorize",
		TokenEndpoint:         "http://op.host:8080/oauth/token",
		JwksURI:               "http://op.host:8080/jwks.json",
	}), nil
}

func (m *MyServer) CodeExchange(ctx context.Context, r *op.ClientRequest[oidc.AccessTokenRequest]) (*op.Response, error) {
	authReq, err := op.AuthRequestByCode(ctx, Storage{}, r.Data.Code)
	if err != nil {
		return nil, err
	}
	if r.Client.AuthMethod() == oidc.AuthMethodNone || r.Data.CodeVerifier != "" {
		if err = op.AuthorizeCodeChallenge(r.Data.CodeVerifier, authReq.GetCodeChallenge()); err != nil {
			return nil, err
		}
	}
	if r.Data.RedirectURI != authReq.GetRedirectURI() {
		return nil, oidc.ErrInvalidGrant().WithDescription("redirect_uri does not correspond")
	}
	resp, err := op.CreateTokenResponse(ctx, authReq, r.Client, TokenCreator{}, true, r.Data.Code, "")
	if err != nil {
		return nil, err
	}
	return op.NewResponse(resp), nil
}

func (m *MyServer) VerifyClient(ctx context.Context, r *op.Request[op.ClientCredentials]) (op.Client, error) {
	if oidc.GrantType(r.Form.Get("grant_type")) == oidc.GrantTypeClientCredentials {
		panic("not implemented")
	}

	if r.Data.ClientAssertionType == oidc.ClientAssertionTypeJWTAssertion {
		panic("not implemented")
	}
	storage := Storage{}
	client, err := storage.GetClientByClientID(ctx, r.Data.ClientID)
	if err != nil {
		return nil, oidc.ErrInvalidClient().WithParent(err)
	}

	switch client.AuthMethod() {
	case oidc.AuthMethodNone:
		return client, nil
	case oidc.AuthMethodPrivateKeyJWT:
		return nil, oidc.ErrInvalidClient().WithDescription("private_key_jwt not allowed for this client")
	case oidc.AuthMethodPost:
		panic("not implemented")
	}

	err = op.AuthorizeClientIDSecret(ctx, r.Data.ClientID, r.Data.ClientSecret, storage)
	if err != nil {
		return nil, err
	}

	return client, nil
}
