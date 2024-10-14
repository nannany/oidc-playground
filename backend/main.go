package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/zitadel/oidc/v3/pkg/op"
	"html/template"
	"myoidc/middleware"
	server2 "myoidc/server"
	"net/http"
)

var homeTmpl = template.Must(template.ParseFiles("templates/home.html"))
var loginTmpl = template.Must(template.ParseFiles("templates/login.html"))

var opSessionID = ""

// 8080で動くサーバーを起動する
func main() {
	fmt.Println("openid provider started!")

	router := chi.NewRouter()
	router.Use(middleware.SessionCheck)

	router.Get("/", homeViewHandler)
	router.Get("/login", loginViewHandler)
	router.Post("/login/username", loginHandler)

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

func homeViewHandler(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := r.Context().Value("sessionID").(string)
	if !ok || sessionID == "" {
		sessionID = "empty"
	}
	data := map[string]string{
		"SessionID": sessionID,
	}

	// テンプレートをレンダリング
	if err := homeTmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginViewHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title":   "Go Template Example",
		"Message": "This is a static HTML page rendered with Go!",
	}

	// テンプレートをレンダリング
	if err := loginTmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	username := r.Form.Get("username")

	fmt.Println("username:", username)

	// クッキーにセッションをセット
	opSessionID = uuid.New().String()
	http.SetCookie(w, &http.Cookie{ // クッキーをセット
		Name:     "op-session",
		Value:    opSessionID,
		Secure:   false,
		HttpOnly: true,
		Path:     "/",
	})

	// rpにリダイレクト
	http.Redirect(w, r, "http://localhost:8081/auth/callback", http.StatusFound)
}
