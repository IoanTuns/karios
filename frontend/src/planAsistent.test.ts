import {describe, expect, it} from 'vitest';

describe('logica filtrare plan AI (F14-T32)', () => {
    const esteExceptie = (stare: string) => stare === 'ocupata' || stare === 'blocata' || stare === 'identica';

    it('identifică corect stările de excepție din plan', () => {
        expect(esteExceptie('ocupata')).toBe(true);
        expect(esteExceptie('blocata')).toBe(true);
        expect(esteExceptie('identica')).toBe(true);
        expect(esteExceptie('noua')).toBe(false);
        expect(esteExceptie('inlocuire')).toBe(false);
    });

    it('calculează numărul de persoane unice afectate cu modificări', () => {
        const randuri = [
            {angajat_id: 'a1', stare: 'noua'},
            {angajat_id: 'a1', stare: 'inlocuire'},
            {angajat_id: 'a2', stare: 'ocupata'},
            {angajat_id: 'a3', stare: 'blocata'},
            {angajat_id: 'a4', stare: 'inlocuire'},
        ];
        const persoaneAfectate = new Set(
            randuri
                .filter(r => r.stare === 'noua' || r.stare === 'inlocuire')
                .map(r => r.angajat_id)
        ).size;
        expect(persoaneAfectate).toBe(2); // doar a1 si a4
    });
});

