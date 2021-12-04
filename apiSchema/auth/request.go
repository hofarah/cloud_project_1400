package auth

type SignUpRequest struct {
	UserName string `json:"username" validate:"required"`
}
