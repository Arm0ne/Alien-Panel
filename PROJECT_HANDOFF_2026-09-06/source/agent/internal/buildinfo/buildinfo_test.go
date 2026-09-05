package buildinfo

import "testing"

func TestDisplayVersion(t *testing.T) {
	originalVersion, originalCommit := Version, Commit
	t.Cleanup(func() {
		Version, Commit = originalVersion, originalCommit
	})

	Version, Commit = "v1.0.1", "abc123"
	if got := DisplayVersion(); got != "v1.0.1" {
		t.Fatalf("DisplayVersion() = %q, want v1.0.1", got)
	}
	if got := BuildIdentity(); got != "v1.0.1 (abc123)" {
		t.Fatalf("BuildIdentity() = %q, want v1.0.1 (abc123)", got)
	}
	Version, Commit = "", ""
	if got := DisplayVersion(); got != "dev" {
		t.Fatalf("DisplayVersion() = %q, want dev", got)
	}
	if got := BuildIdentity(); got != "dev" {
		t.Fatalf("BuildIdentity() = %q, want dev", got)
	}
}
