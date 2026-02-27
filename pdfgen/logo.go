package pdfgen

import (
	"bytes"
	"fmt"

	"github.com/go-pdf/fpdf"
)

// LogoComponent renders an image at a fixed page position without advancing
// the Y cursor, so it floats above the normal content flow.
//
// Set either ImagePath (path to a PNG/JPG file on disk) or ImageData (raw bytes).
// Position values: "top-left", "top-right", "top-center".
type LogoComponent struct {
	ImagePath string  // path to PNG or JPG file on disk
	ImageData []byte  // raw PNG/JPG bytes (alternative to ImagePath)
	Width     float64 // mm
	Height    float64 // mm; 0 = auto-preserve aspect ratio
	Position  string  // "top-left" | "top-right" | "top-center"
	OffsetX   float64 // additional X offset in mm
	OffsetY   float64 // additional Y offset in mm
}

// Render places the image at the configured position.
// The Y cursor is NOT advanced — logos float above the layout flow.
func (l *LogoComponent) Render(doc *Document) error {
	if l.ImagePath == "" && len(l.ImageData) == 0 {
		return fmt.Errorf("pdfgen: LogoComponent requires ImagePath or ImageData")
	}
	if l.Width <= 0 {
		return fmt.Errorf("pdfgen: LogoComponent Width must be > 0")
	}

	pdf := doc.pdf
	pageW, _ := pdf.GetPageSize()

	var x float64
	switch l.Position {
	case "top-right":
		x = pageW - doc.marginR - l.Width + l.OffsetX
	case "top-center":
		x = (pageW-l.Width)/2 + l.OffsetX
	default: // "top-left"
		x = doc.marginL + l.OffsetX
	}
	y := doc.marginT + l.OffsetY

	// Save current Y so we can restore it after ImageOptions (which may move cursor).
	savedY := doc.currentY()

	imgName := l.ImagePath
	opts := fpdf.ImageOptions{}

	if len(l.ImageData) > 0 {
		doc.imageCount++
		imgName = fmt.Sprintf("pdfgen_img_%d", doc.imageCount)
		imgType := detectImageType(l.ImageData)
		opts.ImageType = imgType
		reader := bytes.NewReader(l.ImageData)
		pdf.RegisterImageOptionsReader(imgName, opts, reader)
	}

	pdf.ImageOptions(imgName, x, y, l.Width, l.Height, false, opts, 0, "")

	// Restore Y — logos do not participate in the content flow.
	doc.setY(savedY)
	return nil
}

// detectImageType returns "JPG" or "PNG" based on the file magic bytes.
func detectImageType(data []byte) string {
	if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xD8 {
		return "JPG"
	}
	return "PNG"
}
