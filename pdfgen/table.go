package pdfgen

// OverflowMode controls how cell text is handled when it exceeds the column width.
type OverflowMode int

const (
	// OverflowWrap wraps text across multiple lines, making the row taller.
	OverflowWrap OverflowMode = iota
	// OverflowTruncate clips the text and appends "…".
	OverflowTruncate
)

// ColumnDef defines a single table column.
type ColumnDef struct {
	Header      string
	Width       float64      // mm; 0 = column shares remaining space equally
	Align       string       // "L", "C", "R"
	Overflow    OverflowMode // per-column overflow handling
	HeaderAlign string       // defaults to Align if empty
	Bold        bool         // render cell content bold
}

// TableComponent renders a structured data table with optional header, striping,
// borders, and per-column overflow control.
//
// Note: ShowHeader and RowStriping default to false (Go zero value).
// Set them explicitly to true in your config.
//
// BorderStyle values:
//   "all"     — full grid (border around every cell)
//   "outer"   — border around each row only
//   "columns" — outer border for the whole table + column separators + header bottom line (matches Lucid ELD HTML design)
//   "none"    — no borders
type TableComponent struct {
	Columns      []ColumnDef
	Rows         [][]string
	ShowHeader   bool         // render the column header row
	RowStriping  bool         // alternate row background colors
	CellPaddingH float64      // horizontal cell padding mm; default 3
	CellPaddingV float64      // vertical cell padding mm; default 2
	BorderStyle  string       // "none", "outer", "all", "columns"; default "all"
	HeaderFont   FontConfig   // zero value → theme default, bold
	RowFont      FontConfig   // zero value → theme default
	MinRowHeight float64      // mm; default 8
}

// Render draws the table and advances the Y cursor.
func (t *TableComponent) Render(doc *Document) error {
	if len(t.Columns) == 0 {
		return nil
	}

	paddingH := t.CellPaddingH
	if paddingH == 0 {
		paddingH = 2.8
	}
	paddingV := t.CellPaddingV
	if paddingV == 0 {
		paddingV = 2.1
	}
	minRowH := t.MinRowHeight
	if minRowH == 0 {
		minRowH = 9
	}
	lineH := minRowH - 2*paddingV
	if lineH < 3 {
		lineH = 3
	}
	borderStyle := t.BorderStyle
	if borderStyle == "" {
		borderStyle = "all"
	}

	headerFont := t.HeaderFont
	if headerFont.Family == "" {
		headerFont = FontConfig{
			Family: doc.theme.DefaultFont.Family,
			Size:   doc.theme.DefaultFont.Size,
			Style:  "",
		}
	}
	rowFont := t.RowFont
	if rowFont.Family == "" {
		rowFont = doc.theme.DefaultFont
	}

	widths := t.resolveColumnWidths(doc.usableWidth())

	// "columns" style: outer border + column separators + header bottom line.
	// Borders are drawn per-page-section after rows are rendered.
	if borderStyle == "columns" {
		return t.renderColumnsStyle(doc, widths, paddingH, paddingV, minRowH, lineH, headerFont, rowFont)
	}

	if t.ShowHeader {
		t.renderHeaderRow(doc, widths, paddingH, paddingV, minRowH, headerFont, borderStyle)
	}

	for i, row := range t.Rows {
		bgColor := doc.theme.TableRowOddBg
		if t.RowStriping && i%2 == 0 {
			bgColor = doc.theme.TableRowEvenBg
		}

		rowH := t.calcRowHeight(doc, row, widths, paddingH, paddingV, lineH, rowFont)
		if rowH < minRowH {
			rowH = minRowH
		}

		added := doc.newPageIfNeeded(rowH)
		if added && t.ShowHeader {
			t.renderHeaderRow(doc, widths, paddingH, paddingV, minRowH, headerFont, borderStyle)
		}

		t.renderDataRow(doc, row, bgColor, widths, paddingH, paddingV, rowH, lineH, rowFont, borderStyle)
	}

	return nil
}

// renderColumnsStyle renders with the "columns" border mode:
//   - Outer rect + column separator lines drawn per row
//   - Horizontal line below the header row
//   - No horizontal lines between data rows
func (t *TableComponent) renderColumnsStyle(doc *Document, widths []float64, paddingH, paddingV, minRowH, lineH float64, headerFont, rowFont FontConfig) error {
	startX := doc.marginL

	if t.ShowHeader {
		t.renderColumnsRow(doc, nil, true, doc.theme.TableHeaderBg, widths, paddingH, paddingV, minRowH, lineH, headerFont)
		// Header bottom separator line.
		doc.applyColor(doc.theme.TableBorderColor)
		doc.pdf.Line(startX, doc.currentY(), startX+doc.usableWidth(), doc.currentY())
	}

	for i, row := range t.Rows {
		bgColor := doc.theme.TableRowOddBg
		if t.RowStriping && i%2 == 0 {
			bgColor = doc.theme.TableRowEvenBg
		}

		rowH := t.calcRowHeight(doc, row, widths, paddingH, paddingV, lineH, rowFont)
		if rowH < minRowH {
			rowH = minRowH
		}

		added := doc.newPageIfNeeded(rowH)
		if added && t.ShowHeader {
			t.renderColumnsRow(doc, nil, true, doc.theme.TableHeaderBg, widths, paddingH, paddingV, minRowH, lineH, headerFont)
			doc.applyColor(doc.theme.TableBorderColor)
			doc.pdf.Line(startX, doc.currentY(), startX+doc.usableWidth(), doc.currentY())
		}

		t.renderColumnsRow(doc, row, false, bgColor, widths, paddingH, paddingV, rowH, lineH, rowFont)
	}

	return nil
}

// renderColumnsRow draws one row with per-row outer rect + column separators.
// Pass isHeader=true and row=nil to render the header row.
func (t *TableComponent) renderColumnsRow(doc *Document, row []string, isHeader bool, bgColor Color, widths []float64, paddingH, paddingV, rowH, lineH float64, font FontConfig) {
	startY := doc.currentY()
	startX := doc.marginL
	totalW := doc.usableWidth()

	// Step 1: Fill background.
	doc.applyColor(bgColor)
	doc.pdf.Rect(startX, startY, totalW, rowH, "F")

	// Step 2: Draw cell text (no border).
	if isHeader {
		doc.applyFont(font)
		doc.applyTextColor(doc.theme.HeaderTextColor)
		x := startX
		for i, col := range t.Columns {
			align := col.HeaderAlign
			if align == "" {
				align = col.Align
			}
			if align == "" {
				align = "L"
			}
			doc.pdf.SetXY(x+paddingH, startY+paddingV)
			doc.pdf.CellFormat(widths[i]-2*paddingH, rowH-2*paddingV, col.Header, "", 0, align, false, 0, "")
			x += widths[i]
		}
	} else {
		x := startX
		for i, col := range t.Columns {
			text := ""
			if i < len(row) {
				text = row[i]
			}
			align := col.Align
			if align == "" {
				align = "L"
			}
			cellFont := font
			if col.Bold {
				cellFont.Style = "B"
			}
			doc.applyFont(cellFont)
			doc.applyTextColor(doc.theme.PrimaryText)

			cellW := widths[i] - 2*paddingH

			if col.Overflow == OverflowWrap {
				doc.pdf.SetXY(x+paddingH, startY+paddingV)
				doc.pdf.MultiCell(cellW, lineH, text, "", align, false)
			} else {
				if col.Overflow == OverflowTruncate {
					text = truncateText(doc, text, cellW)
				}
				doc.pdf.SetXY(x+paddingH, startY+paddingV)
				doc.pdf.CellFormat(cellW, rowH-2*paddingV, text, "", 0, align, false, 0, "")
			}
			x += widths[i]
		}
	}

	// Step 3: Draw outer row rect + internal column separators.
	doc.applyColor(doc.theme.TableBorderColor)
	doc.pdf.Rect(startX, startY, totalW, rowH, "D")
	x := startX
	for i := 0; i < len(widths)-1; i++ {
		x += widths[i]
		doc.pdf.Line(x, startY, x, startY+rowH)
	}

	doc.setY(startY + rowH)
}

// resolveColumnWidths distributes usable width among columns.
// Fixed-width columns are allocated first; remaining space is split equally
// among columns with Width == 0.
func (t *TableComponent) resolveColumnWidths(usableWidth float64) []float64 {
	widths := make([]float64, len(t.Columns))
	remaining := usableWidth
	autoCount := 0

	for i, col := range t.Columns {
		if col.Width > 0 {
			widths[i] = col.Width
			remaining -= col.Width
		} else {
			autoCount++
		}
	}

	if autoCount > 0 {
		autoW := remaining / float64(autoCount)
		for i := range t.Columns {
			if t.Columns[i].Width == 0 {
				widths[i] = autoW
			}
		}
	}

	return widths
}

// calcRowHeight returns the required row height in mm, accounting for wrapped
// cells (OverflowWrap). Returns 0 if no cell needs more than a single line.
func (t *TableComponent) calcRowHeight(doc *Document, row []string, widths []float64, paddingH, paddingV, lineH float64, font FontConfig) float64 {
	maxContentH := 0.0
	doc.applyFont(font)

	for i, col := range t.Columns {
		if i >= len(row) {
			break
		}
		if col.Overflow != OverflowWrap {
			continue
		}
		cellW := widths[i] - 2*paddingH
		if cellW <= 0 {
			continue
		}
		lines := doc.pdf.SplitLines([]byte(row[i]), cellW)
		h := float64(len(lines)) * lineH
		if h > maxContentH {
			maxContentH = h
		}
	}

	if maxContentH == 0 {
		return 0
	}
	return maxContentH + 2*paddingV
}

func (t *TableComponent) renderHeaderRow(doc *Document, widths []float64, paddingH, paddingV, rowH float64, font FontConfig, borderStyle string) {
	startY := doc.currentY()
	startX := doc.marginL

	// Background
	doc.applyColor(doc.theme.TableHeaderBg)
	doc.pdf.Rect(startX, startY, doc.usableWidth(), rowH, "F")

	// Text — use HeaderTextColor (muted gray) to match the reference design
	doc.applyFont(font)
	doc.applyTextColor(doc.theme.HeaderTextColor)

	x := startX
	for i, col := range t.Columns {
		align := col.HeaderAlign
		if align == "" {
			align = col.Align
		}
		if align == "" {
			align = "L"
		}
		doc.pdf.SetXY(x+paddingH, startY+paddingV)
		doc.pdf.CellFormat(widths[i]-2*paddingH, rowH-2*paddingV, col.Header, "", 0, align, false, 0, "")
		x += widths[i]
	}

	t.drawBorders(doc, startX, startY, widths, rowH, borderStyle)
	doc.setY(startY + rowH)
}

func (t *TableComponent) renderDataRow(doc *Document, row []string, bgColor Color, widths []float64, paddingH, paddingV, rowH, lineH float64, font FontConfig, borderStyle string) {
	startY := doc.currentY()
	startX := doc.marginL

	// Background
	doc.applyColor(bgColor)
	doc.pdf.Rect(startX, startY, doc.usableWidth(), rowH, "F")

	x := startX
	for i, col := range t.Columns {
		text := ""
		if i < len(row) {
			text = row[i]
		}

		align := col.Align
		if align == "" {
			align = "L"
		}

		cellFont := font
		if col.Bold {
			cellFont.Style = "B"
		}
		doc.applyFont(cellFont)
		doc.applyTextColor(doc.theme.PrimaryText)

		cellW := widths[i] - 2*paddingH

		if col.Overflow == OverflowWrap {
			doc.pdf.SetXY(x+paddingH, startY+paddingV)
			doc.pdf.MultiCell(cellW, lineH, text, "", align, false)
		} else {
			if col.Overflow == OverflowTruncate {
				text = truncateText(doc, text, cellW)
			}
			doc.pdf.SetXY(x+paddingH, startY+paddingV)
			doc.pdf.CellFormat(cellW, rowH-2*paddingV, text, "", 0, align, false, 0, "")
		}

		x += widths[i]
	}

	t.drawBorders(doc, startX, startY, widths, rowH, borderStyle)
	doc.setY(startY + rowH)
}

// drawBorders draws cell borders according to the BorderStyle.
func (t *TableComponent) drawBorders(doc *Document, startX, startY float64, widths []float64, rowH float64, borderStyle string) {
	doc.applyColor(doc.theme.TableBorderColor)

	switch borderStyle {
	case "all":
		x := startX
		for _, w := range widths {
			doc.pdf.Rect(x, startY, w, rowH, "D")
			x += w
		}
	case "outer":
		doc.pdf.Rect(startX, startY, doc.usableWidth(), rowH, "D")
	}
}

// truncateText clips text and appends "…" so it fits within maxW mm using
// the currently active font.
func truncateText(doc *Document, text string, maxW float64) string {
	if doc.pdf.GetStringWidth(text) <= maxW {
		return text
	}
	const ellipsis = "…"
	runes := []rune(text)
	for len(runes) > 0 {
		runes = runes[:len(runes)-1]
		candidate := string(runes) + ellipsis
		if doc.pdf.GetStringWidth(candidate) <= maxW {
			return candidate
		}
	}
	return ellipsis
}
