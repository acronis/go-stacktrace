# Stack Trace Implementation for Go

This repository provides a comprehensive stack trace implementation for Go, allowing for detailed error handling and reporting. It includes features such as error wrapping, severity levels, and position tracking.

## Features

- **Error Wrapping**: Wrap errors with additional context.
- **Severity Levels**: Define the severity of errors.
- **Position Tracking**: Track the position (line and column) of errors in files.
- **Trace Options**: Customize trace generation with options like ensuring duplicates are not printed.

## Installation

To install the package, use:

```sh
go get github.com/acronis/go-stacktrace
```

## Usage

Basic usage
```Go
package main

import (
    "fmt"
    "github.com/acronis/go-stacktrace"
)

func main() {
    err := stacktrace.New("An error occurred", stacktrace.WithLocation("/path/to/file"), stacktrace.WithPosition(stacktrace.NewPosition(10, 1)))
    fmt.Println(err)
    // Output:
    // /path/to/file:10:1: An error occurred
}
```

Wrapping errors
```Go
package main

import (
    "fmt"
    "github.com/acronis/go-stacktrace"
)

func main() {
    baseErr := fmt.Errorf("base error")
    wrappedErr := stacktrace.NewWrapped("an error occurred", baseErr, stacktrace.WithLocation("/path/to/file"), stacktrace.WithPosition(stacktrace.NewPosition(10, 1)))
    fmt.Println(wrappedErr)
    // Output:
    // /path/to/file:10:1: an error occurred: base error
    
    unwrappedErr, ok := stacktrace.Unwrap(wrappedErr)
    if ok {
        fmt.Println(unwrappedErr)
    }
    
    wrappedErr2 := stacktrace.Wrap(baseErr, stacktrace.WithLocation("/path/to/file"), stacktrace.WithPosition(stacktrace.NewPosition(10, 1)))
    fmt.Println(wrappedErr2)
    // Output:
    // /path/to/file:10:1: base error
}
```

Customizing Traces
```Go
package main

import (
    "fmt"
    "github.com/acronis/go-stacktrace"
)

func main() {
    err := stacktrace.New("an error occurred", stacktrace.WithLocation("/path/to/file"), stacktrace.WithPosition(stacktrace.NewPosition(10, 1)))
    traces := err.GetTraces(stacktrace.WithEnsureDuplicates())
    fmt.Println(traces)
}
```

## API

### Types

* **StackTrace**: Represents a stack trace with various attributes.
* **Severity**: Represents the severity level of an error.
* **Type**: Represents the type of an error.
* **Position**: Represents the position (line and column) of an error in a file.
* **Location**: Represents the location (file path) of an error.

### Functions

* `New(message string, opts ...Option) *StackTrace`: Creates a new stack trace.
* `NewWrapped(message string, err error, opts ...Option) *StackTrace`: Creates a new wrapped stack trace.
* `Wrap(err error, opts ...Option) *StackTrace`: Wraps an existing error in a stack trace.
* `Unwrap(err error) (*StackTrace, bool)`: Unwraps a stack trace from an error.

### Options

* `WithLocation(location string) Option`: Sets the location of the error.
* `WithSeverity(severity Severity) Option`: Sets the severity of the error.
* `WithPosition(position *Position) Option`: Sets the position of the error.
* `WithInfo(key string, value fmt.Stringer) Option`: Adds additional information to the error.
* `WithType(errType Type) Option`: Sets the type of the error.
* `WithEnsureDuplicates() TracesOpt`: Ensures that duplicates are not printed in traces.

##  Contributing
Contributions are welcome! Please open an issue or submit a pull request.
