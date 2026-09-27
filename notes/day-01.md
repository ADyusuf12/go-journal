# Day 1

## Environment

- Installed Go 1.25.1
- Verified Go installation with `go version`
- Verified executable path with `which go`
- Installed VS Code Go extension

## New Concepts

### package main

Files belong to packages.

`package main` is special because it produces an executable program.

### func main()

Execution begins here.

The function name does not need to match the filename.

### Modules

Created first module with:

```bash
go mod init hello
```

Generated:

```text
go.mod
```

### Packages

Multiple files can belong to the same package.

Example:

```text
main.go
greetings.go
```

Both can be:

```go
package main
```

and interact directly.

### Variables

Create variable:

```go
name := "Yusuf"
```

Reassign variable:

```go
name = "Ad"
```

### Scope

Variables created inside braces create a new scope.

Example:

```go
{
    name := "Ad"
}
```

This shadows the outer variable.

### Type Inference

```go
age := 30
```

Go inferred:

```text
int
```

### Formatting

Examples:

```go
%s
%d
%f
%T
%v
```

## Biggest Mental Shift

Go cares more about:

- modules
- packages

than filenames.

Coming from Rails, this is a major difference.