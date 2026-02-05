package psdk

type OID int32

type IPType interface {
	// Required methods
	TypInput(input []byte) (PType, error)
	TypOutput(value PType) ([]byte, error)
	TypReceive(input []byte) (PType, error)
	TypSend(value PType) ([]byte, error)
}

type PType struct {
	OID   OID
	Value []byte

	Meta PTypeMeta
}

type PTypeMeta struct {
	TypName        string // Type name
	TypNamespace   OID    // Namespace OID
	TypOwner       int32  // Owner's user ID
	TypLen         int16  // Length in bytes
	TypByVal       bool   // Passed by value
	TypType        rune   // Type type
	TypCategory    rune   // Type category
	TypIsPreferred bool   // Is this a preferred type?
	TypIsDefined   bool   // Is this a user-defined type?
	// TODO check if we need more fields from pg_type
	TypDelim       rune  // Delimiter
	TypRelid       int32 // Relation OID if composite type
	TypSubScript   string
	TypElem        int32
	TypArray       int32
	TypInputFunc   string
	TypOutputFunc  string
	TypReceiveFunc string
	TypSendFunc    string
	TypModIn       string
	TypModOut      string
	TypAnalyzeFunc string
	TypAlign       rune // Alignment requirement
	TypStorage     rune
	TypNotNull     bool   // Is NOT NULL
	TypBaseType    OID    // Base type OID if this is a domain
	TypTypMod      int32  // Type modifier
	TypCollocation int32  // Collation OID
	TypDefaultBin  string // Default value (as a binary string)
	TypDefault     string // Default value (as a string)
}
