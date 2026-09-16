package tosca_1_3_corpus_test

import (
	"fmt"
	"testing"
)

// Normalized TOSCA maps store entries in an unordered list with $key fields.
// Index those entries by their string key for assertions; never depend on Go
// map iteration order or pretend that the entry at index zero is token_type.
func evidenceMapView(value any) (any, error) {
	switch v := value.(type) {
	case map[string]any:
		result := map[string]any{}
		for k, x := range v {
			if k == "$map" {
				entries, ok := x.([]any)
				if !ok {
					return nil, fmt.Errorf("malformed normalized map")
				}
				indexed := map[string]any{}
				for _, raw := range entries {
					entry, ok := raw.(map[string]any)
					if !ok {
						return nil, fmt.Errorf("malformed map entry")
					}
					key, ok := entry["$key"].(map[string]any)
					if !ok {
						return nil, fmt.Errorf("missing normalized map key")
					}
					name, ok := key["$primitive"].(string)
					if !ok { // Non-string keys are not addressable in this evidence view.
						result[k] = x
						indexed = nil
						break
					}
					if _, found := indexed[name]; found {
						return nil, fmt.Errorf("duplicate normalized map key %s", name)
					}
					item, err := evidenceMapView(entry)
					if err != nil {
						return nil, err
					}
					indexed[name] = item
				}
				if indexed != nil {
					result[k] = indexed
				}
			} else {
				item, err := evidenceMapView(x)
				if err != nil {
					return nil, err
				}
				result[k] = item
			}
		}
		return result, nil
	case []any:
		result := make([]any, len(v))
		for i, x := range v {
			item, err := evidenceMapView(x)
			if err != nil {
				return nil, err
			}
			result[i] = item
		}
		return result, nil
	default:
		return value, nil
	}
}
func TestEvidenceMapKeys(t *testing.T) {
	entry := func(k, v string) any { return map[string]any{"$key": map[string]any{"$primitive": k}, "$primitive": v} }
	for _, es := range [][]any{{entry("token", "secret"), entry("token_type", "password")}, {entry("token_type", "password"), entry("token", "secret")}} {
		v, err := evidenceMapView(map[string]any{"$map": es})
		if err != nil {
			t.Fatal(err)
		}
		got, err := jsonPointer(v, "/$map/token_type/$primitive")
		if err != nil || got != "password" {
			t.Fatalf("%v %v", got, err)
		}
		if _, err := jsonPointer(v, "/$map/absent/$primitive"); err == nil {
			t.Fatal("missing key accepted")
		}
	}
	if _, err := evidenceMapView(map[string]any{"$map": []any{entry("x", "a"), entry("x", "b")}}); err == nil {
		t.Fatal("duplicate key accepted")
	}
}
