// Diferențele pe câmpuri afișate în ecranul jurnalului (F12-T02).
const CAMPURI_TEHNICE = new Set(['revizie', 'modificat_la', 'creat_la']);

type Valoare = unknown;

export function parseaza(text: string): Valoare {
    if (!text) return undefined;
    try {
        return JSON.parse(text);
    } catch {
        return text;
    }
}

export function scurt(v: Valoare): string {
    if (v === undefined || v === null || v === '') return '—';
    const s = typeof v === 'string' ? v : JSON.stringify(v);
    return s.length > 80 ? s.slice(0, 77) + '…' : s;
}

// diferente întoarce câmpurile schimbate între valoarea veche și cea nouă.
export function diferente(vechi: string, nou: string): {camp: string; inainte: string; dupa: string}[] {
    const a = parseaza(vechi), b = parseaza(nou);
    const obiect = (x: Valoare): x is Record<string, Valoare> => typeof x === 'object' && x !== null && !Array.isArray(x);
    if (!obiect(a) || !obiect(b)) return [{camp: '', inainte: scurt(a), dupa: scurt(b)}];
    const chei = Array.from(new Set([...Object.keys(a), ...Object.keys(b)])).sort();
    return chei
        .filter(k => !CAMPURI_TEHNICE.has(k) && JSON.stringify(a[k]) !== JSON.stringify(b[k]))
        .map(k => ({camp: k, inainte: scurt(a[k]), dupa: scurt(b[k])}));
}

export function rezumatNou(nou: string): string {
    const b = parseaza(nou);
    if (typeof b !== 'object' || b === null || Array.isArray(b)) return scurt(b);
    return Object.entries(b as Record<string, Valoare>)
        .filter(([k, v]) => !CAMPURI_TEHNICE.has(k) && v !== '' && v !== null && v !== undefined)
        .slice(0, 4).map(([k, v]) => `${k}: ${scurt(v)}`).join(' · ');
}
