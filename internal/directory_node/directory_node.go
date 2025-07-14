package directorynode

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"slices"
)

type DirectoryNode struct {
	Files       []string
	SubDirNodes map[string]*DirectoryNode
}

// TODO: idk if all theese extensions is gonna work with what am going to be working next
var extensions = []string{
	// Unknown
	".webm",

	// SDTV
	".m4v", ".3gp", ".nsv", ".ty", ".strm", ".rm", ".rmvb", ".m3u", ".ifo",
	".mov", ".qt", ".divx", ".xvid", ".bivx", ".nrg", ".pva", ".wmv", ".asf",
	".asx", ".ogm", ".ogv", ".m2v", ".avi", ".bin", ".dat", ".dvr-ms", ".mpg",
	".mpeg", ".mp4", ".avc", ".vp3", ".svq3", ".nuv", ".viv", ".dv", ".fli",
	".flv", ".wpl",

	// DVD
	".img", ".iso", ".vob",

	// HD
	".mkv", ".mk3d", ".ts", ".wtv",

	// Bluray
	".m2ts",
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

		parent := root

		if parent.SubDirNodes == nil {
			parent.SubDirNodes = make(map[string]*DirectoryNode)
		}

		if info.IsDir() {
			if _, exists := parent.SubDirNodes[info.Name()]; !exists {
				parent.SubDirNodes[info.Name()] = New()
			}
			parent.SubDirNodes[info.Name()] = New()
		} else {
			// append current file relative path to files array
			parent.Files = append(parent.Files, filepath.Join(start, rel))
		}

		return nil
	})

	return root, err
}

// get all files on every directory level
func (dn *DirectoryNode) DirFiles() []string {
	res := []string{}
	stack := []*DirectoryNode{dn}
	for len(stack) > 0 {
		current := stack[0]
		stack = stack[1:]
		if len(current.Files) > 0 {
			for _, file := range current.Files {
				if dn.isMedia(file) {
					res = append(res, file)
				}
			}
		}
		for _, node := range current.SubDirNodes {
			stack = append(stack, node)
		}
	}
	return res
}

func (dn *DirectoryNode) isMedia(filename string) bool {
	pattern := regexp.MustCompile(`\.[0-9a-zA-Z]+$`)
	match := pattern.FindString(filename)
	if slices.Contains(extensions, match) {
		return true
	}
	return false
}
