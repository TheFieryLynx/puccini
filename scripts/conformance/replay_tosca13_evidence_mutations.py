#!/usr/bin/env python3
"""Replay reviewed mutation patches, restoring exact bytes even on test failure.

Run only in an isolated task with no concurrent build/test/edit of production.
Compilation failure, timeout, or an unrelated failed test NEVER counts as caught.
"""
from pathlib import Path
import argparse
import hashlib
import json
import re
import subprocess
import tempfile
import yaml

ROOT=Path(__file__).resolve().parents[2]
MAP=ROOT/'docs/conformance/tosca-1.3/mutation-coverage.yaml'
REPORT=ROOT/'docs/conformance/tosca-1.3/re-audit/evidence-mutation-replay.yaml'
CORPUS=ROOT/'tests/corpus/tosca_1_3'

def load(p):return yaml.load(p.read_text(),Loader=yaml.CSafeLoader)
def sha(b):return hashlib.sha256(b).hexdigest()
def dump(p,data):p.write_text(yaml.safe_dump(data,sort_keys=False,width=120))
def run(command):
    result=subprocess.run(command,cwd=ROOT,capture_output=True,text=True,timeout=90)
    events=[]
    for line in result.stdout.splitlines():
        try:events.append(json.loads(line))
        except json.JSONDecodeError:pass
    return result,events

def replay(path,cases):
    patch=path['patch'];source=ROOT/patch['file'];original=source.read_bytes()
    before=patch['before'].encode();after=patch['after'].encode()
    if original.count(before)!=patch.get('expected_occurrences',1):raise RuntimeError(f"{path['id']}: patch match count changed")
    selected=path['detecting_cases']
    command=['go','test','-json','-count=1','-timeout=60s','./tests/corpus/tosca_1_3','-run','^TestTOSCA13Corpus$/('+'|'.join(re.escape(c) for c in selected)+')$']
    baseline,events=run(command)
    if baseline.returncode or any(e.get('Action')=='fail' for e in events):
        raise RuntimeError(f"{path['id']}: baseline failed: {baseline.stdout}{baseline.stderr}")
    started={e.get('Test') for e in events if e.get('Action')=='run'}
    if any('TestTOSCA13Corpus/'+c not in started for c in selected):raise RuntimeError('baseline detector not executed')
    # Keep a recovery copy until restoration succeeds. Do not overwrite user edits.
    temp=Path(tempfile.mkdtemp(prefix='tosca13-mutation-'))
    backup=Path(temp)/source.name;backup.write_bytes(original)
    mutant=original.replace(before,after)
    try:
        source.write_bytes(mutant)
        result,events=run(command)
    finally:
        if source.read_bytes()!=mutant:
            raise RuntimeError(f"concurrent edit detected; recover original from {backup}")
        source.write_bytes(original)
    if source.read_bytes()!=original:raise RuntimeError('restoration failed')
    backup.unlink();temp.rmdir()
    failed={e.get('Test') for e in events if e.get('Action')=='fail'}
    detected=[]
    for assertion in path['detecting_assertions']:
        case,aid=assertion.rsplit('/',1);test='TestTOSCA13Corpus/'+case
        output=''.join(e.get('Output','') for e in events if e.get('Test')==test)
        if test in failed and assertion+':' in output:detected.append(assertion)
    caught=result.returncode==1 and bool(detected) and not any('[build failed]' in e.get('Output','') for e in events)
    path['status']='caught' if caught else 'missed'
    path['source_sha256']=sha(original)
    path['detected_assertions']=detected
    path['fixture_sha256']={c:sha((CORPUS/cases[c]['file']).read_bytes()) for c in selected}
    return {'id':path['id'],'original_mutation_id':path.get('original_mutation_id'),'command':command,'baseline_passed':True,'exit_code':result.returncode,'status':path['status'],'detected_assertions':detected,'restored':True,'source_sha256':sha(original),'log_sha256':sha((result.stdout+result.stderr).encode()),'failure_output':''.join(e.get('Output','') for e in events if e.get('Test') in failed)}

def main():
    parser=argparse.ArgumentParser();parser.add_argument('--only',nargs='*');args=parser.parse_args()
    data=load(MAP);cases={c['id']:c for c in load(CORPUS/'manifest.yaml')['cases']}
    existing=load(REPORT) if REPORT.exists() else {'schema_version':1,'results':[]}
    rows={r['id']:r for r in existing['results']}
    for path in data['production_paths']:
        if args.only and path['id'] not in args.only:continue
        row=replay(path,cases);rows[row['id']]=row
        print(row['id'],row['status'],row['detected_assertions'],flush=True)
        existing['results']=list(rows.values());dump(MAP,data);dump(REPORT,existing)
        if row['status']!='caught':raise SystemExit('mutation survived or failed without an owned assertion')
    existing['summary']={'attempted':len(rows),'caught':sum(r['status']=='caught' for r in rows.values()),'missed':sum(r['status']!='caught' for r in rows.values()),'all_restored':all(r['restored'] for r in rows.values())}
    dump(REPORT,existing)
if __name__=='__main__':main()
