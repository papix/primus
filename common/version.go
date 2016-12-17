package common

import (
	"fmt"
	"runtime"
)

const (
	Version = "0.0.1"
)

func VersionInfo() []byte {
	return []byte(fmt.Sprintf(`primus %s
Compiler: %s %s
Copyright (C) 2016 papix <mail@papix.net>
`,
		Version,
		runtime.Compiler,
		runtime.Version()))
}
