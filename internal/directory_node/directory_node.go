package directorynode

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type DirectoryNode struct {
	Files       []string
	SubDirNodes map[string]*DirectoryNode
}

func New() *DirectoryNode {
	return &DirectoryNode{
		Files:       []string{},
		SubDirNodes: map[string]*DirectoryNode{},
	}
}

func (dn *DirectoryNode) Scan(start string) (*DirectoryNode, error) {
	root := dn

	// walk every directories found from `start` path in the tree
	err := filepath.Walk(start, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(start, path)
		if rel == "." {
			return nil
		}

		// split the current file name into parts, so it's easier to get the directory name
		parent := root
		parts := strings.Split(rel, string(os.PathSeparator))

		for i := 0; i < len(parts)-1; i++ {
			dir := parts[i]
			if _, exists := parent.SubDirNodes[dir]; !exists {
				parent.SubDirNodes[dir] = New()
			}
			// update the parent pointer to the current directory
			parent = parent.SubDirNodes[dir]
		}

		if info.IsDir() {
			// initialize dir node if scanned file is a directory
			parent.SubDirNodes[info.Name()] = New()
		} else {
			// append current file relative path to files array
			parent.Files = append(parent.Files, rel)
		}

		return nil
	})

	return root, err
}
