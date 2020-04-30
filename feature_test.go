package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/mackerel-api/pkg/infrastructure"
	"github.com/itsubaki/mackerel-api/pkg/infrastructure/config"
	"github.com/itsubaki/mackerel-api/pkg/infrastructure/handler"
	"net/http/httptest"
	"reflect"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
)

type APIFeature struct {
	xapikey string
	resp    *httptest.ResponseRecorder
}

func (a *APIFeature) reset(m *messages.Pickle) {
	a.resp = httptest.NewRecorder()
}

func (a *APIFeature) SetXAPIKEY(xapikey string) error {
	a.xapikey = xapikey
	return nil
}

func (a *APIFeature) Request(method, endpoint string) error {
	gin.SetMode(gin.ReleaseMode)

	c := config.New()
	h := handler.New(c)
	r := infrastructure.Router(h)

	req := httptest.NewRequest(method, endpoint, nil)
	req.Header.Add("X-Api-key", a.xapikey)
	r.ServeHTTP(a.resp, req)

	return nil
}

func (a *APIFeature) ResponseCodeShouldBe(code int) error {
	if code == a.resp.Code {
		return nil
	}

	return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, a.resp.Code)
}

func (a *APIFeature) ResponseShouldMatchJson(body *messages.PickleStepArgument_PickleDocString) error {
	var expected, actual interface{}

	if err := json.Unmarshal([]byte(body.Content), &expected); err != nil {
		return err
	}

	if err := json.Unmarshal(a.resp.Body.Bytes(), &actual); err != nil {
		return err
	}

	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expected, actual)
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	a := &APIFeature{}

	s.BeforeScenario(a.reset)
	s.Step(`^I fill the XAPIKEY with "([^"]*)"$`, a.SetXAPIKEY)
	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, a.Request)
	s.Step(`^the response code should be (\d+)$`, a.ResponseCodeShouldBe)
	s.Step(`^the response should match json:$`, a.ResponseShouldMatchJson)

}
