# Go Generic Math Library

A drop‑in replacement for Go's standard `math` package with generic wrappers that eliminate manual type conversions. This library also includes additional user-defined functions and types to extend functionality.

## Features

- **Generic Functions:** Automatically wraps standard math functions to work with a variety of numeric types.
- **Automatic Code Generation:** Use `go:generate` to produce core files (`core_functions.go`, `core_consts.go`, `core_vars.go`, and `core_types.go`) that mirror the standard library.
- **Extended Utilities:** Extra functions such as `Blend`, `Avg`, `Clamp`, `Delta`, `Wrap`, `FormatNumber`, and more.
- **User Extensible:** Easily fork, modify, and extend the library to suit your needs.

## Installation

Clone the repository:

```sh
git clone https://github.com/toxyl/math.git
cd math
```

## Code Generation

The library uses `go:generate` to automatically generate core files. To generate these files, run:

```sh
go generate generate.go
```

This will create:
- `core_functions.go`
- `core_consts.go`
- `core_vars.go`
- `core_types.go`

## Usage Example

```go
package main

import (
	"fmt"
	"github.com/toxyl/math"
)

func main() {
	// Use a generic math function without manual type casting.
	fmt.Println(math.Max(0.0, uint(1)))

	// Use an extended utility function.
	fmt.Println(math.Blend(10, 20, 0.25))
}
```

## License

This project is released under the **UNLICENSE**. Feel free to fork, modify, and extend the library according to your needs.

Happy coding!