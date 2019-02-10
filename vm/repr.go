package vm

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"unsafe"
)

// Represent returns a string containing a printable representation
// of a value. It detects circular references. In that case, it returns
// a "[Circular]".
//
// Map's keys are sorted so that multiple calls on a map should produce
// the same result.
func Represent(x Value) string {
	return represent(x, make(map[unsafe.Pointer]struct{}))
}

func represent(x Value, visited map[unsafe.Pointer]struct{}) string {
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Map, reflect.Slice:
		addr := unsafe.Pointer(v.Pointer())
		if _, seen := visited[addr]; seen {
			return "[Circular]"
		}
		visited[addr] = struct{}{}
		defer delete(visited, addr)
	}

	switch x := x.(type) {
	case []interface{}:
		return representArray(x, visited)
	case map[string]interface{}:
		return representMap(x, visited)
	case string:
		return strconv.Quote(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case float64:
		return strconv.FormatFloat(x, 'G', -1, 64)
	case nil:
		return "null"
	case bool:
		return strconv.FormatBool(x)
	default:
		return fmt.Sprintf("unknown %#v %s", v, reflect.TypeOf(x).Name())
	}
}

func representArray(x []interface{}, visited map[unsafe.Pointer]struct{}) string {
	var buf bytes.Buffer
	buf.WriteString("[")

	needComma := false
	for _, v := range x {
		if needComma {
			buf.WriteString(",")
		} else {
			needComma = true
		}
		buf.WriteString(represent(v, visited))
	}
	buf.WriteString("]")

	return buf.String()
}

func representMap(x map[string]interface{}, visited map[unsafe.Pointer]struct{}) string {
	var buf bytes.Buffer
	buf.WriteString("{")

	keys := make([]string, 0, len(x))
	for k := range x {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	needComma := false
	for _, k := range keys {
		if needComma {
			buf.WriteString(",")
		} else {
			needComma = true
		}
		v := x[k]
		buf.WriteString(represent(k, visited))
		buf.WriteString(":")
		buf.WriteString(represent(v, visited))
	}
	buf.WriteString("}")

	return buf.String()
}
