package api

import "github.com/gin-gonic/gin"

func UriHandler[T any](next func(*Context, *T)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			c = NewContext(ctx)
			r = new(T)
		)
		if err := c.ShouldBindUri(r); err != nil {
			c.BadRequest(err)
			return
		}
		next(c, r)
	}
}
