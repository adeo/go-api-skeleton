package authentication

import (
	"github.com/gin-gonic/gin"
)

type ServiceFake struct {
}

func NewServiceFake() Service {
	return &ServiceFake{}
}

func (s *ServiceFake) Introspect(c *gin.Context) (*IntrospectResponse, error) {
	return &IntrospectResponse{
		Active:  true,
		Subject: "10000000",
	}, nil
}
