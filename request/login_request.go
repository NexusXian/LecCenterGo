package request

type LoginRequest struct {
	Account  string `json:"account"` // 邮箱或手机号
	Password string `json:"password"`
}
