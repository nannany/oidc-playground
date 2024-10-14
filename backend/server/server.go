package server

import (
	"context"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"log/slog"
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
			op.ApplicationTypeWeb,
			[]string{"http://localhost:9999/notapplied"},
			[]string{"http://localhost:9999/**"},
			true,
			[]oidc.ResponseType{oidc.ResponseTypeCode}),
	}, nil

}

func (m *MyServer) Authorize(ctx context.Context, r *op.ClientRequest[oidc.AuthRequest]) (*op.Redirect, error) {
	return op.NewRedirect("http://localhost:8080/login"), nil
}

func (m *MyServer) Discovery(ctx context.Context, r *op.Request[struct{}]) (*op.Response, error) {
	return op.NewResponse(&oidc.DiscoveryConfiguration{
		Issuer: "http://localhost:8080",
	}), nil
}
