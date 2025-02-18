package endpoints

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/RIT3shSapata/todo-list-api/internal/log"
	"go.uber.org/zap"
)

// Endpoint implements ServeHTTP by adding the context of the http request into the
// handler for the endpoint by default, similar to router.HandleFunc(), will be using router.Handle()
type Endpoint struct {
	Handler func(context.Context, http.ResponseWriter, *http.Request)
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	e.Handler(ctx, w, r)
}

// a struct to handle all the writes to response writer
type Responder struct {
	Logger log.Logger
}

// encoderFn is a custom type for specifying a custom json encoding function
type encoderFn func(interface{}) ([]byte, error)

type responderOpts struct {
	statusCode int
	err        error
	encoder    encoderFn
	shouldLog  bool
}

// implementing the options pattern
type ResponderOption func(opts *responderOpts)

// Optional argument function for setting the response status code
func WithStatusCode(code int) ResponderOption {
	return func(options *responderOpts) {
		options.statusCode = code
	}
}

// Optional argument function for setting the returned error
func WithError(err error) ResponderOption {
	return func(options *responderOpts) {
		options.err = err
	}
}

// Optional argument function for logging the error
func WithLog() ResponderOption {
	return func(options *responderOpts) {
		options.shouldLog = true
	}
}

// Optional argument function for setting the returned error
func WithCustomEncoder(encoder encoderFn) ResponderOption {
	return func(options *responderOpts) {
		options.encoder = encoder
	}
}

// setContentType method is used to set the content type for response header
func setContentType(w http.ResponseWriter, body interface{}) {
	contentType := "text/plain; charset=utf-8"
	kind := reflect.TypeOf(body).Kind()

	switch kind {
	case reflect.Struct, reflect.Map, reflect.Ptr:
		contentType = "application/json; charset=utf-8"
	case reflect.String:
		contentType = "text/plain; charset=utf-8"
	default:
		if b, ok := body.([]byte); ok {
			contentType = http.DetectContentType(b)
		} else if kind == reflect.Slice {
			contentType = "application/json; charset=utf-8"
		}
	}

	w.Header().Set("Content-Type", contentType)
}

// Respond writes a response into a header. Takes a custom JSON encoder as an optional argument for when we need to customise the format of
// the http response.
func (r *Responder) Respond(ctx context.Context, w http.ResponseWriter, v interface{}, optionalArgs ...ResponderOption) {
	options := responderOpts{}
	for _, o := range optionalArgs {
		o(&options)
	}
	statusCode := options.statusCode
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	encoderFunc := options.encoder
	if encoderFunc == nil {
		encoderFunc = json.Marshal
	}
	var b []byte
	if v != nil {
		var err error
		if b, err = encoderFunc(v); err != nil {
			const msg = "failed to marshal payload to JSON"
			r.Logger.Error(msg, zap.Error(err))
		}

		setContentType(w, v)
	}
	w.WriteHeader(statusCode)
	if _, err := w.Write(b); err != nil {
		const msg = "failed to write response body"
		r.Logger.Error(msg, zap.Error(err))
	}
}

type ErrorResponse struct {
	Msg string `json:"msg"`
}

// Error responds to the request with the optional error and status code or just responds with the default 500
// if none are provided.
func (r *Responder) Error(ctx context.Context, w http.ResponseWriter, optionalArgs ...ResponderOption) {
	options := responderOpts{}
	for _, o := range optionalArgs {
		o(&options)
	}
	statusCode := options.statusCode
	err := options.err
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	if err == nil {
		err = errors.New("internal server error")
	}

	if options.shouldLog {
		r.Logger.Error("internal server error", zap.Error(err))
	}

	resp := ErrorResponse{
		Msg: err.Error(),
	}
	r.Respond(ctx, w, resp, WithStatusCode(statusCode))
}
