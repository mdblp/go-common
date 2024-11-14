package request

import (
	"context"
	"encoding/json"
	dblcontext "github.com/mdblp/go-common/v2/context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const validHost = "http://ici.ou.labas.org"
const invalidHost = ":thisIsnotAUrl"

func TestDefaultRequestBuilder(t *testing.T) {
	t.Run("Building a request with an empty url should log an error", func(t *testing.T) {
		requestBuilder := NewBuilder("", http.MethodGet)
		assert.Equal(t, "No client host defined", requestBuilder.errorMessage)
		assert.Equal(t, true, requestBuilder.buildError)
	})
	t.Run("Building a request with a malformed url should log an error", func(t *testing.T) {
		requestBuilder := NewBuilder(invalidHost, http.MethodGet)
		assert.Equal(t, "Unable to parse urlString [:thisIsnotAUrl]", requestBuilder.errorMessage)
		assert.Equal(t, true, requestBuilder.buildError)
	})
	t.Run("Building a request with a valid url should create a builder for a GET request", func(t *testing.T) {
		requestBuilder := NewBuilder(validHost, http.MethodGet)
		assert.Equal(t, http.MethodGet, requestBuilder.method)
		assert.Equal(t, &url.URL{Scheme: "http", Host: "ici.ou.labas.org"}, requestBuilder.baseUrl)
		assert.Equal(t, false, requestBuilder.buildError)
	})
}

func TestNewCustomRequest(t *testing.T) {
	requestBuilder := NewBuilder(validHost, http.MethodHead)
	assert.Equal(t, false, requestBuilder.buildError)
	assert.Equal(t, http.MethodHead, requestBuilder.method)
}

func TestNewDeleteRequest(t *testing.T) {
	requestBuilder := NewDeleteBuilder(validHost)
	assert.Equal(t, false, requestBuilder.buildError)
	assert.Equal(t, http.MethodDelete, requestBuilder.method)
}

func TestNewGetRequest(t *testing.T) {
	requestBuilder := NewGetBuilder(validHost)
	assert.Equal(t, false, requestBuilder.buildError)
	assert.Equal(t, http.MethodGet, requestBuilder.method)
}

func TestNewPostRequest(t *testing.T) {
	requestBuilder := NewPostBuilder(validHost)
	assert.Equal(t, false, requestBuilder.buildError)
	assert.Equal(t, http.MethodPost, requestBuilder.method)
}

func TestNewPutRequest(t *testing.T) {
	requestBuilder := NewPutBuilder(validHost)
	assert.Equal(t, false, requestBuilder.buildError)
	assert.Equal(t, http.MethodPut, requestBuilder.method)
}

func TestBuildError(t *testing.T) {
	_, err := NewBuilder("", http.MethodGet).Build(context.TODO())
	expectedError := RequestBuilderError("No client host defined")
	assert.Equal(t, expectedError, err)
}

func TestDefaultGetRequest(t *testing.T) {
	request, err := NewGetBuilder(validHost).Build(context.TODO())
	assert.Equal(t, nil, err)
	assert.Equal(t, http.MethodGet, request.Method)
	assert.Equal(t, &url.URL{Scheme: "http", Host: "ici.ou.labas.org", Path: ""}, request.URL)
	assert.Equal(t, nil, request.Body)
	assert.Equal(t, "", request.Header.Get("Authorization"))
}

func TestRequestBuilderWithTraceSessionId(t *testing.T) {
	ctx := dblcontext.SetTraceSessionId(context.TODO(), "ThisIsSessionId")
	request, err := NewGetBuilder(validHost).Build(ctx)
	assert.Equal(t, nil, err)
	assert.Equal(t, http.MethodGet, request.Method)
	assert.Equal(t, &url.URL{Scheme: "http", Host: "ici.ou.labas.org", Path: ""}, request.URL)
	assert.Equal(t, nil, request.Body)
	assert.Equal(t, "", request.Header.Get("Authorization"))
	assert.Equal(t, "ThisIsSessionId", request.Header.Get(traceSessionHeader))
}

func TestRequestBuilder_BuildWithAuthToken(t *testing.T) {
	t.Run("Building a request with a bearer token", func(t *testing.T) {
		request, err := NewGetBuilder(validHost).WithAuthToken("thisIsAGreatBearerToken").Build(context.TODO())
		assert.Equal(t, nil, err)
		assert.Equal(t, http.MethodGet, request.Method)
		assert.Equal(t, &url.URL{Scheme: "http", Host: "ici.ou.labas.org"}, request.URL)
		assert.Equal(t, nil, request.Body)
		assert.Equal(t, "Bearer thisIsAGreatBearerToken", request.Header.Get("Authorization"))
	})
	t.Run("Building a request with a tidepool token", func(t *testing.T) {
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9-ddskln546hfgr34"
		request, err := NewGetBuilder(validHost).WithAuthToken(token).Build(context.TODO())
		assert.Equal(t, nil, err)
		assert.Equal(t, http.MethodGet, request.Method)
		assert.Equal(t, &url.URL{Scheme: "http", Host: "ici.ou.labas.org"}, request.URL)
		assert.Equal(t, nil, request.Body)
		assert.Equal(t, token, request.Header.Get(LegacyTokenHeader))
	})
}

func TestRequestBuilder_BuildWithPath(t *testing.T) {
	request, err := NewGetBuilder(validHost).WithPath("user", "123456789").Build(context.TODO())
	assert.Equal(t, nil, err)
	assert.Equal(t, http.MethodGet, request.Method)
	assert.Equal(t, &url.URL{Scheme: "http", Host: "ici.ou.labas.org", Path: "/user/123456789"}, request.URL)
}

func TestRequestBuilder_BuildWithQueryParamArray(t *testing.T) {
	request, err := NewGetBuilder(validHost).
		WithQueryParamArray("tableName", []string{"toto", "titi"}).
		Build(context.TODO())
	assert.Equal(t, nil, err)
	assert.Equal(t, http.MethodGet, request.Method)
	assert.Equal(t, &url.URL{Scheme: "http", Host: "ici.ou.labas.org", Path: "", RawQuery: "tableName=toto&tableName=titi"}, request.URL)
}

func TestRequestBuilder_BuildWithQueryParams(t *testing.T) {
	request, err := NewGetBuilder(validHost).
		WithQueryParams(map[string]string{"key1": "toto", "key2": "titi"}).
		Build(context.TODO())
	assert.Equal(t, nil, err)
	assert.Equal(t, http.MethodGet, request.Method)
	assert.Equal(t, &url.URL{Scheme: "http", Host: "ici.ou.labas.org", Path: "", RawQuery: "key1=toto&key2=titi"}, request.URL)
}

func TestRequestEnd2End(t *testing.T) {
	token := "ThisIsAGreatToken"
	payload := "Hello world"
	var server = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			t.Errorf("Incorrect HTTP Method [%s]", req.Method)
		}
		// by default set in Auth header
		if req.Header.Get("Authorization") != "Bearer "+token {
			t.Errorf("auth token not correctly set")
		}
		if req.Header.Get("x-tidepool-trace-session") != "123456789456" {
			t.Errorf("trace session not correctly set")
		}
		if req.Header.Get("Content-Type") != "application/json" {
			t.Errorf("content-type not correctly set")
		}
		switch req.URL.Path {
		case "/metadata/1234/profile":
			var receivedBody *string
			if err := json.NewDecoder(req.Body).Decode(&receivedBody); err != nil {
				t.Error("Error decoding body")
			}
			assert.Equal(t, payload, *receivedBody)
			res.WriteHeader(http.StatusOK)
		default:
			t.Error("Something went wrong with the request")
		}
	}))

	defer server.Close()
	ctx := dblcontext.SetTraceSessionId(context.TODO(), "123456789456")

	request, err := NewPutBuilder(server.URL).
		WithAuthToken(token).
		WithPath("metadata", "1234", "profile").
		WithPayload(payload).
		Build(ctx)
	assert.Equal(t, nil, err)
	client := http.DefaultClient
	response, _ := client.Do(request)
	assert.Equal(t, "200 OK", response.Status)
}
