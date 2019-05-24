package authentication

import (
	"github.com/gin-gonic/gin"
)

type Service interface {
	TokenIntrospect(c *gin.Context) (*TokenIntrospectionResponse, error)
}
