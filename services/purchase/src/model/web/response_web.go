package response

type Web struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type GlobalErrorHandlerResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
