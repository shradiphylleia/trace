package domain
import (
	"testing"
	"time"
)
func TestExpirationTime(t *testing.T) {
	now := time.Date(2026, 6, 16, 12, 0, 0, 0, time.UTC)
	expires, err := ExpirationTime("7d", now)
	if err != nil {
		t.Fatalf("ExpirationTime returned error: %v", err)
	}
	if expires == nil || !expires.Equal(now.Add(7*24*time.Hour)) {
		t.Fatalf("expected 7 day expiration, got %v", expires)
	}

	expires, err = ExpirationTime("never", now)
	if err != nil {
		t.Fatalf("ExpirationTime never returned error: %v", err)
	}
	if expires != nil {
		t.Fatalf("expected nil expiration, got %v", expires)
	}
}
func TestCreateArtifactInputValidate(t *testing.T) {
	input := CreateArtifactInput{
		Title:       "Checkout failure",
		Type:        ArtifactStackTrace,
		ServiceName: "payments",
		Environment: "staging",
		Creator:     "engineer@oracle.com",
		SizeBytes:   10,
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid input: %v", err)
	}
	input.Type = ArtifactType("unknown")
	if err := input.Validate(); err == nil {
		t.Fatal("expected invalid artifact type")
	}
}
