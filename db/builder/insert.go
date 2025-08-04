// File: db/builder/insert.go

package builder

import (
	"fmt"
	"strings"

	"github.com/entiqon/db/driver"
	"github.com/entiqon/db/internal/core/builder/bind"
	core "github.com/entiqon/db/internal/core/errors"
	token2 "github.com/entiqon/db/internal/core/token"
)

// InsertBuilder builds a SQL INSERT statement.
//
// It supports inserting data into a table with specified columns and values,
// and can optionally append a RETURNING clause (PostgreSQL).
type InsertBuilder struct {
	BaseBuilder
	table     string
	columns   []token2.FieldToken
	values    [][]any
	returning []token2.FieldToken
}

// NewInsert creates a new InsertBuilder using the given SQL dialect.
//
// If the provided dialect is nil, it defaults to driver.NewGenericDialect().
// The builder name is automatically set to "insert".
//
// Since: v1.4.0
func NewInsert(dialect driver.Dialect) *InsertBuilder {
	base := NewBaseBuilder("insert", dialect)

	return &InsertBuilder{
		BaseBuilder: base,
		columns:     []token2.FieldToken{},
		values:      [][]any{},
		returning:   []token2.FieldToken{},
	}
}

// Into sets the target table for the INSERT operation.
func (b *InsertBuilder) Into(table string) *InsertBuilder {
	if table == "" {
		b.AddStageError("INTO", fmt.Errorf("requires a target table"))
	} else {
		b.table = table
	}
	return b
}

// Columns sets the column names for the INSERT statement.
//
// All provided column names must be plain identifiers (e.g., "id", "name").
// Column aliasing is not allowed in INSERT statements. Any column passed in the
// form "name AS alias" will be rejected, logged as an error under the "COLUMNS"
// stage, and skipped.
//
// Example:
//     Columns("id", "name")             // ✅ valid
//     Columns("id", "name AS n")        // ❌ invalid, will be rejected
//
// If any alias is detected, the builder will collect an error and the aliased column
// will not be added to the insert operation.
//
// Since: v1.4.0

func (b *InsertBuilder) Columns(names ...string) *InsertBuilder {
	b.columns = []token2.FieldToken{}
	for _, name := range names {
		f := token2.Field(name)
		if f.Alias != "" {
			b.AddStageError("COLUMNS", fmt.Errorf("column aliasing is not allowed: '%s AS %s'", f.Name, f.Alias))
			continue
		}
		b.columns = append(b.columns, f)
	}
	return b
}

// Values appends a row of values using a map of column name to value.
// Each call adds a new row. The map must contain every column defined in Columns().
func (b *InsertBuilder) Values(row ...any) *InsertBuilder {
	b.values = append(b.values, row)
	return b
}

// Returning adds one or more column names to the RETURNING clause.
// It parses string expressions into FieldTokens.
// If called multiple times, it appends to the existing list.
func (b *InsertBuilder) Returning(fields ...string) *InsertBuilder {
	for _, f := range fields {
		b.returning = append(b.returning, token2.FieldsFromExpr(f)...)
	}
	return b
}

// Build compiles the full INSERT SQL statement along with arguments.
// Returns an error if the structure is invalid or values are missing.
// Build compiles the full INSERT SQL statement along with arguments.
// Returns an error if the structure is invalid or values are missing.
func (b *InsertBuilder) Build() (string, []any, error) {
	return b.buildQuery(len(b.returning) > 0)
}

// BuildInsertOnly compiles a full INSERT SQL statement with arguments,
// excluding any RETURNING clause.
//
// This method is useful when you only need the INSERT operation itself,
// without fetching results (e.g., for engines that don't support RETURNING).
//
// It validates the following before building:
//   - A target table must be set using .Into()
//   - At least one column must be defined via .Columns()
//   - At least one row of values must be provided using .Values()
//   - Each row must match the number of defined columns
//
// The dialect determines how placeholders (e.g., ?, $1, :name) and identifiers are rendered.
//
// Example output (Postgres dialect):
//
//	INSERT INTO "users" ("id", "name") VALUES ($1, $2)
//
// Returns the SQL string, the bound arguments, or an error.
//
// Since: v1.4.0
func (b *InsertBuilder) BuildInsertOnly() (string, []any, error) {
	return b.buildQuery(false)
}

func (b *InsertBuilder) buildQuery(withReturning bool) (string, []any, error) {
	if b.table == "" {
		b.Validator.AddStageError(core.StageFrom, fmt.Errorf("requires a target table"))
	}
	if len(b.columns) == 0 {
		b.Validator.AddStageError(core.StageInto, fmt.Errorf("at least one column is required"))
	}
	if len(b.values) == 0 {
		b.Validator.AddStageError(core.StageValues, fmt.Errorf("at least one set of values is required"))
	}

	if withReturning && !b.Dialect.SupportsReturning() {
		b.Validator.AddStageError(core.StageReturning, fmt.Errorf("at least one set of values is required"))
	}

	colCount := len(b.columns)
	binder := bind.NewParamBinder(b.Dialect)

	var args []any
	var rowPlaceholders []string
	quotedCols := make([]string, len(b.columns))
	for i, col := range b.columns {
		quotedCols[i] = b.Dialect.QuoteIdentifier(col.Name)
	}

	for i, row := range b.values {
		if len(row) != colCount {
			b.Validator.AddStageError(
				core.StageReturning,
				fmt.Errorf("row %d has %d values, expected %d", i+1, len(row), colCount),
			)
		}
		placeholders := binder.BindMany(row...)
		args = append(args, row...)
		rowPlaceholders = append(rowPlaceholders, fmt.Sprintf("(%s)", strings.Join(placeholders, ", ")))
	}

	if err := b.Validate(); err != nil {
		return "", nil, err
	}

	tokens := []string{
		"INSERT INTO",
		b.Dialect.RenderFrom(b.table, ""),
		fmt.Sprintf("(%s)", strings.Join(quotedCols, ", ")),
		"VALUES", strings.Join(rowPlaceholders, ", "),
	}

	if withReturning {
		returnCols := make([]string, len(b.returning))
		for i, col := range b.returning {
			returnCols[i] = b.Dialect.QuoteIdentifier(col.Name)
		}
		tokens = append(tokens, "RETURNING", strings.Join(returnCols, ", "))
	}

	return strings.Join(tokens, " "), args, nil
}
