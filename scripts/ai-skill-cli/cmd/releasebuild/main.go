package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

type target struct {
	goos   string
	goarch string
}

var releaseTargets = []target{
	{goos: "windows", goarch: "amd64"},
	{goos: "darwin", goarch: "amd64"},
	{goos: "darwin", goarch: "arm64"},
	{goos: "linux", goarch: "amd64"},
	{goos: "linux", goarch: "arm64"},
}

func main() {
	version := flag.String("version", "dev", "version string embedded into ai-skill")
	commit := flag.String("commit", "unknown", "git commit embedded into ai-skill")
	date := flag.String("date", time.Now().UTC().Format(time.RFC3339), "build date embedded into ai-skill")
	dist := flag.String("dist", "dist", "artifact output directory")
	stableNames := flag.Bool("stable-names", false, "write stable artifact names without embedding the version in filenames")
	flag.Parse()

	if err := os.MkdirAll(*dist, 0o755); err != nil {
		fatal(err)
	}
	checksums := []string{}
	for _, item := range releaseTargets {
		name := artifactName(*version, item, *stableNames)
		path := filepath.Join(*dist, name)
		if err := buildArtifact(path, item, *version, *commit, *date); err != nil {
			fatal(err)
		}
		sum, err := sha256File(path)
		if err != nil {
			fatal(err)
		}
		checksums = append(checksums, fmt.Sprintf("%s  %s", sum, name))
		fmt.Printf("built %s\n", path)
	}
	sort.Strings(checksums)
	if err := os.WriteFile(filepath.Join(*dist, "SHA256SUMS"), []byte(strings.Join(checksums, "\n")+"\n"), 0o644); err != nil {
		fatal(err)
	}
}

func artifactName(version string, item target, stable bool) string {
	separator := "_"
	prefix := "ai-skill_" + version
	if stable {
		separator = "-"
		prefix = "ai-skill"
	}
	name := fmt.Sprintf("%s%s%s%s%s", prefix, separator, item.goos, separator, item.goarch)
	if item.goos == "windows" {
		name += ".exe"
	}
	return name
}

func buildArtifact(path string, item target, version string, commit string, date string) error {
	ldflags := strings.Join([]string{
		"-X", "github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/app.Version=" + version,
		"-X", "github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/app.Commit=" + commit,
		"-X", "github.com/linyihong/Ai-skill/scripts/ai-skill-cli/internal/app.Date=" + date,
	}, " ")
	cmd := exec.Command("go", "build", "-trimpath", "-ldflags", ldflags, "-o", path, "./cmd/ai-skill")
	cmd.Env = append(os.Environ(), "GOOS="+item.goos, "GOARCH="+item.goarch, "CGO_ENABLED=0")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if runtime.GOOS == "windows" {
		cmd.Env = append(cmd.Env, "GOFLAGS=-buildvcs=false")
	}
	return cmd.Run()
}

func sha256File(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:]), nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
