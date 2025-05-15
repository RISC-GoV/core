package main

import (
	"RiscCPU/core"
	"os"
	"path"
)

func main() {
	cpu := core.NewCPU(core.NewMemory())
	core.Kernel.CWD = path.Dir(os.Args[1])
	if len(os.Args) > 2 && os.Args[2] == "-debug" {
		cpu.DebugFile(os.Args[1])
	} else {
		cpu.ExecuteFile(os.Args[1])
	}
}
