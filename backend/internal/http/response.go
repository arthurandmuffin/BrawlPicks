package http

import (
	"encoding/json"
)

type Response struct {
	Code Code            `json:"code"`
	Data json.RawMessage `json:"data"`
}
