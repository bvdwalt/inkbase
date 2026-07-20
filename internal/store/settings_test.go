package store_test

import (
	"context"
	"testing"
)

func TestGetOrCreateAPIKeyIsStableAndGenerated(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)

	key1, err := s.GetOrCreateAPIKey(ctx)
	if err != nil {
		t.Fatalf("GetOrCreateAPIKey: %v", err)
	}
	if len(key1) != 64 {
		t.Errorf("len(key1) = %d, want 64", len(key1))
	}

	key2, err := s.GetOrCreateAPIKey(ctx)
	if err != nil {
		t.Fatalf("GetOrCreateAPIKey: %v", err)
	}
	if key1 != key2 {
		t.Errorf("key1 = %q, key2 = %q, want same key on repeated calls", key1, key2)
	}
}

func TestRegenerateAPIKeyChangesValue(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)

	original, err := s.GetOrCreateAPIKey(ctx)
	if err != nil {
		t.Fatalf("GetOrCreateAPIKey: %v", err)
	}

	regenerated, err := s.RegenerateAPIKey(ctx)
	if err != nil {
		t.Fatalf("RegenerateAPIKey: %v", err)
	}
	if regenerated == original {
		t.Errorf("regenerated key equals original: %q", regenerated)
	}
	if len(regenerated) != 64 {
		t.Errorf("len(regenerated) = %d, want 64", len(regenerated))
	}

	current, err := s.GetOrCreateAPIKey(ctx)
	if err != nil {
		t.Fatalf("GetOrCreateAPIKey after regenerate: %v", err)
	}
	if current != regenerated {
		t.Errorf("current = %q, want %q", current, regenerated)
	}
}
