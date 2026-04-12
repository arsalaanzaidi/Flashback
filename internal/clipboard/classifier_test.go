// internal/clipboard/classifier_test.go
package clipboard_test

import (
	"testing"
	"clipboard-manager/internal/clipboard"
	"clipboard-manager/internal/store"
)

func classify(uti, content string) string {
	return clipboard.Classify(uti, content).Type
}

func TestClassify_UTITypes(t *testing.T) {
	cases := []struct{ uti, want string }{
		{"public.png",                          store.TypeImage},
		{"public.tiff",                         store.TypeImage},
		{"com.adobe.pdf",                       store.TypePDF},
		{"public.rtf",                          store.TypeRTF},
		{"public.html",                         store.TypeHTML},
		{"com.apple.cocoa.pasteboard.color",    store.TypeNSColor},
		{"public.file-url",                     store.TypeFileRef},
	}
	for _, c := range cases {
		if got := classify(c.uti, ""); got != c.want {
			t.Errorf("UTI %q: want %q, got %q", c.uti, c.want, got)
		}
	}
}

func TestClassify_Tier2(t *testing.T) {
	cases := []struct{ content, want string }{
		{"https://github.com/wailsapp/wails",   store.TypeURL},
		{"arsalaan@servicenow.com",              store.TypeEmail},
		{"192.168.1.1",                          store.TypeIP},
		{"2001:db8::1",                          store.TypeIP},
		{"#7c3aed",                              store.TypeColorCode},
		{"rgb(124, 58, 237)",                    store.TypeColorCode},
		{"550e8400-e29b-41d4-a716-446655440000", store.TypeUUID},
		{"/Users/arsalaan/.ssh/id_rsa",          store.TypeFilePath},
		{"~/Projects/squawk/main.go",            store.TypeFilePath},
		{"d41d8cd98f00b204e9800998ecf8427e",     store.TypeHash},
		{"da39a3ee5e6b4b0d3255bfef95601890afd80709", store.TypeHash},
	}
	for _, c := range cases {
		if got := classify("public.utf8-plain-text", c.content); got != c.want {
			t.Errorf("content %q: want %q, got %q", c.content, c.want, got)
		}
	}
}

func TestClassify_Tier3(t *testing.T) {
	cases := []struct{ content, want string }{
		{`{"key": "value", "num": 42}`,  store.TypeJSON},
		{"<?xml version=\"1.0\"?>",       store.TypeXML},
		{"SELECT id FROM items WHERE pinned = 0", store.TypeSQL},
		{"eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1c2VyIn0.signature", store.TypeJWT},
		{"-----BEGIN RSA PRIVATE KEY-----", store.TypeSSHKey},
		{"sk-ant-api03-xK9mP2abc123",      store.TypeAPIKey},
		{"ghp_16C7e42F292c6912E7710c838347Dd9032650",  store.TypeAPIKey},
		{"# Heading\n**bold** and [link](url)\n```code```", store.TypeMarkdown},
	}
	for _, c := range cases {
		if got := classify("public.utf8-plain-text", c.content); got != c.want {
			t.Errorf("content %q: want %q, got %q", c.content[:minInt(30, len(c.content))], c.want, got)
		}
	}
}

func TestClassify_CodeLanguages(t *testing.T) {
	cases := []struct{ content, subtype string }{
		{"func main() {\n\tfmt.Println(\"hi\")\n}", "go"},
		{"def hello():\n    print('hi')\n    return True", "python"},
		{"#!/bin/bash\necho hello", "shell"},
		{"const x = 42\ninterface Foo {\n  bar: string\n}", "typescript"},
		{"fn add(a: i32, b: i32) -> i32 {\n    a + b\n}", "rust"},
	}
	for _, c := range cases {
		r := clipboard.Classify("public.utf8-plain-text", c.content)
		if r.Type != store.TypeCode {
			t.Errorf("expected CODE for %q, got %q", c.content[:20], r.Type)
		}
		if r.Subtype != c.subtype {
			t.Errorf("expected subtype %q, got %q", c.subtype, r.Subtype)
		}
	}
}

func TestClassify_FallsBackToText(t *testing.T) {
	if got := classify("public.utf8-plain-text", "just a plain sentence"); got != store.TypeText {
		t.Errorf("expected TEXT, got %q", got)
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
