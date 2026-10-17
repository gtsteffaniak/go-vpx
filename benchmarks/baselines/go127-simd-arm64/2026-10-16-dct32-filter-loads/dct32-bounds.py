from pathlib import Path
import re
src=Path('internal/vp9/encoder/transform_quant.go').read_text()
constants={k:int(v) for k,v in re.findall(r'(fdctCospi\d+_64)\s*=\s*(\d+)',src)}
body=src.split('func forwardDCT32(input []int, output []int, round bool) {')[1].split('\n}\n')[0]
maxbound=(0,'')
class Aff:
 def __init__(self,c=None,e=0):self.c=c or [0.]*32;self.e=e
 def __add__(self,o):
  if isinstance(o,int):return Aff(self.c,self.e+abs(o))
  v=Aff([a+b for a,b in zip(self.c,o.c)],self.e+o.e)
  global maxbound
  if v.bound()>maxbound[0]:maxbound=(v.bound(),line)
  return v
 def __neg__(self):return self*-1
 def __sub__(self,o):return self+-o
 def __mul__(self,o):
  global maxbound
  v=Aff([a*o for a in self.c],self.e*abs(o));bound=sum(map(abs,v.c))+v.e
  if bound>maxbound[0]:maxbound=(bound,line)
  return v
 def rnd(self,scale):return Aff([a/scale for a in self.c],self.e/scale+0.5)
 def bound(self):return sum(map(abs,self.c))+self.e

def calc(bound,rd):
 global maxbound,line
 maxbound=(0,'')
 inp=[Aff([bound if j==i else 0 for j in range(32)]) for i in range(32)]
 env=dict(constants,input=inp,output=[Aff() for _ in range(32)],step=[Aff() for _ in range(32)],fdctRoundShift=lambda x:x.rnd(16384))
 inround=False
 for line in body.splitlines():
  line=line.strip()
  if line=='if round {':inround=True;continue
  if inround:
   if line=='}':
    if rd:env['output']=[v.rnd(4) for v in env['output']]
    inround=False
   continue
  if re.match(r'(step|output)\[\d+\] = ',line):exec(line,{},env)
 return max(v.bound() for v in env['output']),maxbound
first,fb=calc(1020,False)
intermediate=(first+2)/4
second,sb=calc(intermediate,False)
rd,rb=calc(intermediate,True)
print('first output bound',first,'intermediate bound',intermediate,'max first product',fb)
print('second output bound',second,'max second product',sb,'with rounding',sb[0]+8192)
print('RD second bound',rd,'max RD product',rb)
