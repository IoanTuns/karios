// F09-T20/T24/T25: parcurgerea editorului vizual pdfme cu punte Wails simulata.
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
// Ca in Go: coloanele ascunse lipsesc din document, iar tabelul sablonului
// trebuie sa aiba cate o coloana pentru fiecare coloana vizibila.
export function simularePdfme(fonturi) {
    const App = window.go.main.App;
    const vechi = App.PrevizualizeazaModelRaport;
    const lista = ['Arial', 'Arial Bold', 'Calibri', 'Cambria', 'Segoe UI', 'Times New Roman', 'Verdana'];
    App.FonturiRapoarte = async () => fonturi.go;
    App.ListaFonturiSistem = async () => ({fonturi: lista.map(nume => ({nume, familie: nume.split(' ')[0], stil: 'Regular', subset: true}))});
    App.FontSistem = async nume => fonturi.windows[nume] ? {nume, date: fonturi.windows[nume], subset: true, eroare: null} : {nume, date: '', subset: false, eroare: {cod: 'fonturi.inexistent', param: {nume}}};
    App.PrevizualizeazaModelRaport = async (d, c) => {
        const r = await vechi(d, c);
        if (r.eroare) return r;
        const f = d.continut.flexibil;
        // Coloanele din SELECT-ul simulat: cele din SQL, daca exista, altfel cele de baza.
        const sel = /select\s+(.+?)\s+from/i.exec(f.sql || '');
        const coduri = sel ? sel[1].split(',').map(x => x.trim()) : r.document.coloane.map(c => c.cod);
        const vizibile = coduri.filter(cod => !(f.ascunse || []).includes(cod));
        const doc = {...r.document, coloane: vizibile.map(cod => ({cod, titlu: cod.replaceAll('_', ' '), pondere: 10}))};
        const randuri = Array.from({length: 45}, (_, i) => vizibile.map((cod, j) => cod === 'angajat' ? `Dr. Ștefănescu Ana ${i + 1}` : cod === 'marca' ? `DEMO-${String(i + 1).padStart(3, '0')}` : String((i * 7 + j) % 160)));
        doc.randuri = randuri.slice(0, 10).map(valori => ({valori, culori: []}));
        // F09-T27: ca in Go dupa F09-T23, modul vechi „tabel” nu mai are motor.
        if (f.mod_aspect !== 'PDFME') return window.__coloaneSimulate ? {...r, document: doc} : {eroare: {cod: 'rapoarte.model_neconvertit'}};
        const sablon = JSON.parse(f.sablon);
        window.__ultimSablon = sablon;
        const tabel = sablon.schemas.flat().find(e => e.type === 'table' && e.name === 'tabel' && !e.readOnly);
        if (tabel && tabel.head.length !== vizibile.length) return {eroare: {cod: 'rapoarte.pdfme_coloane', param: {sablon: String(tabel.head.length), date: String(vizibile.length)}}};
        const intrari = {titlu: d.denumire.ro || 'Raport', perioada: `${c.an}-${String(c.luna).padStart(2, '0')}`, an: String(c.an), luna: String(c.luna), generat_la: '05.11.2026 10:30', antet_stanga: 'Spital demo', antet_centru: '', antet_dreapta: 'Secția demo', subsol_stanga: '', subsol_centru: '', subsol_dreapta: '', tabel: JSON.stringify(randuri)};
        // F09-T22: ca in Go, weekendul (aici fiecare al 6-lea si al 7-lea rand) e aldin, o celula e editata.
        const randuriAldine = Object.fromEntries(randuri.map((_, i) => [String(i), {aldin: true}]).filter((_, i) => i % 7 === 5 || i % 7 === 6));
        intrari.pontix_stiluri = JSON.stringify({randuri: randuriAldine, celule: {'0:1': {cursiv: true, aliniere: 'R'}}});
        window.__pdfmeCereri = (window.__pdfmeCereri || 0) + 1;
        return {...r, document: doc, pdf: '', generare_pdfme: true, sablon: f.sablon, intrari, token: 'pdfme-' + window.__pdfmeCereri};
    };
    // F09-T27: coloanele si schema fara PDF, ca ColoaneModelRaport din Go.
    App.ColoaneModelRaport = async (d, c) => {
        window.__coloaneSimulate = true;
        try {
            const r = await App.PrevizualizeazaModelRaport(d, c);
            return r.eroare ? r : {document: {coloane: r.document.coloane, schema: r.document.schema}, eroare: null};
        } finally { window.__coloaneSimulate = false; }
    };
    App.FinalizeazaRaportPdfme = async (token, pdf) => {
        window.__pdfmeFinal = {token, octeti: atob(pdf).length};
        window.__pdfmeB64 = pdf;
        return {document: window.__ultimDoc || null, pdf, token, nume_fisier: 'raport-pdfme.pdf', pagini: 2, eroare: null};
    };
}

export async function verificaRapoartePdfme({cdp, evalueaza, asteapta, apasa, completeaza, cataloage, limba, captura, salveazaPdf}) {
    const t = k => cataloage[limba][k];
    const masuri = {};
    const pauza = ms => new Promise(r => setTimeout(r, ms));
    const masoara = async nume => {
        masuri[nume] = await evalueaza(`(() => {
            const sp = document.querySelector('.dialog-pdfme[open]');
            const ed = document.querySelector('.dialog-pdfme .editor-pdfme')?.getBoundingClientRect();
            const bara = document.querySelector('.dialog-pdfme .bara-pdfme')?.getBoundingClientRect();
            const corp = document.querySelector('.dialog-pdfme .spatiu-pdfme');
            return {latime: innerWidth, inaltime: innerHeight, spatiuDeschis: Boolean(sp),
                dialog: sp ? {latime: Math.round(sp.getBoundingClientRect().width), inaltime: Math.round(sp.getBoundingClientRect().height)} : null,
                editor: ed ? {sus: Math.round(ed.top), latime: Math.round(ed.width), inaltime: Math.round(ed.height), vizibilIntreg: ed.bottom <= innerHeight + 1} : null,
                baraSus: bara ? Math.round(bara.top) : null,
                derulareSpatiu: corp ? corp.scrollHeight > corp.clientHeight + 2 : null,
                derulareOrizontala: document.scrollingElement.scrollWidth > innerWidth,
                texteEngleza: [...document.querySelectorAll('.dialog-pdfme *')].some(e => e.childElementCount === 0 && /^(Field List|Bulk update field names|Edit Field)$/.test(e.textContent.trim()))};
        })()`);
    };

    await asteapta(() => evalueaza("Boolean(document.querySelector('.actiuni-antet'))"), 'Aplicatie lipsa');
    await apasa(t('ui.setari'));
    await apasa(t('ui.definire_rapoarte'));
    await asteapta(() => evalueaza("Boolean(document.querySelector('.tabel-modele-rapoarte'))"), 'Modele lipsa');
    await apasa(t('rapoarte.model_nou'));
    // F09-T31: modelul nou porneste gol, nu ca o copie a primului model din lista.
    await asteapta(() => evalueaza("Boolean(document.querySelector('textarea[name=raport_sql]'))"), 'Formular model nou lipsa');
    masuri.modelNou = await evalueaza(`({sursa: document.querySelector('select[name=sursa]').value, sql: document.querySelector('textarea[name=raport_sql]').value, antet: [...document.querySelectorAll('.formular-model-raport fieldset')][0].querySelectorAll('textarea') ? [...[...document.querySelectorAll('.formular-model-raport fieldset')][0].querySelectorAll('textarea')].map(x => x.value).join('') : '', ascunse: document.querySelector('input[name=coloane_ascunse]').value})`);
    assert.deepEqual(masuri.modelNou, {sursa: 'SQL_LIBER', sql: '', antet: '', ascunse: ''}, 'Modelul nou nu porneste gol');
    await completeaza('input[name=denumire_ro]', 'Gărzi lunare (pdfme)');
    await completeaza('select[name=sursa]', 'SQL_LIBER');
    await asteapta(() => evalueaza("Boolean(document.querySelector('.editor-flexibil'))"), 'Editor flexibil lipsa');
    await completeaza('textarea[name=raport_sql]', 'SELECT angajat_id, angajat, marca, total_ore_lucrate FROM lunar');

    // UIR-03: alegerea aspectului vizual fara verificare face singura sablonul din SQL.
    await completeaza('select[name=flex_aspect]', 'PDFME');
    await asteapta(() => evalueaza(`[...document.querySelectorAll('button')].some(b => b.textContent.trim() === ${JSON.stringify(t('rapoarte.pdfme_deschide'))} && !b.disabled)`), 'Sablonul de pornire nu s-a creat');
    masuri.ascunseDupaAlegere = await evalueaza("document.querySelector('input[name=coloane_ascunse]').value");
    assert.equal(masuri.ascunseDupaAlegere, 'angajat_id', 'angajat_id nu este ascuns implicit (UIR-04)');
    // UIR-02: in modul vizual, stilurile si zilele speciale nu mai apar.
    masuri.formular = await evalueaza(`({stiluri: Boolean(document.querySelector('select[name=stil_tinta_model]')), boldWeekend: Boolean(document.querySelector('input[name=bold_weekend]')), structura: Boolean(document.querySelector('.structura-model'))})`);
    // F09-T22: boldul pe weekend/sarbatori are efect si in modul vizual, deci ramane.
    assert.deepEqual(masuri.formular, {stiluri: false, boldWeekend: true, structura: false}, 'Controale gresite in modul vizual');
    await evalueaza("document.querySelector('.rezumat-pdfme').scrollIntoView({block:'center'})");
    await captura('t25-01-rezumat-formular.png');

    // UIR-01: spatiul de lucru separat, marit.
    await apasa(t('rapoarte.pdfme_deschide'));
    await asteapta(() => evalueaza("Boolean(document.querySelector('.dialog-pdfme .editor-pdfme [class*=ant]'))"), 'Editorul vizual nu s-a deschis', 30000);
    await pauza(1500);
    await masoara('spatiu_1024');
    await captura('t25-02-spatiu-1024.png');

    // UIR-07: campurile se copiaza cu un clic.
    await evalueaza("document.querySelector('.campuri-pdfme').open=true");
    await pauza(200);
    await captura('t25-03-campuri.png');
    await evalueaza("document.querySelector('.lista-campuri-pdfme .eticheta-camp.sql')?.click()");
    await pauza(300);
    masuri.copiere = await evalueaza("document.body.textContent.includes('{angajat}')");
    await evalueaza("document.querySelector('.campuri-pdfme').open=false");

    // UIR-08: fontul se adauga langa editor.
    await evalueaza("document.querySelector('.fonturi-pdfme').open=true");
    await completeaza('input[name=font_windows]', 'Arial');
    await apasa(t('rapoarte.pdfme_adauga_font'));
    await asteapta(() => evalueaza(`document.querySelector('.fonturi-pdfme').textContent.includes('Arial')`), 'Fontul nu apare ca folosit');
    await captura('t25-04-font.png');
    await evalueaza("document.querySelector('.fonturi-pdfme').open=false");

    // UIR-10: verificarea arata PDF-ul in fila lui (butonul din spatiul de lucru).
    await evalueaza(`[...document.querySelectorAll('.dialog-pdfme .actiuni-pdfme button')].find(b => b.textContent.trim() === ${JSON.stringify(t('rapoarte.flex_verifica'))}).click()`);
    await asteapta(() => evalueaza('Boolean(window.__pdfmeFinal)'), 'PDF pdfme negenerat', 30000);
    await asteapta(() => evalueaza("Boolean(document.querySelector('.dialog-pdfme .cadru-pdfme'))"), 'Fila PDF nu s-a deschis');
    masuri.pdfGenerat = await evalueaza('window.__pdfmeFinal');
    await salveazaPdf?.(await evalueaza('window.__pdfmeB64'));
    await captura('t25-05-fila-pdf.png');
    await apasa(t('rapoarte.pdfme_gata'));
    await asteapta(() => evalueaza("!document.querySelector('.dialog-pdfme[open]')"), 'Spatiul nu s-a inchis');

    // UIR-05: o coloana noua in SQL; tabelul se actualizeaza pastrand restul.
    await completeaza('textarea[name=raport_sql]', 'SELECT angajat_id, angajat, marca, total_ore_lucrate, total_ore_lucratoare FROM lunar');
    await apasa(t('rapoarte.flex_verifica'));
    await asteapta(() => evalueaza(`document.body.textContent.includes(${JSON.stringify(t('rapoarte.pdfme_actualizeaza_tabel'))})`), 'Actualizarea tabelului nu este oferita');
    await evalueaza("document.querySelector('.rezumat-pdfme').scrollIntoView({block:'center'})");
    await captura('t25-06-tabel-nepotrivit.png');
    await apasa(t('rapoarte.pdfme_actualizeaza_tabel'));
    await asteapta(() => evalueaza("!document.querySelector('.rezumat-pdfme .mesaj.avertisment')"), 'Tabelul a ramas nepotrivit');
    masuri.dupaActualizare = await evalueaza(`(() => { const d = JSON.parse(localStorage.getItem('pontix-test')||'{}'); return true; })()`);
    await apasa(t('rapoarte.flex_verifica'));
    await asteapta(() => evalueaza('window.__pdfmeCereri >= 2'), 'Verificarea dupa actualizare a esuat', 30000);

    // UIR-09: la 1440 panoul de proprietati porneste deschis.
    await cdp('Emulation.setDeviceMetricsOverride', {width: 1440, height: 960, deviceScaleFactor: 1, mobile: false});
    await pauza(500);
    await apasa(t('rapoarte.pdfme_deschide'));
    await asteapta(() => evalueaza("Boolean(document.querySelector('.dialog-pdfme .editor-pdfme [class*=ant]'))"), 'Editorul vizual nu s-a redeschis', 30000);
    await pauza(1500);
    await masoara('spatiu_1440');
    await captura('t25-07-spatiu-1440.png');
    // F09-T29: antetul (pagina fixa) se editeaza pe pagina; verificarea il pune la loc.
    const butonSpatiu = text => `[...document.querySelectorAll('.dialog-pdfme .actiuni-pdfme button')].find(b => b.textContent.trim() === ${JSON.stringify(text)})`;
    await evalueaza(`${butonSpatiu(t('rapoarte.pdfme_editeaza_antet'))}.click()`);
    await asteapta(() => evalueaza("Boolean(document.querySelector('.dialog-pdfme .indiciu-antet-pdfme'))"), 'Modul antet nu a pornit');
    await pauza(1200);
    masuri.antet = await evalueaza(`({apasat: ${butonSpatiu(t('rapoarte.pdfme_antet_gata'))}?.getAttribute('aria-pressed'), filaPdf: [...document.querySelectorAll('.dialog-pdfme .file-pdfme button')][1].disabled})`);
    await captura('t29-01-antet-editabil.png');
    const cereriInainte = await evalueaza('window.__pdfmeCereri || 0');
    await evalueaza(`${butonSpatiu(t('rapoarte.flex_verifica'))}.click()`);
    await asteapta(() => evalueaza(`(window.__pdfmeCereri || 0) > ${cereriInainte}`), 'Verificarea din modul antet nu a pornit', 30000);
    masuri.antetDupa = await evalueaza(`(() => { const s = window.__ultimSablon; return {fixe: s.basePdf.staticSchema.length, pagina: s.schemas[0].map(e => e.name), marcate: JSON.stringify(s).includes('pontixStatic'), sus: s.basePdf.padding[0], indiciu: Boolean(document.querySelector('.dialog-pdfme .indiciu-antet-pdfme'))}; })()`);
    assert.deepEqual(masuri.antet, {apasat: 'true', filaPdf: true}, 'Modul antet gresit');
    assert.deepEqual(masuri.antetDupa, {fixe: 7, pagina: ['tabel'], marcate: false, sus: 32, indiciu: false}, 'Antetul nu s-a pus la loc');
    for (const k of ['spatiu_1024', 'spatiu_1440']) {
        assert(masuri[k].spatiuDeschis && masuri[k].editor?.vizibilIntreg, `Editorul nu incape in fereastra (${k})`);
        assert(!masuri[k].derulareOrizontala, `Derulare orizontala (${k})`);
        assert(!masuri[k].texteEngleza, `Texte in engleza in editor (${k})`);
    }
    assert(masuri.pdfGenerat?.octeti > 1000, 'PDF prea mic');
    return masuri;
}
