package handlers

import (
	"net/http"

	"github.com/adeo/turbine-go-api-skeleton/api"
	"github.com/adeo/turbine-go-api-skeleton/utils/httputils"
	"github.com/gin-gonic/gin"
)

func (hc *Context) GetOpenAPISchema(c *gin.Context) {
	httputils.YAML(c.Writer, http.StatusOK, api.OpenAPISchema)
}
