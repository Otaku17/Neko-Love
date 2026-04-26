package cache

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type FileMeta struct {
	Size       int64
	Readable   string
	MimeType   string
	ModifiedAt int64
}

type ImageCache struct {
	sync.RWMutex
	files   map[string][]string
	paths   map[string]map[string]string
	metas   map[string]map[string]FileMeta
	root    string
	watcher *fsnotify.Watcher
}

// New creates and initializes a new ImageCache instance using the specified assetsRoot directory.
// It loads all image categories into the cache. If loading fails, it returns an error.
// On success, it returns a pointer to the initialized ImageCache.
func New(assetsRoot string) (*ImageCache, error) {
	cache := &ImageCache{
		files: make(map[string][]string),
		paths: make(map[string]map[string]string),
		metas: make(map[string]map[string]FileMeta),
		root:  assetsRoot,
	}

	if err := cache.loadAllCategories(); err != nil {
		return nil, err
	}

	log.Printf("Image cache initialized with root: %s", assetsRoot)
	return cache, nil
}

// loadAllCategories loads all category directories from the cache root directory into the ImageCache.
// It iterates over the entries in the root directory, loading each category found.
// After loading, it logs the number of categories loaded and starts a file watcher to monitor changes.
// Returns an error if reading the directory or starting the watcher fails.
func (c *ImageCache) loadAllCategories() error {
	entries, err := os.ReadDir(c.root)
	if err != nil {
		return err
	}

	var loadErrors []error
	for _, entry := range entries {
		if entry.IsDir() {
			if err := c.LoadCategory(entry.Name()); err != nil {
				loadErrors = append(loadErrors, fmt.Errorf("%s: %w", entry.Name(), err))
			}
		}
	}

	log.Printf("Loaded %d categories from cache", len(c.files))
	if err := c.StartWatching(); err != nil {
		log.Printf("Error starting file watcher: %v", err)
		return err
	}

	if len(loadErrors) > 0 {
		return errors.Join(loadErrors...)
	}

	return nil
}

// LoadCategory loads all files from the specified category directory into the cache.
// It scans the directory, collects file names and their metadata (size, human-readable size,
// MIME type, and modification time), and stores them in the cache's internal structures.
// If the directory cannot be read, an error is returned. Files that cannot be stat'ed are skipped.
// The function is thread-safe and logs the number of files loaded for the category.
func (c *ImageCache) LoadCategory(category string) error {
	c.Lock()
	defer c.Unlock()

	folder := filepath.Join(c.root, category)
	entries, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	var paths []string
	pathByFile := make(map[string]string, len(entries))
	metaByFile := make(map[string]FileMeta, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		paths = append(paths, name)

		fullPath := filepath.Join(folder, name)
		pathByFile[name] = fullPath
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}

		mime := "application/octet-stream"
		if f, err := os.Open(fullPath); err == nil {
			buf := make([]byte, 512)
			_, _ = f.Read(buf)
			mime = http.DetectContentType(buf)
			_ = f.Close()
		}

		metaByFile[name] = FileMeta{
			Size:       info.Size(),
			Readable:   humanFileSize(info.Size()),
			MimeType:   mime,
			ModifiedAt: info.ModTime().Unix(),
		}
	}

	c.files[category] = paths
	c.paths[category] = pathByFile
	c.metas[category] = metaByFile

	log.Printf("Loaded %d files for category '%s'", len(paths), category)
	return nil
}

// GetRandom returns a random file path from the cache for the specified category.
// If the category does not exist or contains no files, it returns an error (os.ErrNotExist).
// This method is safe for concurrent use.
func (c *ImageCache) GetRandom(category string) (string, error) {
	c.RLock()
	defer c.RUnlock()

	paths, ok := c.files[category]
	if !ok || len(paths) == 0 {
		return "", os.ErrNotExist
	}

	return paths[rand.Intn(len(paths))], nil
}

// GetFiles returns a slice of file paths associated with the specified category.
// If the category does not exist in the cache, it returns nil.
// This method is safe for concurrent use.
func (c *ImageCache) GetFiles(category string) []string {
	c.RLock()
	defer c.RUnlock()

	paths, ok := c.files[category]
	if !ok {
		return nil
	}

	return append([]string(nil), paths...)
}

// Categories returns the currently cached category names in stable order.
func (c *ImageCache) Categories() []string {
	c.RLock()
	defer c.RUnlock()

	categories := make([]string, 0, len(c.files))
	for category := range c.files {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	return categories
}

// GetImagePath returns the full file path of an image with the specified name within the given category.
// It acquires a read lock to ensure thread-safe access to the cache.
// If the image is found, it returns the full path and true; otherwise, it returns an empty string and false.
func (c *ImageCache) GetImagePath(category, name string) (string, bool) {
	c.RLock()
	defer c.RUnlock()

	pathByFile, ok := c.paths[category]
	if !ok {
		return "", false
	}

	fullPath, exists := pathByFile[name]
	return fullPath, exists
}

// GetImageMeta retrieves the metadata (FileMeta) for an image specified by its category and name.
// It returns the FileMeta and a boolean indicating whether the metadata was found in the cache.
// The method is safe for concurrent use.
func (c *ImageCache) GetImageMeta(category, name string) (FileMeta, bool) {
	c.RLock()
	defer c.RUnlock()

	if meta, ok := c.metas[category]; ok {
		if data, exists := meta[name]; exists {
			return data, true
		}
	}
	return FileMeta{}, false
}

// humanFileSize converts a file size given in bytes to a human-readable string
// using binary (base-1024) units. For example, 1536 bytes will be formatted as "1.50 KB".
// The function supports units up to exabytes (EB).
//
// Parameters:
//   - size: the file size in bytes.
//
// Returns:
//   - A string representing the human-readable file size.
func humanFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := unit, 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func (c *ImageCache) Close() error {
	c.Lock()
	defer c.Unlock()

	if c.watcher == nil {
		return nil
	}

	err := c.watcher.Close()
	c.watcher = nil
	return err
}
