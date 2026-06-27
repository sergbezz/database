package compute

import "go.uber.org/zap"

type CommandType string

const (
	CommandSet CommandType = "SET"
	CommandGet CommandType = "GET"
	CommandDel CommandType = "DEL"
)

type Command struct {
	Type CommandType
	Args []string
}

type queryParser interface {
	Parse(query string) (*Command, error)
}

type Compute struct {
	parser queryParser
	logger *zap.Logger
}

func NewCompute(parser queryParser, logger *zap.Logger) *Compute {
	return &Compute{
		parser: parser,
		logger: logger,
	}
}

func (c *Compute) Parse(query string) (*Command, error) {
	c.logger.Debug("parsing query", zap.String("query", query))

	cmd, err := c.parser.Parse(query)
	if err != nil {
		c.logger.Error("failed to parse query", zap.Error(err))
		return nil, err
	}

	c.logger.Debug("parsed command", zap.String("type", string(cmd.Type)))
	return cmd, nil
}
