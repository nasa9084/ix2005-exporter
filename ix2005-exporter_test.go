package main

import "testing"

func TestNewExporter(t *testing.T) {
	candidates := []struct {
		uri string
		ok  bool
	}{
		{uri: "192.168.1.1", ok: true},
		{uri: "http://192.168.1.1", ok: true},
		{uri: "https://192.168.1.1", ok: false},
		{uri: "192.168.1.1:1234", ok: true},
		{uri: "http://192.168.1.1:1234", ok: true},
		{uri: "localhost", ok: true},
		{uri: "localhost:1234", ok: true},
		{uri: "http://localhost", ok: true},
		{uri: "http://localhost:1234", ok: true},
		{uri: "https://localhost", ok: false},
		{uri: "http://example.com", ok: true},
		{uri: "foo://bar/baz", ok: false},
	}
	for _, c := range candidates {
		_, err := NewExporter(c.uri)
		if c.ok && err != nil {
			t.Errorf("expected no error w/ %q, but got %q", c.uri, err)
			continue
		}
		if !c.ok && err == nil {
			t.Errorf("expected error w/ %q, but got %q", c.uri, err)
			continue
		}
	}
}
