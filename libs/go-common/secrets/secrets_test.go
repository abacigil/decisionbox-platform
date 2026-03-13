package secrets

import "testing"

func TestMaskValue(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"sk-ant-api03-abc123def456", "sk-ant***f456"},
		{"short", "***"},
		{"exactly10!", "***"},
		{"12345678901", "123456***8901"},
		{"", "***"},
	}
	for _, tt := range tests {
		got := MaskValue(tt.input)
		if got != tt.want {
			t.Errorf("MaskValue(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	cfg := LoadConfig()
	if cfg.Provider != "mongodb" {
		t.Errorf("default provider = %q, want mongodb", cfg.Provider)
	}
	if cfg.Namespace != "decisionbox" {
		t.Errorf("default namespace = %q, want decisionbox", cfg.Namespace)
	}
}

func TestErrNotFound(t *testing.T) {
	if ErrNotFound == nil {
		t.Error("ErrNotFound should not be nil")
	}
	if ErrNotFound.Error() != "secret not found" {
		t.Errorf("ErrNotFound = %q", ErrNotFound.Error())
	}
}

func TestRegisterAndList(t *testing.T) {
	// Register a test provider
	Register("test-secrets", func(cfg Config) (Provider, error) {
		return nil, nil
	}, ProviderMeta{Name: "Test Provider", Description: "for testing"})

	providers := RegisteredProviders()
	found := false
	for _, p := range providers {
		if p == "test-secrets" {
			found = true
		}
	}
	if !found {
		t.Error("test-secrets not found in registered providers")
	}

	metas := RegisteredProvidersMeta()
	found = false
	for _, m := range metas {
		if m.ID == "test-secrets" {
			found = true
			if m.Name != "Test Provider" {
				t.Errorf("name = %q", m.Name)
			}
		}
	}
	if !found {
		t.Error("test-secrets meta not found")
	}
}

func TestNewProvider_Unknown(t *testing.T) {
	_, err := NewProvider(Config{Provider: "nonexistent"})
	if err == nil {
		t.Error("expected error for unknown provider")
	}
}
