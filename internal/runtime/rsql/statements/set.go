package statements

// SetStatement представляет SET
type SetStatement struct {
	Variable string
	Value    interface{}
}

func (s *SetStatement) Type() string { return "SET" }
