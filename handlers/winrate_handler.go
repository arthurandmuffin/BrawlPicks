package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"BrawlPicks/services"
)

type WinrateHandler struct {
	service services.WinrateService
}

func NewWinrateHandler(service services.WinrateService) *WinrateHandler {
	return &WinrateHandler{service: service}
}

func (h *WinrateHandler) GetTopWinrates(c *gin.Context) {
	winrates, err := h.service.GetTopWinrates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, winrates)
}
