package app

import (
	"reflect"
	"testing"
)

func TestExtractMarkdownLinks(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []Link
	}{
		{
			name: "plain inline link",
			in:   `See [foo](../foo.md) for details.`,
			want: []Link{{Target: "../foo.md", Line: 1, Column: 5}},
		},
		{
			name: "titled inline link",
			in:   `See [foo](../foo.md "title text") for details.`,
			want: []Link{{Target: "../foo.md", Line: 1, Column: 5}},
		},
		{
			name: "escaped parens in path",
			in:   `[a](../a\(b\).md)`,
			want: []Link{{Target: "../a(b).md", Line: 1, Column: 1}},
		},
		{
			name: "multiple links on same line",
			in:   `[a](a.md) and [b](b.md)`,
			want: []Link{
				{Target: "a.md", Line: 1, Column: 1},
				{Target: "b.md", Line: 1, Column: 15},
			},
		},
		{
			name: "absolute URL filtered",
			in:   `[home](https://example.com) and [local](./x.md)`,
			want: []Link{{Target: "./x.md", Line: 1, Column: 33}},
		},
		{
			name: "anchor filtered",
			in:   `[a](#section) and [b](./x.md)`,
			want: []Link{{Target: "./x.md", Line: 1, Column: 19}},
		},
		{
			name: "mailto filtered",
			in:   `[mail](mailto:foo@example.com) and [b](./x.md)`,
			want: []Link{{Target: "./x.md", Line: 1, Column: 36}},
		},
		{
			name: "reference-style link ignored",
			in:   "[foo][ref]\n\n[ref]: ../foo.md\n",
			want: nil,
		},
		{
			name: "shortcut reference ignored",
			in:   `Just [foo] mentioned.`,
			want: nil,
		},
		{
			name: "html anchor ignored",
			in:   `Visit <a href="../foo.md">foo</a>.`,
			want: nil,
		},
		{
			name: "autolink ignored",
			in:   `<https://example.com>`,
			want: nil,
		},
		{
			name: "code-fence excluded (backticks)",
			in:   "Before [a](a.md)\n```\n[b](b.md)\n```\nAfter [c](c.md)",
			want: []Link{
				{Target: "a.md", Line: 1, Column: 8},
				{Target: "c.md", Line: 5, Column: 7},
			},
		},
		{
			name: "code-fence excluded (tildes)",
			in:   "[a](a.md)\n~~~\n[b](b.md)\n~~~\n[c](c.md)",
			want: []Link{
				{Target: "a.md", Line: 1, Column: 1},
				{Target: "c.md", Line: 5, Column: 1},
			},
		},
		{
			name: "code-fence with leading spaces",
			in:   "[a](a.md)\n  ```\n[b](b.md)\n  ```\n[c](c.md)",
			want: []Link{
				{Target: "a.md", Line: 1, Column: 1},
				{Target: "c.md", Line: 5, Column: 1},
			},
		},
		{
			name: "tilde fence does not close backtick fence",
			in:   "```\n[hidden](h.md)\n~~~\n[still hidden](h2.md)\n```\n[visible](v.md)",
			want: []Link{{Target: "v.md", Line: 6, Column: 1}},
		},
		{
			name: "unmatched bracket does not panic",
			in:   `Text [foo without close paren`,
			want: nil,
		},
		{
			name: "unmatched paren in link does not capture",
			in:   `[foo](unclosed`,
			want: nil,
		},
		{
			name: "multiline document with mixed links",
			in: "# Heading\n\n" +
				"Para with [link1](a/b.md) inline.\n\n" +
				"```go\n[not-a-link](nope.md)\n```\n\n" +
				"Another [link2](../c.md).\n",
			want: []Link{
				{Target: "a/b.md", Line: 3, Column: 11},
				{Target: "../c.md", Line: 9, Column: 9},
			},
		},
		{
			name: "empty target ignored",
			in:   `[empty]()`,
			want: nil,
		},
		{
			name: "nested brackets in text",
			in:   `[outer [inner] text](./x.md)`,
			want: []Link{{Target: "./x.md", Line: 1, Column: 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractMarkdownLinks([]byte(tt.in))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractMarkdownLinks() mismatch\n  got:  %#v\n  want: %#v", got, tt.want)
			}
		})
	}
}

func TestShouldKeepLinkTarget(t *testing.T) {
	keep := []string{
		"foo.md",
		"./foo.md",
		"../foo.md",
		"sub/dir/foo.md",
		"a(b).md",
	}
	drop := []string{
		"",
		"#anchor",
		"http://example.com",
		"https://example.com/path",
		"ftp://example.com",
		"mailto:foo@example.com",
		"tel:+12345",
	}
	for _, k := range keep {
		if !shouldKeepLinkTarget(k) {
			t.Errorf("expected to keep %q", k)
		}
	}
	for _, d := range drop {
		if shouldKeepLinkTarget(d) {
			t.Errorf("expected to drop %q", d)
		}
	}
}
