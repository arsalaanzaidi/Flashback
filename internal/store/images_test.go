// internal/store/images_test.go
package store_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func makePNG(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x), G: uint8(y), B: 128, A: 255})
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

func TestSaveImage_CreatesFileAndThumb(t *testing.T) {
	s := openTestStore(t)
	data := makePNG(400, 300)

	path, thumb, phash, err := s.SaveImage(data, "png")
	if err != nil {
		t.Fatalf("SaveImage: %v", err)
	}
	if path == "" {
		t.Fatal("expected non-empty path")
	}
	if len(thumb) == 0 {
		t.Fatal("expected non-empty thumbnail blob")
	}
	if phash == "" {
		t.Fatal("expected non-empty phash")
	}

	// Thumbnail must decode as valid PNG
	_, err = png.Decode(bytes.NewReader(thumb))
	if err != nil {
		t.Fatalf("thumbnail is not valid PNG: %v", err)
	}
}

func TestSaveImage_SameImageSamePhash(t *testing.T) {
	s := openTestStore(t)
	data := makePNG(400, 300)
	_, _, phash1, _ := s.SaveImage(data, "png")
	_, _, phash2, _ := s.SaveImage(data, "png")
	if phash1 != phash2 {
		t.Fatal("same image should produce same phash")
	}
}
