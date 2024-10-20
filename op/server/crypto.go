package server

import "github.com/zitadel/oidc/v3/pkg/op"

type Crypto struct {
}

func (c Crypto) Encrypt(s string) (string, error) {
	//TODO implement me
	panic("implement me")

}

func (c Crypto) Decrypt(s string) (string, error) {
	//TODO implement me
	panic("implement me")
}

var _ op.Crypto = (*Crypto)(nil)
