package request

type EmailVerificationRequest struct {
	Secret string `json:"secret"`
}
