package planner

import (
	"bytes"
	"fmt"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
	"reflect"
)

// applyUnion выполняет операцию UNION
func applyUnion(left, right *table.ExecuteResult, all bool) (*table.ExecuteResult, error) {
	if len(left.Schema.Fields) != len(right.Schema.Fields) {
		return nil, fmt.Errorf("UNION requires equal number of columns, got %d and %d",
			len(left.Schema.Fields), len(right.Schema.Fields))
	}

	// Используем схему из левой части
	result := &table.ExecuteResult{
		Schema: left.Schema,
		Rows:   make([]*ptypes.Row, 0, len(left.Rows)+len(right.Rows)),
	}

	if all {
		// UNION ALL - просто объединяем все строки
		result.Rows = append(result.Rows, left.Rows...)
		result.Rows = append(result.Rows, right.Rows...)
	} else {
		// UNION - объединяем и удаляем дубликаты
		seenRows := make(map[string]bool)

		// Добавляем строки из левой части
		for _, row := range left.Rows {
			key := rowToString(row, left.Schema)
			if !seenRows[key] {
				result.Rows = append(result.Rows, row)
				seenRows[key] = true
			}
		}

		// Добавляем уникальные строки из правой части
		for _, row := range right.Rows {
			key := rowToString(row, right.Schema)
			if !seenRows[key] {
				result.Rows = append(result.Rows, row)
				seenRows[key] = true
			}
		}
	}

	return result, nil
}

// applyIntersect выполняет операцию INTERSECT
func applyIntersect(left, right *table.ExecuteResult, all bool) (*table.ExecuteResult, error) {
	if len(left.Schema.Fields) != len(right.Schema.Fields) {
		return nil, fmt.Errorf("INTERSECT requires equal number of columns, got %d and %d",
			len(left.Schema.Fields), len(right.Schema.Fields))
	}

	result := &table.ExecuteResult{
		Schema: left.Schema,
		Rows:   make([]*ptypes.Row, 0),
	}

	if all {
		// INTERSECT ALL - учитываем количество повторений
		rightRowCounts := make(map[string]int)
		for _, row := range right.Rows {
			key := rowToString(row, right.Schema)
			rightRowCounts[key]++
		}

		leftRowCounts := make(map[string]int)
		for _, row := range left.Rows {
			key := rowToString(row, left.Schema)
			leftRowCounts[key]++
		}

		// Добавляем строки, которые есть в обеих частях
		// Количество = min(left count, right count)
		addedRows := make(map[string]int)
		for _, row := range left.Rows {
			key := rowToString(row, left.Schema)
			if rightRowCounts[key] > 0 {
				if addedRows[key] < min(leftRowCounts[key], rightRowCounts[key]) {
					result.Rows = append(result.Rows, row)
					addedRows[key]++
				}
			}
		}
	} else {
		// INTERSECT - возвращаем уникальные строки, которые есть в обеих частях
		rightRows := make(map[string]bool)
		for _, row := range right.Rows {
			key := rowToString(row, right.Schema)
			rightRows[key] = true
		}

		seenRows := make(map[string]bool)
		for _, row := range left.Rows {
			key := rowToString(row, left.Schema)
			if rightRows[key] && !seenRows[key] {
				result.Rows = append(result.Rows, row)
				seenRows[key] = true
			}
		}
	}

	return result, nil
}

// applyExcept выполняет операцию EXCEPT (разность)
func applyExcept(left, right *table.ExecuteResult, all bool) (*table.ExecuteResult, error) {
	if len(left.Schema.Fields) != len(right.Schema.Fields) {
		return nil, fmt.Errorf("EXCEPT requires equal number of columns, got %d and %d",
			len(left.Schema.Fields), len(right.Schema.Fields))
	}

	result := &table.ExecuteResult{
		Schema: left.Schema,
		Rows:   make([]*ptypes.Row, 0),
	}

	if all {
		// EXCEPT ALL - вычитаем с учетом количества
		rightRowCounts := make(map[string]int)
		for _, row := range right.Rows {
			key := rowToString(row, right.Schema)
			rightRowCounts[key]++
		}

		removedRows := make(map[string]int)
		for _, row := range left.Rows {
			key := rowToString(row, left.Schema)
			// Добавляем строку, если она еще не была "вычтена"
			if removedRows[key] < rightRowCounts[key] {
				removedRows[key]++
			} else {
				result.Rows = append(result.Rows, row)
			}
		}
	} else {
		// EXCEPT - возвращаем уникальные строки из левой части, которых нет в правой
		rightRows := make(map[string]bool)
		for _, row := range right.Rows {
			key := rowToString(row, right.Schema)
			rightRows[key] = true
		}

		seenRows := make(map[string]bool)
		for _, row := range left.Rows {
			key := rowToString(row, left.Schema)
			if !rightRows[key] && !seenRows[key] {
				result.Rows = append(result.Rows, row)
				seenRows[key] = true
			}
		}
	}

	return result, nil
}

// rowToString преобразует строку в строковое представление для сравнения
func rowToString(row *ptypes.Row, schema *serializers.BaseSchema) string {
	if row == nil || schema == nil {
		return "nil"
	}

	var buf bytes.Buffer
	for i := range schema.Fields {
		if i > 0 {
			buf.WriteString("|")
		}
		fieldBuf, oid, err := schema.GetField(row, i)
		if err != nil {
			buf.WriteString("<error>")
			continue
		}
		if fieldBuf == nil {
			buf.WriteString("<null>")
		} else {
			val, err := serializers.DeserializeGeneric(fieldBuf, oid)
			if err != nil {
				buf.WriteString("<deserr>")
			} else {
				buf.WriteString(fmt.Sprintf("%v", val))
			}
		}
	}
	return buf.String()
}

// rowEqual проверяет равенство двух строк
func rowEqual(row1, row2 *ptypes.Row, schema1, schema2 *serializers.BaseSchema) bool {
	if len(schema1.Fields) != len(schema2.Fields) {
		return false
	}
	for i := range schema1.Fields {
		buf1, oid1, err1 := schema1.GetField(row1, i)
		buf2, oid2, err2 := schema2.GetField(row2, i)

		if err1 != nil || err2 != nil {
			return false
		}

		// Сравниваем raw bytes если OID совпадают
		if oid1 != oid2 {
			return false
		}

		if !bytes.Equal(buf1, buf2) {
			return false
		}
	}
	return true
}

// valuesEqual сравнивает два значения
func valuesEqual(v1, v2 interface{}) bool {
	// Обработка nil
	if v1 == nil && v2 == nil {
		return true
	}
	if v1 == nil || v2 == nil {
		return false
	}

	// Используем reflect для сравнения
	return reflect.DeepEqual(v1, v2)
}
