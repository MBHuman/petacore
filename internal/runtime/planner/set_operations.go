package planner

import (
	"fmt"
	"petacore/internal/runtime/rsql/table"
	"reflect"
)

// applyUnion выполняет операцию UNION
func applyUnion(left, right *table.ExecuteResult, all bool) (*table.ExecuteResult, error) {
	if len(left.Columns) != len(right.Columns) {
		return nil, fmt.Errorf("UNION requires equal number of columns, got %d and %d",
			len(left.Columns), len(right.Columns))
	}

	// Используем колонки из левой части
	result := &table.ExecuteResult{
		Columns: left.Columns,
		Rows:    make([][]interface{}, 0, len(left.Rows)+len(right.Rows)),
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
			key := rowToString(row)
			if !seenRows[key] {
				result.Rows = append(result.Rows, row)
				seenRows[key] = true
			}
		}

		// Добавляем уникальные строки из правой части
		for _, row := range right.Rows {
			key := rowToString(row)
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
	if len(left.Columns) != len(right.Columns) {
		return nil, fmt.Errorf("INTERSECT requires equal number of columns, got %d and %d",
			len(left.Columns), len(right.Columns))
	}

	result := &table.ExecuteResult{
		Columns: left.Columns,
		Rows:    make([][]interface{}, 0),
	}

	if all {
		// INTERSECT ALL - учитываем количество повторений
		rightRowCounts := make(map[string]int)
		for _, row := range right.Rows {
			key := rowToString(row)
			rightRowCounts[key]++
		}

		leftRowCounts := make(map[string]int)
		for _, row := range left.Rows {
			key := rowToString(row)
			leftRowCounts[key]++
		}

		// Добавляем строки, которые есть в обеих частях
		// Количество = min(left count, right count)
		addedRows := make(map[string]int)
		for _, row := range left.Rows {
			key := rowToString(row)
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
			key := rowToString(row)
			rightRows[key] = true
		}

		seenRows := make(map[string]bool)
		for _, row := range left.Rows {
			key := rowToString(row)
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
	if len(left.Columns) != len(right.Columns) {
		return nil, fmt.Errorf("EXCEPT requires equal number of columns, got %d and %d",
			len(left.Columns), len(right.Columns))
	}

	result := &table.ExecuteResult{
		Columns: left.Columns,
		Rows:    make([][]interface{}, 0),
	}

	if all {
		// EXCEPT ALL - вычитаем с учетом количества
		rightRowCounts := make(map[string]int)
		for _, row := range right.Rows {
			key := rowToString(row)
			rightRowCounts[key]++
		}

		removedRows := make(map[string]int)
		for _, row := range left.Rows {
			key := rowToString(row)
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
			key := rowToString(row)
			rightRows[key] = true
		}

		seenRows := make(map[string]bool)
		for _, row := range left.Rows {
			key := rowToString(row)
			if !rightRows[key] && !seenRows[key] {
				result.Rows = append(result.Rows, row)
				seenRows[key] = true
			}
		}
	}

	return result, nil
}

// rowToString преобразует строку в строковое представление для сравнения
func rowToString(row []interface{}) string {
	return fmt.Sprintf("%v", row)
}

// rowEqual проверяет равенство двух строк
func rowEqual(row1, row2 []interface{}) bool {
	if len(row1) != len(row2) {
		return false
	}
	for i := range row1 {
		if !valuesEqual(row1[i], row2[i]) {
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
