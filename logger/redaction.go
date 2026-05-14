package logger

import (
	"fmt"
	"strings"
)

// RedactedPlaceholder is the value substituted in place of redacted field values.
const RedactedPlaceholder = "***REDACTED***"

// redactKeySubstrings — credentials. Fully replaced with RedactedPlaceholder.
var redactKeySubstrings = []string{
	"api_key", "apikey",
	"password", "passwd",
	"secret",
	"token", "bearer",
	"authorization", "auth_header",
	"credential", "credentials",
	"private_key",
	"cookie",
	"otp",
}

// piiKeySubstrings — PII. Partially masked (first/last chars kept) so logs
// remain useful for debugging while hiding most of the value.
var piiKeySubstrings = []string{
	"email",
	"phone", "mobile",
	"identifier",
}

// ShouldRedact reports whether a field key looks like a credential and its
// value must be fully replaced with RedactedPlaceholder.
func ShouldRedact(key string) bool {
	lower := strings.ToLower(key)
	for _, needle := range redactKeySubstrings {
		if strings.Contains(lower, needle) {
			return true
		}
	}
	return false
}

// ShouldMask reports whether a field key looks like PII and its value should
// be partially masked.
func ShouldMask(key string) bool {
	lower := strings.ToLower(key)
	for _, needle := range piiKeySubstrings {
		if strings.Contains(lower, needle) {
			return true
		}
	}
	return false
}

// MaskValue partially masks a string, keeping a few characters at each end and
// replacing the middle with asterisks. The number of visible characters scales
// with input length:
//   - length <  8  → fully redacted (too short to safely reveal any part)
//   - length 8–14  → 3 visible at each end
//   - length ≥ 15  → 4 visible at each end
//
// Example: "alice@example.com" → "alic*********e.com",
//
//	"+84912345678" → "+84******678".
func MaskValue(s string) string {
	n := len(s)
	if n < 8 {
		return RedactedPlaceholder
	}
	keep := 3
	if n >= 15 {
		keep = 4
	}
	return s[:keep] + strings.Repeat("*", n-2*keep) + s[n-keep:]
}

// RedactFields returns a deep-copied Fields map with sensitive values replaced
// or masked. Nested maps and slices are walked recursively so secrets buried
// inside structured payloads are still scrubbed. The input map is not mutated.
func RedactFields(in Fields) Fields {
	if in == nil {
		return nil
	}
	out := make(Fields, len(in))
	for k, v := range in {
		if ShouldRedact(k) {
			out[k] = RedactedPlaceholder
			continue
		}
		if ShouldMask(k) {
			out[k] = maskValue(v)
			continue
		}
		out[k] = redactValue(v)
	}
	return out
}

func redactValue(v interface{}) interface{} {
	switch val := v.(type) {
	case Fields:
		return RedactFields(val)
	case map[string]interface{}:
		return RedactFields(Fields(val))
	case []interface{}:
		cp := make([]interface{}, len(val))
		for i, item := range val {
			cp[i] = redactValue(item)
		}
		return cp
	default:
		return v
	}
}

func maskValue(v interface{}) interface{} {
	switch val := v.(type) {
	case string:
		return MaskValue(val)
	case []string:
		cp := make([]string, len(val))
		for i, s := range val {
			cp[i] = MaskValue(s)
		}
		return cp
	case nil:
		return nil
	default:
		return MaskValue(fmt.Sprintf("%v", val))
	}
}
