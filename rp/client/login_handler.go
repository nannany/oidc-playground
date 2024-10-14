package client

import "net/http"

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// レスポンスとしてhello worldを返す
	w.Write([]byte("hello world"))
}
