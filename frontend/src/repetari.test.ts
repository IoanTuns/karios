import {describe, expect, it} from 'vitest';
import {ai} from '../wailsjs/go/models';
import {cheiePropunere, imparteRepetari} from './repetari';

const p = (id: string, cod: string) => new ai.PropunereActiune({id, tip: 'planifica_tura', descriere: '', parametri: {marca: 'M1', data: '2026-10-15', cod_schimb: cod}, necesita_confirmare: true});

describe('imparteRepetari', () => {
    it('cheia nu depinde de id sau de ordinea parametrilor', () => {
        const a = p('a', 'Z7');
        const b = new ai.PropunereActiune({...a, id: 'b', parametri: {cod_schimb: 'Z7', data: '2026-10-15', marca: 'M1'}});
        expect(cheiePropunere(a)).toBe(cheiePropunere(b));
    });
    it('omite duplicatele unei propuneri neconfirmate', () => {
        const r = imparteRepetari([p('nou', 'Z7'), p('alt', 'N12')], [{propuneri: [p('vechi', 'Z7')]}]);
        expect(r.noi.map(x => x.id)).toEqual(['alt']);
        expect(r.omise).toBe(1);
    });
    it('arată, marcată, o propunere identică cu una deja aplicată (F14-T19)', () => {
        const r = imparteRepetari([p('nou', 'Z7')], [{propuneri: [p('vechi', 'Z7')], propuneriAplicata: {vechi: true}}]);
        expect(r.noi.map(x => x.id)).toEqual(['nou']);
        expect(r.reaplicari).toEqual({nou: true});
        expect(r.omise).toBe(0);
    });
});
