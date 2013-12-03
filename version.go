package helpers

import (
	"fmt"
)

var version = "0.0.0"

func VersionGet() string {
	return version
}

func VersionSet(major int, minor int, fixpack int) {
	version = fmt.Sprintf("%d.%d.%d", major, minor, fixpack)
}
