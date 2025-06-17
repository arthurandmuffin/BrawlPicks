package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}

func NewContext(ctx *gin.Context) *Context {
	return &Context{ctx}
}

func (c *Context) Response(httpCode int, code string, data any) {
	c.JSON(httpCode, NewReponse(code, data))
}

func (c *Context) OK(data any) {
	c.Response(http.StatusOK, OK, data)
}

func (c *Context) BadRequest(data any) {
	c.Response(http.StatusBadRequest, EREQ, data)
}

func (c *Context) InternalServerError(data any) {
	c.Response(http.StatusInternalServerError, EINT, data)
}
