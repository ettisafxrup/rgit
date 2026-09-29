// Package version holds build information for rgit.
package version

// Version is the current rgit version. Release builds override it with:
//
//	go build -ldflags "-X github.com/ettisafxrup/rgit/internal/version.Version=2.0.0"
var Version = "2.0.0"
