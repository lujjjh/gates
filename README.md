# Gates

![](https://github.com/lujjjh/gates/workflows/.github/workflows/main.yml/badge.svg)

> An embedded language for Go.

## Features

- Easily embedded in Go.
- JavaScript-like syntax with native int64 support.
- First class functions.
- Execution time limit.

## Comparision

| Features             | Gates | Lua 5.3+                                      | JavaScript |
|----------------------|:-----:|:---------------------------------------------:|:----------:|
| Int64 Support        | Y     | Y                                             | N          |
| Compatible with JSON | Y     | N (hard to distinguish between `[]` and `{}`) | Y          |

## Try Gates in Command Line

```sh
$ go get -u github.com/lujjjh/gates/cmd/gates
$ echo '[1, 2, 3] | map(x => x * x)' | gates
# 1,4,9
```

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
