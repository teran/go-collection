package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func (s *downloadFileTestSuite) TestDownloadFile() {
	payload := []byte("test data blob")

	s.handlersMock.On("Blob").Return(http.StatusOK, payload).Once()

	tempFilename := path.Join(s.T().TempDir(), "TestDownloadFile")
	err := DownloadFile(s.ctx, s.srv.URL+"/blob", tempFilename)
	s.Require().NoError(err)

	fi, err := os.Stat(tempFilename)
	s.Require().NoError(err)
	s.Require().Equal(int64(len(payload)), fi.Size())

	data, err := os.ReadFile(tempFilename)
	s.Require().NoError(err)
	s.Require().Equal(payload, data)
}

func (s *downloadFileTestSuite) TestDefaultHeaders() {
	s.handlersMock.
		On(
			"Headers",
			http.Header{
				"Accept-Encoding": []string{"gzip"},
				"User-Agent":      []string{"download file client/1.0"},
			},
		).
		Return(http.StatusOK).
		Once()

	err := DownloadFile(s.ctx, s.srv.URL+"/headers", path.Join(s.T().TempDir(), "TestHeaders"))
	s.Require().NoError(err)
}

func (s *downloadFileTestSuite) TestCustomUserAgent() {
	s.handlersMock.
		On(
			"Headers",
			http.Header{
				"Accept-Encoding": []string{"gzip"},
				"User-Agent":      []string{"my custom agent/1.0"},
			},
		).
		Return(http.StatusOK).
		Once()

	err := DownloadFile(s.ctx, s.srv.URL+"/headers", path.Join(s.T().TempDir(), "TestHeaders"), UserAgent("my custom agent/1.0"))
	s.Require().NoError(err)
}

// ========================================================================
// Test suite setup
// ========================================================================
type downloadFileTestSuite struct {
	suite.Suite

	ctx         context.Context
	ctxCancelFn context.CancelFunc

	handlersMock *httpHandlerMock
	srv          *httptest.Server
}

func (s *downloadFileTestSuite) SetupTest() {
	s.ctx, s.ctxCancelFn = context.WithTimeout(context.Background(), 5*time.Second)

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	s.handlersMock = &httpHandlerMock{}

	e.GET("/blob", s.handlersMock.Blob)
	e.GET("/headers", s.handlersMock.Headers)

	s.srv = httptest.NewServer(e)
}

func (s *downloadFileTestSuite) TearDownTest() {
	s.ctxCancelFn()
	s.srv.Close()
}

func TestDownloadFileTestSuite(t *testing.T) {
	suite.Run(t, &downloadFileTestSuite{})
}

// ========================================================================
// Mocks
// ========================================================================
type httpHandlerMock struct {
	mock.Mock
}

func (m *httpHandlerMock) Blob(c echo.Context) error {
	args := m.Called()

	statusCode := args.Int(0)
	data := args.Get(1).([]byte)

	return c.Blob(statusCode, "application/octet-stream", data)
}

func (m *httpHandlerMock) Headers(c echo.Context) error {
	args := m.Called(c.Request().Header)
	return c.NoContent(args.Int(0))
}
