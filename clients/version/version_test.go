package version

import (
	"testing"
)

func TestVersion(t *testing.T) {
	ReleaseNumber = "1.2.3"
	ShortCommit = "e0c73b9"
	FullCommit = "e0c73b95646559e9a3696d41711e918398d557fb"

	v := GetVersion()

	shortV := v.GetShort()
	longV := v.String()

	if shortV != "1.2.3+e0c73b9" {
		t.Errorf("Expected short version %s but got %s", "1.2.3+e0c73b9", shortV)
	}
	if longV != "1.2.3+e0c73b95646559e9a3696d41711e918398d557fb" {
		t.Errorf("Expected short version %s but got %s", "1.2.3+e0c73b95646559e9a3696d41711e918398d557fb", longV)
	}
}
