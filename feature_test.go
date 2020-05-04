package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/mackerel-api/pkg/infrastructure"
	"github.com/itsubaki/mackerel-api/pkg/infrastructure/config"
	"github.com/itsubaki/mackerel-api/pkg/infrastructure/handler"
	"github.com/itsubaki/mackerel-api/pkg/interface/database"
	"github.com/jfilipczyk/gomatch"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	os.Setenv("DATABASE_NAME", "mackerel_test")

	c := config.New()
	db, err := sql.Open(c.Driver, c.DataSourceName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	start := time.Now()
	for {
		if time.Since(start) > 10*time.Minute {
			panic("db ping time over")
		}

		if err := db.Ping(); err != nil {
			log.Printf("db ping: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		break
	}

	q0 := fmt.Sprintf("drop database if exists %s", c.DatabaseName)
	if _, err := db.Exec(q0); err != nil {
		panic(err)
	}

	q1 := fmt.Sprintf("create database if not exists %s", c.DatabaseName)
	if _, err := db.Exec(q1); err != nil {
		panic(err)
	}
}

type apiFeature struct {
	header http.Header
	body   io.Reader
	resp   *httptest.ResponseRecorder

	config  *config.Config
	handler database.SQLHandler
	server  *gin.Engine
	keep    map[string]string
}

func (a *apiFeature) start() {
	c := config.New()
	h := handler.New(c)
	r := infrastructure.Router(h)

	a.config = c
	a.handler = h
	a.server = r
	a.keep = make(map[string]string)
}

func (a *apiFeature) stop() {
	if err := a.handler.Close(); err != nil {
		panic(err)
	}
}

func (a *apiFeature) reset(m *messages.Pickle) {
	a.header = make(http.Header)
	a.body = nil
	a.resp = httptest.NewRecorder()
}

func (a *apiFeature) replace(str string) string {
	for k, v := range a.keep {
		str = strings.Replace(str, k, v, -1)
	}

	return str
}

func (a *apiFeature) SetHeader(k, v string) error {
	a.header.Add(k, v)
	return nil
}

func (a *apiFeature) SetRequestBody(b *messages.PickleStepArgument_PickleDocString) error {
	a.body = bytes.NewBuffer([]byte(b.Content))
	return nil
}

func (a *apiFeature) Request(method, endpoint string) error {
	req := httptest.NewRequest(method, a.replace(endpoint), a.body)
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
	expected := a.replace(body.Content)
	actual := a.resp.Body.String()

	ok, err := gomatch.NewDefaultJSONMatcher().Match(expected, actual)
	if err != nil {
		return fmt.Errorf("match: %v", err)
	}

	if !ok {
		return fmt.Errorf("expected JSON does not match actual, %#v vs. %#v", expected, actual)
	}

	return nil
}

func (a *apiFeature) Keep(key, as string) error {
	var actual map[string]string
	if err := json.Unmarshal(a.resp.Body.Bytes(), &actual); err != nil {
		return err
	}

	for k, v := range actual {
		if k != key {
			continue
		}

		a.keep[as] = v
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	a := &apiFeature{}
	s.BeforeSuite(a.start)

	s.BeforeScenario(a.reset)
	s.Step(`^I set "([^"]*)" header with "([^"]*)"$`, a.SetHeader)
	s.Step(`^I set request body:$`, a.SetRequestBody)
	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, a.Request)
	s.Step(`^the response code should be (\d+)$`, a.ResponseCodeShouldBe)
	s.Step(`^the response should match json:$`, a.ResponseShouldMatchJson)
	s.Step(`^I keep the JSON response at "([^"]*)" as "([^"]*)"$`, a.Keep)

	s.AfterSuite(a.stop)
}
