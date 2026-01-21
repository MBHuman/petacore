package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"petacore/internal/core"
	"petacore/internal/distributed"
	"petacore/internal/runtime/wire"
	"petacore/internal/storage"
	"syscall"
)

func main() {
	log.SetOutput(io.Discard)
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

	server := wire.NewWireServer(store, "5432")
	if err := server.Start(); err != nil {
		panic(err)
	}
	defer server.Stop()

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Pg server started")

	<-sigCh
	fmt.Println("Shutting down...")
}
