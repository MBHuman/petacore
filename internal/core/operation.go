package core

type Operation interface {
	Execute(tx *Transaction) error
}

type ReadOperation struct {
	Key []byte
}

func (ro *ReadOperation) Execute(tx *Transaction) error {
	_, _ = tx.Read(ro.Key)
	return nil
}

type WriteOperation struct {
	Key   []byte
	Value string
}

func (wo *WriteOperation) Execute(tx *Transaction) error {
	tx.Write(wo.Key, wo.Value)
	return nil
}
