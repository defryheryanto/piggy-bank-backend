package response

type SuccessResponse struct {
	Success bool `json:"success"`
}

func NewSuccessResponse() *SuccessResponse {
	return &SuccessResponse{true}
}
