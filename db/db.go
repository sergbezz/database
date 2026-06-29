package db

import (
	"database/compute"
	"errors"

	"go.uber.org/zap"
)

type computeLayer interface {
	Parse(query string) (*compute.Command, error)
}

type storageLayer interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

type Database struct {
	compute computeLayer
	storage storageLayer
	logger  *zap.Logger
}

func NewDatabase(compute computeLayer, storage storageLayer, logger *zap.Logger) (*Database, error) {
	if compute == nil {
		return nil, errors.New("compute is invalid")
	}
	if storage == nil {
		return nil, errors.New("storage is invalid")
	}
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}
	return &Database{
		compute: compute,
		storage: storage,
		logger:  logger,
	}, nil
}

func (d *Database) HandleQuery(query string) (string, error) {
	d.logger.Info("handling query", zap.String("query", query))

	cmd, err := d.compute.Parse(query)
	if err != nil {
		return "", err
	}

	switch cmd.Type {
	case compute.CommandSet:
		if err := d.storage.Set(cmd.Args[0], cmd.Args[1]); err != nil {
			d.logger.Error("SET failed", zap.Error(err))
			return "", err
		}
		d.logger.Info("SET successful", zap.String("key", cmd.Args[0]))
		return "OK", nil

	case compute.CommandGet:
		value, err := d.storage.Get(cmd.Args[0])
		if err != nil {
			d.logger.Error("GET failed", zap.Error(err))
			return "", err
		}
		d.logger.Info("GET successful", zap.String("key", cmd.Args[0]))
		return value, nil

	case compute.CommandDel:
		if err := d.storage.Delete(cmd.Args[0]); err != nil {
			d.logger.Error("DEL failed", zap.Error(err))
			return "", err
		}
		d.logger.Info("DEL successful", zap.String("key", cmd.Args[0]))
		return "OK", nil

	default:
		return "", errors.New("unknown command")
	}
}
