import {describe, it, expect} from 'vitest';
import {grupeazaPeCereri} from './grupareCereri';
import type {datebaza} from '../wailsjs/go/models';

const zi = (angajat_id: string, data: string, alte: Partial<datebaza.ZiAbsenta> = {}) =>
    ({angajat_id, data, cod: 'CO', activ: true, stare_aprobare: 'APROBATA', ...alte}) as datebaza.ZiAbsenta;

// F05-T21 (I050): vederea pe cereri grupeaza zilele aceleiasi cereri.
describe('gruparea absentelor pe cereri', () => {
    it('pune la un loc zilele cu acelasi document sau numar, la acelasi angajat', () => {
        const g = grupeazaPeCereri([
            zi('a', '2026-10-07', {nr_inregistrare: '12'}), zi('a', '2026-10-05', {nr_inregistrare: '12'}),
            zi('b', '2026-10-05', {nr_inregistrare: '12'}),
            zi('c', '2026-10-01', {document_id: 'doc1'}), zi('c', '2026-10-20', {document_id: 'doc1'}),
        ]);
        expect(g.map(x => [x.zile[0].angajat_id, x.zile.map(z => z.data)])).toEqual([
            ['c', ['2026-10-01', '2026-10-20']], ['a', ['2026-10-05', '2026-10-07']], ['b', ['2026-10-05']],
        ]);
        expect(g.every(x => !x.faraCerere)).toBe(true);
    });
    it('fara cerere: un interval peste weekend, dar separat la pauza lunga, alt cod sau alta stare', () => {
        const g = grupeazaPeCereri([
            zi('a', '2026-10-02'), zi('a', '2026-10-05'), zi('a', '2026-10-06'),
            zi('a', '2026-10-15'),
            zi('a', '2026-10-07', {cod: 'CM'}),
            zi('a', '2026-10-08', {stare_aprobare: 'IN_ASTEPTARE'}),
        ]);
        expect(g.filter(x => x.faraCerere).map(x => x.zile.map(z => z.data))).toEqual([
            ['2026-10-02', '2026-10-05', '2026-10-06'], ['2026-10-07'], ['2026-10-08'], ['2026-10-15'],
        ]);
    });
});
