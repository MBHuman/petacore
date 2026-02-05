package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVectorClock_Increment(t *testing.T) {
	vc := NewVectorClock()

	v1 := vc.Increment("node1")
	require.Equal(t, uint64(1), v1, "First increment should return 1")

	v2 := vc.Increment("node1")
	require.Equal(t, uint64(2), v2, "Second increment should return 2")

	v3 := vc.Increment("node2")
	require.Equal(t, uint64(1), v3, "First increment for node2 should return 1")
}

func TestVectorClock_HappensBefore(t *testing.T) {
	// Create two vectors: vc1 < vc2
	vc1 := NewVectorClock()
	vc1.Increment("node1") // {node1: 1}

	vc2 := NewVectorClock()
	vc2.Increment("node1") // {node1: 1}
	vc2.Increment("node1") // {node1: 2}

	assert.True(t, vc1.HappensBefore(vc2), "vc1 should happen before vc2")
	assert.False(t, vc2.HappensBefore(vc1), "vc2 should not happen before vc1")
}

func TestVectorClock_ConcurrentWith(t *testing.T) {
	// Create two concurrent vectors
	vc1 := NewVectorClock()
	vc1.Increment("node1") // {node1: 1}

	vc2 := NewVectorClock()
	vc2.Increment("node2") // {node2: 1}

	assert.True(t, vc1.ConcurrentWith(vc2), "vc1 and vc2 should be concurrent")
	assert.True(t, vc2.ConcurrentWith(vc1), "vc2 and vc1 should be concurrent")
}

func TestVectorClock_Update(t *testing.T) {
	vc1 := NewVectorClock()
	vc1.Increment("node1") // {node1: 1}

	vc2 := NewVectorClock()
	vc2.Increment("node1") // {node1: 1}
	vc2.Increment("node1") // {node1: 2}
	vc2.Increment("node2") // {node1: 2, node2: 1}

	vc1.Update(vc2)

	assert.Equal(t, uint64(2), vc1.Get("node1"), "node1 should be 2 after update")
	assert.Equal(t, uint64(1), vc1.Get("node2"), "node2 should be 1 after update")
}

func TestVectorClock_Equals(t *testing.T) {
	vc1 := NewVectorClock()
	vc1.Increment("node1")
	vc1.Increment("node2")

	vc2 := NewVectorClock()
	vc2.Increment("node1")
	vc2.Increment("node2")

	assert.True(t, vc1.Equals(vc2), "vc1 should equal vc2")

	vc2.Increment("node2")
	assert.False(t, vc1.Equals(vc2), "vc1 should not equal vc2 after increment")
}

func TestVectorClock_Clone(t *testing.T) {
	vc1 := NewVectorClock()
	vc1.Increment("node1")
	vc1.Increment("node2")

	vc2 := vc1.Clone()

	assert.True(t, vc1.Equals(vc2), "Clone should be equal to original")

	vc2.Increment("node2")
	assert.False(t, vc1.Equals(vc2), "Modifying clone should not affect original")
}

func TestVectorClock_IsSafeToRead(t *testing.T) {
	vc := NewVectorClock()

	// No acknowledgments
	assert.False(t, vc.IsSafeToRead(2, 3, ""), "Should not be safe with 0 acks")

	// One acknowledgment
	vc.Increment("node1")
	assert.False(t, vc.IsSafeToRead(2, 3, ""), "Should not be safe with 1 ack, need 2")

	// Two acknowledgments - quorum
	vc.Increment("node2")
	assert.True(t, vc.IsSafeToRead(2, 3, ""), "Should be safe with 2 acks (quorum)")

	// Three acknowledgments
	vc.Increment("node3")
	assert.True(t, vc.IsSafeToRead(2, 3, ""), "Should be safe with 3 acks")
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

	assert.Equal(t, uint64(2), merged.Get("node1"), "node1 should be max of both (2)")
	assert.Equal(t, uint64(1), merged.Get("node2"), "node2 should be 1 (from vc1)")
	assert.Equal(t, uint64(1), merged.Get("node3"), "node3 should be 1 (from vc2)")

	// Originals should not be modified
	assert.Equal(t, uint64(1), vc1.Get("node1"), "Original vc1 should not be modified")
	assert.Equal(t, uint64(0), vc2.Get("node2"), "Original vc2 should not be modified")
}

func TestVectorClock_Concurrent(t *testing.T) {
	// Concurrent access test
	vc := NewVectorClock()

	done := make(chan bool)

	// 10 goroutines incrementing different nodes
	for i := 0; i < 10; i++ {
		go func(nodeID string) {
			for j := 0; j < 100; j++ {
				vc.Increment(nodeID)
			}
			done <- true
		}(string(rune('A' + i)))
	}

	// Wait for completion
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all increments were applied
	for i := 0; i < 10; i++ {
		nodeID := string(rune('A' + i))
		assert.Equal(t, uint64(100), vc.Get(nodeID), "Node %s should have 100 increments", nodeID)
	}
}
