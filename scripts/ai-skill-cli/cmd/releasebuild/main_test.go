package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/app"
)

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

func TestRepoLocalBinariesMatchChecksumsAndCurrentSource(t *testing.T) {
	moduleRoot := moduleRootForTest(t)
	binDir := filepath.Join(moduleRoot, "bin")
	buildInfo := readBuildInfo(t, filepath.Join(binDir, "BUILDINFO"))
	sourceCommit := buildInfo["source_commit"]
	if sourceCommit == "" {
		t.Fatalf("BUILDINFO missing source_commit")
	}
	latestSourceCommit := latestCLISourceCommit(t, moduleRoot)
	if sourceCommit != latestSourceCommit {
		t.Fatalf("repo-local binaries built from %s, latest CLI source commit is %s; rebuild bin/", sourceCommit, latestSourceCommit)
	}

	checksums := readChecksums(t, filepath.Join(binDir, "SHA256SUMS"))
	expectedNames := map[string]bool{}
	for _, item := range releaseTargets {
		expectedNames[artifactName("", item, true)] = true
	}
	if len(checksums) != len(expectedNames) {
		t.Fatalf("expected %d checksums, got %d", len(expectedNames), len(checksums))
	}
	for name, want := range checksums {
		if !expectedNames[name] {
			t.Fatalf("unexpected checksum entry: %s", name)
		}
		got := sha256Path(t, filepath.Join(binDir, name))
		if got != want {
			t.Fatalf("checksum mismatch for %s: got %s want %s", name, got, want)
		}
	}

	host := target{goos: runtime.GOOS, goarch: runtime.GOARCH}
	hostName := artifactName("", host, true)
	if !expectedNames[hostName] {
		t.Skipf("no repo-local binary for %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	hostBinary := filepath.Join(binDir, hostName)
	versionResult := runBinaryJSON(t, hostBinary, "version", "--json")
	if versionResult.Command != "version" || versionResult.Status != "success" {
		t.Fatalf("unexpected version result: %#v", versionResult)
	}
	if checkMessage(versionResult.Checks, "commit") != sourceCommit {
		t.Fatalf("binary commit metadata = %q, BUILDINFO source_commit = %q", checkMessage(versionResult.Checks, "commit"), sourceCommit)
	}
	doctorResult := runBinaryJSON(t, hostBinary, "doctor", "--json")
	if doctorResult.Command != "doctor" || doctorResult.Status != "success" {
		t.Fatalf("unexpected doctor result: %#v", doctorResult)
	}
}

func moduleRootForTest(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func readBuildInfo(t *testing.T, path string) map[string]string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	result := map[string]string{}
	for _, line := range strings.Split(string(content), "\n") {
		key, value, ok := strings.Cut(line, "=")
		if ok {
			result[key] = value
		}
	}
	return result
}

func latestCLISourceCommit(t *testing.T, moduleRoot string) string {
	t.Helper()
	output, err := exec.Command("git", "-C", moduleRoot, "log", "-1", "--format=%h", "--", "cmd", "internal", "go.mod", "go.sum").Output()
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(output))
}

func readChecksums(t *testing.T, path string) map[string]string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	result := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(string(content)), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 {
			result[fields[1]] = fields[0]
		}
	}
	return result
}

func sha256Path(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

func runBinaryJSON(t *testing.T, binary string, args ...string) app.Result {
	t.Helper()
	cmd := exec.Command(binary, args...)
	cmd.Dir = moduleRootForTest(t)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s failed: %v\n%s", binary, strings.Join(args, " "), err, string(output))
	}
	var result app.Result
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("decode JSON: %v\n%s", err, string(output))
	}
	return result
}

func checkMessage(checks []app.Check, name string) string {
	for _, check := range checks {
		if check.Name == name {
			return check.Message
		}
	}
	return ""
}
