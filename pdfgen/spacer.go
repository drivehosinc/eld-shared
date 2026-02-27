package pdfgen

// SpacerComponent advances the Y cursor by a fixed amount.
type SpacerComponent struct {
	Height float64 // mm
}

// Render advances the document cursor by Height millimetres.
func (s *SpacerComponent) Render(doc *Document) error {
	doc.setY(doc.currentY() + s.Height)
	return nil
}
