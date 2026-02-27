package pdfgen

// InfoItem is a single label+value pair in an InfoBlockComponent.
type InfoItem struct {
	Label string
	Value string
}

// InfoBlockComponent renders a bordered grid of label+value pairs.
// Typical use: summary boxes like "Total Vehicle: 1 | Total Distance: 7,000 mi".
//
// Column widths from the Lucid ELD HTML reference (A4, 11.3mm margins):
//   Total Vehicle:  37.7mm  (107px × 0.353)
//   Total Distance: 149.6mm (424px × 0.353)
// Set via ColumnWidths; if empty, columns are equal.
type InfoBlockComponent struct {
	Items        []InfoItem
	Columns      int        // items per row; default 2
	ShowBorder   bool       // draw border around each cell
	Width        float64    // mm; 0 = full usable width
	ColumnWidths []float64  // optional per-column widths in mm; must match Columns count
	LabelFont    FontConfig // zero value → theme default, 10pt regular
	ValueFont    FontConfig // zero value → theme default, 10pt bold
}

// Render draws the info block and advances the Y cursor.
func (b *InfoBlockComponent) Render(doc *Document) error {
	if len(b.Items) == 0 {
		return nil
	}

	cols := b.Columns
	if cols <= 0 {
		cols = 2
	}

	totalWidth := b.Width
	if totalWidth == 0 {
		totalWidth = doc.usableWidth()
	}

	// Resolve per-column widths.
	colWidths := make([]float64, cols)
	if len(b.ColumnWidths) == cols {
		copy(colWidths, b.ColumnWidths)
	} else {
		equal := totalWidth / float64(cols)
		for i := range colWidths {
			colWidths[i] = equal
		}
	}

	labelFont := b.LabelFont
	if labelFont.Family == "" {
		labelFont = FontConfig{
			Family: doc.theme.DefaultFont.Family,
			Size:   10,
			Style:  "",
		}
	}
	valueFont := b.ValueFont
	if valueFont.Family == "" {
		valueFont = FontConfig{
			Family: doc.theme.DefaultFont.Family,
			Size:   10,
			Style:  "B",
		}
	}

	// Dimensions from HTML: padding 6px V, 8px H; each text line 14px = 4.94mm
	// Total cell height: 6+14+14+6 = 40px = 14.1mm
	const (
		paddingH  = 2.82 // 8px × 0.352778
		paddingV  = 2.12 // 6px × 0.352778
		labelLineH = 4.94 // 14px line-height
		valueLineH = 4.94
	)
	cellH := paddingV + labelLineH + valueLineH + paddingV

	numRows := (len(b.Items) + cols - 1) / cols
	totalH := float64(numRows) * cellH

	doc.newPageIfNeeded(totalH)
	startX := doc.marginL
	startY := doc.currentY()

	for i, item := range b.Items {
		col := i % cols
		row := i / cols

		// Compute x offset for this column.
		x := startX
		for c := 0; c < col; c++ {
			x += colWidths[c]
		}
		y := startY + float64(row)*cellH
		w := colWidths[col]

		// Cell background (white)
		doc.applyColor(Color{R: 255, G: 255, B: 255})
		doc.pdf.Rect(x, y, w, cellH, "F")

		// Cell border
		if b.ShowBorder {
			doc.applyColor(doc.theme.TableBorderColor)
			doc.pdf.Rect(x, y, w, cellH, "D")
		}

		// Label — small, secondary color
		doc.applyFont(labelFont)
		doc.applyTextColor(doc.theme.SecondaryText)
		doc.pdf.SetXY(x+paddingH, y+paddingV)
		doc.pdf.CellFormat(w-2*paddingH, labelLineH, item.Label, "", 0, "L", false, 0, "")

		// Value — primary color, bold
		doc.applyFont(valueFont)
		doc.applyTextColor(doc.theme.PrimaryText)
		doc.pdf.SetXY(x+paddingH, y+paddingV+labelLineH)
		doc.pdf.CellFormat(w-2*paddingH, valueLineH, item.Value, "", 0, "L", false, 0, "")
	}

	doc.setY(startY + totalH)
	return nil
}
