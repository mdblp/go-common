package portal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mdblp/go-common/clients/disc"
)

func TestGetPatientConfig(t *testing.T) {
	const resultOne = `{
		"_id": "5d7fa3760aad24585462b99e",
		"createdAt": "2019-09-16T15:00:06.348Z",
		"updatedAt": "2020-06-16T08:52:30.422Z",
		"__v": 0,
		"cgm": {
			"apiVersion": "0.1.0",
			"endOfLifeTransmitterDate": "2020-12-31T04:13:00+00:00",
			"expirationDate": "2020-12-31T04:13:00+00:00",
			"manufacturer": "Dexcom",
			"name": "G6",
			"swVersionTransmitter": "0.0.1",
			"transmitterId": "123456789"
		},
		"device": {
			"deviceId": "123456789-ID",
			"imei": "123456789-IMEI",
			"manufacturer": "Diabeloop",
			"name": "DBLG1",
			"swVersion": "1.0.0",
			"historyId": "5d7fa37a0aad24585462b9a0"
		},
		"parameters": {
			"values": [
				{
					"name": "PARAM_1",
					"value": "130",
					"unit": "%",
					"level": 1,
					"effectiveDate": "2020-06-10T09:25:25.000Z"
				},
				{
					"name": "PARAM_2",
					"value": "130",
					"unit": "%",
					"level": 1,
					"effectiveDate": "2020-06-10T10:00:00.000Z"
				}
			],
			"historyId": "5d7fa37a0aad24585462b9a2"
		},
		"pump": {
			"expirationDate": "2020-12-31T04:13:00+00:00",
			"manufacturer": "Vicentra",
			"name": "Kaleido",
			"serialNumber": "123456789",
			"swVersion": "0.1.0"
		},
		"time": "2020-06-10T10:00:00.000Z",
		"timezone": "Europe/Paris",
		"timezoneOffset": 120
	}`
	const resultTwo = `{
		"parameters": {
			"values": [
				{
					"name": "PARAM_1",
					"value": "130",
					"unit": "%",
					"level": "1",
					"effectiveDate": "2020-06-10T09:25:25.000Z"
				}
			]
		}
	}`

	srvr := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case routeV2GetPatientConfig:
			token := req.Header.Get("x-tidepool-session-token")
			res.Header().Set("content-type", "application/json")
			if token == "1-valid-json" {
				res.WriteHeader(200)
				fmt.Fprint(res, resultOne)
			} else if token == "2-invalid-json" {
				res.WriteHeader(200)
				fmt.Fprint(res, resultTwo)
			} else {
				res.WriteHeader(403)
				fmt.Fprint(res, `{"message":"Invalid token"}`)
			}
		default:
			res.WriteHeader(404)
		}
	}))
	defer srvr.Close()

	portalClient := NewPortalClientBuilder().
		WithHostGetter(disc.NewStaticHostGetterFromString(srvr.URL)).
		Build()

	pc, err := portalClient.GetPatientConfig("0")
	if pc != nil || err == nil {
		t.Errorf("Failed GetPatientConfig expected an error")
		return
	}

	pc, err = portalClient.GetPatientConfig("1-valid-json")
	if pc == nil || err != nil {
		t.Errorf("Failed GetPatientConfig with error[%v]", err)
		return
	}
	if pc == nil || pc.Device == nil || pc.Device.IMEI != "123456789-IMEI" {
		t.Errorf("Failed GetPatientConfig: Invalid returned value: %v", pc)
		return
	}

	pc, err = portalClient.GetPatientConfig("2-invalid-json")
	if pc != nil || err == nil {
		t.Errorf("Failed GetPatientConfig expected an error")
		return
	}
}
