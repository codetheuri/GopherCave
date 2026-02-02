package saver

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func SavePage(filename string, data io.Reader) error {
	// Create the directory path before creating the file
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	out, err := os.Create(filename)
	if err != nil { return err }
	defer out.Close()

	_, err = io.Copy(out, data)
	return err
}

func DownloadAsset(hostRoot, assetURL string) error {
	resp, err := http.Get(assetURL)
	if err != nil { return err }
	defer resp.Body.Close()

	name := filepath.Base(assetURL)
	assetDir := filepath.Join(hostRoot, "assets")
	os.MkdirAll(assetDir, 0755)

	content, _ := io.ReadAll(resp.Body)
	if strings.HasSuffix(name, ".css") {
		css := string(content)
		css = strings.ReplaceAll(css, "../fonts/", "")
		css = strings.ReplaceAll(css, "fonts/", "")
		content = []byte(css)
	}

	return os.WriteFile(filepath.Join(assetDir, name), content, 0644)
}