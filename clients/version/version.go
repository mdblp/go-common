package version

import (
	"fmt"
)

var (
	VersionBase        string
	VersionShortCommit string
	VersionFullCommit  string
)

func Base() string {
	if VersionBase != "" {
		return VersionBase
	} else {
		return "N/A"
	}
}

func ShortCommit() string {
	if VersionShortCommit != "" {
		return VersionShortCommit
	} else {
		return "N/A"
	}
}

func FullCommit() string {
	if VersionFullCommit != "" {
		return VersionFullCommit
	} else {
		return ShortCommit()
	}
}

func Short() string {
	return fmt.Sprintf("%s+%s", Base(), ShortCommit())
}

func Long() string {
	return fmt.Sprintf("%s+%s", Base(), FullCommit())
}
