package distributed

import (
	"context"
	"encoding/json"
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

// etcdValue структура для хранения значения с версией в ETCD
type etcdValue struct {
	Value   string `json:"value"`
	Version int64  `json:"version"`
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
	var val etcdValue
	if err := json.Unmarshal(kv.Value, &val); err != nil {
		return nil, fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return &KVEntry{
		Key:      key,
		Value:    val.Value,
		Version:  val.Version,
		Revision: kv.ModRevision,
	}, nil
}

// Put записывает значение по ключу с версией
func (e *ETCDStore) Put(ctx context.Context, key []byte, value string, version int64) error {
	fullKey := e.makeKey(key)

	val := etcdValue{
		Value:   value,
		Version: version,
	}

	data, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	_, err = e.client.Put(ctx, string(fullKey), string(data))
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
			var val etcdValue
			if err := json.Unmarshal(kv.Value, &val); err != nil {
				continue
			}

			key := kv.Key
			if e.prefix != "" && strings.HasPrefix(string(key), e.prefix+"/") {
				key = []byte(string(key)[len(e.prefix)+1:])
			}

			event := &WatchEvent{
				Type: EventTypePut,
				Entry: &KVEntry{
					Key:      key,
					Value:    val.Value,
					Version:  val.Version,
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
					var val etcdValue
					if err := json.Unmarshal(ev.Kv.Value, &val); err != nil {
						continue
					}
					event.Entry = &KVEntry{
						Key:      key,
						Value:    val.Value,
						Version:  val.Version,
						Revision: ev.Kv.ModRevision,
					}

					if ev.PrevKv != nil {
						var prevVal etcdValue
						if err := json.Unmarshal(ev.PrevKv.Value, &prevVal); err == nil {
							event.PrevEntry = &KVEntry{
								Key:      key,
								Value:    prevVal.Value,
								Version:  prevVal.Version,
								Revision: ev.PrevKv.ModRevision,
							}
						}
					}

				case clientv3.EventTypeDelete:
					event.Type = EventTypeDelete
					event.Entry = &KVEntry{
						Key: key,
					}

					if ev.PrevKv != nil {
						var prevVal etcdValue
						if err := json.Unmarshal(ev.PrevKv.Value, &prevVal); err == nil {
							event.PrevEntry = &KVEntry{
								Key:      key,
								Value:    prevVal.Value,
								Version:  prevVal.Version,
								Revision: ev.PrevKv.ModRevision,
							}
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
		var val etcdValue
		if err := json.Unmarshal(kv.Value, &val); err != nil {
			continue // skip invalid
		}

		key := kv.Key
		if e.prefix != "" && strings.HasPrefix(string(key), e.prefix+"/") {
			key = []byte(string(key)[len(e.prefix)+1:])
		}

		entries = append(entries, &KVEntry{
			Key:      key,
			Value:    val.Value,
			Version:  val.Version,
			Revision: kv.ModRevision,
		})
	}
	return entries, nil
}

// Close закрывает соединение с ETCD
func (e *ETCDStore) Close() error {
	return e.client.Close()
}
