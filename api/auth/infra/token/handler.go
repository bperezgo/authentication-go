package token

import (
	"log"
	"net/http"
	"time"

	"github.com/bperezgo/authentication/config"
	"github.com/bperezgo/authentication/pkg/parser"
	"github.com/bperezgo/authentication/serverApp"
	"github.com/golang-jwt/jwt"
)

type TokenClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.StandardClaims
}

type Handler struct{}

// Handler function to sign the token
func (h *Handler) Handle(res http.ResponseWriter, req *http.Request) (*serverApp.SuccessResponse, *serverApp.ErrorResponse) {
	log.Println("[INFO] Handling the token")
	c := config.GetConfig()
	signingKey := []byte(c.AuthJwtSecret)
	request := &Request{}
	err := parser.Json(req.Body, request)
	if err != nil {
		return nil, &serverApp.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Stack:      err.Error(),
		}
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
	log.Printf("[INFO] Response: %+v", response)
	return &serverApp.SuccessResponse{
		StatusCode: http.StatusOK,
		Body:       response,
	}, nil
}
