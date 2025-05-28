package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "ERROR"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func OKWithData(data map[string]interface{}) map[string]interface{} {
	resp := make(map[string]interface{})

	for k, v := range data {
		resp[k] = v
	}
	return resp
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}
