// F09-T20/T24: parcurgerea editorului vizual pdfme cu punte Wails simulata.
// Fonturile sunt cele reale (Go din modulul x/image, Arial din Windows), ca
// editorul si generarea sa lucreze ca in aplicatie. Fara date reale.
import assert from 'node:assert/strict';
import {readFileSync, existsSync} from 'node:fs';
import {join} from 'node:path';
import {spawnSync} from 'node:child_process';

export function fonturiPdfmeTest() {
    const dir = spawnSync('go', ['list', '-m', '-f', '{{.Dir}}', 'golang.org/x/image'], {encoding: 'utf8'}).stdout.trim();
    const ttf = n => readFileSync(join(dir, 'font', 'gofont', 'ttfs', n)).toString('base64');
    const go = {Go: ttf('Go-Regular.ttf'), GoBold: ttf('Go-Bold.ttf'), GoItalic: ttf('Go-Italic.ttf'), GoBoldItalic: ttf('Go-Bold-Italic.ttf')};
    const windows = {};
    for (const [nume, fisier] of [['Arial', 'arial.ttf'], ['Arial Bold', 'arialbd.ttf']]) {
        const cale = join(process.env.WINDIR || 'C:/Windows', 'Fonts', fisier);
        if (existsSync(cale)) windows[nume] = readFileSync(cale).toString('base64');
    }
    return {go, windows};
}

// Se injecteaza dupa simularea de baza si completeaza puntea cu metodele pdfme.
export function simularePdfme(fonturi) {
    const App = window.go.main.App;
    const vechi = App.PrevizualizeazaModelRaport;
    const lista = ['Arial', 'Arial Bold', 'Calibri', 'Cambria', 'Segoe UI', 'Times New Roman', 'Verdana'];
    App.FonturiRapoarte = async () => fonturi.go;
    App.ListaFonturiSistem = async () => ({fonturi: lista.map(nume => ({nume, familie: nume.split(' ')[0], stil: 'Regular', subset: true}))});
    App.FontSistem = async nume => fonturi.windows[nume] ? {nume, date: fonturi.windows[nume], subset: true, eroare: null} : {nume, date: '', subset: false, eroare: {cod: 'fonturi.inexistent', param: {nume}}};
    App.PrevizualizeazaModelRaport = async (d, c) => {
        const r = await vechi(d, c);
        if (r.eroare || d.continut.flexibil.mod_aspect !== 'PDFME') return r;
        const doc = r.document;
        const randuri = Array.from({length: 45}, (_, i) => doc.coloane.map((col, j) => j === 1 ? `Angajat demo ${i + 1} Ștefănescu` : j === 2 ? `DEMO-${String(i + 1).padStart(3, '0')}` : String((i * 7 + j) % 160)));
        const intrari = {titlu: d.denumire.ro || 'Raport', perioada: `${c.an}-${String(c.luna).padStart(2, '0')}`, an: String(c.an), luna: String(c.luna), generat_la: '05.11.2026 10:30', antet_stanga: 'Spital demo', antet_centru: '', antet_dreapta: 'Secția demo', subsol_stanga: '', subsol_centru: '', subsol_dreapta: '', tabel: JSON.stringify(randuri)};
        window.__pdfmeCereri = (window.__pdfmeCereri || 0) + 1;
        return {...r, pdf: '', generare_pdfme: true, sablon: d.continut.flexibil.sablon, intrari, token: 'pdfme-' + window.__pdfmeCereri};
    };
    App.FinalizeazaRaportPdfme = async (token, pdf) => {
        window.__pdfmeFinal = {token, octeti: atob(pdf).length};
        return {document: null, pdf, token, nume_fisier: 'raport-pdfme.pdf', pagini: 0, eroare: null};
    };
}

export async function verificaRapoartePdfme({cdp, evalueaza, asteapta, apasa, completeaza, cataloage, limba, captura}) {
    const t = k => cataloage[limba][k];
    const masuri = {};
    const masoara = async nume => {
        masuri[nume] = await evalueaza(`(() => {
            const sc = document.querySelector('dialog[open] .corp-dialog, dialog[open]') || document.scrollingElement;
            const ed = document.querySelector('.editor-pdfme');
            const r = ed?.getBoundingClientRect();
            return {latime: innerWidth, inaltime: innerHeight, derulareDialog: sc.scrollHeight, vizibilDialog: sc.clientHeight,
                derulareOrizontala: document.scrollingElement.scrollWidth > innerWidth || (sc.scrollWidth > sc.clientWidth + 2),
                editor: r ? {sus: Math.round(r.top + sc.scrollTop), latime: Math.round(r.width), inaltime: Math.round(r.height)} : null,
                campuriInainteDeEditor: ed ? [...document.querySelectorAll('dialog[open] input, dialog[open] select, dialog[open] textarea, dialog[open] button')].filter(e => e.compareDocumentPosition(ed) & Node.DOCUMENT_POSITION_FOLLOWING).length : null};
        })()`);
    };
    const laEditor = async () => evalueaza("document.querySelector('.editor-pdfme')?.scrollIntoView({block:'start'})");

    await asteapta(() => evalueaza("Boolean(document.querySelector('.actiuni-antet'))"), 'Aplicatie lipsa');
    await apasa(t('ui.setari'));
    await apasa(t('ui.definire_rapoarte'));
    await asteapta(() => evalueaza("Boolean(document.querySelector('.tabel-modele-rapoarte'))"), 'Modele lipsa');
    await captura('pdfme-01-lista-modele.png');
    await apasa(t('rapoarte.model_nou'));
    await completeaza('input[name=denumire_ro]', 'Gărzi lunare (pdfme)');
    await completeaza('select[name=sursa]', 'SQL_LIBER');
    await asteapta(() => evalueaza("Boolean(document.querySelector('.editor-flexibil'))"), 'Editor flexibil lipsa');
    await captura('pdfme-02-model-nou.png');
    await completeaza('textarea[name=raport_sql]', 'SELECT angajat_id, angajat, marca, total_ore_lucrate, total_ore_lucratoare FROM lunar');
    await apasa(t('rapoarte.flex_verifica'));
    await asteapta(() => evalueaza("Boolean(document.querySelector('.flex-date tbody tr'))"), 'Previzualizare tabel lipsa');

    await completeaza('select[name=flex_aspect]', 'PDFME');
    await asteapta(() => evalueaza("Boolean(document.querySelector('.editor-pdfme canvas, .editor-pdfme .ant-layout, .editor-pdfme [class*=ant]'))"), 'Designer pdfme neafisat', 30000);
    await new Promise(r => setTimeout(r, 1500));
    const sablon = JSON.parse(await evalueaza("JSON.stringify(JSON.parse(localStorage.getItem('pontix-test')||'{}'))") || '{}');
    void sablon;
    await masoara('editor_1440');
    await captura('pdfme-03-dialog-sus.png');
    await laEditor();
    await new Promise(r => setTimeout(r, 500));
    await captura('pdfme-04-editor.png');

    await completeaza('input[name=font_windows]', 'Arial');
    await apasa(t('rapoarte.pdfme_adauga_font'));
    await asteapta(() => evalueaza(`document.body.textContent.includes(${JSON.stringify(t('rapoarte.pdfme_font_adaugat').replace('{{nume}}', 'Arial'))})`), 'Font Windows neadaugat');
    await laEditor();
    await captura('pdfme-05-font-adaugat.png');

    // Selectia unui element arata panoul de proprietati (lista fonturilor).
    const selectat = await evalueaza(`(() => { const e = [...document.querySelectorAll('.editor-pdfme [title], .editor-pdfme [data-schema-name], .editor-pdfme [id^="antet"], .editor-pdfme [id^="titlu"]')].find(x => (x.id||x.title||'').includes('titlu')); if (!e) return false; e.dispatchEvent(new MouseEvent('mousedown', {bubbles:true})); e.dispatchEvent(new MouseEvent('mouseup', {bubbles:true})); e.click(); return true; })()`);
    masuri.elementSelectat = selectat;
    await new Promise(r => setTimeout(r, 500));
    await laEditor();
    await captura('pdfme-06-element-selectat.png');

    await apasa(t('rapoarte.flex_verifica'));
    await asteapta(() => evalueaza('Boolean(window.__pdfmeFinal)'), 'PDF pdfme negenerat', 30000);
    masuri.pdfGenerat = await evalueaza('window.__pdfmeFinal');
    await asteapta(() => evalueaza("Boolean(document.querySelector('.editor-flexibil .cadru-raport'))"), 'Previzualizare PDF lipsa');
    await evalueaza("document.querySelector('.editor-flexibil .cadru-raport').scrollIntoView({block:'start'})");
    await captura('pdfme-07-previzualizare.png');

    await cdp('Emulation.setDeviceMetricsOverride', {width: 1024, height: 640, deviceScaleFactor: 1, mobile: false});
    await new Promise(r => setTimeout(r, 800));
    await masoara('editor_1024');
    await laEditor();
    await captura('pdfme-08-editor-1024.png');
    await evalueaza("(document.querySelector('dialog[open] .corp-dialog, dialog[open]')||document.scrollingElement).scrollTop=0");
    await captura('pdfme-09-dialog-1024-sus.png');
    await cdp('Emulation.setDeviceMetricsOverride', {width: 1440, height: 960, deviceScaleFactor: 1, mobile: false});
    await evalueaza(`document.querySelector('dialog[open] button[aria-label=${JSON.stringify(t('ui.mareste'))}]')?.click()`);
    await new Promise(r => setTimeout(r, 800));
    await masoara('editor_1440_maximizat');
    await laEditor();
    await captura('pdfme-10-editor-1440-maximizat.png');
    assert(masuri.pdfGenerat?.octeti > 1000, 'PDF prea mic');
    return masuri;
}
