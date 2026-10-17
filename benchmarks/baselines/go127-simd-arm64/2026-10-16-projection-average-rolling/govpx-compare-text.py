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
a=text_section('/private/tmp/govpx-reductions-final-bench');b=text_section('/private/tmp/govpx-reductions-scoped-bench')
# llvm-nm locates go:buildid at __text+0 and the first instruction-bearing
# function internal/abi.BoundsDecode at __text+0x70 in both binaries.
assert a.startswith(b'\xff Go build ID: ') and b.startswith(b'\xff Go build ID: ')
first_code_offset=0x70
r={'before_sha256':hashlib.sha256(a).hexdigest(),'after_sha256':hashlib.sha256(b).hexdigest(),'before_bytes':len(a),'after_bytes':len(b),'ignored_header_bytes':first_code_offset,'before_code_sha256':hashlib.sha256(a[first_code_offset:]).hexdigest(),'after_code_sha256':hashlib.sha256(b[first_code_offset:]).hexdigest(),'code_identical':a[first_code_offset:]==b[first_code_offset:]}
Path('/private/tmp/govpx-reductions-scoped-arm64-text.json').write_text(json.dumps(r,indent=2)+'\n')
print(json.dumps(r));assert r['code_identical'],'ARM executable code changed'
