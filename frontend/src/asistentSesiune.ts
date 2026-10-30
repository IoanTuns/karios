// Sesiunea asistentului AI trăiește în afara componentei (F14-T07): conversația și sarcina
// în curs rămân când utilizatorul schimbă ecranul, iar răspunsul vine prin evenimentul
// Wails „ai:sarcina_terminata”, chiar dacă ecranul asistentului nu e deschis.
import {useSyncExternalStore} from 'react';
import i18next from 'i18next';
import {OpresteSarcinaAI, PornesteSarcinaAI, SarcinaAIActiva} from '../wailsjs/go/main/App';
import {ai, datebaza, main} from '../wailsjs/go/models';
import {EventsOn} from '../wailsjs/runtime/runtime';
import {notifica} from './Toast';
import {imparteRepetari} from './repetari';

export interface MesajUI {
 articoleAjutor?: string[];
 versiuneAjutor?: string;
	intentie?: ai.IntentiePlan;
	rezultatPlan?: datebaza.PlanAI;
    id: string;
    rol: 'user' | 'assistant';
    text: string;
    propuneri?: ai.PropunereActiune[];
    propuneriAplicata?: Record<string, boolean>;
    // reaplicari: propunerile identice cu una deja aplicată mai sus (F14-T19).
    reaplicari?: Record<string, boolean>;
}

// RezultatSarcinaAI vine prin evenimentul Wails, nu printr-o metodă legată, deci generatorul
// de legături nu-l include în models.ts; corespunde lui main.RezultatSarcinaAI din Go.
interface RezultatSarcinaAI {
    id: string;
    raspuns?: ai.RaspunsAsistent;
    eroare?: main.MesajEroare;
    oprita: boolean;
}

export interface SesiuneAI {
    conversatie: MesajUI[];
    sarcina: string; // id-ul sarcinii în curs; '' dacă nu rulează nimic
    eroare: string;
    necitit: boolean; // a venit un răspuns cât timp ecranul asistentului era închis
    perioada: {an: number; luna: number}; // luna la care se referă cererile (F14-T10)
}

const CHEIE_STORAGE = 'pontix_asistent_conversatie';
const EVENIMENT = 'ai:sarcina_terminata';

function mesajBunVenit(): MesajUI {
    return {id: 'bun_venit', rol: 'assistant', text: i18next.t('ai.mesaj_bun_venit')};
}

function incarca(): MesajUI[] {
    try {
        const salvat = sessionStorage.getItem(CHEIE_STORAGE);
        const parsed = salvat ? JSON.parse(salvat) : null;
        if (Array.isArray(parsed) && parsed.length > 0) return parsed;
    } catch {
        // stocarea poate lipsi; conversația pornește de la zero
    }
    return [mesajBunVenit()];
}

const azi = new Date();
let stare: SesiuneAI = {conversatie: [], sarcina: '', eroare: '', necitit: false, perioada: {an: azi.getFullYear(), luna: azi.getMonth() + 1}};
let initializata = false;
let vizibil = 0;
const abonati = new Set<() => void>();
// Un răspuns calculat direct din cod poate sosi înainte ca PornesteSarcinaAI să întoarcă id-ul.
const sositeDevreme = new Map<string, RezultatSarcinaAI>();
const PORNIRE = 'pornire';

function seteaza(schimbare: Partial<SesiuneAI>) {
    stare = {...stare, ...schimbare};
    if (schimbare.conversatie) {
        try {
            sessionStorage.setItem(CHEIE_STORAGE, JSON.stringify(stare.conversatie));
        } catch {
            // ignoră
        }
    }
    abonati.forEach(f => f());
}

function textEroare(e: main.MesajEroare): string {
    const text = i18next.t(e.cod, e.param ?? {});
    return text !== e.cod ? text : (e.param?.detalii ?? e.cod);
}

// faraRepetari scoate propunerile identice cu unele deja făcute sau încă afișate (F14-T09).
function faraRepetari(propuneri: ai.PropunereActiune[] | undefined) {
    return imparteRepetari(propuneri, stare.conversatie);
}

function laTerminare(rez: RezultatSarcinaAI) {
    if (!rez) return;
    if (rez.id !== stare.sarcina) {
        if (stare.sarcina === PORNIRE) sositeDevreme.set(rez.id, rez);
        return;
    }
    const schimbare: Partial<SesiuneAI> = {sarcina: '', eroare: ''};
    if (rez.oprita) {
        schimbare.eroare = i18next.t('ai.sarcina_oprita');
    } else if (rez.eroare) {
        schimbare.eroare = textEroare(rez.eroare);
    } else if (rez.raspuns) {
        const {noi, omise, reaplicari} = faraRepetari(rez.raspuns.propuneri);
        let text = rez.raspuns.mesaj || (noi.length ? i18next.t('ai.propuneri_titlu') : '');
        if (omise > 0) text = [noi.length ? text : '', i18next.t('ai.propuneri_repetate', {numar: String(omise)})].filter(Boolean).join('\n');
        schimbare.conversatie = [...stare.conversatie, {
            id: `a_${Date.now()}`,
            rol: 'assistant',
            text,
            propuneri: noi,
            propuneriAplicata: {},
            reaplicari,
			intentie: rez.raspuns.intentie,
            articoleAjutor: rez.raspuns.articole_ajutor,
            versiuneAjutor: rez.raspuns.versiune_ajutor,
        }];
    }
    if (vizibil === 0 && !rez.oprita) {
        schimbare.necitit = true;
        notifica(i18next.t(rez.eroare ? 'ai.sarcina_eroare' : 'ai.sarcina_terminata'), rez.eroare ? 'eroare' : 'info');
    }
    seteaza(schimbare);
}

// initializeazaSesiuneAI se apelează o dată, din App, după încărcarea traducerilor.
export function initializeazaSesiuneAI() {
    if (initializata) return;
    initializata = true;
    stare = {...stare, conversatie: incarca()};
    EventsOn(EVENIMENT, laTerminare);
    // După o reîncărcare a interfeței, o sarcină pornită anterior poate rula încă.
    SarcinaAIActiva().then(id => { if (id) seteaza({sarcina: id}); }).catch(() => {});
}

function aboneaza(f: () => void) {
    abonati.add(f);
    return () => { abonati.delete(f); };
}

export function useSesiuneAI(): SesiuneAI {
    return useSyncExternalStore(aboneaza, () => stare);
}

// marcheazaVizibil ține evidența ecranului asistentului; la deschidere, răspunsurile devin citite.
export function marcheazaVizibil(da: boolean) {
    vizibil += da ? 1 : -1;
    if (da && stare.necitit) seteaza({necitit: false});
}

export function seteazaConversatie(f: (c: MesajUI[]) => MesajUI[]) {
    seteaza({conversatie: f(stare.conversatie)});
}

export function seteazaPerioada(perioada: {an: number; luna: number}) {
    if (perioada.an !== stare.perioada.an || perioada.luna !== stare.perioada.luna) seteaza({perioada});
}

export function seteazaEroare(eroare: string) {
    seteaza({eroare});
}

export function reseteazaConversatie() {
    seteaza({conversatie: [mesajBunVenit()], eroare: ''});
}

export async function pornesteSarcina(cerere: main.CerereInterogareAI): Promise<void> {
    if (stare.sarcina) return;
    seteaza({eroare: '', sarcina: PORNIRE});
    try {
        const r = await PornesteSarcinaAI(cerere);
        if (r.eroare) {
            seteaza({eroare: textEroare(r.eroare), sarcina: ''});
            return;
        }
        seteaza({sarcina: r.id});
        const sosit = sositeDevreme.get(r.id);
        sositeDevreme.clear();
        if (sosit) laTerminare(sosit);
    } catch {
        seteaza({eroare: i18next.t('date.eroare'), sarcina: ''});
    }
}

export async function opresteSarcina(): Promise<void> {
    try {
        await OpresteSarcinaAI();
    } catch {
        // dacă oprirea eșuează, evenimentul de terminare vine oricum
    }
}
