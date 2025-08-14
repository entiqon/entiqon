package token

// FieldCollection is a mutable, chainable sequence of Field items.
//
// A nil receiver is safe: all methods are nil-safe and will no-op or return
// sensible defaults. Mutating methods (except Clear) return *FieldCollection
// to allow fluent chaining. Clear intentionally ends the chain by returning no value.
type FieldCollection []Field

// Add appends one or more fields to the end of the collection.
//
// The method is chainable. If f is nil or fields is empty, Add is a no-op.
func (f *FieldCollection) Add(fields ...Field) *FieldCollection {
	if f == nil || len(fields) == 0 {
		return f
	}
	*f = append(*f, fields...)
	return f
}

// Clear removes all elements from the collection.
//
// Clear intentionally does not return *FieldCollection, so it ends any fluent
// chain. The receiver is nil-safe; calling Clear on a nil receiver is a no-op.
// Capacity may be retained due to re-slicing to zero length.
func (f *FieldCollection) Clear() {
	if f == nil || len(*f) == 0 {
		return
	}
	// Re-slice to zero length (capacity retained).
	*f = (*f)[:0]
}

// Clone returns a shallow copy of the collection.
//
// The returned slice is independent of the original backing array.
// If the receiver is nil or empty, Clone returns nil.
func (f *FieldCollection) Clone() FieldCollection {
	if f == nil || len(*f) == 0 {
		return nil
	}
	cp := make(FieldCollection, len(*f))
	copy(cp, *f)
	return cp
}

// Contains reports whether field exists in the collection using == comparison.
//
// The method is nil-safe; a nil receiver returns false.
func (f *FieldCollection) Contains(field Field) bool {
	if f == nil {
		return false
	}
	for _, fld := range *f {
		if fld == field {
			return true
		}
	}
	return false
}

// IndexOf returns the zero-based index of the first occurrence of field
// using == comparison, or -1 if the field is not present.
//
// The method is nil-safe; a nil receiver returns -1.
func (f *FieldCollection) IndexOf(field Field) int {
	if f == nil {
		return -1
	}
	for i, fld := range *f {
		if fld == field {
			return i
		}
	}
	return -1
}

// InsertAfter inserts one or more fields immediately after the first occurrence
// of target within the collection.
//
// If target is not found, the fields are appended to the end.
// The method is chainable and nil-safe. If fields is empty, the call is a no-op.
func (f *FieldCollection) InsertAfter(target Field, fields ...Field) *FieldCollection {
	if f == nil || len(fields) == 0 {
		return f
	}
	idx := f.IndexOf(target)
	if idx == -1 {
		// Target not found → append.
		return f.Add(fields...)
	}
	return f.insertManyAt(idx+1, fields)
}

// InsertAt inserts a single field at the specified index.
//
// If index is less than 0, the element is inserted at the beginning.
// If index is greater than the current length, the element is appended.
// The method is chainable and nil-safe.
func (f *FieldCollection) InsertAt(index int, field Field) *FieldCollection {
	if f == nil {
		return f
	}
	return f.insertManyAt(index, []Field{field})
}

// InsertBefore inserts one or more fields immediately before the first
// occurrence of target within the collection.
//
// If target is not found, the fields are prepended at index 0.
// The method is chainable and nil-safe. If fields is empty, the call is a no-op.
func (f *FieldCollection) InsertBefore(target Field, fields ...Field) *FieldCollection {
	if f == nil || len(fields) == 0 {
		return f
	}
	idx := f.IndexOf(target)
	if idx == -1 {
		// Target not found → prepend.
		return f.insertManyAt(0, fields)
	}
	return f.insertManyAt(idx, fields)
}

// IsEmpty reports whether the collection contains no elements.
//
// The method is nil-safe; a nil receiver is considered empty.
func (f *FieldCollection) IsEmpty() bool { return f.Length() == 0 }

// Length reports the number of elements in the collection.
//
// The method is nil-safe; a nil receiver returns 0.
func (f *FieldCollection) Length() int {
	if f == nil {
		return 0
	}
	return len(*f)
}

// Remove deletes the first occurrence of field from the collection using == comparison.
//
// If the field does not exist, the collection is unchanged.
// The method is chainable and nil-safe.
func (f *FieldCollection) Remove(field Field) *FieldCollection {
	if f == nil || len(*f) == 0 {
		return f
	}
	for i, fld := range *f {
		if fld == field {
			*f = append((*f)[:i], (*f)[i+1:]...)
			break
		}
	}
	return f
}

// RemoveAt deletes the element at index if it is within bounds.
//
// If index is out of range, the collection is unchanged.
// The method is chainable and nil-safe.
func (f *FieldCollection) RemoveAt(index int) *FieldCollection {
	if f == nil || index < 0 || index >= len(*f) {
		return f
	}
	*f = append((*f)[:index], (*f)[index+1:]...)
	return f
}

// insertManyAt inserts fields at index, clamping index to [0, len].
//
// This helper is nil-safe and chainable. If fields is empty, the call is a no-op.
// It is unexported to keep the public API focused on Add/InsertAt/InsertAfter/InsertBefore.
func (f *FieldCollection) insertManyAt(index int, fields []Field) *FieldCollection {
	if index < 0 {
		index = 0
	}
	if index > len(*f) {
		index = len(*f)
	}
	*f = append((*f)[:index], append(fields, (*f)[index:]...)...)
	return f
}
