package buildinfo

import "testing"

func TestDisplayVersion(t *testing.T) {
	originalVersion, originalCommit := Version, Commit
	t.Cleanup(func() {
		Version, Commit = originalVersion, originalCommit
	})

	Version, Commit = "release-abc123", "abc123"
	if got := DisplayVersion(); got != "release-abc123" {
		t.Fatalf("DisplayVersion() = %q, want release-abc123", got)
	}
	Version, Commit = "v0.1.0", "abc123"
	if got := DisplayVersion(); got != "v0.1.0 (abc123)" {
		t.Fatalf("DisplayVersion() = %q, want v0.1.0 (abc123)", got)
	}
	Version, Commit = "", ""
	if got := DisplayVersion(); got != "dev" {
		t.Fatalf("DisplayVersion() = %q, want dev", got)
	}
}
