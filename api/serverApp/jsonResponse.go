package serverApp

import (
	"encoding/json"
	"net/http"
)

// Use the http.ResponseWriter to send the information to the user
func jsonResponse(res http.ResponseWriter, response interface{}) {
	resByte, _ := json.Marshal(response)
	res.Write([]byte(resByte))
}
