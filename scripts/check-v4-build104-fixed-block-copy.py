#!/usr/bin/env python3
"""Verify the fixed8 specialization and otherwise unchanged stage dispatch."""
from pathlib import Path
import re
import subprocess
root=Path(__file__).resolve().parent.parent
source=(root/'watermark/experimental_v4_phone_build84_rgb_fetch.go').read_text()
a=source.index('func experimentalV4PhoneBuild84ReadProjectiveBlockValue(')
b=source.index('\ntype experimentalV4PhoneBuild84Gen4Result',a)
expected=source[a:b]
a=source.index('func experimentalV4PhoneBuild84Generation4Fetch(')
b=source.index('\nfunc ',a+5)
expected+='\n'+source[a:b]
for name in ('ReadProjectiveBlockValue','FoldScoreFetch','ContinueFetch','Generation4Fetch'):
 expected=expected.replace('experimentalV4PhoneBuild84'+name,'experimentalV4Build104'+name)
a=expected.index('\nfunc experimentalV4Build104FoldScoreFetch')
reader=expected[:a].replace('originX, originY, size int','originX, originY int').replace('readCosTables[size]','readCosTables[blockSize]').replace('table2, table3 := table[2], table[3]','table2, table3 := (*[blockSize]float64)(table[2]), (*[blockSize]float64)(table[3])').replace('y < size','y < blockSize').replace('x < size','x < blockSize')
expected=reader+expected[a:]
expected=expected.replace('experimentalV4Build104ReadProjectiveBlockValue(plane, h, (baseX+px)*blockSize, (baseY+py)*blockSize, blockSize)','experimentalV4Build104ReadProjectiveBlockValue(plane, h, (baseX+px)*blockSize, (baseY+py)*blockSize)')
actual=(root/'watermark/experimental_v4_build104_fixed_block_candidate_test.go').read_text()
actual=actual[actual.index('func experimentalV4Build104ReadProjectiveBlockValue'):]
def normalize(text):
 text=re.sub(r'//[^\n]*','',text)
 text=subprocess.run(['gofmt'],input='package watermark\n'+text,text=True,capture_output=True,check=True).stdout
 return re.sub(r'\s+','',text)
if normalize(expected)!=normalize(actual): raise SystemExit('FAIL: Build104 differs beyond fixed8 contract and named dispatch')
print('Build104 fixed8 reader/stage structural check: PASS')
