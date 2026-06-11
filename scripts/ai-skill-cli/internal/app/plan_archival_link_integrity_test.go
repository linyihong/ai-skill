package app

import (
	"reflect"
	"testing"
)

func TestParsePlanRenames(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []planRename
	}{
		{
			name: "active to archived",
			in:   "R100\tplans/active/foo.md\tplans/archived/foo.md\n",
			want: []planRename{{OldPath: "plans/active/foo.md", NewPath: "plans/archived/foo.md"}},
		},
		{
			name: "archived to active (reactivation)",
			in:   "R095\tplans/archived/bar.md\tplans/active/bar.md\n",
			want: []planRename{{OldPath: "plans/archived/bar.md", NewPath: "plans/active/bar.md"}},
		},
		{
			name: "multi-archive in same commit",
			in: "R100\tplans/active/a.md\tplans/archived/a.md\n" +
				"R092\tplans/active/b.md\tplans/archived/b.md\n",
			want: []planRename{
				{OldPath: "plans/active/a.md", NewPath: "plans/archived/a.md"},
				{OldPath: "plans/active/b.md", NewPath: "plans/archived/b.md"},
			},
		},
		{
			name: "non-plan rename ignored",
			in:   "R100\tdocs/old.md\tdocs/new.md\n",
			want: nil,
		},
		{
			name: "rename within active ignored",
			in:   "R100\tplans/active/old.md\tplans/active/new.md\n",
			want: nil,
		},
		{
			name: "additions and modifications ignored",
			in: "M\tplans/active/foo.md\n" +
				"A\tplans/archived/bar.md\n" +
				"R100\tplans/active/baz.md\tplans/archived/baz.md\n",
			want: []planRename{{OldPath: "plans/active/baz.md", NewPath: "plans/archived/baz.md"}},
		},
		{
			name: "empty input",
			in:   "",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePlanRenames(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePlanRenames mismatch\n  got:  %#v\n  want: %#v", got, tt.want)
			}
		})
	}
}

func TestIsPlanArchiveMove(t *testing.T) {
	cases := []struct {
		old, new string
		want     bool
	}{
		{"plans/active/foo.md", "plans/archived/foo.md", true},
		{"plans/archived/foo.md", "plans/active/foo.md", true},
		{"plans/active/old.md", "plans/active/new.md", false},
		{"plans/archived/a.md", "plans/archived/b.md", false},
		{"docs/x.md", "docs/y.md", false},
	}
	for _, c := range cases {
		if got := isPlanArchiveMove(c.old, c.new); got != c.want {
			t.Errorf("isPlanArchiveMove(%q, %q) = %v want %v", c.old, c.new, got, c.want)
		}
	}
}

func TestResolveRepoPath(t *testing.T) {
	cases := []struct {
		fromFile, target, want string
	}{
		{"plans/archived/a.md", "../active/sibling.md", "plans/active/sibling.md"},
		{"plans/archived/a.md", "sibling.md", "plans/archived/sibling.md"},
		{"plans/archived/a.md", "./sibling.md", "plans/archived/sibling.md"},
		{"docs/x.md", "../plans/active/foo.md", "plans/active/foo.md"},
		{"plans/active/a.md", "plans/archived/foo.md", "plans/active/plans/archived/foo.md"},
		{"plans/active/a.md", "", ""},
	}
	for _, c := range cases {
		got := resolveRepoPath(c.fromFile, c.target)
		if got != c.want {
			t.Errorf("resolveRepoPath(%q, %q) = %q want %q", c.fromFile, c.target, got, c.want)
		}
	}
}

func TestStripLinkFragment(t *testing.T) {
	cases := []struct{ in, want string }{
		{"foo.md", "foo.md"},
		{"foo.md#section", "foo.md"},
		{"#anchor", ""},
		{"../bar/baz.md#header", "../bar/baz.md"},
	}
	for _, c := range cases {
		if got := stripLinkFragment(c.in); got != c.want {
			t.Errorf("stripLinkFragment(%q) = %q want %q", c.in, got, c.want)
		}
	}
}

func TestPosixRel(t *testing.T) {
	cases := []struct {
		fromDir, toPath, want string
	}{
		{"plans/archived", "plans/archived/sibling.md", "sibling.md"},
		{"plans/archived", "plans/active/sibling.md", "../active/sibling.md"},
		{"docs", "plans/archived/foo.md", "../plans/archived/foo.md"},
		{"plans/archived/sub", "plans/archived/sub/x.md", "x.md"},
		{".", "plans/archived/foo.md", "plans/archived/foo.md"},
	}
	for _, c := range cases {
		got := posixRel(c.fromDir, c.toPath)
		if got != c.want {
			t.Errorf("posixRel(%q, %q) = %q want %q", c.fromDir, c.toPath, got, c.want)
		}
	}
}

func TestSuggestReplacement(t *testing.T) {
	cases := []struct {
		name                                string
		fromFile, newPath, originalTarget, want string
	}{
		{
			name:           "repo-rooted target uses new repo path",
			fromFile:       "metadata/x.yaml",
			newPath:        "plans/archived/foo.md",
			originalTarget: "plans/active/foo.md",
			want:           "plans/archived/foo.md",
		},
		{
			name:           "relative target gets relative suggestion",
			fromFile:       "plans/archived/a.md",
			newPath:        "plans/archived/sibling.md",
			originalTarget: "../active/sibling.md",
			want:           "sibling.md",
		},
		{
			name:           "multi-archive cross-reference same dir",
			fromFile:       "plans/archived/a.md",
			newPath:        "plans/archived/b.md",
			originalTarget: "../active/b.md",
			want:           "b.md",
		},
		{
			name:           "dotted relative target",
			fromFile:       "plans/archived/a.md",
			newPath:        "plans/archived/sibling.md",
			originalTarget: "./sibling.md",
			want:           "sibling.md",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := suggestReplacement(c.fromFile, c.newPath, c.originalTarget)
			if got != c.want {
				t.Errorf("suggestReplacement(%q, %q, %q) = %q want %q",
					c.fromFile, c.newPath, c.originalTarget, got, c.want)
			}
		})
	}
}
