package pdfgen

// GroupedTableComponent is a convenience wrapper that composes a
// SectionLabelComponent followed by a TableComponent and an optional spacer.
type GroupedTableComponent struct {
	Label       string         // left text of the section label
	BadgeText   string         // right text (supports ":" splitting for two-color)
	Table       TableComponent // embedded table
	SpacerAfter float64        // mm of whitespace after the table; default 6
}

// Render draws: SectionLabel → Table → Spacer.
func (g *GroupedTableComponent) Render(doc *Document) error {
	label := &SectionLabelComponent{
		LeftText:  g.Label,
		RightText: g.BadgeText,
	}
	if err := label.Render(doc); err != nil {
		return err
	}

	if err := g.Table.Render(doc); err != nil {
		return err
	}

	gap := g.SpacerAfter
	if gap == 0 {
		gap = 8
	}
	spacer := &SpacerComponent{Height: gap}
	return spacer.Render(doc)
}
