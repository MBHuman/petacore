package serializers

import (
	"fmt"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
	"time"
)

func SerializeGeneric(allocator pmem.Allocator, value any, oid ptypes.OID) ([]byte, error) {
	switch oid {
	case ptypes.PTypeBool:
		boolVal, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("expected bool for bool type, got %T", value)
		}
		return BoolSerializerInstance.Serialize(allocator, boolVal)
	case ptypes.PTypeBytea:
		byteaVal, ok := value.([]byte)
		if !ok {
			return nil, fmt.Errorf("expected []byte for bytea type, got %T", value)
		}
		return ByteaSerializerInstance.Serialize(allocator, byteaVal)
	case ptypes.PTypeFloat4:
		floatVal, ok := value.(float32)
		if !ok {
			return nil, fmt.Errorf("expected float32 for float4 type, got %T", value)
		}
		return Float4SerializerInstance.Serialize(allocator, floatVal)
	case ptypes.PTypeBoolArray:
		values := value.([]bool)
		elementBufs := make([][]byte, len(values))
		for i, v := range values {
			buf, err := BoolSerializerInstance.Serialize(allocator, v)
			if err != nil {
				return nil, fmt.Errorf("serialize bool array element %d: %w", i, err)
			}
			elementBufs[i] = buf
		}
		factory := ptypes.BoolFactory
		return NewArraySerializer(ptypes.PTypeBool, factory).Serialize(allocator, elementBufs)
	case ptypes.PTypeFloat4Array:
		values := value.([]float32)
		elementBufs := make([][]byte, len(values))
		for i, v := range values {
			buf, err := Float4SerializerInstance.Serialize(allocator, v)
			if err != nil {
				return nil, fmt.Errorf("serialize float4 array element %d: %w", i, err)
			}
			elementBufs[i] = buf
		}
		factory := ptypes.Float4Factory
		return NewArraySerializer(ptypes.PTypeFloat4, factory).Serialize(allocator, elementBufs)
	case ptypes.PTypeFloat8Array:
		values := value.([]float64)
		elementBufs := make([][]byte, len(values))
		for i, v := range values {
			buf, err := Float8SerializerInstance.Serialize(allocator, v)
			if err != nil {
				return nil, fmt.Errorf("serialize float8 array element %d: %w", i, err)
			}
			elementBufs[i] = buf
		}
		factory := ptypes.Float8Factory
		return NewArraySerializer(ptypes.PTypeFloat8, factory).Serialize(allocator, elementBufs)
	case ptypes.PTypeInt2Array:
		values := value.([]int16)
		elementBufs := make([][]byte, len(values))
		for i, v := range values {
			buf, err := Int2SerializerInstance.Serialize(allocator, v)
			if err != nil {
				return nil, fmt.Errorf("serialize int2 array element %d: %w", i, err)
			}
			elementBufs[i] = buf
		}
		factory := ptypes.Int2Factory
		return NewArraySerializer(ptypes.PTypeInt2, factory).Serialize(allocator, elementBufs)
	case ptypes.PTypeInt4Array:
		values := value.([]int32)
		elementBufs := make([][]byte, len(values))
		for i, v := range values {
			buf, err := Int4SerializerInstance.Serialize(allocator, v)
			if err != nil {
				return nil, fmt.Errorf("serialize int4 array element %d: %w", i, err)
			}
			elementBufs[i] = buf
		}
		factory := ptypes.Int4Factory
		return NewArraySerializer(ptypes.PTypeInt4, factory).Serialize(allocator, elementBufs)
	case ptypes.PTypeTextArray:
		values := value.([]string)
		elementBufs := make([][]byte, len(values))
		for i, v := range values {
			buf, err := TextSerializerInstance.Serialize(allocator, v)
			if err != nil {
				return nil, fmt.Errorf("serialize text array element %d: %w", i, err)
			}
			elementBufs[i] = buf
		}
		factory := ptypes.TextFactory
		return NewArraySerializer(ptypes.PTypeText, factory).Serialize(allocator, elementBufs)
	case ptypes.PTypeDate:
		dateVal, ok := value.(*time.Time)
		if !ok {
			return nil, fmt.Errorf("expected *time.Time for date type, got %T", value)
		}
		return DateSerializerInstance.Serialize(allocator, dateVal)
	case ptypes.PTypeFloat8:
		floatVal, ok := value.(float64)
		if !ok {
			return nil, fmt.Errorf("expected float64 for float8 type, got %T", value)
		}
		return Float8SerializerInstance.Serialize(allocator, floatVal)
	case ptypes.PTypeInt2:
		intVal, ok := value.(int16)
		if !ok {
			return nil, fmt.Errorf("expected int16 for int2 type, got %T", value)
		}
		return Int2SerializerInstance.Serialize(allocator, intVal)
	case ptypes.PTypeInt4:
		intVal, ok := value.(int32)
		if !ok {
			return nil, fmt.Errorf("expected int32 for int4 type, got %T", value)
		}
		return Int4SerializerInstance.Serialize(allocator, intVal)
	case ptypes.PTypeInt8:
		intVal, ok := value.(int64)
		if !ok {
			return nil, fmt.Errorf("expected int64 for int8 type, got %T", value)
		}
		return Int8SerializerInstance.Serialize(allocator, intVal)
	case ptypes.PTypeNumeric:
		numVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("expected string for numeric type, got %T", value)
		}
		return NumericSerializerInstance.Serialize(allocator, numVal)
	case ptypes.PTypeText, ptypes.PTypeVarchar, ptypes.PTypeName, ptypes.PTypeChar:
		strVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("expected string for text type, got %T", value)
		}
		return TextSerializerInstance.Serialize(allocator, strVal)
	case ptypes.PTypeTime:
		timeVal, ok := value.(*time.Time)
		if !ok {
			return nil, fmt.Errorf("expected *time.Time for time type, got %T", value)
		}
		return TimeSerializerInstance.Serialize(allocator, timeVal)
	case ptypes.PTypeTimestamp:
		timeVal, ok := value.(*time.Time)
		if !ok {
			return nil, fmt.Errorf("expected *time.Time for timestamp type, got %T", value)
		}
		return TimestampSerializerInstance.Serialize(allocator, timeVal)
	case ptypes.PTypeTimestampz:
		timeVal, ok := value.(*time.Time)
		if !ok {
			return nil, fmt.Errorf("expected *time.Time for timestamptz type, got %T", value)
		}
		return TimestampzSerializerInstance.Serialize(allocator, timeVal)

	default:
		return nil, fmt.Errorf("unsupported OID: %d", oid)
	}
}

func DeserializeGeneric(data []byte, oid ptypes.OID) (ptypes.BaseType[any], error) {
	switch oid {
	case ptypes.PTypeBool:
		val, err := BoolSerializerInstance.Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize bool: %w", err)
		}
		return ptypes.NewAnyWrapper[bool](val), nil
	case ptypes.PTypeBytea:
		val, err := ByteaSerializerInstance.Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize bytea: %w", err)
		}
		return ptypes.NewAnyWrapper[[]byte](val), nil
	case ptypes.PTypeFloat4:
		val, err := Float4SerializerInstance.Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize float4: %w", err)
		}
		return ptypes.NewAnyWrapper[float32](val), nil
	case ptypes.PTypeFloat8:
		val, err := Float8SerializerInstance.Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize float8: %w", err)
		}
		return ptypes.NewAnyWrapper[float64](val), nil
	case ptypes.PTypeBoolArray:
		val, err := NewArraySerializer(ptypes.PTypeBool, ptypes.BoolFactory).Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize bool array: %w", err)
		}
		return ptypes.NewAnyWrapper(val), nil
	case ptypes.PTypeFloat4Array:
		val, err := NewArraySerializer(ptypes.PTypeFloat4, ptypes.Float4Factory).Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize float4 array: %w", err)
		}
		return ptypes.NewAnyWrapper(val), nil
	case ptypes.PTypeFloat8Array:
		val, err := NewArraySerializer(ptypes.PTypeFloat8, ptypes.Float8Factory).Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize float8 array: %w", err)
		}
		return ptypes.NewAnyWrapper(val), nil
	case ptypes.PTypeInt2Array:
		val, err := NewArraySerializer(ptypes.PTypeInt2, ptypes.Int2Factory).Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize int2 array: %w", err)
		}
		return ptypes.NewAnyWrapper(val), nil
	case ptypes.PTypeInt4Array:
		val, err := NewArraySerializer(ptypes.PTypeInt4, ptypes.Int4Factory).Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize int4 array: %w", err)
		}
		return ptypes.NewAnyWrapper(val), nil
	case ptypes.PTypeTextArray:
		val, err := NewArraySerializer(ptypes.PTypeText, ptypes.TextFactory).Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize text array: %w", err)
		}
		return ptypes.NewAnyWrapper(val), nil
	case ptypes.PTypeInt2:
		val, err := Int2SerializerInstance.Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize int2: %w", err)
		}
		return ptypes.NewAnyWrapper(val), nil
	case ptypes.PTypeInt4:
		val, err := Int4SerializerInstance.Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize int4: %w", err)
		}
		return ptypes.NewAnyWrapper(val), nil
	case ptypes.PTypeInt8:
		val, err := Int8SerializerInstance.Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize int8: %w", err)
		}
		return ptypes.NewAnyWrapper(val), nil
	case ptypes.PTypeNumeric:
		val, err := NumericSerializerInstance.Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize numeric: %w", err)
		}
		return ptypes.NewAnyWrapper(val), nil
	case ptypes.PTypeText, ptypes.PTypeVarchar, ptypes.PTypeName, ptypes.PTypeChar:
		val, err := TextSerializerInstance.Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize text: %w", err)
		}
		return ptypes.NewAnyWrapper(val), nil
	case ptypes.PTypeDate:
		val, err := DateSerializerInstance.Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize date: %w", err)
		}
		return ptypes.NewAnyWrapper(val), nil
	case ptypes.PTypeTime:
		val, err := TimeSerializerInstance.Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize time: %w", err)
		}
		return ptypes.NewAnyWrapper(val), nil
	case ptypes.PTypeTimestamp:
		val, err := TimestampSerializerInstance.Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize timestamp: %w", err)
		}
		return ptypes.NewAnyWrapper(val), nil
	case ptypes.PTypeTimestampz:
		val, err := TimestampzSerializerInstance.Deserialize(data)
		if err != nil {
			return nil, fmt.Errorf("deserialize timestampz: %w", err)
		}
		return ptypes.NewAnyWrapper(val), nil
	default:
		return nil, fmt.Errorf("unsupported OID: %d", oid)
	}
}

func ValidateGeneric(value []byte, oid ptypes.OID) error {
	switch oid {
	case ptypes.PTypeBool:
		v := ptypes.BoolFactory(value)
		return BoolSerializerInstance.Validate(v)
	case ptypes.PTypeBytea:
		v := ptypes.ByteaFactory(value)
		return ByteaSerializerInstance.Validate(v)
	case ptypes.PTypeFloat4:
		v := ptypes.Float4Factory(value)
		return Float4SerializerInstance.Validate(v)
	case ptypes.PTypeFloat8:
		v := ptypes.Float8Factory(value)
		return Float8SerializerInstance.Validate(v)

	case ptypes.PTypeInt2:
		v := ptypes.Int2Factory(value)
		return Int2SerializerInstance.Validate(v)

	case ptypes.PTypeInt4:
		v := ptypes.Int4Factory(value)
		return Int4SerializerInstance.Validate(v)

	case ptypes.PTypeInt8:
		v := ptypes.Int8Factory(value)
		return Int8SerializerInstance.Validate(v)

	// case ptypes.PTypeNumeric: // TODO придумать как использовать без meta (тут надо брать meta из определения колонки)
	// 	v := ptypes.NumericFactory(value)
	// 	return NumericSerializerInstance.Validate(v)
	case ptypes.PTypeText, ptypes.PTypeVarchar, ptypes.PTypeName, ptypes.PTypeChar:
		// v, ok := value.(ptypes.TypeText)
		// if !ok {
		// 	return fmt.Errorf("expected TypeText, got %T", value)
		// }
		v := ptypes.TextFactory(value)
		return TextSerializerInstance.Validate(v)

	// case ptypes.PTypeVarchar: // TODO придумать как использовать без meta
	// 	// v, ok := value.(ptypes.TypeVarchar)
	// 	// if !ok {
	// 	// 	return fmt.Errorf("expected TypeVarchar, got %T", value)
	// 	// }
	// 	v := ptypes.VarcharFactory(value)
	// 	return VarcharSerializerInstance.Validate(v)
	case ptypes.PTypeDate:
		v := ptypes.DateFactory(value)
		return DateSerializerInstance.Validate(v)
	case ptypes.PTypeTime:
		v := ptypes.TimeFactory(value)
		return TimeSerializerInstance.Validate(v)
	case ptypes.PTypeTimestamp:
		v := ptypes.TimestampFactory(value)
		return TimestampSerializerInstance.Validate(v)
	case ptypes.PTypeTimestampz:
		v := ptypes.TimestampzFactory(value)
		return TimestampzSerializerInstance.Validate(v)
	// TODO Придумать как валидировать массивы без meta (там же надо брать meta из определения колонки) и только со знанием OID и []byte приходящим
	// case ptypes.PTypeBoolArray:

	// 	v, ok := value.(ptypes.TypeArray[bool, ptypes.TypeBool])
	// 	if !ok {
	// 		return fmt.Errorf("expected TypeArray[bool, TypeBool], got %T", value)
	// 	}
	// 	return NewArraySerializer(ptypes.PTypeBool, ptypes.BoolFactory).Validate(&v)

	// case ptypes.PTypeFloat4Array:
	// 	v, ok := value.(ptypes.TypeArray[float32, ptypes.TypeFloat4])
	// 	if !ok {
	// 		return fmt.Errorf("expected TypeArray[float32, TypeFloat4], got %T", value)
	// 	}
	// 	return NewArraySerializer(ptypes.PTypeFloat4, ptypes.Float4Factory).Validate(&v)

	// case ptypes.PTypeFloat8Array:
	// 	v, ok := value.(ptypes.TypeArray[float64, ptypes.TypeFloat8])
	// 	if !ok {
	// 		return fmt.Errorf("expected TypeArray[float64, TypeFloat8], got %T", value)
	// 	}
	// 	return NewArraySerializer(ptypes.PTypeFloat8, ptypes.Float8Factory).Validate(&v)

	// case ptypes.PTypeInt2Array:
	// 	v, ok := value.(ptypes.TypeArray[int16, ptypes.TypeInt2])
	// 	if !ok {
	// 		return fmt.Errorf("expected TypeArray[int16, TypeInt2], got %T", value)
	// 	}
	// 	return NewArraySerializer(ptypes.PTypeInt2, ptypes.Int2Factory).Validate(&v)

	// case ptypes.PTypeInt4Array:
	// 	v, ok := value.(ptypes.TypeArray[int32, ptypes.TypeInt4])
	// 	if !ok {
	// 		return fmt.Errorf("expected TypeArray[int32, TypeInt4], got %T", value)
	// 	}
	// 	return NewArraySerializer(ptypes.PTypeInt4, ptypes.Int4Factory).Validate(&v)

	// case ptypes.PTypeTextArray:
	// 	v, ok := value.(ptypes.TypeArray[string, ptypes.TypeText])
	// 	if !ok {
	// 		return fmt.Errorf("expected TypeArray[string, TypeText], got %T", value)
	// 	}
	// 	return NewArraySerializer(ptypes.PTypeText, ptypes.TextFactory).Validate(&v)

	default:
		return fmt.Errorf("unsupported OID: %d", oid)
	}
}
