package client

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func TestClientInit(t *testing.T) {
	r := require.New(t)

	c := New()
	r.Equal(Client{
		headers:     make(http.Header),
		queryParams: make(url.Values),
		cli:         http.DefaultClient,
	}, c)
}

func TestRequestMethods(t *testing.T) {
	r := require.New(t)

	c := New()

	r.Equal(
		Client{
			cli:         http.DefaultClient,
			uri:         "https://example.com",
			headers:     make(http.Header),
			queryParams: make(url.Values),
		},
		c.Basef("https://%s", "example.com"),
	)

	r.Equal(
		Client{
			cli:         http.DefaultClient,
			method:      "BLAH",
			headers:     make(http.Header),
			queryParams: make(url.Values),
			path:        "/someuri/blah",
		},
		c.Request("BLAH", "/someuri/%s", "blah"),
	)

	r.Equal(
		Client{
			cli:         http.DefaultClient,
			method:      "DELETE",
			headers:     make(http.Header),
			queryParams: make(url.Values),
			path:        "/someuri/blah",
		},
		c.Delete("/someuri/%s", "blah"),
	)

	r.Equal(
		Client{
			cli:         http.DefaultClient,
			method:      "GET",
			headers:     make(http.Header),
			queryParams: make(url.Values),
			path:        "/someuri/blah",
		},
		c.Get("/someuri/%s", "blah"),
	)
}

func TestHeaderMethods(t *testing.T) {
	r := require.New(t)
	c := New()

	r.Equal(Client{
		headers: map[string][]string{
			"Blahname": {"blahvalue"},
		},
		queryParams: make(url.Values),
		cli:         http.DefaultClient,
	}, c.Header("blahname", "blahvalue"))

	r.Equal(Client{
		headers: map[string][]string{
			"Authorization": {"test creds"},
		},
		queryParams: make(url.Values),
		cli:         http.DefaultClient,
	}, c.Auth("test", "creds"))

	r.Equal(Client{
		headers: map[string][]string{
			"Authorization": {"Basic dGVzdCB1c2VyOnRlc3QgcGFzc3dvcmQ="},
		},
		queryParams: make(url.Values),
		cli:         http.DefaultClient,
	}, c.BasicAuth("test user", "test password"))

	r.Equal(Client{
		headers: map[string][]string{
			"User-Agent": {"Some user agent/1.0"},
		},
		queryParams: make(url.Values),
		cli:         http.DefaultClient,
	}, c.UserAgent("Some user agent/1.0"))
}

func TestQueryParams(t *testing.T) {
	r := require.New(t)

	c := New()
	r.Equal(Client{
		headers:     make(http.Header),
		queryParams: make(url.Values),
		cli:         http.DefaultClient,
	}, c)

	c = c.QueryParam("some_key", "some_value")
	r.Equal(Client{
		headers: make(http.Header),
		queryParams: map[string][]string{
			"some_key": {"some_value"},
		},
		cli: http.DefaultClient,
	}, c)
}

func TestMisconfigDetection(t *testing.T) {
	r := require.New(t)

	c := New()
	_, err := c.Do(context.TODO(), nil)
	r.Error(err)
	r.Equal(ErrMisconfig, errors.Cause(err))
}
