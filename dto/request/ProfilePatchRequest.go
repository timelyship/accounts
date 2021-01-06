package request

type ProfilePatchRequest struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}
