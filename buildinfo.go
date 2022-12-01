package buildinfo

import (
	"reflect"
	"runtime/debug"
)

var (
	ok      bool
	info    *debug.BuildInfo
	version = "dev"
)

func init() {
	info, ok = debug.ReadBuildInfo()
}

// Runtime returns the current Go runtime version
func Runtime() string {
	return info.GoVersion
}

// Path returns the import path
func Path() string {
	return info.Path
}

// Inspect whatever type and get module information in return
func Inspect(t any) *debug.Module {
	if info == nil {
		return nil
	}

	path := reflect.TypeOf(t).PkgPath()

	for _, dep := range info.Deps {
		if path == dep.Path {
			return dep
		}
	}

	return &debug.Module{
		Path:    path,
		Version: version,
	}
}

// Version returns the vcs.revision from buildinfo or injected version if set at
// built time. A version can be injected with the command below:
//
//	go build -ldflags "-X github.com/quartercastle/buildinfo.version="
func Version() string {
	if info == nil {
		return version
	}

	if version != "dev" {
		return version
	}

	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			return setting.Value[:7]
		}
	}

	return version
}
