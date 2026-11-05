// Rapoarte cu sablon pdfme (D087): editorul si generarea PDF-ului ruleaza in
// interfata; Go pregateste datele din SQL si adauga anexa la finalizare.
import type {Font, Template} from '@pdfme/common';
import {FinalizeazaRaportPdfme, FonturiRapoarte} from '../wailsjs/go/main/App';
import {main} from '../wailsjs/go/models';

// pdfme (cativa MB) se incarca la cerere, numai cand un raport il foloseste.
// Aceleasi tipuri ca lista acceptata de Go (internal/rapoarte/pdfme.go).
export async function pluginuriPdfme() {
    const {ellipse, image, line, multiVariableText, rectangle, table, text} = await import('@pdfme/schemas');
    return {Text: text, 'Multi-Variable Text': multiVariableText, Table: table, Image: image, Line: line, Rectangle: rectangle, Ellipse: ellipse};
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

export async function genereazaPdfPdfme(sablon: string, intrari: Record<string, string>): Promise<Uint8Array> {
    const template = JSON.parse(sablon) as Template;
    const [{generate}, plugins, font] = await Promise.all([import('@pdfme/generator'), pluginuriPdfme(), fonturiPdfme()]);
    return generate({template, inputs: [intrari], plugins, options: {font}});
}

// completeazaRaport transforma un raspuns care cere generare pdfme intr-un
// raspuns cu PDF final; celelalte raspunsuri trec neschimbate.
export async function completeazaRaport(r: main.RaspunsRaport): Promise<main.RaspunsRaport> {
    if (r.eroare || !r.generare_pdfme) return r;
    let pdf: Uint8Array;
    try {
        pdf = await genereazaPdfPdfme(r.sablon, r.intrari || {});
    } catch {
        return new main.RaspunsRaport({...r, eroare: {cod: 'rapoarte.pdfme_eroare_generare'}});
    }
    return FinalizeazaRaportPdfme(r.token, inBase64(pdf));
}

export interface PaginaRaport {format: string; orientare: string; margine_mm: number; font_pt: number}
export interface ColoanaSablon {cod: string; titlu: string}

type SchemaPdfme = Record<string, unknown>;

function dinPlugin(plugin: {propPanel: {defaultSchema: SchemaPdfme}}, extra: SchemaPdfme): SchemaPdfme {
    return {...structuredClone(plugin.propPanel.defaultSchema), ...extra};
}

// sablonDinColoane face un sablon de pornire: antet in trei zone, titlu,
// tabelul cu coloanele SQL si subsol cu numarul paginii pe fiecare pagina.
export async function sablonDinColoane(coloane: ColoanaSablon[], pagina: PaginaRaport, textPagina: string): Promise<string> {
    const {table, text} = await pluginuriPdfme().then(p => ({table: p.Table, text: p.Text}));
    const [l, h] = pagina.format === 'A3' ? [297, 420] : [210, 297];
    const [W, H] = pagina.orientare === 'PEISAJ' ? [h, l] : [l, h];
    const m = Math.max(5, pagina.margine_mm || 10);
    const util = W - 2 * m;
    const font = Math.max(6, pagina.font_pt || 9);
    const textFix = (name: string, content: string, x: number, y: number, width: number, alignment: string, extra: SchemaPdfme = {}) =>
        dinPlugin(text as never, {name, content, readOnly: true, position: {x, y}, width, height: 8, alignment, fontSize: font, fontName: 'Go', ...extra});
    const treime = util / 3;
    const antet = [
        textFix('antet_stanga_txt', '{antet_stanga}', m, m, treime, 'left'),
        textFix('antet_centru_txt', '{antet_centru}', m + treime, m, treime, 'center'),
        textFix('antet_dreapta_txt', '{antet_dreapta}', m + 2 * treime, m, treime, 'right'),
        textFix('titlu_txt', '{titlu}', m, m + 10, util, 'center', {height: 10, fontSize: font + 5, fontName: 'GoBold'}),
    ];
    const latime = coloane.length ? 100 / coloane.length : 100;
    const tabel = dinPlugin(table as never, {
        name: TABEL_PDFME,
        position: {x: m, y: m + 24},
        width: util,
        height: 20,
        showHead: true,
        repeatHead: true,
        head: coloane.map(c => c.titlu),
        headWidthPercentages: coloane.map(() => latime),
        content: JSON.stringify([coloane.map(() => '')]),
    });
    for (const [cheie, fontName] of [['headStyles', 'GoBold'], ['bodyStyles', 'Go']] as const) {
        tabel[cheie] = {...(tabel[cheie] as object), fontName, fontSize: font};
    }
    const jos = H - m - 6;
    const staticSchema = [
        textFix('subsol_stanga_txt', '{subsol_stanga}', m, jos, treime, 'left', {fontSize: 8}),
        textFix('subsol_centru_txt', '{subsol_centru}', m + treime, jos, treime, 'center', {fontSize: 8}),
        textFix('pagina_txt', textPagina, m + 2 * treime, jos, treime, 'right', {fontSize: 8}),
    ];
    const sablon = {basePdf: {width: W, height: H, padding: [m, m, m + 8, m], staticSchema}, schemas: [[...antet, tabel]]};
    return JSON.stringify(sablon);
}
