package parser

import (
	"testing"
)

func TestParseSetCommand(t *testing.T) {
	p := NewParser()

	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{"valid SET", "SET key value", false},
		{"valid SET with underscores", "SET weather_2_pm cold_moscow_weather", false},
		{"valid SET with special chars", "SET user_**** data123", false},
		{"invalid SET - no args", "SET", true},
		{"invalid SET - one arg", "SET key", true},
		{"invalid SET - too many args", "SET key value extra", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := p.Parse(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && cmd.Type != CommandSet {
				t.Errorf("Expected SET command, got %v", cmd.Type)
			}
		})
	}
}

func TestParseGetCommand(t *testing.T) {
	p := NewParser()

	tests := []struct {
		name    string
		query   string
		wantErr bool
		wantKey string
	}{
		{"valid GET", "GET mykey", false, "mykey"},
		{"valid GET with path", "GET /etc/nginx/config", false, "/etc/nginx/config"},
		{"invalid GET - no arg", "GET", true, ""},
		{"invalid GET - too many args", "GET key extra", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := p.Parse(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if cmd.Type != CommandGet {
					t.Errorf("Expected GET command, got %v", cmd.Type)
				}
				if cmd.Args[0] != tt.wantKey {
					t.Errorf("Expected key %v, got %v", tt.wantKey, cmd.Args[0])
				}
			}
		})
	}
}

func TestParseDelCommand(t *testing.T) {
	p := NewParser()

	cmd, err := p.Parse("DEL user_****")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if cmd.Type != CommandDel {
		t.Errorf("Expected DEL command, got %v", cmd.Type)
	}

	if len(cmd.Args) != 1 || cmd.Args[0] != "user_****" {
		t.Errorf("Expected args [user_****], got %v", cmd.Args)
	}
}

func TestParseCaseInsensitive(t *testing.T) {
	p := NewParser()

	tests := []string{"set key value", "Set key value", "SET key value"}

	for _, query := range tests {
		cmd, err := p.Parse(query)
		if err != nil {
			t.Errorf("Parse(%q) unexpected error: %v", query, err)
		}
		if cmd.Type != CommandSet {
			t.Errorf("Parse(%q) expected SET, got %v", query, cmd.Type)
		}
	}
}

func TestParseInvalidCharacters(t *testing.T) {
	p := NewParser()

	_, err := p.Parse("SET key@ value")
	if err == nil {
		t.Error("Expected error for invalid character @")
	}
}
