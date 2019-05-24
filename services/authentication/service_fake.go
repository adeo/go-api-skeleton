package authentication

import (
	"errors"
	"strings"

	"github.com/adeo/turbine-go-api-skeleton/utils/httputils"
	"github.com/gin-gonic/gin"
)

type ServiceFake struct {
}

func NewServiceFake() Service {
	return &ServiceFake{}
}

func (s *ServiceFake) TokenIntrospect(c *gin.Context) (*TokenIntrospectionResponse, error) {
	auth := c.GetHeader(httputils.HeaderNameAuthorization)
	authSplitted := strings.SplitN(auth, " ", 2)
	if len(authSplitted) != 2 {
		return nil, newAuthenticationError(ErrTypeAuthorizationMalformed, nil)
	}

	if strings.ToUpper(authSplitted[0]) != strings.ToUpper("Bearer") {
		return nil, newAuthenticationError(ErrTypeAuthorizationMalformed, errors.New("unsupported authentication scheme"))
	}

	return &TokenIntrospectionResponse{
		Active:  true,
		Subject: "10000000",
	}, nil
}
