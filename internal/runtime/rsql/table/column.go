package table

type ColType int

const (
	ColTypeString ColType = iota
	ColTypeInt
	ColTypeFloat
	ColTypeBool
)

func (c ColType) String() string {
	switch c {
	case ColTypeString:
		return "text"
	case ColTypeInt:
		return "integer"
	case ColTypeFloat:
		return "real"
	case ColTypeBool:
		return "boolean"
	default:
		return "text"
	}
}

// ColumnDef определяет колонку
type ColumnDef struct {
	Name         string
	Type         ColType
	IsPrimaryKey bool
	IsNullable   bool
	IsUnique     bool
	IsSerial     bool
	DefaultValue interface{}
}
