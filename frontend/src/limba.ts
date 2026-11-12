// F03-T05 (D103): limba interfetei si limba implicita a rapoartelor se retin local, in
// stocarea WebView din karios/webview, langa .exe, ca tema (nu sunt date personale si trebuie
// sa fie cunoscute inainte de deschiderea bazei).
import i18next from 'i18next';
import {SeteazaLimbaInterfata} from '../wailsjs/go/main/App';

export const CHEIE_LIMBA_INTERFATA = 'pontix-limba-interfata';
export const CHEIE_LIMBA_RAPOARTE = 'pontix-limba-rapoarte';
/** INTERFATA: rapoartele folosesc limba interfetei (implicit). */
export const LIMBA_RAPOARTE_INTERFATA = 'interfata';

function citeste(cheie: string): string | null {
    try { return localStorage.getItem(cheie); } catch { return null; }
}
function scrie(cheie: string, valoare: string) {
    try { localStorage.setItem(cheie, valoare); } catch { /* fara stocare: alegerea tine pana la inchidere */ }
}

/** limbaInterfataSalvata intoarce limba aleasa anterior, daca e una disponibila. */
export function limbaInterfataSalvata(disponibile: string[]): string | undefined {
    const l = citeste(CHEIE_LIMBA_INTERFATA);
    return l && disponibile.includes(l) ? l : undefined;
}

/** alegeLimbaInterfata schimba limba interfetei si o retine. */
export function alegeLimbaInterfata(limba: string) {
    void i18next.changeLanguage(limba);
    document.documentElement.lang = limba;
    document.title = i18next.t('aplicatie.titlu');
    scrie(CHEIE_LIMBA_INTERFATA, limba);
    SeteazaLimbaInterfata(limba).catch(() => {});
}

/** preferintaLimbaRapoarte: 'interfata', 'ro' sau 'en'. */
export function preferintaLimbaRapoarte(): string {
    const l = citeste(CHEIE_LIMBA_RAPOARTE);
    return l === 'ro' || l === 'en' ? l : LIMBA_RAPOARTE_INTERFATA;
}
export function alegeLimbaRapoarte(valoare: string) {
    scrie(CHEIE_LIMBA_RAPOARTE, valoare === 'ro' || valoare === 'en' ? valoare : LIMBA_RAPOARTE_INTERFATA);
}

/** limbaRapoarte intoarce limba in care se genereaza implicit rapoartele. */
export function limbaRapoarte(): string {
    const p = preferintaLimbaRapoarte();
    if (p === 'ro' || p === 'en') return p;
    return i18next.language === 'en' ? 'en' : 'ro';
}
