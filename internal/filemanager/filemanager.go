package filemanager

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileEntry struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	IsDir   bool   `json:"is_dir"`
	Size    int64  `json:"size"`
	Mode    string `json:"mode"`
	ModTime int64  `json:"mod_time"`
}

func safePath(path string) (string, error) {
	cleaned := filepath.Clean(path)
	if !filepath.IsAbs(cleaned) {
		return "", fmt.Errorf("path must be absolute")
	}
	// prevent directory traversal
	if strings.Contains(cleaned, "..") {
		return "", fmt.Errorf("invalid path")
	}
	return cleaned, nil
}

func List(dir string) ([]FileEntry, error) {
	dir, err := safePath(dir)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var result []FileEntry
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		result = append(result, FileEntry{
			Name:    entry.Name(),
			Path:    filepath.Join(dir, entry.Name()),
			IsDir:   entry.IsDir(),
			Size:    info.Size(),
			Mode:    info.Mode().String(),
			ModTime: info.ModTime().Unix(),
		})
	}
	return result, nil
}

func ReadFile(path string) (*os.File, os.FileInfo, error) {
	path, err := safePath(path)
	if err != nil {
		return nil, nil, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, nil, err
	}
	if info.IsDir() {
		return nil, nil, fmt.Errorf("cannot download a directory")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	return f, info, nil
}

func Upload(dir, filename string, src io.Reader) error {
	dir, err := safePath(dir)
	if err != nil {
		return err
	}

	destPath := filepath.Join(dir, filepath.Base(filename))
	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, src)
	return err
}

func Delete(path string) error {
	path, err := safePath(path)
	if err != nil {
		return err
	}
	return os.RemoveAll(path)
}

func Mkdir(path string) error {
	path, err := safePath(path)
	if err != nil {
		return err
	}
	return os.MkdirAll(path, 0755)
}

func Rename(oldPath, newPath string) error {
	old, err := safePath(oldPath)
	if err != nil {
		return err
	}
	newP, err := safePath(newPath)
	if err != nil {
		return err
	}
	return os.Rename(old, newP)
}
