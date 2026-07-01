package storage

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const (
	maxSize  = 100 * 1024 * 1024
	maxTotal = 2 * 1024 * 1024 * 1024
	maxAge   = 7 * 24 * time.Hour
)

func Save(dir string, r io.Reader, originalName string) (name string, size int64, err error) {
	tmp, err := os.CreateTemp("", "upload-*")
	if err != nil {
		return "", 0, err
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	hash := md5.New()
	tee := io.TeeReader(io.LimitReader(r, maxSize+1), hash)
	written, err := io.Copy(tmp, tee)
	if err != nil {
		return "", 0, err
	}
	if written > maxSize {
		return "", 0, fmt.Errorf("file too large, max 100MB")
	}

	ext := ""
	if dot := filepath.Ext(originalName); dot != "" {
		ext = dot
	}
	filename := fmt.Sprintf("%x%s", hash.Sum(nil), ext)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", 0, err
	}

	dst := filepath.Join(dir, filename)
	if _, err := tmp.Seek(0, 0); err != nil {
		return "", 0, err
	}
	data, _ := io.ReadAll(tmp)
	if err := os.WriteFile(dst, data, 0644); err != nil {
		return "", 0, err
	}

	return filename, written, nil
}

func Evict(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	type fileInfo struct {
		name  string
		size  int64
		mtime time.Time
	}

	var files []fileInfo
	for _, e := range entries {
		if !e.Type().IsRegular() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		files = append(files, fileInfo{
			name:  e.Name(),
			size:  info.Size(),
			mtime: info.ModTime(),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].mtime.Before(files[j].mtime)
	})

	now := time.Now()

	for _, f := range files {
		if now.Sub(f.mtime) > maxAge {
			os.Remove(filepath.Join(dir, f.name))
		}
	}

	var remaining []fileInfo
	for _, f := range files {
		if now.Sub(f.mtime) <= maxAge {
			remaining = append(remaining, f)
		}
	}

	var total int64
	for _, f := range remaining {
		total += f.size
	}

	for _, f := range remaining {
		if total <= maxTotal {
			break
		}
		os.Remove(filepath.Join(dir, f.name))
		total -= f.size
	}

	return nil
}
