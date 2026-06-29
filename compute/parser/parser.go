package parser

import (
	"database/compute"
	"errors"
	"strings"
	"unicode"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(query string) (*compute.Command, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, errors.New("empty query")
	}

	tokens, err := p.tokenize(query)
	if err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		return nil, errors.New("no tokens found")
	}

	cmdType := strings.ToUpper(tokens[0])

	switch compute.CommandType(cmdType) {
	case compute.CommandSet:
		if len(tokens) != 3 {
			return nil, errors.New("SET requires exactly 2 arguments")
		}
		return &compute.Command{Type: compute.CommandSet, Args: tokens[1:]}, nil

	case compute.CommandGet:
		if len(tokens) != 2 {
			return nil, errors.New("GET requires exactly 1 argument")
		}
		return &compute.Command{Type: compute.CommandGet, Args: tokens[1:]}, nil

	case compute.CommandDel:
		if len(tokens) != 2 {
			return nil, errors.New("DEL requires exactly 1 argument")
		}
		return &compute.Command{Type: compute.CommandDel, Args: tokens[1:]}, nil

	default:
		return nil, errors.New("unknown command: " + cmdType)
	}
}

func (p *Parser) tokenize(query string) ([]string, error) {
	var tokens []string
	var current strings.Builder

	for _, r := range query {
		if unicode.IsSpace(r) {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		} else if p.isValidChar(r) {
			current.WriteRune(r)
		} else {
			return nil, errors.New("invalid character in query")
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens, nil
}

func (p *Parser) isValidChar(r rune) bool {
	return unicode.IsLetter(r) ||
		unicode.IsDigit(r) ||
		r == '*' || r == '/' || r == '_' ||
		r == '-' || r == '.' || r == ':'
}
