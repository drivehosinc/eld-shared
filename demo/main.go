package main

import (
	"fmt"
	"os"

	"github.com/drivehosinc/eld-shared/pdfgen"
)

// All measurements derived from demo/sample.html at 1px = 1pt = 0.352778mm
//
//	Page:          595×842px → A4
//	Margin:        32px = 11.29mm
//	Content width: 531px = 187.3mm
//	Border color:  #E2E8F0
//	Stripe color:  #F1F5F9 (slate-100)
//	Primary text:  #181D27
//	Label text:    #535862
const margin = 11.3 // mm (32px × 0.352778)

func main() {
	// Use exact theme values from HTML source.
	theme := pdfgen.DefaultTheme() // already set to correct values

	doc := pdfgen.New(pdfgen.DocumentConfig{
		PageSize:     "A4",
		Orientation:  "portrait",
		MarginTop:    margin,
		MarginBottom: margin,
		MarginLeft:   margin,
		MarginRight:  margin,
		Theme:        theme,
	})

	// Footer with page numbers on every page.
	doc.SetFooter(&pdfgen.FooterComponent{
		LeftText:  "QGM Express — IFTA Q4 2025",
		RightText: "Page {page} of {total}",
	})

	// Column layout for all state-distance tables.
	// HTML uses flex:1 for "No" and "State" — equal share of remaining after nothing fixed.
	// Looking at the PDF, "No" is narrow. We'll use fixed 15mm for No, flex for State, fixed 40mm for Distance.
	stateCols := []pdfgen.ColumnDef{
		{Header: "No", Width: 15, Align: "L"},
		{Header: "State", Width: 0, Align: "L"},
		{Header: "Distance", Width: 40, Align: "L"},
	}

	// Column layout for the "By vehicle" table.
	// HTML: Vehicle=flex:1, VIN=199px=70.2mm, State=flex:1, Total miles=flex:1
	// Remaining after VIN: 187.3 - 70.2 = 117.1mm → 3 flex cols × 39mm each
	vehicleCols := []pdfgen.ColumnDef{
		{Header: "Vehicle", Width: 39, Align: "L"},
		{Header: "VIN", Width: 70.2, Align: "L", Overflow: pdfgen.OverflowTruncate},
		{Header: "State", Width: 39, Align: "L"},
		{Header: "Total miles", Width: 0, Align: "L"},
	}

	// Table config shared across sections.
	makeTable := func(cols []pdfgen.ColumnDef, rows [][]string) pdfgen.TableComponent {
		return pdfgen.TableComponent{
			ShowHeader:   true,
			RowStriping:  true,
			BorderStyle:  "columns", // matches HTML: outer + column separators, no row borders
			CellPaddingH: 2.82,      // 8px × 0.352778
			CellPaddingV: 2.12,      // 6px × 0.352778
			MinRowHeight: 9.17,      // 6+14+6px = 26px... actually paddingV+lineH+paddingV = 2.12+4.94+2.12 = 9.18mm
			HeaderFont: pdfgen.FontConfig{
				Family: theme.DefaultFont.Family,
				Size:   10,
				Style:  "", // header is regular weight in HTML (font-weight: 400)
			},
			RowFont: pdfgen.FontConfig{
				Family: theme.DefaultFont.Family,
				Size:   10,
				Style:  "",
			},
			Columns: cols,
			Rows:    rows,
		}
	}

	// 35 rows — enough to overflow page 1 mid-table and continue on page 2.
	// Page 1 fits ~23 data rows; rows 24-35 spill to page 2 with header repeated.
	sampleStateRows := [][]string{
		{"1", "California", "1,240 mi"},
		{"2", "Oregon", "870 mi"},
		{"3", "Washington", "530 mi"},
		{"4", "Nevada", "410 mi"},
		{"5", "Arizona", "620 mi"},
		{"6", "Utah", "390 mi"},
		{"7", "Colorado", "480 mi"},
		{"8", "New Mexico", "310 mi"},
		{"9", "Texas", "1,050 mi"},
		{"10", "Oklahoma", "280 mi"},
		{"11", "Kansas", "340 mi"},
		{"12", "Nebraska", "290 mi"},
		{"13", "South Dakota", "210 mi"},
		{"14", "North Dakota", "175 mi"},
		{"15", "Montana", "260 mi"},
		{"16", "Idaho", "320 mi"},
		{"17", "Wyoming", "195 mi"},
		{"18", "Minnesota", "430 mi"},
		{"19", "Iowa", "370 mi"},
		{"20", "Missouri", "440 mi"},
		{"21", "Arkansas", "220 mi"},
		{"22", "Louisiana", "310 mi"},
		{"23", "Mississippi", "190 mi"},
		{"24", "Alabama", "240 mi"},
		{"25", "Tennessee", "380 mi"},
		{"26", "Kentucky", "290 mi"},
		{"27", "Indiana", "350 mi"},
		{"28", "Ohio", "420 mi"},
		{"29", "Michigan", "480 mi"},
		{"30", "Wisconsin", "310 mi"},
		{"31", "Illinois", "400 mi"},
		{"32", "Georgia", "510 mi"},
		{"33", "Florida", "730 mi"},
		{"34", "South Carolina", "280 mi"},
		{"35", "North Carolina", "460 mi"},
	}

	// 20 rows across 4 vehicles — fills most of the remaining space on pages 2–3.
	sampleVehicleRows := [][]string{
		{"7070", "1HGCM82633A123456", "CA", "1,240"},
		{"7070", "1HGCM82633A123456", "OR", "870"},
		{"7070", "1HGCM82633A123456", "WA", "530"},
		{"7070", "1HGCM82633A123456", "NV", "410"},
		{"7070", "1HGCM82633A123456", "AZ", "620"},
		{"8120", "2T1BURHE0JC043821", "TX", "1,050"},
		{"8120", "2T1BURHE0JC043821", "OK", "280"},
		{"8120", "2T1BURHE0JC043821", "KS", "340"},
		{"8120", "2T1BURHE0JC043821", "MO", "440"},
		{"8120", "2T1BURHE0JC043821", "IL", "400"},
		{"9341", "3VWFE21C04M000001", "GA", "510"},
		{"9341", "3VWFE21C04M000001", "FL", "730"},
		{"9341", "3VWFE21C04M000001", "SC", "280"},
		{"9341", "3VWFE21C04M000001", "NC", "460"},
		{"9341", "3VWFE21C04M000001", "TN", "380"},
		{"6050", "4T1BF3EK8AU123789", "OH", "420"},
		{"6050", "4T1BF3EK8AU123789", "MI", "480"},
		{"6050", "4T1BF3EK8AU123789", "IN", "350"},
		{"6050", "4T1BF3EK8AU123789", "KY", "290"},
		{"6050", "4T1BF3EK8AU123789", "WI", "310"},
	}

	doc.Add(
		// ── Logo: logos/logo_lucid.png (512×128 px, aspect 4:1) ─────────────────
		// HTML: ~112px wide × 20px tall → 39.5mm wide; height = 39.5/4 ≈ 9.9mm
		&pdfgen.LogoComponent{
			ImagePath: "demo/logos/logo_lucid.png",
			Width:     39.5,
			Height:    9.9, // 39.5 / 4 (preserves 512×128 aspect ratio)
			Position:  "top-right",
			OffsetY:   0,
		},

		// ── Header (HTML lines 2-10) ──────────────────────────────────────────
		// Title:    16px bold #181D27
		// Subtitle: 10px semibold #181D27  (not accent blue)
		// Date:      8px regular #535862
		&pdfgen.HeaderComponent{
			Title:         "IFTA REPORT",
			Subtitle:      "QGM EXPRESS",
			SubtitleColor: pdfgen.Color{R: 24, G: 29, B: 39}, // #181D27 — dark, not accent
			Lines:         []string{"Oct 1, 2025 - Dec 31, 2025"},
		},

		&pdfgen.SpacerComponent{Height: 3},

		// ── Info block (HTML lines 12-21) ─────────────────────────────────────
		// Height: 40px = 14.1mm | border: #E2E8F0
		// Left (Total Vehicle):  flex:1 = 107px = 37.74mm
		// Right (Total Distance): 424px = 149.59mm
		&pdfgen.InfoBlockComponent{
			Items: []pdfgen.InfoItem{
				{Label: "Total Vehicle", Value: "4"},
				{Label: "Total Distance", Value: "14,123 mi"},
			},
			Columns:      2,
			ShowBorder:   true,
			ColumnWidths: []float64{37.74, 149.59},
		},

		&pdfgen.SpacerComponent{Height: 6},

		// ── Total distance per state (HTML lines 22-29, 30+) ─────────────────
		&pdfgen.GroupedTableComponent{
			Label:       "Total distance per state",
			BadgeText:   "Total Distance: 14,123 mi",
			Table:       makeTable(stateCols, sampleStateRows),
			SpacerAfter: 6,
		},

		// ── By vehicle (HTML lines 22+) ───────────────────────────────────────
		// Left label uses no colon → plain primary text
		// Right uses slate-700 / slate-400 colors
		&pdfgen.GroupedTableComponent{
			Label:       "By vehicle",
			BadgeText:   "Total Distance: 14,123 mi",
			Table:       makeTable(vehicleCols, sampleVehicleRows),
			SpacerAfter: 0,
		},
	)

	output := "demo/sample_output.pdf"
	if err := doc.Save(output); err != nil {
		fmt.Printf("save failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("PDF generated: %s\n", output)
}
