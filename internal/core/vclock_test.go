package core

import (
	"testing"
)

func TestVectorClock_Increment(t *testing.T) {
	vc := NewVectorClock()

	v1 := vc.Increment("node1")
	if v1 != 1 {
		t.Errorf("Expected 1, got %d", v1)
	}

	v2 := vc.Increment("node1")
	if v2 != 2 {
		t.Errorf("Expected 2, got %d", v2)
	}

	v3 := vc.Increment("node2")
	if v3 != 1 {
		t.Errorf("Expected 1, got %d", v3)
	}
}

func TestVectorClock_HappensBefore(t *testing.T) {
	// Создаём два вектора: vc1 < vc2
	vc1 := NewVectorClock()
	vc1.Increment("node1") // {node1: 1}

	vc2 := NewVectorClock()
	vc2.Increment("node1") // {node1: 1}
	vc2.Increment("node1") // {node1: 2}

	if !vc1.HappensBefore(vc2) {
		t.Error("vc1 should happen before vc2")
	}

	if vc2.HappensBefore(vc1) {
		t.Error("vc2 should not happen before vc1")
	}
}

func TestVectorClock_ConcurrentWith(t *testing.T) {
	// Создаём два конкурентных вектора
	vc1 := NewVectorClock()
	vc1.Increment("node1") // {node1: 1}

	vc2 := NewVectorClock()
	vc2.Increment("node2") // {node2: 1}

	if !vc1.ConcurrentWith(vc2) {
		t.Error("vc1 and vc2 should be concurrent")
	}

	if !vc2.ConcurrentWith(vc1) {
		t.Error("vc2 and vc1 should be concurrent")
	}
}

func TestVectorClock_Update(t *testing.T) {
	vc1 := NewVectorClock()
	vc1.Increment("node1") // {node1: 1}

	vc2 := NewVectorClock()
	vc2.Increment("node1") // {node1: 1}
	vc2.Increment("node1") // {node1: 2}
	vc2.Increment("node2") // {node1: 2, node2: 1}

	vc1.Update(vc2)

	if vc1.Get("node1") != 2 {
		t.Errorf("Expected node1=2, got %d", vc1.Get("node1"))
	}

	if vc1.Get("node2") != 1 {
		t.Errorf("Expected node2=1, got %d", vc1.Get("node2"))
	}
}

func TestVectorClock_Equals(t *testing.T) {
	vc1 := NewVectorClock()
	vc1.Increment("node1")
	vc1.Increment("node2")

	vc2 := NewVectorClock()
	vc2.Increment("node1")
	vc2.Increment("node2")

	if !vc1.Equals(vc2) {
		t.Error("vc1 should equal vc2")
	}

	vc2.Increment("node2")
	if vc1.Equals(vc2) {
		t.Error("vc1 should not equal vc2 after increment")
	}
}

func TestVectorClock_Clone(t *testing.T) {
	vc1 := NewVectorClock()
	vc1.Increment("node1")
	vc1.Increment("node2")

	vc2 := vc1.Clone()

	if !vc1.Equals(vc2) {
		t.Error("Clone should be equal to original")
	}

	vc2.Increment("node2")
	if vc1.Equals(vc2) {
		t.Error("Modifying clone should not affect original")
	}
}

func TestVectorClock_IsSafeToRead(t *testing.T) {
	vc := NewVectorClock()

	// Нет подтверждений
	if vc.IsSafeToRead(2, 3) {
		t.Error("Should not be safe with 0 acks")
	}

	// Одно подтверждение
	vc.Increment("node1")
	if vc.IsSafeToRead(2, 3) {
		t.Error("Should not be safe with 1 ack, need 2")
	}

	// Два подтверждения - quorum
	vc.Increment("node2")
	if !vc.IsSafeToRead(2, 3) {
		t.Error("Should be safe with 2 acks (quorum)")
	}

	// Три подтверждения
	vc.Increment("node3")
	if !vc.IsSafeToRead(2, 3) {
		t.Error("Should be safe with 3 acks")
	}
}

func TestVectorClock_MergeMax(t *testing.T) {
	vc1 := NewVectorClock()
	vc1.Increment("node1") // {node1: 1}
	vc1.Increment("node2") // {node1: 1, node2: 1}

	vc2 := NewVectorClock()
	vc2.Increment("node1") // {node1: 1}
	vc2.Increment("node1") // {node1: 2}
	vc2.Increment("node3") // {node1: 2, node3: 1}

	merged := vc1.MergeMax(vc2)

	if merged.Get("node1") != 2 {
		t.Errorf("Expected node1=2, got %d", merged.Get("node1"))
	}
	if merged.Get("node2") != 1 {
		t.Errorf("Expected node2=1, got %d", merged.Get("node2"))
	}
	if merged.Get("node3") != 1 {
		t.Errorf("Expected node3=1, got %d", merged.Get("node3"))
	}

	// Оригиналы не должны измениться
	if vc1.Get("node1") != 1 {
		t.Error("Original vc1 should not be modified")
	}
	if vc2.Get("node2") != 0 {
		t.Error("Original vc2 should not be modified")
	}
}

func TestVectorClock_Concurrent(t *testing.T) {
	// Тест конкурентной работы
	vc := NewVectorClock()

	done := make(chan bool)

	// 10 горутин инкрементируют разные узлы
	for i := 0; i < 10; i++ {
		go func(nodeID string) {
			for j := 0; j < 100; j++ {
				vc.Increment(nodeID)
			}
			done <- true
		}(string(rune('A' + i)))
	}

	// Ждём завершения
	for i := 0; i < 10; i++ {
		<-done
	}

	// Проверяем, что все инкременты применились
	for i := 0; i < 10; i++ {
		nodeID := string(rune('A' + i))
		if vc.Get(nodeID) != 100 {
			t.Errorf("Expected 100 for %s, got %d", nodeID, vc.Get(nodeID))
		}
	}
}
