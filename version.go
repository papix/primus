package primus

import (
	"fmt"
	"runtime"
)

func PrintVersion() {
	fmt.Printf(`primus %s
Compiler: %s %s
Copyright (C) 2016 papix <mail@papix.net>
`,
		Version,
		runtime.Compiler,
		runtime.Version())
}
