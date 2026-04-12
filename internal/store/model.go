// internal/store/model.go
package store

// Item represents one clipboard entry, serialised to JSON for the frontend.
type Item struct {
	ID          string `json:"id"`
	Content     string `json:"content"`      // raw text; empty for binary items
	ContentHash string `json:"contentHash"`  // SHA256 (text) or pHash (image)
	Type        string `json:"type"`         // see ContentType constants below
	Subtype     string `json:"subtype"`      // language for CODE, ext for IMAGE
	Pinned      bool   `json:"pinned"`
	CopiedAt    int64  `json:"copiedAt"`     // unix ms — updated on re-copy
	CreatedAt   int64  `json:"createdAt"`    // unix ms — set on first capture
	CharCount   int    `json:"charCount"`
	ImagePath   string `json:"imagePath"`    // full image on disk (images only)
	ThumbBase64 string `json:"thumbBase64"`  // data:image/png;base64,... (images only)
	ThumbBlob   []byte `json:"-"`            // raw bytes stored in SQLite, not sent to frontend
}

// Settings is persisted as a single JSON row in the settings table.
type Settings struct {
	RetentionMode    string `json:"retentionMode"`    // "unlimited" | "count" | "days"
	RetentionValue   int    `json:"retentionValue"`   // N items or N days
	GlobalShortcut   string `json:"globalShortcut"`   // default: "option+space"
	LaunchAtLogin    bool   `json:"launchAtLogin"`
	FirstRunComplete bool   `json:"firstRunComplete"`
}

// Content type constants — kept as plain strings so they serialise cleanly to JSON.
const (
	TypeURL       = "URL"
	TypeEmail     = "EMAIL"
	TypeImage     = "IMAGE"
	TypePDF       = "PDF"
	TypeRTF       = "RICH_TEXT"
	TypeHTML      = "HTML"
	TypeNSColor   = "COLOR"
	TypeColorCode = "COLOR_CODE"
	TypeFileRef   = "FILE_REF"
	TypeIP        = "IP"
	TypeUUID      = "UUID"
	TypeFilePath  = "FILE_PATH"
	TypeHash      = "HASH"
	TypePhone     = "PHONE"
	TypeJSON      = "JSON"
	TypeXML       = "XML"
	TypeYAML      = "YAML"
	TypeSQL       = "SQL"
	TypeJWT       = "JWT"
	TypeSSHKey    = "SSH_KEY"
	TypeAPIKey    = "API_KEY"
	TypeBase64    = "BASE64"
	TypeMarkdown  = "MARKDOWN"
	TypeCode      = "CODE"
	TypeText      = "TEXT"
)

func DefaultSettings() Settings {
	return Settings{
		RetentionMode:  "count",
		RetentionValue: 1000,
		GlobalShortcut: "option+space",
		LaunchAtLogin:  false,
	}
}
