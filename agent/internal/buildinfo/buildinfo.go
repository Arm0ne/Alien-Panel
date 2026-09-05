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

// DisplayVersion is the released semantic version stored on the central
// service and shown to operators.
func DisplayVersion() string {
	version := strings.TrimSpace(Version)
	if version == "" {
		version = "dev"
	}
	return version
}

// BuildIdentity adds the source commit to the release version for diagnostic
// logs. It is intentionally separate from DisplayVersion so the UI does not
// expose an implementation-specific build identifier as the Agent version.
func BuildIdentity() string {
	version := DisplayVersion()
	commit := strings.TrimSpace(Commit)
	if commit == "" || strings.Contains(version, commit) {
		return version
	}
	return version + " (" + commit + ")"
}
