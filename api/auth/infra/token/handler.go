package token

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bperezgo/authentication/config"
	"github.com/bperezgo/authentication/pkg/parser"
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
	expiresAt := time.Now().Unix() + 15000000
	claims := TokenClaims{
		Email: request.Email,
		Name:  request.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			Subject:   request.Username,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(signingKey)
	if err != nil {
		log.Println("[ERROR] Failed signing the token")
	}
	response := Response{
		AccessToken: accessToken,
	}
	resByte, _ := json.Marshal(response)
	log.Printf("[INFO] Response: %+v", response)
	res.Write([]byte(resByte))
}
