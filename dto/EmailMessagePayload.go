package dto

type destination struct {
	ToAddresses  []string `json:"toAddresses"`
	CcAddresses  []string `json:"ccAddresses"`
	BccAddresses []string `json:"bccAddresses"`
}

type EmailVerificationMsgPayload struct {
	Context          string                 `json:"context"`
	Destination      destination            `json:"destination"`
	EmailBodyContent map[string]interface{} `json:"emailBodyContent"`
	ReplyToAddresses []string               `json:"replyToAddresses"`
	ReturnPath       string                 `json:"returnPath"`
}

func NewEmailVerificationMsgPayload(
	context string, to []string, cc []string, bcc []string,
	emailBodyContent map[string]interface{}, replyTo []string, returnPath string) *EmailVerificationMsgPayload {
	return &EmailVerificationMsgPayload{
		Context: context,
		Destination: destination{
			ToAddresses:  to,
			CcAddresses:  cc,
			BccAddresses: bcc,
		},
		EmailBodyContent: emailBodyContent,
		ReplyToAddresses: replyTo,
		ReturnPath:       returnPath,
	}
}
