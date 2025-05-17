# core

## example usage
`
go get github.com/RISC-GoV/core
`
if the latest commit was not taken then go on branch dirty,
take the latest commit hash then
`
go get github.com/RISC-GoV/core@[HASH]
`

```go
package main

import (
	"github.com/RISC-GoV/core"
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
```