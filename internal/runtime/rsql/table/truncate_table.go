package table

import (
	"log"
	"petacore/internal/core"
	"petacore/internal/storage"
)

// TruncateTable удаляет все строки из таблицы, но оставляет структуру
func (t *Table) TruncateTable(name string) error {
	return t.Storage.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		// Get table metadata to find SERIAL columns
		meta, err := t.GetTableMetadataInTx(tx)
		if err == nil {
			// Reset sequences for SERIAL columns
			for colName, colMeta := range meta.Columns {
				if colMeta.IsSerial {
					seqKey := t.getSequencePrefixKey(colName)
					log.Printf("Resetting sequence for column %s: %s", colName, seqKey)
					tx.Write([]byte(seqKey), "1")
				}
			}
		}

		// Удаляем все строки таблицы
		prefix := t.getRowPrefixKey()
		kvMap, err := tx.Scan([]byte(prefix), core.IteratorTypeAll, -1)
		if err != nil {
			return err
		}
		for key := range kvMap {
			tx.Delete([]byte(key))
		}

		return nil
	})
}
