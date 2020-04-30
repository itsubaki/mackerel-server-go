package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/mackerel-api/pkg/infrastructure"
	"net/http/httptest"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
)

type APIFeature struct {
	resp *httptest.ResponseRecorder
}

func (a *APIFeature) reset(m *messages.Pickle) {
	a.resp = httptest.NewRecorder()
}

func (a *APIFeature) Request(method, endpoint string) error {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	infrastructure.Status(r)

	req := httptest.NewRequest(method, endpoint, nil)
	r.ServeHTTP(a.resp, req)

	return nil
}

func (a *APIFeature) ResponseCodeShouldBe(code int) error {
	if code == a.resp.Code {
		return nil
	}

	return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, a.resp.Code)
}

func FeatureContext(s *godog.Suite) {
	a := &APIFeature{}

	s.BeforeScenario(a.reset)
	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, a.Request)
	s.Step(`^the response code should be (\d+)$`, a.ResponseCodeShouldBe)
}
