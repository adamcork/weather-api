package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIHandler struct {
}

func NewAPIHandler() *APIHandler {
	return &APIHandler{}
}

func (h *APIHandler) Weather(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello")
}
