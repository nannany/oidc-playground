package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/zitadel/oidc/v3/pkg/op"
	"html/template"
	"log/slog"
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
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.SessionCheck)
	router.Use(middleware.IssuerSetter)

	router.Get("/", homeViewHandler)
	router.Get("/login", loginViewHandler)
	router.Get("/check_session_iframe", checkSessionIframeHandler)
	router.Get("/jwks.json", jwksHandler)
	router.Get("/auto-login", autoLoginHandler)
	router.Get("/webauthn/login/challenge", webauthnLoginChallengeHandler)
	router.Post("/register-passkey", registerPasskeyHandler)
	router.Post("/finish-register-passkey", finishRegisterPasskeyHandler)
	router.Post("/login/username", loginHandler)
	router.Post("/logout", logoutHandler)
	router.Post("/webauthn/login", webauthnLoginHandler)

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

func webauthnLoginHandler(writer http.ResponseWriter, request *http.Request) {
	// webauthn-loginからsessionDataを取得する
	webauthnLoginSession, _ := session.Store.Get(request, "webauthn-login")
	sessionDataID := webauthnLoginSession.Values["sessionDataID"]
	sessionData := server2.SessionData[sessionDataID.(string)]

	findUserHandler := func(rawID, userHandle []byte) (user webauthn.User, err error) {
		targetUserID := base64.RawURLEncoding.EncodeToString(rawID)
		userHandleStr := base64.RawURLEncoding.EncodeToString(userHandle)
		fmt.Println("targetUserID:", targetUserID)
		fmt.Println("userHandle:", userHandleStr)
		retUser := server2.WebAuthnIDUserMap[targetUserID]

		if retUser == nil {
			return nil, fmt.Errorf("user not found")
		}
		return retUser, nil
	}

	credential, err := server2.WebAuthn.FinishDiscoverableLogin(findUserHandler, *sessionData, request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// credentialをとりあえずログに入れる
	slog.Info("credential: %v", credential)

	slog.Info("login success")

	// 200返す
	writer.WriteHeader(http.StatusOK)
}

func webauthnLoginChallengeHandler(writer http.ResponseWriter, request *http.Request) {
	credentialAssertion, sessionData, err := server2.WebAuthn.BeginDiscoverableLogin()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// sessionをセッションに保存
	webauthnLoginSession, _ := session.Store.New(request, "webauthn-login")
	sessionDataID := uuid.New().String()
	server2.SessionData[sessionDataID] = sessionData
	webauthnLoginSession.Values["sessionDataID"] = sessionDataID
	err = session.Store.Save(request, writer, webauthnLoginSession)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// credentialAssertion をjsonとして、クライアントに返す
	writer.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(credentialAssertion); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func finishRegisterPasskeyHandler(writer http.ResponseWriter, request *http.Request) {
	webauthnSession, _ := session.Store.Get(request, "webauthn-session")
	sessionDataID := webauthnSession.Values["sessionDataID"]
	sessionData := server2.SessionData[sessionDataID.(string)]
	userID := webauthnSession.Values["userID"]
	user := server2.Users[userID.(string)]
	webauthnSession.Flashes("sessionData", "user")

	credential, err := server2.WebAuthn.FinishRegistration(user, *sessionData, request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	user.AddCredential(credential)

	server2.WebAuthnIDUserMap[base64.RawURLEncoding.EncodeToString(credential.ID)] = user

	// userの様子をログで見る
	slog.Info("user: %v", user)
	// credential.idの様子をbase64url decodeしてログで見る
	slog.Info("credential.id: %v", base64.RawURLEncoding.EncodeToString(credential.ID))

	// 200返す
	writer.WriteHeader(http.StatusOK)
}

func registerPasskeyHandler(writer http.ResponseWriter, request *http.Request) {
	user := server2.Users["21e204ab-b1f4-4a37-b4cf-28cffabdfe49"]
	options, sessionData, err := server2.WebAuthn.BeginRegistration(user)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// todo:session をセッションに保存
	slog.Info("session: %v", sessionData)
	webAuthnSession, _ := session.Store.New(request, "webauthn-session")
	sessionDataID := uuid.New().String()
	server2.SessionData[sessionDataID] = sessionData
	webAuthnSession.Values["sessionDataID"] = sessionDataID
	userID := "21e204ab-b1f4-4a37-b4cf-28cffabdfe49"
	webAuthnSession.Values["userID"] = userID
	err = session.Store.Save(request, writer, webAuthnSession)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// options をjsonとして、クライアントに返す
	writer.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(options); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
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
	http.ServeFile(w, r, "static/check_session.html")
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
	http.SetCookie(w, &http.Cookie{
		Name:     "op_session_state",
		Value:    "",
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "op-session",
		Value:    "",
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   1,
	})
	http.Redirect(w, r, "https://satyr-ample-supposedly.ngrok-free.app/", http.StatusFound)
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

	if id == "" {
		w.Header().Add("Set-Cookie", "op_session_state="+sid+"; Path=/; SameSite=None; Secure;")
		http.Redirect(w, r, "https://satyr-ample-supposedly.ngrok-free.app/", http.StatusFound)
		return
	} else {
		// rpにリダイレクト
		authReq := server2.AuthRequests[id]
		authReq.UserID = "21e204ab-b1f4-4a37-b4cf-28cffabdfe49" // 可変に
		// authReqののリダイレクトuriに session_state をセット
		// 本当はこんなことしたくないけど、これ以外zitadelを使ってどうやるのかわからん
		// authReqのディープコピーを作る
		copyAuthReq := authReq.DeepCopy()
		copyAuthReq.CallbackURI = authReq.CallbackURI + "?session_state=" + sid

		// cookie にop_session_stateをセットする
		w.Header().Add("Set-Cookie", "op_session_state="+sid+"; Path=/; SameSite=None; Secure;")
		op.AuthResponse(copyAuthReq, authorizer, w, r)
	}
}
