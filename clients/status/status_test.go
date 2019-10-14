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

func TestVersionDisplay(t *testing.T) {
	//set the application version
	version.VersionBase = "1.0.1"
	version.VersionFullCommit = "5fa7cdd19bd6eb8c9082c5f456f806d4cfd0f438"
	expectedVersion := "1.0.1+5fa7cdd19bd6eb8c9082c5f456f806d4cfd0f438"
	s := NewStatus(200, "OK")
	if s.Code != 200 {
		t.Error("Expected status code to be 200")
	}
	if s.Reason != "OK" {
		t.Error("Expected status reason to be 'OK'")
	}
	if s.Version != expectedVersion {
		t.Errorf("Expected the version to be %s but got %s", expectedVersion, s.Version)
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
	if *s.Error != 123 {
		t.Errorf("Expected error code to be 123, but got %d", *s.Error)
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
