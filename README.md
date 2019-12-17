# Gates

![](https://github.com/lujjjh/gates/workflows/.github/workflows/main.yml/badge.svg)

> An embedded language interacting with Go.

## Features

- Easily embedded in and interacting with Go.
- JavaScript-like syntax with native int64 support.
- First class functions.
- Execution time limit.

## Why?

### Why Gates?

Gates is designed to be an interpreted language embedded in Go. It aims to providing a
relatively controllable VM so that multiple VMs could run _untrusted_ code in the same process
with maximum execution time set.

### Why not Lua?

Lua is great except that arrays and maps are both represented as `table`. In our use cases,
we need to distinguish between empty arrays and empty maps to produce a correct JSON string.

### Why not JavaScript?

int64 :).

## Data Types

- number (int64 / float64)
- string
- bool
- map
- array
- function

## Examples

[View Examples](/examples/)

## Credits

- https://github.com/dop251/goja/
- https://golang.org/pkg/go/
