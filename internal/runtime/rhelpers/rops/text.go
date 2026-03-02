package rops

import (
	"fmt"
	"petacore/internal/utils"
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
	"strconv"
)

// CoerceToString конвертирует любое значение в строку для конкатенации
func CoerceToString(allocator pmem.Allocator, buf []byte, oid ptypes.OID) (string, error) {
	switch oid {
	case ptypes.PTypeText, ptypes.PTypeVarchar:
		// Zero-copy conversion через utils.BytesToString
		return utils.BytesToString(buf), nil

	case ptypes.PTypeInt2:
		v, err := serializers.Int2SerializerInstance.Deserialize(buf)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(int64(v.IntoGo()), 10), nil

	case ptypes.PTypeInt4:
		v, err := serializers.Int4SerializerInstance.Deserialize(buf)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(int64(v.IntoGo()), 10), nil

	case ptypes.PTypeInt8:
		v, err := serializers.Int8SerializerInstance.Deserialize(buf)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(v.IntoGo(), 10), nil

	case ptypes.PTypeFloat4:
		v, err := serializers.Float4SerializerInstance.Deserialize(buf)
		if err != nil {
			return "", err
		}
		return strconv.FormatFloat(float64(v.IntoGo()), 'f', -1, 32), nil

	case ptypes.PTypeFloat8:
		v, err := serializers.Float8SerializerInstance.Deserialize(buf)
		if err != nil {
			return "", err
		}
		return strconv.FormatFloat(v.IntoGo(), 'f', -1, 64), nil

	// TODO добавить поддержку, надо чтобы в строку переводился
	// case ptypes.PTypeNumeric:
	// 	v, err := serializers.NumericSerializerInstance.Deserialize(buf)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	f, err := v.
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	return f.Text('f', int(v.Meta.Scale)), nil

	case ptypes.PTypeBool:
		v, err := serializers.BoolSerializerInstance.Deserialize(buf)
		if err != nil {
			return "", err
		}
		if v.IntoGo() {
			return "true", nil
		}
		return "false", nil

	case ptypes.PTypeDate:
		v, err := serializers.DateSerializerInstance.Deserialize(buf)
		if err != nil {
			return "", err
		}
		tm := v.IntoGo()
		if tm == nil {
			return "", fmt.Errorf("coerceToString: nil date")
		}
		return tm.Format("2006-01-02"), nil

	case ptypes.PTypeTimestamp:
		v, err := serializers.TimestampSerializerInstance.Deserialize(buf)
		if err != nil {
			return "", err
		}
		tm := v.IntoGo()
		if tm == nil {
			return "", fmt.Errorf("coerceToString: nil timestamp")
		}
		return tm.Format("2006-01-02 15:04:05"), nil

	case ptypes.PTypeTimestampz:
		v, err := serializers.TimestampzSerializerInstance.Deserialize(buf)
		if err != nil {
			return "", err
		}
		tm := v.IntoGo()
		if tm == nil {
			return "", fmt.Errorf("coerceToString: nil timestampz")
		}
		return tm.UTC().Format("2006-01-02 15:04:05+00"), nil

	default:
		return "", fmt.Errorf("coerceToString: unsupported OID %d", oid)
	}
}
