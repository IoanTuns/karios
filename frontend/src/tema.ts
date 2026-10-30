// Tema interfeței (F16-T13, UI-1): până când utilizatorul alege una, urmează setarea din Windows
// (WebView2 transmite tema aplicațiilor prin prefers-color-scheme).
export type Tema = 'light' | 'dark';

export const cheieTema = 'pontix_tema';

// temaInitiala: tema salvată, dacă utilizatorul a ales una; altfel cea a sistemului.
export function temaInitiala(salvata: string | null, sistemIntunecat: boolean): Tema {
    if (salvata === 'light' || salvata === 'dark') return salvata;
    return sistemIntunecat ? 'dark' : 'light';
}

// aplicaTemaInitiala pune tema pe <html> înainte de primul ecran, ca s-o primească și ecranele
// de dinaintea autentificării (deblocare, cheie de recuperare), care nu trec prin App (F16-T16).
export function aplicaTemaInitiala(): Tema {
    let salvata: string | null = null;
    try { salvata = localStorage.getItem(cheieTema); } catch { /* fără stocare */ }
    const tema = temaInitiala(salvata, window.matchMedia?.('(prefers-color-scheme: dark)').matches ?? false);
    document.documentElement.setAttribute('data-theme', tema);
    return tema;
}
