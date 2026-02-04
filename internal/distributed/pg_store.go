package distributed

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"strings"

// 	"github.com/jackc/pgx/v5"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// // PGStore реализация KVStore для PostgreSQL
// type PGStore struct {
// 	client *pgxpool.Pool
// 	prefix string // Префикс для всех ключей (для namespace isolation)
// }

// // pgValue структура для хранения значения с версией в PostgreSQL
// type pgValue struct {
// 	Value   string `json:"value"`
// 	Version int64  `json:"version"`
// }

// // NewPGStore создает новое подключение к PostgreSQL
// func NewPGStore(connString string, prefix string) (*PGStore, error) {
// 	config, err := pgxpool.ParseConfig(connString)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse pg connection string: %w", err)
// 	}

// 	pool, err := pgxpool.NewWithConfig(context.Background(), config)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create pg connection pool: %w", err)
// 	}

// 	store := &PGStore{
// 		client: pool,
// 		prefix: prefix,
// 	}

// 	// Создаем таблицу если её нет
// 	if err := store.initSchema(context.Background()); err != nil {
// 		pool.Close()
// 		return nil, fmt.Errorf("failed to init schema: %w", err)
// 	}

// 	return store, nil
// }

// // initSchema создает необходимые таблицы и триггеры
// func (ps *PGStore) initSchema(ctx context.Context) error {
// 	queries := []string{
// 		`CREATE TABLE IF NOT EXISTS kv_store (
// 			key TEXT NOT NULL,
// 			value JSONB NOT NULL,
// 			revision BIGSERIAL NOT NULL,
// 			PRIMARY KEY (revision, key)
// 		)`,
// 		`CREATE INDEX IF NOT EXISTS idx_kv_store_key ON kv_store (key)`,
// 		`CREATE OR REPLACE FUNCTION notify_kv_change() RETURNS TRIGGER AS $$
// 		DECLARE
// 			payload JSON;
// 		BEGIN
// 			IF (TG_OP = 'DELETE') THEN
// 				payload = json_build_object(
// 					'type', 'DELETE',
// 					'key', OLD.key,
// 					'prev_value', OLD.value,
// 					'prev_revision', OLD.revision
// 				);
// 				PERFORM pg_notify('kv_changes', payload::text);
// 				RETURN OLD;
// 			ELSIF (TG_OP = 'INSERT') THEN
// 				payload = json_build_object(
// 					'type', 'PUT',
// 					'key', NEW.key,
// 					'value', NEW.value,
// 					'revision', NEW.revision
// 				);
// 				PERFORM pg_notify('kv_changes', payload::text);
// 				RETURN NEW;
// 			ELSIF (TG_OP = 'UPDATE') THEN
// 				payload = json_build_object(
// 					'type', 'PUT',
// 					'key', NEW.key,
// 					'value', NEW.value,
// 					'revision', NEW.revision,
// 					'prev_value', OLD.value,
// 					'prev_revision', OLD.revision
// 				);
// 				PERFORM pg_notify('kv_changes', payload::text);
// 				RETURN NEW;
// 			END IF;
// 		END;
// 		$$ LANGUAGE plpgsql`,
// 		`DROP TRIGGER IF EXISTS kv_change_trigger ON kv_store`,
// 		`CREATE TRIGGER kv_change_trigger
// 			AFTER INSERT OR UPDATE OR DELETE ON kv_store
// 			FOR EACH ROW EXECUTE FUNCTION notify_kv_change()`,
// 	}

// 	for _, query := range queries {
// 		if _, err := ps.client.Exec(ctx, query); err != nil {
// 			return fmt.Errorf("failed to execute query: %w", err)
// 		}
// 	}

// 	return nil
// }

// // makeKey создает полный ключ с префиксом
// func (ps *PGStore) makeKey(key string) string {
// 	if ps.prefix == "" {
// 		return key
// 	}
// 	return ps.prefix + "/" + key
// }

// // Get получает значение по ключу
// func (ps *PGStore) Get(ctx context.Context, key string) (*KVEntry, error) {
// 	fullKey := ps.makeKey(key)

// 	var valueJSON []byte
// 	var revision int64

// 	err := ps.client.QueryRow(ctx,
// 		"SELECT value, revision FROM kv_store WHERE key = $1",
// 		fullKey,
// 	).Scan(&valueJSON, &revision)

// 	if err != nil {
// 		if err == pgx.ErrNoRows {
// 			return nil, fmt.Errorf("key not found: %s", key)
// 		}
// 		return nil, fmt.Errorf("pg get error: %w", err)
// 	}

// 	var val pgValue
// 	if err := json.Unmarshal(valueJSON, &val); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal value: %w", err)
// 	}

// 	return &KVEntry{
// 		Key:      key,
// 		Value:    val.Value,
// 		Version:  val.Version,
// 		Revision: revision,
// 	}, nil
// }

// // Put записывает значение по ключу с версией
// func (ps *PGStore) Put(ctx context.Context, key string, value string, version int64) error {
// 	fullKey := ps.makeKey(key)

// 	val := pgValue{
// 		Value:   value,
// 		Version: version,
// 	}

// 	data, err := json.Marshal(val)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal value: %w", err)
// 	}

// 	_, err = ps.client.Exec(ctx,
// 		`INSERT INTO kv_store (key, value) VALUES ($1, $2)
// 		 ON CONFLICT (key) DO UPDATE SET value = $2`,
// 		fullKey, data,
// 	)

// 	if err != nil {
// 		return fmt.Errorf("pg put error: %w", err)
// 	}

// 	return nil
// }

// // Delete удаляет ключ
// func (ps *PGStore) Delete(ctx context.Context, key string) error {
// 	fullKey := ps.makeKey(key)
// 	_, err := ps.client.Exec(ctx, "DELETE FROM kv_store WHERE key = $1", fullKey)
// 	if err != nil {
// 		return fmt.Errorf("pg delete error: %w", err)
// 	}
// 	return nil
// }

// // SyncIterator возвращает итератор для синхронизации ключей с префиксом
// // Сначала отправляет все существующие ключи как PUT события, затем переключается на listen
// func (ps *PGStore) SyncIterator(ctx context.Context, prefix string) (<-chan *WatchEvent, error) {
// 	fullPrefix := ps.makeKey(prefix)
// 	eventChan := make(chan *WatchEvent, 1000)

// 	go func() {
// 		defer close(eventChan)

// 		// Фаза 1: Получить все существующие ключи пачками
// 		lastRev := int64(0)
// 		batchSize := 200

// 		for {
// 			rows, err := ps.client.Query(ctx,
// 				"SELECT key, value, revision FROM kv_store WHERE key LIKE $1 AND revision > $2 ORDER BY revision LIMIT $3",
// 				fullPrefix+"%", lastRev, batchSize,
// 			)
// 			if err != nil {
// 				return
// 			}

// 			hasRows := false
// 			for rows.Next() {
// 				hasRows = true
// 				var key string
// 				var valueJSON []byte
// 				var revision int64

// 				if err := rows.Scan(&key, &valueJSON, &revision); err != nil {
// 					rows.Close()
// 					return
// 				}

// 				var val pgValue
// 				if err := json.Unmarshal(valueJSON, &val); err != nil {
// 					continue
// 				}

// 				// Извлекаем ключ без префикса
// 				if ps.prefix != "" && strings.HasPrefix(key, ps.prefix+"/") {
// 					key = key[len(ps.prefix)+1:]
// 				}

// 				event := &WatchEvent{
// 					Type: EventTypePut,
// 					Entry: &KVEntry{
// 						Key:      key,
// 						Value:    val.Value,
// 						Version:  val.Version,
// 						Revision: revision,
// 					},
// 				}

// 				select {
// 				case eventChan <- event:
// 				case <-ctx.Done():
// 					rows.Close()
// 					return
// 				}

// 				lastRev = revision
// 			}

// 			rows.Close()

// 			if err := rows.Err(); err != nil {
// 				return
// 			}

// 			if !hasRows {
// 				break
// 			}
// 		}

// 		// Фаза 2: Запустить listen для инкрементальных обновлений
// 		conn, err := ps.client.Acquire(ctx)
// 		if err != nil {
// 			return
// 		}
// 		defer conn.Release()

// 		_, err = conn.Exec(ctx, "LISTEN kv_changes")
// 		if err != nil {
// 			return
// 		}

// 		for {
// 			notification, err := conn.Conn().WaitForNotification(ctx)
// 			if err != nil {
// 				if ctx.Err() != nil {
// 					return
// 				}
// 				continue
// 			}

// 			// Парсим событие
// 			var payload struct {
// 				Type         string          `json:"type"`
// 				Key          string          `json:"key"`
// 				Value        json.RawMessage `json:"value,omitempty"`
// 				Revision     int64           `json:"revision,omitempty"`
// 				PrevValue    json.RawMessage `json:"prev_value,omitempty"`
// 				PrevRevision int64           `json:"prev_revision,omitempty"`
// 			}

// 			if err := json.Unmarshal([]byte(notification.Payload), &payload); err != nil {
// 				continue
// 			}

// 			// Проверяем префикс
// 			if !strings.HasPrefix(payload.Key, fullPrefix) {
// 				continue
// 			}

// 			// Извлекаем ключ без префикса
// 			key := payload.Key
// 			if ps.prefix != "" && strings.HasPrefix(key, ps.prefix+"/") {
// 				key = key[len(ps.prefix)+1:]
// 			}

// 			event := &WatchEvent{}

// 			switch payload.Type {
// 			case "PUT":
// 				event.Type = EventTypePut

// 				var val pgValue
// 				if err := json.Unmarshal(payload.Value, &val); err != nil {
// 					continue
// 				}

// 				event.Entry = &KVEntry{
// 					Key:      key,
// 					Value:    val.Value,
// 					Version:  val.Version,
// 					Revision: payload.Revision,
// 				}

// 				// Обрабатываем предыдущее значение, если есть
// 				if len(payload.PrevValue) > 0 {
// 					var prevVal pgValue
// 					if err := json.Unmarshal(payload.PrevValue, &prevVal); err == nil {
// 						event.PrevEntry = &KVEntry{
// 							Key:      key,
// 							Value:    prevVal.Value,
// 							Version:  prevVal.Version,
// 							Revision: payload.PrevRevision,
// 						}
// 					}
// 				}

// 			case "DELETE":
// 				event.Type = EventTypeDelete
// 				event.Entry = &KVEntry{
// 					Key: key,
// 				}

// 				// Предыдущее значение при удалении
// 				if len(payload.PrevValue) > 0 {
// 					var prevVal pgValue
// 					if err := json.Unmarshal(payload.PrevValue, &prevVal); err == nil {
// 						event.PrevEntry = &KVEntry{
// 							Key:      key,
// 							Value:    prevVal.Value,
// 							Version:  prevVal.Version,
// 							Revision: payload.PrevRevision,
// 						}
// 					}
// 				}
// 			}

// 			select {
// 			case eventChan <- event:
// 			case <-ctx.Done():
// 				return
// 			}
// 		}
// 	}()

// 	return eventChan, nil
// }

// // DeleteAll удаляет все ключи с префиксом (для тестирования)
// func (ps *PGStore) DeleteAll(ctx context.Context, prefix string) error {
// 	fullPrefix := ps.makeKey(prefix)
// 	_, err := ps.client.Exec(ctx, "DELETE FROM kv_store WHERE key LIKE $1", fullPrefix+"%")
// 	return err
// }

// // Close закрывает соединение с хранилищем
// func (ps *PGStore) Close() error {
// 	ps.client.Close()
// 	return nil
// }
