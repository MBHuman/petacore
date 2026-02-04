package distributed

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// InMemoryStore реализация KVStore в памяти для тестирования
type InMemoryStore struct {
	mu        sync.RWMutex
	data      map[string]*KVEntry
	revision  int64 // Глобальная ревизия для имитации распределенного хранилища
	watchers  map[string][]*watcher
	watcherMu sync.RWMutex
}

// watcher представляет подписчика на изменения
type watcher struct {
	prefix string
	ch     chan *WatchEvent
	ctx    context.Context
}

// NewInMemoryStore создает новый in-memory store
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data:     make(map[string]*KVEntry),
		watchers: make(map[string][]*watcher),
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

	keyStr := string(key)
	prevEntry := s.data[keyStr]

	s.revision++
	newEntry := &KVEntry{
		Key:      append([]byte(nil), key...),
		Value:    value,
		Version:  version,
		Revision: s.revision,
	}
	s.data[keyStr] = newEntry
	s.mu.Unlock()

	// Уведомляем наблюдателей
	s.notifyWatchers(keyStr, EventTypePut, newEntry, prevEntry)

	return nil
}

// Delete удаляет ключ
func (s *InMemoryStore) Delete(ctx context.Context, key []byte) error {
	s.mu.Lock()

	keyStr := string(key)
	prevEntry := s.data[keyStr]
	delete(s.data, keyStr)
	s.revision++
	s.mu.Unlock()

	// Уведомляем наблюдателей
	s.notifyWatchers(keyStr, EventTypeDelete, nil, prevEntry)

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

func (s *InMemoryStore) SyncIterator(ctx context.Context, prefix []byte) (<-chan *WatchEvent, error) {
	ch := make(chan *WatchEvent, 100)
	prefixStr := string(prefix)

	// Регистрируем watcher
	w := &watcher{
		prefix: prefixStr,
		ch:     ch,
		ctx:    ctx,
	}

	s.watcherMu.Lock()
	s.watchers[prefixStr] = append(s.watchers[prefixStr], w)
	s.watcherMu.Unlock()

	go func() {
		defer func() {
			// Удаляем watcher при завершении
			s.watcherMu.Lock()
			watchers := s.watchers[prefixStr]
			for i, watcher := range watchers {
				if watcher == w {
					s.watchers[prefixStr] = append(watchers[:i], watchers[i+1:]...)
					break
				}
			}
			if len(s.watchers[prefixStr]) == 0 {
				delete(s.watchers, prefixStr)
			}
			s.watcherMu.Unlock()
			close(ch)
		}()

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

		// Сигнализируем о завершении начальной синхронизации
		select {
		case ch <- &WatchEvent{
			Type: EventTypeSyncComplete,
		}:
		case <-ctx.Done():
			return
		}

		// Продолжаем слушать изменения через context
		<-ctx.Done()
	}()

	return ch, nil
}

// notifyWatchers уведомляет всех подписчиков об изменении
func (s *InMemoryStore) notifyWatchers(key string, eventType WatchEventType, entry, prevEntry *KVEntry) {
	s.watcherMu.RLock()
	defer s.watcherMu.RUnlock()

	for prefix, watchers := range s.watchers {
		if !strings.HasPrefix(key, prefix) {
			continue
		}

		event := &WatchEvent{
			Type:  eventType,
			Entry: entry,
		}

		if prevEntry != nil {
			event.PrevEntry = &KVEntry{
				Key:      append([]byte(nil), prevEntry.Key...),
				Value:    prevEntry.Value,
				Version:  prevEntry.Version,
				Revision: prevEntry.Revision,
			}
		}

		if entry != nil {
			event.Entry = &KVEntry{
				Key:      append([]byte(nil), entry.Key...),
				Value:    entry.Value,
				Version:  entry.Version,
				Revision: entry.Revision,
			}
		}

		for _, w := range watchers {
			select {
			case w.ch <- event:
			case <-w.ctx.Done():
			default:
				// Канал переполнен, пропускаем событие
			}
		}
	}
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
