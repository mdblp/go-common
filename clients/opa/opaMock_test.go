package opa

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestMockGetOpaAuth(t *testing.T) {
	mock := NewMock()

	// Init mock data
	auth := &Authorization{
		Result: &opaResult{
			Authorized: true,
			Data:       map[string]interface{}{"test": 5, "test2": "1234"},
			Route:      "myopa-route",
		},
	}
	mock.SetMockOpaAuth("test0/url", auth, nil)

	auth = &Authorization{
		Result: &opaResult{
			Authorized: false,
			Data:       map[string]interface{}{"test": 6, "test2": "4567"},
			Route:      "myopa-route",
		},
	}
	mock.SetMockOpaAuth("test0/url2", auth, nil)

	// Test return
	var opaReq http.Request
	var err error
	opaReq.Host = "test0"
	opaReq.URL, err = url.Parse("http://test0/url")

	auth, err = mock.GetOpaAuth(&opaReq)
	if auth == nil || auth.Result == nil || auth.Result.Authorized != true || err != nil {
		t.Error("Invalid mock return for request 1")
		fmt.Printf("%v \n", err)
		return
	}
	opaReq.URL, err = url.Parse("http://test0/url2")
	auth, err = mock.GetOpaAuth(&opaReq)
	if auth == nil || auth.Result == nil || auth.Result.Authorized != false || err != nil {
		t.Error("Invalid mock return for request 2")
		return
	}
	opaReq.URL, err = url.Parse("http://test1/urlError")
	auth, err = mock.GetOpaAuth(&opaReq)
	if auth != nil || err == nil {
		t.Error("Invalid mock return for request 3")
		return
	}
}
