// Package buildinfo exposes the identity embedded in a released Agent binary.
package buildinfo

import "strings"

// These values are set by deploy/build-bundle.ps1 with go build -ldflags.
// Development builds deliberately report "dev" instead of pretending to be a
// released version.
var (
	Version   = "dev"
	Commit    = ""
	BuildTime = ""
)

// DisplayVersion is stable enough to store on the central service and show to
// operators. Version normally contains a Git describe value; commit is added
// only when it is not already part of that value.
func DisplayVersion() string {
	version := strings.TrimSpace(Version)
	if version == "" {
		version = "dev"
	}
	commit := strings.TrimSpace(Commit)
	if commit == "" || strings.Contains(version, commit) {
		return version
	}
	return version + " (" + commit + ")"
}
