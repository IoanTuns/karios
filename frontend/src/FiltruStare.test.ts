import {describe,it,expect} from 'vitest';
import {corespundeStare} from './FiltruStare';
describe('filtrarea starii',()=>{
 it('separa activele de inactive si permite istoricul complet',()=>{
  const lista=[{id:'a',activ:true},{id:'i',activ:false}];
  expect(lista.filter(v=>corespundeStare(v.activ,'ACTIVE')).map(v=>v.id)).toEqual(['a']);
  expect(lista.filter(v=>corespundeStare(v.activ,'INACTIVE')).map(v=>v.id)).toEqual(['i']);
  expect(lista.filter(v=>corespundeStare(v.activ,'TOATE')).map(v=>v.id)).toEqual(['a','i']);
  expect(lista).toHaveLength(2);
 });
});
