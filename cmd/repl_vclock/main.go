package main

import (
	"bufio"
	"context"
	"flag"
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
	storage *storage.DistributedStorageVClock
	scanner *bufio.Scanner
	nodeID  string
}

func NewREPL(storage *storage.DistributedStorageVClock, nodeID string) *REPL {
	return &REPL{
		storage: storage,
		scanner: bufio.NewScanner(os.Stdin),
		nodeID:  nodeID,
	}
}

func (r *REPL) printHelp() {
	fmt.Println("\n=== PetaCore VClock REPL - Команды ===")
	fmt.Println("  set <key> <value>     - Записать значение")
	fmt.Println("  get <key>             - Прочитать значение (с quorum проверкой)")
	fmt.Println("  tx                    - Начать интерактивную транзакцию")
	fmt.Println("  status                - Показать статус синхронизации")
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

	err := r.storage.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
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

	err := r.storage.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		if value, ok := tx.Read(key); ok {
			fmt.Printf("✓ %s = %s\n", key, value)
		} else {
			fmt.Printf("⚠ Ключ %s не найден (или нет quorum)\n", key)
			fmt.Printf("  [DEBUG] minAcks=%d, totalNodes=%d\n", r.storage.GetMinAcks(), r.storage.GetTotalNodes())
		}
		return nil
	})

	if err != nil {
		fmt.Printf("❌ Ошибка чтения: %v\n", err)
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
			err := r.storage.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
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
	isSynced := r.storage.IsSynced()

	fmt.Println("\n=== Статус синхронизации ===")
	fmt.Printf("NodeID: %s\n", r.nodeID)
	fmt.Printf("Синхронизирован: %v\n", isSynced)

	if isSynced {
		fmt.Println("Статус: ✓ Узел готов к работе")
	} else {
		fmt.Println("Статус: ⏳ Синхронизация...")
	}
	fmt.Println()
}

func (r *REPL) Run() {
	fmt.Printf("\n🚀 PetaCore VClock REPL [Node: %s]\n", r.nodeID)
	r.printHelp()

	for {
		fmt.Print("vclock> ")

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

		case "tx", "transaction":
			r.handleTransaction()

		case "status":
			r.handleStatus()

		case "exit", "quit":
			fmt.Println("👋 До свидания!")
			return

		default:
			fmt.Printf("❌ Неизвестная команда: %s (используйте 'help' для справки)\n", cmd)
		}
	}
}

func main() {
	// Параметры командной строки
	nodeID := flag.String("node", "node1", "ID узла")
	totalNodes := flag.Int("nodes", 3, "Общее количество узлов в кластере")
	etcdEndpoints := flag.String("etcd", "localhost:2379,localhost:2479,localhost:2579", "ETCD endpoints через запятую")
	flag.Parse()

	fmt.Println("=== PetaCore VClock REPL ===")
	fmt.Printf("NodeID: %s, Total Nodes: %d, Quorum: %d\n", *nodeID, *totalNodes, (*totalNodes/2)+1)
	fmt.Println()

	// Подключаемся к ETCD кластеру
	fmt.Println("📡 Подключение к ETCD кластеру...")
	endpoints := strings.Split(*etcdEndpoints, ",")

	kvStore, err := distributed.NewETCDStore(endpoints, "petacore-vclock")
	if err != nil {
		log.Printf("❌ Не удалось подключиться к ETCD: %v\n", err)
		fmt.Println("\nПодсказка: Убедитесь, что ETCD кластер запущен:")
		fmt.Println("  docker-compose up -d")
		fmt.Println("\nИли используйте один эндпоинт:")
		fmt.Println("  go run cmd/repl_vclock/main.go -node=node1 -nodes=1 -etcd=localhost:2379")
		os.Exit(1)
	}
	defer kvStore.Close()

	fmt.Println("✓ Подключено к ETCD")

	// Создаем распределенное хранилище с VClock
	fmt.Println("🔧 Создание VClock хранилища...")
	storageVClock := storage.NewDistributedStorageVClock(kvStore, *nodeID, *totalNodes, core.ReadCommitted, 0)

	// Запускаем синхронизацию
	if err := storageVClock.Start(); err != nil {
		log.Fatalf("❌ Не удалось запустить синхронизацию: %v\n", err)
	}
	defer storageVClock.Stop()

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
			if storageVClock.IsSynced() {
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
	repl := NewREPL(storageVClock, *nodeID)
	repl.Run()
}
