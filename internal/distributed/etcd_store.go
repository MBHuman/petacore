package distributed

import (
	"context"
	"encoding/json"
	"fmt"
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
func (e *ETCDStore) makeKey(key string) string {
	if e.prefix == "" {
		return key
	}
	return e.prefix + "/" + key
}

// Get получает значение по ключу
func (e *ETCDStore) Get(ctx context.Context, key string) (*KVEntry, error) {
	fullKey := e.makeKey(key)
	resp, err := e.client.Get(ctx, fullKey)
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
func (e *ETCDStore) Put(ctx context.Context, key string, value string, version int64) error {
	fullKey := e.makeKey(key)

	val := etcdValue{
		Value:   value,
		Version: version,
	}

	data, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	_, err = e.client.Put(ctx, fullKey, string(data))
	if err != nil {
		return fmt.Errorf("etcd put error: %w", err)
	}

	return nil
}

// Delete удаляет ключ
func (e *ETCDStore) Delete(ctx context.Context, key string) error {
	fullKey := e.makeKey(key)
	_, err := e.client.Delete(ctx, fullKey)
	if err != nil {
		return fmt.Errorf("etcd delete error: %w", err)
	}
	return nil
}

// Watch следит за изменениями ключей с указанным префиксом
func (e *ETCDStore) Watch(ctx context.Context, prefix string) (<-chan *WatchEvent, error) {
	fullPrefix := e.makeKey(prefix)
	watchChan := e.client.Watch(ctx, fullPrefix, clientv3.WithPrefix())

	eventChan := make(chan *WatchEvent, 100)

	go func() {
		defer close(eventChan)

		for wresp := range watchChan {
			if wresp.Canceled {
				return
			}

			for _, ev := range wresp.Events {
				event := &WatchEvent{}

				// Извлекаем ключ без префикса
				key := string(ev.Kv.Key)
				if e.prefix != "" {
					key = key[len(e.prefix)+1:]
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

					// Обрабатываем предыдущее значение, если есть
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

					// Предыдущее значение при удалении
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

// GetAll получает все ключи с указанным префиксом
func (e *ETCDStore) GetAll(ctx context.Context, prefix string) ([]*KVEntry, error) {
	fullPrefix := e.makeKey(prefix)
	resp, err := e.client.Get(ctx, fullPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("etcd get all error: %w", err)
	}

	entries := make([]*KVEntry, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var val etcdValue
		if err := json.Unmarshal(kv.Value, &val); err != nil {
			continue
		}

		// Извлекаем ключ без префикса
		key := string(kv.Key)
		if e.prefix != "" {
			key = key[len(e.prefix)+1:]
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

func (e *ETCDStore) DeleteAll(ctx context.Context, prefix string) error {
	fullPrefix := e.makeKey(prefix)
	_, err := e.client.Delete(ctx, fullPrefix, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("etcd delete all error: %w", err)
	}
	return nil
}

// Close закрывает соединение с ETCD
func (e *ETCDStore) Close() error {
	return e.client.Close()
}
