package pdfgen

import (
	"fmt"
	"strings"
)

// FooterComponent renders a one-line footer on every page.
// Register it via doc.SetFooter() — do NOT pass to doc.Add().
//
// Placeholders in text fields:
//   {page}  → current page number
//   {total} → total page count
type FooterComponent struct {
	LeftText   string     // static text, left-aligned
	CenterText string     // center-aligned; supports {page} and {total}
	RightText  string     // static text, right-aligned
	ShowBorder bool       // draw a thin top border line
	Font       FontConfig // zero value → theme default at 8pt
	TextColor  Color      // zero value → theme SecondaryText
}

// Render implements Component so FooterComponent can also be used standalone,
// but it is normally called automatically via the fpdf footer callback.
func (f *FooterComponent) Render(doc *Document) error {
	f.render(doc)
	return nil
}

// render is called by the fpdf footer callback on every page.
func (f *FooterComponent) render(doc *Document) {
	pdf := doc.pdf

	// Position at the bottom margin area.
	pdf.SetY(-doc.marginB)

	font := f.Font
	if font.Family == "" {
		font = doc.theme.DefaultFont
		font.Size = 8
	}
	doc.applyFont(font)

	color := f.TextColor
	if color.R == 0 && color.G == 0 && color.B == 0 {
		color = doc.theme.SecondaryText
	}
	doc.applyTextColor(color)

	if f.ShowBorder {
		y := pdf.GetY()
		doc.applyColor(doc.theme.TableBorderColor)
		pdf.Line(doc.marginL, y, doc.marginL+doc.pageWidth, y)
		pdf.SetY(y + 1)
		doc.applyTextColor(color)
	}

	// Replace {page} with the current page number.
	// {total} is left as-is; fpdf replaces it at output time via AliasNbPages.
	pageNum := fmt.Sprintf("%d", pdf.PageNo())
	center := strings.ReplaceAll(f.CenterText, "{page}", pageNum)

	w := doc.pageWidth
	h := 5.0

	pdf.SetX(doc.marginL)
	pdf.CellFormat(w/3, h, f.LeftText, "", 0, "L", false, 0, "")
	pdf.CellFormat(w/3, h, center, "", 0, "C", false, 0, "")
	pdf.CellFormat(w/3, h, f.RightText, "", 0, "R", false, 0, "")
}
