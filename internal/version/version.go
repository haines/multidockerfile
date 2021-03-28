package version

import "runtime"

var (
	version   = "unknown"
	gitCommit = "unknown"
	built     = "unknown"
)

type Info struct {
	Version   string
	GitCommit string
	Built     string
	GoVersion string
	OS        string
	Arch      string
}

func Get() Info {
	return Info{
		Version:   version,
		GitCommit: gitCommit,
		Built:     built,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}
