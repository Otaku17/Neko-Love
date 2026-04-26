package cache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fsnotify/fsnotify"
)

func TestLoadCategoryBuildsMetadata(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	categoryDir := filepath.Join(root, "neko")
	if err := os.Mkdir(categoryDir, 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	imagePath := filepath.Join(categoryDir, "cat.png")
	if err := os.WriteFile(imagePath, pngFixture(), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	cache := &ImageCache{
		files: make(map[string][]string),
		paths: make(map[string]map[string]string),
		metas: make(map[string]map[string]FileMeta),
		root:  root,
	}

	if err := cache.LoadCategory("neko"); err != nil {
		t.Fatalf("LoadCategory failed: %v", err)
	}

	meta, ok := cache.GetImageMeta("neko", "cat.png")
	if !ok {
		t.Fatal("expected metadata for cat.png")
	}

	if meta.MimeType != "image/png" {
		t.Fatalf("expected image/png, got %q", meta.MimeType)
	}

	files := cache.GetFiles("neko")
	if len(files) != 1 || files[0] != "cat.png" {
		t.Fatalf("unexpected files slice: %#v", files)
	}
}

func TestGetFilesReturnsCopy(t *testing.T) {
	t.Parallel()

	cache := &ImageCache{
		files: map[string][]string{
			"neko": {"a.png"},
		},
		paths: map[string]map[string]string{
			"neko": {"a.png": filepath.Join("assets", "neko", "a.png")},
		},
		metas: make(map[string]map[string]FileMeta),
	}

	files := cache.GetFiles("neko")
	files[0] = "mutated.png"

	again := cache.GetFiles("neko")
	if again[0] != "a.png" {
		t.Fatalf("expected cached slice to remain unchanged, got %#v", again)
	}
}

func TestHandleWatchEventLoadsNewCategory(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatalf("watcher failed: %v", err)
	}
	defer watcher.Close()

	cache := &ImageCache{
		files:   make(map[string][]string),
		paths:   make(map[string]map[string]string),
		metas:   make(map[string]map[string]FileMeta),
		root:    root,
		watcher: watcher,
	}

	categoryDir := filepath.Join(root, "hug")
	if err := os.Mkdir(categoryDir, 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(categoryDir, "01.png"), pngFixture(), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	cache.handleWatchEvent(fsnotify.Event{
		Name: filepath.Join(root, "hug"),
		Op:   fsnotify.Create,
	})

	files := cache.GetFiles("hug")
	if len(files) != 1 || files[0] != "01.png" {
		t.Fatalf("expected new category to be loaded, got %#v", files)
	}
}

func pngFixture() []byte {
	return []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
		0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
		0x89, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9c, 0x63, 0xf8, 0xcf, 0xc0, 0xf0,
		0x1f, 0x00, 0x05, 0x00, 0x01, 0xff, 0x89, 0x99,
		0x3d, 0x1d, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45,
		0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
	}
}
