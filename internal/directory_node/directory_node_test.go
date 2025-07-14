package directorynode

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestDirectoryNode(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.MkdirAll(filepath.Join(tmpDir, "root/dir1/dir2/dir3"), 0755)
	if err != nil {
		t.Fatalf("failed to create dirs: %v", err)
	}

	files := map[string]string{
		fmt.Sprintf("%s", filepath.Join(tmpDir, "root", "file1.txt")):                         "This is file1 in root",
		fmt.Sprintf("%s", filepath.Join(tmpDir, "root", "dir1", "file2.txt")):                 "This is file2 in dir1",
		fmt.Sprintf("%s", filepath.Join(tmpDir, "root", "dir1", "dir2", "file3.log")):         "Log file in dir2",
		fmt.Sprintf("%s", filepath.Join(tmpDir, "root", "dir1", "dir2", "dir3", "file4.md")):  "Markdown file in dir3",
		fmt.Sprintf("%s", filepath.Join(tmpDir, "root", "dir1", "dir2", "video.mp4")):         "FAKE_MP4_CONTENT",
		fmt.Sprintf("%s", filepath.Join(tmpDir, "root", "dir1", "dir2", "audio.webm")):        "FAKE_WEBM_CONTENT",
		fmt.Sprintf("%s", filepath.Join(tmpDir, "root", "dir1", "dir2", "dir3", "clip.webm")): "FAKE_CLIP",
	}

	for p, content := range files {
		err := os.WriteFile(p, []byte(content), 0644)
		if err != nil {
			t.Fatalf("failed to write file %s: %v", p, err)
		}
	}

	rootNode := &DirectoryNode{}
	_, err = rootNode.Scan(filepath.Join(tmpDir, "root"))
	if err != nil {
		t.Errorf("Scan failed: %v", err)
	}

	for _, file := range rootNode.DirFiles() {
		fullPath := filepath.Join(file)
		_, exists := files[fullPath]
		if !exists {
			t.Errorf("expected %s got %s", file, fullPath)
		}
		if !rootNode.isMedia(fullPath) {
			t.Errorf("expected %s to be media", file)
		}
	}
}
