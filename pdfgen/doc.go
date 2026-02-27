package pdfgen

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/go-pdf/fpdf"
)

// DocumentConfig controls page layout and theme for a new document.
type DocumentConfig struct {
	PageSize     string      // "A4" or "Letter"; default "A4"
	Orientation  string      // "portrait" or "landscape"; default "portrait"
	MarginTop    float64     // mm; default 15
	MarginBottom float64     // mm; default 15
	MarginLeft   float64     // mm; default 15
	MarginRight  float64     // mm; default 15
	Theme        ThemeConfig // zero value → DefaultTheme()
}

// Document is the root object that manages the fpdf instance and renders components.
type Document struct {
	pdf        *fpdf.Fpdf
	theme      ThemeConfig
	marginL    float64
	marginR    float64
	marginT    float64
	marginB    float64
	footer     *FooterComponent
	pageWidth  float64 // usable width = page width − left margin − right margin
	err        error   // first component error encountered
	imageCount int     // used to generate unique image names for inline images
}

// New creates a new Document with the given configuration.
func New(cfg DocumentConfig) *Document {
	if cfg.PageSize == "" {
		cfg.PageSize = "A4"
	}
	if cfg.Orientation == "" {
		cfg.Orientation = "portrait"
	}
	if cfg.MarginTop == 0 {
		cfg.MarginTop = 11.3
	}
	if cfg.MarginBottom == 0 {
		cfg.MarginBottom = 15
	}
	if cfg.MarginLeft == 0 {
		cfg.MarginLeft = 11.3
	}
	if cfg.MarginRight == 0 {
		cfg.MarginRight = 11.3
	}

	orientation := "P"
	if strings.EqualFold(cfg.Orientation, "landscape") {
		orientation = "L"
	}

	theme := cfg.Theme
	if theme.DefaultFont.Family == "" {
		theme = DefaultTheme()
	}

	pdf := fpdf.New(orientation, "mm", cfg.PageSize, "")
	pdf.SetMargins(cfg.MarginLeft, cfg.MarginTop, cfg.MarginRight)
	// Disable automatic page breaks; components call newPageIfNeeded themselves.
	pdf.SetAutoPageBreak(false, cfg.MarginBottom)
	// Register {total} as the alias for the total page count.
	pdf.AliasNbPages("{total}")

	w, _ := pdf.GetPageSize()
	d := &Document{
		pdf:       pdf,
		theme:     theme,
		marginL:   cfg.MarginLeft,
		marginR:   cfg.MarginRight,
		marginT:   cfg.MarginTop,
		marginB:   cfg.MarginBottom,
		pageWidth: w - cfg.MarginLeft - cfg.MarginRight,
	}

	pdf.SetFooterFunc(func() {
		if d.footer != nil {
			d.footer.render(d)
		}
	})

	pdf.AddPage()
	return d
}

// Add renders one or more components onto the document in order.
// Returns the document for method chaining.
// On the first component error, subsequent components are skipped and the
// error is surfaced by Save or Bytes.
func (d *Document) Add(components ...Component) *Document {
	for _, c := range components {
		if d.err != nil {
			return d
		}
		if err := c.Render(d); err != nil {
			d.err = fmt.Errorf("pdfgen: %T render: %w", c, err)
		}
	}
	return d
}

// SetFooter registers the footer component that will render automatically on
// every page. Call before adding content for best results.
func (d *Document) SetFooter(f *FooterComponent) *Document {
	d.footer = f
	return d
}

// Save writes the PDF to the given file path.
func (d *Document) Save(path string) error {
	if d.err != nil {
		return d.err
	}
	if err := d.pdf.Error(); err != nil {
		return fmt.Errorf("pdfgen: fpdf internal error: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("pdfgen: create file %q: %w", path, err)
	}
	defer f.Close()
	return d.pdf.Output(f)
}

// Bytes returns the PDF as a byte slice.
func (d *Document) Bytes() ([]byte, error) {
	if d.err != nil {
		return nil, d.err
	}
	if err := d.pdf.Error(); err != nil {
		return nil, fmt.Errorf("pdfgen: fpdf internal error: %w", err)
	}
	var buf bytes.Buffer
	if err := d.pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("pdfgen: output error: %w", err)
	}
	return buf.Bytes(), nil
}

// ── internal helpers used by components ──────────────────────────────────────

func (d *Document) currentY() float64 {
	return d.pdf.GetY()
}

func (d *Document) setY(y float64) {
	d.pdf.SetY(y)
}

func (d *Document) usableWidth() float64 {
	return d.pageWidth
}

// newPageIfNeeded adds a new page when the remaining vertical space is less
// than requiredHeight. Returns true if a new page was added.
func (d *Document) newPageIfNeeded(requiredHeight float64) bool {
	_, pageH := d.pdf.GetPageSize()
	remaining := pageH - d.marginB - d.currentY()
	if remaining < requiredHeight {
		d.pdf.AddPage()
		return true
	}
	return false
}

// applyFont sets the active font, falling back to theme defaults for zero values.
func (d *Document) applyFont(f FontConfig) {
	family := f.Family
	if family == "" {
		family = d.theme.DefaultFont.Family
	}
	size := f.Size
	if size == 0 {
		size = d.theme.DefaultFont.Size
	}
	d.pdf.SetFont(family, f.Style, size)
}

// applyColor sets both the draw color (lines, rect borders) and fill color.
func (d *Document) applyColor(c Color) {
	d.pdf.SetDrawColor(c.R, c.G, c.B)
	d.pdf.SetFillColor(c.R, c.G, c.B)
}

// applyTextColor sets the text rendering color.
func (d *Document) applyTextColor(c Color) {
	d.pdf.SetTextColor(c.R, c.G, c.B)
}
