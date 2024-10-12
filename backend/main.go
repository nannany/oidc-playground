package main

import (
	"fmt"
	"github.com/zitadel/oidc/v3/pkg/op"
	"net/http"
)

// 8080で動くサーバーを起動する
func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})
	http.ListenAndServe(":8080", nil)
}

type MyServer struct {
	op.UnimplementedServer
}

func NewMyServer() *MyServer {
	return &MyServer{UnimplementedServer: op.UnimplementedServer{}}
}
