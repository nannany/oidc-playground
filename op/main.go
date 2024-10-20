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

var authorizer = server2.Authorizer{}

// 8080で動くサーバーを起動する
func main() {
	fmt.Println("openid provider started!")

	router := chi.NewRouter()
	router.Use(middleware.SessionCheck)

	router.Get("/", homeViewHandler)
	router.Get("/login", loginViewHandler)
	router.Get("/check_session_iframe", checkSessionIframeHandler)
	router.Post("/login/username", loginHandler)
	router.Post("/logout", logoutHandler)

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

func checkSessionIframeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	js := `
    window.addEventListener("message", receiveMessage, false);

    function receiveMessage(e){ // e.data has client_id and session_state

      var client_id = e.data.substr(0, e.data.lastIndexOf(' '));
      var session_state = e.data.substr(e.data.lastIndexOf(' ') + 1);
      var salt = session_state.split('.')[1];

      // if message is syntactically invalid
      //     postMessage('error', e.origin) and return

      // if message comes an unexpected origin
      //     postMessage('error', e.origin) and return

      // get_op_user_agent_state() is an OP defined function
      // that returns the User Agent's login status at the OP.
      // How it is done is entirely up to the OP.
      var opuas = get_op_user_agent_state();

      // Here, the session_state is calculated in this particular way,
      // but it is entirely up to the OP how to do it under the
      // requirements defined in this specification.
      var ss = CryptoJS.SHA256(client_id + ' ' + e.origin + ' ' +
        opuas + ' ' + salt) + "." + salt;

      var stat = '';
      if (session_state === ss) {
        stat = 'unchanged';
      } else {
        stat = 'changed';
      }

      e.source.postMessage(stat, e.origin);
    };
	`
	_, _ = w.Write([]byte(js))
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// セッションを削除する
	opSessionID = ""
	http.SetCookie(w, &http.Cookie{ // クッキーをセット
		Name:     "op-session",
		Value:    opSessionID,
		Secure:   false,
		HttpOnly: true,
		Path:     "/",
	})

	// rpにリダイレクト
	http.Redirect(w, r, "http://localhost:8080/", http.StatusFound)
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
	// クエリパラメータのrequest_idを取得
	requestID := r.URL.Query().Get("request_id")
	// IDに上記で取得したrequestIDをセット
	data := map[string]string{
		"ID": requestID,
	}

	// テンプレートをレンダリング
	if err := loginTmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	id := r.Form.Get("id")
	username := r.Form.Get("username")

	fmt.Println("id:", id)
	fmt.Println("username:", username)

	// todo: username, password検証

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
	authReq := server2.AuthRequests[id]
	authReq.UserID = "21e204ab-b1f4-4a37-b4cf-28cffabdfe49" // 可変に
	op.AuthResponse(authReq, authorizer, w, r)
	http.Redirect(w, r, "http://localhost:8081/auth/callback", http.StatusFound)
}
