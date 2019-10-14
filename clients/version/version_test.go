package version_test

import (
	"testing"

	"github.com/tidepool-org/go-common/clients/version"
)

func TestVersion(t *testing.T) {
	version.VersionBase = "1.2.3"
	version.VersionShortCommit = "e0c73b9"
	version.VersionFullCommit = "e0c73b95646559e9a3696d41711e918398d557fb"

	shortV := version.Short()
	longV := version.Long()

	if shortV != "1.2.3+e0c73b9" {
		t.Errorf("Expected short version %s but got %s", "1.2.3+e0c73b9", shortV)
	}
	if longV != "1.2.3+e0c73b95646559e9a3696d41711e918398d557fb" {
		t.Errorf("Expected short version %s but got %s", "1.2.3+e0c73b95646559e9a3696d41711e918398d557fb", longV)
	}
}
