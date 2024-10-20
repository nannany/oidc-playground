package server

import "github.com/zitadel/oidc/v3/pkg/op"

type TokenCreator struct {
}

func (t TokenCreator) Storage() op.Storage {
	return Storage{}
}

func (t TokenCreator) Crypto() op.Crypto {
	return Crypto{}
}

var _ op.TokenCreator = (*TokenCreator)(nil)
