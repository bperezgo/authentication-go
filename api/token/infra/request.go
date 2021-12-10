package infra

type Request struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
}
