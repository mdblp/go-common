package version

import (
	"testing"
)

func TestVersion(t *testing.T) {
	releaseNumber = "1.2.3"
	shortCommit = "e0c73b9"
	fullCommit = "e0c73b95646559e9a3696d41711e918398d557fb"

	forceInit()
	shortV := Short()
	longV := Long()

	if shortV != "1.2.3+e0c73b9" {
		t.Errorf("Expected short version %s but got %s", "1.2.3+e0c73b9", shortV)
	}
	if longV != "1.2.3+e0c73b95646559e9a3696d41711e918398d557fb" {
		t.Errorf("Expected short version %s but got %s", "1.2.3+e0c73b95646559e9a3696d41711e918398d557fb", longV)
	}
}
