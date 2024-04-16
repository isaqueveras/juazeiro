package juazeiro

import (
	"fmt"
	"net"
	"net/http"
)

// Error ...
type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Method  string `json:"-"`
	Host    string `json:"-"`
	Cause   error  `json:"-"`
	Status  int    `json:"-"`
}

// Error enable Error type to implements error type
func (e *Error) Error() string {
	return fmt.Sprintf("Request failed to (%v) %v with status %v: %+v", e.Method, e.Host, e.Status, e.Cause)
}

// ToString ...
func (e *Error) ToString() string {
	var desc string
	switch e.Method {
	case http.MethodGet:
		desc = "fetch data"
	case http.MethodPost:
		desc = "register data"
	case http.MethodPut:
		desc = "update data"
	case http.MethodDelete:
		desc = "delete data"
	}

	var status string
	switch {
	case e.Status >= http.StatusMultipleChoices && e.Status < http.StatusBadRequest:
		status = "not accepting requests"
	case e.Status == http.StatusNotFound:
		status = "did not find the data"
	case e.Status >= http.StatusBadRequest && e.Status < http.StatusInternalServerError:
		status = "received incorrect parameters"
	case e.Status == http.StatusServiceUnavailable:
		status = "is unavailable"
	case e.Status >= http.StatusInternalServerError:
		status = "found an internal problem"
	case e.Status == http.StatusNotImplemented:
		status = "method not implemented"
	}

	if cause, ok := e.Cause.(net.Error); ok {
		if cause.Timeout() {
			status = "server took longer than expected"
		}
	}

	return fmt.Sprintf("It was not possible %s. The external service %s, try again.", desc, status)
}
