package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bperezgo/authentication/auth/infra/token"
	"github.com/bperezgo/authentication/auth/infra/verify"
	"github.com/bperezgo/authentication/config"
)

func main() {
	envs := config.GetConfig()
	appOpts := &AppOpts{
		PrefixRoute: "/api",
	}
	app := NewApp(appOpts)
	app.SetHandler("/api/auth/token", "POST", token.Handler)
	app.SetHandler("/api/auth/verify", "POST", verify.Handler)
	serverUrl := fmt.Sprintf(":%d", envs.Port)
	log.Println("[INFO] Starting server in", serverUrl)
	if err := http.ListenAndServe(serverUrl, app); err != nil {
		log.Fatal(err)
	}
}
