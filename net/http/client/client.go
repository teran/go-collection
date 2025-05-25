// Package client implements HTTP utilities to ease multiple HTTP operations
//
// Deprecated: since first library version was implemented long ago Go ecosystem
// grew enough with more mature, tested and stable libraries. Please consider using
// another library for your purpose.
//
// In favor of backward compatibility this code won't be removed but it's frozen
// for adding new features. Only major issues will be fixed.
package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	method      string
	uri         string
	path        string
	queryParams url.Values
	headers     http.Header
	cli         *http.Client
}

func New() Client {
	return NewWithHTTPClient(http.DefaultClient)
}

func NewWithHTTPClient(cli *http.Client) Client {
	return Client{
		cli:         cli,
		headers:     make(http.Header),
		queryParams: make(url.Values),
	}
}

func (c Client) Clone() Client {
	h := make(http.Header)

	for k, v := range c.headers {
		h[k] = append(h[k], v...)
	}

	qp := make(url.Values)
	for k, v := range c.queryParams {
		qp[k] = append(qp[k], v...)
	}

	return Client{
		cli:         c.cli,
		method:      c.method,
		uri:         c.uri,
		path:        c.path,
		queryParams: qp,
		headers:     h,
	}
}

func (c Client) Base(uri string) Client {
	newc := c.Clone()
	newc.uri = uri
	return newc
}

func (c Client) Basef(uri string, args ...any) Client {
	newc := c.Clone()
	newc.uri = fmt.Sprintf(uri, args...)
	return newc
}

func (c Client) Header(name, value string) Client {
	n := c.Clone()
	n.headers.Set(name, value)
	return n
}

func (c Client) UserAgent(name string) Client {
	return c.Header("User-Agent", name)
}

func (c Client) Auth(kind, credential string) Client {
	return c.Header("Authorization", fmt.Sprintf("%s %s", kind, credential))
}

func (c Client) BasicAuth(username, password string) Client {
	cred := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	return c.Auth("Basic", cred)
}

func (c Client) Request(method, path string, args ...any) Client {
	newc := c.Clone()
	newc.method = method
	newc.path = fmt.Sprintf(path, args...)
	return newc
}

func (c Client) Get(path string, args ...any) Client {
	return c.Request(http.MethodGet, path, args...)
}

func (c Client) Post(path string, args ...any) Client {
	return c.Request(http.MethodPost, path, args...)
}

func (c Client) Put(path string, args ...any) Client {
	return c.Request(http.MethodPut, path, args...)
}

func (c Client) Delete(path string, args ...any) Client {
	return c.Request(http.MethodDelete, path, args...)
}

func (c Client) Head(path string, args ...any) Client {
	return c.Request(http.MethodHead, path, args...)
}

func (c Client) Options(path string, args ...any) Client {
	return c.Request(http.MethodOptions, path, args...)
}

func (c Client) QueryParam(key, value string) Client {
	newc := c.Clone()
	newc.queryParams.Add(key, value)
	return newc
}

func (c Client) Do(ctx context.Context, body io.Reader) (*http.Response, error) {
	if c.uri == "" {
		return nil, errors.Wrap(ErrMisconfig, "Base() should be called before Do()")
	}

	uri := c.uri + c.path
	if len(c.queryParams) > 0 {
		uri += "?" + c.queryParams.Encode()
	}

	log.WithFields(log.Fields{
		"uri":     uri,
		"method":  c.method,
		"headers": c.headers,
	}).Trace("sending request")

	req, err := http.NewRequestWithContext(ctx, c.method, uri, body)
	if err != nil {
		return nil, err
	}

	req.Header = c.headers.Clone()

	return c.cli.Do(req)
}

func (c Client) DoJSON(ctx context.Context, body io.Reader, model, errorModel any) (code int, err error) {
	resp, err := c.Do(ctx, body)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	log.WithFields(log.Fields{
		"status":  resp.StatusCode,
		"headers": resp.Header,
	}).Trace("response received")

	mime, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return resp.StatusCode, err
	}

	if strings.ToLower(mime) != "application/json" {
		return resp.StatusCode, errors.Wrapf(ErrUnsupportedMediaType, "expected application/json but got %s", mime)
	}

	decoder := json.NewDecoder(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return resp.StatusCode, decoder.Decode(model)
	}
	return resp.StatusCode, decoder.Decode(errorModel)
}
