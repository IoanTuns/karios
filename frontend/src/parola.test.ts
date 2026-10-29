import {describe, expect, it} from 'vitest';
import {parolaValida} from './parola';

describe('parolaValida', () => {
    it('acceptă 8 caractere cu literă și cifră', () => {
        expect(parolaValida('parola12')).toBe(true);
        expect(parolaValida('Șapte7ăă')).toBe(true);
    });
    it('respinge parolele scurte sau fără literă ori cifră', () => {
        for (const p of ['scurt1', 'faracifre', '12345678', '']) expect(parolaValida(p)).toBe(false);
    });
});
