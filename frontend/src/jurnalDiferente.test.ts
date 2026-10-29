import {describe, expect, it} from 'vitest';
import {diferente, rezumatNou} from './jurnalDiferente';

describe('diferente', () => {
    it('arată doar câmpurile schimbate, fără cele tehnice', () => {
        const d = diferente('{"valoare":"D8","revizie":1,"nume":"Ana"}', '{"valoare":"N12","revizie":2,"nume":"Ana"}');
        expect(d).toEqual([{camp: 'valoare', inainte: 'D8', dupa: 'N12'}]);
    });
    it('compară texte simple când valorile nu sunt obiecte', () => {
        expect(diferente('"a"', '"b"')).toEqual([{camp: '', inainte: 'a', dupa: 'b'}]);
    });
    it('rezumă o creare fără câmpurile goale sau tehnice', () => {
        expect(rezumatNou('{"cod":"CO","creat_la":"x","motiv":""}')).toBe('cod: CO');
    });
});
