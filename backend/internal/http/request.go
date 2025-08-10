package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	AuthoriazationHeader = "Authorization"
)

func NewBearerAuthRequest(method, url, token string, body any) (req *http.Request, err error) {
	if req, err = NewRequest(method, url, body); err != nil {
		return
	}
	req.Header.Add(AuthoriazationHeader, toBearerToken(token))
	return
}

func NewRequest(method, url string, body any) (req *http.Request, err error) {
	b, err := newRequestBody(body)
	if err != nil {
		err = wraps(EJSON, err)
		return
	}
	if req, err = http.NewRequest(method, url, b); err != nil {
		err = wraps(EHTTP, err)
		return
	}
	return
}

func toBearerToken(token string) string {
	return fmt.Sprintf("Bearer %s", token)
}

func newRequestBody(body any) (rv io.Reader, err error) {
	if body == nil {
		return
	}
	rv, ok := body.(io.Reader)
	if ok {
		return
	}
	buffer := new(bytes.Buffer)
	err = json.NewEncoder(buffer).Encode(body)
	return
}
