package apierrors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	MissingField = NewOutput("missing_field", "A required field is missing.")
	NotFound     = NewOutput("not_found", "What you are looking for cannot be found.")
	BadRequest   = NewOutput("bad_request", "The request is invalid.")
	InvalidJSON  = NewOutput("invalid_json", "The JSON value provided is invalid.")
	Invalid      = NewOutput("invalid", "The value provided is invalid.")
	Internal     = NewOutput("internal_error", "An internal error occurred.")
)

type Output struct {
	Value string
	Desc  string
}

func NewOutput(value, desc string) *Output {
	return &Output{
		Value: value,
		Desc:  desc,
	}
}

func (o Output) String() string {
	return o.Value
}

func (o Output) ToUpper() string {
	return strings.ToUpper(o.String())
}

func (o Output) Description() string {
	return o.Desc
}

type FieldError struct {
	Field     string `json:"field,omitempty"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"message,omitempty"`
}

type Error struct {
	Message     string        `json:"message"`
	Resource    string        `json:"resource"`
	Description string        `json:"description"`
	Errors      []*FieldError `json:"errors,omitempty"`
}

func (e Error) Error() string {
	return e.Message
}

func write(w http.ResponseWriter, code int, data interface{}) {
	switch {
	case code < 200:
		panic(fmt.Sprintf("status code %d must be >= 200", code))
	case code == 204:
		w.WriteHeader(code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if data != nil {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "")
		enc.SetEscapeHTML(false)
		if err := enc.Encode(data); err != nil {
			fields := log.Fields{"data": data, "code": code}
			log.WithFields(fields).Errorf("%+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func Err(w http.ResponseWriter, e *Error, code int) {
	if e == nil {
		panic("error must not be empty")
	}
	write(w, code, e)
}

func (e *Error) AddInvalidError(field string) {
	e.Errors = append(e.Errors, &FieldError{
		ErrorCode: Invalid.String(),
		Field:     field,
	})
}

func New(message, resource string, output *Output, e ...*FieldError) *Error {
	err := &Error{
		Message:     message,
		Resource:    resource,
		Description: output.Description(),
	}
	if len(e) != 0 && e[0] != nil {
		err.Errors = e
	}
	return err
}

func NewInvalidError(message, resource, field string) *Error {
	return New(message, resource, Invalid, &FieldError{
		Field:     field,
		ErrorCode: Invalid.String(),
	})
}

func NewMissingFieldError(message, resource, field string) *Error {
	return New(message, resource, MissingField, &FieldError{
		Field:     field,
		ErrorCode: MissingField.String(),
	})
}

func NewInvalidJSONError(message, resource string) *Error {
	return New(message, resource, InvalidJSON, nil)
}

func NewInternalError(resource string) *Error {
	return New(Internal.ToUpper(), resource, Internal, nil)
}

func NewNotFoundError(resource string) *Error {
	return New(NotFound.ToUpper(), resource, NotFound, nil)
}

func OK200(w http.ResponseWriter, data interface{}) {
	write(w, 200, data)
}

func Created201(w http.ResponseWriter, data interface{}) {
	write(w, 201, data)
}

func BadRequest400(w http.ResponseWriter, resource, field string) {
	Err(w, NewInvalidError("BAD_REQUEST", resource, field), 400)
}

func NotFound404(w http.ResponseWriter, resource string) {
	Err(w, NewNotFoundError(resource), 404)
}

func InternalError500(w http.ResponseWriter, resource string, err error) {
	Err(w, NewInternalError(resource), 500)
}
