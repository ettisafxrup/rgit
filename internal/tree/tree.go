// Package tree renders a list of file paths as a directory tree:
//
//	.
//	├── cmd/
//	│   └── main.go
//	└── go.mod
package tree

import (
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// Node is a file or directory in the tree.
type Node struct {
	Name     string
	IsDir    bool
	Children map[string]*Node
}

// Build creates a tree from slash-separated relative paths. A trailing
// slash marks an (empty) directory.
func Build(paths []string) *Node {
	root := newDir(".")
	for _, p := range paths {
		isDir := strings.HasSuffix(p, "/")
		parts := strings.Split(strings.Trim(path.Clean(filepath.ToSlash(p)), "/"), "/")
		node := root
		for i, part := range parts {
			if part == "" || part == "." {
				continue
			}
			last := i == len(parts)-1
			child, ok := node.Children[part]
			if !ok {
				if last && !isDir {
					child = &Node{Name: part}
				} else {
					child = newDir(part)
				}
				node.Children[part] = child
			}
			node = child
		}
	}
	return root
}

func newDir(name string) *Node {
	return &Node{Name: name, IsDir: true, Children: map[string]*Node{}}
}

// Render draws the tree. maxDepth limits how deep it goes; 0 means no limit.
func Render(root *Node, maxDepth int) string {
	var b strings.Builder
	b.WriteString(".\n")
	renderChildren(&b, root, "", 1, maxDepth)
	return b.String()
}

func renderChildren(b *strings.Builder, node *Node, prefix string, depth, maxDepth int) {
	if maxDepth > 0 && depth > maxDepth {
		return
	}
	children := sortedChildren(node)
	for i, child := range children {
		connector, nextPrefix := "├── ", prefix+"│   "
		if i == len(children)-1 {
			connector, nextPrefix = "└── ", prefix+"    "
		}
		name := child.Name
		if child.IsDir {
			name += "/"
		}
		b.WriteString(prefix + connector + name + "\n")
		if child.IsDir {
			renderChildren(b, child, nextPrefix, depth+1, maxDepth)
		}
	}
}

// sortedChildren orders directories before files, then alphabetically.
func sortedChildren(node *Node) []*Node {
	children := make([]*Node, 0, len(node.Children))
	for _, child := range node.Children {
		children = append(children, child)
	}
	sort.Slice(children, func(i, j int) bool {
		a, b := children[i], children[j]
		if a.IsDir != b.IsDir {
			return a.IsDir
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})
	return children
}

// Walk lists the files below root, skipping entries for which skip returns
// true. Paths are relative, slash-separated, and empty directories end in "/".
func Walk(root string, skip func(name string, isDir bool) bool) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // unreadable entries are simply left out
		}
		if p == root {
			return nil
		}
		if skip(d.Name(), d.IsDir()) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		rel, _ := filepath.Rel(root, p)
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			rel += "/"
		}
		paths = append(paths, rel)
		return nil
	})
	return paths, err
}

// noise lists directories and files that rarely belong in a project overview.
var noise = map[string]bool{
	"node_modules": true, "bower_components": true, "vendor": true,
	"dist": true, "build": true, "bin": true, "obj": true, "out": true, "target": true,
	"coverage": true, "logs": true, "env": true, "venv": true, "__pycache__": true,
	"Thumbs.db": true, "desktop.ini": true,
}

var noiseExtensions = []string{".log", ".tmp", ".swp", ".swo", ".bak", ".pyc"}

// IsNoise reports whether a file or directory is usually clutter: hidden
// entries, dependency folders, build output, logs and editor leftovers.
func IsNoise(name string, isDir bool) bool {
	if strings.HasPrefix(name, ".") || noise[name] {
		return true
	}
	if isDir {
		return false
	}
	for _, ext := range noiseExtensions {
		if strings.HasSuffix(strings.ToLower(name), ext) {
			return true
		}
	}
	return false
}
