// Verificarea contrastului în grila de pontaj, pe tema luminoasă și pe cea întunecată (F16-T12, UI-3),
// în Edge headless, cu o punte Wails simulată și date fictive care produc toate stările celulelor:
// weekend, sărbătoare, selecție, avertisment, absențe (aprobată, în așteptare, respinsă), înlocuitor,
// „—” și totaluri necalculate. Pentru fiecare celulă se citesc culorile calculate de browser, cu
// fondul primului părinte opac, și se raportează contrastul minim pe fiecare stare (prag WCAG AA 4,5:1).
// Folosire: npm run build în frontend, apoi node verificari/contrast-grila.mjs
import {createServer} from 'node:http';
import {readFileSync, existsSync, mkdtempSync, rmSync} from 'node:fs';
import {resolve, dirname, join, extname, sep} from 'node:path';
import {fileURLToPath} from 'node:url';
import {tmpdir} from 'node:os';
import {spawn, spawnSync} from 'node:child_process';
import assert from 'node:assert/strict';

const radacina = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const dist = join(radacina, 'frontend', 'dist');
const browser = ['C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe',
    'C:/Program Files/Microsoft/Edge/Application/msedge.exe'].find(existsSync);
assert(browser, 'Microsoft Edge nu este disponibil.');
assert(existsSync(join(dist, 'index.html')), 'Rulează întâi npm run build în frontend.');
const cataloage = Object.fromEntries(['ro', 'en'].map(l => [l, JSON.parse(readFileSync(join(radacina, `internal/i18n/locales/${l}.json`), 'utf8'))]));
const schimburi = JSON.parse(readFileSync(join(radacina, 'internal/datebaza/initiale/schimburi.json'), 'utf8'));

const server = createServer((cerere, raspuns) => {
    const cale = resolve(dist, '.' + new URL(cerere.url, 'http://localhost').pathname);
    if (!cale.startsWith(dist + sep) && cale !== dist) { raspuns.writeHead(403).end(); return; }
    try {
        const fisier = cale === dist ? join(dist, 'index.html') : cale;
        raspuns.setHeader('Content-Type', {'.html': 'text/html', '.js': 'application/javascript', '.css': 'text/css'}[extname(fisier)] ?? 'application/octet-stream');
        raspuns.end(readFileSync(fisier));
    } catch { raspuns.writeHead(404).end(); }
});
await new Promise(r => server.listen(0, '127.0.0.1', r));
const url = `http://127.0.0.1:${server.address().port}`;
const pauza = ms => new Promise(r => setTimeout(r, ms));
// Edge lasă procese copil după oprirea celui principal (136 rămase pe 2026-10-30), iar taskkill /T
// nu le găsește pe toate. Se opresc toate procesele Edge care folosesc profilul temporar al rulării:
// profilul e unic, deci alte ferestre Edge nu sunt atinse.
const opresteBrowserul = (p, profil) => {
    p.kill();
    if (process.platform !== 'win32') return;
    const filtru = profil.replace(/'/g, "''");
    spawnSync('powershell', ['-NoProfile', '-Command', `Get-CimInstance Win32_Process -Filter "Name='msedge.exe'" | Where-Object { $_.CommandLine -like '*${filtru}*' } | ForEach-Object { Stop-Process -Id $_.ProcessId -Force -ErrorAction SilentlyContinue }`], {stdio: 'ignore'});
};

// Puntea simulată: 6 angajați fictivi; tema se alege prin localStorage înainte de încărcare.
function punte(tema, cataloage, schimburi) {
    try { localStorage.setItem('pontix_tema', tema); } catch { /* fără stocare */ }
    const zile = Array.from({length: 31}, (_, i) => `2026-10-${String(i + 1).padStart(2, '0')}`);
    const angajati = Array.from({length: 6}, (_, i) => ({id: `a${i}`, nume: `Angajat${i}`, prenume: 'Test', marca: `M-${i}`, functie: '',
        sectie_id: 's', departament_id: '', categorie_program_id: '', tip_contract: 'INTREAGA', norma_partiala_minute: null,
        norma_minute: 480, concediu_zile_an: null, activ: true}));
    const coduri = ['Z7', 'S7', 'G24', 'N12', ''];
    const celule = angajati.flatMap((a, i) => zile.map((d, z) => ({angajat_id: a.id, data: d, valoare: i === 5 ? '' : coduri[(i + z) % coduri.length], revizie: '1'})))
        .filter(c => c.valoare);
    // Calcul: ore pentru angajații 0–3, „—” (necalculabil) pentru angajatul 4, peste 24 h într-o zi (avertisment).
    const ore = celule.filter(c => c.angajat_id !== 'a4').map(c => ({angajat_id: c.angajat_id, data: c.data, minute: c.angajat_id === 'a3' && c.data === '2026-10-07' ? 1500 : 420, calculabil: true, surse: []}))
        .concat(zile.map(d => ({angajat_id: 'a4', data: d, minute: 0, calculabil: false, surse: []})));
    const absente = [['a1', '2026-10-12', 'APROBATA'], ['a1', '2026-10-13', 'IN_ASTEPTARE'], ['a1', '2026-10-14', 'RESPINSA'], ['a2', '2026-10-20', 'APROBATA']]
        .map(([angajat_id, data, stare_aprobare]) => ({angajat_id, data, cod: 'CO', activ: true, revizie: '1', nr_inregistrare: '1', inlocuitor_id: angajat_id === 'a2' ? 'a0' : '', stare_aprobare, aprobat_de: '', document_id: '', document_nume: '', document_mime: ''}));
    const schita = (an, luna) => ({schita: {an, luna, stare: 'DESCHISA', celule, absente, calcul: {celule: ore, totaluri: []},
        configuratie: {atribuire: 'ZI_INCEPUT', reguli: [{coloana: 'T', sursa: 'NEDEFINITA'}], revizie: ''}}});
    const cont = {id: 'c', nume: 'test', nume_afisat: 'Test', rol: 'ADMINISTRATOR', activ: true, intrebari: [], blocat_pana_la: ''};
    const metode = {
        Limbi: async () => ['ro', 'en'],
        Catalog: async l => cataloage[l],
        StareSeif: async () => ({blocat: false}),
        StareRecuperare: async () => ({necesara: false}),
        StareCont: async () => ({necesita_cont: false, autentificat: true, cont, ultimul_nume: 'test', intrebari: []}),
        DateInitiale: async () => ({date: {angajati, sectii: [{id: 's', denumire: {ro: 'ATI'}}], departamente: [], categorii_program: [], norma_intreaga_minute: 480, demo_disponibil: false}}),
        Informatii: async () => ({versiune: 'contrast', commit: ''}),
        StareAudit: async () => ({rezultat: {stare: 'ok', tabele: [], rand_rupt: 0, verificat_la: '', durata_ms: 0, copie_veche: false}}),
        SarcinaAIActiva: async () => '',
        CitesteSchita: async (an, luna) => schita(an, luna),
        PrevizualizeazaSchita: async (an, luna) => schita(an, luna),
        CitesteCalendar: async an => ({calendar: {an, sarbatori: [{data: '2026-10-26', denumire: {ro: 'Sărbătoare de probă'}}], import_disponibil: false}, adaugate: 0}),
        ListaSchimburi: async () => ({schimburi, adaugate: 0}),
        CitesteStareLuna: async (an, luna) => ({stare: {an, luna, stare: 'DESCHISA', modificat_la: '', modificat_de: '', motiv: '', amprenta: '', amprenta_corespunde: false}}),
        UltimaModificareCelula: async () => ({}),
    };
    window.go = {main: {App: new Proxy(metode, {get: (t, k) => t[k] ?? (async () => ({}))})}};
    window.runtime = new Proxy({}, {get: () => () => () => {}});
}

// Rulează în pagină: contrastul fiecărei celule, grupat pe stare.
function masoara() {
    const rgb = s => { const m = s.match(/[\d.]+/g).map(Number); return {r: m[0], g: m[1], b: m[2], a: m.length > 3 ? m[3] : 1}; };
    const lum = c => { const f = x => { x /= 255; return x <= 0.03928 ? x / 12.92 : ((x + 0.055) / 1.055) ** 2.4; }; return 0.2126 * f(c.r) + 0.7152 * f(c.g) + 0.0722 * f(c.b); };
    const fond = el => { for (let e = el; e; e = e.parentElement) { const c = rgb(getComputedStyle(e).backgroundColor); if (c.a > 0.5) return c; } return {r: 255, g: 255, b: 255, a: 1}; };
    const hex = c => '#' + [c.r, c.g, c.b].map(x => Math.round(x).toString(16).padStart(2, '0')).join('');
    const grupuri = {};
    const adauga = (stare, el) => {
        const text = rgb(getComputedStyle(el).color), f = fond(el);
        const [a, b] = [lum(text), lum(f)].sort((x, y) => y - x);
        const r = (a + 0.05) / (b + 0.05);
        if (!grupuri[stare] || r < grupuri[stare].contrast) grupuri[stare] = {contrast: Math.round(r * 100) / 100, text: hex(text), fond: hex(f), exemplu: el.value ?? el.textContent.trim().slice(0, 12)};
    };
    for (const input of document.querySelectorAll('.grila-pontaj tbody input[data-zi]')) {
        if (!input.value) continue;
        const td = input.parentElement, cls = td.className + ' ' + input.className;
        const stare = [/celula-absenta-(\w+)/.exec(cls)?.[0] ?? (cls.includes('celula-absenta') ? 'celula-absenta' : ''),
            cls.includes('celula-avertisment') ? 'avertisment' : '', cls.includes('celula-selectata') ? 'selecție' : '',
            cls.includes('zi-sarbatoare') ? 'sărbătoare' : cls.includes('zi-weekend') ? 'weekend' : '', input.value === '—' ? '„—”' : ''].filter(Boolean).join(' + ') || 'normală';
        adauga('celulă ' + stare, input);
    }
    for (const td of document.querySelectorAll('.grila-pontaj tbody td:not(:has(input))')) if (td.textContent.trim()) adauga(td.className.includes('total-necalculat') ? 'total necalculat' : (td.className.includes('fix-marca') ? 'marcă' : 'total'), td);
    for (const th of document.querySelectorAll('.grila-pontaj thead th')) adauga('antet ' + (th.className.includes('zi-sarbatoare') ? 'sărbătoare' : th.className.includes('zi-weekend') ? 'weekend' : 'zi'), th);
    for (const el of document.querySelectorAll('.legenda-grila span, .indiciu-grila')) adauga('legendă / indiciu', el);
    return grupuri;
}

const rezultate = {};
for (const tema of ['light', 'dark']) {
    const profil = mkdtempSync(join(tmpdir(), 'pontix-contrast-'));
    const proces = spawn(browser, ['--headless', '--disable-gpu', '--no-first-run', '--no-default-browser-check',
        '--remote-debugging-port=0', `--user-data-dir=${profil}`, '--window-size=1440,960', 'about:blank'], {windowsHide: true, stdio: 'ignore'});
    let port;
    for (let i = 0; i < 100 && !port; i++) {
        await pauza(100);
        try { port = readFileSync(join(profil, 'DevToolsActivePort'), 'utf8').split('\n')[0]; } catch { /* încă pornește */ }
    }
    assert(port, 'Edge nu a pornit');
    const tinte = await (await fetch(`http://127.0.0.1:${port}/json`)).json();
    const ws = new WebSocket(tinte.find(t => t.type === 'page').webSocketDebuggerUrl);
    await new Promise(r => ws.addEventListener('open', r));
    let id = 0;
    const asteptari = new Map();
    const erori = [];
    ws.addEventListener('message', e => {
        const m = JSON.parse(e.data);
        if (m.id && asteptari.has(m.id)) { asteptari.get(m.id)(m.result ?? m); asteptari.delete(m.id); }
        if (m.method === 'Runtime.exceptionThrown') erori.push(m.params.exceptionDetails.exception?.description ?? 'excepție');
    });
    const cdp = (method, params = {}) => new Promise(r => { asteptari.set(++id, r); ws.send(JSON.stringify({id, method, params})); });
    const evalueaza = async expr => {
        const r = await cdp('Runtime.evaluate', {expression: expr, returnByValue: true, awaitPromise: true});
        if (r.exceptionDetails) throw new Error(r.exceptionDetails.exception?.description ?? r.exceptionDetails.text);
        return r.result?.value;
    };
    await cdp('Runtime.enable');
    await cdp('Page.enable');
    await cdp('Page.addScriptToEvaluateOnNewDocument', {source: `(${punte})(${JSON.stringify(tema)}, ${JSON.stringify(cataloage)}, ${JSON.stringify(schimburi)})`});
    await cdp('Page.navigate', {url});
    const asteapta = async (expr, mesaj, limita = 30000) => {
        const sfarsit = Date.now() + limita;
        while (Date.now() < sfarsit) { if (await evalueaza(expr)) return; await pauza(50); }
        throw new Error(mesaj + (erori.length ? ': ' + erori.join(' | ') : ''));
    };
    await asteapta(`Boolean([...document.querySelectorAll('button.navigare')].find(b => b.textContent.includes('Pontaj')))`, 'aplicația nu s-a încărcat');
    await evalueaza(`[...document.querySelectorAll('button.navigare')].find(b => b.textContent.includes('Pontaj')).click()`);
    await asteapta(`document.querySelectorAll('.grila-pontaj tbody tr').length === 6 && [...document.querySelectorAll('.grila-pontaj tbody input')].some(i => i.value && i.value !== '…')`, 'grila nu s-a încărcat');
    assert.equal(await evalueaza(`document.documentElement.getAttribute('data-theme')`), tema);
    // O selecție de 3×4 celule, ca starea „selecție” să apară peste celelalte.
    await evalueaza(`(async () => {
        const c = (r, z) => document.querySelector('.grila-pontaj input[data-rand="' + r + '"][data-zi="' + z + '"]');
        c(0, 10).focus(); c(0, 10).dispatchEvent(new MouseEvent('mousedown', {bubbles: true, buttons: 1}));
        c(2, 13).dispatchEvent(new MouseEvent('mouseover', {bubbles: true, buttons: 1})); c(2, 13).dispatchEvent(new MouseEvent('mouseenter', {buttons: 1}));
        await new Promise(r => setTimeout(r, 300));
    })()`);
    rezultate[tema] = await evalueaza(`(${masoara})()`);
    rezultate[tema].erori = erori.length;
    ws.close();
    opresteBrowserul(proces, profil);
    await pauza(300);
    try { rmSync(profil, {recursive: true, force: true}); } catch { /* Edge poate ține fișiere câteva clipe */ }
}
server.close();
let sub = 0;
for (const [tema, grupuri] of Object.entries(rezultate)) {
    console.log(`\n${tema === 'light' ? 'Tema luminoasă' : 'Tema întunecată'} (erori JS: ${grupuri.erori})`);
    delete grupuri.erori;
    for (const [stare, g] of Object.entries(grupuri).sort((a, b) => a[1].contrast - b[1].contrast)) {
        const ok = g.contrast >= 4.5;
        if (!ok) sub++;
        console.log(`  ${ok ? 'OK     ' : 'SUB 4,5'} ${String(g.contrast).padStart(5)}:1  ${g.text} pe ${g.fond}  ${stare} (${g.exemplu})`);
    }
}
console.log(sub ? `\n${sub} stări sub 4,5:1` : '\nToate stările au cel puțin 4,5:1');
process.exitCode = sub ? 1 : 0;
