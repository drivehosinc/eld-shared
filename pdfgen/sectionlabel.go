package pdfgen

import "strings"

// SectionLabelComponent renders a single horizontal line with left label and
// optional right text.
//
// Colon-split coloring applies to BOTH sides independently:
//   LeftText  "Vehicle number: 7981"  → "Vehicle number:" (SecondaryText) + " 7981" (AccentColor)
//   RightText "Total Distance: 7,000" → "Total Distance:" (RightLabelColor) + " 7,000" (RightColor)
//
// Color defaults:
//   LeftText colon-split value  → AccentColor
//   RightLabelColor             → SecondaryText
//   RightColor (value after ":") → SecondaryText (lighter, set explicitly for the reference design)
type SectionLabelComponent struct {
	LeftText        string
	RightText       string     // optional; split on ":" for two-color rendering
	LeftFont        FontConfig // zero value → theme default, bold
	RightFont       FontConfig // zero value → theme default
	RightLabelColor Color      // color for the label part before ":" in RightText; zero = SecondaryText
	RightColor      Color      // color for the value part after ":" in RightText; zero = AccentColor
	MarginBottom    float64    // mm below the line; default 2
}

// Render draws the section label row and advances the Y cursor.
func (s *SectionLabelComponent) Render(doc *Document) error {
	leftFont := s.LeftFont
	if leftFont.Family == "" {
		leftFont = FontConfig{
			Family: doc.theme.DefaultFont.Family,
			Size:   doc.theme.DefaultFont.Size,
			Style:  "B",
		}
	}
	rightFont := s.RightFont
	if rightFont.Family == "" {
		rightFont = doc.theme.DefaultFont
	}

	const lineH = 7.0
	startY := doc.currentY()

	// ── Left text ────────────────────────────────────────────────────────────
	// When LeftText contains ":", split into label (secondary) + value (accent).
	if s.LeftText != "" {
		leftParts := strings.SplitN(s.LeftText, ":", 2)
		if len(leftParts) == 2 {
			labelPart := leftParts[0] + ":"
			valuePart := leftParts[1]

			doc.applyFont(leftFont)
			labelW := doc.pdf.GetStringWidth(labelPart) + 1

			accentFont := leftFont
			accentFont.Style = "B"
			doc.applyFont(accentFont)
			valueW := doc.pdf.GetStringWidth(valuePart) + 1

			// Label part → secondary color
			doc.applyFont(leftFont)
			doc.applyTextColor(doc.theme.SecondaryText)
			doc.pdf.SetXY(doc.marginL, startY)
			doc.pdf.CellFormat(labelW, lineH, labelPart, "", 0, "L", false, 0, "")

			// Value part → accent color
			doc.applyFont(accentFont)
			doc.applyTextColor(doc.theme.AccentColor)
			doc.pdf.SetXY(doc.marginL+labelW, startY)
			doc.pdf.CellFormat(valueW+2, lineH, valuePart, "", 0, "L", false, 0, "")
		} else {
			// No colon → plain bold primary text
			doc.applyFont(leftFont)
			doc.applyTextColor(doc.theme.PrimaryText)
			doc.pdf.SetXY(doc.marginL, startY)
			doc.pdf.CellFormat(doc.usableWidth()*0.5, lineH, s.LeftText, "", 0, "L", false, 0, "")
		}
	}

	// ── Right text ───────────────────────────────────────────────────────────
	if s.RightText != "" {
		rightEdge := doc.marginL + doc.usableWidth()

		rightLabelColor := s.RightLabelColor
		if rightLabelColor.R == 0 && rightLabelColor.G == 0 && rightLabelColor.B == 0 {
			rightLabelColor = doc.theme.SectionLabelLeft
		}
		rightValueColor := s.RightColor
		if rightValueColor.R == 0 && rightValueColor.G == 0 && rightValueColor.B == 0 {
			rightValueColor = doc.theme.SectionLabelValue
		}

		parts := strings.SplitN(s.RightText, ":", 2)
		if len(parts) == 2 {
			labelPart := parts[0] + ":"
			valuePart := parts[1]

			doc.applyFont(rightFont)
			labelW := doc.pdf.GetStringWidth(labelPart) + 1

			doc.applyFont(rightFont)
			valueW := doc.pdf.GetStringWidth(valuePart) + 2 // +2 for gap

			startRX := rightEdge - labelW - valueW

			// Label part
			doc.applyFont(rightFont)
			doc.applyTextColor(rightLabelColor)
			doc.pdf.SetXY(startRX, startY)
			doc.pdf.CellFormat(labelW, lineH, labelPart, "", 0, "L", false, 0, "")

			// Value part
			doc.applyFont(rightFont)
			doc.applyTextColor(rightValueColor)
			doc.pdf.SetXY(startRX+labelW, startY)
			doc.pdf.CellFormat(valueW+1, lineH, valuePart, "", 0, "L", false, 0, "")
		} else {
			// No colon → right-aligned, right label color
			doc.applyFont(rightFont)
			doc.applyTextColor(rightLabelColor)
			doc.pdf.SetXY(doc.marginL, startY)
			doc.pdf.CellFormat(doc.usableWidth(), lineH, s.RightText, "", 0, "R", false, 0, "")
		}
	}

	mb := s.MarginBottom
	if mb == 0 {
		mb = 2
	}
	doc.setY(startY + lineH + mb)
	return nil
}
