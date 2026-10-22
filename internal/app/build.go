// Default build info. Values are overridden at link time by `make build`
// via -ldflags -X, so plain `go build` (e.g. in CI) still compiles.
package app

var (
	Version = "0.0.0"
	Git     = "not a git repo"
	GoVer   = "unknown"
	BuildAt = "unknown"
	Branch  = "unknown"
	Distro  = "Unknown"
	Kernel  = "Unknown"
)
