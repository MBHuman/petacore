package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"petacore/internal/core"
	"petacore/internal/distributed"
	"petacore/internal/storage"
	"strings"
	"time"
)

type REPL struct {
	ds      *storage.DistributedStorage
	scanner *bufio.Scanner
}

func NewREPL(ds *storage.DistributedStorage) *REPL {
	return &REPL{
		ds:      ds,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (r *REPL) printHelp() {
	fmt.Println("\n=== Команды ===")
	fmt.Println("  set <key> <value>     - Записать значение")
	fmt.Println("  get <key>             - Прочитать значение")
	fmt.Println("  del <key>             - Удалить ключ (записать пустую строку)")
	fmt.Println("  tx                    - Начать интерактивную транзакцию")
	fmt.Println("  status                - Показать статус синхронизации")
	fmt.Println("  isolation <level>     - Установить уровень изоляции (rc/si)")
	fmt.Println("  help                  - Показать эту справку")
	fmt.Println("  exit                  - Выйти")
	fmt.Println()
}

func (r *REPL) handleSet(args []string) {
	if len(args) < 2 {
		fmt.Println("❌ Использование: set <key> <value>")
		return
	}

	key := args[0]
	value := strings.Join(args[1:], " ")

	err := r.ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		tx.Write(key, value)
		return nil
	})

	if err != nil {
		fmt.Printf("❌ Ошибка записи: %v\n", err)
	} else {
		fmt.Printf("✓ Записано: %s = %s\n", key, value)
	}
}

func (r *REPL) handleGet(args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Использование: get <key>")
		return
	}

	key := args[0]

	err := r.ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		if value, ok := tx.Read(key); ok {
			fmt.Printf("✓ %s = %s\n", key, value)
		} else {
			fmt.Printf("⚠ Ключ %s не найден\n", key)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("❌ Ошибка чтения: %v\n", err)
	}
}

func (r *REPL) handleDelete(args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Использование: del <key>")
		return
	}

	key := args[0]

	err := r.ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		tx.Write(key, "")
		return nil
	})

	if err != nil {
		fmt.Printf("❌ Ошибка удаления: %v\n", err)
	} else {
		fmt.Printf("✓ Ключ %s удалён\n", key)
	}
}

func (r *REPL) handleTransaction() {
	fmt.Println("\n=== Интерактивная транзакция ===")
	fmt.Println("Команды внутри транзакции:")
	fmt.Println("  read <key>            - Прочитать значение")
	fmt.Println("  write <key> <value>   - Записать значение")
	fmt.Println("  commit                - Зафиксировать транзакцию")
	fmt.Println("  rollback              - Отменить транзакцию")
	fmt.Println()

	type txOp struct {
		opType string
		key    string
		value  string
	}

	operations := []txOp{}
	inTransaction := true

	for inTransaction {
		fmt.Print("tx> ")
		if !r.scanner.Scan() {
			break
		}

		line := strings.TrimSpace(r.scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		cmd := parts[0]
		args := parts[1:]

		switch cmd {
		case "read":
			if len(args) < 1 {
				fmt.Println("❌ Использование: read <key>")
				continue
			}
			operations = append(operations, txOp{opType: "read", key: args[0]})
			fmt.Printf("→ Запланировано чтение: %s\n", args[0])

		case "write":
			if len(args) < 2 {
				fmt.Println("❌ Использование: write <key> <value>")
				continue
			}
			key := args[0]
			value := strings.Join(args[1:], " ")
			operations = append(operations, txOp{opType: "write", key: key, value: value})
			fmt.Printf("→ Запланирована запись: %s = %s\n", key, value)

		case "commit":
			fmt.Println("\n→ Выполнение транзакции...")
			err := r.ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
				for _, op := range operations {
					switch op.opType {
					case "read":
						if value, ok := tx.Read(op.key); ok {
							fmt.Printf("  [READ] %s = %s\n", op.key, value)
						} else {
							fmt.Printf("  [READ] %s = <не найдено>\n", op.key)
						}
					case "write":
						tx.Write(op.key, op.value)
						fmt.Printf("  [WRITE] %s = %s\n", op.key, op.value)
					}
				}
				return nil
			})

			if err != nil {
				fmt.Printf("❌ Ошибка транзакции: %v\n", err)
			} else {
				fmt.Println("✓ Транзакция успешно зафиксирована")
			}
			inTransaction = false

		case "rollback":
			fmt.Println("✓ Транзакция отменена")
			inTransaction = false

		default:
			fmt.Printf("❌ Неизвестная команда: %s\n", cmd)
		}
	}
}

func (r *REPL) handleStatus() {
	isSynced := r.ds.IsSynced()

	fmt.Println("\n=== Статус синхронизации ===")
	fmt.Printf("Синхронизирован: %v\n", isSynced)

	if isSynced {
		fmt.Println("Статус: ✓ Узел готов к работе")
	} else {
		fmt.Println("Статус: ⏳ Синхронизация...")
	}
}

func (r *REPL) handleIsolation(args []string, currentLevel *core.IsolationLevel) {
	if len(args) < 1 {
		levelName := "ReadCommitted"
		if *currentLevel == core.SnapshotIsolation {
			levelName = "SnapshotIsolation"
		}
		fmt.Printf("Текущий уровень изоляции: %s\n", levelName)
		fmt.Println("Использование: isolation <rc|si>")
		return
	}

	level := strings.ToLower(args[0])
	switch level {
	case "rc", "readcommitted":
		*currentLevel = core.ReadCommitted
		fmt.Println("✓ Уровень изоляции установлен: Read Committed")
	case "si", "snapshotisolation":
		*currentLevel = core.SnapshotIsolation
		fmt.Println("✓ Уровень изоляции установлен: Snapshot Isolation")
	default:
		fmt.Printf("❌ Неизвестный уровень изоляции: %s\n", level)
		fmt.Println("Доступные уровни: rc (Read Committed), si (Snapshot Isolation)")
	}
}

func (r *REPL) Run() {
	fmt.Println("\n🚀 PetaCore Distributed REPL")
	r.printHelp()

	isolationLevel := core.ReadCommitted

	for {
		fmt.Print("petacore> ")

		if !r.scanner.Scan() {
			break
		}

		line := strings.TrimSpace(r.scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		cmd := parts[0]
		args := parts[1:]

		switch cmd {
		case "help":
			r.printHelp()

		case "set":
			r.handleSet(args)

		case "get":
			r.handleGet(args)

		case "del", "delete":
			r.handleDelete(args)

		case "tx", "transaction":
			r.handleTransaction()

		case "status":
			r.handleStatus()

		case "isolation":
			r.handleIsolation(args, &isolationLevel)

		case "exit", "quit":
			fmt.Println("👋 До свидания!")
			return

		default:
			fmt.Printf("❌ Неизвестная команда: %s (используйте 'help' для справки)\n", cmd)
		}
	}
}

func main() {
	fmt.Println("=== PetaCore Распределённая СУБД ===")
	fmt.Println()

	// Подключаемся к ETCD кластеру
	fmt.Println("📡 Подключение к ETCD кластеру...")
	etcdEndpoints := []string{"localhost:2379", "localhost:2479", "localhost:2579"}

	// Проверяем переменные окружения для ETCD endpoints
	if envEndpoints := os.Getenv("ETCD_ENDPOINTS"); envEndpoints != "" {
		etcdEndpoints = strings.Split(envEndpoints, ",")
	}

	kvStore, err := distributed.NewETCDStore(etcdEndpoints, "petacore")
	if err != nil {
		log.Fatalf("❌ Не удалось подключиться к ETCD: %v\n", err)
		fmt.Println("\nПодсказка: Убедитесь, что ETCD кластер запущен:")
		fmt.Println("  docker-compose up -d")
		os.Exit(1)
	}
	defer kvStore.Close()

	fmt.Println("✓ Подключено к ETCD")

	// Создаем распределенное хранилище с Read Committed
	fmt.Println("🔧 Создание распределенного хранилища...")
	ds := storage.NewDistributedStorage(kvStore, core.ReadCommitted)

	// Запускаем синхронизацию
	if err := ds.Start(); err != nil {
		log.Fatalf("❌ Не удалось запустить синхронизацию: %v\n", err)
	}
	defer ds.Stop()

	fmt.Println("✓ Синхронизация запущена")

	// Ждем завершения начальной синхронизации
	fmt.Print("⏳ Ожидание синхронизации")
	syncContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	synced := false
	for !synced {
		select {
		case <-syncContext.Done():
			fmt.Println("\n⚠ Таймаут ожидания синхронизации, продолжаем...")
			goto skipSync
		default:
			if ds.IsSynced() {
				synced = true
				fmt.Println(" ✓")
			} else {
				fmt.Print(".")
				time.Sleep(200 * time.Millisecond)
			}
		}
	}

skipSync:
	// Запускаем REPL
	repl := NewREPL(ds)
	repl.Run()
}
