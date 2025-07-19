package request

type RegisterRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Grade     string `json:"grade"`
	Major     string `json:"major"`
	Direction string `json:"direction"`
}
