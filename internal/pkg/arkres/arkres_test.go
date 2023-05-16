package arkres

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/flandiayingman/arkwaifu/internal/pkg/test"
)

var (
	filterRegexp = regexp.MustCompile("^(avg/(imgs|bg))|(gamedata/(excel|levels/obt/main))")
)

func TestGet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	dir := t.TempDir()
	err := GetFromHGAPI(ctx, "21-12-31-15-44-39-814f71", dir, filterRegexp)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	expected := "h1:lrJBULywG33Q3Rfu7i1f/KjojSnvG+HsfkW5wnK0CLk="
	actual := test.HashDir(dir)
	if expected != actual {
		t.Fatalf("expected: %v; actual: %v", expected, actual)
	}
}

func TestGetUpdate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	dirGet := t.TempDir()
	err := GetFromHGAPI(ctx, "21-12-31-15-44-39-814f71", dirGet, filterRegexp)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	dirUpdate := t.TempDir()
	err = GetFromHGAPI(ctx, "21-12-01-03-53-27-2e01ea", dirUpdate, filterRegexp)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	err = GetFromHGAPIIncrementally(ctx, "21-12-01-03-53-27-2e01ea", "21-12-31-15-44-39-814f71", dirUpdate, filterRegexp)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	err = test.AssertAllIn(dirGet, dirUpdate)
	if err != nil {
		t.Fatalf("%+v", err)
	}
}
