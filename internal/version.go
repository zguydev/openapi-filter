package internal

import (
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
)

const pkg = "github.com/zguydev/openapi-filter"

type VersionInfo struct {
	Version   string
	GoVersion string
	Commit    string
}

func (vi VersionInfo) String() string {
	var s strings.Builder
	s.WriteString("openapi-filter version ")
	if v := vi.Version; v != "" {
		s.WriteString(v)
	} else {
		s.WriteString("unknown")
	}
	if commit := vi.Commit; commit != "" {
		s.WriteByte('-')
		s.WriteString(commit)
	}
	if v := vi.GoVersion; v != "" {
		s.WriteByte(' ')
		if v != "" {
			s.WriteString(v)
		}
	}
	const osArch = " " + runtime.GOOS + "/" + runtime.GOARCH
	s.WriteString(osArch)
	return s.String()
}

var getOnce struct {
	info VersionInfo
	ok   bool
	once sync.Once
}

func getToolVersion(m *debug.Module) (string, bool) {
	if m == nil || m.Path != pkg {
		return "", false
	}
	return m.Version, true
}

func getVersionInfo() (info VersionInfo, ok bool) {
	getOnce.once.Do(func() {
		bi, ok := debug.ReadBuildInfo()
		getOnce.ok = ok
		if !ok {
			return
		}

		var isDep bool
		if v, ok := getToolVersion(&bi.Main); ok {
			getOnce.info.Version = v
		} else {
			isDep = true
			for _, m := range bi.Deps {
				if v, ok := getToolVersion(m); ok {
					getOnce.info.Version = v
					break
				}
			}
		}
		getOnce.info.GoVersion = bi.GoVersion
		if !isDep {
			for _, s := range bi.Settings {
				switch s.Key {
				case "vcs.revision":
					getOnce.info.Commit = s.Value
				}
			}
		}
	})
	return getOnce.info, getOnce.ok
}

func GetInfo() (info VersionInfo, ok bool) {
	return getVersionInfo()
}
