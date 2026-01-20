package statements

// ShowStatement представляет SHOW
type ShowStatement struct {
	Parameter string
}

func (s *ShowStatement) Type() string { return "SHOW" }
