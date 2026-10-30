// Verificarea contrastului pe toate ecranele, pe tema luminoasă și pe cea întunecată (F16-T16, UI-2),
// în Edge headless, cu o punte Wails simulată și date fictive. Pentru fiecare ecran din meniu, fiecare
// secțiune din Setări și ecranele de deblocare și de cheie de recuperare se citesc culorile calculate
// ale fiecărui text vizibil, cu fondul primului părinte opac. Se raportează textele sub 4,5:1 (WCAG AA;
// textele mari, de cel puțin 18,66 px aldin sau 24 px, au pragul 3:1) și, pe tema întunecată,
// suprafețele mari rămase luminoase. Câmpurile dezactivate sunt exceptate, ca în WCAG.
// Folosire: npm run build în frontend, apoi node verificari/contrast-ecrane.mjs [--detalii]
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
const detalii = process.argv.includes('--detalii');
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
// Edge lasă procese copil după oprirea celui principal; se opresc toate procesele cu profilul rulării.
const opresteBrowserul = (p, profil) => {
    p.kill();
    if (process.platform !== 'win32') return;
    const filtru = profil.replace(/'/g, "''");
    spawnSync('powershell', ['-NoProfile', '-Command', `Get-CimInstance Win32_Process -Filter "Name='msedge.exe'" | Where-Object { $_.CommandLine -like '*${filtru}*' } | ForEach-Object { Stop-Process -Id $_.ProcessId -Force -ErrorAction SilentlyContinue }`], {stdio: 'ignore'});
};

// Puntea simulată. `poarta` alege ecranul de pornire: aplicația, deblocarea sau cheia de recuperare.
function punte(tema, poarta, cataloage, schimburi) {
    try { localStorage.setItem('pontix_tema', tema); } catch { /* fără stocare */ }
    const zile = Array.from({length: 31}, (_, i) => `2026-10-${String(i + 1).padStart(2, '0')}`);
    const angajati = Array.from({length: 5}, (_, i) => ({id: `a${i}`, nume: `Angajat${i}`, prenume: 'Test', marca: `M-${i}`, functie: 'Asistent',
        sectie_id: 's', departament_id: '', categorie_program_id: '', tip_contract: 'INTREAGA', norma_partiala_minute: null,
        norma_minute: 480, concediu_zile_an: null, activ: i !== 4}));
    const celule = angajati.flatMap((a, i) => zile.map((d, z) => ({angajat_id: a.id, data: d, valoare: ['Z7', 'S7', 'G24', ''][(i + z) % 4], revizie: '1'}))).filter(c => c.valoare);
    const absente = [['a1', '2026-10-12', 'APROBATA'], ['a1', '2026-10-13', 'IN_ASTEPTARE'], ['a2', '2026-10-14', 'RESPINSA']]
        .map(([angajat_id, data, stare_aprobare]) => ({angajat_id, data, cod: 'CO', activ: true, revizie: '1', nr_inregistrare: '12', inlocuitor_id: '', stare_aprobare, aprobat_de: '', document_id: '', document_nume: '', document_mime: ''}));
    const schita = (an, luna) => ({schita: {an, luna, stare: 'DESCHISA', celule, absente, calcul: {celule: [], totaluri: []}, configuratie: {atribuire: 'ZI_INCEPUT', reguli: [], revizie: ''}}});
    const cont = {id: 'c', nume: 'admin', nume_afisat: 'Admin', rol: 'ADMINISTRATOR', activ: true, intrebari: ['intrebare_1', 'intrebare_2', 'intrebare_3'], blocat_pana_la: ''};
    const metode = {
        Limbi: async () => ['ro', 'en'],
        Catalog: async l => cataloage[l],
        StareSeif: async () => ({blocat: poarta === 'deblocare'}),
        StareRecuperare: async () => poarta === 'recuperare' ? {necesara: true, cod: 'ABCDE-FGHJK-MNPQR-STUVW-XYZ23'} : {necesara: false},
        StareCont: async () => ({necesita_cont: false, autentificat: true, cont, ultimul_nume: 'admin', intrebari: ['intrebare_1', 'intrebare_2', 'intrebare_3']}),
        DateInitiale: async () => ({date: {angajati, sectii: [{id: 's', denumire: {ro: 'ATI'}}], departamente: [], categorii_program: [], norma_intreaga_minute: 480, demo_disponibil: false}}),
        Informatii: async () => ({versiune: 'contrast', commit: ''}),
        StareAudit: async () => ({rezultat: {stare: 'ok', tabele: [], rand_rupt: 0, verificat_la: '2026-10-30T08:00:00Z', durata_ms: 0, copie_veche: false}}),
        SarcinaAIActiva: async () => '',
        CitesteSchita: async (an, luna) => schita(an, luna),
        PrevizualizeazaSchita: async (an, luna) => schita(an, luna),
        CitesteCalendar: async an => ({calendar: {an, sarbatori: [{data: '2026-10-26', denumire: {ro: 'Sărbătoare de probă'}}], import_disponibil: false}, adaugate: 0}),
        ListaSchimburi: async () => ({schimburi, adaugate: 0}),
        CitesteStareLuna: async (an, luna) => ({stare: {an, luna, stare: 'DESCHISA', modificat_la: '', modificat_de: '', motiv: '', amprenta: '', amprenta_corespunde: false}}),
        CitesteAbsente: async an => ({absente: {an, zile: absente, cereri: [], tipuri: []}}),
        ListaConturi: async () => ({conturi: [cont]}),
        ListaBackupuri: async () => ({backupuri: [{nume: 'pontix-2026-10-29-080000-v0.88.11.db', dimensiune: 1048576, modificat: '2026-10-29T08:00:00Z', manual: false}]}),
        UltimaModificareCelula: async () => ({}),
        PregatesteRaportProblema: async c => ({titlu: '[Problemă] ' + c.descriere, text: `## Problemă\n\n${c.descriere}\n\n### Mediu\n\n- Versiune: contrast`}),
    };
    window.go = {main: {App: new Proxy(metode, {get: (t, k) => t[k] ?? (async () => ({}))})}};
    window.runtime = new Proxy({}, {get: () => () => () => {}});
}

// Rulează în pagină: textele sub prag și, pe tema întunecată, suprafețele luminoase.
function masoara(intunecata) {
    const rgb = s => { const m = (s.match(/[\d.]+/g) || [0, 0, 0, 0]).map(Number); return {r: m[0], g: m[1], b: m[2], a: m.length > 3 ? m[3] : 1}; };
    const lum = c => { const f = x => { x /= 255; return x <= 0.03928 ? x / 12.92 : ((x + 0.055) / 1.055) ** 2.4; }; return 0.2126 * f(c.r) + 0.7152 * f(c.g) + 0.0722 * f(c.b); };
    const fondPropriu = e => rgb(getComputedStyle(e).backgroundColor);
    const fond = el => { for (let e = el; e; e = e.parentElement) { const c = fondPropriu(e); if (c.a > 0.5) return c; } return intunecata ? {r: 18, g: 25, b: 23, a: 1} : {r: 255, g: 255, b: 255, a: 1}; };
    const hex = c => '#' + [c.r, c.g, c.b].map(x => Math.round(x).toString(16).padStart(2, '0')).join('');
    const nume = el => el.tagName.toLowerCase() + (el.classList.length ? '.' + [...el.classList].slice(0, 2).join('.') : '');
    const vizibil = el => { const r = el.getBoundingClientRect(); const s = getComputedStyle(el); return r.width > 0 && r.height > 0 && s.visibility !== 'hidden' && Number(s.opacity) > 0.1; };
    const probleme = {};
    const adauga = (cheie, info) => { if (!probleme[cheie] || info.valoare < probleme[cheie].valoare) probleme[cheie] = {...info, numar: (probleme[cheie]?.numar ?? 0) + 1}; else probleme[cheie].numar++; };
    for (const el of document.querySelectorAll('body *')) {
        if (!vizibil(el) || el.closest('[hidden], [aria-hidden="true"]')) continue;
        const text = [...el.childNodes].some(n => n.nodeType === 3 && n.textContent.trim()) || (el.matches('input, select, textarea') && el.value);
        if (text && !el.matches(':disabled') && !el.closest(':disabled, [aria-disabled="true"]')) {
            const s = getComputedStyle(el), c = rgb(s.color), f = fond(el);
            const [a, b] = [lum(c), lum(f)].sort((x, y) => y - x);
            const raport = (a + 0.05) / (b + 0.05), marime = parseFloat(s.fontSize), aldin = Number(s.fontWeight) >= 700;
            const prag = marime >= 24 || (aldin && marime >= 18.66) ? 3 : 4.5;
            if (raport < prag) adauga('text ' + nume(el), {valoare: Math.round(raport * 100) / 100, prag, text: hex(c), fond: hex(f), exemplu: (el.value || el.textContent).trim().slice(0, 30)});
        }
        if (intunecata) {
            const c = fondPropriu(el), r = el.getBoundingClientRect();
            if (c.a > 0.5 && lum(c) > 0.45 && r.width * r.height > 600 && !el.matches('input[type="checkbox"], input[type="radio"], .mostra-weekend, .mostra-sarbatoare, .mostra-selectie, .mostra-modificat, .mostra-avertisment, .mostra-inlocuire'))
                adauga('suprafață luminoasă ' + nume(el), {valoare: Math.round(lum(c) * 100) / 100, prag: 0.45, text: '', fond: hex(c), exemplu: el.textContent.trim().slice(0, 30)});
        }
    }
    const panou = document.querySelector('.panou-eroare-componenta');
    if (panou) probleme['ecran căzut (ErrorBoundary)'] = {valoare: 0, prag: 0, text: '', fond: '', exemplu: panou.textContent.trim().slice(0, 120), numar: 1};
    return probleme;
}

async function sesiune(tema, poarta, actiune) {
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
    const asteapta = async (expr, mesaj, limita = 30000) => {
        const sfarsit = Date.now() + limita;
        while (Date.now() < sfarsit) { if (await evalueaza(expr)) return; await pauza(50); }
        throw new Error(mesaj + (erori.length ? ': ' + erori.join(' | ') : ''));
    };
    try {
        await cdp('Runtime.enable');
        await cdp('Page.enable');
        await cdp('Page.addScriptToEvaluateOnNewDocument', {source: `(${punte})(${JSON.stringify(tema)}, ${JSON.stringify(poarta)}, ${JSON.stringify(cataloage)}, ${JSON.stringify(schimburi)})`});
        await cdp('Page.navigate', {url});
        return await actiune({evalueaza, asteapta, erori});
    } finally {
        ws.close();
        opresteBrowserul(proces, profil);
        await pauza(300);
        try { rmSync(profil, {recursive: true, force: true}); } catch { /* Edge poate ține fișiere câteva clipe */ }
    }
}

const rezultate = [];
const masuratoare = async (tema, ecran, evalueaza) => {
    await pauza(400);
    rezultate.push({tema, ecran, probleme: await evalueaza(`(${masoara})(${tema === 'dark'})`)});
};
for (const tema of ['light', 'dark']) {
    await sesiune(tema, 'aplicatie', async ({evalueaza, asteapta}) => {
        await asteapta(`document.querySelectorAll('button.navigare').length > 3`, 'aplicația nu s-a încărcat');
        const ecrane = await evalueaza(`[...document.querySelectorAll('button.navigare')].map(b => b.textContent.trim())`);
        for (const [i, ecran] of ecrane.entries()) {
            await evalueaza(`document.querySelectorAll('button.navigare')[${i}].click()`);
            await masuratoare(tema, ecran, evalueaza);
            const sub = await evalueaza(`[...document.querySelectorAll('button.subnavigare')].map(b => b.textContent.trim())`);
            for (const [j, s] of sub.entries()) {
                await evalueaza(`document.querySelectorAll('button.subnavigare')[${j}].click()`);
                await masuratoare(tema, `${ecran} › ${s}`, evalueaza);
            }
        }
        // Dialogul „Raportează o problemă” (F16-T15): formularul și raportul pregătit.
        await evalueaza(`document.querySelector('.buton-raportare').click()`);
        await asteapta(`Boolean(document.querySelector('.formular-raportare textarea'))`, 'dialogul de raportare nu s-a deschis');
        await masuratoare(tema, 'Dialog raportare › formular', evalueaza);
        await evalueaza(`(() => { const el = document.querySelector('.formular-raportare textarea');
            Object.getOwnPropertyDescriptor(HTMLTextAreaElement.prototype, 'value').set.call(el, 'Grila nu se salvează');
            el.dispatchEvent(new Event('input', {bubbles: true})); })()`);
        await evalueaza(`[...document.querySelectorAll('.subsol-dialog .buton.primar')].pop().click()`);
        await asteapta(`Boolean(document.querySelector('.formular-raportare .text-raport')) && document.querySelector('.text-raport').value.includes('Grila nu se salvează')`, 'raportul nu s-a pregătit');
        await masuratoare(tema, 'Dialog raportare › raport pregătit', evalueaza);
    });
    for (const poarta of ['deblocare', 'recuperare']) {
        await sesiune(tema, poarta, async ({evalueaza, asteapta}) => {
            await asteapta(`Boolean(document.querySelector('.card-cont'))`, `ecranul ${poarta} nu s-a încărcat`);
            await masuratoare(tema, poarta === 'deblocare' ? 'Deblochează datele' : 'Cheia de recuperare', evalueaza);
        });
    }
}
server.close();

let total = 0;
for (const tema of ['light', 'dark']) {
    console.log(`\n${tema === 'light' ? 'Tema luminoasă' : 'Tema întunecată'}`);
    for (const r of rezultate.filter(x => x.tema === tema)) {
        const p = Object.entries(r.probleme);
        total += p.length;
        console.log(`  ${p.length ? '✗' : '✓'} ${r.ecran}${p.length ? ` — ${p.length} probleme` : ''}`);
        if (detalii || p.length <= 6) for (const [cheie, x] of p.sort((a, b) => a[1].valoare - b[1].valoare))
            console.log(`      ${cheie}: ${x.text ? `${x.valoare}:1 (prag ${x.prag}) ${x.text} pe ${x.fond}` : `fond ${x.fond}`} ×${x.numar} „${x.exemplu}”`);
    }
}
console.log(total ? `\n${total} probleme` : '\nNicio problemă de contrast');
process.exitCode = total ? 1 : 0;
