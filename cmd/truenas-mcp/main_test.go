package main

import (
	"reflect"
	"testing"
)

func TestParseDomains(t *testing.T) {
	t.Run("empty input returns nil", func(t *testing.T) {
		got, err := parseDomains("")
		if err != nil {
			t.Fatalf("parseDomains: %v", err)
		}
		if got != nil {
			t.Errorf("got %v, want nil", got)
		}
	})

	t.Run("single domain", func(t *testing.T) {
		got, err := parseDomains("iscsi")
		if err != nil {
			t.Fatalf("parseDomains: %v", err)
		}
		want := []string{"iscsi"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("multiple domains with whitespace", func(t *testing.T) {
		got, err := parseDomains("iscsi, nvmeof , pool")
		if err != nil {
			t.Fatalf("parseDomains: %v", err)
		}
		want := []string{"iscsi", "nvmeof", "pool"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("unknown domain returns error", func(t *testing.T) {
		_, err := parseDomains("iscsi,not-a-real-domain")
		if err == nil {
			t.Fatal("expected error for unknown domain, got nil")
		}
	})

	t.Run("blank entries are skipped", func(t *testing.T) {
		got, err := parseDomains("iscsi,,pool,")
		if err != nil {
			t.Fatalf("parseDomains: %v", err)
		}
		want := []string{"iscsi", "pool"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
