// internal/clipboard/classifier.go
package clipboard

import (
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
)

// utiTypes maps NSPasteboard UTI strings directly to a ContentType.
var utiTypes = map[string]string{
	"public.png":                       store.TypeImage,
	"public.jpeg":                      store.TypeImage,
	"public.tiff":                      store.TypeImage,
	"com.compuserve.gif":               store.TypeImage,
	"public.heic":                      store.TypeImage,
	"public.bmp":                       store.TypeImage,
	"com.adobe.pdf":                    store.TypePDF,
	"public.rtf":                       store.TypeRTF,
	"public.html":                      store.TypeHTML,
	"com.apple.cocoa.pasteboard.color": store.TypeNSColor,
	"public.file-url":                  store.TypeFileRef,
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

	// Tier 2 — narrow regex checks. Only patterns with strong, unambiguous
	// signals live here; heuristic-based content sniffing was removed because
	// it produced too many false positives.
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
