// internal/store/images.go
package store

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	_ "image/gif"
	_ "image/jpeg"
	"os"
	"path/filepath"

	"github.com/corona10/goimagehash"
	"golang.org/x/image/draw"
)

// SaveImage writes a full image to disk and returns:
//   - path: absolute path to the saved file
//   - thumbBlob: 128x128 PNG bytes for SQLite storage
//   - phash: perceptual hash string for deduplication
func (s *Store) SaveImage(data []byte, ext string) (path string, thumbBlob []byte, phash string, err error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", nil, "", fmt.Errorf("images: decode: %w", err)
	}

	// Perceptual hash (64-bit) for deduplication
	h, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return "", nil, "", fmt.Errorf("images: phash: %w", err)
	}
	phash = fmt.Sprintf("p%016x", h.GetHash())

	// 128×128 thumbnail
	thumb := image.NewRGBA(image.Rect(0, 0, 128, 128))
	draw.CatmullRom.Scale(thumb, thumb.Bounds(), img, img.Bounds(), draw.Over, nil)
	var buf bytes.Buffer
	if err = png.Encode(&buf, thumb); err != nil {
		return "", nil, "", fmt.Errorf("images: thumb encode: %w", err)
	}
	thumbBlob = buf.Bytes()

	// Persist full image — filename is the phash so duplicates overwrite cleanly
	path = filepath.Join(s.imageDir, phash+"."+ext)
	if err = os.WriteFile(path, data, 0644); err != nil {
		return "", nil, "", fmt.Errorf("images: write: %w", err)
	}

	return path, thumbBlob, phash, nil
}

// DeleteImage removes the full image file from disk. Safe to call with "".
func (s *Store) DeleteImage(path string) error {
	if path == "" {
		return nil
	}
	err := os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
