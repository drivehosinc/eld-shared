package pdfgen

// Component is implemented by every renderable element.
// Render writes the component into the document starting at the current cursor
// position and advances the Y cursor unless documented otherwise.
type Component interface {
	Render(doc *Document) error
}
