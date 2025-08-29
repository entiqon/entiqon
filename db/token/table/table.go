package table

import (
	"fmt"
	"strings"

	"github.com/entiqon/entiqon/db/token"
)

// Table represents a SQL table token.
//
// A Table encapsulates the user-provided input, the normalized
// table name, and an optional alias. It exposes multiple forms:
//
//   - Raw(): generic SQL fragment (dialect-agnostic).
//   - Render(): canonical, machine-facing representation for query building.
//   - String(): human-facing representation for logs and audits.
//   - Debug(): developer-facing internal state.
//
// Table also supports semantic cloning via contract.Clonable,
// error inspection via contract.Errorable, and raw-state
// inspection via contract.Rawable.
type Table struct {
	kind token.ExpressionKind
	// input is the exact user-provided string, never modified.
	input string

	// name is the normalized base table name (derived from input).
	name string

	// alias is the optional alias (parsed or explicitly set).
	alias string

	// err holds any construction or validation error.
	err error

	// isRaw reports whether the table was constructed via the
	// explicit two-argument form or is a subquery.
	isRaw bool
}

// New constructs a new Table from user input.
//
// Accepted forms:
//   - table.New("users")          → name="users"
//   - table.New("users u")        → name="users", alias="u"
//   - table.New("users AS u")     → name="users", alias="u"
//   - table.New("users", "u")     → name="users", alias="u", isRaw=true
//   - table.New("(SELECT ...)", "t")
//     → name="(SELECT ...)", alias="t", isRaw=true
//   - table.New("(SELECT ...) AS t")
//     → name="(SELECT ...)", alias="t", isRaw=true
//
// The first argument is always preserved in input (verbatim or normalized).
// If construction fails, the returned Table is errored but still
// carries the original input for diagnostics.
func New(args ...string) Token {
	t := &Table{}
	if len(args) > 0 {
		t.input = strings.TrimSpace(args[0])
	}

	switch len(args) {
	case 0:
		t.err = fmt.Errorf("requires at least one argument")
		return t

	case 1:
		raw := strings.TrimSpace(args[0])
		if raw == "" {
			t.err = fmt.Errorf("empty table name")
			return t
		}
		t.input = raw

		// Subquery form
		if strings.HasPrefix(raw, "(") {
			upper := strings.ToUpper(raw)
			if strings.Contains(upper, " AS ") {
				idx := strings.LastIndex(upper, " AS ")
				t.name = strings.TrimSpace(raw[:idx])
				t.alias = strings.TrimSpace(raw[idx+4:])
				t.isRaw = true
			} else {
				t.name = raw
				t.isRaw = true // ✅ ensure all subqueries are marked raw
				t.err = fmt.Errorf("subquery source must have an alias")
			}
			return t
		}

		// Plain table parsing
		parts := strings.Fields(raw)
		switch len(parts) {
		case 1:
			t.name = parts[0]
		case 2:
			if strings.EqualFold(parts[1], "AS") {
				t.err = fmt.Errorf("invalid format %q", raw)
			} else {
				t.name, t.alias = parts[0], parts[1]
			}
		case 3:
			if strings.EqualFold(parts[1], "AS") {
				t.name, t.alias = parts[0], parts[2]
			} else {
				t.err = fmt.Errorf("invalid format %q", raw)
			}
		default:
			t.err = fmt.Errorf("too many tokens in %q", raw)
		}

	case 2:
		name := strings.TrimSpace(args[0])
		alias := strings.TrimSpace(args[1])
		if name == "" || alias == "" {
			t.input = strings.Join(args, " ")
			t.err = fmt.Errorf("table and alias must be non-empty")
			return t
		}

		t.name, t.alias, t.isRaw = name, alias, true
		t.input = fmt.Sprintf("%s AS %s", name, alias)

	default:
		t.input = strings.Join(args, " ")
		t.err = fmt.Errorf("too many arguments")
	}

	return t
}

func (t *Table) ExpressionKind() token.ExpressionKind {
	return t.kind
}

func (t *Table) Input() string {
	return t.input
}

func (t *Table) Expr() string {
	return t.name
}

func (t *Table) Name() string {
	return t.name
}

func (t *Table) Alias() string {
	return t.name
}

// IsAliased reports whether the table has an alias.
func (t *Table) IsAliased() bool { return t.alias != "" }

// Clone returns a semantic copy of the Table.
func (t *Table) Clone() Token {
	return &Table{
		input: t.input,
		name:  t.name,
		alias: t.alias,
		err:   t.err,
		isRaw: t.isRaw,
	}
}

// Debug returns a developer-facing representation of the Table.
//
// The output is verbose and intended for diagnostics, showing the
// original input and internal state flags (raw, aliased, errored).
// If the Table is errored, Debug also appends the error message
// inside a separate { } block.
func (t *Table) Debug() string {
	flags := fmt.Sprintf(
		"[raw:%v, aliased:%v, errored:%v]",
		t.IsRaw(),
		t.IsAliased(),
		t.IsErrored(),
	)
	if !t.IsValid() {
		return fmt.Sprintf("❌ Table(%q): %s {err=%v}", t.input, flags, t.err)
	}
	return fmt.Sprintf("✅ Table(%q): %s", t.input, flags)
}

// IsErrored reports whether the Table was constructed with an error.
func (t *Table) IsErrored() bool { return t.err != nil }

// Error returns the underlying construction error, if any.
func (t *Table) Error() error { return t.err }

// SetError assigns an error to the table. Intended for use during
// construction/parsing to capture validation failures.
func (t *Table) SetError(err error) Token {
	t.err = err
	return t
}

// IsRaw reports whether the Table was explicitly constructed as raw
// (via the two-argument form or as a subquery).
func (t *Table) IsRaw() bool { return t.isRaw }

// Raw returns the generic SQL fragment of the Table, including alias if present.
//
// Unlike Render(), Raw() does not apply dialect-specific quoting or rewriting.
// It simply reflects the normalized fragment as a plain SQL string.
func (t *Table) Raw() string {
	if !t.IsValid() {
		return ""
	}
	if t.alias != "" {
		return fmt.Sprintf("%s AS %s", t.name, t.alias)
	}
	return t.name
}

// Render returns the canonical SQL representation of the Table.
//
// Render() differs from Raw() in that it represents the resolved form
// the builder will actually use. It may later be extended to apply
// dialect-specific quoting or rewriting.
//
// If the Table is invalid or errored, Render() returns an empty string.
func (t *Table) Render() string {
	if !t.IsValid() {
		return ""
	}
	if t.alias != "" {
		return fmt.Sprintf("%s AS %s", t.name, t.alias)
	}
	return t.name
}

// String returns the human-facing representation of the Table.
//
// Unlike Render(), which is used in query building, String() is
// intended for logs and audits. It produces a concise summary of
// the table state.
//
// If the Table is invalid or errored, String() reports the user input
// and the associated error. If valid, it reports the normalized form
// with alias if present, prefixed with a ✅ marker.
func (t *Table) String() string {
	if !t.IsValid() {
		// Always show original input and error
		return fmt.Sprintf("❌ Table(%q): %v", t.input, t.err)
	}
	if t.alias != "" {
		return fmt.Sprintf("✅ Table(%s AS %s)", t.name, t.alias)
	}
	return fmt.Sprintf("✅ Table(%s)", t.name)
}

// IsValid reports whether the table has a non-empty name and no error.
func (t *Table) IsValid() bool { return !t.IsErrored() && t.name != "" }
