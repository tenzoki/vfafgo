package gov

import (
    "os"
    "testing"
    "strings"
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
