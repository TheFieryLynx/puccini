#!/usr/bin/env python3
"""Recompute evidence for historical IDs without changing catalog membership."""
from pathlib import Path
import argparse
import collections
import hashlib
import json
import os
import re
import subprocess
import yaml

ROOT=Path(__file__).resolve().parents[2]
BASE=ROOT/'docs/conformance/tosca-1.3'
CORPUS=ROOT/'tests/corpus/tosca_1_3'
EXECUTION=BASE/'re-audit/evidence-execution.yaml'
REPORT=BASE/'re-audit/evidence-recheck.yaml'
def load(path):return yaml.load(path.read_text(),Loader=yaml.CSafeLoader)
def sha(path):return hashlib.sha256(path.read_bytes()).hexdigest()
def encoded(data):return yaml.safe_dump(data,sort_keys=False,width=120)
def corpus_cases():return load(CORPUS/'manifest.yaml')['cases']+load(CORPUS/'non_must/manifest.yaml')['cases']
def fingerprint():
    # Include manifests, reviewed assertions and transitive imported fixtures.
    paths={p for p in CORPUS.rglob('*') if p.is_file() and p.suffix in {'.go','.yaml','.yml','.json','.zip','.csar','.js'}}
    for folder in ['tosca','normal','clout','assets/tosca']:
        paths.update(p for p in (ROOT/folder).rglob('*') if p.is_file() and p.suffix in {'.go','.yaml','.js'})
    return hashlib.sha256(''.join(str(p.relative_to(ROOT))+':'+sha(p)+'\n' for p in sorted(paths)).encode()).hexdigest()
def record_execution():
    command=['go','test','-json','-count=1','./tests/corpus/tosca_1_3','-run','^(TestTOSCA13Corpus|TestTOSCA13NonMUSTCorpus)$']
    env={k:v for k,v in os.environ.items() if not k.startswith(('TOSCA_CORPUS_','TOSCA_NON_MUST_','TOSCA_EVIDENCE_'))}
    result=subprocess.run(command,cwd=ROOT,env=env,capture_output=True,text=True,timeout=120)
    if result.returncode:raise SystemExit(result.stdout+result.stderr)
    events=[json.loads(line) for line in result.stdout.splitlines() if line.startswith('{')]
    passed={e.get('Test') for e in events if e.get('Action')=='pass'}
    rows=[]
    for c in corpus_cases():
        prefix='TestTOSCA13NonMUSTCorpus/' if c['file'].startswith('non_must/') else 'TestTOSCA13Corpus/'
        name=prefix+c['id']; output=''.join(e.get('Output','') for e in events if e.get('Test')==name)
        match=re.search(r'EVIDENCE phase=(\w+) assertions=(\d+)',output)
        if name not in passed:raise SystemExit('case not executed successfully: '+name)
        if not match:raise SystemExit('missing observed phase: '+name)
        if int(match[2])!=len(c['assertions']):raise SystemExit('owned assertion count mismatch')
        rows.append({'id':c['id'],'status':'passed','phase':match[1],'executed_assertions':[a['id'] for a in c['assertions']],'fixture_sha256':sha(CORPUS/c['file'])})
    EXECUTION.write_text(encoded({'schema_version':1,'command':command,'source_fingerprint':fingerprint(),'cases':rows}))

def reasons(req,review,cases,mutations,execution):
    r=review.get(req)
    if not r:return ['no-reviewed-primary-predicate','historical-labels-or-unreviewed-direct-tests-only']
    linked=[c for c in cases if req in c['coverage']['primary_requirements']]
    result=[]
    if not r.get('complete_predicate'):result.append('only-subset-of-historical-predicate-reviewed')
    if not linked:result.append('no-primary-fixture')
    pos=neg=False
    for c in linked:
        owned=[a for a in c['assertions'] if req in a['requirement_ids']]
        if not owned:result.append('primary-without-owned-assertion')
        e=execution.get(c['id'],{})
        if e.get('status')!='passed' or not {a['id'] for a in owned}<=set(e.get('executed_assertions',[])):result.append('not-executed-owned-assertions')
        if c['expected']['accepted']:
            pos=True
            if r['evidence_kind'].startswith('semantic-') or r['evidence_kind']=='normalization':
                if not any(a['kind'] in {'value','normalization'} for a in owned):result.append('semantic-positive-lacks-value')
        else:
            neg=True
            if e.get('phase')!=c['expected']['phase'] or not any(a['kind']=='phase' for a in owned):result.append('phase-not-proven')
            if not all(any(a['kind']=='diagnostic' and a['path']==key for a in owned) for key in ['category','entity']):result.append('targeted-diagnostic-not-proven')
            nearest=c.get('related_cases',{}).get('nearest_valid')
            sibling=next((x for x in cases if x['id']==nearest),None)
            if not sibling or not sibling['expected']['accepted'] or execution.get(nearest,{}).get('status')!='passed':result.append('nearest-valid-not-proven')
    if r.get('positive_required') and not pos:result.append('missing-positive')
    if r.get('negative_required') and not neg:result.append('missing-negative')
    if not r.get('mutation_paths'):result.append('no-reviewed-mutation-path')
    for path in r.get('mutation_paths',[]):
        m=mutations.get(path,{})
        if m.get('status')!='caught' or req not in m.get('requirements',[]):result.append('relevant-owned-mutation-not-caught')
    return sorted(set(result))

def generate():
    cases=corpus_cases();review={r['id']:r for r in load(CORPUS/'evidence-predicates.yaml')['requirements']}
    mutations={p['id']:p for p in load(BASE/'mutation-coverage.yaml')['production_paths']}
    run=load(EXECUTION);executed={e['id']:e for e in run['cases']}
    if run['source_fingerprint']!=fingerprint():raise SystemExit('stale execution fingerprint; run --record-run')
    for c in cases:
        if executed.get(c['id'],{}).get('fixture_sha256')!=sha(CORPUS/c['file']):raise SystemExit('stale fixture execution: '+c['id'])
    frozen=[r for r in load(BASE/'coverage.yaml')['requirements'] if r.get('included_in_must_denominator')]
    assert len(frozen)==223
    rows=[]
    for record in frozen:
        req=record['requirement_id'];why=reasons(req,review,cases,mutations,executed)
        rows.append({'requirement_id':req,'section':record['section'],'verification_status':'insufficient-evidence' if why else 'verified','reasons':why,'primary_cases':[c['id'] for c in cases if req in c['coverage']['primary_requirements']],'supporting_case_count':sum(req in c['coverage']['supporting_requirements'] for c in cases)})
    statuses=collections.Counter(r['verification_status'] for r in rows)
    return {'schema_version':1,'scope':'Historical 223 records only. Verification is evidence status, not proof of catalog completeness, applicability or atomicity. Unreviewed direct Go tests are not silently promoted.','historical_catalog_sha256':sha(BASE/'requirements.yaml'),'historical_coverage_sha256':sha(BASE/'coverage.yaml'),'source_fingerprint':run['source_fingerprint'],'summary':{'historical_records':223,'verified':statuses['verified'],'insufficient_evidence':statuses['insufficient-evidence'],'insufficiency_reasons':dict(collections.Counter(why for r in rows for why in r['reasons'])),'executable_cases':len(cases),'phase_checked_negative_cases':sum(not c['expected']['accepted'] for c in cases),'primary_links':sum(len(c['coverage']['primary_requirements']) for c in cases),'supporting_links':sum(len(c['coverage']['supporting_requirements']) for c in cases),'reviewed_primary_predicates':len(review),'owned_assertions':sum(len(c['assertions']) for c in cases)},'requirements':rows}
def main():
    p=argparse.ArgumentParser();p.add_argument('--record-run',action='store_true');p.add_argument('--check',action='store_true');args=p.parse_args()
    if args.record_run:record_execution()
    text=encoded(generate())
    if args.check:
        if not REPORT.exists() or REPORT.read_text()!=text:raise SystemExit('stale evidence-recheck.yaml')
    else:REPORT.write_text(text)
    print(load(REPORT)['summary'])
if __name__=='__main__':main()
