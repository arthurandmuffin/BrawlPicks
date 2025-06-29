package handlers

import (
	"BrawlPicks/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PLHandler struct {
	service *services.PLService
}

func NewPLHandler(service *services.PLService) *PLHandler {
	return &PLHandler{service: service}
}

func (h *PLHandler) GetTopWinrates(c *gin.Context) {
	data, err := h.service.GetTop5BrawlersByMap()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
