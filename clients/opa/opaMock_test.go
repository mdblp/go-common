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
	auth := mock.GetMockedAuth(true, map[string]interface{}{"test": 5, "test2": "1234"}, "myopa-route1")
	mock.SetMockOpaAuth("test0/url", &auth, nil)

	auth2 := mock.GetMockedAuth(false, map[string]interface{}{"test": 6, "test2": "4567"}, "myopa-route2")
	mock.SetMockOpaAuth("test0/url2", &auth2, nil)

	// Test return
	var opaReq http.Request
	var err error
	opaReq.Host = "test0"
	opaReq.URL, err = url.Parse("http://test0/url")

	testAuth, err := mock.GetOpaAuth(&opaReq)
	if testAuth == nil || testAuth.Result == nil || testAuth.Result.Authorized != true || err != nil {
		t.Error("Invalid mock return for request 1")
		fmt.Printf("%v \n", err)
		fmt.Printf("%v \n", testAuth)
		fmt.Printf("%v \n", testAuth.Result)
		fmt.Printf("%v \n", testAuth.Result.Authorized)
		fmt.Printf("%v \n", testAuth.Result.Route)
		return
	}
	opaReq.URL, err = url.Parse("http://test0/url2")
	testAuth, err = mock.GetOpaAuth(&opaReq)
	if testAuth == nil || testAuth.Result == nil || testAuth.Result.Authorized != false || err != nil {
		t.Error("Invalid mock return for request 2")
		return
	}
	opaReq.URL, err = url.Parse("http://test1/urlError")
	testAuth, err = mock.GetOpaAuth(&opaReq)
	if testAuth != nil || err == nil {
		t.Error("Invalid mock return for request 3")
		return
	}
}
