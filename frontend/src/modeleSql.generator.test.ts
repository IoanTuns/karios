// Generator reproductibil pentru modelele predefinite SQL + pdfme (D090, F09-T21).
// Ruleaza doar la cerere:  GENEREAZA_MODELE=<cale.json> npx vitest run src/modeleSql.generator.test.ts
// Rezultatul se transforma in migrarea internal/bd/migrari/0033_modele_standard_sql.sql.
import {describe, expect, it, vi} from 'vitest';
import {getDefaultFont} from '@pdfme/common';

vi.mock('../wailsjs/go/main/App', () => {
    const roboto = Object.values(getDefaultFont())[0].data as Uint8Array;
    let bin = '';
    for (const o of roboto) bin += String.fromCharCode(o);
    const b64 = btoa(bin);
    return {FonturiRapoarte: async () => ({Go: b64, GoBold: b64, GoItalic: b64, GoBoldItalic: b64}), FontSistem: vi.fn(), FinalizeazaRaportPdfme: vi.fn()};
});
vi.mock('./Toast', () => ({notifica: () => {}}));
vi.mock('../wailsjs/go/models', () => ({main: {RaspunsRaport: class {}}}));

import {sablonDinColoane} from './raportPdfme';
import modele from '../../internal/rapoarte/testdata/modele-v6.json';
import pontaj from '../../internal/rapoarte/testdata/raport-pontaj-v17.json';
import detaliat from '../../internal/rapoarte/testdata/raport-pontaj-detaliat-v18.json';
import fisa from '../../internal/rapoarte/testdata/raport-fisa-v17.json';
import absente from '../../internal/rapoarte/testdata/raport-absente-v14.json';
import garzi from '../../internal/rapoarte/testdata/raport-garzi-v14.json';

type Doc = {document: {coloane: {cod: string; titlu: string; pondere: number}[]}};
const surse: [string, string, string, Doc][] = [
    ['raport-pontaj', 'pontaj_lunar', 'raport-pontaj-sql', pontaj as Doc],
    ['raport-pontaj-detaliat', 'pontaj_detaliat', 'raport-pontaj-detaliat-sql', detaliat as Doc],
    ['raport-fisa', 'fisa_individuala', 'raport-fisa-sql', fisa as Doc],
    ['raport-absente', 'absente_anuale', 'raport-absente-sql', absente as Doc],
    ['raport-garzi', 'grafic_garzi', 'raport-garzi-sql', garzi as Doc],
];
const codSQL = (cod: string) => (/^\d{4}-\d{2}-\d{2}$/.test(cod) ? 'z' + cod.slice(8) : cod);
const proces = (globalThis as {process?: {env?: Record<string, string | undefined>}}).process;
const tinta = proces?.env?.GENEREAZA_MODELE;

describe('modele predefinite SQL + pdfme', () => {
    it.skipIf(!tinta)('genereaza definitiile', async () => {
        const rez = [];
        for (const [vechi, tabel, id, doc] of surse) {
            const m = (modele as {modele: Record<string, unknown>[]}).modele.find(x => x.id === vechi) as Record<string, any>;
            expect(m).toBeTruthy();
            const coloane = doc.document.coloane.map(c => ({cod: codSQL(c.cod), titlu: c.titlu}));
            const sablon = await sablonDinColoane(coloane, m.pagina, 'Pagina {currentPage} din {totalPages}',
                {antet_stanga: 'Antet stânga', antet_centru: 'Antet centru', antet_dreapta: 'Antet dreapta', titlu: 'Titlu', subsol_stanga: 'Subsol stânga', subsol_centru: 'Subsol centru', pagina: 'Număr pagină'},
                {ponderi: doc.document.coloane.map(c => c.pondere || 1), titluriDinRaport: true, stilClasic: true});
            const fara = (b: Record<string, unknown>) => { const c = {...b}; delete c.imagini; return c; };
            rez.push({
                id, denumire: {ro: m.denumire.ro + ' (editor vizual)', en: m.denumire.en + ' (visual editor)'},
                antet: fara(m.antet), subsol: fara(m.subsol), pagina: m.pagina,
                continut: {sursa: 'SQL_LIBER', coloane: ['rezultat'], text: m.continut.text, ...(m.continut.ascunde_explicatii ? {ascunde_explicatii: true} : {}),
                    flexibil: {inaltime_rand: 18, ascunse: ['rand'], sql: `SELECT rand, ${coloane.map(c => c.cod).join(', ')} FROM ${tabel} ORDER BY rand`,
                        mod_aspect: 'PDFME', markdown: '', sablon, elemente: [], calculate: [], concedii_tipuri: [], concedii_categorii: [], weekend: true, sarbatori: true}},
            });
        }
        const fs = await import(/* @vite-ignore */ 'node:' + 'fs') as {writeFileSync: (p: string, d: string) => void};
        fs.writeFileSync(tinta!, JSON.stringify(rez, null, 1));
    }, 60000);
});
