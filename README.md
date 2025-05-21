# RISC-GoV Core Library Documentation

## Overview

[RISC-GoV Core](https://github.com/RISC-GoV/core) is a lightweight, embeddable RISC-V CPU emulator written in Go. It provides a simple API for loading and executing RISC-V programs. It is ideal for building tools such as IDEs, debuggers, educational platforms, and embedded emulators for RISC-V instruction sets.

This documentation explains:

* What RISC-GoV Core is
* How to use it
* Integration in a sample Qt-based IDE
* Available CPU states and memory management

---

## Features

* Implements a basic RISC-V CPU emulator
* Single-instruction step execution
* Full memory emulation support
* Program loading from assembled `.exe` binaries
* Clean API for embedding and debugging

---

## Installation

```bash
go get github.com/RISC-GoV/core
```

---

## Basic Usage

```go
import rcore "github.com/RISC-GoV/core"

func main() {
    cpu := rcore.NewCPU(rcore.NewMemory())
    err := cpu.LoadFile("path/to/output.exe")
    if err != nil {
        log.Fatalf("Error loading program: %v", err)
    }

    for {
        state := cpu.ExecuteSingle()
        if state == rcore.PROGRAM_EXIT || state == rcore.PROGRAM_EXIT_FAILURE {
            break
        }
    }
}
```
OR ExecuteFile directly to run fully the code from A to Z
```go
debugInfo.cpu = rcore.NewCPU(rcore.NewMemory())
rcore.Kernel.Init()
err := debugInfo.cpu.ExecuteFile("output.exe")
```
OR use integrated CLI via DebugFile
```go
debugInfo.cpu = rcore.NewCPU(rcore.NewMemory())
rcore.Kernel.Init()
err := debugInfo.cpu.DebugFile("output.exe")
```

### CPU States

RISC-GoV returns a `State` enum after each instruction execution:

* `rcore.PROGRAM_CONTINUE` - Execution can proceed.
* `rcore.PROGRAM_EXIT` - Program exited normally.
* `rcore.PROGRAM_EXIT_FAILURE` - Program exited with an error.
* `rcore.E_BREAK` - An `ebreak` instruction was encountered (e.g., breakpoint).

---

## Memory Handling

Memory is managed via the `Memory` type:

```go
mem := rcore.NewMemory()
cpu := rcore.NewCPU(mem)
```

All memory accesses (load/store) go through this emulated memory object.

---

## Example Integration with GUI/IDE

In the [example IDE](https://github.com/RISC-GoV/gui) integration (written with `therecipe/qt`), the following features are demonstrated:

* Assembling source files with a custom assembler (`risc-assembler`)
* Inserting `ebreak` instructions as breakpoints
* Highlighting the currently executing line
* Register display updates

### Starting a Session

```go
debugInfo.cpu = rcore.NewCPU(rcore.NewMemory())
rcore.Kernel.Init()
err := debugInfo.cpu.LoadFile("output.exe")
```

### Executing Instructions

```go
state := debugInfo.cpu.ExecuteSingle()
```

### Step and Continue Functions

The IDE provides these:

* `stepDebugCode()` - Executes one instruction and updates the UI
* `continueDebugCode()` - Continuously runs until `ebreak`, exit, or failure
* `hotReloadCode()` - allows to replace memory while code is running

### Breakpoint Handling

During preprocessing, lines with breakpoints are prefixed with `ebreak` before being assembled:

```go
if debugInfo.breakpoints[lineIndex] {
    modifiedContent.WriteString("ebreak\n")
}
```

---

## API Reference

### `func NewCPU(mem *Memory) *CPU`

Creates a new CPU instance with attached memory.

### `func (c *CPU) LoadFile(path string) error`

Loads a compiled `.exe` binary into memory.

### `func (c *CPU) ExecuteSingle() State`

Executes one instruction and returns the new CPU state.

### `func NewMemory() *Memory`

Creates a new emulated memory object.

### `var Kernel kernelStruct`

A global kernel object to initialize system call handling.

---

## Error Handling

Typical error sources:

* Invalid program files
* Unsupported instructions (incomplete ISA implementation)

Use Go error checks to handle file loading and assembler integration issues.

---

## Contributing

The RISC-GoV project is open to contributions. Visit the [GitHub repository](https://github.com/RISC-GoV/core) and open issues or pull requests for fixes or new features.

---

## Related Projects

* [risc-assembler](https://github.com/RISC-GoV/risc-assembler): Assembler for generating `.exe` files
* [risc-IDE](https://github.com/RISC-GoV/gui) using Qt (`therecipe/qt`) with integration for assembling, debugging, and highlighting

---

## License

[UnLicensed](https://github.com/RISC-GoV/core/blob/master/LICENSE)
