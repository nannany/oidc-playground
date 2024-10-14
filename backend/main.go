package main

import (
	"fmt"
	"github.com/zitadel/oidc/v3/pkg/op"
	server2 "myoidc/server"
	"net/http"
)

// 8080で動くサーバーを起動する
func main() {
	fmt.Println("openid provider started!")

	server := &http.Server{
		Addr:    ":8080",
		Handler: op.RegisterServer(server2.NewMyServer(), *op.DefaultEndpoints),
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
