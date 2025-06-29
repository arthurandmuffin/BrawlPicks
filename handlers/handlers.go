package handlers

import (
	"BrawlPicks/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BrawlStarsHandler struct {
	service *services.BrawlStarsService
}

func NewBrawlStarsHandler(service *services.BrawlStarsService) *BrawlStarsHandler {
	return &BrawlStarsHandler{service: service}
}

func (h *BrawlStarsHandler) GetPlayer(c *gin.Context) {
	tag := c.Param("tag")
	playerData, err := h.service.GetPlayer(tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, playerData)
}

func (h *BrawlStarsHandler) GetPower11Brawlers(c *gin.Context) {
	tag := c.Param("tag")

	brawlers, err := h.service.GetPower11Brawlers(tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"power11_brawlers": brawlers,
	})
}

/*
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
*/
