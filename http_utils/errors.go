package http_utils

import (
	"fmt"
	"log"
	"net/http"
)

type ServiceError interface {
	ErrorPublic() string
	ErrorInternal() string
}

var oauthErrorMap = map[int]string{
	http.StatusBadRequest:          "invalid_request",
	http.StatusUnauthorized:        "unauthorized_client",
	http.StatusForbidden:           "access_denied",
	http.StatusInternalServerError: "server_error",
	http.StatusServiceUnavailable:  "temporarily_unavailable",
}

// OAuthError is the JSON handler for OAuth2 error responses
type OAuthError struct {
	Err             string `json:"error"`
	Description     string `json:"error_description,omitempty"`
	InternalError   error  `json:"-"`
	InternalMessage string `json:"-"`
}

func (e *OAuthError) ErrorPublic() string {
	return fmt.Sprintf("%s: %s", e.Err, e.Description)
}
func (e *OAuthError) ErrorInternal() string {
	return fmt.Sprintf("%s: %s :%v", e.Err, e.Description, e.InternalError)
}

func (e *OAuthError) Error() string {
	if e.InternalMessage != "" {
		return e.InternalMessage
	}
	return fmt.Sprintf("%s: %s", e.Err, e.Description)
}

// WithInternalError adds internal error information to the error
func (e *OAuthError) WithInternalError(err error) *OAuthError {
	e.InternalError = err
	return e
}

// WithInternalMessage adds internal message information to the error
func (e *OAuthError) WithInternalMessage(fmtString string, args ...interface{}) *OAuthError {
	e.InternalMessage = fmt.Sprintf(fmtString, args...)
	return e
}

// Cause returns the root cause error
func (e *OAuthError) Cause() error {
	if e.InternalError != nil {
		return e.InternalError
	}
	return e
}

func OauthError(err string, description string) *OAuthError {
	return &OAuthError{Err: err, Description: description}
}

func BadRequestError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusBadRequest, fmtString, args...)
}

func InternalServerError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusInternalServerError, fmtString, args...)
}

func NotFoundError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusNotFound, fmtString, args...)
}

func UnauthorizedError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusUnauthorized, fmtString, args...)
}

func ForbiddenError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusForbidden, fmtString, args...)
}

func UnprocessableEntityError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusUnprocessableEntity, fmtString, args...)
}

// HTTPError is an error with a message and an HTTP status code.
type HTTPError struct {
	Code            int    `json:"code"`
	Message         string `json:"msg"`
	InternalError   error  `json:"-"`
	InternalMessage string `json:"-"`
	ErrorID         string `json:"error_id,omitempty"`
}

func (e *HTTPError) ErrorPublic() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func (e *HTTPError) ErrorInternal() string {
	return fmt.Sprintf("%d: %s :%s", e.Code, e.Message, e.InternalError)
}

func (e *HTTPError) Error() string {
	if e.InternalMessage != "" {
		return e.InternalMessage
	}
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// Cause returns the root cause error
func (e *HTTPError) Cause() error {
	if e.InternalError != nil {
		return e.InternalError
	}
	return e
}

// WithInternalError adds internal error information to the error
func (e *HTTPError) WithInternalError(err error) *HTTPError {
	e.InternalError = err
	return e
}

// WithInternalMessage adds internal message information to the error
func (e *HTTPError) WithInternalMessage(fmtString string, args ...interface{}) *HTTPError {
	e.InternalMessage = fmt.Sprintf(fmtString, args...)
	return e
}

func httpError(code int, fmtString string, args ...interface{}) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: fmt.Sprintf(fmtString, args...),
	}
}
func LogInternalError(err error) {
	switch e := err.(type) {
	case *HTTPError:
		//Log the message
		if e.InternalError != nil {

			log.Printf("HTTPError %v", e)
			log.Printf("Internal error %+v", e.InternalError)
		}
	case *OAuthError:
		if e.InternalError != nil {
			log.Printf("OAuthError %v", e)
			log.Printf("Internal error %+v", e.InternalError)
		}
	default:
		log.Printf("Unhandled server error: %v", err)
		// hide real error details from response to prevent info leaks
	}
}

func SendErrorResponse(err error, w http.ResponseWriter, r *http.Request) {

	var jsonerr error

	switch e := err.(type) {
	case *HTTPError:
		log.Printf("HTTPError: code %v message %v internal error %v internal message %v", e.Code, e.Message, e.InternalError, e.InternalMessage)
		if e.Code == http.StatusInternalServerError {
			log.Printf("Error Stack: %+v", e)
		}
		jsonerr = SendJSON(w, e.Code, e)
	case *OAuthError:
		log.Printf("OAuthError: %v, description %v  internal error %v internal message %v", e.Err, e.Description, e.InternalError, e.InternalMessage)

		jsonerr = SendJSON(w, http.StatusBadRequest, e)
	default:
		log.Printf("Unknown error type: %+v", e)
		// hide real error details from response to prevent info leaks
		w.WriteHeader(http.StatusInternalServerError)
		_, jsonerr = w.Write([]byte(`{"code":500,"msg":"Internal server error"}`))
	}

	if jsonerr != nil {
		log.Printf("\n error while sending error message %v orig error: %v", jsonerr, err)
	}

}
