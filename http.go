package vfafgo

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	crypto "github.com/tenzoki/cryptogo"
)

func encrypt(buf bytes.Buffer, key []byte) (*bytes.Buffer, error) {
	encrypted, err := crypto.Encrypt(buf.Bytes(), key)
	if err != nil {
		fmt.Printf("Encryption error: %v", err)
		return nil, err
	}
	return bytes.NewBuffer(encrypted), nil
}

func PutStreamEncrypted(remoteURL, rel string, buf bytes.Buffer, key []byte) error {
	encrypted, err := encrypt(buf, key)
	if err != nil {
		fmt.Printf("Encryption error: %v", err)
		return err
	}
	return PutStream(remoteURL, rel, *encrypted)
}

func PutStream(remoteURL, rel string, buf bytes.Buffer) error {

	target := strings.TrimSuffix(remoteURL, "/")
	if rel != "" {
		target += "/" + strings.TrimPrefix(rel, "/")
	}
	req, err := http.NewRequest(http.MethodPut, target, bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("remote error: %s", resp.Status)
	}
	fmt.Println("Push OK")
	return nil
}
