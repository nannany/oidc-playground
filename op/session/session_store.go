package session

import "github.com/gorilla/sessions"

var Store = sessions.NewCookieStore([]byte("OP_SESSION_KEY")) // todo: 本当は秘密にすべきkeyを渡す

func init() {
	Store.Options.Secure = false
	Store.Options.SameSite = 3 // strict
}
