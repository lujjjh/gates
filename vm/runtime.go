package vm

import (
	"math"
	"strconv"
	"strings"
)

func toNumber(value interface{}) interface{} {
	switch value := value.(type) {
	default:
		return math.NaN()
	case int64:
		return value
	case float64:
		return value
	case bool:
		if value {
			return int64(1)
		} else {
			return int64(0)
		}
	case string:
		intValue, err := strconv.ParseInt(value, 0, 64)
		if err == nil {
			return intValue
		}
		floatValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return floatValue
		}
		return math.NaN()
	}
}

func negate(value interface{}) interface{} {
	switch value := toNumber(value).(type) {
	default:
		return math.NaN()
	case int64:
		return -value
	case float64:
		return -value
	}
}

func toBool(value interface{}) bool {
	switch value := value.(type) {
	default:
		return true
	case nil:
		return false
	case int64:
		return value != 0
	case float64:
		return value != 0
	case bool:
		return value
	case string:
		return value != ""
	}
}

func toString(value interface{}) string {
	switch value := value.(type) {
	default:
		return "<object>"
	case nil:
		return "null"
	case int64:
		return strconv.FormatInt(value, 10)
	case float64:
		return strconv.FormatFloat(value, 'g', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []interface{}:
		return "<array>"
	}
}

func add(x, y interface{}) interface{} {
	if x, ok := x.(string); ok {
		return x + toString(y)
	}
	if y, ok := y.(string); ok {
		return toString(x) + y
	}
	switch x := x.(type) {
	case int64:
		switch y := y.(type) {
		case int64:
			return x + y
		case float64:
			return float64(x) + y
		}
	case float64:
		switch y := y.(type) {
		case int64:
			return x + float64(y)
		case float64:
			return x + y
		}
	}
	return math.NaN()
}

func sub(x, y interface{}) interface{} {
	switch x := x.(type) {
	case int64:
		switch y := y.(type) {
		case int64:
			return x - y
		case float64:
			return float64(x) - y
		}
	case float64:
		switch y := y.(type) {
		case int64:
			return x - float64(y)
		case float64:
			return x - y
		}
	}
	return math.NaN()
}

func mul(x, y interface{}) interface{} {
	switch x := x.(type) {
	case int64:
		switch y := y.(type) {
		case int64:
			r := x * y
			// overflow
			if x != 0 && r/x != y {
				return float64(x) * float64(y)
			}
			return r
		case float64:
			return float64(x) * y
		}
	case float64:
		switch y := y.(type) {
		case int64:
			return x * float64(y)
		case float64:
			return x * y
		}
	}
	return math.NaN()
}

func div(x, y interface{}) interface{} {
	switch x := x.(type) {
	case int64:
		switch y := y.(type) {
		case int64:
			return float64(x) / float64(y)
		case float64:
			return float64(x) / y
		}
	case float64:
		switch y := y.(type) {
		case int64:
			return x / float64(y)
		case float64:
			return x / y
		}
	}
	return math.NaN()
}

func mod(x, y interface{}) interface{} {
	switch x := x.(type) {
	case int64:
		switch y := y.(type) {
		case int64:
			if y != 0 {
				return x % y
			}
			return math.NaN()
		case float64:
			return math.Mod(float64(x), y)
		}
	case float64:
		switch y := y.(type) {
		case int64:
			return math.Mod(x, float64(y))
		case float64:
			return math.Mod(x, y)
		}
	}
	return math.NaN()
}

func less(x, y interface{}, defaults bool) bool {
	if x, ok := x.(string); ok {
		if y, ok := y.(string); ok {
			return strings.Compare(x, y) == -1
		}
	}
	nx, ny := toNumber(x), toNumber(y)
	switch nx := nx.(type) {
	case int64:
		switch ny := ny.(type) {
		case int64:
			return nx < ny
		case float64:
			if math.IsNaN(ny) {
				return defaults
			}
			return float64(nx) < ny
		}
	case float64:
		if math.IsNaN(nx) {
			return defaults
		}
		switch ny := ny.(type) {
		case int64:
			return nx < float64(ny)
		case float64:
			if math.IsNaN(ny) {
				return defaults
			}
			return nx < ny
		}
	}
	return defaults
}
