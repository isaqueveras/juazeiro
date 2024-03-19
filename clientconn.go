package juazeiro

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type ClientConnInterface interface {
	Invoke(ctx context.Context, method, uri string, args interface{}, reply interface{}) error
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

func (c *ClientConn) Invoke(ctx context.Context, method, uri string, args interface{}, reply interface{}) error {
	url, err := c.url.Parse(path.Join(c.url.Path, uri))
	if err != nil {
		return err
	}

	var body []byte
	if body, err = json.Marshal(&args); err != nil {
		return err
	}

	var inner *http.Request
	if inner, err = http.NewRequestWithContext(ctx, method, url.String(), bytes.NewBuffer(body)); err != nil {
		return err
	}

	if err = c.do(inner, reply); err != nil {
		return err
	}

	return nil
}

func (c *ClientConn) do(request *http.Request, reply interface{}) error {
	resp, err := c.clt.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNoContent:
	case http.StatusOK, http.StatusCreated, http.StatusAccepted,
		http.StatusNonAuthoritativeInfo, http.StatusPartialContent:
		switch value := reply.(type) {
		case *[]byte:
			if value == nil {
				return errors.New("invalid data type, try *[]byte")
			}
			if *value, err = ioutil.ReadAll(resp.Body); err != nil {
				return err
			}
			return resp.Body.Close()
		default:
			return json.NewDecoder(resp.Body).Decode(reply)
		}
	default:
		if err = json.NewDecoder(resp.Body).Decode(err); err == nil {
			return err
		}
	}

	return nil
}
