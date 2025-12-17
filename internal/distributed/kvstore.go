package distributed

import (
	"context"
)

// KVEntry представляет запись в KV хранилище с метаданными
type KVEntry struct {
	Key      string
	Value    string
	Version  int64 // Версия записи (для MVCC)
	Revision int64 // Ревизия в distributed store (например, ETCD revision)
}

// KVStore - интерфейс для работы с распределенным KV хранилищем
type KVStore interface {
	// Get получает значение по ключу
	Get(ctx context.Context, key string) (*KVEntry, error)

	// Put записывает значение по ключу с версией
	Put(ctx context.Context, key string, value string, version int64) error

	// Delete удаляет ключ
	Delete(ctx context.Context, key string) error

	// SyncIterator возвращает итератор для синхронизации ключей с префиксом
	// После того как Revision достигает ревизии в distributed хранилище,
	// состояние узла считается синхронизированным и он принимать участие в кворуме
	// и принимать запросы
	SyncIterator(ctx context.Context, prefix string) (<-chan *WatchEvent, error)

	// Close закрывает соединение с хранилищем
	Close() error
}

// KVStoreT - тестовый интерфейс, расширяющий KVStore методами для тестирования
type KVStoreT interface {
	KVStore

	// DeleteAll удаляет все ключи с префиксом (для тестирования)
	DeleteAll(ctx context.Context, prefix string) error
}

// WatchEventType тип события изменения
type WatchEventType int

const (
	// EventTypePut событие записи/обновления
	EventTypePut WatchEventType = iota
	// EventTypeDelete событие удаления
	EventTypeDelete
)

// WatchEvent событие изменения в KV хранилище
type WatchEvent struct {
	Type      WatchEventType
	Entry     *KVEntry
	PrevEntry *KVEntry // Предыдущее значение (если было)
}
