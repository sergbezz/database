package engine

import (
	"errors"

	"go.uber.org/zap"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type InMemoryEngine struct {
	data   map[string]string
	logger *zap.Logger
}

func NewInMemoryEngine(logger *zap.Logger) *InMemoryEngine {
	logger.Info("Initializing in-memory engine")
	return &InMemoryEngine{
		data:   make(map[string]string),
		logger: logger,
	}
}

func (e *InMemoryEngine) Set(key, value string) error {

	e.data[key] = value
	e.logger.Debug("Engine: key set",
		zap.String("key", key),
		zap.Int("total_keys", len(e.data)))

	return nil
}

func (e *InMemoryEngine) Get(key string) (string, error) {

	value, exists := e.data[key]
	if !exists {
		e.logger.Debug("Engine: key not found", zap.String("key", key))
		return "", ErrKeyNotFound
	}

	e.logger.Debug("Engine: key retrieved", zap.String("key", key))
	return value, nil
}

func (e *InMemoryEngine) Delete(key string) error {

	if _, exists := e.data[key]; !exists {
		e.logger.Debug("Engine: key not found for deletion", zap.String("key", key))
		return ErrKeyNotFound
	}

	delete(e.data, key)
	e.logger.Debug("Engine: key deleted",
		zap.String("key", key),
		zap.Int("total_keys", len(e.data)))

	return nil
}
