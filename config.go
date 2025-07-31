package vfafgo

import (
	"log"
	"os"
	"path/filepath"
)

const storage_root_env_key = "STORAGE_ROOT"
const storage_root_default = "/data/storage"

type Config struct {
	LocalDir    string
	RemoteURL   string
	StorageRoot string
}

func LoadStorageConfig() Config {
	root := os.Getenv(storage_root_env_key)
	if root == "" {
		root = storage_root_default
	}

	return Config{
		StorageRoot: root,
	}
}

func New(localDir, remoteURL string) Config {
	absPath, err := filepath.Abs(localDir)
	if err != nil {
		log.Fatalf("invalid path: %v", err)
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Fatalf("directory does not exist: %s", absPath)
	}
	return Config{
		LocalDir:  absPath,
		RemoteURL: remoteURL,
	}
}
