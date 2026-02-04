package distributed

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// InMemoryStore реализация KVStore в памяти для тестирования
type InMemoryStore struct {
	mu       sync.RWMutex
	data     map[string]*KVEntry
	revision int64 // Глобальная ревизия для имитации распределенного хранилища
}

// NewInMemoryStore создает новый in-memory store
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data:     make(map[string]*KVEntry),
		revision: 0,
	}
}

// Get получает значение по ключу
func (s *InMemoryStore) Get(ctx context.Context, key []byte) (*KVEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.data[string(key)]
	if !exists {
		return nil, fmt.Errorf("key not found: %s", string(key))
	}

	// Возвращаем копию, чтобы избежать race conditions
	return &KVEntry{
		Key:      append([]byte(nil), entry.Key...),
		Value:    entry.Value,
		Version:  entry.Version,
		Revision: entry.Revision,
	}, nil
}

// Put записывает значение по ключу с версией
func (s *InMemoryStore) Put(ctx context.Context, key []byte, value string, version int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.revision++
	s.data[string(key)] = &KVEntry{
		Key:      append([]byte(nil), key...),
		Value:    value,
		Version:  version,
		Revision: s.revision,
	}

	return nil
}

// Delete удаляет ключ
func (s *InMemoryStore) Delete(ctx context.Context, key []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, string(key))
	s.revision++

	return nil
}

// ScanPrefix сканирует все ключи с префиксом
func (s *InMemoryStore) ScanPrefix(ctx context.Context, prefix []byte) ([]*KVEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	prefixStr := string(prefix)
	var entries []*KVEntry

	for key, entry := range s.data {
		if strings.HasPrefix(key, prefixStr) {
			// Возвращаем копию, чтобы избежать race conditions
			entries = append(entries, &KVEntry{
				Key:      append([]byte(nil), entry.Key...),
				Value:    entry.Value,
				Version:  entry.Version,
				Revision: entry.Revision,
			})
		}
	}

	return entries, nil
}

// SyncIterator возвращает итератор для синхронизации ключей с префиксом
// Для in-memory реализации просто возвращаем все существующие данные и сразу сигнализируем о завершении
func (s *InMemoryStore) SyncIterator(ctx context.Context, prefix []byte) (<-chan *WatchEvent, error) {
	ch := make(chan *WatchEvent, 100)

	go func() {
		defer close(ch)

		// Получаем все существующие ключи
		entries, err := s.ScanPrefix(ctx, prefix)
		if err != nil {
			return
		}

		// Отправляем все существующие записи как PUT события
		for _, entry := range entries {
			select {
			case ch <- &WatchEvent{
				Type:      EventTypePut,
				Entry:     entry,
				PrevEntry: nil,
			}:
			case <-ctx.Done():
				return
			}
		}

		// Сигнализируем о завершении синхронизации
		select {
		case ch <- &WatchEvent{
			Type: EventTypeSyncComplete,
		}:
		case <-ctx.Done():
			return
		}
	}()

	return ch, nil
}

// Close закрывает соединение с хранилищем (no-op для in-memory)
func (s *InMemoryStore) Close() error {
	return nil
}

// DeleteAll удаляет все ключи с префиксом (для тестирования)
func (s *InMemoryStore) DeleteAll(ctx context.Context, prefix string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key := range s.data {
		if strings.HasPrefix(key, prefix) {
			delete(s.data, key)
		}
	}

	s.revision++
	return nil
}

// Clear очищает все данные (для тестирования)
func (s *InMemoryStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]*KVEntry)
	s.revision = 0
}

// Size возвращает количество записей в хранилище (для отладки)
func (s *InMemoryStore) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}
