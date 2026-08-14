package copy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopyFilesToDirectoryFlattensAndRenamesDuplicates(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	mustWrite(t, filepath.Join(src, "root.txt"), "root")
	mustWrite(t, filepath.Join(src, "sub", "nested.txt"), "nested")
	mustWrite(t, filepath.Join(src, "a", "dup.txt"), "aaa")
	mustWrite(t, filepath.Join(src, "b", "dup.txt"), "bbb")

	n, err := CopyFilesToDirectory(src, dst, FilterAll, nil)
	if err != nil {
		t.Fatalf("CopyFilesToDirectory: %v", err)
	}
	if n != 4 {
		t.Fatalf("copied %d files, want 4", n)
	}

	entries, err := os.ReadDir(dst)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 4 {
		t.Fatalf("dest has %d files, want 4", len(entries))
	}

	names := map[string]bool{}
	dupCount := 0
	for _, e := range entries {
		if e.IsDir() {
			t.Fatalf("dest should be flat, found dir %s", e.Name())
		}
		names[e.Name()] = true
		if e.Name() == "dup.txt" || (strings.HasSuffix(e.Name(), "-dup.txt") && e.Name() != "dup.txt") {
			dupCount++
		}
	}

	if !names["root.txt"] || !names["nested.txt"] || !names["dup.txt"] {
		t.Fatalf("dest names = %v, missing expected files", names)
	}
	if dupCount != 2 {
		t.Fatalf("expected 2 dup.txt variants, got %d in %v", dupCount, names)
	}
}

func TestCopyFilesToDirectoryByExtension(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	mustWrite(t, filepath.Join(src, "a.JPG"), "jpg")
	mustWrite(t, filepath.Join(src, "b.png"), "png")
	mustWrite(t, filepath.Join(src, "c.txt"), "txt")
	mustWrite(t, filepath.Join(src, "sub", "d.jpeg"), "jpeg")

	n, err := CopyFilesToDirectory(src, dst, FilterInclude, ParseExtensions(".jpg, png"))
	if err != nil {
		t.Fatalf("CopyFilesToDirectory: %v", err)
	}
	if n != 2 {
		t.Fatalf("copied %d files, want 2", n)
	}

	entries, err := os.ReadDir(dst)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for _, e := range entries {
		got[e.Name()] = true
	}
	if !got["a.JPG"] || !got["b.png"] || got["c.txt"] || got["d.jpeg"] {
		t.Fatalf("dest names = %v", got)
	}
}

func TestCopyFilesToDirectoryExcludeExtension(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	mustWrite(t, filepath.Join(src, "a.JPG"), "jpg")
	mustWrite(t, filepath.Join(src, "b.png"), "png")
	mustWrite(t, filepath.Join(src, "c.txt"), "txt")
	mustWrite(t, filepath.Join(src, "sub", "d.pdf"), "pdf")

	n, err := CopyFilesToDirectory(src, dst, FilterExclude, ParseExtensions("jpg, txt"))
	if err != nil {
		t.Fatalf("CopyFilesToDirectory: %v", err)
	}
	if n != 2 {
		t.Fatalf("copied %d files, want 2", n)
	}

	entries, err := os.ReadDir(dst)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for _, e := range entries {
		got[e.Name()] = true
	}
	if !got["b.png"] || !got["d.pdf"] || got["a.JPG"] || got["c.txt"] {
		t.Fatalf("dest names = %v", got)
	}
}

func TestCopyFilesToDirectoryMissingSource(t *testing.T) {
	n, err := CopyFilesToDirectory(filepath.Join(t.TempDir(), "missing"), t.TempDir(), FilterAll, nil)
	if err == nil || n != -1 {
		t.Fatalf("got (%d, %v), want (-1, error)", n, err)
	}
}

func TestParseExtensions(t *testing.T) {
	got := ParseExtensions(" .JPG,png; PDF  txt ")
	want := []string{"jpg", "png", "pdf", "txt"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
