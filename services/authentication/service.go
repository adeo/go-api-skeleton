package authentication

import (
	"github.com/gin-gonic/gin"
)

type Service interface {
	Introspect(c *gin.Context) (*IntrospectResponse, error)
}
