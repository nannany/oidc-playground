package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/zitadel/oidc/v3/pkg/op"
	"html/template"
	server2 "myoidc/server"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("templates/login.html"))

// 8080で動くサーバーを起動する
func main() {
	fmt.Println("openid provider started!")

	router := chi.NewRouter()

	router.HandleFunc("/login", loginViewHandler)

	router.Mount("/", op.RegisterServer(server2.NewMyServer(), *op.DefaultEndpoints))

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func loginViewHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title":   "Go Template Example",
		"Message": "This is a static HTML page rendered with Go!",
	}

	// テンプレートをレンダリング
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
