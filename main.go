package main

import "RiscCPU/core"

func main() {
	cpu := core.NewCPU(core.NewMemory())
	cpu.ExecuteFile("output.exe")
}
