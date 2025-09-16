package services

import (
	"context"
	"testing"
)

func TestGetAPIStatsActiveConnectionsNonNegative(t *testing.T) {
	svc := NewStatsService()
	stats, err := svc.GetAPIStats(context.Background())
	if err != nil {
		t.Fatalf("GetAPIStats returned error: %v", err)
	}

	raw, ok := stats["active_connections"]
	if !ok {
		t.Fatalf("expected active_connections key to be present")
	}

	val, ok := raw.(int)
	if !ok {
		t.Fatalf("expected active_connections to be int, got %T", raw)
	}

	if val < 0 {
		t.Fatalf("expected active_connections to be non-negative, got %d", val)
	}
}
