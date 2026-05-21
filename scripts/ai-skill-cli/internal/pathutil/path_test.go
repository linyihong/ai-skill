package pathutil

import "testing"

func TestNormalizeForReport(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "windows drive path",
			input: `C:\Workspace\Repo\file.txt`,
			want:  "C:/Workspace/Repo/file.txt",
		},
		{
			name:  "mixed separators collapse",
			input: `C:/Workspace//Repo\..\Repo\file.txt`,
			want:  "C:/Workspace/Repo/file.txt",
		},
		{
			name:  "unc path",
			input: `\\server\share\folder\file.txt`,
			want:  "//server/share/folder/file.txt",
		},
		{
			name:  "spaces preserved",
			input: `C:\Workspace\My Documents\Ai-skill`,
			want:  "C:/Workspace/My Documents/Ai-skill",
		},
		{
			name:  "relative path",
			input: `..\repo\scripts\ai-skill-cli`,
			want:  "../repo/scripts/ai-skill-cli",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeForReport(tt.input)
			if err != nil {
				t.Fatalf("NormalizeForReport returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("NormalizeForReport(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeForReportRejectsEmptyPath(t *testing.T) {
	if _, err := NormalizeForReport("   "); err == nil {
		t.Fatal("expected empty path error")
	}
}

func TestSummarizePathList(t *testing.T) {
	summary := SummarizePathList("")
	if summary.Entries != 0 || summary.EmptyEntries != 0 {
		t.Fatalf("expected empty summary, got %#v", summary)
	}
}
