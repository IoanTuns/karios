// Rapoarte cu sablon pdfme (D087): editorul si generarea PDF-ului ruleaza in
// interfata; Go pregateste datele din SQL si adauga anexa la finalizare.
import type {Font, Template} from '@pdfme/common';
import i18next from 'i18next';
import {FinalizeazaRaportPdfme, FontSistem, FonturiRapoarte} from '../wailsjs/go/main/App';
import {notifica} from './Toast';
import {main} from '../wailsjs/go/models';

import cheiPdfme from './cheiPdfme.json';
import {areStiluri, INTRARE_STILURI, tabelCuStiluri, TITLURI_DIN_RAPORT, variantaNecesare, type StiluriTabel} from './pdfmeTabel';

// Tipurile de elemente, in ordinea din bara editorului; numele afisate vin din cataloage (UIR-11).
export const TIPURI_PDFME = ['text', 'multiVariableText', 'table', 'image', 'line', 'rectangle', 'ellipse'] as const;
export type TipPdfme = typeof TIPURI_PDFME[number];

// pdfme (cativa MB) se incarca la cerere, numai cand un raport il foloseste.
// Aceleasi tipuri ca lista acceptata de Go (internal/rapoarte/pdfme.go). Cheile
// obiectului sunt etichetele din bara editorului; generarea nu le foloseste.
export async function pluginuriPdfme(etichete?: Partial<Record<TipPdfme, string>>) {
    const m = await import('@pdfme/schemas');
    const plugin: Record<TipPdfme, unknown> = {text: m.text, multiVariableText: m.multiVariableText, table: m.table, image: m.image, line: m.line, rectangle: m.rectangle, ellipse: m.ellipse};
    return Object.fromEntries(TIPURI_PDFME.map(tip => [etichete?.[tip] || tip, plugin[tip]])) as Record<string, never>;
}

// etichetePdfme traduce interfata editorului din cataloagele aplicatiei (UIR-06, D004).
export function etichetePdfme(t: (k: string) => string): Record<string, string> {
    return Object.fromEntries((cheiPdfme as string[]).map(k => [k, t('pdfme.' + k)]));
}

export const TABEL_PDFME = 'tabel';

export function dinBase64(b64: string): Uint8Array<ArrayBuffer> {
    const s = atob(b64);
    const r = new Uint8Array(s.length);
    for (let i = 0; i < s.length; i++) r[i] = s.charCodeAt(i);
    return r;
}

export function inBase64(octeti: Uint8Array): string {
    let s = '';
    for (let i = 0; i < octeti.length; i += 0x8000) s += String.fromCharCode(...octeti.subarray(i, i + 0x8000));
    return btoa(s);
}

let fonturi: Promise<Font> | null = null;

// Fonturile Go vin din aplicatie, nu de pe internet; Go Regular este rezerva.
export function fonturiPdfme(): Promise<Font> {
    if (!fonturi) {
        fonturi = FonturiRapoarte().then(f => {
            const rez: Font = {};
            for (const [nume, b64] of Object.entries(f)) rez[nume] = {data: dinBase64(b64), fallback: nume === 'Go', subset: true};
            return rez;
        }).catch(e => {
            fonturi = null;
            throw e;
        });
    }
    return fonturi;
}

const fonturiSistem = new Map<string, Promise<Font[string] | null>>();

// fontSistem incarca un font Windows dupa nume, o singura data (D088).
export function fontSistem(nume: string): Promise<Font[string] | null> {
    let f = fonturiSistem.get(nume);
    if (!f) {
        f = FontSistem(nume).then(r => (r.eroare || !r.date ? null : {data: dinBase64(r.date), subset: r.subset, fallback: false}));
        f.then(v => { if (!v) fonturiSistem.delete(nume); }, () => fonturiSistem.delete(nume));
        fonturiSistem.set(nume, f);
    }
    return f;
}

// fonturiDinSablon aduna toate valorile fontName din sablon (elemente, pagina fixa, stiluri de tabel).
export function fonturiDinSablon(v: unknown, rez = new Set<string>()): Set<string> {
    if (Array.isArray(v)) v.forEach(x => fonturiDinSablon(x, rez));
    else if (v && typeof v === 'object') {
        for (const [k, x] of Object.entries(v)) {
            if (k === 'fontName' && typeof x === 'string' && x) rez.add(x);
            else fonturiDinSablon(x, rez);
        }
    }
    return rez;
}

function inlocuiesteFonturi(v: unknown, lipsa: Set<string>): void {
    if (Array.isArray(v)) v.forEach(x => inlocuiesteFonturi(x, lipsa));
    else if (v && typeof v === 'object') {
        const o = v as Record<string, unknown>;
        for (const [k, x] of Object.entries(o)) {
            if (k === 'fontName' && typeof x === 'string' && lipsa.has(x)) o[k] = 'Go';
            else inlocuiesteFonturi(x, lipsa);
        }
    }
}

// fonturiPentruSablon pregateste fonturile Go si fonturile Windows folosite de
// sablon. Cele care lipsesc pe acest calculator se inlocuiesc cu Go, in sablon.
export async function fonturiPentruSablon(template: Template): Promise<{font: Font; lipsa: string[]}> {
    const font: Font = {...await fonturiPdfme()};
    const lipsa = new Set<string>();
    for (const nume of fonturiDinSablon(template)) {
        if (font[nume]) continue;
        const f = await fontSistem(nume).catch(() => null);
        if (f) font[nume] = f;
        else lipsa.add(nume);
    }
    inlocuiesteFonturi(template, lipsa);
    return {font, lipsa: [...lipsa]};
}

// stiluriDinIntrari citeste stilurile pe rand/celula trimise de Go (F09-T22).
export function stiluriDinIntrari(intrari: Record<string, string>): Record<string, StiluriTabel> {
    const text = intrari[INTRARE_STILURI];
    if (!text) return {};
    try {
        const v = JSON.parse(text) as StiluriTabel;
        // Titlurile singure conteaza doar pentru tabelele marcate pontixTitluri.
        return areStiluri(v) || v.titluri?.length ? {[TABEL_PDFME]: v} : {};
    } catch {
        return {};
    }
}

export async function genereazaPdfPdfme(sablon: string, intrari: Record<string, string>): Promise<{pdf: Uint8Array; lipsa: string[]}> {
    const template = JSON.parse(sablon) as Template;
    const [{generate}, schemasMod, tabeleMod, plugins, {font, lipsa}] = await Promise.all([import('@pdfme/generator'), import('@pdfme/schemas'), import('@pdfme/schemas/tables'), pluginuriPdfme(), fonturiPentruSablon(template)]);
    const stiluri = stiluriDinIntrari(intrari);
    if (stiluri[TABEL_PDFME]) {
        // Variantele aldin/cursiv ale fonturilor Windows din tabel se incarca la cerere.
        const tabel = template.schemas.flat().find(e => e.type === 'table' && e.name === TABEL_PDFME) as {headStyles?: {fontName?: string}; bodyStyles?: {fontName?: string}} | undefined;
        const fonturiTabel = [tabel?.headStyles?.fontName, tabel?.bodyStyles?.fontName].filter(Boolean) as string[];
        for (const nume of variantaNecesare(fonturiTabel, stiluri[TABEL_PDFME])) {
            if (font[nume]) continue;
            const f = await fontSistem(nume).catch(() => null);
            if (f) font[nume] = f;
        }
        plugins.table = tabelCuStiluri(schemasMod.table as never, {text: schemasMod.text.pdf, rectangle: schemasMod.rectangle.pdf, line: schemasMod.line.pdf} as never,
            schemasMod.getDynamicHeightsForTable as never, tabeleMod.getTableBodyRange as never, stiluri) as never;
    }
    return {pdf: await generate({template, inputs: [intrari], plugins, options: {font}}), lipsa};
}

// completeazaRaport transforma un raspuns care cere generare pdfme intr-un
// raspuns cu PDF final; celelalte raspunsuri trec neschimbate.
export async function completeazaRaport(r: main.RaspunsRaport): Promise<main.RaspunsRaport> {
    if (r.eroare || !r.generare_pdfme) return r;
    let rez: {pdf: Uint8Array; lipsa: string[]};
    try {
        rez = await genereazaPdfPdfme(r.sablon, r.intrari || {});
    } catch {
        return new main.RaspunsRaport({...r, eroare: {cod: 'rapoarte.pdfme_eroare_generare'}});
    }
    if (rez.lipsa.length) notifica(i18next.t('rapoarte.pdfme_fonturi_lipsa', {fonturi: rez.lipsa.join(', ')}), 'avertisment', 8000);
    return FinalizeazaRaportPdfme(r.token, inBase64(rez.pdf));
}

export interface PaginaRaport {format: string; orientare: string; margine_mm: number; font_pt: number}
export interface ColoanaSablon {cod: string; titlu: string}
// Numele elementelor din sablonul de pornire, afisate in lista elementelor (UIR-06).
export type NumeElementeSablon = Partial<Record<'antet_stanga' | 'antet_centru' | 'antet_dreapta' | 'titlu' | 'subsol_stanga' | 'subsol_centru' | 'pagina', string>>;

// Prima litera mare, fara a schimba restul ("angajat id" -> "Angajat id").
export function titluColoana(c: ColoanaSablon): string {
    const t = (c.titlu || c.cod).trim();
    return t.charAt(0).toLocaleUpperCase() + t.slice(1);
}

type SchemaPdfme = Record<string, unknown>;

function dinPlugin(plugin: {propPanel: {defaultSchema: SchemaPdfme}}, extra: SchemaPdfme): SchemaPdfme {
    return {...structuredClone(plugin.propPanel.defaultSchema), ...extra};
}

// sablonDinColoane face un sablon de pornire: antet in trei zone, titlu,
// tabelul cu coloanele SQL si subsol cu numarul paginii pe fiecare pagina.
// Optiuni pentru modelele predefinite (D090): ponderile coloanelor din raportul
// standard si titlurile coloanelor luate din raport, in limba lui.
export interface OptiuniSablon {ponderi?: number[]; titluriDinRaport?: boolean; stilClasic?: boolean}

export async function sablonDinColoane(coloane: ColoanaSablon[], pagina: PaginaRaport, textPagina: string, nume: NumeElementeSablon = {}, optiuni: OptiuniSablon = {}): Promise<string> {
    const p = await pluginuriPdfme();
    const text = p.text, table = p.table;
    const n = (k: keyof NumeElementeSablon, implicit: string) => nume[k] || implicit;
    const [l, h] = pagina.format === 'A3' ? [297, 420] : [210, 297];
    const [W, H] = pagina.orientare === 'PEISAJ' ? [h, l] : [l, h];
    const m = Math.max(5, pagina.margine_mm || 10);
    const util = W - 2 * m;
    const font = Math.max(6, pagina.font_pt || 9);
    const textFix = (name: string, content: string, x: number, y: number, width: number, alignment: string, extra: SchemaPdfme = {}) =>
        dinPlugin(text as never, {name, content, readOnly: true, position: {x, y}, width, height: 8, alignment, fontSize: font, fontName: 'Go', ...extra});
    const treime = util / 3;
    const antet = [
        textFix(n('antet_stanga', 'antet_stanga_txt'), '{antet_stanga}', m, m, treime, 'left'),
        textFix(n('antet_centru', 'antet_centru_txt'), '{antet_centru}', m + treime, m, treime, 'center'),
        textFix(n('antet_dreapta', 'antet_dreapta_txt'), '{antet_dreapta}', m + 2 * treime, m, treime, 'right'),
        textFix(n('titlu', 'titlu_txt'), '{titlu}', m, m + 10, util, 'center', {height: 10, fontSize: font + 5, fontName: 'GoBold'}),
    ];
    const ponderi = optiuni.ponderi && optiuni.ponderi.length === coloane.length ? optiuni.ponderi : coloane.map(() => 1);
    const totalPonderi = ponderi.reduce((a, b) => a + b, 0) || 1;
    const tabel = dinPlugin(table as never, {
        name: TABEL_PDFME,
        position: {x: m, y: m + 24},
        width: util,
        height: 20,
        showHead: true,
        repeatHead: true,
        head: coloane.map(titluColoana),
        headWidthPercentages: ponderi.map(p => (p * 100) / totalPonderi),
        content: JSON.stringify([coloane.map(() => '')]),
        [COLOANE_SQL]: coloane.map(c => c.cod),
        ...(optiuni.titluriDinRaport ? {[TITLURI_DIN_RAPORT]: true} : {}),
    });
    // Spatiere mica pentru tabelele cu multe coloane (zilele lunii), ca textul sa incapa.
    const spatiu = coloane.length > 12 ? 0.8 : 2;
    const padding = {top: spatiu, right: spatiu, bottom: spatiu, left: spatiu};
    for (const [cheie, fontName] of [['headStyles', 'GoBold'], ['bodyStyles', 'Go']] as const) {
        tabel[cheie] = {...(tabel[cheie] as object), fontName, fontSize: font, padding};
    }
    // Zilele (z01-z31) se centreaza, ca in raportul standard.
    const centrate = Object.fromEntries(coloane.map((c, i) => [i, c.cod]).filter(([, cod]) => /^z\d{2}$/.test(String(cod))).map(([i]) => [i, 'center']));
    if (Object.keys(centrate).length) tabel.columnStyles = {...(tabel.columnStyles as object), alignment: centrate};
    if (optiuni.stilClasic) {
        // Aspectul raportului standard vechi: cap gri-verzui cu text negru, fara alternanta, chenare gri.
        tabel.headStyles = {...(tabel.headStyles as object), backgroundColor: '#E7EFEA', fontColor: '#000000', borderColor: '#8C9690', borderWidth: {top: 0.15, right: 0.15, bottom: 0.15, left: 0.15}};
        tabel.bodyStyles = {...(tabel.bodyStyles as object), alternateBackgroundColor: '', borderColor: '#8C9690', borderWidth: {top: 0.15, right: 0.15, bottom: 0.15, left: 0.15}};
        tabel.tableStyles = {borderColor: '#8C9690', borderWidth: 0.15};
    }
    const jos = H - m - 6;
    const staticSchema = [
        textFix(n('subsol_stanga', 'subsol_stanga_txt'), '{subsol_stanga}', m, jos, treime, 'left', {fontSize: 8}),
        textFix(n('subsol_centru', 'subsol_centru_txt'), '{subsol_centru}', m + treime, jos, treime, 'center', {fontSize: 8}),
        textFix(n('pagina', 'pagina_txt'), textPagina, m + 2 * treime, jos, treime, 'right', {fontSize: 8}),
    ];
    // Antetul si titlul se repeta pe fiecare pagina, ca in rapoartele standard vechi (F09-T26).
    const sablon = {basePdf: {width: W, height: H, padding: [m + 22, m, m + 8, m], staticSchema: [...antet, ...staticSchema]}, schemas: [[tabel]]};
    return JSON.stringify(sablon);
}

// Codurile coloanelor SQL se pastreaza pe tabel, ca tabelul sa poata fi
// actualizat dupa o schimbare a SQL-ului fara a pierde aspectul (UIR-05).
export const COLOANE_SQL = 'pontixColoane';

type TabelSablon = SchemaPdfme & {head?: string[]; headWidthPercentages?: number[]; content?: string; readOnly?: boolean};

function tabelDate(template: {schemas?: SchemaPdfme[][]}): TabelSablon | undefined {
    for (const pagina of template.schemas || []) {
        for (const e of pagina) if (e.type === 'table' && e.name === TABEL_PDFME && !e.readOnly) return e as TabelSablon;
    }
    return undefined;
}

// StareTabel compara tabelul din sablon cu coloanele curente ale SQL-ului.
export interface StareTabel {exista: boolean; potrivit: boolean; adaugate: string[]; eliminate: string[]}

function coduriTabel(tabel: TabelSablon, coloane: ColoanaSablon[]): string[] {
    const salvate = tabel[COLOANE_SQL];
    if (Array.isArray(salvate) && salvate.length === (tabel.head || []).length) return salvate as string[];
    // Sabloane fara coduri: potrivire dupa titlul afisat, apoi o cheie de pozitie.
    return (tabel.head || []).map((h, i) => coloane.find(c => titluColoana(c) === h || c.titlu === h || c.cod === h)?.cod ?? '#' + i);
}

export function stareTabel(sablon: string, coloane: ColoanaSablon[]): StareTabel {
    const tabel = tabelDate(JSON.parse(sablon));
    if (!tabel) return {exista: false, potrivit: false, adaugate: coloane.map(c => c.cod), eliminate: []};
    // F09-T27: tabelele marcate pontixTitluri (modelele standard si cele convertite) urmeaza
    // singure coloanele SQL la generare (SablonPentruDocument), deci nu sunt nepotrivite.
    if ((tabel as Record<string, unknown>)['pontixTitluri'] === true) return {exista: true, potrivit: true, adaugate: [], eliminate: []};
    const vechi = coduriTabel(tabel, coloane), noi = coloane.map(c => c.cod);
    const adaugate = noi.filter(c => !vechi.includes(c)), eliminate = vechi.filter(c => !noi.includes(c));
    return {exista: true, potrivit: !adaugate.length && !eliminate.length && vechi.join('|') === noi.join('|'), adaugate, eliminate};
}

// actualizeazaTabel aduce tabelul la coloanele SQL curente, in ordinea lor:
// coloanele pastrate isi pastreaza titlul si latimea, cele noi primesc latimea
// medie, iar totalul se normalizeaza la 100%. Restul paginii ramane neatins.
export function actualizeazaTabel(sablon: string, coloane: ColoanaSablon[]): string {
    const template = JSON.parse(sablon);
    const tabel = tabelDate(template);
    if (!tabel) throw new Error('fara tabel');
    const vechi = coduriTabel(tabel, coloane);
    const head = tabel.head || [], latimi = tabel.headWidthPercentages || [];
    const medie = latimi.length ? latimi.reduce((a, b) => a + b, 0) / latimi.length : 100 / Math.max(1, coloane.length);
    const nouHead: string[] = [], nouLatimi: number[] = [];
    for (const c of coloane) {
        const i = vechi.indexOf(c.cod);
        nouHead.push(i >= 0 ? head[i] : titluColoana(c));
        nouLatimi.push(i >= 0 && latimi[i] > 0 ? latimi[i] : medie);
    }
    const total = nouLatimi.reduce((a, b) => a + b, 0) || 1;
    tabel.head = nouHead;
    tabel.headWidthPercentages = nouLatimi.map(l => (l * 100) / total);
    tabel.content = JSON.stringify([coloane.map(() => '')]);
    tabel[COLOANE_SQL] = coloane.map(c => c.cod);
    return JSON.stringify(template);
}

// La zoom 1, pdfme potriveste pagina pe inaltimea zonei (masurat in Edge:
// latimea paginii la 100% = inaltimea utila * latime/inaltime). Pentru
// „potriveste latimea” zoomul este raportul dintre latimea disponibila si
// aceasta latime; limitele sunt cele ale editorului (UIR-09).
export function zoomPotrivit(latimeZona: number, inaltimeZona: number, latimePaginaMm: number, inaltimePaginaMm: number, panouDeschis: boolean): number {
    const disponibil = latimeZona - (panouDeschis ? 420 : 40) - 90;
    const latimeLa100 = Math.max(1, inaltimeZona - 40) * (latimePaginaMm / Math.max(1, inaltimePaginaMm));
    const z = disponibil / latimeLa100;
    return Math.max(0.5, Math.min(1.5, Math.floor(z * 20) / 20));
}

// F09-T29: antetul si subsolul stau in basePdf.staticSchema (se repeta pe fiecare
// pagina), unde Designer-ul pdfme nu le poate selecta. Pentru editare, elementele
// trec temporar pe pagina, marcate STATIC_PDFME, cu marginile reduse la margine.
export const STATIC_PDFME = 'pontixStatic';

type SablonEditabil = {basePdf: {width?: number; height?: number; padding?: number[]; staticSchema?: SchemaPdfme[]}; schemas: SchemaPdfme[][]};

function marginePagina(b: SablonEditabil['basePdf']): number {
    const stanga = b.padding?.[3];
    return typeof stanga === 'number' && stanga >= 0 ? stanga : 10;
}

// antetEditabil muta elementele fixe pe prima pagina, ca sa poata fi mutate,
// modificate sau sterse. Sablonul fara basePdf obiect (PDF de fundal) ramane neschimbat.
export function antetEditabil(sablon: string): string {
    const s = JSON.parse(sablon) as SablonEditabil;
    if (!s.basePdf || typeof s.basePdf !== 'object' || !Array.isArray(s.basePdf.staticSchema) || !s.basePdf.staticSchema.length) return sablon;
    const m = marginePagina(s.basePdf);
    const fixe = s.basePdf.staticSchema.map(e => ({...e, [STATIC_PDFME]: true}));
    s.schemas = s.schemas?.length ? s.schemas : [[]];
    s.schemas[0] = [...fixe, ...s.schemas[0]];
    s.basePdf = {...s.basePdf, staticSchema: [], padding: [m, m, m, m]};
    return JSON.stringify(s);
}

// antetFix pune la loc elementele marcate si pe cele adaugate in timpul editarii
// antetului (numele lipsesc din `inainte`). Marginea de sus si cea de jos se
// recalculeaza din elementele fixe, iar tabelul coboara sub antet daca s-ar suprapune.
export function antetFix(sablon: string, inainte?: Set<string>): string {
    const s = JSON.parse(sablon) as SablonEditabil;
    const pagina = s.schemas?.[0];
    if (!pagina || !s.basePdf || typeof s.basePdf !== 'object') return sablon;
    const eFix = (e: SchemaPdfme) => e[STATIC_PDFME] === true || (!!inainte && !inainte.has(String(e.name)));
    if (!pagina.some(eFix)) return sablon;
    const fixe: SchemaPdfme[] = [], rest: SchemaPdfme[] = [];
    for (const e of pagina) {
        if (!eFix(e)) {
            rest.push(e);
            continue;
        }
        const {[STATIC_PDFME]: _, ...curat} = e;
        // In staticSchema textul se scrie din continut, cu {variabile}, deci e numai citire.
        fixe.push({...curat, readOnly: true});
    }
    const H = s.basePdf.height || 297, m = marginePagina(s.basePdf);
    const p = s.basePdf.padding || [m, m, m, m];
    let sus = m, jos = m;
    for (const e of fixe) {
        const y = Number((e.position as {y?: number})?.y) || 0, h = Number(e.height) || 0;
        if (y + h / 2 < H / 2) sus = Math.max(sus, y + h + 2);
        else jos = Math.max(jos, H - y + 2);
    }
    for (const e of rest) {
        const poz = e.position as {x: number; y: number} | undefined;
        if (e.type === 'table' && poz && poz.y < sus) e.position = {...poz, y: sus};
    }
    s.schemas[0] = rest;
    s.basePdf = {...s.basePdf, staticSchema: [...fixe, ...(s.basePdf.staticSchema || [])], padding: [round2(sus), p[1] ?? m, round2(jos), p[3] ?? m]};
    return JSON.stringify(s);
}

function round2(v: number): number {
    return Math.round(v * 100) / 100;
}

// F09-T32: dimensiunea paginii din Setări de pagină (format, orientare) este cea a
// raportului; sablonul pdfme o urmeaza. La schimbare, elementele se intind pe noua
// latime (x si latimea proportional), iar cele din jumatatea de jos pastreaza distanta
// fata de marginea de jos (subsolul ramane jos). Fonturile si marginile raman.
export function dimensiunePagina(pagina: Pick<PaginaRaport, 'format' | 'orientare'>): [number, number] {
    const [l, h] = pagina.format === 'A3' ? [297, 420] : [210, 297];
    return pagina.orientare === 'PEISAJ' ? [h, l] : [l, h];
}

export function sablonPePagina(sablon: string, pagina: Pick<PaginaRaport, 'format' | 'orientare'>): string {
    const s = JSON.parse(sablon) as {basePdf?: {width?: number; height?: number; staticSchema?: SchemaPdfme[]}; schemas?: SchemaPdfme[][]};
    const b = s.basePdf;
    if (!b || typeof b !== 'object' || !b.width || !b.height) return sablon;
    const [W, H] = dimensiunePagina(pagina);
    if (Math.abs(b.width - W) < 0.01 && Math.abs(b.height - H) < 0.01) return sablon;
    // Se scaleaza zona dintre marginile laterale, nu toata pagina, ca elementele sa nu intre in margini.
    const p = (b as {padding?: number[]}).padding || [];
    const ml = typeof p[3] === 'number' ? p[3] : 10, md = typeof p[1] === 'number' ? p[1] : 10;
    const kx = (W - ml - md) / Math.max(1, b.width - ml - md), dy = H - b.height, jumatate = b.height / 2;
    const r2 = (v: number) => Math.round(v * 100) / 100;
    const muta = (e: SchemaPdfme) => {
        const poz = e.position as {x: number; y: number} | undefined;
        if (!poz) return e;
        const h = Number(e.height) || 0, w = Number(e.width) || 0;
        const jos = poz.y + h / 2 > jumatate;
        // In margini (de ex. o imagine pusa pe margine) elementul ramane la aceeasi distanta de marginea paginii.
        const x = poz.x < ml ? poz.x : ml + (poz.x - ml) * kx;
        const latime = Math.min(w * kx, Math.max(0, W - md - Math.max(x, ml)) || w * kx);
        return {...e, position: {x: r2(x), y: r2(Math.max(0, jos ? poz.y + dy : poz.y))}, width: r2(latime)};
    };
    s.basePdf = {...b, width: W, height: H, ...(b.staticSchema ? {staticSchema: b.staticSchema.map(muta)} : {})};
    s.schemas = (s.schemas || []).map(pg => pg.map(muta));
    return JSON.stringify(s);
}
