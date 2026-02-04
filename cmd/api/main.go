package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"petacore/internal/core"
	"petacore/internal/distributed"
	"petacore/internal/storage"
	"strings"
	"syscall"
	"time"
)

// APIServer представляет HTTP API сервер
type APIServer struct {
	storage *storage.DistributedStorageVClock
	server  *http.Server
}

// WriteRequest запрос на запись
type WriteRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ReadResponse ответ на чтение
type ReadResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Found bool   `json:"found"`
}

// StatusResponse статус узла
type StatusResponse struct {
	NodeID      string `json:"node_id"`
	IsSynced    bool   `json:"is_synced"`
	TotalNodes  int    `json:"total_nodes"`
	MinAcks     int    `json:"min_acks"`
	Description string `json:"description"`
}

// HealthResponse проверка здоровья
type HealthResponse struct {
	Status string `json:"status"`
	NodeID string `json:"node_id"`
}

// ErrorResponse ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

// SetMinAcksRequest запрос на изменение minAcks
type SetMinAcksRequest struct {
	MinAcks int `json:"min_acks"`
}

// SetMinAcksResponse ответ на изменение minAcks
type SetMinAcksResponse struct {
	Status      string `json:"status"`
	OldMinAcks  int    `json:"old_min_acks"`
	NewMinAcks  int    `json:"new_min_acks"`
	TotalNodes  int    `json:"total_nodes"`
	Description string `json:"description"`
}

// NewAPIServer создает новый API сервер
func NewAPIServer(storage *storage.DistributedStorageVClock, port string) *APIServer {
	mux := http.NewServeMux()

	server := &APIServer{
		storage: storage,
		server: &http.Server{
			Addr:         ":" + port,
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	// Регистрируем обработчики
	mux.HandleFunc("/health", server.handleHealth)
	mux.HandleFunc("/status", server.handleStatus)
	mux.HandleFunc("/write", server.handleWrite)
	mux.HandleFunc("/read", server.handleRead)
	mux.HandleFunc("/config/min_acks", server.handleSetMinAcks)
	mux.HandleFunc("/", server.handleRoot)

	return server
}

// handleRoot корневой обработчик
func (s *APIServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"service":   "PetaCore VClock API",
		"version":   "1.0",
		"endpoints": "/health, /status, /write, /read, /config/min_acks",
	})
}

// handleHealth проверка здоровья
func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := HealthResponse{
		Status: "ok",
		NodeID: os.Getenv("NODE_ID"),
	}

	if !s.storage.IsSynced() {
		response.Status = "syncing"
	}

	json.NewEncoder(w).Encode(response)
}

// handleStatus статус узла
func (s *APIServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := StatusResponse{
		NodeID:      os.Getenv("NODE_ID"),
		IsSynced:    s.storage.IsSynced(),
		TotalNodes:  s.storage.GetTotalNodes(),
		MinAcks:     s.storage.GetMinAcks(),
		Description: fmt.Sprintf("VClock distributed storage with quorum (min %d/%d nodes)", s.storage.GetMinAcks(), s.storage.GetTotalNodes()),
	}

	json.NewEncoder(w).Encode(response)
}

// handleWrite обработчик записи
func (s *APIServer) handleWrite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req WriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON: " + err.Error()})
		return
	}

	if req.Key == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Key is required"})
		return
	}

	// Выполняем транзакцию записи
	err := s.storage.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		tx.Write([]byte(req.Key), req.Value)
		return nil
	})

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Write failed: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "Value written successfully",
		"key":     req.Key,
	})
}

// handleRead обработчик чтения
func (s *APIServer) handleRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Key parameter is required"})
		return
	}

	var response ReadResponse
	response.Key = key

	// Выполняем транзакцию чтения с quorum проверкой
	err := s.storage.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		if value, ok := tx.Read([]byte(key)); ok {
			response.Value = value
			response.Found = true
		} else {
			response.Found = false
		}
		return nil
	})

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Read failed: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !response.Found {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(response)
}

// handleSetMinAcks обработчик изменения minAcks
func (s *APIServer) handleSetMinAcks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SetMinAcksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON: " + err.Error()})
		return
	}

	oldMinAcks := s.storage.GetMinAcks()
	s.storage.SetMinAcks(req.MinAcks)
	newMinAcks := s.storage.GetMinAcks()

	var description string
	switch req.MinAcks {
	case 0:
		description = fmt.Sprintf("Set to default quorum (%d/%d + 1)", s.storage.GetTotalNodes(), 2)
	case -1:
		description = fmt.Sprintf("Set to all nodes (%d/%d)", newMinAcks, s.storage.GetTotalNodes())
	default:
		description = fmt.Sprintf("Set to custom value (%d/%d)", newMinAcks, s.storage.GetTotalNodes())
	}

	response := SetMinAcksResponse{
		Status:      "ok",
		OldMinAcks:  oldMinAcks,
		NewMinAcks:  newMinAcks,
		TotalNodes:  s.storage.GetTotalNodes(),
		Description: description,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	log.Printf("MIN_ACKS changed: %d -> %d (%s)", oldMinAcks, newMinAcks, description)
}

// Start запускает сервер
func (s *APIServer) Start() error {
	log.Printf("Starting API server on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown корректно останавливает сервер
func (s *APIServer) Shutdown(ctx context.Context) error {
	log.Println("Shutting down API server...")
	return s.server.Shutdown(ctx)
}

func main() {
	log.Println("=== PetaCore VClock API Server ===")

	// Получаем параметры из переменных окружения
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		nodeID = "node-1"
		log.Printf("NODE_ID not set, using default: %s", nodeID)
	}

	etcdEndpointsStr := os.Getenv("ETCD_ENDPOINTS")
	if etcdEndpointsStr == "" {
		etcdEndpointsStr = "localhost:2379"
		log.Printf("ETCD_ENDPOINTS not set, using default: %s", etcdEndpointsStr)
	}
	etcdEndpoints := strings.Split(etcdEndpointsStr, ",")

	totalNodesStr := os.Getenv("TOTAL_NODES")
	totalNodes := 3
	if totalNodesStr != "" {
		fmt.Sscanf(totalNodesStr, "%d", &totalNodes)
	}
	log.Printf("Total nodes: %d", totalNodes)

	// Настраиваемый минимум подтверждений для quorum
	minAcksStr := os.Getenv("MIN_ACKS")
	minAcks := 0 // 0 означает использовать значение по умолчанию (N/2 + 1)
	if minAcksStr != "" {
		fmt.Sscanf(minAcksStr, "%d", &minAcks)
		if minAcks == 0 {
			log.Printf("MIN_ACKS=0: using default quorum %d/%d + 1 = %d", totalNodes, 2, totalNodes/2+1)
		} else if minAcks == -1 {
			log.Printf("MIN_ACKS=-1: using all nodes = %d", totalNodes)
		} else {
			log.Printf("MIN_ACKS=%d: custom value", minAcks)
		}
	} else {
		log.Printf("MIN_ACKS not set, using default quorum: %d/%d + 1 = %d", totalNodes, 2, totalNodes/2+1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("PORT not set, using default: %s", port)
	}

	// Подключаемся к ETCD
	log.Printf("Connecting to ETCD cluster: %v", etcdEndpoints)
	kvStore, err := distributed.NewETCDStore(etcdEndpoints, "petacore-vclock")
	if err != nil {
		log.Fatalf("Failed to connect to ETCD: %v", err)
	}
	defer kvStore.Close()
	log.Println("✓ Connected to ETCD")

	// Создаем распределенное хранилище с VClock и Read Committed
	log.Printf("Creating distributed storage with VClock for node: %s", nodeID)
	distStorage := storage.NewDistributedStorageVClock(
		kvStore,
		nodeID,
		totalNodes,
		core.ReadCommitted,
		minAcks,
	)

	// Запускаем синхронизацию
	log.Println("Starting synchronization...")
	if err := distStorage.Start(); err != nil {
		log.Fatalf("Failed to start synchronization: %v", err)
	}
	defer distStorage.Stop()
	log.Println("✓ Synchronization started")

	// Ждем начальной синхронизации
	log.Println("Waiting for initial sync...")
	syncTimeout := time.After(10 * time.Second)
	syncTicker := time.NewTicker(200 * time.Millisecond)
	defer syncTicker.Stop()

syncLoop:
	for {
		select {
		case <-syncTicker.C:
			if distStorage.IsSynced() {
				log.Println("✓ Node synced")
				break syncLoop
			}
		case <-syncTimeout:
			log.Println("Warning: Initial sync timeout, continuing anyway...")
			break syncLoop
		}
	}

	// Создаем и запускаем API сервер
	apiServer := NewAPIServer(distStorage, port)

	// Обработка graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := apiServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("✓ API Server is running on port %s", port)
	log.Println("Ready to accept requests!")

	// Ждем сигнала завершения
	<-done
	log.Println("Received shutdown signal")

	// Graceful shutdown с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
