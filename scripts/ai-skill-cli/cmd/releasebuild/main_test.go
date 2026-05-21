package main

import "testing"

func TestArtifactNameSupportsVersionedAndStableNames(t *testing.T) {
	item := target{goos: "darwin", goarch: "arm64"}
	if got := artifactName("v1", item, false); got != "ai-skill_v1_darwin_arm64" {
		t.Fatalf("unexpected versioned name: %s", got)
	}
	if got := artifactName("v1", item, true); got != "ai-skill-darwin-arm64" {
		t.Fatalf("unexpected stable name: %s", got)
	}
	win := target{goos: "windows", goarch: "amd64"}
	if got := artifactName("v1", win, true); got != "ai-skill-windows-amd64.exe" {
		t.Fatalf("unexpected Windows stable name: %s", got)
	}
}
