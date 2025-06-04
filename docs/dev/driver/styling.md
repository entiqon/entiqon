# SQL Styling System

The Entiqon driver layer includes a modular styling system for rendering SQL output consistently across supported dialects.

## Components

| Style Type         | Enum Source                          | Used In                      | Example Output        |
|--------------------|--------------------------------------|-------------------------------|------------------------|
| `AliasStyle`       | `driver/styling/alias_style.go`      | `RenderTable`, `RenderColumn`| `users AS u`           |
| `QuoteStyle`       | `driver/styling/quote_style.go`      | `QuoteIdentifier`            | `"users"`              |
| `PlaceholderStyle` | `driver/styling/placeholder_style.go`| `RenderPlaceholder`          | `$1`, `?`, `:id`       |

## Behavior

Each style has both basic and dialect-aware rendering support.

### AliasStyle

- `Format(base, alias)` — local style-based formatting
- `FormatWith(quoter, base, alias)` — dialect-aware using `IdentifierQuoter`

### QuoteStyle

- `Quote(identifier)` — applies configured quoting style (e.g., `"user"`, `[user]`)

### PlaceholderStyle

- `Format(index)` — for positional placeholders
- `FormatNamed(expr)` — for named placeholders

## Integration

These styles are used by `BaseDialect`, which supports:

- `TableAliasStyle`
- `ColumnAliasStyle`
- `QuoteStyle`
- `PlaceholderStyle`

## Since

Introduced and normalized in **v1.5.0**
