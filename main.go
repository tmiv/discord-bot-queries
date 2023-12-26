package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-collections/collections/set"
	"github.com/rs/cors"
	"github.com/xenitab/go-oidc-middleware/oidctoken"
	"github.com/xenitab/go-oidc-middleware/options"
)

type EmailClaims struct {
	Email string `json:"email"`
}

func emailAllowedValidator() options.ClaimsValidationFn[EmailClaims] {
	allow_env := strings.Split(os.Getenv("SECURITY_ALLOW"), ",")
	allow_list := set.New()
	for _, s := range allow_env {
		allow_list.Insert(s)
	}
	return func(claims *EmailClaims) error {
		if allow_list.Has(claims.Email) {
			return nil
		} else {
			return fmt.Errorf("%s is not on the allow list", claims.Email)
		}
	}
}

func SetupOIDCMiddleware() func(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	oidctok, err := oidctoken.New[EmailClaims](
		emailAllowedValidator(),
		options.WithIssuer(os.Getenv("SECURITY_ISSUER")),
		options.WithRequiredTokenType("JWT"),
		options.WithRequiredAudience(os.Getenv("SECURITY_AUDIENCE")),
	)
	if err != nil {
		log.Fatalf("Error creating token parser %+v\n", err)
	}
	oidcmiddle := func(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				log.Printf("No bearer %s\n", auth)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			_, err = oidctok.ParseToken(r.Context(), auth[7:])
			if err != nil {
				log.Printf("Unauthorized %v\n", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next(w, r)
		}
	}

	return oidcmiddle
}

func setupcors() *cors.Cors {
	options := cors.Options{
		AllowedMethods:   []string{http.MethodPost},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization"},
	}
	originsenv := os.Getenv("CORS_ORIGINS")
	if len(originsenv) > 0 {
		origins := strings.Split(originsenv, "'")
		options.AllowedOrigins = origins
	}
	return cors.New(options)
}

func passthrough(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) { next(w, r) }
}

func main() {
	err := authBot()
	if err != nil {
		log.Fatalf("Could not init bot %v\n", err)
		return
	}
	defer closeBot()

	var middleware func(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request)
	if len(os.Getenv("SKIP_OIDC")) > 0 {
		middleware = passthrough
	} else {
		middleware = SetupOIDCMiddleware()
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/VerifyMembership", middleware(verifyMembership))
	corsobj := setupcors()
	handler := corsobj.Handler(mux)
	fmt.Println("Starting Listen")
	err = http.ListenAndServe("0.0.0.0:8080", handler)
	fmt.Printf("Listen Error %v\n", err)
}
