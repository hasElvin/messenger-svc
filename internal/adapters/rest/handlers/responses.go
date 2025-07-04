package handlers

type SuccessResponse struct {
	Message string `json:"message"`
}

type FailResponse struct {
	Message string `json:"error"`
}
