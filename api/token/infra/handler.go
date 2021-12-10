package infra

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bperezgo/authentication/config"
	"github.com/bperezgo/authentication/parser"
	"github.com/golang-jwt/jwt"
)

type TokenClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.StandardClaims
}

// Handler function to sign the token
func Handler(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] Handling the token")
	c := config.GetConfig()
	signingKey := []byte(c.AuthJwtSecret)
	request := &Request{}
	err := parser.Json(req.Body, request)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("[INFO] Request: %+v", request)
	claims := TokenClaims{
		Email: request.Email,
		Name:  request.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 15000000,
			Subject:   request.Username,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signingKey)
	if err != nil {
		log.Println("[ERROR] Failed signing the token")
	}
	response := fmt.Sprintf("{\"accessToken\":\"%s\"}", ss)
	log.Printf("[INFO] Response: %+v", response)
	res.Write([]byte(response))
}
