# pdfgen — Release Notes

## v0.1.0

### Overview

First release of `github.com/drivehosinc/eld-shared/pdfgen` — a structured
PDF generation package for Lucid ELD business reports. Wraps
[go-pdf/fpdf](https://github.com/go-pdf/fpdf) behind a component-based API
so callers never touch the low-level PDF primitives directly.

### Components

| Component | Description |
|---|---|
| `Document` | Root object. Configures page size, orientation, margins, and theme. |
| `HeaderComponent` | Title / subtitle / date lines. Reserves right side for logo. |
| `LogoComponent` | Places a PNG image at a fixed page position (top-right). Does not affect content flow. |
| `InfoBlockComponent` | Bordered grid of label+value pairs (e.g. "Total Vehicle: 4"). Supports unequal column widths. |
| `SectionLabelComponent` | Horizontal row with bold left label and two-color right badge (colon-split). |
| `TableComponent` | Striped data table. Header repeats on page break. Per-column overflow: wrap or truncate. |
| `GroupedTableComponent` | Convenience wrapper: SectionLabel → Table → Spacer. |
| `FooterComponent` | Renders on every page. Supports `{page}` and `{total}` placeholders. |
| `SpacerComponent` | Adds vertical whitespace between components. |

### Theme

`DefaultTheme()` is calibrated to the Lucid ELD HTML design system:

| Token | Value | Usage |
|---|---|---|
| PrimaryText | `#181D27` | Body text |
| SecondaryText | `#535862` | Labels, table headers |
| TableBorderColor | `#E2E8F0` | All borders |
| TableRowEvenBg | `#F1F5F9` | Alternating row stripe |
| SectionLabelLeft | `#334155` | Right-badge label part |
| SectionLabelValue | `#94A3B8` | Right-badge value part |

### Border Styles

`TableComponent.BorderStyle` accepts:

- `"columns"` — outer rect + column separator lines per row, no horizontal row lines (matches Lucid ELD HTML)
- `"all"` — full grid
- `"outer"` — outline around each row only
- `"none"` — no borders

### Usage

```go
import "github.com/drivehosinc/eld-shared/pdfgen"

doc := pdfgen.New(pdfgen.DocumentConfig{
    PageSize:    "A4",
    Orientation: "portrait",
})

doc.SetFooter(&pdfgen.FooterComponent{
    RightText: "Page {page} of {total}",
})

doc.Add(
    &pdfgen.LogoComponent{ImagePath: "logo.png", Width: 39.5, Height: 9.9, Position: "top-right"},
    &pdfgen.HeaderComponent{Title: "IFTA REPORT", Subtitle: "Company Name"},
    &pdfgen.GroupedTableComponent{
        Label:     "By state",
        BadgeText: "Total Distance: 7,000 mi",
        Table: pdfgen.TableComponent{
            ShowHeader:  true,
            RowStriping: true,
            BorderStyle: "columns",
            Columns: []pdfgen.ColumnDef{
                {Header: "State",    Width: 0,  Align: "L"},
                {Header: "Distance", Width: 40, Align: "L"},
            },
            Rows: [][]string{
                {"California", "1,240 mi"},
                {"Oregon",     "870 mi"},
            },
        },
    },
)

doc.Save("output.pdf")
```

### Demos

| File | Description |
|---|---|
| `demo/main.go` | A4 portrait — IFTA quarterly report, 3 pages, page-number footer |
| `demo/movement/main.go` | A4 landscape — Movement report, 9-column table, 8pt font |

Run from the `eld-shared` root:

```bash
go run ./demo/main.go
go run ./demo/movement/main.go
```

### Dependencies

- `github.com/go-pdf/fpdf v0.9.0`
