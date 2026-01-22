package statements

// SQLStatement представляет SQL statement
// TODO перенести на enum с int и корвентацией через String()
type SQLStatement interface {
	Type() string
}
