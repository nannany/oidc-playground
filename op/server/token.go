package server

import "time"

type Token struct {
	ID             string
	ApplicationID  string
	Subject        string
	RefreshTokenID string
	Audience       []string
	Expiration     time.Time
	Scopes         []string
}

// Token を保存しておくメモリ領域
var Tokens = make(map[string]*Token)
