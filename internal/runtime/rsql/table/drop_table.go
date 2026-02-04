package table

import (
	"petacore/internal/core"
	"petacore/internal/storage"
)

// DropTable удаляет таблицу
func (t *Table) DropTable(name string) error {
	return t.Storage.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		// Удаляем метаданные таблицы
		metaPrefixKey := t.getMetadataPrefixKey()
		tx.Delete([]byte(metaPrefixKey))

		// Удаляем все строки таблицы
		prefix := t.getRowPrefixKey()
		kvMap, err := tx.Scan([]byte(prefix), core.IteratorTypeAll, -1)
		if err != nil {
			return err
		}
		for key := range kvMap {
			tx.Delete([]byte(key))
		}

		// Удаляем sequences
		seqPrefix := t.getAllSequencePrefixKey()
		seqMap, err := tx.Scan([]byte(seqPrefix), core.IteratorTypeAll, -1)
		if err != nil {
			return err
		}
		for key := range seqMap {
			tx.Delete([]byte(key))
		}

		return nil
	})
}
