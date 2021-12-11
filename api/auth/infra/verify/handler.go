package verify

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bperezgo/authentication/config"
	"github.com/golang-jwt/jwt"
)

type TokenClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.StandardClaims
}

// Handler function to verify the token
func Handler(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] Handling the token")
	c := config.GetConfig()
	signingKey := []byte(c.AuthJwtSecret)
	urlValues := req.URL.Query()
	accessToken := urlValues.Get("access_token")
	if accessToken == "" {
		res.WriteHeader(http.StatusBadRequest)
		response := Response{
			Message: "invalid access_token",
		}
		resByte, _ := json.Marshal(response)
		log.Printf("[INFO] Response: %+v", response)
		res.Write([]byte(resByte))
		return
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
		res.WriteHeader(http.StatusOK)
		response := Response{
			Message: "valid access_token",
		}
		resByte, _ := json.Marshal(response)
		log.Printf("[INFO] Response: %+v", response)
		res.Write([]byte(resByte))
		return
	}
	log.Println("[ERROR]", err)
	res.WriteHeader(http.StatusBadRequest)
	response := Response{
		Message: "invalid access_token",
	}
	resByte, _ := json.Marshal(response)
	log.Printf("[INFO] Response: %+v", response)
	res.Write([]byte(resByte))
}
