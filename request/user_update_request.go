package request

type UserUpdatePayload struct {
	Email     string `json:"email"` // Email is required and used for lookup
	Username  string `json:"username"`
	Grade     string `json:"grade"`
	Major     string `json:"major"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	Direction string `json:"direction"`
}
