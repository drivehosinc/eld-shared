package pdfgen

// Color represents an RGB color with components in range [0, 255].
type Color struct {
	R, G, B int
}

// FontConfig describes a font family, size, and style.
type FontConfig struct {
	Family string  // e.g. "Arial", "Helvetica"
	Size   float64 // points
	Style  string  // "", "B", "I", "BI"
}

// ThemeConfig holds document-wide visual settings.
type ThemeConfig struct {
	PrimaryText        Color
	SecondaryText      Color
	AccentColor        Color
	TableHeaderBg      Color
	TableRowEvenBg     Color
	TableRowOddBg      Color
	TableBorderColor   Color
	HeaderTextColor    Color
	SectionLabelLeft   Color // label part before ":" in section label right text; #334155 slate-700
	SectionLabelValue  Color // value part after ":" in section label right text;  #94A3B8 slate-400
	DefaultFont        FontConfig
}

// DefaultTheme returns values derived from the Lucid ELD HTML design source.
//
// Color reference:
//   #181D27 → PrimaryText   (text-primary-900)
//   #535862 → SecondaryText (text-tertiary-600)
//   #E2E8F0 → borders
//   #F1F5F9 → stripe rows   (slate-100)
//   #94A3B8 → AccentColor   (slate-400, used for badge values)
//   #334155 → SectionLabelLeft  (slate-700)
//   #94A3B8 → SectionLabelValue (slate-400)
func DefaultTheme() ThemeConfig {
	return ThemeConfig{
		PrimaryText:        Color{R: 24, G: 29, B: 39},    // #181D27
		SecondaryText:      Color{R: 83, G: 88, B: 98},    // #535862
		AccentColor:        Color{R: 148, G: 163, B: 184}, // #94A3B8 slate-400
		TableHeaderBg:      Color{R: 255, G: 255, B: 255}, // no fill — white
		TableRowEvenBg:     Color{R: 241, G: 245, B: 249}, // #F1F5F9 slate-100
		TableRowOddBg:      Color{R: 255, G: 255, B: 255}, // white
		TableBorderColor:   Color{R: 226, G: 232, B: 240}, // #E2E8F0
		HeaderTextColor:    Color{R: 83, G: 88, B: 98},    // same as SecondaryText
		SectionLabelLeft:   Color{R: 51, G: 65, B: 85},    // #334155 slate-700
		SectionLabelValue:  Color{R: 148, G: 163, B: 184}, // #94A3B8 slate-400
		DefaultFont: FontConfig{
			Family: "Arial",
			Size:   10,
			Style:  "",
		},
	}
}
