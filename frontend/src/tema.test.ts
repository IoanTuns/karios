import {describe, expect, it} from 'vitest';
import {temaInitiala} from './tema';

describe('temaInitiala', () => {
    it('urmează sistemul cât timp nu e aleasă o temă', () => {
        expect(temaInitiala(null, true)).toBe('dark');
        expect(temaInitiala(null, false)).toBe('light');
    });
    it('păstrează alegerea utilizatorului', () => {
        expect(temaInitiala('light', true)).toBe('light');
        expect(temaInitiala('dark', false)).toBe('dark');
    });
    it('ignoră o valoare salvată necunoscută', () => {
        expect(temaInitiala('albastru', true)).toBe('dark');
    });
});
