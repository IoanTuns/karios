// Verificare UI cu Edge existent și punte Wails simulată. Backendul real este
// verificat separat prin testele Go; acest script nu deschide date operaționale.
// Folosire: npm run build în frontend, apoi node verificari/interfata.mjs.
import {createServer} from 'node:http';
import {readFileSync, existsSync, mkdtempSync, rmSync, mkdirSync, writeFileSync} from 'node:fs';
import {resolve, dirname, join, basename, extname, sep} from 'node:path';
import {fileURLToPath} from 'node:url';
import {tmpdir} from 'node:os';
import {spawn} from 'node:child_process';
import assert from 'node:assert/strict';

const radacina = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const dist = join(radacina, 'frontend', 'dist');
const browser = ['C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe',
    'C:/Program Files/Microsoft/Edge/Application/msedge.exe'].find(existsSync);
assert(browser, 'Microsoft Edge nu este disponibil.');
assert(existsSync(join(dist, 'index.html')), 'Rulează întâi npm run build în frontend.');
const demo = JSON.parse(readFileSync(join(radacina, 'internal/datebaza/demo/angajati.json'), 'utf8'));
const calendarSursa = JSON.parse(readFileSync(join(radacina, 'internal/datebaza/calendare/romania-2026.json'), 'utf8'));
const versiuneAplicatie = readFileSync(join(radacina, 'internal/versiune/versiune.go'), 'utf8').match(/const Aplicatie = "([^"]+)"/)[1];
const cataloage = Object.fromEntries(['ro', 'en'].map(limba => [limba,
    JSON.parse(readFileSync(join(radacina, `internal/i18n/locales/${limba}.json`), 'utf8'))]));
const profil = mkdtempSync(join(tmpdir(), 'pontix-ui-'));
const server = createServer((cerere, raspuns) => {
    const cale = resolve(dist, '.' + new URL(cerere.url, 'http://localhost').pathname);
    if (!cale.startsWith(dist + sep) && cale !== dist) { raspuns.writeHead(403).end(); return; }
    try {
        const fisier = cale === dist ? join(dist, 'index.html') : cale;
        const mime = {'.html': 'text/html', '.js': 'application/javascript', '.css': 'text/css'}[extname(fisier)];
        raspuns.setHeader('Content-Type', mime ?? 'application/octet-stream');
        raspuns.end(readFileSync(fisier));
    } catch { raspuns.writeHead(404).end(); }
});
await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
const url = `http://127.0.0.1:${server.address().port}`;
const proces = spawn(browser, ['--headless', '--disable-gpu', '--no-first-run', '--no-default-browser-check',
    '--remote-debugging-port=0', `--user-data-dir=${profil}`, '--window-size=1440,960', 'about:blank'],
    {windowsHide: true, stdio: 'ignore'});
let ws;
let numar = 0;
const asteptari = new Map();
const erori = [];
const pauza = ms => new Promise(resolve => setTimeout(resolve, ms));

async function asteapta(conditie, mesaj, limita = 15000) {
    const sfarsit = Date.now() + limita;
    while (Date.now() < sfarsit) {
        if (await conditie()) return;
        await pauza(100);
    }
    throw new Error(mesaj);
}

function cdp(method, params = {}) {
    return new Promise((resolve, reject) => {
        const id = ++numar;
        const timer = setTimeout(() => { asteptari.delete(id); reject(new Error(`CDP: ${method}`)); }, 15000);
        asteptari.set(id, {resolve: rezultat => { clearTimeout(timer); resolve(rezultat); }, reject});
        ws.send(JSON.stringify({id, method, params}));
    });
}

async function evalueaza(expression) {
    const r = await cdp('Runtime.evaluate', {expression, returnByValue: true, awaitPromise: true});
    if (r.exceptionDetails) throw new Error(r.exceptionDetails.text);
    return r.result?.value;
}

async function apasa(text) {
    assert(await evalueaza(`(() => { const b = [...document.querySelectorAll('button')].find(b => { const copie = b.cloneNode(true); copie.querySelectorAll('[aria-hidden="true"]').forEach(e => e.remove()); return copie.textContent.trim() === ${JSON.stringify(text)} && !b.disabled; }); if (!b) return false; b.click(); return true; })()`), `Buton lipsă: ${text}`);
}

async function completeaza(selector, valoare) {
    await evalueaza(`(() => { const e = document.querySelector(${JSON.stringify(selector)}); if (!e) throw new Error('Câmp lipsă');
        Object.getOwnPropertyDescriptor(e.tagName === 'SELECT' ? HTMLSelectElement.prototype : HTMLInputElement.prototype, 'value').set.call(e, ${JSON.stringify(valoare)});
        e.dispatchEvent(new Event(e.tagName === 'SELECT' ? 'change' : 'input', {bubbles:true})); })()`);
}

// Se injectează înaintea React; toate mutațiile sunt limitate la localStorage
// din profilul temporar al testului, cu datele fictive incluse în repo.
function simulare(demo, cataloage, limba, calendarSursa, versiuneAplicatie) {
    let stare = JSON.parse(localStorage.getItem('pontix-test') || 'null') || {angajati: [], sectii: [], norma_intreaga_minute: 420, demo_disponibil: true};
    stare.sarbatori ??= structuredClone(calendarSursa);
    stare.schimburi ??= [];
    const raspuns = () => {
        stare.angajati = stare.angajati.map(a => ({...a, norma_minute: a.tip_contract === 'PARTIALA' ? a.norma_partiala_minute : stare.norma_intreaga_minute}));
        stare.demo_disponibil = stare.angajati.length === 0;
        localStorage.setItem('pontix-test', JSON.stringify(stare));
        return {date: structuredClone(stare), eroare: null};
    };
    const calendar = an => ({calendar: {an, sarbatori: structuredClone(stare.sarbatori.filter(s => Number(s.data.slice(0, 4)) === an).sort((a, b) => a.data.localeCompare(b.data))), import_disponibil: an === 2026}, adaugate: 0, eroare: null});
    window.go = {main: {App: {
        Limbi: async () => [limba, ...Object.keys(cataloage).filter(l => l !== limba)],
        Catalog: async l => cataloage[l], Informatii: async () => ({versiune: versiuneAplicatie, commit: 'verificare'}),
        DateInitiale: async () => raspuns(),
        IncarcaDemo: async () => { stare.angajati = structuredClone(demo.angajati); stare.sectii = structuredClone(demo.sectii); return raspuns(); },
        SalveazaAngajat: async a => {
            if (stare.angajati.some(b => b.id !== a.id && b.marca === a.marca)) return {date: null, eroare: {cod: 'date.marca_duplicata', param: {}}};
            const index = stare.angajati.findIndex(b => b.id === a.id);
            if (index < 0) stare.angajati.push({...a, id: 'nou-' + Date.now()}); else stare.angajati[index] = a;
            return raspuns();
        },
        SchimbaActivAngajat: async (id, activ) => { stare.angajati.find(a => a.id === id).activ = activ; return raspuns(); },
        SalveazaNormaIntreaga: async minute => { stare.norma_intreaga_minute = minute; return raspuns(); },
        AdaugaSectie: async denumire => { stare.sectii.push({id: 'sectie-' + Date.now(), denumire}); return raspuns(); },
        CitesteCalendar: async an => calendar(an),
        SalveazaSarbatoare: async (initiala, s, an) => {
            if (stare.sarbatori.some(v => v.data === s.data && v.data !== initiala)) return {calendar: null, eroare: {cod: 'sarbatori.data_duplicata', param: {}}};
            const i = stare.sarbatori.findIndex(v => v.data === initiala);
            if (i < 0) stare.sarbatori.push(s); else stare.sarbatori[i] = s;
            raspuns(); return calendar(an);
        },
        CompleteazaSarbatori2026: async () => {
            const lipsa = calendarSursa.filter(s => !stare.sarbatori.some(v => v.data === s.data));
            stare.sarbatori.push(...structuredClone(lipsa)); raspuns(); return {...calendar(2026), adaugate: lipsa.length};
        },
        ListaSchimburi: async () => ({schimburi: structuredClone(stare.schimburi), eroare: null}),
        SalveazaSchimb: async s => {
            s.cod = s.cod.trim().toUpperCase();
            if (stare.schimburi.some(v => v.id !== s.id && v.cod === s.cod)) return {schimburi: [], eroare: {cod: 'schimburi.cod_duplicat', param: {}}};
            const i = stare.schimburi.findIndex(v => v.id === s.id);
            if (i < 0) stare.schimburi.push({...s, id: 'schimb-' + Date.now()}); else stare.schimburi[i] = s;
            raspuns(); return {schimburi: structuredClone(stare.schimburi), eroare: null};
        },
        SchimbaActivSchimb: async (id, activ) => { stare.schimburi.find(s => s.id === id).activ = activ; raspuns(); return {schimburi: structuredClone(stare.schimburi), eroare: null}; }
    }}};
}

try {
    await asteapta(() => existsSync(join(profil, 'DevToolsActivePort')), 'Edge nu a pornit.');
    const port = readFileSync(join(profil, 'DevToolsActivePort'), 'utf8').split('\n')[0];
    const pagini = await (await fetch(`http://127.0.0.1:${port}/json/list`)).json();
    ws = new WebSocket(pagini.find(p => p.type === 'page').webSocketDebuggerUrl);
    await new Promise((resolve, reject) => { ws.onopen = resolve; ws.onerror = reject; });
    ws.onmessage = eveniment => {
        const r = JSON.parse(eveniment.data);
        if (r.id && asteptari.has(r.id)) {
            const cerere = asteptari.get(r.id); asteptari.delete(r.id);
            if (r.error) cerere.reject(new Error(r.error.message)); else cerere.resolve(r.result);
        } else if (r.method === 'Runtime.exceptionThrown') erori.push(r.params.exceptionDetails.text);
        else if (r.method === 'Page.javascriptDialogOpening') cdp('Page.handleJavaScriptDialog', {accept: true}).catch(err => erori.push(err.message));
    };
    await cdp('Page.enable'); await cdp('Runtime.enable');
    await cdp('Emulation.setDeviceMetricsOverride', {width: 1440, height: 960, deviceScaleFactor: 1, mobile: false});
    const injectat = await cdp('Page.addScriptToEvaluateOnNewDocument', {source: `(${simulare})(${JSON.stringify(demo)},${JSON.stringify(cataloage)},'ro',${JSON.stringify(calendarSursa)},${JSON.stringify(versiuneAplicatie)})`});
    await cdp('Page.navigate', {url});
    await asteapta(() => evalueaza(`document.body.textContent.includes('Adaugă primii membri')`), 'Starea goală nu apare.');
    await apasa('Încarcă date demo');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog[open]'))`), 'Dialog demo lipsă.');
    await apasa('Încarcă exemplele');
    await asteapta(() => evalueaza(`document.querySelectorAll('tbody tr').length === 12`), 'Demo nu apare în tabel.');
    await completeaza('select[aria-label="Stare"]', 'inactivi');
    await asteapta(() => evalueaza(`document.querySelectorAll('tbody tr').length === 2`), 'Filtrul de stare nu funcționează.');
    await completeaza('select[aria-label="Stare"]', 'toti');
    await completeaza('select[aria-label="Secție"]', demo.sectii[0].id);
    await asteapta(() => evalueaza(`document.querySelectorAll('tbody tr').length === 4`), 'Filtrul de secție nu funcționează.');
    await completeaza('select[aria-label="Secție"]', '');
    mkdirSync(join(radacina, 'build/bin'), {recursive: true});
    const captura = await cdp('Page.captureScreenshot', {format: 'png'});
    writeFileSync(join(radacina, 'build/bin/interfata-demo.png'), Buffer.from(captura.data, 'base64'));
    await completeaza('.cautare input', 'DEMO-002');
    await asteapta(() => evalueaza(`document.querySelectorAll('tbody tr').length === 1`), 'Căutarea nu filtrează.');
    await apasa('Editează');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog[open]'))`), 'Formular lipsă.');
    await completeaza('dialog input[type="number"]', '4.5');
    await apasa('Salvează');
    await asteapta(() => evalueaza(`!document.querySelector('dialog') && document.querySelector('tbody').textContent.includes('4 h 30 min')`), 'Norma parțială nu se salvează.');
    await apasa('Dezactivează');
    await asteapta(() => evalueaza(`document.querySelector('tbody').textContent.includes('Inactiv')`), 'Dezactivarea nu se vede.');
    await apasa('Reactivează');
    await apasa('Setări');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('.card-setari'))`), 'Setări lipsă.');
    await completeaza('.card-setari input[type="number"]', '8'); await apasa('Salvează');
    await asteapta(() => evalueaza(`document.querySelector('[role="status"]')?.textContent.includes('salvate')`), 'Norma globală nu se salvează.');
    await completeaza('.card-setari:nth-child(2) input', 'Secție de test'); await apasa('Adaugă secție');
    await asteapta(() => evalueaza(`document.querySelectorAll('.lista-sectii li').length === 4`), 'Secția nu se adaugă.');
    await apasa('Angajați'); await completeaza('.cautare input', 'DEMO-001');
    await asteapta(() => evalueaza(`document.querySelector('tbody')?.textContent.includes('8 h')`), 'Norma întreagă nu se actualizează.');
    await completeaza('.cautare input', 'DEMO-002');
    await asteapta(() => evalueaza(`document.querySelector('tbody')?.textContent.includes('4 h 30 min')`), 'Norma globală a schimbat norma parțială.');
    await apasa('Adaugă angajat');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog[open]'))`), 'Formular nou lipsă.');
    const campuri = ['Exemplu', 'Test', 'DEMO-001', 'Operator test'];
    for (let i = 0; i < campuri.length; i++) await completeaza(`dialog .campuri label:nth-child(${i + 1}) input`, campuri[i]);
    await apasa('Salvează');
    await asteapta(() => evalueaza(`document.querySelector('dialog [role="alert"]')?.textContent.includes('deja folosit')`), 'Eroarea de marcă duplicată nu apare.');
    await completeaza('dialog .campuri label:nth-child(3) input', 'TEST-013'); await apasa('Salvează');
    await asteapta(() => evalueaza(`!document.querySelector('dialog')`), 'Angajatul nou nu se salvează.');
    await apasa('Adaugă angajat');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog[open]'))`), 'Formular lipsă la anulare.');
    await completeaza('dialog .campuri label:first-child input', 'Modificat');
    await evalueaza(`window.confirm = () => false`); await apasa('Anulează');
    assert(await evalueaza(`Boolean(document.querySelector('dialog[open]'))`), 'Modificări nesalvate abandonate fără acord.');
    await evalueaza(`window.confirm = () => true`); await apasa('Anulează');
    await cdp('Page.reload');
    await asteapta(() => evalueaza(`document.querySelectorAll('tbody tr').length === 13`), 'Reîncărcarea UI nu păstrează starea simulată.');
    // Verificare și la lățimea minimă a ferestrei desktop.
    await cdp('Emulation.setDeviceMetricsOverride', {width: 1024, height: 640, deviceScaleFactor: 1, mobile: false});
    assert(await evalueaza(`document.documentElement.scrollWidth <= 1024`), 'Fereastra are depășire orizontală în afara tabelului.');
    await apasa('Sărbători');
    await asteapta(() => evalueaza(`document.querySelectorAll('.tabel-sarbatori tbody tr').length === 16`), 'Calendarul 2026 nu apare.');
    const capturaCalendar = await cdp('Page.captureScreenshot', {format: 'png'});
    writeFileSync(join(radacina, 'build/bin/calendar-2026.png'), Buffer.from(capturaCalendar.data, 'base64'));
    await evalueaza(`document.querySelector('.tabel-sarbatori button').click()`);
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog[open]'))`), 'Editare sărbătoare lipsă.');
    await completeaza('dialog .campuri label:nth-child(2) input', 'An nou personalizat'); await apasa('Salvează');
    await asteapta(() => evalueaza(`!document.querySelector('dialog') && document.querySelector('.tabel-sarbatori').textContent.includes('An nou personalizat')`), 'Editarea sărbătorii nu se salvează.');
    await apasa('Completează lista legală 2026');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog[open]'))`), 'Confirmare import lipsă.');
    await apasa('Completează datele lipsă');
    await asteapta(() => evalueaza(`!document.querySelector('dialog') && document.querySelector('.calendar-sarbatori [role="status"]').textContent.includes('0 date')`), 'Import repetat duplicat.');
    assert(await evalueaza(`document.querySelector('.tabel-sarbatori').textContent.includes('An nou personalizat')`), 'Importul a suprascris denumirea editată.');
    await apasa('Adaugă sărbătoare');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog[open]'))`), 'Formular sărbătoare lipsă.');
    await completeaza('dialog input[type="date"]', '2026-01-01');
    await completeaza('dialog .campuri label:nth-child(2) input', 'Duplicat'); await apasa('Salvează');
    await asteapta(() => evalueaza(`document.querySelector('dialog [role="alert"]')?.textContent.includes('Există deja')`), 'Data duplicată nu este semnalată.');
    await completeaza('dialog input[type="date"]', '2026-02-01'); await apasa('Salvează');
    await asteapta(() => evalueaza(`!document.querySelector('dialog') && document.querySelectorAll('.tabel-sarbatori tbody tr').length === 17`), 'Sărbătoarea nouă nu apare.');
    await completeaza('.unelte-calendar input', '2027'); await apasa('Arată anul');
    await asteapta(() => evalueaza(`document.querySelectorAll('.tabel-sarbatori tbody tr').length === 0 && Boolean(document.querySelector('.fara-rezultate'))`), 'Filtrarea anului nu funcționează.');
    await apasa('Schimburi');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('.tabel-schimburi'))`), 'Nomenclator lipsă.');
    await apasa('Adaugă schimb');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog[open]'))`), 'Formular schimb lipsă.');
    await completeaza('dialog .campuri label:nth-child(1) input', 'n12');
    await completeaza('dialog .campuri label:nth-child(2) select', 'NOAPTE');
    await completeaza('dialog .campuri label:nth-child(3) input', 'Tură de test');
    await completeaza('dialog input[type="time"]', '19:00');
    await completeaza('dialog .campuri label:nth-child(6) input', '720');
    await completeaza('dialog .campuri label:nth-child(7) input', '30'); await apasa('Salvează');
    await asteapta(() => evalueaza(`!document.querySelector('dialog') && document.querySelector('.tabel-schimburi').textContent.includes('N12')`), 'Schimb nesalvat.');
    assert(await evalueaza(`document.querySelector('.tabel-schimburi').textContent.includes('19:00') && document.querySelector('.tabel-schimburi').textContent.includes('12 h')`), 'Intervalul schimbului diferă.');
    await apasa('Dezactivează');
    await asteapta(() => evalueaza(`document.querySelector('.tabel-schimburi').textContent.includes('Inactiv')`), 'Schimbul nu se dezactivează.');
    await apasa('Reactivează');
    await asteapta(() => evalueaza(`document.querySelector('.tabel-schimburi').textContent.includes('Activ')`), 'Schimbul nu se reactivează.');
    await apasa('Adaugă schimb');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog[open]'))`), 'Formular duplicat lipsă.');
    await completeaza('dialog .campuri label:nth-child(1) input', 'N12');
    await completeaza('dialog .campuri label:nth-child(3) input', 'Duplicat');
    await completeaza('dialog input[type="time"]', '20:00');
    await completeaza('dialog .campuri label:nth-child(6) input', '720'); await apasa('Salvează');
    await asteapta(() => evalueaza(`document.querySelector('dialog [role="alert"]')?.textContent.includes('acest cod')`), 'Codul duplicat nu este semnalat.');
    await apasa('Anulează');
    await cdp('Page.removeScriptToEvaluateOnNewDocument', {identifier: injectat.identifier});
    await cdp('Page.addScriptToEvaluateOnNewDocument', {source: `(${simulare})(${JSON.stringify(demo)},${JSON.stringify(cataloage)},'en',${JSON.stringify(calendarSursa)},${JSON.stringify(versiuneAplicatie)})`});
    await cdp('Page.reload');
    await asteapta(() => evalueaza(`document.querySelector('h1')?.textContent === 'Employees'`), 'Traducerea EN nu apare.');
    assert.deepEqual(erori, [], 'Erori JavaScript în interfață.');
    console.log('UI verificată: angajați/demo/setări, calendar 2026, schimburi N12, duplicate, activare, import, ani, 1024px și EN.');
    console.log('Captură: build/bin/interfata-demo.png. Puntea Wails a fost simulată; baza reală este verificată de testele Go.');
} finally {
    if (ws?.readyState === WebSocket.OPEN) {
        try { await cdp('Browser.close'); } catch { /* Browserul poate închide conexiunea înaintea răspunsului. */ }
        ws.close();
    }
    proces.kill();
    await new Promise(resolve => server.close(resolve));
    // Profil creat exclusiv de acest test; verificăm ținta înainte de ștergere.
    if (dirname(resolve(profil)) === resolve(tmpdir()) && basename(profil).startsWith('pontix-ui-')) {
        await pauza(500);
        try { rmSync(profil, {recursive: true, force: true, maxRetries: 3, retryDelay: 300}); } catch { /* Poate rămâne blocat temporar de Edge. */ }
    }
}
