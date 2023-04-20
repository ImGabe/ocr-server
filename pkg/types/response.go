package types

import "net/http"

type Response struct {
	Result  string `json:"result"`
	Version string `json:"version"`
}

func (rp *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewResponse(result string) *Response {
	return &Response{Result: result, Version: "0.0.1"}
}
