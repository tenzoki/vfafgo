package gov

import (
	"os"
	"path/filepath"
)

type Vfs struct {
	root string
}

// NewFS initializes the file system abstraction with a root folder
func NewFS(root string) *Vfs {
	abs, err := filepath.Abs(root)
	if err != nil {
		panic("invalid root path")
	}
	return &Vfs{root: abs}
}

// Path returns the full path for a given relative path
func (fs *Vfs) Path(parts ...string) string {
	rel := filepath.Join(parts...)
	return filepath.Join(fs.root, rel)
}

// Read reads a file relative to the root
func (fs *Vfs) Read(parts ...string) (string, error) {
	data, err := os.ReadFile(fs.Path(parts...))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Write writes content to a file relative to the root
func (fs *Vfs) Write(content string, parts ...string) error {
	path := fs.Path(parts...)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// Exists checks if a file or dir exists
func (fs *Vfs) Exists(parts ...string) bool {
	_, err := os.Stat(fs.Path(parts...))
	return !os.IsNotExist(err)
}
