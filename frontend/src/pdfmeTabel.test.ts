import {describe, expect, it, vi} from 'vitest';
import {getDefaultFont} from '@pdfme/common';
import {PDFDocument} from '@pdfme/pdf-lib';

// Fonturile aplicatiei; in test, fontul inclus in pdfme sub numele Go*.
vi.mock('../wailsjs/go/main/App', () => {
    const roboto = Object.values(getDefaultFont())[0].data as Uint8Array;
    let bin = '';
    for (const o of roboto) bin += String.fromCharCode(o);
    const b64 = btoa(bin);
    return {
        FonturiRapoarte: async () => ({Go: b64, GoBold: b64, GoItalic: b64, GoBoldItalic: b64}),
        FontSistem: vi.fn(async (nume: string) => ({nume, date: '', subset: false, eroare: {cod: 'fonturi.inexistent'}})),
        FinalizeazaRaportPdfme: vi.fn(),
    };
});
vi.mock('./Toast', () => ({notifica: () => {}}));
vi.mock('../wailsjs/go/models', () => ({main: {RaspunsRaport: class {}}}));

import {areStiluri, coloaneVizibile, numeVarianta, stilFinal, tabelCuStiluri, variantaNecesare} from './pdfmeTabel';
import {genereazaPdfPdfme, sablonDinColoane, stiluriDinIntrari} from './raportPdfme';

const baza = {fontName: 'Go', alignment: 'left', verticalAlignment: 'middle', fontSize: 9, lineHeight: 1, characterSpacing: 0, fontColor: '#000000', backgroundColor: '', borderColor: '#888888', borderWidth: {top: .1, right: .1, bottom: .1, left: .1}, padding: {top: 1, right: 1, bottom: 1, left: 1}};

describe('stiluri pe rand si celula in tabelul pdfme', () => {
    it('alege varianta aldin/cursiv a fontului', () => {
        expect(numeVarianta('Go', true, false, {})).toBe('GoBold');
        expect(numeVarianta('GoBold', false, true, {})).toBe('GoItalic');
        expect(numeVarianta('Go', true, true, {})).toBe('GoBoldItalic');
        const f = {'Arial': {data: ''}, 'Arial Bold': {data: ''}} as never;
        expect(numeVarianta('Arial', true, false, f)).toBe('Arial Bold');
        // Varianta lipsa: ramane fontul initial, fara eroare.
        expect(numeVarianta('Arial', false, true, f)).toBe('Arial');
        expect(variantaNecesare(['Arial', 'Go'], {randuri: {'1': {aldin: true}}, celule: {'0:0': {cursiv: true}}}).sort()).toEqual(['Arial Bold', 'Arial Italic']);
    });

    it('combina stilurile de la sectiune la celula, cel mai specific castiga', () => {
        const s = stilFinal(baza, [{aldin: true, font_pt: 11}, {aliniere: 'R'}, {aldin: false, cursiv: true, fundal: '#eeeeee'}], {});
        expect(s.fontName).toBe('GoItalic');
        expect(s.fontSize).toBe(11);
        expect(s.alignment).toBe('right');
        expect(s.backgroundColor).toBe('#eeeeee');
        expect(stilFinal(baza, [undefined], {})).toEqual(baza);
        // Fundal deschis peste text alb: textul trece pe negru.
        expect(stilFinal({...baza, fontColor: '#ffffff', backgroundColor: '#2980ba'}, [{fundal: '#E7EAE7'}], {}).fontColor).toBe('#000000');
        expect(stilFinal({...baza, fontColor: '#ffffff', backgroundColor: '#2980ba'}, [{aldin: true}], {}).fontColor).toBe('#ffffff');
        expect(areStiluri({})).toBe(false);
        expect(areStiluri({randuri: {'3': {aldin: true}}})).toBe(true);
    });

    it('foloseste desenarea originala cand tabelul nu are stiluri', async () => {
        const original = {pdf: vi.fn(async () => {}), ui: vi.fn(), propPanel: {}} as never;
        const p = tabelCuStiluri(original, {} as never, vi.fn() as never, vi.fn() as never, {});
        await p.pdf({schema: {name: 'tabel', type: 'table'}} as never);
        expect((original as {pdf: ReturnType<typeof vi.fn>}).pdf).toHaveBeenCalledOnce();
    });

    it('citeste stilurile trimise de Go doar pentru tabelul de date', () => {
        expect(stiluriDinIntrari({})).toEqual({});
        expect(stiluriDinIntrari({pontix_stiluri: '{'})).toEqual({});
        expect(stiluriDinIntrari({pontix_stiluri: '{"randuri":{"2":{"aldin":true}}}'})).toEqual({tabel: {randuri: {'2': {aldin: true}}}});
    });

    it('genereaza PDF cu randuri aldine pe mai multe pagini, cu aceeasi paginare', async () => {
        const coloane = [{cod: 'data', titlu: 'Data'}, {cod: 'angajat', titlu: 'Angajat'}, {cod: 'ore', titlu: 'Ore'}];
        const sablon = await sablonDinColoane(coloane, {format: 'A4', orientare: 'PORTRET', margine_mm: 10, font_pt: 9}, '{currentPage}/{totalPages}');
        const randuri = Array.from({length: 150}, (_, i) => [`2026-11-${String(i % 30 + 1).padStart(2, '0')}`, `Ștefănescu Ana ${i}`, String(i % 24)]);
        const intrari = {titlu: 'Gărzi', tabel: JSON.stringify(randuri)};
        const randuriAldine = Object.fromEntries(randuri.map((_, i) => [String(i), {aldin: true}]).filter((_, i) => i % 7 === 0 || i % 7 === 6));
        const stiluri = JSON.stringify({randuri: randuriAldine, celule: {'0:2': {aliniere: 'R', fundal: '#ffe6e6'}}, cap: {aldin: true}});
        const simplu = await genereazaPdfPdfme(sablon, intrari);
        const stilat = await genereazaPdfPdfme(sablon, {...intrari, pontix_stiluri: stiluri});
        const nSimplu = (await PDFDocument.load(simplu.pdf)).getPageCount();
        const nStilat = (await PDFDocument.load(stilat.pdf)).getPageCount();
        expect(nSimplu).toBeGreaterThan(1);
        expect(nStilat).toBe(nSimplu);
        // Desenarea proprie a produs alt continut decat cea originala.
        expect(stilat.pdf.length).not.toBe(simplu.pdf.length);
    }, 60000);

    it('scoate coloanele ascunse si pastreaza proportiile celorlalte (D090)', () => {
        const {indici, latimi} = coloaneVizibile([20, 10, 10, 10, 50], 200, [2, 3]);
        expect(indici).toEqual([0, 1, 4]);
        expect(latimi.reduce((a, b) => a + b, 0)).toBeCloseTo(200);
        expect(latimi[0] / latimi[1]).toBeCloseTo(2);
        expect(coloaneVizibile([50, 50], 100).latimi).toEqual([50, 50]);
        expect(areStiluri({ascunse: [3]})).toBe(true);
    });

    it('genereaza modelul predefinit de pontaj pentru februarie, cu zilele 29-31 ascunse', async () => {
        const proces = (globalThis as {process?: {env?: Record<string, string | undefined>}}).process;
        const coloane = [{cod: 'angajat', titlu: 'Angajat'}, {cod: 'marca', titlu: 'Marca'}, ...Array.from({length: 31}, (_, i) => ({cod: 'z' + String(i + 1).padStart(2, '0'), titlu: String(i + 1)}))];
        const sablon = JSON.parse(await sablonDinColoane(coloane, {format: 'A3', orientare: 'PEISAJ', margine_mm: 10, font_pt: 8}, 'Pagina {currentPage} din {totalPages}', {}, {ponderi: [12, 5, ...Array(31).fill(3)], titluriDinRaport: true, stilClasic: true}));
        const tabel = sablon.schemas[0].find((e: {name: string}) => e.name === 'tabel');
        expect(tabel.pontixTitluri).toBe(true);
        const randuri = Array.from({length: 60}, (_, i) => [`Ștefănescu Ana ${i}`, `DEMO-${i}`, ...Array.from({length: 31}, (_, z) => (z >= 28 ? '' : z % 7 === 5 ? 'CO' : '8'))]);
        const celule: Record<string, {fundal: string}> = {};
        randuri.forEach((_, i) => { for (let z = 0; z < 28; z++) if (z % 7 === 5 || z % 7 === 6) celule[`${i}:${z + 2}`] = {fundal: '#E7EAE7'}; });
        const capColoane = Object.fromEntries(Array.from({length: 28}, (_, z) => z).filter(z => z % 7 === 5 || z % 7 === 6).map(z => [String(z + 2), {fundal: '#E7EAE7'}]));
        const stiluri = {celule, cap_coloane: capColoane, randuri: {'0': {aldin: true, fundal: '#E7EFEA'}}, ascunse: [30, 31, 32], titluri: ['Nume angajat', 'Marcă', ...Array.from({length: 31}, (_, i) => String(i + 1))]};
        const {pdf} = await genereazaPdfPdfme(JSON.stringify(sablon), {titlu: 'Pontaj februarie', antet_centru: 'PONTAJ · 2026-02', tabel: JSON.stringify(randuri), pontix_stiluri: JSON.stringify(stiluri)});
        expect((await PDFDocument.load(pdf)).getPageCount()).toBeGreaterThan(1);
        const iesire = proces?.env?.PDF_TEST;
        if (iesire) {
            const fs = await import(/* @vite-ignore */ 'node:' + 'fs') as {writeFileSync: (p: string, d: Uint8Array) => void};
            fs.writeFileSync(iesire, pdf);
        }
    }, 60000);
});

