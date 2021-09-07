package main_test

import (
	"bytes"
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
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"github.com/itsubaki/mackerel-server-go/pkg/infrastructure"
	"github.com/itsubaki/mackerel-server-go/pkg/infrastructure/config"
	"github.com/itsubaki/mackerel-server-go/pkg/infrastructure/handler"
	"github.com/itsubaki/mackerel-server-go/pkg/interface/database"
	"github.com/jfilipczyk/gomatch"
)

var api = &apiFeature{}

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
	h, err := handler.New(c.Driver, c.Host, c.Database, handler.Opt{
		SQLMode:         c.SQLMode,
		Timeout:         &c.Timeout,
		Sleep:           &c.Sleep,
		MaxIdleConns:    &c.MaxIdleConns,
		MaxOpenConns:    &c.MaxOpenConns,
		ConnMaxLifetime: &c.ConnMaxLifetime,
	})
	if err != nil {
		panic(err)
	}
	log.Printf("db connected")

	if err := infrastructure.RunFixture(h); err != nil {
		panic(err)
	}

	a.config = c
	a.handler = h
	a.server = infrastructure.APIv0(infrastructure.Default(), h)
	a.keep = make(map[string]interface{})
}

func (a *apiFeature) stop() {
	if err := a.handler.Close(); err != nil {
		panic(err)
	}
}

func (a *apiFeature) reset(sc *godog.Scenario) {
	a.header = make(http.Header)
	a.body = nil
	a.resp = httptest.NewRecorder()
}

func (a *apiFeature) replace(str string) string {
	for k, v := range a.keep {
		switch val := v.(type) {
		case string:
			str = strings.Replace(str, k, val, -1)
		}
	}

	return str
}

func (a *apiFeature) SetHeader(k, v string) error {
	a.header.Add(k, v)
	return nil
}

func (a *apiFeature) SetRequestBody(body *godog.DocString) error {
	r := a.replace(body.Content)
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

	return fmt.Errorf("got=%v, want=%v", a.resp.Code, code)
}

func (a *apiFeature) ResponseShouldMatchJSON(body *godog.DocString) error {
	want := a.replace(body.Content)
	got := a.resp.Body.String()

	ok, err := gomatch.NewDefaultJSONMatcher().Match(want, got)
	if err != nil {
		return fmt.Errorf("got=%v, want=%v, match: %v", got, want, err)
	}

	if !ok {
		return fmt.Errorf("got=%v, want=%v", got, want)
	}

	return nil
}

func (a *apiFeature) Keep(key, as string) error {
	var resposne map[string]interface{}
	if err := json.Unmarshal(a.resp.Body.Bytes(), &resposne); err != nil {
		return fmt.Errorf("body=%s, unmarshal: %v", a.resp.Body.String(), err)
	}

	if v, ok := resposne[key]; ok {
		a.keep[as] = v
	}

	return nil
}

func (a *apiFeature) AlertsExists(alerts *godog.Table) error {
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

func (a *apiFeature) UsersExists(users *godog.Table) error {
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

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	before := func() {
		gin.SetMode(gin.ReleaseMode)
		os.Setenv("DATABASE", "mackerel_test")

		c := config.New()
		if err := handler.Exec(c.Driver, c.Host, []string{
			fmt.Sprintf("drop database if exists %s", c.Database),
			fmt.Sprintf("create database if not exists %s", c.Database),
		}); err != nil {
			panic(err)
		}

		api.start()
	}
	after := func() {
		api.stop()
	}

	ctx.BeforeSuite(before)
	ctx.AfterSuite(after)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(api.reset)

	ctx.Step(`^the following alerts exist:$`, api.AlertsExists)
	ctx.Step(`^the following users exist:$`, api.UsersExists)
	ctx.Step(`^I set "([^"]*)" header with "([^"]*)"$`, api.SetHeader)
	ctx.Step(`^I set request body:$`, api.SetRequestBody)
	ctx.Step(`^I send "([^"]*)" request to "([^"]*)"$`, api.Request)
	ctx.Step(`^the response code should be (\d+)$`, api.ResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, api.ResponseShouldMatchJSON)
	ctx.Step(`^I keep the JSON response at "([^"]*)" as "([^"]*)"$`, api.Keep)
}
