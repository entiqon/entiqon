package operator

import (
	"sort"
	"strings"
)

// Type classifies a SQL comparison or predicate operator used in WHERE, ON, or
// HAVING clauses. The String() form of a Type is its canonical SQL spelling
// (e.g., "IN", "BETWEEN", ">=", "IS NOT NULL").
//
// Example:
//
//	var t Type = In
//	fmt.Println(t.String()) // "IN"
type Type int

const (
	// Invalid marks an unrecognized or unsupported operator.
	Invalid Type = iota

	// Equal represents the "=" operator.
	//
	// Example:
	//   Condition: "status = ?"
	//   Operator:  Equal
	//   Output:    "="
	Equal

	// NotEqual represents the "!=" operator (the parser also accepts "<>").
	//
	// Example:
	//   Condition: "status != ?"
	//   Operator:  NotEqual
	//   Output:    "!="
	NotEqual

	// GreaterThan represents the ">" operator.
	//
	// Example:
	//   Condition: "age > ?"
	//   Operator:  GreaterThan
	//   Output:    ">"
	GreaterThan

	// GreaterThanOrEqual represents the ">=" operator.
	//
	// Example:
	//   Condition: "age >= ?"
	//   Operator:  GreaterThanOrEqual
	//   Output:    ">="
	GreaterThanOrEqual

	// LessThan represents the "<" operator.
	//
	// Example:
	//   Condition: "age < ?"
	//   Operator:  LessThan
	//   Output:    "<"
	LessThan

	// LessThanOrEqual represents the "<=" operator.
	//
	// Example:
	//   Condition: "age <= ?"
	//   Operator:  LessThanOrEqual
	//   Output:    "<="
	LessThanOrEqual

	// In represents the "IN" operator.
	//
	// Example:
	//   Condition: "id IN (?)"
	//   Operator:  In
	//   Output:    "IN"
	In

	// NotIn represents the "NOT IN" operator.
	//
	// Example:
	//   Condition: "id NOT IN (?)"
	//   Operator:  NotIn
	//   Output:    "NOT IN"
	NotIn

	// Between represents the "BETWEEN" operator.
	//
	// Example:
	//   Condition: "created BETWEEN ? AND ?"
	//   Operator:  Between
	//   Output:    "BETWEEN"
	Between

	// Like represents the "LIKE" operator.
	//
	// Example:
	//   Condition: "name LIKE ?"
	//   Operator:  Like
	//   Output:    "LIKE"
	Like

	// NotLike represents the "NOT LIKE" operator.
	//
	// Example:
	//   Condition: "name NOT LIKE ?"
	//   Operator:  NotLike
	//   Output:    "NOT LIKE"
	NotLike

	// IsNull represents the "IS NULL" predicate.
	//
	// Example:
	//   Condition: "deleted_at IS NULL"
	//   Operator:  IsNull
	//   Output:    "IS NULL"
	IsNull

	// IsNotNull represents the "IS NOT NULL" predicate.
	//
	// Example:
	//   Condition: "deleted_at IS NOT NULL"
	//   Operator:  IsNotNull
	//   Output:    "IS NOT NULL"
	IsNotNull

	// IsDistinctFrom represents the "IS DISTINCT FROM" predicate.
	//
	// Example (PostgreSQL):
	//   Condition: "a IS DISTINCT FROM b"
	//   Operator:  IsDistinctFrom
	//   Output:    "IS DISTINCT FROM"
	IsDistinctFrom

	// NotIsDistinctFrom represents the "IS NOT DISTINCT FROM" predicate.
	//
	// Example (PostgreSQL):
	//   Condition: "a IS NOT DISTINCT FROM b"
	//   Operator:  NotIsDistinctFrom
	//   Output:    "IS NOT DISTINCT FROM"
	NotIsDistinctFrom
)

// Meta describes a supported operator: its canonical spelling, alias, and a
// deterministic position for GetKnownOperators ordering (lower = earlier).
// Synonyms are accepted by ParseFrom in addition to String and Alias.
type Meta struct {
	String   string   // canonical SQL spelling (render form)
	Alias    string   // short stable mnemonic (logs/JSON/flags)
	Position int      // deterministic order key (lower ranks first)
	Synonyms []string // additional parse tokens (case/space-insensitive)
}

// registry declares all supported operators in one place.
// Positions are chosen to match the expected order:
//
// [IS NOT DISTINCT FROM IS DISTINCT FROM IS NOT NULL NOT LIKE BETWEEN IS NULL NOT IN LIKE IN != >= <= > < =]
var registry = map[Type]Meta{
	NotIsDistinctFrom:  {String: "IS NOT DISTINCT FROM", Alias: "notdistinct", Position: 1, Synonyms: []string{"is not distinct from", "notdistinct"}},
	IsDistinctFrom:     {String: "IS DISTINCT FROM", Alias: "isdistinct", Position: 2, Synonyms: []string{"is distinct from", "isdistinct"}},
	IsNotNull:          {String: "IS NOT NULL", Alias: "notnull", Position: 3, Synonyms: []string{"is not null", "notnull"}},
	NotLike:            {String: "NOT LIKE", Alias: "nlike", Position: 4, Synonyms: []string{"not like", "nlike"}},
	Between:            {String: "BETWEEN", Alias: "between", Position: 5, Synonyms: []string{"between"}},
	IsNull:             {String: "IS NULL", Alias: "isnull", Position: 6, Synonyms: []string{"is null", "isnull"}},
	NotIn:              {String: "NOT IN", Alias: "nin", Position: 7, Synonyms: []string{"not in", "nin"}},
	Like:               {String: "LIKE", Alias: "like", Position: 8, Synonyms: []string{"like"}},
	In:                 {String: "IN", Alias: "in", Position: 9, Synonyms: []string{"in"}},
	NotEqual:           {String: "!=", Alias: "neq", Position: 10, Synonyms: []string{"!=", "<>", "neq"}},
	GreaterThanOrEqual: {String: ">=", Alias: "gte", Position: 11, Synonyms: []string{">=", "gte"}},
	LessThanOrEqual:    {String: "<=", Alias: "lte", Position: 12, Synonyms: []string{"<=", "lte"}},
	GreaterThan:        {String: ">", Alias: "gt", Position: 13, Synonyms: []string{">", "gt"}},
	LessThan:           {String: "<", Alias: "lt", Position: 14, Synonyms: []string{"<", "lt"}},
	Equal:              {String: "=", Alias: "eq", Position: 15, Synonyms: []string{"=", "eq"}},
}

// reverse index for ParseFrom; built once at init
var parseIndex map[string]Type

func init() {
	parseIndex = make(map[string]Type, len(registry)*3)
	for t, m := range registry {
		addParseToken(t, m.String)
		if m.Alias != "" {
			addParseToken(t, m.Alias)
		}
		for _, s := range m.Synonyms {
			addParseToken(t, s)
		}
	}
}

func addParseToken(t Type, raw string) {
	key := normalize(raw)
	parseIndex[key] = t
}

func normalize(s string) string {
	// case-insensitive; collapse internal whitespace to single spaces
	key := strings.TrimSpace(strings.ToLower(s))
	if strings.IndexByte(key, ' ') >= 0 || strings.IndexByte(key, '\t') >= 0 {
		key = strings.Join(strings.Fields(key), " ")
	}
	return key
}

// cachedKnownOperators holds canonical spellings ordered deterministically by Position.
var cachedKnownOperators = func() []string {
	type kv struct {
		str string
		pos int
	}
	tmp := make([]kv, 0, len(registry))
	for _, m := range registry {
		tmp = append(tmp, kv{str: m.String, pos: m.Position})
	}
	sort.Slice(tmp, func(i, j int) bool { return tmp[i].pos < tmp[j].pos })
	out := make([]string, len(tmp))
	for i, it := range tmp {
		out[i] = it.str
	}
	return out
}()

// Alias returns a short stable mnemonic (e.g., "lte", "nin").
// Invalid → "invalid".
func (t Type) Alias() string {
	if m, ok := registry[t]; ok {
		if m.Alias != "" {
			return m.Alias
		}
	}
	return "invalid"
}

// GetKnownOperators returns all supported operators in canonical spelling,
// ordered deterministically (see cachedKnownOperators above). The returned
// slice is a copy; callers can reorder safely.
func GetKnownOperators() []string {
	out := make([]string, len(cachedKnownOperators))
	copy(out, cachedKnownOperators)
	return out
}

// IsValid reports whether t is a recognized operator.
func (t Type) IsValid() bool {
	_, ok := registry[t]
	return ok && t != Invalid
}

// ParseFrom converts a value into a Type. It accepts:
//   - Type         : returned as-is
//   - string/[]byte: case-insensitive; whitespace-normalized; supports
//     canonical symbols/words and aliases.
//
// Unrecognized inputs yield Invalid.
func ParseFrom(v any) Type {
	switch x := v.(type) {
	case Type:
		return x
	case string:
		if t, ok := parseIndex[normalize(x)]; ok {
			return t
		}
	case []byte:
		if t, ok := parseIndex[normalize(string(x))]; ok {
			return t
		}
	}
	return Invalid
}

// String returns the canonical SQL spelling (render form). Invalid → "Invalid".
func (t Type) String() string {
	if m, ok := registry[t]; ok {
		return m.String
	}
	return "Invalid"
}
