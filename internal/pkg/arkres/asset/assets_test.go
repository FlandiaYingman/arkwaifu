package asset

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/pkg/test"
	"path/filepath"
	"regexp"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	dir := t.TempDir()
	err := Get(ctx, "21-12-31-15-44-39-814f71", dir, regexp.MustCompile("^avg/(imgs|bg)"))
	if err != nil {
		t.Fatalf("%+v", err)
	}

	expected := "h1:R6xBu0bIzs4Kx1hKAO63f6bku9PTdoX4XQAs6K99ZZo="
	actual := test.HashDir(dir)
	if expected != actual {
		t.Fatalf("expected: %v; actual: %v", expected, actual)
	}
}

func TestGetUpdate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var err error
	resRegexp := regexp.MustCompile("^avg/(imgs|bg)")

	dirGet := filepath.Join(t.TempDir(), "Get")
	err = Get(ctx, "21-12-31-15-44-39-814f71", dirGet, resRegexp)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	dirUpdate := filepath.Join(t.TempDir(), "Update")
	err = Get(ctx, "21-12-01-03-53-27-2e01ea", dirUpdate, resRegexp)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	err = Update(ctx, "21-12-01-03-53-27-2e01ea", "21-12-31-15-44-39-814f71", dirUpdate, resRegexp)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	expected := test.HashDir(dirGet)
	actual := test.HashDir(dirUpdate)
	if expected != actual {
		t.Fatalf("expected: %v; actual: %v", expected, actual)
	}
}
