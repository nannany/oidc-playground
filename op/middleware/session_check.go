package middleware

import (
	"context"
	"net/http"
)

// セッションがあればcontextにセットする
func SessionCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// セッションがあればcontextにセット
		// なければ次のハンドラを呼び出す
		ctx := r.Context()
		// cookieからセッションIDを取得

		cookie, err := r.Cookie("op-session")
		if err == nil {
			ctx = context.WithValue(ctx, "sessionID", cookie.Value)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
