package verify

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bperezgo/authentication/config"
	"github.com/bperezgo/authentication/serverApp"
	"github.com/golang-jwt/jwt"
)

type Handler struct{}

// Handler function to verify the token
func (h *Handler) Handle(ctx context.Context, req *http.Request) (*serverApp.SuccessResponse, *serverApp.ErrorResponse) {
	log.Println("[INFO] Handling the token")
	c := config.GetConfig()
	signingKey := []byte(c.AuthJwtSecret)
	urlValues := req.URL.Query()
	accessToken := urlValues.Get("access_token")
	if accessToken == "" {
		// res.WriteHeader(http.StatusBadRequest)
		response := Response{
			Message: "invalid access_token",
		}
		log.Printf("[INFO] Response: %+v", response)
		return nil, &serverApp.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Body:       response,
		}
	}
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return signingKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Println("[INFO]", claims["email"], claims["sub"], claims["name"])
		response := Response{
			Message: "valid access_token",
		}
		log.Printf("[INFO] Response: %+v", response)
		return &serverApp.SuccessResponse{
			StatusCode: http.StatusOK,
		}, nil
	}
	log.Println("[ERROR]", err)

	response := Response{
		Message: "invalid access_token",
	}
	log.Printf("[INFO] Response: %+v", response)
	return nil, &serverApp.ErrorResponse{
		StatusCode: http.StatusBadRequest,
	}
}
