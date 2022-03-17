package asset

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/pkg/test"
	"regexp"
	"testing"
	"time"
)

var (
	filterRegexp = regexp.MustCompile("^(avg/(imgs|bg))|(gamedata/(excel|levels/obt/main))")
)

func TestGet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	SetChatMask([]byte("UITpAi82pHAWwnzqHRMCwPonJLIB3WCl"))

	dir := t.TempDir()
	err := Get(ctx, "21-12-31-15-44-39-814f71", dir, filterRegexp)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	expected := "h1:8g0nB4HSnlGiv7taFGo0gY9wZJJcINzi64VGp1ICO/Q="
	actual := test.HashDir(dir)
	if expected != actual {
		t.Fatalf("expected: %v; actual: %v", expected, actual)
	}
}

func TestGetUpdate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	SetChatMask([]byte("UITpAi82pHAWwnzqHRMCwPonJLIB3WCl"))

	dirGet := t.TempDir()
	err := Get(ctx, "21-12-31-15-44-39-814f71", dirGet, filterRegexp)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	dirUpdate := t.TempDir()
	err = Get(ctx, "21-12-01-03-53-27-2e01ea", dirUpdate, filterRegexp)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	err = Update(ctx, "21-12-01-03-53-27-2e01ea", "21-12-31-15-44-39-814f71", dirUpdate, filterRegexp)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	err = test.AssertAllIn(dirGet, dirUpdate)
	if err != nil {
		t.Fatalf("%+v", err)
	}
}
