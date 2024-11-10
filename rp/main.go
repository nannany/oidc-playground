package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"os"
	"rp/domain"
	"rp/middleware"
	"rp/session"
	"time"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

var sessionID = ""

var (
	callbackPath = "/auth/callback"
	key          = []byte("test1234test1234")
)

// ポートが8081で動くサーバーを起動する
func main() {

	router := chi.NewRouter()
	router.Use(middleware.LoggingMiddleware)

	// ハンドラを設定
	router.Get("/", middleware.SessionCheck(handler))

	clientID := "web"
	clientSecret := "secret"
	issuer := "http://op.host:8080"
	port := "8081"
	scopes := []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail}
	responseMode := "query"

	redirectURI := fmt.Sprintf("http://localhost:%v%v", port, callbackPath)
	cookieHandler := httphelper.NewCookieHandler(key, key, httphelper.WithUnsecure())

	logger := slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}),
	)
	client := &http.Client{
		Timeout:   time.Minute,
		Transport: &LoggingRoundTripper{Transport: http.DefaultTransport},
	}
	// enable outgoing request logging
	logging.EnableHTTPClient(client,
		logging.WithClientGroup("client"),
	)

	options := []rp.Option{
		rp.WithCookieHandler(cookieHandler),
		rp.WithVerifierOpts(rp.WithIssuedAtOffset(5 * time.Second)),
		rp.WithHTTPClient(client),
		rp.WithLogger(logger),
		rp.WithSigningAlgsFromDiscovery(),
	}

	// One can add a logger to the context,
	// pre-defining log attributes as required.
	ctx := logging.ToContext(context.TODO(), logger)
	provider, err := rp.NewRelyingPartyOIDC(ctx, issuer, clientID, clientSecret, redirectURI, scopes, options...)
	if err != nil {
		logrus.Fatalf("error creating provider %s", err.Error())
	}

	// generate some state (representing the state of the user in your application,
	// e.g. the page where he was before sending him to login
	state := func() string {
		return uuid.New().String()
	}

	var urlOptions []rp.URLParamOpt
	urlOptions = append(urlOptions, rp.WithResponseModeURLParam(oidc.ResponseMode(responseMode)))

	router.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("prompt") == "none" {
			urlOptions = append(urlOptions, rp.WithPromptURLParam(oidc.PromptNone))
		}

		rp.AuthURLHandler(
			state,
			provider,
			urlOptions...,
		)(w, r)
	})

	// rp.CodeExchangeCallback を定義
	f := func(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens[*oidc.IDTokenClaims], state string, rp rp.RelyingParty) {
		// tokensの情報から、userをcreate
		if domain.Users[tokens.IDTokenClaims.Subject] == nil {
			domain.Users[tokens.IDTokenClaims.Subject] = &domain.User{
				ID:            tokens.IDTokenClaims.Subject,
				Email:         tokens.IDTokenClaims.Email,
				EmailVerified: bool(tokens.IDTokenClaims.EmailVerified),
				FamilyName:    tokens.IDTokenClaims.FamilyName,
				GivenName:     tokens.IDTokenClaims.GivenName,
			}
		}

		rpSession, _ := session.Store.Get(r, "rp-session")
		rpSession.Values["user"] = tokens.IDTokenClaims.Subject
		err = rpSession.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// リクエストのクエリからsession_stateを取得して、cookieにセット
		sessionState := r.URL.Query().Get("session_state")
		if sessionState != "" {
			w.Header().Add("Set-Cookie", fmt.Sprintf("rp_session_state=%s; Path=/;", sessionState))
		}
		fmt.Print("callback success")
		http.Redirect(w, r, "/", http.StatusFound)
	}

	router.Get(callbackPath, rp.CodeExchangeHandler[*oidc.IDTokenClaims](f, provider))

	// 8081ポートでサーバーを起動
	log.Println("Server started at http://localhost:8081")
	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	// リクエストコンテキストからuserIDを取得
	userID, ok := r.Context().Value("userID").(string)
	var user *domain.User
	if !ok {
		user = nil
	} else {
		user = domain.Users[userID]
	}

	fmt.Sprintf("user: %v", user)

	var message string

	if user == nil {
		message = "Hello, Guest"
	} else {
		message = "Hello, " + user.GivenName + " " + user.FamilyName
	}

	data := map[string]string{
		"Title":   "RP Page",
		"Message": message,
	}

	// テンプレートをレンダリング
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type LoggingRoundTripper struct {
	Transport http.RoundTripper
}

func (lrt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log the request
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Printf("Error dumping request: %v", err)
	}
	log.Printf("Request:\n%s", string(reqDump))

	// Perform the request
	resp, err := lrt.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Log the response
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Printf("Error dumping response: %v", err)
	}
	log.Printf("Response:\n%s", string(respDump))

	return resp, nil
}
