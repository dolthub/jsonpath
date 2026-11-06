// Copyright 2026 DoltHub, Inc.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file or at
// https://opensource.org/licenses/MIT.

package jsonpath

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyArrayWildcard(t *testing.T) {
	for _, path := range []string{"$[*]", "$[*].email", "$[*].email.name", "$[*][0]"} {
		t.Run(path, func(t *testing.T) {
			for _, obj := range []interface{}{[]interface{}{}, []interface{}(nil), []string{}} {
				value, err := JsonPathLookup(obj, path)
				assert.Nil(t, value)
				assert.True(t, ErrKeyError.Is(err), "expected no match, got %v", err)
			}
		})
	}
	for _, path := range []string{"$.users[*]", "$.users[*].email"} {
		value, err := JsonPathLookup(map[string]interface{}{"users": []interface{}{}}, path)
		assert.Nil(t, value)
		assert.True(t, ErrKeyError.Is(err), "expected no match, got %v", err)
	}
}

func TestArrayWildcardPreservesValues(t *testing.T) {
	for _, tt := range []struct{ document, path, expected string }{
		{`[]`, `$`, `[]`},
		{`{"users":[]}`, `$.users`, `[]`},
		{`[{"email":"jordan@example.com"}]`, `$[*].email`, `["jordan@example.com"]`},
		{`[{"email":null}]`, `$[*].email`, `[null]`},
		{`[{"email":[]}]`, `$[*].email`, `[[]]`},
		{`[[],[1]]`, `$[*]`, `[[],[1]]`},
	} {
		t.Run(tt.document+tt.path, func(t *testing.T) {
			var document, expected interface{}
			assert.NoError(t, json.Unmarshal([]byte(tt.document), &document))
			assert.NoError(t, json.Unmarshal([]byte(tt.expected), &expected))
			value, err := JsonPathLookup(document, tt.path)
			assert.NoError(t, err)
			assert.Equal(t, expected, value)
		})
	}
}

func TestEmptyArrayWildcardDoesNotHideInvalidPaths(t *testing.T) {
	for _, path := range []string{"$[*].email[", "$[*][invalid]"} {
		_, err := JsonPathLookup([]interface{}{}, path)
		assert.Error(t, err)
		assert.False(t, ErrKeyError.Is(err))
	}
	// Explicit ranges retain their bounds errors, even when both ends are omitted.
	for _, path := range []string{"$[:]", "$[0:1]", "$[:1]", "$[1:]"} {
		_, err := JsonPathLookup([]interface{}{}, path)
		assert.Error(t, err)
		assert.False(t, ErrKeyError.Is(err))
	}
}
