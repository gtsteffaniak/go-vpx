from pathlib import Path
import struct,hashlib,json

def text_section(path):
 b=Path(path).read_bytes()
 assert struct.unpack_from('<I',b)[0]==0xfeedfacf
 n=struct.unpack_from('<I',b,16)[0];at=32
 for _ in range(n):
  cmd,size=struct.unpack_from('<II',b,at)
  if cmd==0x19:
   count=struct.unpack_from('<I',b,at+64)[0]
   for i in range(count):
    s=at+72+80*i
    if b[s:s+16].split(b'\0')[0]==b'__text':
     length=struct.unpack_from('<Q',b,s+40)[0];offset=struct.unpack_from('<I',b,s+48)[0]
     return b[offset:offset+length]
  at+=size
 raise ValueError('text section missing')
import subprocess
files=['/private/tmp/govpx-reductions-scoped-amd64.test','/private/tmp/govpx-fused-final-amd64.test']
sections=[text_section(f) for f in files];offsets=[]
for f in files:
 output=subprocess.check_output(['/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/llvm-nm','--numeric-sort','--defined-only',f],text=True)
 lines=output.splitlines();base=int(lines[0].split()[0],16)
 for line in lines:
  a=line.split()
  if a[-1] not in ['_go:buildid','_runtime.text']:
   offsets.append(int(a[0],16)-base);break
r={'files':files,'header_bytes':offsets,'code_bytes':[len(x)-o for x,o in zip(sections,offsets)],'code_sha256':[hashlib.sha256(x[o:]).hexdigest() for x,o in zip(sections,offsets)],'code_identical':sections[0][offsets[0]:]==sections[1][offsets[1]:]}
print(json.dumps(r,indent=2));assert r['code_identical']
