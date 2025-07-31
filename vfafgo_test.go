package vfafgo

import (
    "bytes"
    "os"
    "strings"
    "testing"
)

func TestVfsBasic(t *testing.T) {
	tmp := t.TempDir()
	fs := NewFS(tmp)

	// Test file doesn't exist
	if fs.Exists("nofile.txt") {
		t.Error("Exists should be false for missing file")
	}

	// Write and Exists
	if err := fs.Write("hello", "test.txt"); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if !fs.Exists("test.txt") {
		t.Error("Exists should be true after write")
	}

	// Read
	data, err := fs.Read("test.txt")
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if data != "hello" {
		t.Errorf("Read wrong content: got %q", data)
	}

	// Relative paths
	subdir := []string{"foo", "bar.txt"}
	txt := "subfile"
	if err := fs.Write(txt, subdir...); err != nil {
		t.Fatal(err)
	}
	got, err := fs.Read(subdir...)
	if err != nil || got != txt {
		t.Errorf("subdir Read failed: got %q, err=%v", got, err)
	}
}

func TestVcrBasic(t *testing.T) {
	tmp := t.TempDir()

	// Write a file using Vfs
	fs := NewFS(tmp)
	if err := fs.Write("init", "f1.txt"); err != nil {
		t.Fatal(err)
	}

	// Init Vcr
	vcr := NewVcr("tester", tmp)

	// Change f1.txt, then commit changes
	if err := fs.Write("first", "f1.txt"); err != nil {
		t.Fatal(err)
	}
	commit1 := vcr.Commit("initial commit")
	if commit1 == "" || commit1 == "?" {
		t.Error("Commit should return hash")
	}

	// Change and commit again
	if err := fs.Write("second", "f1.txt"); err != nil {
		t.Fatal(err)
	}
	commit2 := vcr.Commit("second commit")
	if commit2 == "" || commit2 == "?" {
		t.Error("Second commit should return hash")
	}
	if commit1 == commit2 {
		t.Error("Commit hashes should be different")
	}

	// Get history
	hist := vcr.GetHistory()
	joined := strings.Join(hist, ",")
	if !strings.Contains(joined, "initial commit") || !strings.Contains(joined, "second commit") {
		t.Errorf("Unexpected history: %v", hist)
	}

	// Purge should remove .git
	if err := vcr.Purge(); err != nil {
		t.Errorf("Purge error: %v", err)
	}
	if _, err := os.Stat(fs.Path(".git")); !os.IsNotExist(err) {
		t.Error(".git folder not removed after Purge")
	}
}


func TestConfig(t *testing.T) {
	// Test with STORAGE_ROOT unset
	os.Unsetenv("STORAGE_ROOT")
	c := LoadStorageConfig()
	if c.StorageRoot == "" {
		t.Error("StorageRoot should be set to default if env is unset")
	}

	tmp := t.TempDir()
	_ = os.MkdirAll(tmp, 0755)
	r := New(tmp, "http://test-remote")
	if r.LocalDir == "" || r.RemoteURL == "" {
		t.Error("New should set LocalDir and RemoteURL")
	}
}

func TestIgnore(t *testing.T) {
	tmp := t.TempDir()
	def := ".qignore"
	ignoreFile := tmp + string(os.PathSeparator) + def
	os.WriteFile(ignoreFile, []byte("*.tmp\nignoreme.txt\n# comment\n\n"), 0644)
	m := LoadIgnoreMatcher(tmp)

	shouldIgnore := []string{"foo.tmp", "ignoreme.txt"}
	shouldKeep := []string{"not.tmp.jpg", "keepme.txt"}
	for _, f := range shouldIgnore {
		if !m.Ignore(f) {
			t.Errorf("Should ignore %v", f)
		}
	}
	for _, f := range shouldKeep {
		if m.Ignore(f) {
			t.Errorf("Should not ignore %v", f)
		}
	}
}

func TestHTTPHelpers(t *testing.T) {
	// This will generate an error because no HTTP server exists - that's normal for test.
	var b bytes.Buffer
	b.WriteString("test data")
	err := PutStream("http://localhost:9999/notexist", "", b)
	if err == nil {
		t.Error("PutStream expected error with dummy host")
	}

	// test encrypt function
	key := []byte("0123456789abcdef")
	enc, err := encrypt(b, key)
	if err != nil || enc == nil {
		t.Error("encrypt should not fail with sample key")
	}
}

func TestZipUnzip(t *testing.T) {
	tmp := t.TempDir()
	srcdir := tmp + "/src"
	os.Mkdir(srcdir, 0755)
	err := os.WriteFile(srcdir+"/a.txt", []byte("file a contents"), 0644)
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	os.WriteFile(srcdir+"/b.txt", []byte("file b!"), 0644)

	destzip := tmp + "/out.zip"
	b, err := Zip(tmp, "src", destzip)
	if err != nil {
		t.Fatalf("Zip failed: %v", err)
	}
	err = os.WriteFile(destzip, b.Bytes(), 0644)
	if err != nil {
		t.Fatalf("Failed to write zip: %v", err)
	}
	unzipdir := tmp + "/unzipped"
	os.Mkdir(unzipdir, 0755)
	err = Unzip(destzip, unzipdir)
	if err != nil {
		t.Fatalf("Unzip failed: %v", err)
	}
	f1, err := os.ReadFile(unzipdir + "/src/a.txt")
	if err != nil || string(f1) != "file a contents" {
		t.Errorf("Unzipped file wrong or missing: %s", string(f1))
	}
}
