package tosca_v1_3

import "testing"

func TestExternalSchemaCompileCacheIsOperationLocal(t *testing.T) {
	calls := 0
	validator := newExternalSchemaValidator()
	validator.compile = func(schemaType string, schema string) string {
		calls++
		return externalSchemaCompilationFailure
	}

	for range 2 {
		if got := validator.compileCached("json", `{"type":7}`); got != externalSchemaCompilationFailure {
			t.Fatalf("cached category = %q", got)
		}
	}
	if calls != 1 {
		t.Fatalf("identical schema compiled %d times, want 1", calls)
	}

	independent := newExternalSchemaValidator()
	independent.compile = validator.compile
	independent.compileCached("json", `{"type":7}`)
	if calls != 2 {
		t.Fatalf("cache crossed processing contexts: compile count = %d, want 2", calls)
	}
}
