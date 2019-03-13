package handlers

import (
	"net/http"

	"github.com/adeo/go-api-skeleton/utils"
	"github.com/gin-gonic/gin"
)

func (hc *Context) GetOpenAPISchema(c *gin.Context) {
	utils.YAML(c.Writer, http.StatusOK, utils.OpenAPISchema)
}
