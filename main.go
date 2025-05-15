package main

import (
	"RiscCPU/core"
	"os"
)

func main() {
	cpu := core.NewCPU(core.NewMemory())
	path, _ := os.Getwd()
	core.Kernel.CWD = path //path.Dir(os.Args[1])
	cpu.ExecuteFile("./output.exe")
	// if len(os.Args) < 2 {
	// 	fmt.Println("Usage: go run main.go <filename> [-debug]")
	// 	return
	// }

	// if os.Args[2] == "-debug" {
	// 	cpu.DebugFile(os.Args[1])
	// } else {
	// 	cpu.ExecuteFile(os.Args[1])
	// }
}
