package cache

import (
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkGetImagePath(b *testing.B) {
	root := b.TempDir()
	categoryDir := filepath.Join(root, "neko")
	if err := os.Mkdir(categoryDir, 0o755); err != nil {
		b.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(categoryDir, "cat.png"), pngFixture(), 0o644); err != nil {
		b.Fatalf("write file failed: %v", err)
	}

	imageCache, err := New(root)
	if err != nil {
		b.Fatalf("cache init failed: %v", err)
	}
	defer imageCache.Close()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, ok := imageCache.GetImagePath("neko", "cat.png"); !ok {
			b.Fatal("expected image path")
		}
	}
}
