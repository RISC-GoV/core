package main

import (
	"RiscCPU/core"
	"fmt"
	"os"
	"path"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename> [-debug]")
		return
	}

	cpu := core.NewCPU(core.NewMemory())
	core.Kernel.CWD = path.Dir(os.Args[1])
	if os.Args[2] == "-debug" {
		cpu.DebugFile(os.Args[1])
	} else {
		cpu.LoadFile(os.Args[1])
	}
}
