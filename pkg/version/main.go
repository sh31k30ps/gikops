package version

var (
	// Version is the current version of gikopsctl
	Version = "0.1.0"

	// GitCommit is the git commit hash, this will be filled by the build system
	GitCommit = "unknown"

	// BuildTime is the build timestamp, this will be filled by the build system
	BuildTime = "unknown"

	// GoVersion is the go version used to build the binary
	GoVersion = "unknown"
)
