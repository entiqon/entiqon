// Package errors provides extended error handling utilities for Entiqon.
//
// It introduces the CausableError interface for errors that carry both
// a machine-readable cause and a human-readable reason, making it easier
// to distinguish between error categories programmatically while keeping
// detailed descriptions for users and logs.
//
// It also includes ProcessStageError, which associates errors with specific
// processing stages, supporting structured error reporting and diagnostics.
package errors
