// internal/clipboard/classifier.go
package clipboard

import (
	"encoding/json"
	"net"
	"net/url"
	"regexp"
	"strings"

	"clipboard-manager/internal/store"
)

type ClassifyResult struct {
	Type    string
	Subtype string
}

var (
	reEmail    = regexp.MustCompile(`(?i)^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	reUUID     = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	reHexColor = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
	reCSSColor = regexp.MustCompile(`(?i)^(rgb|rgba|hsl|hsla|oklch)\(`)
	reFilePath = regexp.MustCompile(`^(~/|/)`)
	rePhone    = regexp.MustCompile(`^[\+\d\s\-\(\)\.]{7,20}$`)
	reJWT      = regexp.MustCompile(`^[A-Za-z0-9\-_]+\.[A-Za-z0-9\-_]+\.[A-Za-z0-9\-_]+$`)
	reBase64   = regexp.MustCompile(`^[A-Za-z0-9+/]+=*$`)
	reMD5      = regexp.MustCompile(`^[0-9a-fA-F]{32}$`)
	reSHA1     = regexp.MustCompile(`^[0-9a-fA-F]{40}$`)
	reSHA256   = regexp.MustCompile(`^[0-9a-fA-F]{64}$`)

	sqlKeywords    = []string{"SELECT ", "INSERT INTO", "UPDATE ", "DELETE FROM", "CREATE TABLE", "DROP TABLE", "ALTER TABLE"}
	apiKeyPrefixes = []string{"sk-", "pk_", "ghp_", "ghs_", "xoxb-", "xoxp-", "AKIA", "rk_live_", "rk_test_"}
)

// utiTypes maps NSPasteboard UTI strings directly to a ContentType.
var utiTypes = map[string]string{
	"public.png":                        store.TypeImage,
	"public.jpeg":                       store.TypeImage,
	"public.tiff":                       store.TypeImage,
	"com.compuserve.gif":                store.TypeImage,
	"public.heic":                       store.TypeImage,
	"public.bmp":                        store.TypeImage,
	"com.adobe.pdf":                     store.TypePDF,
	"public.rtf":                        store.TypeRTF,
	"public.html":                       store.TypeHTML,
	"com.apple.cocoa.pasteboard.color":  store.TypeNSColor,
	"public.file-url":                   store.TypeFileRef,
}

func Classify(uti, content string) ClassifyResult {
	// Tier 1 — UTI
	if t, ok := utiTypes[uti]; ok {
		return ClassifyResult{Type: t}
	}

	s := strings.TrimSpace(content)
	if s == "" {
		return ClassifyResult{Type: store.TypeText}
	}

	// Tier 2 — regex / pattern
	if isURL(s) {
		return ClassifyResult{Type: store.TypeURL}
	}
	if reEmail.MatchString(s) {
		return ClassifyResult{Type: store.TypeEmail}
	}
	if net.ParseIP(s) != nil {
		return ClassifyResult{Type: store.TypeIP}
	}
	if reHexColor.MatchString(s) || reCSSColor.MatchString(s) {
		return ClassifyResult{Type: store.TypeColorCode}
	}
	if reUUID.MatchString(s) {
		return ClassifyResult{Type: store.TypeUUID}
	}
	if reFilePath.MatchString(s) && !strings.Contains(s, "\n") && len(s) < 512 {
		return ClassifyResult{Type: store.TypeFilePath}
	}
	if reMD5.MatchString(s) || reSHA1.MatchString(s) || reSHA256.MatchString(s) {
		return ClassifyResult{Type: store.TypeHash}
	}
	if rePhone.MatchString(s) {
		return ClassifyResult{Type: store.TypePhone}
	}

	// Tier 3 — heuristics
	if isJSON(s) {
		return ClassifyResult{Type: store.TypeJSON}
	}
	if strings.HasPrefix(s, "<?xml") || strings.HasPrefix(s, "<svg") {
		return ClassifyResult{Type: store.TypeXML}
	}
	if isSQL(s) {
		return ClassifyResult{Type: store.TypeSQL}
	}
	if reJWT.MatchString(s) && strings.Count(s, ".") == 2 {
		return ClassifyResult{Type: store.TypeJWT}
	}
	if strings.HasPrefix(s, "-----BEGIN ") {
		return ClassifyResult{Type: store.TypeSSHKey}
	}
	if isAPIKey(s) {
		return ClassifyResult{Type: store.TypeAPIKey}
	}
	if isYAML(s) {
		return ClassifyResult{Type: store.TypeYAML}
	}
	if isBase64(s) {
		return ClassifyResult{Type: store.TypeBase64}
	}
	if isMarkdown(s) {
		return ClassifyResult{Type: store.TypeMarkdown}
	}
	if lang := detectCode(s); lang != "" {
		return ClassifyResult{Type: store.TypeCode, Subtype: lang}
	}

	return ClassifyResult{Type: store.TypeText}
}

func isURL(s string) bool {
	for _, prefix := range []string{"http://", "https://", "ftp://", "ssh://", "git://"} {
		if strings.HasPrefix(s, prefix) {
			u, err := url.ParseRequestURI(s)
			return err == nil && u.Host != ""
		}
	}
	return false
}

func isJSON(s string) bool {
	if len(s) < 2 || (s[0] != '{' && s[0] != '[') {
		return false
	}
	var v interface{}
	return json.Unmarshal([]byte(s), &v) == nil
}

func isSQL(s string) bool {
	u := strings.ToUpper(s)
	for _, kw := range sqlKeywords {
		if strings.Contains(u, kw) {
			return true
		}
	}
	return false
}

func isAPIKey(s string) bool {
	if len(s) < 16 || strings.ContainsAny(s, " \n\t") {
		return false
	}
	for _, prefix := range apiKeyPrefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

func isYAML(s string) bool {
	lines := strings.Split(s, "\n")
	if len(lines) < 2 {
		return false
	}
	count := 0
	for _, l := range lines {
		t := strings.TrimSpace(l)
		if strings.Contains(t, ": ") || (strings.HasSuffix(t, ":") && !strings.HasPrefix(t, "#")) {
			count++
		}
	}
	return count >= 2
}

func isBase64(s string) bool {
	return len(s) >= 32 && len(s)%4 == 0 &&
		!strings.ContainsAny(s, " \n\t") && reBase64.MatchString(s)
}

func isMarkdown(s string) bool {
	score := 0
	for _, sig := range []string{"# ", "**", "- [", "](", "```", "---", "*"} {
		if strings.Contains(s, sig) {
			score++
		}
	}
	return score >= 2
}

func detectCode(s string) string {
	switch {
	case strings.Contains(s, "func ") && (strings.Contains(s, "package ") || strings.Contains(s, " error")):
		return "go"
	// Go: also detect func with fmt.Println (common snippet without package keyword)
	case strings.Contains(s, "func ") && strings.Contains(s, "fmt."):
		return "go"
	case strings.Contains(s, "def ") && strings.Contains(s, ":\n"):
		return "python"
	case strings.HasPrefix(s, "#!/bin/bash") || strings.HasPrefix(s, "#!/usr/bin/env bash"):
		return "shell"
	// TypeScript: interface keyword is a strong signal; also catch const/let with arrow or import
	case strings.Contains(s, "interface ") && (strings.Contains(s, "const ") || strings.Contains(s, "let ") || strings.Contains(s, "type ")):
		return "typescript"
	case (strings.Contains(s, "const ") || strings.Contains(s, "let ")) &&
		(strings.Contains(s, "=>") || strings.Contains(s, "import ")):
		return "typescript"
	case strings.Contains(s, "fn ") && strings.Contains(s, "->"):
		return "rust"
	default:
		lines := strings.Split(s, "\n")
		if len(lines) >= 3 {
			symbols := strings.Count(s, "{") + strings.Count(s, "}") +
				strings.Count(s, "(") + strings.Count(s, ")")
			if symbols >= 4 {
				return "code"
			}
		}
		return ""
	}
}
