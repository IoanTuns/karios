import {describe, expect, it} from 'vitest';
import {formatData, formateazaDataSigura} from './formatare';

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
