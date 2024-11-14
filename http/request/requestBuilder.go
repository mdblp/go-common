package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mdblp/go-common/v2/blperr"
	dblcontext "github.com/mdblp/go-common/v2/context"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const traceSessionHeader = "x-tidepool-trace-session"
const apiErrorKind = "request-builder"
const LegacyTokenHeader = "x-tidepool-session-token"

func RequestBuilderError(err string) blperr.StackError {
	details := map[string]interface{}{}
	details["error"] = err
	return blperr.NewWithDetails(apiErrorKind, "failed to build the http request", details)
}

type RequestBuilder struct {
	errorMessage string
	buildError   bool
	baseUrl      *url.URL
	method       string
	token        string
	payload      interface{}
}

func NewBuilder(host string, method string) *RequestBuilder {
	defaultUrl, _ := url.Parse("http://go.nowhere")
	if host == "" {
		return &RequestBuilder{
			baseUrl:      defaultUrl,
			buildError:   true,
			errorMessage: "No client host defined"}
	}
	baseUrl, err := url.Parse(host)
	if err != nil {
		return &RequestBuilder{
			baseUrl:      defaultUrl,
			buildError:   true,
			errorMessage: fmt.Sprintf("Unable to parse urlString [%s]", host)}
	}
	return &RequestBuilder{baseUrl: baseUrl, buildError: false, method: method}
}

// NewGetBuilder initialize a request builder for a 'GET' request
func NewGetBuilder(host string) *RequestBuilder {
	builder := NewBuilder(host, http.MethodGet)
	return builder
}

// NewPostBuilder initialize a request builder for a 'POST' request
func NewPostBuilder(host string) *RequestBuilder {
	builder := NewBuilder(host, http.MethodPost)
	return builder
}

// NewPutBuilder initialize a request builder for a 'PUT' request
func NewPutBuilder(host string) *RequestBuilder {
	builder := NewBuilder(host, http.MethodPut)
	return builder
}

// NewDeleteBuilder initialize a request builder for a 'DELETE' request
func NewDeleteBuilder(host string) *RequestBuilder {
	builder := NewBuilder(host, http.MethodDelete)
	return builder
}

// WithPath sets the path of the request url.
// For example: WithPath("abc", "456") will generate a url like http://host/abc/456
func (b *RequestBuilder) WithPath(pathParams ...string) *RequestBuilder {
	pathFragments := append([]string{b.baseUrl.Path}, pathParams...)
	b.baseUrl.Path = path.Join(pathFragments...)
	return b
}

// WithPayload add an object which will be serialized and add in the request
func (b *RequestBuilder) WithPayload(payload interface{}) *RequestBuilder {
	b.payload = payload
	return b
}

// WithAuthToken add an authentication token in the request
func (b *RequestBuilder) WithAuthToken(token string) *RequestBuilder {
	b.token = token
	return b
}

// WithQueryParams allows you to add multiple query params (in the form of key/value pairs) in the request
// For example: WithQueryParams(map[string]string{"key1": "val1", "key2": "val2"}) will generate a url like http://host?key1=val1&key2=val2
func (b *RequestBuilder) WithQueryParams(queryParams map[string]string) *RequestBuilder {
	q := b.baseUrl.Query()
	for key, value := range queryParams {
		if value != "" {
			q.Set(key, value)
		}
	}
	b.baseUrl.RawQuery = q.Encode()
	return b
}

// WithQueryParamArray allows you to add an array as an url query param in the request
// For example: WithQueryParamArray("colors", []string{"red", "green"}) will generate a url like http://host?colors=red&colors=green
// There is no http standard to pass arrays as query params, but this is a common convention and this how Gin (our preferred http server engine) understand it.
func (b *RequestBuilder) WithQueryParamArray(tableName string, values []string) *RequestBuilder {
	q := b.baseUrl.Query()
	for _, value := range values {
		if value != "" {
			q.Add(tableName, value)
		}
	}
	b.baseUrl.RawQuery = q.Encode()
	return b
}

// Build instantiates the http.request based on the parameters provided to the builder previously
func (b *RequestBuilder) Build(ctx context.Context) (*http.Request, error) {
	var err error
	if b.buildError {
		return nil, RequestBuilderError(b.errorMessage)
	}
	var req *http.Request
	if b.payload != nil {
		body, err := json.Marshal(b.payload)
		if err != nil {
			return nil, RequestBuilderError(err.Error())
		}
		req, err = http.NewRequestWithContext(ctx, b.method, b.baseUrl.String(), bytes.NewBuffer(body))
	} else {

		req, err = http.NewRequestWithContext(ctx, b.method, b.baseUrl.String(), nil)
	}
	if err != nil {
		return nil, RequestBuilderError(err.Error())
	}
	if b.token != "" {
		setAuthHeader(req, b.token)
	}
	req.Header.Set("Content-Type", "application/json")
	if traceSessionId, ok := dblcontext.GetTraceSessionId(ctx); ok {
		req.Header.Set(traceSessionHeader, traceSessionId)
	}

	return req, nil
}

func setAuthHeader(req *http.Request, authToken string) {
	/*eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9 is the default shoreline token header*/
	if strings.Index(authToken, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9") == 0 {
		req.Header.Add(LegacyTokenHeader, authToken)
	} else {
		req.Header.Add("Authorization", "Bearer "+authToken)
	}
}
