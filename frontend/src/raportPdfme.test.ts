import {describe, expect, it, vi} from 'vitest';
import {getDefaultFont} from '@pdfme/common';
import {PDFDocument} from '@pdfme/pdf-lib';

// Fonturile Go vin din aplicatie; in test folosim fontul inclus in pdfme sub aceleasi nume.
vi.mock('../wailsjs/go/main/App', () => {
    const roboto = Object.values(getDefaultFont())[0].data as Uint8Array;
    let bin = '';
    for (const o of roboto) bin += String.fromCharCode(o);
    const b64 = btoa(bin);
    return {
        FonturiRapoarte: async () => ({Go: b64, GoBold: b64, GoItalic: b64, GoBoldItalic: b64}),
        FontSistem: vi.fn(async (nume: string) => (nume === 'Arial' ? {nume, date: b64, subset: true} : {nume, date: '', subset: false, eroare: {cod: 'fonturi.inexistent'}})),
        FinalizeazaRaportPdfme: vi.fn(async (token: string, pdf: string) => ({token, pdf, pagini: 1})),
    };
});
const notificari: string[] = [];
vi.mock('./Toast', () => ({notifica: (m: string) => notificari.push(m)}));
vi.mock('i18next', () => ({default: {t: (k: string, p?: Record<string, string>) => k + (p ? ' ' + JSON.stringify(p) : '')}}));
vi.mock('../wailsjs/go/models', () => ({main: {RaspunsRaport: class { constructor(v: object) { Object.assign(this, v); } }}}));

import {completeazaRaport, fonturiDinSablon, fonturiPentruSablon, genereazaPdfPdfme, inBase64, dinBase64, sablonDinColoane} from './raportPdfme';

const coloane = [{cod: 'angajat', titlu: 'Angajat'}, {cod: 'data', titlu: 'Data'}, {cod: 'ore', titlu: 'Ore'}];
const pagina = {format: 'A4', orientare: 'PORTRET', margine_mm: 10, font_pt: 9};

describe('sablon pdfme', () => {
    it('construieste tabelul din coloanele SQL si pagina din model', async () => {
        const s = JSON.parse(await sablonDinColoane(coloane, {...pagina, orientare: 'PEISAJ'}, 'Pagina {currentPage} din {totalPages}'));
        expect(s.basePdf.width).toBe(297);
        expect(s.basePdf.height).toBe(210);
        const tabel = s.schemas[0].find((e: {name: string}) => e.name === 'tabel');
        expect(tabel.type).toBe('table');
        expect(tabel.head).toEqual(['Angajat', 'Data', 'Ore']);
        expect(tabel.headWidthPercentages.reduce((a: number, b: number) => a + b, 0)).toBeCloseTo(100);
        expect(s.basePdf.staticSchema.some((e: {content: string}) => e.content.includes('{currentPage}'))).toBe(true);
    }, 30000);

    it('genereaza un PDF cu diacritice si continua tabelul pe pagini noi', async () => {
        const sablon = await sablonDinColoane(coloane, pagina, 'Pagina {currentPage} din {totalPages}');
        const randuri = Array.from({length: 120}, (_, i) => [`Șerban Țîru ${i}`, '2026-11-05', String(i % 24)]);
        const {pdf, lipsa} = await genereazaPdfPdfme(sablon, {titlu: 'Gărzi noiembrie', antet_stanga: 'Spital', antet_centru: '', antet_dreapta: '', subsol_stanga: '', subsol_centru: '', tabel: JSON.stringify(randuri)});
        const text = new TextDecoder('latin1').decode(pdf);
        expect(text.startsWith('%PDF-')).toBe(true);
        expect(lipsa).toEqual([]);
        expect((await PDFDocument.load(pdf)).getPageCount()).toBeGreaterThan(1);
    }, 30000);

    it('trimite PDF-ul generat la finalizare si lasa neschimbate raspunsurile fara pdfme', async () => {
        const fara = {token: 't', generare_pdfme: false} as never;
        expect(await completeazaRaport(fara)).toBe(fara);
        const sablon = await sablonDinColoane(coloane, pagina, '{currentPage}');
        const r = await completeazaRaport({token: 'abc', generare_pdfme: true, sablon, intrari: {tabel: '[["a","b","c"]]'}} as never) as unknown as {token: string; pdf: string};
        expect(r.token).toBe('abc');
        expect(new TextDecoder('latin1').decode(dinBase64(r.pdf)).startsWith('%PDF-')).toBe(true);
    }, 30000);

    it('codifica base64 pe bucati mari fara depasirea stivei', () => {
        const date = new Uint8Array(300000).map((_, i) => i % 256);
        expect(dinBase64(inBase64(date))).toEqual(date);
    });

    it('intoarce eroare traductibila cand sablonul nu se poate genera', async () => {
        const r = await completeazaRaport({token: 'x', generare_pdfme: true, sablon: '{', intrari: {}} as never) as unknown as {eroare: {cod: string}};
        expect(r.eroare.cod).toBe('rapoarte.pdfme_eroare_generare');
    });

    it('incarca fonturile Windows folosite si inlocuieste cu Go pe cele lipsa', async () => {
        const s = JSON.parse(await sablonDinColoane(coloane, pagina, '{currentPage}'));
        s.schemas[0][0].fontName = 'Arial';
        s.schemas[0][1].fontName = 'Font Inexistent';
        expect([...fonturiDinSablon(s)].sort()).toEqual(['Arial', 'Font Inexistent', 'Go', 'GoBold']);
        const {font, lipsa} = await fonturiPentruSablon(s);
        expect(Object.keys(font)).toContain('Arial');
        expect(lipsa).toEqual(['Font Inexistent']);
        expect(s.schemas[0][1].fontName).toBe('Go');
        expect(s.schemas[0][0].fontName).toBe('Arial');
    }, 30000);

    it('genereaza raportul si avertizeaza cand un font lipseste', async () => {
        const s = JSON.parse(await sablonDinColoane(coloane, pagina, '{currentPage}'));
        s.schemas[0][3].fontName = 'Font Inexistent';
        notificari.length = 0;
        const r = await completeazaRaport({token: 'f', generare_pdfme: true, sablon: JSON.stringify(s), intrari: {titlu: 'T', tabel: '[["a","b","c"]]'}} as never) as unknown as {pdf: string};
        expect(r.pdf).toBeTruthy();
        expect(notificari).toEqual(['rapoarte.pdfme_fonturi_lipsa {"fonturi":"Font Inexistent"}']);
    }, 30000);
});
