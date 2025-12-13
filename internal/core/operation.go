package core

type Operation interface {
	Execute(tx *Transaction) error
}

type ReadOperation struct {
	Key string
}

func (ro *ReadOperation) Execute(tx *Transaction) error {
	_, _ = tx.Read(ro.Key)
	return nil
}

type WriteOperation struct {
	Key   string
	Value string
}

func (wo *WriteOperation) Execute(tx *Transaction) error {
	tx.Write(wo.Key, wo.Value)
	return nil
}
