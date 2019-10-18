package status

import (
	"net/http"
	"testing"

	"github.com/tidepool-org/go-common/clients/version"
)

func TestNewStatus(t *testing.T) {
	s := NewStatus(200, "OK")
	if s.Code != 200 {
		t.Error("Expected status code to be 200")
	}
	if s.Reason != "OK" {
		t.Error("Expected status reason to be 'OK'")
	}
}

func TestDefaultReason(t *testing.T) {
	s := NewStatus(200, "")
	if s.Reason != "OK" {
		t.Error("Expected status reason to be 'OK'")
	}
}

func TestStatusWithError(t *testing.T) {
	s := NewStatusWithError(504, 123, "Internal error")
	if s.Code != 504 {
		t.Error("Expected status code to be 504")
	}
	if s.Reason != "Internal error" {
		t.Error("Expected status reason to be 'Internal error'")
	}
}

func TestNewStatusf(t *testing.T) {
	s := NewStatusf(404, "Ressource %d was not found on route %s", 123, "users")
	expectedReason := "Ressource 123 was not found on route users"
	if s.Code != 404 {
		t.Error("Expected status code to be 404")
	}
	if s.Reason != expectedReason {
		t.Errorf("Expected status reason to be '%s' but got %s", expectedReason, s.Reason)
	}
}

func TestNewStatusFromResponse(t *testing.T) {
	s := StatusFromResponse(&http.Response{StatusCode: 200, Status: "OK"})
	if s.Code != 200 {
		t.Error("Expected status code to be 200")
	}
	if s.Reason != "OK" {
		t.Error("Expected status reason to be 'OK'")
	}
}

func TestNewApiStatus(t *testing.T) {
	//set the application version
	version.ReleaseNumber = "1.2.3"
	version.FullCommit = "e0c73b95646559e9a3696d41711e918398d557fb"
	expectedVersion := "1.2.3+e0c73b95646559e9a3696d41711e918398d557fb"
	s := NewApiStatus(200, "OK")
	if s.Status.Code != 200 {
		t.Error("Expected status code to be 200")
	}
	if s.Status.Reason != "OK" {
		t.Error("Expected status reason to be 'OK'")
	}
	if s.Version != expectedVersion {
		t.Errorf("Expected the version to be %s but got %s", expectedVersion, s.Version)
	}
}
