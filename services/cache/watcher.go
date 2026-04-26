package cache

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// StartWatching initializes a new file system watcher for the image cache root directory.
// It adds all subdirectories (categories) within the root directory to the watcher,
// enabling monitoring for file system events. The method starts a background goroutine
// to handle watch events. Returns an error if the watcher cannot be created or if
// reading the root directory fails.
func (c *ImageCache) StartWatching() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	c.watcher = watcher

	entries, err := os.ReadDir(c.root)
	if err != nil {
		return err
	}

	if err := watcher.Add(c.root); err != nil {
		_ = watcher.Close()
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			categoryPath := filepath.Join(c.root, entry.Name())
			if err := watcher.Add(categoryPath); err != nil {
				_ = watcher.Close()
				return err
			}
		}
	}

	go c.watchLoop()
	return nil
}

// watchLoop continuously listens for filesystem events and errors from the cache's watcher.
// It processes events such as file creation, removal, and renaming within the cache root directory,
// updating the cache for the affected category as needed. If a category directory is detected as missing,
// it attempts to re-add it to the watcher. Any watcher errors are logged. The loop exits when the watcher
// channels are closed.
func (c *ImageCache) watchLoop() {
	for {
		select {
		case event, ok := <-c.watcher.Events:
			if !ok {
				return
			}
			c.handleWatchEvent(event)
		case err, ok := <-c.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Cache watcher error: %v", err)
		}
	}
}

func (c *ImageCache) handleWatchEvent(event fsnotify.Event) {
	if event.Op&(fsnotify.Create|fsnotify.Remove|fsnotify.Rename|fsnotify.Write) == 0 {
		return
	}

	relPath, err := filepath.Rel(c.root, event.Name)
	if err != nil || relPath == "." {
		return
	}

	parts := strings.Split(relPath, string(os.PathSeparator))
	if len(parts) < 1 || parts[0] == "" {
		return
	}

	category := parts[0]
	fullCategoryPath := filepath.Join(c.root, category)

	info, statErr := os.Stat(fullCategoryPath)
	switch {
	case statErr == nil && info.IsDir():
		if event.Op&fsnotify.Create != 0 && len(parts) == 1 {
			if err := c.watcher.Add(fullCategoryPath); err != nil {
				log.Printf("failed to watch new category '%s': %v", category, err)
			}
		}
		if err := c.LoadCategory(category); err != nil {
			log.Printf("failed to reload category '%s': %v", category, err)
			return
		}
		log.Printf("Cache watcher reloaded category '%s' after %s", category, event.Name)
	case os.IsNotExist(statErr):
		c.Lock()
		delete(c.files, category)
		delete(c.paths, category)
		delete(c.metas, category)
		c.Unlock()
		log.Printf("Cache watcher removed category '%s'", category)
	default:
		if statErr != nil {
			log.Printf("cache watcher stat error for '%s': %v", category, statErr)
		}
	}
}
