/*
Package dialect defines the core SQL dialect contract and capability matrix.

# Overview

A dialect is responsible for rendering vendor-specific SQL fragments such as:
  - Identifier quoting
  - Parameter placeholders
  - Pagination (LIMIT/OFFSET)
  - Feature availability (RETURNING, MERGE, UPSERT, CTEs, window functions)

The root package declares:

  - SQLDialect — the shared interface every dialect must implement
  - Options    — the capability matrix advertised by each dialect

Concrete implementations live in subpackages such as:

  - generic   — ANSI-compliant fallback
  - postgres  — PostgreSQL-specific rules
  - mysql     — MySQL-specific rules
  - mariadb   — MariaDB-specific rules
  - sqlite    — SQLite-specific rules
  - mssql     — Microsoft SQL Server-specific rules
  - oracle    — Oracle-specific rules
  - db2       — IBM DB2-specific rules
  - firebird  — Firebird-specific rules
  - informix  — Informix-specific rules
  - cockroach — CockroachDB (Postgres-compatible)
  - tidb      — TiDB (MySQL-compatible)
  - hana      — SAP HANA-specific rules
  - snowflake — Snowflake SQL
  - redshift  — Amazon Redshift
  - teradata  — Teradata SQL
  - clickhouse — ClickHouse SQL-like syntax

# SQLDialect

The SQLDialect interface defines the minimum set of behaviors required:

	type SQLDialect interface {
	    Name() string
	    Options() Options
	    QuoteIdentifier(name string) string
	    QuoteLiteral(literal any) string
	    PaginationSyntax(limit, offset int) string
	    Placeholder(index int) string
	}

# Options

The Options struct advertises the capabilities of each dialect:

	type Options struct {
	    Name                  string
	    QuoteStyle            string
	    PlaceholderStyle      string
	    AllowMerge            bool
	    AllowUpsert           bool
	    ForcedAliasing        bool
	    EnableReturning       bool
	    SupportsCTE           bool
	    SupportsWindowFunctions bool
	    MaxPlaceholderIndex   int
	}

# Usage

Clients typically obtain a dialect from a subpackage:

	d := generic.New()
	sql := fmt.Sprintf("SELECT %s FROM %s%s",
	    d.QuoteIdentifier("id"),
	    d.QuoteIdentifier("users"),
	    d.PaginationSyntax(10, 0),
	)
	// SELECT "id" FROM "users" LIMIT 10
*/
package dialect
