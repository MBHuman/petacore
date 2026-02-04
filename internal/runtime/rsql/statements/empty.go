package statements

// EmptyStatement представляет пустой statement
type EmptyStatement struct{}

func (e *EmptyStatement) Type() string { return "EMPTY" }
