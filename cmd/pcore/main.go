package main

import (
	"flag"
	"fmt"
	"runtime/debug"
	"strings"

	"net/http"
	_ "net/http/pprof"
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
	"petacore/sdk/pmem"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	debug.SetGCPercent(400)
	debug.SetMemoryLimit(2 * 1024 * 1024 * 1024)
	// Флаги командной строки
	storeType := flag.String("store", "etcd", "Type of store to use: etcd, inmemory")
	etcdEndpoints := flag.String("etcd-endpoints", "localhost:2379", "ETCD endpoints (comma-separated)")
	etcdPrefix := flag.String("etcd-prefix", "pcore_cluster", "ETCD key prefix")
	nodeID := flag.String("node-id", "node1", "Node identifier for this instance")
	logFile := flag.String("log-file", "", "Log file path (default: stdout)")
	pprofEnabled := flag.Bool("pprof", false, "Enable pprof HTTP server")
	pprofAddr := flag.String("pprof-addr", "localhost:6060", "pprof listen address")
	flag.Parse()

	// Allow overriding flags via environment variables for containerized runs
	if v := os.Getenv("STORE"); v != "" {
		*storeType = v
	}
	if v := os.Getenv("ETCD_ENDPOINTS"); v != "" {
		*etcdEndpoints = v
	}
	if v := os.Getenv("ETCD_PREFIX"); v != "" {
		*etcdPrefix = v
	}
	if v := os.Getenv("NODE_ID"); v != "" {
		*nodeID = v
	}
	if v := os.Getenv("LOG_FILE"); v != "" {
		*logFile = v
	}
	if v := os.Getenv("PPROF"); v != "" {
		lower := strings.ToLower(v)
		*pprofEnabled = (lower == "1" || lower == "true" || lower == "yes")
	}
	if v := os.Getenv("PPROF_ADDR"); v != "" {
		*pprofAddr = v
	}

	if *logFile != "" {
		logger.Init(true, *logFile)
	} else {
		logger.Init(true)
	}
	level := zap.NewAtomicLevel()
	level.SetLevel(zap.ErrorLevel)
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
		// support comma-separated endpoints via env or flag
		eps := strings.Split(*etcdEndpoints, ",")
		for i := range eps {
			eps[i] = strings.TrimSpace(eps[i])
		}
		kv, err = distributed.NewETCDStore(eps, *etcdPrefix)
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Sprintf("Unknown store type: %s", *storeType))
	}

	store := storage.NewDistributedStorageVClock(kv, *nodeID, 1, core.SnapshotIsolation, 1)
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

	// Start pprof server if requested
	if *pprofEnabled {
		go func() {
			logger.Infof("Starting pprof server at %s", *pprofAddr)
			if err := http.ListenAndServe(*pprofAddr, nil); err != nil {
				logger.Errorf("pprof server stopped: %v", err)
			}
		}()
	}

	systemTablesAllocator, err := pmem.NewMmapArena(1028 * 1028)
	if err != nil {
		panic("Failed to create arena allocator: " + err.Error())
	}

	// Initialize system tables
	if err := system.InitializeSystemTables(systemTablesAllocator, store); err != nil {
		logger.Warnf("Failed to initialize system tables: %v", err)
	}
	systemTablesAllocator.Close()

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
