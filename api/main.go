package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bperezgo/authentication/auth/infra/token"
	"github.com/bperezgo/authentication/auth/infra/verify"
	"github.com/bperezgo/authentication/config"
	"github.com/bperezgo/authentication/serverApp"
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
	if err := http.ListenAndServe(serverUrl, app); err != nil {
		log.Fatal(err)
	}
}
