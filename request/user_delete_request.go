package request

type UserDeletePayload struct {
	Email string `json:"email"` // Email is required for deletion
}
