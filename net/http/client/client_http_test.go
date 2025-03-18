package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func (s *testSuite) TestResponse() {
	s.handlers.On("JSONResponder").Return(http.StatusOK, echo.Map{
		"status": "ok",
	}).Once()

	resp, err := New().
		Base(s.srv.URL).
		Get("/json").
		Do(s.ctx, nil)
	s.Require().NoError(err)

	data, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Require().JSONEq(`{"status":"ok"}`, string(data))
}

func (s *testSuite) TestJSON() {
	s.handlers.On("JSONResponder").Return(http.StatusOK, echo.Map{
		"status": "ok",
	}).Once()

	resp := map[string]string{}
	errResp := map[string]string{}
	statusCode, err := New().
		Base(s.srv.URL).
		Get("/json").
		DoJSON(s.ctx, nil, &resp, &errResp)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, statusCode)
	s.Require().Equal(map[string]string{
		"status": "ok",
	}, resp)
}

func (s *testSuite) TestJSONError() {
	s.handlers.On("JSONResponder").Return(http.StatusNotImplemented, echo.Map{
		"status": "not implemented",
	}).Once()

	resp := map[string]string{}
	errResp := map[string]string{}
	statusCode, err := New().
		Base(s.srv.URL).
		Get("/json").
		DoJSON(s.ctx, nil, &resp, &errResp)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotImplemented, statusCode)
	s.Require().Equal(map[string]string{
		"status": "not implemented",
	}, errResp)
}

func (s *testSuite) TestJSONWrongMimeType() {
	s.handlers.On("CustomResponder").Return(http.StatusOK, "text/plain", []byte(`{"status":"ok"}`)).Once()

	resp := map[string]string{}
	errResp := map[string]string{}
	statusCode, err := New().
		Base(s.srv.URL).
		Get("/custom").
		DoJSON(s.ctx, nil, &resp, &errResp)
	s.Require().Error(err)
	s.Require().Equal(http.StatusOK, statusCode)
	s.Require().Equal(ErrUnsupportedMediaType, errors.Cause(err))
}

func (s *testSuite) TestHeaders() {
	s.handlers.On("HeadersResponder", http.Header{
		"Accept-Encoding": []string{"gzip"},
		"User-Agent":      []string{"test-agent/1.0"},
		"X-Test-Header":   []string{"test value"},
	}).Return(nil).Once()

	resp, err := New().
		Base(s.srv.URL).
		UserAgent("test-agent/1.0").
		Header("X-Test-Header", "test value").
		Get("/headers").
		Do(s.ctx, nil)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

// ========================================================================
// Test suite setup
// ========================================================================
type testSuite struct {
	suite.Suite

	ctx      context.Context
	handlers *handlers
	srv      *httptest.Server
}

func (s *testSuite) SetupTest() {
	s.ctx = context.Background()

	s.handlers = &handlers{}

	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/json", s.handlers.JSONResponder)
	e.GET("/custom", s.handlers.CustomResponder)
	e.GET("/headers", s.handlers.HeadersResponder)

	s.srv = httptest.NewServer(e)
}

func (s *testSuite) TearDownTest() {
	s.srv.Close()

	s.handlers.AssertExpectations(s.T())
}

func TestSuite(t *testing.T) {
	suite.Run(t, &testSuite{})
}

// ========================================================================
// HTTP handlers mock
// ========================================================================

type handlers struct {
	mock.Mock
}

func (m *handlers) JSONResponder(c echo.Context) error {
	args := m.Called()
	return c.JSON(args.Int(0), args.Get(1).(echo.Map))
}

func (m *handlers) CustomResponder(c echo.Context) error {
	args := m.Called()
	return c.Blob(args.Int(0), args.String(1), args.Get(2).([]byte))
}

func (m *handlers) HeadersResponder(c echo.Context) error {
	args := m.Called(c.Request().Header)
	return args.Error(0)
}
