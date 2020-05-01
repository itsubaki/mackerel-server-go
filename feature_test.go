package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/mackerel-api/pkg/infrastructure"
	"github.com/itsubaki/mackerel-api/pkg/infrastructure/config"
	"github.com/itsubaki/mackerel-api/pkg/infrastructure/handler"
	"github.com/itsubaki/mackerel-api/pkg/interface/database"
)

type apiFeature struct {
	header http.Header
	body   io.Reader
	resp   *httptest.ResponseRecorder

	config  *config.Config
	handler database.SQLHandler
	server  *gin.Engine
	random  map[string]string
}

func (a *apiFeature) start() {
	gin.SetMode(gin.ReleaseMode)
	os.Setenv("DATABASE_NAME", "mackerel_feature_test")

	c := config.New()
	h := handler.New(c)
	r := infrastructure.Router(h)

	a.config = c
	a.handler = h
	a.server = r
	a.random = make(map[string]string)
}

func (a *apiFeature) stop() {
	if err := a.handler.Transact(func(tx database.Tx) error {
		q := fmt.Sprintf("drop database if exists %s", a.config.DatabaseName)
		if _, err := tx.Exec(q); err != nil {
			return fmt.Errorf("drop database: %v", err)
		}
		return nil
	}); err != nil {
		panic(err)
	}

	if err := a.handler.Close(); err != nil {
		panic(err)
	}
}

func (a *apiFeature) reset(m *messages.Pickle) {
	a.header = make(http.Header)
	a.body = nil
	a.resp = httptest.NewRecorder()
}

func (a *apiFeature) SetXAPIKEY(k string) error {
	a.header.Add("X-Api-key", k)
	return nil
}

func (a *apiFeature) SetContentType(t string) error {
	a.header.Add("Content-Type", t)
	return nil
}

func (a *apiFeature) SetRequestBody(b *messages.PickleStepArgument_PickleDocString) error {
	a.body = bytes.NewBuffer([]byte(b.Content))
	return nil
}

func (a *apiFeature) Request(method, endpoint string) error {
	for k, v := range a.random {
		endpoint = strings.Replace(endpoint, k, v, -1)
	}

	req := httptest.NewRequest(method, endpoint, a.body)
	req.Header = a.header
	a.server.ServeHTTP(a.resp, req)

	return nil
}

func (a *apiFeature) ResponseCodeShouldBe(code int) error {
	if code == a.resp.Code {
		return nil
	}

	return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, a.resp.Code)
}

func (a *apiFeature) ResponseShouldMatchJson(body *messages.PickleStepArgument_PickleDocString) error {
	var expected, actual map[string]interface{}

	if err := json.Unmarshal([]byte(body.Content), &expected); err != nil {
		return err
	}

	if err := json.Unmarshal(a.resp.Body.Bytes(), &actual); err != nil {
		return err
	}

	for k, v := range expected {
		vv, ok := actual[k]
		if !ok {
			return fmt.Errorf("expected JSON does not match actual, %#v vs. %#v", expected, actual)
		}

		switch variable := v.(type) {
		case string:
			if variable == "<host_id>" {
				a.random["<host_id>"] = vv.(string)
				expected[k] = vv
			}
		}
	}

	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected JSON does not match actual, %#v vs. %#v", expected, actual)
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	a := &apiFeature{}
	s.BeforeSuite(a.start)

	s.BeforeScenario(a.reset)
	s.Step(`^I set X-Api-Key header with "([^"]*)"$`, a.SetXAPIKEY)
	s.Step(`^I set Content-Type header with "([^"]*)"$`, a.SetContentType)
	s.Step(`^I set request body with:$`, a.SetRequestBody)
	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, a.Request)
	s.Step(`^the response code should be (\d+)$`, a.ResponseCodeShouldBe)
	s.Step(`^the response should match json:$`, a.ResponseShouldMatchJson)

	s.AfterSuite(a.stop)
}
