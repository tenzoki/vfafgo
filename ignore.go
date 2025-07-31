package vfafgo

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

const ignore_def_file = ".qignore"

type IgnoreMatcher struct {
	patterns []string
}

func LoadIgnoreMatcher(baseDir string) IgnoreMatcher {
	path := filepath.Join(baseDir, ignore_def_file)
	f, err := os.Open(path)
	if err != nil {
		return IgnoreMatcher{}
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var pats []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		pats = append(pats, line)
	}
	return IgnoreMatcher{pats}
}

func (m IgnoreMatcher) Ignore(relPath string) bool {
	for _, p := range m.patterns {
		matched, _ := filepath.Match(p, relPath)
		if matched {
			return true
		}
	}
	return false
}
