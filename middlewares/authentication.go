package middlewares

import (
	"github.com/adeo/turbine-go-api-skeleton/services/authentication"
	"github.com/adeo/turbine-go-api-skeleton/storage/model"
	"github.com/adeo/turbine-go-api-skeleton/utils"
	"github.com/adeo/turbine-go-api-skeleton/utils/httputils"
	"github.com/gin-gonic/gin"
)

func GetAuthenticationMiddleware(authService authentication.Service) gin.HandlerFunc {
	unauthError := model.ErrInvalidCredentials
	unauthError.Headers = map[string][]string{
		httputils.HeaderNameWWWAuthenticate: {
			"Bearer",
		},
	}

	return func(c *gin.Context) {
		i, err := authService.Introspect(c)
		if err != nil {
			httputils.JSONErrorWithMessage(c.Writer, unauthError, err.Error())
			c.Abort()
			return
		}

		if !i.Active {
			httputils.JSONErrorWithMessage(c.Writer, unauthError, "token expired")
			c.Abort()
			return
		}

		c.Set(utils.ContextKeyAuthIntrospect, i)
		c.Next()
	}
}
