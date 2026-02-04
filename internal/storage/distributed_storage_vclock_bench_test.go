package storage_test

import (
	"fmt"
	"petacore/internal/core"
	"petacore/internal/storage"
	"sync/atomic"
	"testing"
	"time"
)

// Unified, idiomatic benchmarks for VClock-based distributed storage.
// These are intended to be clear and reproducible for plotting with plots.py.

func benchSetup(b *testing.B, nodes int) *storage.DistributedStorageVClock {
	return SetupDistributedStorageVClock(b, "node1", nodes, core.ReadCommitted, 0)
}

func BenchmarkVClock_SingleWrite(b *testing.B) {
	ds := benchSetup(b, 1)
	if ds == nil {
		b.Skip("setup failed")
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key%d", i))
		ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			tx.Write(key, "value")
			return nil
		})
	}
}

func BenchmarkVClock_SingleRead(b *testing.B) {
	ds := benchSetup(b, 1)
	if ds == nil {
		b.Skip("setup failed")
	}

	const prepopulate = 1000
	for i := 0; i < prepopulate; i++ {
		key := []byte(fmt.Sprintf("key%d", i))
		ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			tx.Write(key, "value")
			return nil
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key%d", i%prepopulate))
		ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			tx.Read(key)
			return nil
		})
	}
}

func BenchmarkVClock_ReadWriteMix(b *testing.B) {
	ds := benchSetup(b, 1)
	if ds == nil {
		b.Skip("setup failed")
	}

	const keys = 100
	for i := 0; i < keys; i++ {
		k := []byte(fmt.Sprintf("key%d", i))
		ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			tx.Write(k, "init")
			return nil
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := []byte(fmt.Sprintf("key%d", i%keys))
		if i%2 == 0 {
			ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
				tx.Read(k)
				return nil
			})
		} else {
			ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
				tx.Write(k, "val")
				return nil
			})
		}
	}
}

func BenchmarkVClock_ConcurrentWrites(b *testing.B) {
	ds := benchSetup(b, 1)
	if ds == nil {
		b.Skip("setup failed")
	}

	var ctr uint64
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := atomic.AddUint64(&ctr, 1) - 1
			key := []byte(fmt.Sprintf("key%d", id))
			ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
				tx.Write(key, "value")
				return nil
			})
		}
	})
}

func BenchmarkVClock_ConcurrentReads(b *testing.B) {
	ds := benchSetup(b, 1)
	if ds == nil {
		b.Skip("setup failed")
	}

	const prepopulate = 1000
	for i := 0; i < prepopulate; i++ {
		key := []byte(fmt.Sprintf("key%d", i))
		ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			tx.Write(key, "value")
			return nil
		})
	}

	var idx uint64
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := atomic.AddUint64(&idx, 1) - 1
			key := []byte(fmt.Sprintf("key%d", i%prepopulate))
			ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
				tx.Read(key)
				return nil
			})
		}
	})
}

func BenchmarkVClock_Transaction(b *testing.B) {
	ds := benchSetup(b, 1)
	if ds == nil {
		b.Skip("setup failed")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := []byte(fmt.Sprintf("key%d", i))
		ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			tx.Write(k, "value")
			tx.Read(k)
			return nil
		})
	}
}

func BenchmarkVClock_LargeTransaction(b *testing.B) {
	ds := benchSetup(b, 1)
	if ds == nil {
		b.Skip("setup failed")
	}

	const opsPerTx = 10
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			for j := 0; j < opsPerTx; j++ {
				key := []byte(fmt.Sprintf("key%d_%d", i, j))
				tx.Write(key, "value")
			}
			return nil
		})
	}
}

func BenchmarkVClock_HotKey(b *testing.B) {
	ds := benchSetup(b, 1)
	if ds == nil {
		b.Skip("setup failed")
	}

	ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		tx.Write([]byte("hotkey"), "0")
		return nil
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
				tx.Read([]byte("hotkey"))
				tx.Write([]byte("hotkey"), "updated")
				return nil
			})
		}
	})
}

func BenchmarkVClock_QuorumOverhead(b *testing.B) {
	scenarios := []struct {
		name  string
		nodes int
	}{{"1node", 1}, {"3nodes", 3}, {"5nodes", 5}}

	for _, s := range scenarios {
		b.Run(s.name, func(b *testing.B) {
			ds := benchSetup(b, s.nodes)
			if ds == nil {
				b.Skip("setup failed")
			}

			const pre = 100
			for i := 0; i < pre; i++ {
				k := []byte(fmt.Sprintf("key%d", i))
				ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
					tx.Write(k, "v")
					return nil
				})
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := []byte(fmt.Sprintf("key%d", i%pre))
				ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
					tx.Read(k)
					return nil
				})
			}
		})
	}
}

func BenchmarkVClock_VectorClockSize(b *testing.B) {
	// measure cost when vector clock has more nodes
	ds := benchSetup(b, 5)
	if ds == nil {
		b.Skip("setup failed")
	}

	var ctr uint64
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := atomic.AddUint64(&ctr, 1) - 1
			key := []byte(fmt.Sprintf("key%d", id))
			ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
				tx.Write(key, "v")
				return nil
			})
		}
	})
}

func BenchmarkVClock_LatencyAvg(b *testing.B) {
	ds := benchSetup(b, 1)
	if ds == nil {
		b.Skip("setup failed")
	}

	// sample latencies (avoid storing huge slices for large b.N)
	var total time.Duration
	var sampleCount int64
	const maxSamples = 10000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := time.Now()
		ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			tx.Write([]byte(fmt.Sprintf("key%d", i)), "v")
			return nil
		})
		if sampleCount < maxSamples {
			total += time.Since(start)
			sampleCount++
		}
	}
	b.StopTimer()

	if sampleCount > 0 {
		avg := total / time.Duration(sampleCount)
		b.ReportMetric(float64(avg.Microseconds()), "avg_latency_us")
	}
}

func BenchmarkVClock_Throughput(b *testing.B) {
	ds := benchSetup(b, 1)
	if ds == nil {
		b.Skip("setup failed")
	}

	var ops int64
	start := time.Now()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
				tx.Write([]byte("tput_key"), "v")
				return nil
			})
			atomic.AddInt64(&ops, 1)
		}
	})
	elapsed := time.Since(start)
	if elapsed > 0 {
		b.ReportMetric(float64(ops)/elapsed.Seconds(), "ops/s")
	}
}

func BenchmarkVClock_MemoryUsage(b *testing.B) {
	ds := benchSetup(b, 1)
	if ds == nil {
		b.Skip("setup failed")
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key%d", i))
		val := fmt.Sprintf("value_%d_with_some_data", i)
		ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			tx.Write(key, val)
			return nil
		})
	}
}

func BenchmarkVClock_Scalability(b *testing.B) {
	goroutines := []int{1, 2, 4, 8, 16}
	for _, g := range goroutines {
		b.Run(fmt.Sprintf("goroutines_%d", g), func(b *testing.B) {
			ds := benchSetup(b, 1)
			if ds == nil {
				b.Skip("setup failed")
			}

			b.SetParallelism(g)
			var idx uint64
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					i := atomic.AddUint64(&idx, 1) - 1
					k := []byte(fmt.Sprintf("key%d", i))
					ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
						tx.Write(k, "v")
						return nil
					})
				}
			})
		})
	}
}

func BenchmarkVClock_BatchWrite(b *testing.B) {
	batchSizes := []int{1, 10, 50, 100}

	for _, bs := range batchSizes {
		b.Run(fmt.Sprintf("batch_%d", bs), func(b *testing.B) {
			ds := benchSetup(b, 1)
			if ds == nil {
				b.Skip("setup failed")
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
					for j := 0; j < bs; j++ {
						tx.Write([]byte(fmt.Sprintf("key%d_%d", i, j)), "v")
					}
					return nil
				})
			}
		})
	}
}

func BenchmarkVClock_ContentionLevel(b *testing.B) {
	contentionLevels := []struct {
		name     string
		keyRange int
	}{
		{"low_contention", 10000},
		{"medium_contention", 100},
		{"high_contention", 10},
	}

	for _, level := range contentionLevels {
		b.Run(level.name, func(b *testing.B) {
			ds := benchSetup(b, 1)
			if ds == nil {
				b.Skip("setup failed")
			}

			var idx uint64
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					i := atomic.AddUint64(&idx, 1) - 1
					key := []byte(fmt.Sprintf("key%d", i%uint64(level.keyRange)))
					ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
						tx.Write(key, "v")
						return nil
					})
				}
			})
		})
	}
}

func BenchmarkVClock_SynchronizationOverhead(b *testing.B) {
	ds := benchSetup(b, 1)
	if ds == nil {
		b.Skip("setup failed")
	}

	// Measure synchronization overhead by running parallel transactions
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
				tx.Write([]byte("sync_key"), "v")
				return nil
			})
		}
	})
}
