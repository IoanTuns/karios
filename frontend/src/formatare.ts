// Formatterele de dată, create o singură dată pentru fiecare limbă și set de opțiuni (F16-T14,
// PERF-2): crearea unui Intl.DateTimeFormat e scumpă, iar antetele grilelor o făceau pentru fiecare
// zi la fiecare randare, adică la fiecare pas de selecție.
const cache = new Map<string, Intl.DateTimeFormat>();

export function formatData(limba: string, optiuni?: Intl.DateTimeFormatOptions): Intl.DateTimeFormat {
    const cheie = limba + '|' + JSON.stringify(optiuni ?? {});
    let f = cache.get(cheie);
    if (!f) {
        f = new Intl.DateTimeFormat(limba, optiuni);
        cache.set(cheie, f);
    }
    return f;
}
