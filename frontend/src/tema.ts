// Tema interfeței (F16-T13, UI-1): până când utilizatorul alege una, urmează setarea din Windows
// (WebView2 transmite tema aplicațiilor prin prefers-color-scheme).
export type Tema = 'light' | 'dark';

export const cheieTema = 'pontix_tema';

// temaInitiala: tema salvată, dacă utilizatorul a ales una; altfel cea a sistemului.
export function temaInitiala(salvata: string | null, sistemIntunecat: boolean): Tema {
    if (salvata === 'light' || salvata === 'dark') return salvata;
    return sistemIntunecat ? 'dark' : 'light';
}
