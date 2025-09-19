/*
Package generic provides the ANSI-compliant SQL dialect implementation.

# Overview

The generic dialect is a safe default for generating SQL statements when no
vendor-specific behavior is required or when the target database is unknown.
It follows the ANSI SQL standard as closely as possible:

  - Identifiers are quoted using double quotes (").
  - Placeholders are always "?" and are not positional.
  - Pagination is rendered using LIMIT and OFFSET clauses.
  - RETURNING, MERGE, and UPSERT clauses are not supported.

# Usage

The generic dialect is not instantiated directly. Instead, call New() to obtain
a ready-to-use instance that implements the dialect.SQLDialect interface:

	d := generic.New()
	sql := fmt.Sprintf(
	    "SELECT %s FROM %s%s",
	    d.QuoteIdentifier("id"),
	    d.QuoteIdentifier("users"),
	    d.PaginationSyntax(10, 0),
	)
	// Produces: SELECT "id" FROM "users" LIMIT 10

# Capabilities

The generic dialect advertises its supported features through the Options()
method. This allows builders to conditionally include clauses:

	opts := d.Options()
	fmt.Println(opts.SupportsCTE)          // true
	fmt.Println(opts.EnableReturning)      // false
	fmt.Println(opts.AllowUpsert)          // false

# When to Use

Use this dialect as:
  - A fallback when no database-specific dialect is configured.
  - A baseline for testing SQL builders in a dialect-agnostic way.
  - A template for implementing new dialects (e.g., Postgres, MySQL).

For vendor-specific features such as RETURNING (Postgres, Oracle) or
vendor-specific placeholder formats, use the corresponding dialect package.
*/
package generic
