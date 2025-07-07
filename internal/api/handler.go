package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"order-service-wb/internal/service"
)

type Handler struct {
	serv service.OrderService
}

func NewHandler(serv service.OrderService) *Handler {
	return &Handler{
		serv: serv,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/order/:uid", h.GetOrderByID)
	r.Static("/web", "./web/static")

	return r
}

func (h *Handler) GetOrderByID(c *gin.Context) {
	orderID := c.Param("uid")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order ID is required"})
		return
	}

	order, err := h.serv.GetOrderByID(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get order"})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}
