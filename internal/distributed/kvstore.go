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

	// Watch следит за изменениями ключей с указанным префиксом
	// Возвращает канал с событиями изменений
	Watch(ctx context.Context, prefix string) (<-chan *WatchEvent, error)

	// GetAll получает все ключи с указанным префиксом
	GetAll(ctx context.Context, prefix string) ([]*KVEntry, error)

	// Close закрывает соединение с хранилищем
	Close() error
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
