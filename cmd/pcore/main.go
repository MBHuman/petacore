package main

import (
	"fmt"

	"os"
	"os/signal"
	"petacore/internal/core"
	"petacore/internal/distributed"
	"petacore/internal/logger"
	"petacore/internal/runtime/system"
	"petacore/internal/runtime/wire"
	"petacore/internal/storage"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	logger.Init(true)
	level := zap.NewAtomicLevel()
	level.SetLevel(zap.DebugLevel)
	logger.SetLevel(level)
	var kv distributed.KVStore
	var err error

	kv, err = distributed.NewETCDStore([]string{"localhost:2379"}, "pcore_cluster")
	if err != nil {
		panic(err)
	}

	store := storage.NewDistributedStorageVClock(kv, "node1", 1, core.SnapshotIsolation, 1)
	if err := store.Start(); err != nil {
		panic(err)
	}
	defer store.Stop()

	// Initialize system tables
	if err := system.InitializeSystemTables(store); err != nil {
		logger.Warnf("Failed to initialize system tables: %v", err)
	}

	server := wire.NewWireServer(store, "5432")
	if err := server.Start(); err != nil {
		panic(err)
	}
	defer server.Stop()

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Pg server started")

	<-sigCh
	fmt.Println("Shutting down...")
}
