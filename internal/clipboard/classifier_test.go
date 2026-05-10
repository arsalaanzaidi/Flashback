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
		{"public.png", store.TypeImage},
		{"public.tiff", store.TypeImage},
		{"com.adobe.pdf", store.TypePDF},
		{"public.rtf", store.TypeRTF},
		{"public.html", store.TypeHTML},
		{"com.apple.cocoa.pasteboard.color", store.TypeNSColor},
		{"public.file-url", store.TypeFileRef},
	}
	for _, c := range cases {
		if got := classify(c.uti, ""); got != c.want {
			t.Errorf("UTI %q: want %q, got %q", c.uti, c.want, got)
		}
	}
}

func TestClassify_RegexTypes(t *testing.T) {
	cases := []struct{ content, want string }{
		{"https://github.com/wailsapp/wails", store.TypeURL},
		{"arsalaan@example.com", store.TypeEmail},
		{"192.168.1.1", store.TypeIP},
		{"2001:db8::1", store.TypeIP},
		{"#7c3aed", store.TypeColorCode},
		{"rgb(124, 58, 237)", store.TypeColorCode},
		{"550e8400-e29b-41d4-a716-446655440000", store.TypeUUID},
		{"/Users/arsalaan/.ssh/id_rsa", store.TypeFilePath},
		{"~/Projects/squawk/main.go", store.TypeFilePath},
	}
	for _, c := range cases {
		if got := classify("public.utf8-plain-text", c.content); got != c.want {
			t.Errorf("content %q: want %q, got %q", c.content, c.want, got)
		}
	}
}

func TestClassify_FallsBackToText(t *testing.T) {
	cases := []string{
		"just a plain sentence",
		"REQ0292850",
		"v1.0.3",
		"2.19.1",
		"abc.def.ghi",
		"1234567890",
		`{"key": "value"}`,
		"SELECT * FROM users",
		"# Heading",
		"key: value\nfoo: bar",
		"d41d8cd98f00b204e9800998ecf8427e",
	}
	for _, content := range cases {
		if got := classify("public.utf8-plain-text", content); got != store.TypeText {
			t.Errorf("expected TEXT for %q, got %q", content, got)
		}
	}
}
