package handlers

import (
	"BrawlPicks/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandlerTest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": h.service.ServiceTest()})
}
