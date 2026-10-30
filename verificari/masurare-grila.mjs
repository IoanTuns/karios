// Măsurarea grilei de pontaj cu 50 și 300 de angajați (F06-T10, PERF-1), în Edge headless,
// cu o punte Wails simulată și date fictive generate. Nu deschide date operaționale.
// Folosire: npm run build în frontend, apoi node verificari/masurare-grila.mjs [numere...]
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
const numere = process.argv.slice(2).map(Number).filter(n => n > 0);
if (!numere.length) numere.push(50, 300);

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

// Puntea simulată: doar metodele folosite de autentificare, aplicație și pontaj; restul întorc {}.
function punte(numar, cataloage, schimburi) {
    const zile = Array.from({length: 31}, (_, i) => `2026-10-${String(i + 1).padStart(2, '0')}`);
    const coduri = ['Z7', 'S7', 'N12', 'G24', ''];
    const angajati = Array.from({length: numar}, (_, i) => ({id: `a${i}`, nume: `Angajat${String(i).padStart(3, '0')}`, prenume: 'Test',
        marca: `M-${i}`, functie: '', sectie_id: 's', departament_id: '', categorie_program_id: '', tip_contract: 'INTREAGA',
        norma_partiala_minute: null, norma_minute: 480, concediu_zile_an: null, activ: true}));
    const celule = angajati.flatMap((a, i) => zile.map((d, z) => ({angajat_id: a.id, data: d, valoare: coduri[(i + z) % coduri.length], revizie: '1'})))
        .filter(c => c.valoare);
    const cont = {id: 'c', nume: 'test', nume_afisat: 'Test', rol: 'ADMINISTRATOR', activ: true, intrebari: [], blocat_pana_la: ''};
    const metode = {
        Limbi: async () => ['ro', 'en'],
        Catalog: async l => cataloage[l],
        StareCont: async () => ({necesita_cont: false, autentificat: true, cont, ultimul_nume: 'test', intrebari: []}),
        DateInitiale: async () => ({date: {angajati, sectii: [{id: 's', denumire: {ro: 'ATI'}}], departamente: [], categorii_program: [], norma_intreaga_minute: 480, demo_disponibil: false}}),
        Informatii: async () => ({versiune: 'masurare', commit: ''}),
        StareAudit: async () => ({rezultat: {stare: 'ok', tabele: [], rand_rupt: 0, verificat_la: '', durata_ms: 0, copie_veche: false}}),
        SarcinaAIActiva: async () => '',
        CitesteSchita: async (an, luna) => ({schita: {an, luna, stare: 'DESCHISA', celule, absente: [], calcul: {celule: [], totaluri: []}, configuratie: {atribuire: 'ZI_INCEPUT', reguli: [], revizie: ''}}}),
        PrevizualizeazaSchita: async (an, luna) => ({schita: {an, luna, stare: 'DESCHISA', celule, absente: [], calcul: {celule: [], totaluri: []}, configuratie: {atribuire: 'ZI_INCEPUT', reguli: [], revizie: ''}}}),
        CitesteCalendar: async an => ({calendar: {an, sarbatori: [], import_disponibil: false}, adaugate: 0}),
        ListaSchimburi: async () => ({schimburi, adaugate: 0}),
        CitesteStareLuna: async (an, luna) => ({stare: {an, luna, stare: 'DESCHISA', modificat_la: '', modificat_de: '', motiv: '', amprenta: '', amprenta_corespunde: false}}),
        UltimaModificareCelula: async () => ({}),
    };
    window.go = {main: {App: new Proxy(metode, {get: (t, k) => t[k] ?? (async () => ({}))})}};
    window.runtime = new Proxy({}, {get: () => () => () => {}});
}

const rezultate = [];
for (const numar of numere) {
    const profil = mkdtempSync(join(tmpdir(), 'pontix-masurare-'));
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
    await cdp('Page.addScriptToEvaluateOnNewDocument', {source: `(${punte})(${numar}, ${JSON.stringify(cataloage)}, ${JSON.stringify(schimburi)})`});
    await cdp('Page.navigate', {url});
    const asteapta = async (expr, mesaj, limita = 30000) => {
        const sfarsit = Date.now() + limita;
        while (Date.now() < sfarsit) { if (await evalueaza(expr)) return; await pauza(20); }
        throw new Error(mesaj + (erori.length ? ': ' + erori.join(' | ') : ''));
    };
    await asteapta(`Boolean([...document.querySelectorAll('button.navigare')].find(b => b.textContent.includes('Pontaj')))`, 'aplicația nu s-a încărcat');
    const randare = await evalueaza(`(async () => {
        const t0 = performance.now();
        [...document.querySelectorAll('button.navigare')].find(b => b.textContent.includes('Pontaj')).click();
        while (document.querySelectorAll('.grila-pontaj tbody tr').length < ${numar}) await new Promise(r => setTimeout(r, 5));
        await new Promise(r => requestAnimationFrame(() => setTimeout(r)));
        return performance.now() - t0;
    })()`);
    // Fiecare pas: un eveniment, apoi lucrul React terminat și aspectul calculat. Se raportează mediana.
    const pasi = async cod => evalueaza(`(async () => {
        // requestAnimationFrame e încetinit în Edge headless; se așteaptă golirea cozii React
        // (planificatorul folosește MessageChannel) și se forțează calculul aspectului.
        const tic = () => new Promise(r => { const c = new MessageChannel(); c.port1.onmessage = () => r(); c.port2.postMessage(0); });
        const cadru = async () => { for (let i = 0; i < 3; i++) await tic(); document.body.getBoundingClientRect(); };
        const timpi = [];
        ${cod}
        timpi.sort((a, b) => a - b);
        return {mediana: timpi[Math.floor(timpi.length / 2)], maxim: timpi[timpi.length - 1]};
    })()`);
    const tragere = await pasi(`
        const rand = [...document.querySelectorAll('.grila-pontaj tbody tr')][0].querySelectorAll('input[data-zi]');
        rand[0].dispatchEvent(new MouseEvent('mousedown', {bubbles: true, buttons: 1})); rand[0].focus(); await cadru();
        for (let z = 1; z < 25; z++) { const t0 = performance.now(); rand[z].dispatchEvent(new MouseEvent('mouseenter', {buttons: 1})); rand[z].dispatchEvent(new MouseEvent('mouseover', {bubbles: true, buttons: 1})); await cadru(); timpi.push(performance.now() - t0); }
        rand[24].dispatchEvent(new MouseEvent('mouseup', {bubbles: true}));`);
    const taste = await pasi(`
        const prima = document.querySelector('.grila-pontaj tbody tr input[data-zi="0"]'); prima.focus(); await cadru();
        for (let i = 0; i < 25; i++) { const t0 = performance.now(); document.activeElement.dispatchEvent(new KeyboardEvent('keydown', {key: 'ArrowRight', bubbles: true})); await cadru(); timpi.push(performance.now() - t0); }`);
    rezultate.push({angajati: numar, celule: numar * 31, randare_ms: Math.round(randare), tragere_ms: tragere, sageata_ms: taste, erori: erori.length});
    ws.close();
    opresteBrowserul(proces, profil);
    await pauza(300);
    try { rmSync(profil, {recursive: true, force: true}); } catch { /* Edge poate ține fișiere câteva clipe */ }
}
server.close();
const r1 = x => Math.round(x * 10) / 10;
console.log('angajați | celule | afișare grilă | tragere (mediană / maxim) | săgeată (mediană / maxim) | erori');
for (const r of rezultate) console.log(`${r.angajati} | ${r.celule} | ${r.randare_ms} ms | ${r1(r.tragere_ms.mediana)} / ${r1(r.tragere_ms.maxim)} ms | ${r1(r.sageata_ms.mediana)} / ${r1(r.sageata_ms.maxim)} ms | ${r.erori}`);
