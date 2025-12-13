package main

import (
	"context"
	"fmt"
	"log"
	"petacore/internal/core"
	"petacore/internal/distributed"
	"petacore/internal/storage"
	"time"
)

func main() {
	fmt.Println("=== Распределенная СУБД с ETCD и MVCC ===")
	fmt.Println()

	// Подключаемся к ETCD кластеру
	fmt.Println("Подключение к ETCD кластеру...")
	etcdEndpoints := []string{"localhost:2379", "localhost:2479", "localhost:2579"}

	kvStore, err := distributed.NewETCDStore(etcdEndpoints, "petacore")
	if err != nil {
		log.Fatalf("Не удалось подключиться к ETCD: %v", err)
	}
	defer kvStore.Close()

	fmt.Println("✓ Подключено к ETCD")
	fmt.Println()

	// Создаем распределенное хранилище с Read Committed
	fmt.Println("Создание распределенного хранилища...")
	ds := storage.NewDistributedStorage(kvStore, core.ReadCommitted)

	// Запускаем синхронизацию
	if err := ds.Start(); err != nil {
		log.Fatalf("Не удалось запустить синхронизацию: %v", err)
	}
	defer ds.Stop()

	fmt.Println("✓ Синхронизация запущена")
	fmt.Println()

	// Ждем завершения начальной синхронизации
	fmt.Println("Ожидание синхронизации...")
	for !ds.IsSynced() {
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("✓ Узел синхронизирован")
	fmt.Println()

	// Демонстрация работы: запись
	fmt.Println("=== Запись данных ===")
	err = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		tx.Write("user:1", "Alice")
		tx.Write("user:2", "Bob")
		tx.Write("balance:1", "1000")
		tx.Write("balance:2", "2000")
		return nil
	})
	if err != nil {
		log.Fatalf("Ошибка при записи: %v", err)
	}
	fmt.Println("✓ Данные записаны в ETCD и синхронизированы")
	fmt.Println()

	// Даем время на синхронизацию между узлами
	time.Sleep(500 * time.Millisecond)

	// Демонстрация работы: чтение из локального MVCC кеша
	fmt.Println("=== Чтение данных из локального кеша ===")
	err = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		if user, ok := tx.Read("user:1"); ok {
			fmt.Printf("user:1 = %s\n", user)
		}
		if balance, ok := tx.Read("balance:1"); ok {
			fmt.Printf("balance:1 = %s\n", balance)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Ошибка при чтении: %v", err)
	}
	fmt.Println()

	// Демонстрация обновления
	fmt.Println("=== Обновление баланса ===")
	err = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		balance, ok := tx.Read("balance:1")
		if !ok {
			return fmt.Errorf("balance not found")
		}
		fmt.Printf("Текущий баланс: %s\n", balance)

		// Обновляем баланс
		tx.Write("balance:1", "1500")
		fmt.Println("Новый баланс: 1500")
		return nil
	})
	if err != nil {
		log.Fatalf("Ошибка при обновлении: %v", err)
	}
	fmt.Println("✓ Баланс обновлен")
	fmt.Println()

	// Демонстрация чтения обновленных данных
	time.Sleep(500 * time.Millisecond)

	fmt.Println("=== Проверка обновленных данных ===")
	err = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		balance, ok := tx.Read("balance:1")
		if !ok {
			return fmt.Errorf("balance not found")
		}
		fmt.Printf("Обновленный баланс: %s\n", balance)
		return nil
	})
	if err != nil {
		log.Fatalf("Ошибка при чтении: %v", err)
	}
	fmt.Println()

	// Демонстрация прямого доступа к ETCD
	fmt.Println("=== Прямое чтение из ETCD ===")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	entry, err := kvStore.Get(ctx, "user:1")
	if err != nil {
		log.Printf("Ошибка при чтении из ETCD: %v", err)
	} else {
		fmt.Printf("Из ETCD: user:1 = %s (version: %d, revision: %d)\n",
			entry.Value, entry.Version, entry.Revision)
	}

	fmt.Println()
	fmt.Println("=== Архитектура ===")
	fmt.Println("1. ETCD - источник истины (source of truth)")
	fmt.Println("2. Локальный MVCC - кеш для быстрого чтения")
	fmt.Println("3. Запись: через ETCD → синхронизация на все узлы")
	fmt.Println("4. Чтение: из локального MVCC кеша (быстро)")
	fmt.Println("5. CP модель: консистентность + устойчивость к разделению")
	fmt.Println()

	fmt.Println("Статус синхронизации:", ds.GetSyncStatus())
}
