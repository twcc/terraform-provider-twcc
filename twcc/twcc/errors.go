package twcc

import (
    "fmt"
    "net/http"
)

// BaseError is an error type that all other error types embed.
type BaseError struct {
    DefaultErrString string
    Info             string
}

func (e BaseError) Error() string {
    e.DefaultErrString = "An error occurred while executing a request."
    return e.choseErrString()
}

func (e BaseError) choseErrString() string {
    if e.Info != "" {
        return e.Info
    }
    return e.DefaultErrString
}

// ErrUnexpectedResponseCode is returned by the Request method when a response code other than
// those listed in OkCodes is encountered.
type ErrUnexpectedResponseCode struct {
    BaseError
    URL            string
    Method         string
    Expected       []int
    Actual         int
    Body           []byte
    ResponseHeader http.Header
}

func (e ErrUnexpectedResponseCode) Error() string {
    e.DefaultErrString = fmt.Sprintf(
        "Expected HTTP response code %v when accessing [%s %s], but got %d instead\n%s",
        e.Expected, e.Method, e.URL, e.Actual, e.Body,
    )
    return e.choseErrString()
}

// GetStatusCode returns the actual status code of the error.
func (e ErrUnexpectedResponseCode) GetStatusCode() int {
	return e.Actual
}

// ErrDefault400 is the default error type returned on a 400 HTTP response code.
type ErrDefault400 struct {
	ErrUnexpectedResponseCode
}

// ErrDefault401 is the default error type returned on a 401 HTTP response code.
type ErrDefault401 struct {
        ErrUnexpectedResponseCode
}

// ErrDefault403 is the default error type returned on a 403 HTTP response code.
type ErrDefault403 struct {
        ErrUnexpectedResponseCode
}

// ErrDefault404 is the default error type returned on a 404 HTTP response code.
type ErrDefault404 struct {
    ErrUnexpectedResponseCode
}

// ErrDefault409 is the default error type returned on a 409 HTTP response code.
type ErrDefault409 struct {
        ErrUnexpectedResponseCode
}

// ErrDefault500 is the default error type returned on a 500 HTTP response code.
type ErrDefault500 struct {
        ErrUnexpectedResponseCode
}

// ErrDefault503 is the default error type returned on a 503 HTTP response code.
type ErrDefault503 struct {
        ErrUnexpectedResponseCode
}

func (e ErrDefault400) Error() string {
    e.DefaultErrString = fmt.Sprintf(
        "Bad request with: [%s %s], error message: %s",
        e.Method, e.URL, e.Body,
    )
    return e.choseErrString()
}

func (e ErrDefault401) Error() string {
    e.DefaultErrString = fmt.Sprintf(
        "Unauthorized with: [%s %s], error message: %s",
        e.Method, e.URL, e.Body,
    )
    return e.choseErrString()
}

func (e ErrDefault403) Error() string {
    e.DefaultErrString = fmt.Sprintf(
        "Permission denied with: [%s %s], error message: %s",
        e.Method, e.URL, e.Body,
    )
    return e.choseErrString()
}

func (e ErrDefault404) Error() string {
    e.DefaultErrString = fmt.Sprintf(
        "Resource not found with: [%s %s], error message: %s",
        e.Method, e.URL, e.Body,
    )
    return e.choseErrString()
}

func (e ErrDefault409) Error() string {
    e.DefaultErrString = fmt.Sprintf(
        "Resource conflict with: [%s %s], error message: %s",
        e.Method, e.URL, e.Body,
    )
    return e.choseErrString()
}

func (e ErrDefault500) Error() string {
    e.DefaultErrString = fmt.Sprintf(
        "Internal server error with: [%s %s], error message: %s",
        e.Method, e.URL, e.Body,
    )
    return e.choseErrString()
}

func (e ErrDefault503) Error() string {
    e.DefaultErrString = fmt.Sprintf(
        "Service Unavailable with: [%s %s], error message: %s",
        e.Method, e.URL, e.Body,
    )
    return e.choseErrString()
}

// Err400er is the interface resource error types implement to override the error message
// from a 400 error.
type Err400er interface {
    Error400(ErrUnexpectedResponseCode) error
}

// Err401er is the interface resource error types implement to override the error message
// from a 401 error.
type Err401er interface {
    Error401(ErrUnexpectedResponseCode) error
}

// Err403er is the interface resource error types implement to override the error message
// from a 403 error.
type Err403er interface {
    Error403(ErrUnexpectedResponseCode) error
}

// Err404er is the interface resource error types implement to override the error message
// from a 404 error.
type Err404er interface {
    Error404(ErrUnexpectedResponseCode) error
}

// Err409er is the interface resource error types implement to override the error message
// from a 409 error.
type Err409er interface {
    Error409(ErrUnexpectedResponseCode) error
}

// Err500er is the interface resource error types implement to override the error message
// from a 500 error.
type Err500er interface {
    Error500(ErrUnexpectedResponseCode) error
}

// Err503er is the interface resource error types implement to override the error message
// from a 503 error.
type Err503er interface {
    Error503(ErrUnexpectedResponseCode) error
}
