package authentication

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/adeo/turbine-models/utils/httputils"
	"github.com/gin-gonic/gin"
)

type ServiceHTTP struct {
	httpClient *http.Client
	url        string
}

func NewServiceHTTP(authenticationServiceURL string) Service {
	return &ServiceHTTP{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		url: authenticationServiceURL,
	}
}

func (s *ServiceHTTP) TokenIntrospect(c *gin.Context) (*TokenIntrospectionResponse, error) {
	auth := c.GetHeader(httputils.HeaderNameAuthorization)
	authSplitted := strings.SplitN(auth, " ", 2)
	if len(authSplitted) != 2 {
		return nil, newAuthenticationError(ErrTypeAuthorizationMalformed, nil)
	}

	if strings.ToUpper(authSplitted[0]) != strings.ToUpper("Bearer") && strings.ToUpper(authSplitted[0]) != strings.ToUpper("JWT") {
		return nil, newAuthenticationError(ErrTypeAuthorizationMalformed, errors.New("unsupported authentication scheme"))
	}

	token := authSplitted[1]

	// create request
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/auth/introspect?domainroles=true", s.url), strings.NewReader(fmt.Sprintf("token=%s", token)))
	if err != nil {
		return nil, err
	}

	req.Header.Set(httputils.HeaderNameContentType, httputils.HeaderValueApplicationXWWWFormURLEncoded)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, newAuthenticationError(ErrTypeServerError, fmt.Errorf("server answered with status %d and body `%s`", resp.StatusCode, string(bodyBytes)))
	}

	// build json result
	result := &TokenIntrospectionResponse{}
	if json.Unmarshal(bodyBytes, result) != nil {
		return nil, errors.New(fmt.Sprintf("unable to parse the introspect result: %s", string(bodyBytes)))
	}

	return result, nil
}
