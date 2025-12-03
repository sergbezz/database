package engine

import (
	"testing"

	"go.uber.org/zap"
)

func TestInMemoryEngineSetGet(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	engine := NewInMemoryEngine(logger)

	key := "test_key"
	value := "test_value"

	err := engine.Set(key, value)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	retrieved, err := engine.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved != value {
		t.Errorf("Expected %v, got %v", value, retrieved)
	}
}

func TestInMemoryEngineGetNotFound(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	engine := NewInMemoryEngine(logger)

	_, err := engine.Get("nonexistent")
	if err != ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}
}

func TestInMemoryEngineDelete(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	engine := NewInMemoryEngine(logger)

	key := "test_key"
	value := "test_value"

	engine.Set(key, value)

	err := engine.Delete(key)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = engine.Get(key)
	if err != ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound after delete, got %v", err)
	}
}

func TestInMemoryEngineDeleteNotFound(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	engine := NewInMemoryEngine(logger)

	err := engine.Delete("nonexistent")
	if err != ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}
}

func TestInMemoryEngineOverwrite(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	engine := NewInMemoryEngine(logger)

	key := "test_key"

	engine.Set(key, "value1")
	engine.Set(key, "value2")

	retrieved, _ := engine.Get(key)
	if retrieved != "value2" {
		t.Errorf("Expected value2, got %v", retrieved)
	}
}
