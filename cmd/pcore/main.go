package main

import (
	"flag"
	"fmt"

	"os"
	"os/signal"
	"petacore/internal/core"
	"petacore/internal/distributed"
	"petacore/internal/logger"
	"petacore/internal/runtime/functions"
	"petacore/internal/runtime/system"
	"petacore/internal/runtime/wire"
	"petacore/internal/storage"
	baseplugin "petacore/plugins/base"
	psdk "petacore/sdk"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	// Флаги командной строки
	storeType := flag.String("store", "etcd", "Type of store to use: etcd, inmemory")
	etcdEndpoints := flag.String("etcd-endpoints", "localhost:2379", "ETCD endpoints (comma-separated)")
	etcdPrefix := flag.String("etcd-prefix", "pcore_cluster", "ETCD key prefix")
	logFile := flag.String("log-file", "", "Log file path (default: stdout)")
	flag.Parse()

	if *logFile != "" {
		logger.Init(true, *logFile)
	} else {
		logger.Init(true)
	}
	level := zap.NewAtomicLevel()
	level.SetLevel(zap.DebugLevel)
	logger.SetLevel(level)
	var kv distributed.KVStore
	var err error

	// Выбор типа хранилища
	switch *storeType {
	case "inmemory":
		logger.Info("Using in-memory store")
		kv = distributed.NewInMemoryStore()
	case "etcd":
		logger.Infof("Using ETCD store with endpoints: %s", *etcdEndpoints)
		endpoints := []string{*etcdEndpoints}
		kv, err = distributed.NewETCDStore(endpoints, *etcdPrefix)
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Sprintf("Unknown store type: %s", *storeType))
	}

	store := storage.NewDistributedStorageVClock(kv, "node1", 1, core.SnapshotIsolation, 1)
	if err := store.Start(); err != nil {
		panic(err)
	}
	defer store.Stop()

	// Initialize function registry
	funcRegistry := psdk.NewFunctionRegistry()
	functions.SetFunctionRegistry(funcRegistry)

	// Register base plugin
	pluginRegistry := psdk.NewPetaPluginRegistry()
	basePlugin := baseplugin.Plugin
	if err := basePlugin.Init(nil); err != nil {
		logger.Errorf("Failed to init base plugin: %v", err)
		panic(err)
	}
	if err := pluginRegistry.Register(basePlugin); err != nil {
		logger.Errorf("Failed to register base plugin: %v", err)
		panic(err)
	}
	if err := baseplugin.RegisterFunctions(funcRegistry); err != nil {
		logger.Errorf("Failed to register base functions: %v", err)
		panic(err)
	}
	logger.Info("Plugins registered")

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
