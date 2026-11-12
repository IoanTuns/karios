import {afterEach, describe, expect, it} from 'vitest';
import i18next from 'i18next';
import {alegeLimbaRapoarte, CHEIE_LIMBA_INTERFATA, limbaInterfataSalvata, limbaRapoarte, preferintaLimbaRapoarte} from './limba';

// F03-T05 (D103): preferintele de limba, cu o stocare locala simulata.
const stocare = new Map<string, string>();
(globalThis as {localStorage?: unknown}).localStorage = {
    getItem: (k: string) => stocare.get(k) ?? null,
    setItem: (k: string, v: string) => { stocare.set(k, v); },
};

describe('limba', () => {
    afterEach(() => stocare.clear());
    it('foloseste limba salvata doar daca e disponibila', () => {
        expect(limbaInterfataSalvata(['ro', 'en'])).toBeUndefined();
        stocare.set(CHEIE_LIMBA_INTERFATA, 'en');
        expect(limbaInterfataSalvata(['ro', 'en'])).toBe('en');
        stocare.set(CHEIE_LIMBA_INTERFATA, 'hu');
        expect(limbaInterfataSalvata(['ro', 'en'])).toBeUndefined();
    });
    it('rapoartele urmeaza interfata pana la o alegere explicita', async () => {
        await i18next.init({lng: 'en', resources: {}});
        expect(preferintaLimbaRapoarte()).toBe('interfata');
        expect(limbaRapoarte()).toBe('en');
        alegeLimbaRapoarte('ro');
        expect(limbaRapoarte()).toBe('ro');
        alegeLimbaRapoarte('altceva');
        expect(preferintaLimbaRapoarte()).toBe('interfata');
    });
});
