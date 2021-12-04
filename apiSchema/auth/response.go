package auth

type SignUpResponse struct {
	Token  string `json:"token"`
	Secret string `json:"secret"`
}
