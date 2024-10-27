package middleware

import (
	"context"
	"myoidc/session"
	"net/http"
)

// セッションがあればcontextにセットする
func SessionCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// セッションがあればcontextにセット
		// なければ次のハンドラを呼び出す
		ctx := r.Context()
		// cookieからセッションIDを取得
		opSession, _ := session.Store.Get(r, "op-session")
		userID := opSession.Values["userID"]
		if userID != nil {
			ctx = context.WithValue(ctx, "userID", userID)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
