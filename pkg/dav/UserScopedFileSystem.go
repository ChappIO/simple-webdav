package dav

import (
	"context"
	"golang.org/x/net/webdav"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type UserScopedFileSystem struct {
	FileSystem webdav.FileSystem
	RootDir    string
	SubDir string
}

func slashClean(name string) string {
	if name == "" || name[0] != '/' {
		name = "/" + name
	}
	return path.Clean(name)
}

func (fs *UserScopedFileSystem) getUserPath(ctx context.Context, name string) string {
	// This implementation is based on Dir.Open's code in the standard net/http package.
	if filepath.Separator != '/' && strings.IndexRune(name, filepath.Separator) >= 0 ||
		strings.Contains(name, "\x00") {
		return ""
	}
	dir := fs.RootDir
	if dir == "" {
		dir = "."
	}
	return filepath.Join(dir, ctx.Value("Username").(string), "Files", filepath.FromSlash(slashClean(name)))
}

func (fs *UserScopedFileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	if name = fs.getUserPath(ctx, name); name == "" {
		return os.ErrNotExist
	}
	return os.MkdirAll(name, perm)
}
func (fs *UserScopedFileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	if name = fs.getUserPath(ctx, name); name == "" {
		return nil, os.ErrNotExist
	}
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}
func (fs *UserScopedFileSystem) RemoveAll(ctx context.Context, name string) error {
	if name = fs.getUserPath(ctx, name); name == "" {
		return os.ErrNotExist
	}
	if name == filepath.Clean(fs.getUserPath(ctx, ".")) {
		// Prohibit removing the virtual root directory.
		return os.ErrInvalid
	}
	return os.RemoveAll(name)
}
func (fs *UserScopedFileSystem) Rename(ctx context.Context, oldName, newName string) error {
	if oldName = fs.getUserPath(ctx, oldName); oldName == "" {
		return os.ErrNotExist
	}
	if newName = fs.getUserPath(ctx, newName); newName == "" {
		return os.ErrNotExist
	}
	if root := filepath.Clean(fs.getUserPath(ctx, ".")); root == oldName || root == newName {
		// Prohibit renaming from or to the virtual root directory.
		return os.ErrInvalid
	}
	return os.Rename(oldName, newName)
}
func (fs *UserScopedFileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	if name = fs.getUserPath(ctx, name); name == "" {
		return nil, os.ErrNotExist
	}
	return os.Stat(name)
}
