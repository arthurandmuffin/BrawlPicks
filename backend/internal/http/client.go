package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

type Client struct {
	client *http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) Do(req *http.Request, respData any) (err error) {
	httpResp, err := c.client.Do(req)
	if err != nil {
		err = wraps(EHTTP, err)
		return
	}
	defer httpResp.Body.Close()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		err = wraps(ERES, err)
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		codeStr := strconv.Itoa(httpResp.StatusCode)
		err = wraps(errors.New(codeStr), bytesToErr(body))
		return
	}
	if respData == nil {
		return
	}
	if err = json.Unmarshal(body, respData); err != nil {
		err = wraps(ERES, bytesToErr(body))
		return
	}
	return
}
