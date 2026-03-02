package distributed

import (
	"context"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// ETCDStore реализация KVStore для ETCD
type ETCDStore struct {
	client *clientv3.Client
	prefix string // Префикс для всех ключей (для namespace isolation)
}

// NewETCDStore создает новое подключение к ETCD
func NewETCDStore(endpoints []string, prefix string) (*ETCDStore, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to etcd: %w", err)
	}

	return &ETCDStore{
		client: cli,
		prefix: prefix,
	}, nil
}

// makeKey создает полный ключ с префиксом
func (e *ETCDStore) makeKey(key []byte) []byte {
	if e.prefix == "" {
		return key
	}
	return []byte(e.prefix + "/" + string(key))
}

// Get получает значение по ключу
func (e *ETCDStore) Get(ctx context.Context, key []byte) (*KVEntry, error) {
	fullKey := e.makeKey(key)
	resp, err := e.client.Get(ctx, string(fullKey), clientv3.WithSerializable())
	if err != nil {
		return nil, fmt.Errorf("etcd get error: %w", err)
	}

	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	kv := resp.Kvs[0]
	// Читаем raw bytes - это vclock формат: [4-byte metaLen][avro(meta)][value bytes]
	valueBytes := []byte(kv.Value)

	// Извлекаем version из метаданных vclock
	version := int64(0)
	if len(valueBytes) >= 4 {
		metaLen := binary.BigEndian.Uint32(valueBytes[:4])
		if int(metaLen) <= len(valueBytes)-4 {
			// В vclock формате version это timestamp, но мы уже не используем это поле
			// Можно оставить 0 или попытаться извлечь из avro метаданных
			version = 0
		}
	}

	return &KVEntry{
		Key:      key,
		Value:    valueBytes,
		Version:  version,
		Revision: kv.ModRevision,
	}, nil
}

// Put записывает значение по ключу с версией
func (e *ETCDStore) Put(ctx context.Context, key []byte, value []byte, version int64) error {
	fullKey := e.makeKey(key)

	// Записываем raw bytes напрямую (vclock формат уже включает все метаданные)
	_, err := e.client.Put(ctx, string(fullKey), string(value))
	if err != nil {
		return fmt.Errorf("etcd put error: %w", err)
	}

	return nil
}

// Delete удаляет ключ
func (e *ETCDStore) Delete(ctx context.Context, key []byte) error {
	fullKey := e.makeKey(key)
	_, err := e.client.Delete(ctx, string(fullKey))
	if err != nil {
		return fmt.Errorf("etcd delete error: %w", err)
	}
	return nil
}

// SyncIterator возвращает итератор для синхронизации ключей с префиксом
// Сначала отправляет все существующие ключи как PUT события, затем переключается на watch
func (e *ETCDStore) SyncIterator(ctx context.Context, prefix []byte) (<-chan *WatchEvent, error) {
	fullPrefix := e.makeKey(prefix)
	eventChan := make(chan *WatchEvent, 1000) // Увеличен буфер для большого количества ключей

	go func() {
		defer close(eventChan)

		// Фаза 1: Получить все существующие ключи
		resp, err := e.client.Get(ctx, string(fullPrefix), clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
		if err != nil {
			// Отправить ошибку как событие? Или просто закрыть канал
			return
		}

		// Отправить существующие ключи как PUT события
		for _, kv := range resp.Kvs {
			// Читаем raw bytes - это vclock формат
			valueBytes := []byte(kv.Value)

			key := kv.Key
			if e.prefix != "" && strings.HasPrefix(string(key), e.prefix+"/") {
				key = []byte(string(key)[len(e.prefix)+1:])
			}

			event := &WatchEvent{
				Type: EventTypePut,
				Entry: &KVEntry{
					Key:      key,
					Value:    valueBytes,
					Version:  0, // Version включен в vclock метаданные
					Revision: kv.ModRevision,
				},
			}

			select {
			case eventChan <- event:
			case <-ctx.Done():
				return
			}
		}

		// Отправить событие завершения синхронизации начальных данных
		syncCompleteEvent := &WatchEvent{
			Type: EventTypeSyncComplete,
		}
		select {
		case eventChan <- syncCompleteEvent:
		case <-ctx.Done():
			return
		}

		// Фаза 2: Запустить watch для инкрементальных обновлений
		watchChan := e.client.Watch(ctx, string(fullPrefix), clientv3.WithPrefix(), clientv3.WithRev(resp.Header.Revision+1))

		for wresp := range watchChan {
			if wresp.Canceled {
				return
			}

			for _, ev := range wresp.Events {
				event := &WatchEvent{}

				key := ev.Kv.Key
				if e.prefix != "" && strings.HasPrefix(string(key), e.prefix+"/") {
					key = []byte(string(key)[len(e.prefix)+1:])
				}

				switch ev.Type {
				case clientv3.EventTypePut:
					event.Type = EventTypePut
					// Читаем raw bytes - vclock формат
					valueBytes := []byte(ev.Kv.Value)
					event.Entry = &KVEntry{
						Key:      key,
						Value:    valueBytes,
						Version:  0, // Version включен в vclock метаданные
						Revision: ev.Kv.ModRevision,
					}

					if ev.PrevKv != nil {
						prevValueBytes := []byte(ev.PrevKv.Value)
						event.PrevEntry = &KVEntry{
							Key:      key,
							Value:    prevValueBytes,
							Version:  0,
							Revision: ev.PrevKv.ModRevision,
						}
					}

				case clientv3.EventTypeDelete:
					event.Type = EventTypeDelete
					event.Entry = &KVEntry{
						Key: key,
					}

					if ev.PrevKv != nil {
						prevValueBytes := []byte(ev.PrevKv.Value)
						event.PrevEntry = &KVEntry{
							Key:      key,
							Value:    prevValueBytes,
							Version:  0,
							Revision: ev.PrevKv.ModRevision,
						}
					}
				}

				select {
				case eventChan <- event:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return eventChan, nil
}

// DeleteAll удаляет все ключи с префиксом (для тестирования)
func (e *ETCDStore) DeleteAll(ctx context.Context, prefix []byte) error {
	fullPrefix := e.makeKey(prefix)
	_, err := e.client.Delete(ctx, string(fullPrefix), clientv3.WithPrefix())
	return err
}

// ScanPrefix сканирует все ключи с префиксом
func (e *ETCDStore) ScanPrefix(ctx context.Context, prefix []byte) ([]*KVEntry, error) {
	fullPrefix := e.makeKey(prefix)
	resp, err := e.client.Get(ctx, string(fullPrefix), clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		return nil, fmt.Errorf("etcd scan error: %w", err)
	}

	var entries []*KVEntry
	for _, kv := range resp.Kvs {
		// Читаем raw bytes - vclock формат
		valueBytes := []byte(kv.Value)

		key := kv.Key
		if e.prefix != "" && strings.HasPrefix(string(key), e.prefix+"/") {
			key = []byte(string(key)[len(e.prefix)+1:])
		}

		entries = append(entries, &KVEntry{
			Key:      key,
			Value:    valueBytes,
			Version:  0, // Version включен в vclock метаданные
			Revision: kv.ModRevision,
		})
	}
	return entries, nil
}

// Close закрывает соединение с ETCD
func (e *ETCDStore) Close() error {
	return e.client.Close()
}
