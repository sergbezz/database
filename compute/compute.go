package compute

import (
	"database/compute/parser"
	"errors"

	"go.uber.org/zap"
)

type StorageLayer interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

type Compute struct {
	storage StorageLayer
	parser  *parser.Parser
	logger  *zap.Logger
}

func NewCompute(storage StorageLayer, logger *zap.Logger) *Compute {
	return &Compute{
		storage: storage,
		parser:  parser.NewParser(),
		logger:  logger,
	}
}

func (c *Compute) Execute(query string) (string, error) {
	c.logger.Info("Executing query", zap.String("query", query))

	cmd, err := c.parser.Parse(query)
	if err != nil {
		c.logger.Error("Failed to parse query", zap.Error(err))
		return "", err
	}

	c.logger.Debug("Parsed command", zap.String("type", string(cmd.Type)))

	switch cmd.Type {
	case parser.CommandSet:
		if len(cmd.Args) != 2 {
			return "", errors.New("SET requires exactly 2 arguments")
		}
		err := c.storage.Set(cmd.Args[0], cmd.Args[1])
		if err != nil {
			c.logger.Error("SET failed", zap.Error(err))
			return "", err
		}
		c.logger.Info("SET successful", zap.String("key", cmd.Args[0]))
		return "OK", nil

	case parser.CommandGet:
		if len(cmd.Args) != 1 {
			return "", errors.New("GET requires exactly 1 argument")
		}
		value, err := c.storage.Get(cmd.Args[0])
		if err != nil {
			c.logger.Error("GET failed", zap.Error(err))
			return "", err
		}
		c.logger.Info("GET successful", zap.String("key", cmd.Args[0]))
		return value, nil

	case parser.CommandDel:
		if len(cmd.Args) != 1 {
			return "", errors.New("DEL requires exactly 1 argument")
		}
		err := c.storage.Delete(cmd.Args[0])
		if err != nil {
			c.logger.Error("DEL failed", zap.Error(err))
			return "", err
		}
		c.logger.Info("DEL successful", zap.String("key", cmd.Args[0]))
		return "OK", nil

	default:
		return "", errors.New("unknown command")
	}
}
