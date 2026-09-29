package tree

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestRender(t *testing.T) {
	root := Build([]string{
		"go.mod",
		"cmd/rgit/main.go",
		"internal/ui/ui.go",
		"internal/git/git.go",
		"README.md",
		"empty/",
	})
	want := `.
├── cmd/
│   └── rgit/
│       └── main.go
├── empty/
├── internal/
│   ├── git/
│   │   └── git.go
│   └── ui/
│       └── ui.go
├── go.mod
└── README.md
`
	if got := Render(root, 0); got != want {
		t.Errorf("Render mismatch.\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestRenderDepth(t *testing.T) {
	root := Build([]string{"a/b/c.txt", "top.txt"})
	want := ".\n├── a/\n└── top.txt\n"
	if got := Render(root, 1); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestWalkSkipsNoise(t *testing.T) {
	dir := t.TempDir()
	for _, f := range []string{"main.go", "src/app.go", "node_modules/x/index.js", ".git/HEAD", "debug.log", ".env"} {
		path := filepath.Join(dir, filepath.FromSlash(f))
		os.MkdirAll(filepath.Dir(path), 0o755)
		os.WriteFile(path, nil, 0o644)
	}
	got, err := Walk(dir, IsNoise)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(got)
	want := []string{"main.go", "src/", "src/app.go"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Walk = %v, want %v", got, want)
	}
}
