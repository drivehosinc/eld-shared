package pdfgen

// HeaderComponent renders the document title block at the top-left of the page.
// It intentionally uses ~65% of the usable width so a top-right logo has room.
//
// Default font sizes are derived from the Lucid ELD HTML source:
//   Title:    16pt bold
//   Subtitle: 10pt bold
//   Lines:     8pt regular
type HeaderComponent struct {
	Title         string     // large bold text, e.g. "IFTA REPORT"
	Subtitle      string     // medium bold, e.g. "QGM EXPRESS"
	Lines         []string   // additional detail lines (address, date range, etc.)
	TitleFont     FontConfig // zero value → 16pt bold
	SubtitleFont  FontConfig // zero value → 10pt bold
	LineFont      FontConfig // zero value → 8pt regular
	SubtitleColor Color      // zero value → theme AccentColor
}

// Render draws the header and advances the Y cursor past the rendered block.
func (h *HeaderComponent) Render(doc *Document) error {
	titleFont := h.TitleFont
	if titleFont.Family == "" {
		titleFont = FontConfig{
			Family: doc.theme.DefaultFont.Family,
			Size:   16,
			Style:  "B",
		}
	}
	subtitleFont := h.SubtitleFont
	if subtitleFont.Family == "" {
		subtitleFont = FontConfig{
			Family: doc.theme.DefaultFont.Family,
			Size:   10,
			Style:  "B",
		}
	}
	lineFont := h.LineFont
	if lineFont.Family == "" {
		lineFont = FontConfig{
			Family: doc.theme.DefaultFont.Family,
			Size:   8,
			Style:  "",
		}
	}

	// Use 65% of usable width to leave space for a top-right logo.
	contentW := doc.usableWidth() * 0.65
	titleLineH := 8.5  // pt 16 → ~5.6mm + leading
	bodyLineH  := 4.9  // pt 10 → ~3.5mm + leading

	startY := doc.currentY()
	x := doc.marginL

	// Title
	if h.Title != "" {
		doc.applyFont(titleFont)
		doc.applyTextColor(doc.theme.PrimaryText)
		doc.pdf.SetXY(x, startY)
		doc.pdf.CellFormat(contentW, titleLineH, h.Title, "", 1, "L", false, 0, "")
	}

	// Subtitle
	if h.Subtitle != "" {
		subtitleColor := h.SubtitleColor
		if subtitleColor.R == 0 && subtitleColor.G == 0 && subtitleColor.B == 0 {
			subtitleColor = doc.theme.AccentColor
		}
		doc.applyFont(subtitleFont)
		doc.applyTextColor(subtitleColor)
		doc.pdf.SetX(x)
		doc.pdf.CellFormat(contentW, bodyLineH, h.Subtitle, "", 1, "L", false, 0, "")
	}

	// Additional lines (date range, address, etc.)
	if len(h.Lines) > 0 {
		doc.applyFont(lineFont)
		doc.applyTextColor(doc.theme.SecondaryText)
		for _, line := range h.Lines {
			doc.pdf.SetX(x)
			doc.pdf.CellFormat(contentW, bodyLineH-0.5, line, "", 1, "L", false, 0, "")
		}
	}

	// Small bottom gap
	doc.setY(doc.currentY() + 3)
	return nil
}
