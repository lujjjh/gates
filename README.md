# Gates

[![Build Status](https://travis-ci.org/gates/gates.svg?branch=master)](https://travis-ci.org/gates/gates)

> A very simple embedded language, designed for configuration.

## 特性

* 弱类型
* 语法类似 JavaScript

## 内置类型

* Number
* String
* Bool
* Map
* Array

## 语法

### Number

与 Go 的 int64 / float64 语法一致。

`1.5`、`42`、`-1e5`

### String

必须使用双引号。

`"Hello\nworld"`

`assert("Hello"[0] == "H")`
`assert("Hello".length == 5)`

### Bool

`true` / `false`

### 逻辑表达式

惰性求值

`assert(true && 0 || "hello" == "hello")`

## 参考

* https://golang.org/pkg/go/
* https://github.com/dop251/goja/
