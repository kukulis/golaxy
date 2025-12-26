package util

import "testing"

func TestNewIndexMapPool(t *testing.T) {
	keys := []string{"alpha", "beta", "gamma"}
	pool := NewIndexMapPool(keys)

	if pool == nil {
		t.Fatal("NewIndexMapPool returned nil")
	}

	if pool.count != 3 {
		t.Errorf("Expected count 3, got %d", pool.count)
	}

	if len(pool.keys) != 3 {
		t.Errorf("Expected keys length 3, got %d", len(pool.keys))
	}

	if len(pool.indexMap) != 3 {
		t.Errorf("Expected indexMap length 3, got %d", len(pool.indexMap))
	}
}

func TestNewIndexMapPoolEmpty(t *testing.T) {
	pool := NewIndexMapPool([]string{})

	if pool.count != 0 {
		t.Errorf("Expected count 0, got %d", pool.count)
	}
}

func TestGetKey(t *testing.T) {
	keys := []string{"alpha", "beta", "gamma"}
	pool := NewIndexMapPool(keys)

	if pool.GetKey(0) != "alpha" {
		t.Errorf("Expected 'alpha', got '%s'", pool.GetKey(0))
	}

	if pool.GetKey(1) != "beta" {
		t.Errorf("Expected 'beta', got '%s'", pool.GetKey(1))
	}

	if pool.GetKey(2) != "gamma" {
		t.Errorf("Expected 'gamma', got '%s'", pool.GetKey(2))
	}
}

func TestGetIndex(t *testing.T) {
	keys := []string{"alpha", "beta", "gamma"}
	pool := NewIndexMapPool(keys)

	if pool.GetIndex("alpha") != 0 {
		t.Errorf("Expected index 0, got %d", pool.GetIndex("alpha"))
	}

	if pool.GetIndex("beta") != 1 {
		t.Errorf("Expected index 1, got %d", pool.GetIndex("beta"))
	}

	if pool.GetIndex("gamma") != 2 {
		t.Errorf("Expected index 2, got %d", pool.GetIndex("gamma"))
	}
}

func TestGetIndexNonExistent(t *testing.T) {
	keys := []string{"alpha", "beta"}
	pool := NewIndexMapPool(keys)

	// Non-existent key should return 0 (default value for int)
	if pool.GetIndex("nonexistent") != 0 {
		t.Errorf("Expected index 0 for non-existent key, got %d", pool.GetIndex("nonexistent"))
	}
}

func TestRemoveKeyMiddle(t *testing.T) {
	keys := []string{"alpha", "beta", "gamma"}
	pool := NewIndexMapPool(keys)

	pool.RemoveKey("beta")

	// After removing "beta", "gamma" should be moved to index 1
	if pool.GetKey(1) != "gamma" {
		t.Errorf("Expected 'gamma' at index 1, got '%s'", pool.GetKey(1))
	}

	// The index of "gamma" should now be 1
	if pool.GetIndex("gamma") != 1 {
		t.Errorf("Expected 'gamma' to be at index 1, got %d", pool.GetIndex("gamma"))
	}

	// "alpha" should still be at index 0
	if pool.GetKey(0) != "alpha" {
		t.Errorf("Expected 'alpha' at index 0, got '%s'", pool.GetKey(0))
	}

	// Count should be decremented
	if pool.count != 2 {
		t.Errorf("Expected count 2 after removal, got %d", pool.count)
	}

	// The removed key should no longer be in the index map or should have invalid index
	idx := pool.GetIndex("beta")
	if idx >= 0 && idx < pool.count {
		t.Errorf("Removed key 'beta' should not have valid index, got %d", idx)
	}
}

func TestRemoveKeyFirst(t *testing.T) {
	keys := []string{"alpha", "beta", "gamma"}
	pool := NewIndexMapPool(keys)

	pool.RemoveKey("alpha")

	// After removing "alpha", "gamma" should be moved to index 0
	if pool.GetKey(0) != "gamma" {
		t.Errorf("Expected 'gamma' at index 0, got '%s'", pool.GetKey(0))
	}

	// The index of "gamma" should now be 0
	if pool.GetIndex("gamma") != 0 {
		t.Errorf("Expected 'gamma' to be at index 0, got %d", pool.GetIndex("gamma"))
	}

	// Count should be decremented
	if pool.count != 2 {
		t.Errorf("Expected count 2 after removal, got %d", pool.count)
	}
}

func TestRemoveKeyLast(t *testing.T) {
	keys := []string{"alpha", "beta", "gamma"}
	pool := NewIndexMapPool(keys)

	pool.RemoveKey("gamma")

	// After removing "gamma", elements at 0 and 1 should remain unchanged
	if pool.GetKey(0) != "alpha" {
		t.Errorf("Expected 'alpha' at index 0, got '%s'", pool.GetKey(0))
	}

	if pool.GetKey(1) != "beta" {
		t.Errorf("Expected 'beta' at index 1, got '%s'", pool.GetKey(1))
	}

	// Count should be decremented
	if pool.count != 2 {
		t.Errorf("Expected count 2 after removal, got %d", pool.count)
	}
}

func TestRemoveMultipleKeys(t *testing.T) {
	keys := []string{"alpha", "beta", "gamma", "delta", "epsilon"}
	pool := NewIndexMapPool(keys)

	// Remove "beta" (index 1)
	pool.RemoveKey("beta")
	if pool.count != 4 {
		t.Errorf("Expected count 4 after first removal, got %d", pool.count)
	}

	// Remove "delta" (which might have moved)
	pool.RemoveKey("delta")
	if pool.count != 3 {
		t.Errorf("Expected count 3 after second removal, got %d", pool.count)
	}

	// Remove "alpha"
	pool.RemoveKey("alpha")
	if pool.count != 2 {
		t.Errorf("Expected count 2 after third removal, got %d", pool.count)
	}

	// Verify remaining keys can still be accessed correctly
	// We should still be able to get valid keys at indices 0 and 1
	key0 := pool.GetKey(0)
	key1 := pool.GetKey(1)

	if key0 == "" || key1 == "" {
		t.Errorf("Keys at valid indices should not be empty: key0='%s', key1='%s'", key0, key1)
	}

	// Verify indices are correct for remaining keys
	if pool.GetIndex(key0) != 0 {
		t.Errorf("Key '%s' should be at index 0, got %d", key0, pool.GetIndex(key0))
	}

	if pool.GetIndex(key1) != 1 {
		t.Errorf("Key '%s' should be at index 1, got %d", key1, pool.GetIndex(key1))
	}
}

func TestRemoveKeySingleElement(t *testing.T) {
	keys := []string{"alpha"}
	pool := NewIndexMapPool(keys)

	pool.RemoveKey("alpha")

	if pool.count != 0 {
		t.Errorf("Expected count 0 after removing single element, got %d", pool.count)
	}
}
