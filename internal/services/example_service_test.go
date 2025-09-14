package services

import (
	"context"
	"testing"
	"time"
)

func TestExampleService_Echo(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
		wantErr bool
	}{
		{
			name:    "simple message",
			message: "hello",
			want:    "hello",
			wantErr: false,
		},
		{
			name:    "empty message",
			message: "",
			want:    "",
			wantErr: false,
		},
		{
			name:    "unicode message",
			message: "你好 世界 🌍",
			want:    "你好 世界 🌍",
			wantErr: false,
		},
	}

	svc := NewExampleService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := svc.Echo(ctx, tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Echo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Echo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExampleService_Echo_ContextCancellation(t *testing.T) {
	svc := NewExampleService()

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := svc.Echo(ctx, "test")
	if err == nil {
		t.Error("Echo() expected error for cancelled context")
	}
}

func TestExampleService_GetStatus(t *testing.T) {
	svc := NewExampleService()
	ctx := context.Background()

	status, err := svc.GetStatus(ctx)
	if err != nil {
		t.Errorf("GetStatus() unexpected error: %v", err)
	}
	if status != "ok" {
		t.Errorf("GetStatus() = %v, want 'ok'", status)
	}
}

func BenchmarkExampleService_Echo(b *testing.B) {
	svc := NewExampleService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.Echo(ctx, "benchmark test message")
	}
}

func TestExampleService_Echo_Timeout(t *testing.T) {
	svc := NewExampleService()

	// Test with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Should complete before timeout
	result, err := svc.Echo(ctx, "quick test")
	if err != nil {
		t.Errorf("Echo() unexpected error: %v", err)
	}
	if result != "quick test" {
		t.Errorf("Echo() = %v, want 'quick test'", result)
	}
}