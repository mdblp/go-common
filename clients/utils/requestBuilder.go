package utils

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

func defaultBuilder(host string) *RequestBuilder {
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
	return &RequestBuilder{baseUrl: baseUrl, buildError: false, method: http.MethodGet}
}

func NewGetRequest(host string) *RequestBuilder {
	return defaultBuilder(host)
}

func NewPostRequest(host string) *RequestBuilder {
	builder := defaultBuilder(host)
	builder.method = http.MethodPost
	return builder
}

func NewPutRequest(host string) *RequestBuilder {
	builder := defaultBuilder(host)
	builder.method = http.MethodPut
	return builder
}

func NewDeleteRequest(host string) *RequestBuilder {
	builder := defaultBuilder(host)
	builder.method = http.MethodDelete
	return builder
}

func NewCustomRequest(host string, method string) *RequestBuilder {
	builder := defaultBuilder(host)
	builder.method = method
	return builder
}

func (b *RequestBuilder) WithPath(pathParams ...string) *RequestBuilder {
	pathFragments := append([]string{b.baseUrl.Path}, pathParams...)
	b.baseUrl.Path = path.Join(pathFragments...)
	return b
}

func (b *RequestBuilder) WithPayload(payload interface{}) *RequestBuilder {
	b.payload = payload
	return b
}

func (b *RequestBuilder) WithAuthToken(token string) *RequestBuilder {
	b.token = token
	return b
}

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
