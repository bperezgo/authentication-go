package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bperezgo/authentication/auth/infra/token"
	"github.com/bperezgo/authentication/auth/infra/verify"
	"github.com/bperezgo/authentication/config"
	"github.com/bperezgo/authentication/serverApp"
	"github.com/rs/cors"
)

func main() {
	envs := config.GetConfig()
	appOpts := &serverApp.AppOpts{
		PrefixRoute: "/api",
	}
	app := serverApp.NewApp(appOpts)
	handlerToken := &token.Handler{}
	verifyToken := &verify.Handler{}
	app.SetHandler("/api/auth/token", "POST", handlerToken)
	app.SetHandler("/api/auth/verify", "POST", verifyToken)
	serverUrl := fmt.Sprintf(":%d", envs.Port)
	log.Println("[INFO] Starting server in", serverUrl)
	// Cors implementation
	handler := cors.Default().Handler(app)
	if err := http.ListenAndServe(serverUrl, handler); err != nil {
		log.Fatal(err)
	}
}
