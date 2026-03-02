package serializers

import (
	"fmt"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
	"strconv"
)

type BaseSerializer[GoT any, T any] interface {
	Serializer[GoT]
	Deserializer[T]
	Validator[T]
}

type Serializer[T any] interface {
	Serialize(allocator pmem.Allocator, value T) ([]byte, error)
}

type Deserializer[T any] interface {
	Deserialize([]byte) (T, error)
}

type Validator[T any] interface {
	Validate(T) error
}

func CoerceToString(buf []byte, oid ptypes.OID) (string, error) {
	switch oid {
	case ptypes.PTypeText, ptypes.PTypeVarchar:
		return string(buf), nil

	case ptypes.PTypeInt2:
		v, err := Int2SerializerInstance.Deserialize(buf)
		if err != nil {
			return "", fmt.Errorf("coerce to string: int2: %w", err)
		}
		return strconv.FormatInt(int64(v.IntoGo()), 10), nil

	case ptypes.PTypeInt4:
		v, err := Int4SerializerInstance.Deserialize(buf)
		if err != nil {
			return "", fmt.Errorf("coerce to string: int4: %w", err)
		}
		return strconv.FormatInt(int64(v.IntoGo()), 10), nil

	case ptypes.PTypeInt8:
		v, err := Int8SerializerInstance.Deserialize(buf)
		if err != nil {
			return "", fmt.Errorf("coerce to string: int8: %w", err)
		}
		return strconv.FormatInt(v.IntoGo(), 10), nil

	case ptypes.PTypeFloat4:
		v, err := Float4SerializerInstance.Deserialize(buf)
		if err != nil {
			return "", fmt.Errorf("coerce to string: float4: %w", err)
		}
		return strconv.FormatFloat(float64(v.IntoGo()), 'f', -1, 32), nil

	case ptypes.PTypeFloat8:
		v, err := Float8SerializerInstance.Deserialize(buf)
		if err != nil {
			return "", fmt.Errorf("coerce to string: float8: %w", err)
		}
		return strconv.FormatFloat(v.IntoGo(), 'f', -1, 64), nil

	case ptypes.PTypeNumeric:
		v, err := NumericSerializerInstance.Deserialize(buf)
		if err != nil {
			return "", fmt.Errorf("coerce to string: numeric: %w", err)
		}
		nv, err := v.ToNumericValue()
		if err != nil {
			return "", fmt.Errorf("coerce to string: numeric value: %w", err)
		}
		f := nv.ToBigFloat()
		return f.Text('f', int(v.Meta.Scale)), nil

	case ptypes.PTypeBool:
		v, err := BoolSerializerInstance.Deserialize(buf)
		if err != nil {
			return "", fmt.Errorf("coerce to string: bool: %w", err)
		}
		if v.IntoGo() {
			return "true", nil
		}
		return "false", nil

	case ptypes.PTypeDate:
		v, err := DateSerializerInstance.Deserialize(buf)
		if err != nil {
			return "", fmt.Errorf("coerce to string: date: %w", err)
		}
		tm := v.IntoGo()
		if tm == nil {
			return "", fmt.Errorf("coerce to string: nil date")
		}
		return tm.Format("2006-01-02"), nil

	case ptypes.PTypeTime:
		v, err := TimeSerializerInstance.Deserialize(buf)
		if err != nil {
			return "", fmt.Errorf("coerce to string: time: %w", err)
		}
		tm := v.IntoGo()
		if tm == nil {
			return "", fmt.Errorf("coerce to string: nil time")
		}
		return tm.Format("15:04:05"), nil

	case ptypes.PTypeTimestamp:
		v, err := TimestampSerializerInstance.Deserialize(buf)
		if err != nil {
			return "", fmt.Errorf("coerce to string: timestamp: %w", err)
		}
		tm := v.IntoGo()
		if tm == nil {
			return "", fmt.Errorf("coerce to string: nil timestamp")
		}
		return tm.Format("2006-01-02 15:04:05"), nil

	case ptypes.PTypeTimestampz:
		v, err := TimestampzSerializerInstance.Deserialize(buf)
		if err != nil {
			return "", fmt.Errorf("coerce to string: timestampz: %w", err)
		}
		tm := v.IntoGo()
		if tm == nil {
			return "", fmt.Errorf("coerce to string: nil timestampz")
		}
		return tm.UTC().Format("2006-01-02 15:04:05+00"), nil

	default:
		return "", fmt.Errorf("coerce to string: unsupported OID %d", oid)
	}
}
