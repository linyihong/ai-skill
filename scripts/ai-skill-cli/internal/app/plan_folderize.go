package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// runPlansFolderize converts flat multi-file plan clusters into folder layout.
func runPlansFolderize(args []string, stdout io.Writer, stderr io.Writer) int {
	opts := struct {
		root    string
		cluster string
		dryRun  bool
		apply   bool
	}{}
	fs := newFlagSet("plans folderize", stderr)
	fs.StringVar(&opts.root, "root", ".", "repository root")
	fs.StringVar(&opts.cluster, "cluster", "", "folderize only this plan base slug (default: all detected clusters)")
	fs.BoolVar(&opts.dryRun, "dry-run", false, "print planned moves without writing")
	fs.BoolVar(&opts.apply, "apply", false, "write folder layout and remove flat source files")
	if err := fs.Parse(args); err != nil {
		return ExitInvalidUsage
	}
	if opts.dryRun == opts.apply {
		_, _ = fmt.Fprintln(stderr, "exactly one of --dry-run or --apply is required")
		return ExitInvalidUsage
	}

	clusters := scanFlatPlanClusters(opts.root)
	if opts.cluster != "" {
		var filtered []flatPlanCluster
		for _, c := range clusters {
			if c.Base == opts.cluster {
				filtered = append(filtered, c)
			}
		}
		clusters = filtered
		if len(clusters) == 0 {
			_, _ = fmt.Fprintf(stderr, "no flat multi-file cluster found for %q\n", opts.cluster)
			return ExitValidationFailed
		}
	}
	if len(clusters) == 0 {
		_, _ = fmt.Fprintln(stdout, "no flat multi-file plan clusters detected")
		return ExitSuccess
	}

	var actions []string
	for _, c := range clusters {
		acts, err := planFolderizeCluster(opts.root, c, opts.dryRun)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "folderize %s: %v\n", c.Base, err)
			return ExitValidationFailed
		}
		actions = append(actions, acts...)
	}
	for _, a := range actions {
		_, _ = fmt.Fprintln(stdout, a)
	}
	if opts.dryRun {
		_, _ = fmt.Fprintf(stdout, "clusters=%d actions=%d (dry-run)\n", len(clusters), len(actions))
	} else {
		_, _ = fmt.Fprintf(stdout, "clusters=%d actions=%d (applied)\n", len(clusters), len(actions))
	}
	return ExitSuccess
}

// planFolderizeCluster performs or previews migration for one cluster.
func planFolderizeCluster(root string, c flatPlanCluster, dryRun bool) ([]string, error) {
	var log []string
	folderAbs := filepath.Join(root, filepath.FromSlash(c.folderRel()))
	mainAbs := filepath.Join(root, filepath.FromSlash(c.MainRel))
	mainData, err := os.ReadFile(mainAbs)
	if err != nil {
		return nil, fmt.Errorf("read main %s: %w", c.MainRel, err)
	}
	mainContent := rewriteClusterLinks(string(mainData), c)

	type move struct {
		srcRel string
		dstRel string
		body   string
	}
	var moves []move
	moves = append(moves, move{
		c.MainRel,
		filepath.ToSlash(filepath.Join(c.folderRel(), "_plan.md")),
		mainContent,
	})
	log = append(log, fmt.Sprintf("%s -> %s", moves[0].srcRel, moves[0].dstRel))

	for i, comp := range c.Companions {
		srcAbs := filepath.Join(root, filepath.FromSlash(comp.RelPath))
		data, readErr := os.ReadFile(srcAbs)
		if readErr != nil {
			return nil, fmt.Errorf("read companion %s: %w", comp.RelPath, readErr)
		}
		dstName := companionTargetName(i+1, comp.Suffix)
		dstRel := filepath.ToSlash(filepath.Join(c.folderRel(), dstName))
		moves = append(moves, move{comp.RelPath, dstRel, rewriteClusterLinks(string(data), c)})
		log = append(log, fmt.Sprintf("%s -> %s", comp.RelPath, dstRel))
	}

	if dryRun {
		return log, nil
	}

	if err := os.MkdirAll(folderAbs, 0o755); err != nil {
		return nil, err
	}
	for _, m := range moves {
		dstAbs := filepath.Join(root, filepath.FromSlash(m.dstRel))
		if err := os.WriteFile(dstAbs, []byte(m.body), 0o644); err != nil {
			return nil, fmt.Errorf("write %s: %w", m.dstRel, err)
		}
	}
	for _, m := range moves {
		srcAbs := filepath.Join(root, filepath.FromSlash(m.srcRel))
		if err := os.Remove(srcAbs); err != nil {
			return nil, fmt.Errorf("remove %s: %w", m.srcRel, err)
		}
	}
	return log, nil
}

// rewriteRepoPlanPathRefs replaces markdown links pointing at old flat plan paths.
func rewriteRepoPlanPathRefs(root, oldRel, newRel string) (int, error) {
	oldRel = filepath.ToSlash(oldRel)
	newRel = filepath.ToSlash(newRel)
	oldBase := filepath.Base(oldRel)
	newBase := filepath.Base(newRel)
	updated := 0
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		content := string(data)
		if !strings.Contains(content, oldRel) && !strings.Contains(content, oldBase) {
			return nil
		}
		newContent := strings.ReplaceAll(content, oldRel, newRel)
		newContent = strings.ReplaceAll(newContent, oldBase, newBase)
		if newContent == content {
			return nil
		}
		if err := os.WriteFile(path, []byte(newContent), 0o644); err != nil {
			return err
		}
		updated++
		return nil
	})
	return updated, err
}
