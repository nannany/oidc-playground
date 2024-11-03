package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/zitadel/oidc/v3/pkg/op"
	"html/template"
	"myoidc/middleware"
	server2 "myoidc/server"
	"myoidc/session"
	"net/http"
)

var homeTmpl = template.Must(template.ParseFiles("templates/home.html"))
var loginTmpl = template.Must(template.ParseFiles("templates/login.html"))

var authorizer = server2.Authorizer{}

// 8080で動くサーバーを起動する
func main() {
	fmt.Println("openid provider started!")

	router := chi.NewRouter()
	router.Use(middleware.SessionCheck)
	router.Use(middleware.IssuerSetter)

	router.Get("/", homeViewHandler)
	router.Get("/login", loginViewHandler)
	router.Get("/check_session_iframe", checkSessionIframeHandler)
	router.Get("/jwks.json", jwksHandler)
	router.Get("/auto-login", autoLoginHandler)
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

func autoLoginHandler(w http.ResponseWriter, r *http.Request) {
	// クエリパラメータからauthReqIDを取得
	authReqID := r.URL.Query().Get("auth_req_id")
	authReq := server2.AuthRequests[authReqID]

	sid := r.Context().Value("sid").(string)

	copiedAuthReq := authReq.DeepCopy()
	copiedAuthReq.CallbackURI = authReq.CallbackURI + "?session_state=" + sid
	op.AuthResponse(copiedAuthReq, authorizer, w, r)
}

func jwksHandler(writer http.ResponseWriter, request *http.Request) {
	w := writer
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`
    {
    	"keys": [
    		{
    			"kty": "RSA",
    			"n": "5srSQZgeolXrjpTvw1OuHxBrHiBKnBEOxeeOgBDaB_61Dm1nr39rnbjd7CmuVel9o1CQof26741AoqxFxDAzc1KtnG2pysT32kcKVLBYQYSyXl860jrMXBgs-eR2Gz_YJl5UmmMvexYmnJ1CAhDUxMK23MeeR0_llTUIRDPrE1JFgE033gvUF8PfNxSUzeI5FHu6PjbLrwiatg3sOhUAkxQhC5IPGJoSVuS0_taU72lRoSEKT2Ij32HnhLWx7dAZ_PXcSZGU3L86AGksenF-bfDes6_OXIWkCBtlcpXGo51WNWzmVgX1KBVe48SCWwO9qIr8F6oRNe0zxcIvaSWHKpfMw711uF8OT8XpF9jOvlMxXGOASpAJ8eDVh4DK4YHfG4GFg4mlzQ6wr7_MHl8yXLj5v_-03XS3-AzskLs86haHi91U6zoA2zGkQ6f_KsBa5Mi7Yn9XkjT3LqIdE2Eq6PzLkXa0_BPyoA4yu1AQiZ0UneNCZpxqD_1UVzU2ZmoyvNprAd1Y5RK7pimWx8NAEkcZfLg3OjsQvxho4l0YeyqZPrnmYy2G61BvCWkgzpjoHIxn9IJgdXsS80ugJKOWF-hfKUYwyW5iWuO285WvZbF_jSoqfGvKk21bsyf5_4Pj0i_5lY5OmTrYnDHhGKNcO_FrZXKHvEVTLFC7h1FJOo8",
    			"e": "AQAB",
    			"alg": "RS256",
    			"use": "sig",
    			"kid": "avUja_OmJ6soJ6KUnmM_IWoPLxny3Ph-uWLZnFxrpuE"
    		}
    	]
    }
`))
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
	opSession, _ := session.Store.New(r, "op-session")
	opSession.Values["userID"] = ""
	err := opSession.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// cookie から op_session_state を削除する
	w.Header().Set("Set-Cookie", "op_session_state=; Path=/; Secure; SameSite=None; Max-Age=0")
	w.Header().Set("Set-Cookie", "op-session=; Path=/; Secure; SameSite=None; Max-Age=0")
	http.Redirect(w, r, "http://localhost:8080/", http.StatusFound)
}

func homeViewHandler(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := r.Context().Value("sid").(string)
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
	opSession, _ := session.Store.Get(r, "op-session")
	opSession.Values["userID"] = "21e204ab-b1f4-4a37-b4cf-28cffabdfe49"
	sid := uuid.New().String()
	opSession.Values["sid"] = sid
	err := opSession.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = session.Store.Save(r, w, opSession)

	// rpにリダイレクト
	authReq := server2.AuthRequests[id]
	authReq.UserID = "21e204ab-b1f4-4a37-b4cf-28cffabdfe49" // 可変に
	// authReqののリダイレクトuriに session_state をセット
	// 本当はこんなことしたくないけど、これ以外zitadelを使ってどうやるのかわからん
	// authReqのディープコピーを作る
	copyAuthReq := authReq.DeepCopy()
	copyAuthReq.CallbackURI = authReq.CallbackURI + "?session_state=" + sid

	// cookie にop_session_stateをセットする
	w.Header().Set("Set-Cookie", "op_session_state="+sid+"; Path=/; Secure; SameSite=None")
	op.AuthResponse(copyAuthReq, authorizer, w, r)
}
