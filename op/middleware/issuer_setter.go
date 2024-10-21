package middleware

import (
	"github.com/zitadel/oidc/v3/pkg/op"
	"net/http"
)

func IssuerSetter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// requestのschemeとhostを取得
		r = r.WithContext(op.ContextWithIssuer(r.Context(), "http://"+r.Host)) // schemeを取る方法がわからん
		next.ServeHTTP(w, r)
	})
}
