package compiler

import (
	"fmt"
	"strconv"

	"github.com/lujjjh/gates"
	"github.com/lujjjh/gates/syntax"
	"github.com/lujjjh/gates/vm"
)

type Compiler struct {
	file      *syntax.File
	constants []gates.Value
	scope     *scope
}

type Error struct {
	Message string
	File    *syntax.File
	Pos     syntax.Pos
}

type SyntaxError struct {
	err *Error
}

func NewSyntaxError(err *Error) *SyntaxError {
	return &SyntaxError{
		err: err,
	}
}

func (e *SyntaxError) Error() string {
	if e.err.File != nil {
		return fmt.Sprintf("SyntaxError: %s at %s", e.err.Message, e.err.File.Position(e.err.Pos))
	}
	return fmt.Sprintf("SyntaxError: %s", e.err.Message)
}

func (c *Compiler) emit(instructions ...byte) {
	c.Instructions = append(c.curFunction.Instructions, instructions...)
}

func (c *Compiler) throwSyntaxError(pos syntax.Pos, format string, args ...interface{}) {
	panic(NewSyntaxError(&Error{
		File:    c.file,
		Pos:     pos,
		Message: fmt.Sprintf(format, args...),
	}))
}

func (c *Compiler) openScope() {
	c.scope = newScope(c.scope)
}

func (c *Compiler) closeScope() {
	c.scope = c.scope.outer
}

func (c *Compiler) compile(node syntax.Node) {
	switch node := node.(type) {
	case *syntax.Ident:
		return

	case *syntax.Lit:
		switch node.Kind {
		case syntax.NUMBER:
			i, err := strconv.ParseInt(node.Value, 0, 64)
			if err == nil {
				c.emit(vm.OpLoadConst)
				c.emit(uint16Operand(c.addConstant(gates.Int(i)))...)
			} else {
				f, _ := strconv.ParseFloat(node.Value, 64)
				c.emit(vm.OpLoadConst)
				c.emit(uint16Operand(c.addConstant(gates.Float(f)))...)
			}

		case syntax.STRING:
			s, _ := strconv.Unquote(node.Value)
			c.emit(vm.OpLoadConst)
			c.emit(uint16Operand(c.addConstant(gates.String(s)))...)

		case syntax.BOOL:
			c.emit(vm.OpLoadConst)
			c.emit(uint16Operand(c.addConstant(gates.Bool(node.Value == "true")))...)

		case syntax.NULL:
			c.emit(vm.OpLoadNull)

		default:
			panic(fmt.Errorf("compile: unknown token type: %v", node.Kind))
		}

	case *syntax.ArrayLit:
		numSegments := 0
		elemList := node.ElemList
		l := len(elemList)
		for i := 0; i < l; {
			elem := elemList[i]
			if elem.Expanded {
				c.compile(elem.Value)
				numSegments++
				i++
				continue
			}
			numElements := 0
			for ; i < l; i++ {
				if elemList[i].Expanded {
					i--
					break
				}
				c.compile(elem.Value)
				numElements++
			}
			c.emit(vm.OpArray)
			c.emit(uint16Operand(numElements)...)
		}
		if numSegments > 1 {
			c.emit(vm.OpMergeArray)
			c.emit(byteOperand(numSegments)...)
		}

	case *syntax.MapLit:
		numSegments := 0
		elemList := node.Entries
		l := len(elemList)
		for i := 0; i < l; {
			elem := elemList[i]
			if elem.Expanded {
				c.compile(elem.Value)
				numSegments++
				i++
				continue
			}
			numElements := 0
			for ; i < l; i++ {
				if elemList[i].Expanded {
					i--
					break
				}
				c.compile(elem.Key)
				c.compile(elem.Value)
				numElements++
			}
			c.emit(vm.OpMap)
			c.emit(uint16Operand(numElements)...)
		}
		if numSegments > 1 {
			c.emit(vm.OpMergeMap)
			c.emit(byteOperand(numSegments)...)
		}

	case *syntax.FunctionLit:
		c.openScope()

		for _, parameter := range node.ParameterList.List {
			c.scope.bindName(parameter.Name)
		}

		c.compile(node.Body.StmtList)

	case *syntax.UnaryExpr:
		return c.compileUnaryExpr(e)
	case *syntax.BinaryExpr:
		return c.compileBinaryExpr(e)
	case *syntax.ParenExpr:
		return c.compileExpr(e.X)
	case *syntax.SelectorExpr:
		return c.compileSelectorExpr(e.X, gates.String(e.Sel.Name), e.Sel.NamePos)
	case *syntax.IndexExpr:
		return c.compileIndexExpr(e.X, e.Index, e.Lbrack)
	case *syntax.CallExpr:
		return c.compileCallExpr(e)
	case *syntax.VarDeclExpr:
		return c.compileVarDeclExpr(e)
	default:
		panic(fmt.Errorf("unknown expression type: %T", e))
	}
}

func byteOperand(i int) []byte {
	return []byte{byte(i)}
}

func uint16Operand(i int) []byte {
	u := uint16(i)
	return []byte{
		byte(u >> 8),
		byte(u),
	}
}

func (c *Compiler) addConstant(value gates.Value) int {
	for idx := range c.constants {
		if c.constants[idx].SameAs(value) {
			return idx
		}
	}
	idx := len(c.constants)
	c.constants = append(c.constants, value)
	return idx
}
