import {beforeAll, describe, expect, it, vi} from 'vitest';
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

import {actualizeazaTabel, antetEditabil, antetFix, COLOANE_SQL, completeazaRaport, etichetePdfme, fonturiDinSablon, fonturiPentruSablon, genereazaPdfPdfme, inBase64, dinBase64, pluginuriPdfme, sablonDinColoane, stareTabel, titluColoana, zoomPotrivit} from './raportPdfme';
import cheiPdfme from './cheiPdfme.json';
import catalogRo from '../../internal/i18n/locales/ro.json';
import catalogEn from '../../internal/i18n/locales/en.json';

const coloane = [{cod: 'angajat', titlu: 'Angajat'}, {cod: 'data', titlu: 'Data'}, {cod: 'ore', titlu: 'Ore'}];
const pagina = {format: 'A4', orientare: 'PORTRET', margine_mm: 10, font_pt: 9};

describe('sablon pdfme', () => {
    // I116: pdfme (cativa MB) se incarca o singura data, inaintea testelor, cu o
    // limita generoasa; sub sarcina incarcarea poate depasi limita unui test.
    beforeAll(async () => {
        await Promise.all([import('@pdfme/schemas'), import('@pdfme/generator')]);
    }, 120000);

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

    it('F09-T29: antetul trece pe pagina pentru editare si inapoi, cu marginile recalculate', async () => {
        const sablon = await sablonDinColoane(coloane, pagina, 'Pagina {currentPage} din {totalPages}');
        const ed = JSON.parse(antetEditabil(sablon));
        expect(ed.basePdf.staticSchema).toEqual([]);
        expect(ed.basePdf.padding).toEqual([10, 10, 10, 10]);
        expect(ed.schemas[0].filter((e: Record<string, unknown>) => e.pontixStatic === true).map((e: {name: string}) => e.name)).toEqual(['antet_stanga_txt', 'antet_centru_txt', 'antet_dreapta_txt', 'titlu_txt', 'subsol_stanga_txt', 'subsol_centru_txt', 'pagina_txt']);
        // Fara modificari, drumul dus-intors da acelasi sablon.
        expect(JSON.parse(antetFix(JSON.stringify(ed)))).toEqual(JSON.parse(sablon));
        // Titlul sters, data adaugata langa antetul din dreapta, element nou in antet.
        // Ca in editor: numai elementele corpului; marcajul pontixStatic poate lipsi dupa editare.
        const inainte = new Set<string>(['tabel']);
        for (const e of ed.schemas[0]) delete e.pontixStatic;
        ed.schemas[0] = ed.schemas[0].filter((e: {name: string}) => e.name !== 'titlu_txt');
        ed.schemas[0].find((e: {name: string}) => e.name === 'antet_dreapta_txt').content = '{data_curenta}';
        ed.schemas[0].push({name: 'nou', type: 'text', content: 'Aprobat', position: {x: 10, y: 30}, width: 50, height: 8});
        const fix = JSON.parse(antetFix(JSON.stringify(ed), inainte));
        const nume = fix.basePdf.staticSchema.map((e: {name: string}) => e.name);
        expect(nume).toContain('nou');
        expect(nume).not.toContain('titlu_txt');
        expect(fix.basePdf.staticSchema.every((e: Record<string, unknown>) => e.readOnly === true && !('pontixStatic' in e))).toBe(true);
        expect(fix.schemas[0].map((e: {name: string}) => e.name)).toEqual(['tabel']);
        // Elementul nou se termina la 38 mm: marginea de sus 40, tabelul coborat sub antet.
        expect(fix.basePdf.padding[0]).toBe(40);
        expect(fix.schemas[0][0].position.y).toBe(40);
        expect(fix.basePdf.padding[2]).toBe(18);
        const {pdf} = await genereazaPdfPdfme(JSON.stringify(fix), {data_curenta: '07.11.2026', antet_stanga: '', antet_centru: '', antet_dreapta: '', subsol_stanga: '', subsol_centru: '', titlu: '', tabel: JSON.stringify([['a', 'b', 'c']])});
        expect((await PDFDocument.load(pdf)).getPageCount()).toBe(1);
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
        // Antetul (pe fiecare pagina) e in staticSchema: 0 stanga, 1 centru, 2 dreapta, 3 titlu.
        s.basePdf.staticSchema[0].fontName = 'Arial';
        s.basePdf.staticSchema[1].fontName = 'Font Inexistent';
        expect([...fonturiDinSablon(s)].sort()).toEqual(['Arial', 'Font Inexistent', 'Go', 'GoBold']);
        const {font, lipsa} = await fonturiPentruSablon(s);
        expect(Object.keys(font)).toContain('Arial');
        expect(lipsa).toEqual(['Font Inexistent']);
        expect(s.basePdf.staticSchema[1].fontName).toBe('Go');
        expect(s.basePdf.staticSchema[0].fontName).toBe('Arial');
    }, 30000);

    it('genereaza raportul si avertizeaza cand un font lipseste', async () => {
        const s = JSON.parse(await sablonDinColoane(coloane, pagina, '{currentPage}'));
        s.basePdf.staticSchema[3].fontName = 'Font Inexistent';
        notificari.length = 0;
        const r = await completeazaRaport({token: 'f', generare_pdfme: true, sablon: JSON.stringify(s), intrari: {titlu: 'T', tabel: '[["a","b","c"]]'}} as never) as unknown as {pdf: string};
        expect(r.pdf).toBeTruthy();
        expect(notificari).toEqual(['rapoarte.pdfme_fonturi_lipsa {"fonturi":"Font Inexistent"}']);
    }, 30000);

    it('tine minte codurile SQL pe tabel si da nume lizibile elementelor (UIR-04, UIR-06)', async () => {
        const s = JSON.parse(await sablonDinColoane([{cod: 'angajat', titlu: 'angajat'}, {cod: 'total_ore', titlu: 'Total ore'}], pagina, '{currentPage}', {titlu: 'Titlu', antet_stanga: 'Antet stânga'}));
        const tabel = s.schemas[0].find((e: {name: string}) => e.name === 'tabel');
        expect(tabel[COLOANE_SQL]).toEqual(['angajat', 'total_ore']);
        expect(tabel.head).toEqual(['Angajat', 'Total ore']);
        // Antetul si titlul se repeta pe fiecare pagina (F09-T26).
        expect(s.basePdf.staticSchema.map((e: {name: string}) => e.name)).toContain('Titlu');
        expect(s.basePdf.staticSchema.map((e: {name: string}) => e.name)).toContain('Antet stânga');
        expect(s.schemas[0].map((e: {name: string}) => e.name)).toEqual(['tabel']);
        expect(s.basePdf.padding[0]).toBeGreaterThan(tabel.position.y - 3);
        expect(titluColoana({cod: 'angajat_id', titlu: ''})).toBe('Angajat_id');
    }, 30000);

    it('actualizeaza tabelul pastrand titlurile si latimile coloanelor ramase (UIR-05)', async () => {
        const initial = JSON.parse(await sablonDinColoane([{cod: 'angajat', titlu: 'angajat'}, {cod: 'marca', titlu: 'marca'}, {cod: 'ore', titlu: 'ore'}], pagina, '{currentPage}'));
        const tabel = initial.schemas[0].find((e: {name: string}) => e.name === 'tabel');
        tabel.head = ['Nume și prenume', 'Marcă', 'Ore'];
        tabel.headWidthPercentages = [50, 20, 30];
        tabel.position = {x: 33, y: 44};
        const sablon = JSON.stringify(initial);
        const noi = [{cod: 'angajat', titlu: 'angajat'}, {cod: 'ore', titlu: 'ore'}, {cod: 'garzi', titlu: 'garzi'}];
        const stare = stareTabel(sablon, noi);
        expect(stare).toEqual({exista: true, potrivit: false, adaugate: ['garzi'], eliminate: ['marca']});
        const dupa = JSON.parse(actualizeazaTabel(sablon, noi));
        const t2 = dupa.schemas[0].find((e: {name: string}) => e.name === 'tabel');
        expect(t2.head).toEqual(['Nume și prenume', 'Ore', 'Garzi']);
        expect(t2.headWidthPercentages.reduce((a: number, b: number) => a + b, 0)).toBeCloseTo(100);
        // Latimile pastrate raman in aceeasi proportie (50:30).
        expect(t2.headWidthPercentages[0] / t2.headWidthPercentages[1]).toBeCloseTo(50 / 30);
        expect(t2.position).toEqual({x: 33, y: 44});
        expect(t2[COLOANE_SQL]).toEqual(['angajat', 'ore', 'garzi']);
        expect(stareTabel(JSON.stringify(dupa), noi).potrivit).toBe(true);
        // Restul paginii nu se schimba.
        expect(dupa.basePdf).toEqual(initial.basePdf);
    }, 30000);

    it('recunoaste tabelele vechi fara coduri dupa titlu', async () => {
        const s = JSON.parse(await sablonDinColoane([{cod: 'angajat', titlu: 'angajat'}], pagina, '{currentPage}'));
        delete s.schemas[0].find((e: {name: string}) => e.name === 'tabel')[COLOANE_SQL];
        expect(stareTabel(JSON.stringify(s), [{cod: 'angajat', titlu: 'angajat'}]).potrivit).toBe(true);
        expect(stareTabel(JSON.stringify({basePdf: {}, schemas: [[]]}), [{cod: 'a', titlu: 'a'}]).exista).toBe(false);
    }, 30000);

    it('nu raporteaza nepotrivire la tabelele adaptive pontixTitluri (F09-T27)', async () => {
        // Ca la modelele standard si cele convertite: tabelul urmeaza singur SQL-ul la generare.
        const s = JSON.parse(await sablonDinColoane([], pagina, '{currentPage}'));
        const tabel = s.schemas[0].find((e: {name: string}) => e.name === 'tabel');
        tabel.pontixTitluri = true;
        const coloane = ['data', 'ziua', 'nume_si_prenume', 'observatii'].map(cod => ({cod, titlu: cod}));
        expect(stareTabel(JSON.stringify(s), coloane)).toEqual({exista: true, potrivit: true, adaugate: [], eliminate: []});
        delete tabel.pontixTitluri;
        expect(stareTabel(JSON.stringify(s), coloane).potrivit).toBe(false);
    }, 30000);

    it('are traduceri RO si EN pentru toate etichetele editorului (UIR-06)', () => {
        const ro = catalogRo as Record<string, string>, en = catalogEn as Record<string, string>;
        const lipsa = (cheiPdfme as string[]).filter(k => !ro['pdfme.' + k] || !en['pdfme.' + k]);
        expect(lipsa).toEqual([]);
        const etichete = etichetePdfme(k => ro[k]);
        expect(etichete.fieldsList).toBe('Elementele paginii');
        expect(Object.keys(etichete).length).toBe((cheiPdfme as string[]).length);
    });

    it('foloseste etichete traduse pentru tipurile de elemente (UIR-11)', async () => {
        const p = await pluginuriPdfme({text: 'Text cu câmpuri', table: 'Tabel'});
        expect(Object.keys(p)).toContain('Text cu câmpuri');
        expect(Object.keys(p)).toContain('Tabel');
        expect(Object.keys(p)).toContain('ellipse');
    }, 30000);

    it('potriveste zoomul la latimea zonei (UIR-09)', () => {
        // 1024x640: zona 958x407, A4 peisaj -> pagina mai lata decat la 100%.
        expect(zoomPotrivit(958, 407, 297, 210, false)).toBe(1.5);
        // 1440x960 cu panou: zona 1374x776, A4 peisaj.
        const z = zoomPotrivit(1374, 776, 297, 210, true);
        expect(z).toBeGreaterThan(0.75);
        expect(z).toBeLessThanOrEqual(1);
        // Pagina incape: latimea la zoom nu depaseste zona disponibila.
        expect((776 - 40) * (297 / 210) * z).toBeLessThanOrEqual(1374 - 420 - 90);
        expect(zoomPotrivit(300, 900, 210, 297, true)).toBe(0.5);
    });
});
