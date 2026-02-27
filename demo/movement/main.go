package main

// Measurements derived from movement report HTML (842×595px A4 landscape):
//
//	Page:          842×595px → A4 landscape
//	Margin:        32px = 11.3mm
//	Content width: 778px = 274.5mm
//	Table top:     80px  = 28.2mm from page top
//
// Column widths (px → mm, × 0.352778):
//
//	Vehicle:        51px  = 18.00mm  fixed
//	Date & Time:    flex:1           auto
//	Location:       flex:1           auto
//	Lst:           100px  = 35.28mm  fixed
//	Lng:           100px  = 35.28mm  fixed
//	Speed:          48px  = 16.93mm  fixed
//	Heading:        49px  = 17.29mm  fixed  (data: right-aligned)
//	Odometer Miles: 81px  = 28.57mm  fixed
//	Eng Hours:      57px  = 20.11mm  fixed
//
// Fixed total: 171.46mm  Remaining: 274.5 - 171.46 = 103.04mm → 51.52mm each flex col
//
// Row: 8px font, paddingV 6px = 2.12mm, lineH 14px = 4.94mm, rowH 26px = 9.17mm
// Logo: logos/logo_lucid.png 512×128 px (4:1) → 39.5mm × 9.9mm

import (
	"fmt"
	"os"

	"github.com/drivehosinc/eld-shared/pdfgen"
)

const margin = 11.3 // mm (32px × 0.352778)

func main() {
	theme := pdfgen.DefaultTheme()

	doc := pdfgen.New(pdfgen.DocumentConfig{
		PageSize:     "A4",
		Orientation:  "landscape",
		MarginTop:    margin,
		MarginBottom: margin,
		MarginLeft:   margin,
		MarginRight:  margin,
		Theme:        theme,
	})

	// 9 columns: 2 flex (Date&Time, Location), 7 fixed.
	cols := []pdfgen.ColumnDef{
		{Header: "Vehicle",        Width: 18.00, Align: "L"},
		{Header: "Date & Time (EST)", Width: 0,  Align: "L", Overflow: pdfgen.OverflowTruncate},
		{Header: "Location",       Width: 0,     Align: "L", Overflow: pdfgen.OverflowTruncate},
		{Header: "Lst",            Width: 35.28, Align: "L"},
		{Header: "Lng",            Width: 35.28, Align: "L"},
		{Header: "Speed",          Width: 16.93, Align: "L"},
		{Header: "Heading",        Width: 17.29, Align: "R", HeaderAlign: "L"},
		{Header: "Odometer Miles", Width: 28.57, Align: "L"},
		{Header: "Eng Hours",      Width: 20.11, Align: "L"},
	}

	rows := [][]string{
		{"7000", "02/10/2026 07:44:53", "0.8mi NE of Greenfield Manor, PA", "40.202002", "-40.202002", "64.0", "344", "0.0", "0.0"},
		{"7000", "02/10/2026 07:44:53", "0.8mi NE of Greenfield Manor, PA", "40.202002", "-40.202002", "64.0", "344", "0.0", "0.0"},
		{"7000", "02/10/2026 07:44:53", "0.8mi NE of Greenfield Manor, PA", "40.202002", "-40.202002", "64.0", "344", "0.0", "0.0"},
		{"7000", "02/10/2026 07:44:53", "0.8mi NE of Greenfield Manor, PA", "40.202002", "-40.202002", "64.0", "344", "0.0", "0.0"},
		{"7000", "02/10/2026 07:44:53", "0.8mi NE of Greenfield Manor, PA", "40.202002", "-40.202002", "64.0", "344", "0.0", "0.0"},
		{"7000", "02/10/2026 07:44:53", "0.8mi NE of Greenfield Manor, PA", "40.202002", "-40.202002", "64.0", "344", "0.0", "0.0"},
		{"7000", "02/10/2026 07:44:53", "0.8mi NE of Greenfield Manor, PA", "40.202002", "-40.202002", "64.0", "344", "0.0", "0.0"},
		{"7000", "02/10/2026 07:44:53", "0.8mi NE of Greenfield Manor, PA", "40.202002", "-40.202002", "64.0", "344", "0.0", "0.0"},
		{"7000", "02/10/2026 07:44:53", "0.8mi NE of Greenfield Manor, PA", "40.202002", "-40.202002", "64.0", "344", "0.0", "0.0"},
		{"7000", "02/10/2026 07:44:53", "0.8mi NE of Greenfield Manor, PA", "40.202002", "-40.202002", "64.0", "344", "0.0", "0.0"},
		{"7000", "02/10/2026 07:44:53", "0.8mi NE of Greenfield Manor, PA", "40.202002", "-40.202002", "64.0", "344", "0.0", "0.0"},
		{"7000", "02/10/2026 07:44:53", "0.8mi NE of Greenfield Manor, PA", "40.202002", "-40.202002", "64.0", "344", "0.0", "0.0"},
	}

	doc.Add(
		// Logo: top-right, 512×128 PNG rendered at 39.5×9.9mm (4:1 aspect preserved).
		&pdfgen.LogoComponent{
			ImagePath: "demo/logos/logo_lucid.png",
			Width:     39.5,
			Height:    9.9,
			Position:  "top-right",
			OffsetY:   0,
		},

		// Title: 16pt bold #181D27, left-aligned.
		// No subtitle or date lines — table follows immediately.
		&pdfgen.HeaderComponent{
			Title: "MOVEMENT REPORT",
		},

		// HTML: title top=32px, table top=80px → gap = (80-32-24)px = 24px ≈ 5.4mm
		// HeaderComponent advances ~11.5mm from marginTop, so spacer brings us to ~28.2mm.
		&pdfgen.SpacerComponent{Height: 5.4},

		// Movement table: 8pt font, "columns" border style, same padding as HTML.
		&pdfgen.TableComponent{
			ShowHeader:   true,
			RowStriping:  true,
			BorderStyle:  "columns",
			CellPaddingH: 2.82, // 8px × 0.352778
			CellPaddingV: 2.12, // 6px × 0.352778
			MinRowHeight: 9.17, // 26px × 0.352778
			HeaderFont: pdfgen.FontConfig{
				Family: theme.DefaultFont.Family,
				Size:   8,
				Style:  "", // font-weight: 400 in HTML
			},
			RowFont: pdfgen.FontConfig{
				Family: theme.DefaultFont.Family,
				Size:   8,
				Style:  "",
			},
			Columns: cols,
			Rows:    rows,
		},
	)

	output := "demo/movement/movement_report.pdf"
	if err := doc.Save(output); err != nil {
		fmt.Printf("save failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("PDF generated: %s\n", output)
}
