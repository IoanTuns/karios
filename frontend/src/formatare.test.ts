import {describe, expect, it} from 'vitest';
import {formatData} from './formatare';

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
