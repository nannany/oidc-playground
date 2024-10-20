package server

import "github.com/zitadel/oidc/v3/pkg/op"

type Crypto struct {
}

func (c Crypto) Encrypt(s string) (string, error) {
	return s, nil
}

func (c Crypto) Decrypt(s string) (string, error) {
	return s, nil
}

var _ op.Crypto = (*Crypto)(nil)
