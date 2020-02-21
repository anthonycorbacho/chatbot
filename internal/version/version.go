package version

import "fmt"

var (
	version      string
	gitCommit    string
	gitTreeState string
	buildDate    string
)

// Info represents the structure of the version information.
type Info struct {
	Version      string `json:"version"`
	GitCommit    string `json:"git_commit"`
	GitTreeState string `json:"git_tree_state"`
	BuildDate    string `json:"build_date"`
}

// Get returns the version and buildtime information about the binary
func Get() *Info {
	// These variables typically come from -ldflags settings to `go build`, See Makefile
	return &Info{
		Version:      version,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
	}
}

// String returns a version as a string
func (i *Info) String() string {
	return fmt.Sprintf("Version { version: %v, gitCommit: %v, gitTree: %v, buildDate: %v }",
		i.Version,
		i.GitCommit,
		i.GitTreeState,
		i.BuildDate)
}
