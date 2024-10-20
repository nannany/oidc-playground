package server

import (
	"context"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"log/slog"
)

// Authorizer is a minimal implementation of op.Authorizer
// zitadel に合わせるために作った
type Authorizer struct {
}

func (a Authorizer) Storage() op.Storage {
	return Storage{}
}

func (a Authorizer) Decoder() httphelper.Decoder {
	//TODO implement me
	panic("implement me")
}

// これ何に使うんだ。。？
func (a Authorizer) Encoder() httphelper.Encoder {
	return httphelper.Encoder(oidc.NewEncoder())
}

func (a Authorizer) IDTokenHintVerifier(ctx context.Context) *op.IDTokenHintVerifier {
	//TODO implement me
	panic("implement me")
}

// Crypto returns the crypto implementation
// これでauthorization code grantのcodeを作ってる
func (a Authorizer) Crypto() op.Crypto {
	return Crypto{}
}

func (a Authorizer) RequestObjectSupported() bool {
	//TODO implement me
	panic("implement me")
}

func (a Authorizer) Logger() *slog.Logger {
	//TODO implement me
	panic("implement me")
}

var _ op.Authorizer = (*Authorizer)(nil)
