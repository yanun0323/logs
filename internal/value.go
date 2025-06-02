package internal

import (
	"fmt"
	"strconv"
)

const (
	_emptyString = ""
	_trueString  = "true"
	_falseString = "false"
)

func ValueToString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'g', -1, 32)
	case bool:
		if v {
			return _trueString
		}
		return _falseString
	case nil:
		return _emptyString
	default:
		return fmt.Sprintf("%+v", v)
	}
}
