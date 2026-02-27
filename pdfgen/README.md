# pdfgen — Agent Reference Guide

> **For AI agents**: Read this file. You do not need to read any `.go` source file.
> Every type, field, default, and usage pattern you need to generate any PDF is here.

---

## Import

```go
import "github.com/drivehosinc/eld-pfd-generator/pdfgen"
```

---

## Quick Start (minimal working PDF)

```go
doc := pdfgen.New(pdfgen.DocumentConfig{
    PageSize: "A4",  // or "Letter"
    Theme:    pdfgen.DefaultTheme(),
})

doc.Add(
    &pdfgen.HeaderComponent{Title: "MY REPORT"},
    &pdfgen.SpacerComponent{Height: 5},
    &pdfgen.TableComponent{
        ShowHeader:  true,
        RowStriping: true,
        Columns: []pdfgen.ColumnDef{
            {Header: "Name",  Width: 60, Align: "L"},
            {Header: "Value", Width: 0,  Align: "R"},
        },
        Rows: [][]string{
            {"Row one", "100"},
            {"Row two", "200"},
        },
    },
)

doc.Save("output.pdf")
```

---

## Document Setup

### `DocumentConfig`

| Field          | Type         | Default      | Notes                              |
|----------------|--------------|--------------|------------------------------------|
| `PageSize`     | `string`     | `"A4"`       | `"A4"` or `"Letter"`               |
| `Orientation`  | `string`     | `"portrait"` | `"portrait"` or `"landscape"`      |
| `MarginTop`    | `float64`    | `15`         | mm                                 |
| `MarginBottom` | `float64`    | `15`         | mm                                 |
| `MarginLeft`   | `float64`    | `15`         | mm                                 |
| `MarginRight`  | `float64`    | `15`         | mm                                 |
| `Theme`        | `ThemeConfig`| DefaultTheme | call `pdfgen.DefaultTheme()`       |

```go
doc := pdfgen.New(pdfgen.DocumentConfig{
    PageSize:     "Letter",
    Orientation:  "landscape",
    MarginTop:    20,
    MarginBottom: 20,
    MarginLeft:   15,
    MarginRight:  15,
    Theme:        pdfgen.DefaultTheme(),
})
```

### `ThemeConfig` — all color/font overrides

```go
theme := pdfgen.DefaultTheme()          // start from sensible defaults
theme.AccentColor = pdfgen.Color{220, 50, 50}  // override one field
theme.DefaultFont = pdfgen.FontConfig{Family: "Helvetica", Size: 10}

doc := pdfgen.New(pdfgen.DocumentConfig{Theme: theme})
```

**ThemeConfig fields**

| Field              | Default RGB        | Used for                             |
|--------------------|--------------------|--------------------------------------|
| `PrimaryText`      | `{30, 30, 30}`     | Body text, table cell text           |
| `SecondaryText`    | `{120, 120, 120}`  | Labels, footer, muted lines          |
| `AccentColor`      | `{66, 133, 244}`   | Subtitle, value highlights           |
| `TableHeaderBg`    | `{245, 245, 245}`  | Table header row background          |
| `TableRowEvenBg`   | `{250, 250, 252}`  | Even rows when striping is on        |
| `TableRowOddBg`    | `{255, 255, 255}`  | Odd rows / default row background    |
| `TableBorderColor` | `{220, 220, 220}`  | Table cell borders, footer line      |
| `HeaderTextColor`  | `{100, 100, 100}`  | Available for custom use             |
| `DefaultFont`      | Arial 10pt         | Base font for all components         |

### `FontConfig`

```go
pdfgen.FontConfig{
    Family: "Arial",       // "Arial", "Helvetica", "Times", "Courier"
    Size:   12,            // points
    Style:  "B",           // "" normal | "B" bold | "I" italic | "BI" bold+italic
}
```

Zero-value `FontConfig{}` falls back to the document theme's `DefaultFont`.

### `Color`

```go
pdfgen.Color{R: 66, G: 133, B: 244}   // blue
pdfgen.Color{255, 0, 0}               // red (positional)
pdfgen.Color{0, 0, 0}                 // black
pdfgen.Color{255, 255, 255}           // white
```

---

## Document Methods

```go
doc := pdfgen.New(cfg)              // create document
doc.SetFooter(&pdfgen.FooterComponent{...})  // register footer (call before Add)
doc.Add(component1, component2, ...)         // add components, chainable
doc.Save("output.pdf")             // write to file → returns error
data, err := doc.Bytes()           // write to []byte
```

---

## Components

Components are added via `doc.Add(...)`. Each advances the Y cursor **unless noted otherwise**.

---

### 1. `LogoComponent` — place an image (does NOT advance Y)

Renders a PNG or JPG at a fixed page position. Logos float — they do not push content down.

```go
&pdfgen.LogoComponent{
    ImagePath: "path/to/logo.png",  // path to file on disk
    // ImageData: []byte{...},       // OR raw PNG/JPG bytes (one of the two)
    Width:    42,       // mm — required, must be > 0
    Height:   17,       // mm — 0 = auto-preserve aspect ratio
    Position: "top-right",   // "top-left" | "top-right" | "top-center"
    OffsetX:  0,        // mm, additional X shift
    OffsetY:  2,        // mm, additional Y shift from top margin
}
```

| Field       | Type      | Default     | Notes                                |
|-------------|-----------|-------------|--------------------------------------|
| `ImagePath` | `string`  | —           | Path to PNG/JPG file                 |
| `ImageData` | `[]byte`  | —           | Raw bytes alternative to ImagePath   |
| `Width`     | `float64` | —           | mm, required                         |
| `Height`    | `float64` | `0`         | mm; 0 = preserve aspect ratio        |
| `Position`  | `string`  | `"top-left"`| `"top-left"`, `"top-right"`, `"top-center"` |
| `OffsetX`   | `float64` | `0`         | mm extra X offset                    |
| `OffsetY`   | `float64` | `0`         | mm extra Y offset from top margin    |

**Rules:**
- Add logo **before** header so it's placed at the right Y (top of page).
- Multiple logos can be added independently (e.g., one top-left, one top-right).
- Y cursor is **restored** after rendering — content flow is unaffected.

---

### 2. `HeaderComponent` — document title block

Left-aligned title, subtitle, and detail lines. Uses ~65% of page width to avoid overlapping a top-right logo.

```go
&pdfgen.HeaderComponent{
    Title:    "IFTA REPORT",      // large bold, ~20pt
    Subtitle: "QGM EXPRESS",      // medium bold, ~13pt, accent color
    Lines: []string{              // small regular text lines
        "1234 Trucker Lane, Nashville, TN 37201",
        "Period: Q1 2024  |  Jan 1 – Mar 31, 2024",
        "Report No: IFTA-2024-001",
    },
    // Optional font overrides (zero value = theme defaults above):
    TitleFont:    pdfgen.FontConfig{},
    SubtitleFont: pdfgen.FontConfig{},
    LineFont:     pdfgen.FontConfig{},
}
```

| Field          | Type         | Default              |
|----------------|--------------|----------------------|
| `Title`        | `string`     | —                    |
| `Subtitle`     | `string`     | —                    |
| `Lines`        | `[]string`   | —                    |
| `TitleFont`    | `FontConfig` | Arial 20pt Bold      |
| `SubtitleFont` | `FontConfig` | Arial 13pt Bold      |
| `LineFont`     | `FontConfig` | Arial 9pt            |

---

### 3. `InfoBlockComponent` — summary grid

A bordered grid of label+value pairs. Each cell stacks the label (small, gray) over the value (bold, larger).

```go
&pdfgen.InfoBlockComponent{
    Items: []pdfgen.InfoItem{
        {Label: "Total Vehicle",   Value: "1"},
        {Label: "Total Distance",  Value: "7,000 mi"},
        {Label: "Report Period",   Value: "Q1 2024"},
        {Label: "Generated",       Value: "2024-04-01"},
    },
    Columns:    2,     // items per row; default 2
    ShowBorder: true,  // draw border around each cell
    Width:      0,     // mm; 0 = full usable page width
    // Optional font overrides:
    LabelFont:  pdfgen.FontConfig{},  // default: Arial 8pt
    ValueFont:  pdfgen.FontConfig{},  // default: Arial 11pt Bold
}
```

| Field        | Type         | Default  | Notes                                 |
|--------------|--------------|----------|---------------------------------------|
| `Items`      | `[]InfoItem` | —        | `{Label string; Value string}`        |
| `Columns`    | `int`        | `2`      | Items per row                         |
| `ShowBorder` | `bool`       | `false`  | Border around each cell               |
| `Width`      | `float64`    | `0`      | 0 = full usable width                 |
| `LabelFont`  | `FontConfig` | 8pt      | Small gray label above value          |
| `ValueFont`  | `FontConfig` | 11pt Bold| Large bold value below label          |

---

### 4. `SectionLabelComponent` — section header row

A single text line with an optional right-side label. When `RightText` contains `:`, everything before it is secondary color and everything after is accent color.

```go
&pdfgen.SectionLabelComponent{
    LeftText:  "Total distance per state",
    RightText: "Total Distance: 7,000 mi",  // colon splits colors
    // RightText: "6 records",              // no colon = plain secondary text
    MarginBottom: 2,  // mm below; default 2

    // Optional overrides:
    LeftFont:   pdfgen.FontConfig{},  // default: bold
    RightFont:  pdfgen.FontConfig{},  // default: normal
    RightColor: pdfgen.Color{},       // default: theme AccentColor (for value after ":")
}
```

Output with colon split:
```
Total distance per state              Total Distance:  7,000 mi
                                      ───────────────  ─────────
                                      secondary gray   accent blue
```

---

### 5. `SpacerComponent` — vertical whitespace

```go
&pdfgen.SpacerComponent{Height: 6}  // 6 mm gap
```

---

### 6. `TableComponent` — data table

The most complex component. Before adding, understand two fields that default to `false` in Go but should typically be `true`:

> **Always set `ShowHeader: true` and `RowStriping: true` explicitly.**

```go
&pdfgen.TableComponent{
    ShowHeader:   true,   // render column header row
    RowStriping:  true,   // alternate row background colors
    CellPaddingH: 3,      // mm horizontal padding; default 3
    CellPaddingV: 2,      // mm vertical padding; default 2
    MinRowHeight:  8,     // mm minimum row height; default 8
    BorderStyle:  "all",  // "all" | "outer" | "none"; default "all"

    // Optional font overrides:
    HeaderFont: pdfgen.FontConfig{},  // default: Bold
    RowFont:    pdfgen.FontConfig{},  // default: Normal

    Columns: []pdfgen.ColumnDef{
        {Header: "No",       Width: 15, Align: "C"},
        {Header: "State",    Width: 55, Align: "L"},
        {Header: "Miles",    Width: 0,  Align: "R"},  // Width:0 = fill remaining space
    },
    Rows: [][]string{
        {"1", "Tennessee", "1,200"},
        {"2", "Kentucky",  "850"},
    },
}
```

#### `ColumnDef` fields

| Field         | Type          | Default | Notes                                              |
|---------------|---------------|---------|----------------------------------------------------|
| `Header`      | `string`      | —       | Column header text                                 |
| `Width`       | `float64`     | `0`     | mm; **0 = shares remaining space equally**         |
| `Align`       | `string`      | `"L"`   | `"L"` left \| `"C"` center \| `"R"` right         |
| `Overflow`    | `OverflowMode`| `OverflowWrap` | See below                                 |
| `HeaderAlign` | `string`      | =Align  | Override alignment for header cell only            |
| `Bold`        | `bool`        | `false` | Render cell content bold                           |

#### Column width rules

- `Width > 0` → fixed mm width
- `Width == 0` → column takes an **equal share of remaining space** after fixed columns
- Multiple `Width: 0` columns → each gets equal slice of what's left
- Widths do not need to sum to page width — use `Width: 0` for at least one column

**Example: 4 columns, mixed widths on A4 portrait (180mm usable)**

```go
// Fixed: 15 + 55 + 30 = 100mm used. Remaining: 80mm → one Width:0 column gets 80mm.
{Header: "No",    Width: 15, Align: "C"},
{Header: "State", Width: 55, Align: "L"},
{Header: "City",  Width: 30, Align: "L"},
{Header: "Miles", Width: 0,  Align: "R"},  // gets 80mm
```

#### `OverflowMode`

| Constant           | Behaviour                                                  |
|--------------------|------------------------------------------------------------|
| `OverflowWrap`     | Text wraps; row grows taller to fit all lines *(default)*  |
| `OverflowTruncate` | Text is clipped and `…` is appended to fit the column      |

```go
{Header: "VIN", Width: 55, Align: "L", Overflow: pdfgen.OverflowTruncate}
{Header: "Notes", Width: 60, Align: "L", Overflow: pdfgen.OverflowWrap}
```

#### `BorderStyle`

| Value     | Effect                            |
|-----------|-----------------------------------|
| `"all"`   | Border around every cell (grid)   |
| `"outer"` | Border around the whole row only  |
| `"none"`  | No borders                        |

#### Auto page break

When a row would overflow the page, a new page is added automatically. If `ShowHeader: true`, the header row is re-rendered at the top of the new page.

---

### 7. `GroupedTableComponent` — labeled table section

Convenience wrapper: renders a `SectionLabelComponent` then a `TableComponent` then a spacer. Use this instead of composing the three manually.

```go
&pdfgen.GroupedTableComponent{
    Label:       "Total distance per state",   // left text of section label
    BadgeText:   "Total Distance: 7,000 mi",   // right text (supports ":" color split)
    SpacerAfter: 6,                            // mm after table; default 6

    Table: pdfgen.TableComponent{
        ShowHeader:  true,
        RowStriping: true,
        Columns: []pdfgen.ColumnDef{
            {Header: "No",       Width: 15, Align: "C"},
            {Header: "State",    Width: 55, Align: "L"},
            {Header: "Distance", Width: 0,  Align: "R"},
        },
        Rows: [][]string{
            {"1", "Tennessee", "1,200 mi"},
            {"2", "Kentucky",  "850 mi"},
        },
    },
}
```

---

### 8. `FooterComponent` — page footer (register, don't Add)

Call `doc.SetFooter(...)` **before** `doc.Add(...)`. The footer renders automatically on every page.

```go
doc.SetFooter(&pdfgen.FooterComponent{
    LeftText:   "Confidential",           // static, left side
    CenterText: "Page {page} of {total}", // {page} and {total} are placeholders
    RightText:  "Lucid ELD",             // static, right side
    ShowBorder: true,                    // thin top border line above footer
    Font:       pdfgen.FontConfig{},     // default: 8pt
    TextColor:  pdfgen.Color{},          // default: theme SecondaryText
})
```

| Field        | Type         | Default           | Notes                                      |
|--------------|--------------|-------------------|--------------------------------------------|
| `LeftText`   | `string`     | —                 | Left-aligned static text                   |
| `CenterText` | `string`     | —                 | Centered; `{page}` and `{total}` replaced  |
| `RightText`  | `string`     | —                 | Right-aligned static text                  |
| `ShowBorder` | `bool`       | `false`           | Thin line above footer                     |
| `Font`       | `FontConfig` | 8pt               | Footer font                                |
| `TextColor`  | `Color`      | SecondaryText     | Footer text color                          |

**Placeholders:**
- `{page}` → current page number
- `{total}` → total page count

---

## Complete Patterns

### IFTA Report

```go
doc := pdfgen.New(pdfgen.DocumentConfig{PageSize: "A4", Theme: pdfgen.DefaultTheme()})

doc.SetFooter(&pdfgen.FooterComponent{
    CenterText: "Page {page} of {total}",
    RightText:  "Lucid ELD",
    ShowBorder: true,
})

doc.Add(
    &pdfgen.LogoComponent{ImagePath: "logo.png", Width: 42, Height: 17, Position: "top-right"},
    &pdfgen.HeaderComponent{
        Title:    "IFTA REPORT",
        Subtitle: "QGM EXPRESS",
        Lines:    []string{"1234 Trucker Lane, Nashville TN", "Q1 2024 | Jan 1 – Mar 31"},
    },
    &pdfgen.SpacerComponent{Height: 5},
    &pdfgen.InfoBlockComponent{
        Items:      []pdfgen.InfoItem{{Label: "Total Vehicle", Value: "3"}, {Label: "Total Distance", Value: "21,450 mi"}},
        Columns:    2,
        ShowBorder: true,
    },
    &pdfgen.SpacerComponent{Height: 7},
    &pdfgen.GroupedTableComponent{
        Label:     "Distance per state",
        BadgeText: "Total: 21,450 mi",
        Table: pdfgen.TableComponent{
            ShowHeader: true, RowStriping: true,
            Columns: []pdfgen.ColumnDef{
                {Header: "State",    Width: 60, Align: "L"},
                {Header: "Distance", Width: 0,  Align: "R"},
            },
            Rows: [][]string{{"Tennessee", "7,200 mi"}, {"Kentucky", "6,850 mi"}},
        },
    },
)
doc.Save("ifta.pdf")
```

### Invoice

```go
doc := pdfgen.New(pdfgen.DocumentConfig{PageSize: "A4", Theme: pdfgen.DefaultTheme()})

doc.SetFooter(&pdfgen.FooterComponent{
    LeftText:   "Thank you for your business",
    RightText:  "Page {page} of {total}",
    ShowBorder: true,
})

doc.Add(
    &pdfgen.LogoComponent{ImagePath: "logo.png", Width: 50, Position: "top-right"},
    &pdfgen.HeaderComponent{
        Title:    "INVOICE",
        Subtitle: "Acme Trucking LLC",
        Lines:    []string{"Invoice #: INV-2024-0042", "Date: April 1, 2024", "Due: April 30, 2024"},
    },
    &pdfgen.SpacerComponent{Height: 4},
    &pdfgen.InfoBlockComponent{
        Items: []pdfgen.InfoItem{
            {Label: "Bill To",    Value: "Customer Corp"},
            {Label: "Ship To",   Value: "123 Main St"},
            {Label: "Terms",     Value: "Net 30"},
            {Label: "PO Number", Value: "PO-9981"},
        },
        Columns: 2, ShowBorder: true,
    },
    &pdfgen.SpacerComponent{Height: 6},
    &pdfgen.GroupedTableComponent{
        Label: "Line Items",
        Table: pdfgen.TableComponent{
            ShowHeader: true, RowStriping: true,
            Columns: []pdfgen.ColumnDef{
                {Header: "Description", Width: 0,  Align: "L", Overflow: pdfgen.OverflowWrap},
                {Header: "Qty",         Width: 20, Align: "C"},
                {Header: "Unit Price",  Width: 35, Align: "R"},
                {Header: "Total",       Width: 35, Align: "R"},
            },
            Rows: [][]string{
                {"Freight service Nashville → Chicago", "1", "$1,200.00", "$1,200.00"},
                {"Fuel surcharge", "1", "$180.00", "$180.00"},
                {"", "", "Subtotal", "$1,380.00"},
                {"", "", "Tax (0%)", "$0.00"},
                {"", "", "Total Due", "$1,380.00"},
            },
        },
    },
)
doc.Save("invoice.pdf")
```

---

## Agent Cheatsheet

### "How do I…"

| Task | Answer |
|---|---|
| Set page to landscape | `Orientation: "landscape"` in `DocumentConfig` |
| Change accent color | `theme.AccentColor = pdfgen.Color{R,G,B}` |
| Make a column fill remaining width | `Width: 0` in `ColumnDef` |
| Truncate long text in a cell | `Overflow: pdfgen.OverflowTruncate` |
| Wrap text in a cell (taller rows) | `Overflow: pdfgen.OverflowWrap` *(default)* |
| Add vertical space | `&pdfgen.SpacerComponent{Height: N}` |
| Show page numbers | `doc.SetFooter(...)` with `CenterText: "Page {page} of {total}"` |
| Put logo top-right | `&pdfgen.LogoComponent{Position: "top-right", ...}` |
| Make header row bold | It's bold by default; override with `HeaderFont: FontConfig{Style: "B"}` |
| Remove table borders | `BorderStyle: "none"` |
| Add a section label above a table | Use `GroupedTableComponent` instead of `TableComponent` directly |
| Get the PDF as bytes | `data, err := doc.Bytes()` |
| Repeat header on new page | Automatic when `ShowHeader: true` |

### Common Mistakes

```
❌  &pdfgen.TableComponent{Columns: [...], Rows: [...]}
✅  &pdfgen.TableComponent{ShowHeader: true, RowStriping: true, Columns: [...], Rows: [...]}
    // ShowHeader and RowStriping default to false in Go — always set them explicitly

❌  doc.Add(&pdfgen.FooterComponent{...})
✅  doc.SetFooter(&pdfgen.FooterComponent{...})
    // Footer must be registered via SetFooter, not Add

❌  &pdfgen.LogoComponent{Width: 0, ...}
✅  &pdfgen.LogoComponent{Width: 42, ...}
    // Width must be > 0

❌  Column widths that sum to more than the usable page width
✅  Leave at least one column with Width: 0 to absorb remaining space safely

❌  doc.Add(logo, header) where logo is added after header
✅  doc.Add(logo, header)  // logo first — both render at top, logo floats
```

### Page size reference

| Size     | Orientation | Usable width (default 15mm margins) |
|----------|-------------|--------------------------------------|
| A4       | portrait    | 180 mm                               |
| A4       | landscape   | 267 mm                               |
| Letter   | portrait    | 186 mm (8.5in − 2×15mm)              |
| Letter   | landscape   | 246 mm (11in − 2×15mm)               |
