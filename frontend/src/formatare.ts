import i18next from 'i18next';
// Formatterele de dată, create o singură dată pentru fiecare limbă și set de opțiuni (F16-T14,
// PERF-2): crearea unui Intl.DateTimeFormat e scumpă, iar antetele grilelor o făceau pentru fiecare
// zi la fiecare randare, adică la fiecare pas de selecție.
const cache = new Map<string, Intl.DateTimeFormat>();

export function formatData(limba: string, optiuni?: Intl.DateTimeFormatOptions): Intl.DateTimeFormat {
    const configuratie: Intl.DateTimeFormatOptions = {...optiuni, hour12: undefined, hourCycle: 'h23'};
    const cheie = limba + '|' + JSON.stringify(configuratie);
    let f = cache.get(cheie);
    if (!f) {
        f = new Intl.DateTimeFormat(limba, configuratie);
        cache.set(cheie, f);
    }
    return f;
}

// formateazaDataSigura (I094): formatează defensiv o dată fără a arunca excepții (precum
// RangeError la date vide sau invalide), returnând valoarea de fallback (implicit '—').
export function formateazaDataSigura(
    valoare: string | number | Date | null | undefined,
    limba: string,
    optiuni?: Intl.DateTimeFormatOptions,
    fallback: string = '—'
): string {
    if (valoare === null || valoare === undefined || valoare === '') {
        return fallback;
    }
    const d = valoare instanceof Date ? valoare : new Date(valoare);
    if (Number.isNaN(d.getTime())) {
        return fallback;
    }
    try {
        return formatData(limba, optiuni).format(d);
    } catch {
        return fallback;
    }
}

// dataAfisata (R5, analiza UX 0.92.59): o data AAAA-LL-ZZ afisata in limba interfetei,
// „1 oct. 2026”, sau cu ziua saptamanii („joi, 1 oct. 2026”). Valorile invalide raman cum sunt.
export function dataAfisata(iso: string, cuZi = false): string {
    if (!/^\d{4}-\d{2}-\d{2}$/.test(iso || '')) return iso || '';
    const d = new Date(`${iso}T12:00:00Z`);
    if (Number.isNaN(d.getTime())) return iso;
    const opt: Intl.DateTimeFormatOptions = cuZi
        ? {weekday: 'short', day: 'numeric', month: 'short', year: 'numeric', timeZone: 'UTC'}
        : {dateStyle: 'medium', timeZone: 'UTC'};
    try { return formatData(i18next.language || 'ro', opt).format(d); } catch { return iso; }
}

