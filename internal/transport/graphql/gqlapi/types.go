package gqlapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// MarshalInt8 Int8的序列化
func MarshalInt8(v int8) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(fmt.Sprintf("%d", v)))
	})
}

// UnmarshalInt8 Int8的反序列化
func UnmarshalInt8(v interface{}) (int8, error) {
	switch v := v.(type) {
	case string:
		i, err := strconv.ParseInt(v, 10, 8)
		if err != nil {
			return 0, err
		}
		return int8(i), nil
	// 类型相同直接返回
	case int8:
		return v, nil
	// 高精度判断算术溢出
	case uint8:
		if v > math.MaxInt8 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int8", v)
		}
		return int8(v), nil
	// 高精度判断算术溢出
	case int16:
		if v > math.MaxInt8 || v < math.MinInt8 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int8", v)
		}
		return int8(v), nil
	// 高精度判断算术溢出
	case uint16:
		if v > math.MaxInt8 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int8", v)
		}
		return int8(v), nil
	// 高精度判断算术溢出
	case int:
		if v > math.MaxInt8 || v < math.MinInt8 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int8", v)
		}
		return int8(v), nil
	// 高精度判断算术溢出
	case uint:
		if v > math.MaxInt8 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int8", v)
		}
		return int8(v), nil
	// 高精度判断算术溢出
	case int32:
		if v > math.MaxInt8 || v < math.MinInt8 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int8", v)
		}
		return int8(v), nil
	// 高精度判断算术溢出
	case uint32:
		if v > math.MaxInt8 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int8", v)
		}
		return int8(v), nil
	// 高精度判断算术溢出
	case int64:
		if v > math.MaxInt8 || v < math.MinInt8 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int8", v)
		}
		return int8(v), nil
	// 高精度判断算术溢出
	case uint64:
		if v > math.MaxInt8 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int8", v)
		}
		return int8(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case json.Number:
		i, err := strconv.ParseInt(v.String(), 10, 8)
		if err != nil {
			return 0, err
		}
		return int8(i), nil
	default:
		return 0, fmt.Errorf("%T is not a int8", v)
	}
}

// MarshalUInt8 UInt8的序列化
func MarshalUInt8(v uint8) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(fmt.Sprintf("%d", v)))
	})
}

// UnmarshalUInt8 UInt8的反序列化
func UnmarshalUInt8(v interface{}) (uint8, error) {
	switch v := v.(type) {
	case string:
		i, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			return 0, err
		}
		return uint8(i), nil
	// 类型相同直接返回
	case int8:
		return uint8(v), nil
	// 高精度判断算术溢出
	case uint8:
		return v, nil
	// 高精度判断算术溢出
	case int16:
		if v > math.MaxUint8 || v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint8", v)
		}
		return uint8(v), nil
	// 高精度判断算术溢出
	case uint16:
		if v > math.MaxUint8 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint8", v)
		}
		return uint8(v), nil
	// 高精度判断算术溢出
	case int:
		if v > math.MaxUint8 || v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint8", v)
		}
		return uint8(v), nil
	// 高精度判断算术溢出
	case uint:
		if v > math.MaxUint8 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint8", v)
		}
		return uint8(v), nil
	// 高精度判断算术溢出
	case int32:
		if v > math.MaxUint8 || v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint8", v)
		}
		return uint8(v), nil
	// 高精度判断算术溢出
	case uint32:
		if v > math.MaxUint8 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint8", v)
		}
		return uint8(v), nil
	// 高精度判断算术溢出
	case int64:
		if v > math.MaxUint8 || v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint8", v)
		}
		return uint8(v), nil
	// 高精度判断算术溢出
	case uint64:
		if v > math.MaxUint8 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint8", v)
		}
		return uint8(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case json.Number:
		i, err := strconv.ParseUint(v.String(), 10, 8)
		if err != nil {
			return 0, err
		}
		return uint8(i), nil
	default:
		return 0, fmt.Errorf("%T is not a uint8", v)
	}
}

// MarshalInt16 Int16的序列化
func MarshalInt16(v int16) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(fmt.Sprintf("%d", v)))
	})
}

// UnmarshalInt16 Int16的反序列化
func UnmarshalInt16(v interface{}) (int16, error) {
	switch v := v.(type) {
	case string:
		i, err := strconv.ParseInt(v, 10, 16)
		if err != nil {
			return 0, err
		}
		return int16(i), nil
	// 类型相同直接返回
	case int8, uint8:
		return v.(int16), nil
	// 高精度判断算术溢出
	case int16:
		return v, nil
	// 高精度判断算术溢出
	case uint16:
		if v > math.MaxInt16 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int16", v)
		}
		return int16(v), nil
	// 高精度判断算术溢出
	case int:
		if v > math.MaxInt16 || v < math.MinInt16 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int16", v)
		}
		return int16(v), nil
	// 高精度判断算术溢出
	case uint:
		if v > math.MaxInt16 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int16", v)
		}
		return int16(v), nil
	// 高精度判断算术溢出
	case int32:
		if v > math.MaxInt16 || v < math.MinInt16 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int16", v)
		}
		return int16(v), nil
	// 高精度判断算术溢出
	case uint32:
		if v > math.MaxInt16 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int16", v)
		}
		return int16(v), nil
	// 高精度判断算术溢出
	case int64:
		if v > math.MaxInt16 || v < math.MinInt16 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int16", v)
		}
		return int16(v), nil
	// 高精度判断算术溢出
	case uint64:
		if v > math.MaxInt16 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int16", v)
		}
		return int16(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case json.Number:
		i, err := strconv.ParseInt(v.String(), 10, 16)
		if err != nil {
			return 0, err
		}
		return int16(i), nil
	default:
		return 0, fmt.Errorf("%T is not a int16", v)
	}
}

// MarshalUInt16 UInt16的序列化
func MarshalUInt16(v uint16) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(fmt.Sprintf("%d", v)))
	})
}

// UnmarshalUInt16 UInt16的反序列化
func UnmarshalUInt16(v interface{}) (uint16, error) {
	switch v := v.(type) {
	case string:
		i, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			return 0, err
		}
		return uint16(i), nil
	case uint8:
		return uint16(v), nil
	case int8:
		if v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint16", v)
		}
		return uint16(v), nil
	case int16:
		if v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint16", v)
		}
		return uint16(v), nil
	case uint16:
		return v, nil
	// 高精度判断算术溢出
	case int:
		if v > math.MaxUint16 || v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint16", v)
		}
		return uint16(v), nil
	// 高精度判断算术溢出
	case uint:
		if v > math.MaxUint16 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint16", v)
		}
		return uint16(v), nil
	// 高精度判断算术溢出
	case int32:
		if v > math.MaxUint16 || v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint16", v)
		}
		return uint16(v), nil
	// 高精度判断算术溢出
	case uint32:
		if v > math.MaxUint16 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint16", v)
		}
		return uint16(v), nil
	// 高精度判断算术溢出
	case int64:
		if v > math.MaxUint16 || v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint16", v)
		}
		return uint16(v), nil
	// 高精度判断算术溢出
	case uint64:
		if v > math.MaxUint16 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint16", v)
		}
		return uint16(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case json.Number:
		i, err := strconv.ParseUint(v.String(), 10, 16)
		if err != nil {
			return 0, err
		}
		return uint16(i), nil
	default:
		return 0, fmt.Errorf("%T is not a uint8", v)
	}
}

// MarshalUInt UInt的序列化
func MarshalUInt(v uint) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(fmt.Sprintf("%d", v)))
	})
}

// UnmarshalUInt UInt的反序列化
func UnmarshalUInt(v interface{}) (uint, error) {
	switch v := v.(type) {
	case string:
		i, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(i), nil
	// 低精度直接转换
	case uint8, uint16, uint32:
		return v.(uint), nil
	case int:
		if v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint", v)
		}
		return uint(v), nil
	case int8:
		if v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint", v)
		}
		return uint(v), nil
	case int16:
		if v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint", v)
		}
		return uint(v), nil
	case int32:
		if v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint", v)
		}
		return uint(v), nil
	// 类型相同直接返回
	case uint:
		return v, nil
	// 高精度判断算术溢出
	case int64:
		if v > math.MaxUint32 || v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint", v)
		}
		return uint(v), nil
	// 高精度判断算术溢出
	case uint64:
		if v > math.MaxUint32 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint", v)
		}
		return uint(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case json.Number:
		i, err := strconv.ParseUint(v.String(), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(i), nil
	default:
		return 0, fmt.Errorf("%T is not a uint", v)
	}
}

// MarshalInt32 Int32的序列化
func MarshalInt32(v int32) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(fmt.Sprintf("%d", v)))
	})
}

// UnmarshalInt32 Int32的反序列化
func UnmarshalInt32(v interface{}) (int32, error) {
	switch v := v.(type) {
	case string:
		i, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return 0, err
		}
		return int32(i), nil
	// 低精度直接转换
	case int, int8, int16, uint, uint8, uint16:
		return v.(int32), nil
	// 类型相同直接返回
	case int32:
		return v, nil
	// 高精度判断算术溢出
	case uint32:
		if v > math.MaxInt32 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int32", v)
		}
		return int32(v), nil
	// 高精度判断算术溢出
	case int64:
		if v > math.MaxInt32 || v < math.MinInt32 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int32", v)
		}
		return int32(v), nil
	// 高精度判断算术溢出
	case uint64:
		if v > math.MaxInt32 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int32", v)
		}
		return int32(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case json.Number:
		i, err := strconv.ParseInt(v.String(), 10, 32)
		if err != nil {
			return 0, err
		}
		return int32(i), nil
	default:
		return 0, fmt.Errorf("%T is not a int32", v)
	}
}

// MarshalUInt32 UInt32的序列化
func MarshalUInt32(v uint32) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(fmt.Sprintf("%d", v)))
	})
}

// UnmarshalUInt32 UInt32的反序列化
func UnmarshalUInt32(v interface{}) (uint32, error) {
	switch v := v.(type) {
	case string:
		i, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(i), nil
	// 低精度直接转换
	case uint, uint8, uint16:
		return v.(uint32), nil
	case int:
		if v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint32", v)
		}
		return uint32(v), nil
	case int8:
		if v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint32", v)
		}
		return uint32(v), nil
	case int16:
		if v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint32", v)
		}
		return uint32(v), nil
	case int32:
		if v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint32", v)
		}
		return uint32(v), nil
	// 类型相同直接返回
	case uint32:
		return v, nil
	// 高精度判断算术溢出
	case int64:
		if v > math.MaxUint32 || v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint32", v)
		}
		return uint32(v), nil
	// 高精度判断算术溢出
	case uint64:
		if v > math.MaxUint32 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint32", v)
		}
		return uint32(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case json.Number:
		i, err := strconv.ParseUint(v.String(), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(i), nil
	default:
		return 0, fmt.Errorf("%T is not a uint32", v)
	}
}

// MarshalInt64 Int64的序列化
func MarshalInt64(v int64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(fmt.Sprintf("%d", v)))
	})
}

// UnmarshalInt64 Int64的反序列化
func UnmarshalInt64(v interface{}) (int64, error) {
	switch v := v.(type) {
	case string:
		return strconv.ParseInt(v, 10, 64)
	// 低精度直接转换
	case int, int8, int16, int32, uint, uint8, uint16, uint32:
		return v.(int64), nil
	// 类型相同直接返回
	case int64:
		return v, nil
	// 高精度判断算术溢出
	case uint64:
		if v > math.MaxInt64 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to int64", v)
		}
		return int64(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case json.Number:
		i, err := strconv.ParseInt(v.String(), 10, 64)
		if err != nil {
			return 0, err
		}
		return int64(i), nil
	default:
		return 0, fmt.Errorf("%T is not a int64", v)
	}
}

// MarshalUInt64 UInt64的序列化
func MarshalUInt64(v uint64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		w.Write([]byte(fmt.Sprintf("%d", v)))
	})
}

// UnmarshalUInt64 UInt64的反序列化
func UnmarshalUInt64(v interface{}) (uint64, error) {
	switch v := v.(type) {
	case string:
		return strconv.ParseUint(v, 10, 64)
	// 低精度直接转换
	case uint, uint8, uint16, uint32:
		return v.(uint64), nil
	case int:
		if v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint64", v)
		}
		return uint64(v), nil
	case int8:
		if v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint64", v)
		}
		return uint64(v), nil
	case int16:
		if v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint64", v)
		}
		return uint64(v), nil
	case int32:
		if v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint64", v)
		}
		return uint64(v), nil
	case int64:
		if v < 0 {
			return 0, fmt.Errorf("Arithmetic overflow error converting numeric %d to uint64", v)
		}
		return uint64(v), nil
	// 类型相同直接返回
	case uint64:
		return v, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case json.Number:
		i, err := strconv.ParseUint(v.String(), 10, 64)
		if err != nil {
			return 0, err
		}
		return uint64(i), nil
	default:
		return 0, fmt.Errorf("%T is not a uint64", v)
	}
}

// MarshalTimestamp Timestamp的序列化
func MarshalTimestamp(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.FormatInt(t.Unix(), 10))
	})
}

// UnmarshalTimestamp Timestamp的反序列化
func UnmarshalTimestamp(v interface{}) (time.Time, error) {
	if tmpStr, ok := v.(int64); ok {
		return time.Unix(tmpStr, 0), nil
	}
	return time.Time{}, errors.New("time should be a unix timestamp")
}

// MarshalBytes Bytes的序列化
func MarshalBytes(b []byte) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = fmt.Fprintf(w, "%q", string(b))
	})
}

// UnmarshalBytes Bytes的反序列化
func UnmarshalBytes(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case string:
		return []byte(v), nil
	case *string:
		return []byte(*v), nil
	case []byte:
		return v, nil
	default:
		return nil, fmt.Errorf("%T is not []byte", v)
	}
}
