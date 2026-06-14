package opentrivia

import (
	"testing"
)

// These tests are offline: they exercise the URI driver's pure string functions.

func TestDomainInfo(t *testing.T) {
	info := Domain{}.Info()
	if info.Scheme != "opentrivia" {
		t.Errorf("Scheme = %q, want opentrivia", info.Scheme)
	}
	if len(info.Hosts) == 0 || info.Hosts[0] != Host {
		t.Errorf("Hosts = %v, want [%s]", info.Hosts, Host)
	}
	if info.Identity.Binary != "opentrivia" {
		t.Errorf("Identity.Binary = %q, want opentrivia", info.Identity.Binary)
	}
}

func TestClassify(t *testing.T) {
	typ, id, err := Domain{}.Classify("General Knowledge")
	if err != nil {
		t.Fatalf("Classify error: %v", err)
	}
	if typ != "questions" {
		t.Errorf("typ = %q, want questions", typ)
	}
	if id != "General Knowledge" {
		t.Errorf("id = %q, want General Knowledge", id)
	}
}

func TestClassifyEmpty(t *testing.T) {
	_, _, err := Domain{}.Classify("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestLocate(t *testing.T) {
	got, err := Domain{}.Locate("questions", "foo")
	if err != nil {
		t.Fatalf("Locate error: %v", err)
	}
	if got == "" {
		t.Error("Locate returned empty URL")
	}
}

func TestLocateUnknownType(t *testing.T) {
	_, err := Domain{}.Locate("unknown", "foo")
	if err == nil {
		t.Error("expected error for unknown resource type")
	}
}
