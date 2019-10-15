package version

import (
	"fmt"
)

type version struct {
	initialized bool
	base        string
	shortCommit string
	fullCommit  string
}

func (v version) getBase() string {
	if v.base != "" {
		return v.base
	}
	return "N/A"
}

func (v version) getShortCommit() string {
	if v.shortCommit != "" {
		return v.shortCommit
	}
	return "N/A"
}

func (v version) getFullCommit() string {
	if v.fullCommit != "" {
		return v.fullCommit
	}
	return v.getShortCommit()
}

func (v version) GetShort() string {
	return fmt.Sprintf("%s+%s", v.getBase(), v.getShortCommit())
}

func (v version) String() string {
	return fmt.Sprintf("%s+%s", v.getBase(), v.getFullCommit())
}

/**
	Initialization wrapper
**/

// Singleton version instance
var instance version

// Variables to be injected at build time
var (
	releaseNumber string //Release number. i.e. 1.2.3
	shortCommit   string
	fullCommit    string //Full commit id. i.e. e0c73b95646559e9a3696d41711e918398d557fb
)

func init() {
	if !instance.initialized {
		instance = version{base: releaseNumber, fullCommit: fullCommit, shortCommit: shortCommit, initialized: true}
	}
}

// Initialization function, only used for tests
func forceInit() {
	instance = version{base: releaseNumber, fullCommit: fullCommit, shortCommit: shortCommit, initialized: true}
}

func Long() string {
	return instance.String()
}

func Short() string {
	return instance.GetShort()
}
