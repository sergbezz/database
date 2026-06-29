package storage

import (
	"go.uber.org/zap"
)

type Engine interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

type Storage struct {
	engine Engine
	logger *zap.Logger
}

func NewStorage(engine Engine, logger *zap.Logger) *Storage {
	return &Storage{
		engine: engine,
		logger: logger,
	}
}

func (s *Storage) Set(key, value string) error {
	s.logger.Debug("Storage: setting key", zap.String("key", key))
	return s.engine.Set(key, value)
}

func (s *Storage) Get(key string) (string, error) {
	s.logger.Debug("Storage: getting key", zap.String("key", key))
	return s.engine.Get(key)
}

func (s *Storage) Delete(key string) error {
	s.logger.Debug("Storage: deleting key", zap.String("key", key))
	return s.engine.Delete(key)
}
