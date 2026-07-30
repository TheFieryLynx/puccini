package tosca_1_3_corpus_test

import (
	"os"
	"path/filepath"
	"testing"
)

func FuzzTOSCA13Parser(f *testing.F) {
	root := corpusRoot(f)
	seeds, err := filepath.Glob(filepath.Join(root, "fuzz-seeds", "*.yaml"))
	if err != nil {
		f.Fatal(err)
	}
	for _, seed := range seeds {
		data, readErr := os.ReadFile(seed)
		if readErr != nil {
			f.Fatal(readErr)
		}
		f.Add(data)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 1<<20 {
			t.Skip()
		}
		path := filepath.Join(t.TempDir(), "seed.yaml")
		if err := os.WriteFile(path, data, 0o600); err != nil {
			t.Fatal(err)
		}
		_, first, firstErr := parseCorpusFile(t, path)
		_, second, secondErr := parseCorpusFile(t, path)
		if (firstErr == nil) != (secondErr == nil) || first != second {
			t.Fatalf("nondeterministic parse result:\nfirst:\n%s\nsecond:\n%s", first, second)
		}
	})
}
