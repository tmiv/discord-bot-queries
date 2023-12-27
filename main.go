package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rs/cors"
	ovm "github.com/tmiv/oidc-verify-middleware"
)

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
		middleware = ovm.SetupOIDCMiddleware("")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/VerifyMembership", middleware(verifyMembership))
	corsobj := setupcors()
	handler := corsobj.Handler(mux)
	fmt.Println("Starting Listen")
	err = http.ListenAndServe("0.0.0.0:8080", handler)
	fmt.Printf("Listen Error %v\n", err)
}
