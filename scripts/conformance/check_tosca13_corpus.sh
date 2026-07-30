#!/bin/sh
set -eu

GOCACHE=${GOCACHE:-/tmp/puccini-go-cache}
GOMODCACHE=${GOMODCACHE:-/tmp/puccini-go-modcache}
export GOCACHE GOMODCACHE

python3 scripts/conformance/generate_tosca13_corpus_docs.py
python3 scripts/conformance/generate_tosca13_non_must.py
git diff --exit-code -- \
  docs/conformance/tosca-1.3/template-corpus-sections.yaml \
  docs/conformance/tosca-1.3/template-corpus-pairwise.yaml \
  docs/conformance/tosca-1.3/template-corpus-summary.yaml \
  docs/conformance/tosca-1.3/non-must-requirements.yaml \
  docs/conformance/tosca-1.3/non-must-coverage.yaml \
  docs/conformance/tosca-1.3/non-must-summary.md \
  docs/conformance/tosca-1.3/implementation-policies.md \
  docs/conformance/tosca-1.3/full-section-coverage.md \
  tests/corpus/tosca_1_3/non_must/manifest.yaml
python3 scripts/conformance/check_tosca13_non_must.py
python3 scripts/conformance/generate_tosca13_oasis_corpus.py --check
python3 scripts/conformance/audit_tosca13_profile.py --check
python3 scripts/conformance/check_tosca13_profile.py --self-test
go test -count=1 ./tests/corpus/tosca_1_3
go test -count=1 ./tests/conformance/tosca_2_0
