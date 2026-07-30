#!/bin/sh
set -eu

python3 scripts/conformance/generate_tosca13_corpus_docs.py
git diff --exit-code -- \
  docs/conformance/tosca-1.3/template-corpus-sections.yaml \
  docs/conformance/tosca-1.3/template-corpus-pairwise.yaml \
  docs/conformance/tosca-1.3/template-corpus-summary.yaml
go test -count=1 ./tests/corpus/tosca_1_3
go test -count=1 ./tests/conformance/tosca_2_0
