package opa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetOpaAuth(t *testing.T) {
	const resultOne = `{
		"result": {
			"authorized": true,
			"data": {
				"userIds": [
					"00005"
				]
			},
			"route": "tidewhisperer-get"
		}
	}`
	const resultTwo = `{
		"result": {
			"authorized": false,
			"data": {
				"userIds": [
					"00004"
				]
			},
			"route": "tidewhisperer-get"
		}
	}`

	srvr := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case routeAuth:
			var inputRequest HTTPInput
			json.NewDecoder(req.Body).Decode(&inputRequest)

			res.Header().Set("content-type", "application/json")
			if inputRequest.Input.Request.Host == "authorized" {
				res.WriteHeader(200)
				fmt.Fprint(res, resultOne)
			} else if inputRequest.Input.Request.Host == "unauthorized" {
				res.WriteHeader(200)
				fmt.Fprint(res, resultTwo)
			} else {
				res.WriteHeader(403)
				fmt.Fprint(res, `{"message":"Invalid call"}`)
			}
		default:
			res.WriteHeader(404)
		}
	}))
	defer srvr.Close()

	opaClient := NewOpaClientBuilder().
		WithHost(srvr.URL).
		WithRequestingService("test").
		Build()
	var opaReq http.Request
	opaReq.Host = "authorized"
	opaReq.URL, _ = url.Parse("http://authorized/url2")
	auth, err := opaClient.GetOpaAuth(&opaReq)
	if auth == nil || err != nil {
		t.Errorf("Failed GetOpaAuth with error[%v]", err)
		return
	}
	if auth.Result == nil {
		t.Errorf("Failed GetOpaAuth: Invalid returned value: %v", auth)
		return
	}
	if auth.Result.Authorized != true || auth.Result.Route != "tidewhisperer-get" {
		t.Errorf("Failed GetOpaAuth: Invalid returned value: %v", auth)
		return
	}
	opaReq.Host = "unauthorized"
	auth, err = opaClient.GetOpaAuth(&opaReq)
	if auth == nil || err != nil {
		t.Errorf("Failed GetOpaAuth with error[%v]", err)
		return
	}
	if auth.Result == nil {
		t.Errorf("Failed GetOpaAuth: Invalid returned value: %v", auth)
		return
	}
	if auth.Result.Authorized != false || auth.Result.Route != "tidewhisperer-get" {
		t.Errorf("Failed GetOpaAuth: Invalid returned value: %v", auth)
		return
	}
	opaReq.Host = "error"
	auth, err = opaClient.GetOpaAuth(&opaReq)
	if auth != nil || err == nil {
		t.Errorf("Failed GetOpaAuth expected an error")
		return
	}
}
