import {describe, expect, it} from 'vitest';
import {dataAfisata, formatData, formateazaDataSigura} from './formatare';

describe('formatData', () => {
    it('refolosește formatterul pentru aceeași limbă și aceleași opțiuni', () => {
        expect(formatData('ro', {weekday: 'short', timeZone: 'UTC'})).toBe(formatData('ro', {weekday: 'short', timeZone: 'UTC'}));
        expect(formatData('ro', {weekday: 'short', timeZone: 'UTC'})).not.toBe(formatData('en', {weekday: 'short', timeZone: 'UTC'}));
        expect(formatData('ro')).not.toBe(formatData('ro', {month: 'long'}));
    });
    it('formatează ca Intl.DateTimeFormat', () => {
        const d = new Date('2026-10-31T12:00:00Z');
        expect(formatData('en', {weekday: 'long', timeZone: 'UTC'}).format(d)).toBe('Saturday');
    });
});

describe('formateazaDataSigura', () => {
    it('returnează fallback la valori vide sau nule', () => {
        expect(formateazaDataSigura('', 'ro')).toBe('—');
        expect(formateazaDataSigura(null, 'ro')).toBe('—');
        expect(formateazaDataSigura(undefined, 'ro')).toBe('—');
        expect(formateazaDataSigura('', 'ro', undefined, 'N/A')).toBe('N/A');
    });

    it('returnează fallback la stringuri de dată invalide fără a arunca RangeError', () => {
        expect(formateazaDataSigura('data-invalida', 'ro')).toBe('—');
        expect(formateazaDataSigura('2026-99-99', 'en')).toBe('—');
        expect(formateazaDataSigura(new Date('invalid'), 'ro')).toBe('—');
    });

    it('formatează corect date valide', () => {
        const d = new Date('2026-10-31T12:00:00Z');
        expect(formateazaDataSigura(d, 'en', {weekday: 'long', timeZone: 'UTC'})).toBe('Saturday');
        expect(formateazaDataSigura('2026-10-31T12:00:00Z', 'en', {weekday: 'long', timeZone: 'UTC'})).toBe('Saturday');
    });
});


describe('ore in format 24h', () => {
    for (const limba of ['ro','en','en-US']) {
        it('foloseste 00-23 fara AM/PM pentru '+limba, () => {
            for (const [ora, asteptat] of [['00','00:05'],['07','07:05'],['19','19:05'],['23','23:05']]) {
                const data = new Date(`2026-10-31T${ora}:05:00Z`);
                expect(formatData(limba,{hour:'2-digit',minute:'2-digit',timeZone:'UTC',hour12:true}).format(data)).toBe(asteptat);
                const parti = formatData(limba,{dateStyle:'short',timeStyle:'medium',timeZone:'UTC'}).formatToParts(data);
                expect(parti.some(p=>p.type==='dayPeriod')).toBe(false);
                expect(parti.find(p=>p.type==='hour')?.value).toBe(ora);
            }
        });
    }
});

describe('dataAfisata (R5)', () => {
    it('afișează data în limba interfeței, nu AAAA-LL-ZZ', () => {
        const asteptat = new Intl.DateTimeFormat('ro', {dateStyle: 'medium', timeZone: 'UTC'}).format(new Date('2026-10-01T12:00:00Z'));
        expect(dataAfisata('2026-10-01')).toBe(asteptat);
        expect(dataAfisata('2026-10-01')).not.toContain('2026-10');
        expect(dataAfisata('2026-10-01', true)).toContain(new Intl.DateTimeFormat('ro', {weekday: 'short', timeZone: 'UTC'}).format(new Date('2026-10-01T12:00:00Z')));
    });
    it('lasă neschimbate valorile care nu sunt date', () => {
        expect(dataAfisata('')).toBe('');
        expect(dataAfisata('—')).toBe('—');
        expect(dataAfisata('2026-13-45')).toBe('2026-13-45');
    });
});

