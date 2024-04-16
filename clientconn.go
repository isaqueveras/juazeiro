package juazeiro

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"
)

type ClientConnInterface interface {
	Invoke(ctx context.Context, method, uri string, status int, args interface{}, reply interface{}) error
}

var _ ClientConnInterface = (*ClientConn)(nil)

type ClientConn struct {
	clt *http.Client
	url *url.URL
}

// NewClient creates a client connection to the given target.
func NewClient(baseURL string) (*ClientConn, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &ClientConn{url: url, clt: &http.Client{}}, nil
}

func (c *ClientConn) Invoke(ctx context.Context, method, uri string, status int, args interface{}, reply interface{}) error {
	url, err := c.url.Parse(path.Join(c.url.Path, uri))
	if err != nil {
		return &Error{Host: c.url.Host, Method: method, Cause: err, Status: http.StatusBadRequest}
	}

	var body []byte
	if body, err = json.Marshal(&args); err != nil {
		return &Error{Host: c.url.Host, Method: method, Cause: err, Status: http.StatusBadRequest}
	}

	var inner *http.Request
	if inner, err = http.NewRequestWithContext(ctx, method, url.String(), bytes.NewBuffer(body)); err != nil {
		return &Error{Host: c.url.Host, Method: method, Cause: err, Status: http.StatusBadRequest}
	}

	var response *http.Response
	if response, err = c.clt.Do(inner); err != nil {
		return &Error{Host: c.url.Host, Method: method, Cause: err, Status: http.StatusServiceUnavailable}
	}
	defer response.Body.Close()

	if response.StatusCode == status {
		return json.NewDecoder(response.Body).Decode(reply)
	}

	var errServer = &Error{}
	if err = json.NewDecoder(response.Body).Decode(&errServer); err != nil {
		return &Error{Cause: err, Method: method, Host: c.url.Host, Status: response.StatusCode}
	}

	return &Error{Method: method, Host: c.url.Host, Code: errServer.Code, Status: response.StatusCode, Cause: errors.New(errServer.Message)}
}
