package handlers

import (
	"net/http"

	"github.com/adeo/go-api-skeleton/utils"
	"github.com/gin-gonic/gin"
)

func (hc *Context) GetHealth(c *gin.Context) {
	utils.JSON(c.Writer, http.StatusNoContent, nil)
}
