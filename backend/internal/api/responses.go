package api

type Response struct {
	Code string `json:"code"`
	Data any    `json:"data"`
}

func NewReponse(code string, data any) *Response {
	if err, ok := data.(error); ok {
		data = err.Error()
	}
	return &Response{Code: code, Data: data}
}
