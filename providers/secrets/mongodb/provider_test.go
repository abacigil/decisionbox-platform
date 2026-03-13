package mongodb

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/decisionbox-io/decisionbox/libs/go-common/secrets"
)

func TestEncryptDecrypt(t *testing.T) {
	// Generate a random 32-byte key
	keyBytes := make([]byte, 32)
	rand.Read(keyBytes)
	encKey := base64.StdEncoding.EncodeToString(keyBytes)

	p, err := NewMongoProvider(nil, "test", encKey)
	if err != nil {
		t.Fatal(err)
	}

	plaintext := "sk-ant-api03-very-secret-key-12345"
	encrypted, err := p.encrypt(plaintext)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	if encrypted == plaintext {
		t.Error("encrypted should differ from plaintext")
	}

	decrypted, err := p.decrypt(encrypted)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("decrypted = %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptDecrypt_DifferentCiphertexts(t *testing.T) {
	keyBytes := make([]byte, 32)
	rand.Read(keyBytes)
	encKey := base64.StdEncoding.EncodeToString(keyBytes)

	p, _ := NewMongoProvider(nil, "test", encKey)

	// Same plaintext should produce different ciphertexts (random nonce)
	enc1, _ := p.encrypt("same-secret")
	enc2, _ := p.encrypt("same-secret")

	if enc1 == enc2 {
		t.Error("same plaintext should produce different ciphertexts (random nonce)")
	}

	// Both should decrypt to the same value
	dec1, _ := p.decrypt(enc1)
	dec2, _ := p.decrypt(enc2)
	if dec1 != dec2 || dec1 != "same-secret" {
		t.Errorf("decrypted values differ: %q vs %q", dec1, dec2)
	}
}

func TestNewMongoProvider_InvalidKeyLength(t *testing.T) {
	shortKey := base64.StdEncoding.EncodeToString([]byte("tooshort"))
	_, err := NewMongoProvider(nil, "test", shortKey)
	if err == nil {
		t.Error("expected error for short key")
	}
}

func TestNewMongoProvider_InvalidBase64(t *testing.T) {
	_, err := NewMongoProvider(nil, "test", "not-valid-base64!!!")
	if err == nil {
		t.Error("expected error for invalid base64")
	}
}

func TestNewMongoProvider_NoEncryption(t *testing.T) {
	p, err := NewMongoProvider(nil, "test", "")
	if err != nil {
		t.Fatal(err)
	}
	if p.gcm != nil {
		t.Error("gcm should be nil when no encryption key")
	}
}

func TestNewMongoProvider_DefaultNamespace(t *testing.T) {
	p, _ := NewMongoProvider(nil, "", "")
	if p.namespace != "decisionbox" {
		t.Errorf("namespace = %q, want decisionbox", p.namespace)
	}
}

func TestSecretDoc_Fields(t *testing.T) {
	doc := secretDoc{
		Namespace: "decisionbox",
		ProjectID: "proj-123",
		Key:       "llm-api-key",
		Value:     "encrypted-value",
		Encrypted: true,
	}
	if doc.Namespace != "decisionbox" || doc.Key != "llm-api-key" {
		t.Error("fields not set correctly")
	}
}

// --- Interface compliance ---

func TestMongoProvider_ImplementsInterface(t *testing.T) {
	var _ secrets.Provider = (*MongoProvider)(nil)
}
