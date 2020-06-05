package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/cucumber/godog"
	messages "github.com/cucumber/messages-go/v10"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"github.com/itsubaki/mackerel-server-go/pkg/infrastructure"
	"github.com/itsubaki/mackerel-server-go/pkg/infrastructure/config"
	"github.com/itsubaki/mackerel-server-go/pkg/infrastructure/handler"
	"github.com/itsubaki/mackerel-server-go/pkg/interface/database"
	"github.com/jfilipczyk/gomatch"
)

type apiFeature struct {
	header http.Header
	body   io.Reader
	resp   *httptest.ResponseRecorder

	config  *config.Config
	handler database.SQLHandler
	server  *gin.Engine
	keep    map[string]interface{}
}

func (a *apiFeature) start() {
	c := config.New()
	h, err := handler.New(c)
	if err != nil {
		panic(err)
	}

	if err := infrastructure.RunFixture(h); err != nil {
		panic(err)
	}

	a.config = c
	a.handler = h
	a.server = infrastructure.Router(h)
	a.keep = make(map[string]interface{})
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
		switch val := v.(type) {
		case string:
			str = strings.Replace(str, k, val, -1)
		default:
			continue
		}
	}

	return str
}

func (a *apiFeature) SetHeader(k, v string) error {
	a.header.Add(k, v)
	return nil
}

func (a *apiFeature) SetRequestBody(b *messages.PickleStepArgument_PickleDocString) error {
	r := a.replace(b.Content)
	a.body = bytes.NewBuffer([]byte(r))
	return nil
}

func (a *apiFeature) Request(method, endpoint string) error {
	r := a.replace(endpoint)
	req := httptest.NewRequest(method, r, a.body)
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

func (a *apiFeature) ResponseShouldMatchJSON(body *messages.PickleStepArgument_PickleDocString) error {
	expected := a.replace(body.Content)
	actual := a.resp.Body.String()

	ok, err := gomatch.NewDefaultJSONMatcher().Match(expected, actual)
	if err != nil {
		return fmt.Errorf("actual=%s, match: %v", actual, err)
	}

	if !ok {
		return fmt.Errorf("expected JSON does not match actual, %s vs. %s", expected, actual)
	}

	return nil
}

func (a *apiFeature) Keep(key, as string) error {
	var actual map[string]interface{}
	if err := json.Unmarshal(a.resp.Body.Bytes(), &actual); err != nil {
		return fmt.Errorf("body=%s, unmarshal: %v", a.resp.Body.String(), err)
	}

	for k, v := range actual {
		if k != key {
			continue
		}

		a.keep[as] = v
	}

	return nil
}

func (a *apiFeature) AlertsExists(alerts *messages.PickleStepArgument_PickleTable) error {
	list := make([]domain.Alert, 0)
	for i := 1; i < len(alerts.Rows); i++ {
		list = append(list, domain.Alert{
			OrgID:     alerts.Rows[i].Cells[0].Value,
			ID:        alerts.Rows[i].Cells[1].Value,
			Status:    alerts.Rows[i].Cells[2].Value,
			MonitorID: alerts.Rows[i].Cells[3].Value,
			Type:      alerts.Rows[i].Cells[4].Value,
		})
	}

	r := database.NewAlertRepository(a.handler)
	for i := range list {
		if _, err := r.Save(list[i].OrgID, &list[i]); err != nil {
			return fmt.Errorf("save alert: %v", err)
		}
	}

	return nil
}

func (a *apiFeature) UsersExists(users *messages.PickleStepArgument_PickleTable) error {
	list := make([]domain.User, 0)
	for i := 1; i < len(users.Rows); i++ {
		list = append(list, domain.User{
			OrgID:                   users.Rows[i].Cells[0].Value,
			ID:                      users.Rows[i].Cells[1].Value,
			ScreenName:              users.Rows[i].Cells[2].Value,
			Email:                   users.Rows[i].Cells[3].Value,
			Authority:               "owner",
			IsInRegistrationProcess: true,
			IsMFAEnabled:            true,
			AuthenticationMethods:   []string{"google"},
			JoinedAt:                time.Now().Unix(),
		})
	}

	r := database.NewUserRepository(a.handler)
	for i := range list {
		if err := r.Save(list[i].OrgID, &list[i]); err != nil {
			return fmt.Errorf("save user: %v", err)
		}
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	gin.SetMode(gin.ReleaseMode)
	os.Setenv("DATABASE", "mackerel_test")

	c := config.New()
	if err := handler.Query(c.Driver, c.Host, []string{
		fmt.Sprintf("drop database if exists %s", c.Database),
		fmt.Sprintf("create database if not exists %s", c.Database),
	}); err != nil {
		panic(err)
	}

	a := &apiFeature{}
	s.BeforeSuite(a.start)

	s.BeforeScenario(a.reset)
	s.Step(`^the following alerts exist:$`, a.AlertsExists)
	s.Step(`^the following users exist:$`, a.UsersExists)
	s.Step(`^I set "([^"]*)" header with "([^"]*)"$`, a.SetHeader)
	s.Step(`^I set request body:$`, a.SetRequestBody)
	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, a.Request)
	s.Step(`^the response code should be (\d+)$`, a.ResponseCodeShouldBe)
	s.Step(`^the response should match json:$`, a.ResponseShouldMatchJSON)
	s.Step(`^I keep the JSON response at "([^"]*)" as "([^"]*)"$`, a.Keep)

	s.AfterSuite(a.stop)
}
