package statements

// SQLStatement представляет SQL statement
type SQLStatement interface {
	Type() string
}
