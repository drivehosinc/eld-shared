# eld-shared

Shared Go packages used across Lucid ELD services.

## Import

```bash
go get github.com/drivehosinc/eld-shared@latest
```

Then import the specific package you need:

```go
import "github.com/drivehosinc/eld-shared/pdfgen"
import "github.com/drivehosinc/eld-shared/logger"
import "github.com/drivehosinc/eld-shared/pagination"
// etc.
```

## Packages

| Package | Description |
|---|---|
| [`pdfgen`](./pdfgen) | Component-based PDF builder for business reports (IFTA, movement, invoices) |
| [`logger`](./logger) | Structured logger wrapping `slog` with service name, level, and attribute options |
| [`pagination`](./pagination) | gRPC-compatible pagination request/response calculator with validation |
| [`metadata`](./metadata) | gRPC metadata key constants and context helpers |
| [`argo2id`](./argo2id) | Argon2id password hashing and verification |
| [`random`](./random) | Cryptographically secure random byte generation |
| [`postgresql`](./postgresql) | PostgreSQL connection helpers via `pgx/v5` |

## Quick Examples

### pdfgen

Full docs: [`pdfgen/README.md`](./pdfgen/README.md)

```go
doc := pdfgen.New(pdfgen.DocumentConfig{PageSize: "A4", Orientation: "portrait"})
doc.Add(
    &pdfgen.HeaderComponent{Title: "IFTA REPORT", Subtitle: "Company Name"},
    &pdfgen.TableComponent{
        ShowHeader: true, RowStriping: true, BorderStyle: "columns",
        Columns: []pdfgen.ColumnDef{{Header: "State", Width: 0}, {Header: "Miles", Width: 40}},
        Rows:    [][]string{{"California", "1,240 mi"}},
    },
)
doc.Save("report.pdf")
```

### logger

```go
log := logger.New(
    logger.WithServiceName("my-service"),
    logger.WithLevel(slog.LevelInfo),
)
log.Info("server started", "port", 8080)
```

### pagination

```go
calc := pagination.NewCalculator()
resp := calc.Calculate(pagination.FromGRPCRequest(req.Page, req.Limit), totalCount)
result := pagination.CreateResult(items, resp)
```

### argo2id

```go
hash, err := argo2id.CreateHashFromPassword(password, &argo2id.Params{...})
match, err := argo2id.ComparePasswordAndHash(password, hash)
```

## Requirements

- Go 1.24.5+

## Versioning

This module follows [Semantic Versioning](https://semver.org).

| Version | Description |
|---|---|
| `v1.3.4` | Added `pdfgen` — PDF generation package |

To pin a specific version:

```bash
go get github.com/drivehosinc/eld-shared@v1.3.4
```
