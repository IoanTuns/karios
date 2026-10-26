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
const schimburiSursa = JSON.parse(readFileSync(join(radacina, 'internal/datebaza/initiale/schimburi.json'), 'utf8'));
const versiuneAplicatie = readFileSync(join(radacina, 'internal/versiune/versiune.go'), 'utf8').match(/const Aplicatie = "([^"]+)"/)[1];
const cataloage = Object.fromEntries(['ro', 'en'].map(limba => [limba,
    JSON.parse(readFileSync(join(radacina, `internal/i18n/locales/${limba}.json`), 'utf8'))]));
const profil = mkdtempSync(join(tmpdir(), 'pontix-ui-'));
const imaginiDemo = join(profil, 'rapoarte-demo');
mkdirSync(imaginiDemo, {recursive: true});
const pdfDemoBuffer = Buffer.from('%PDF-1.4\n1 0 obj<</Type/Catalog>>endobj\n%%EOF\n');
writeFileSync(join(imaginiDemo, 'raport-pontaj.pdf'), pdfDemoBuffer);
writeFileSync(join(imaginiDemo, 'cerere-demo.pdf'), pdfDemoBuffer);
writeFileSync(join(imaginiDemo, 'antet-demo.png'), Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR4nGP4z8DwHwAFAAH/iZk9HQAAAABJRU5ErkJggg==', 'base64'));
writeFileSync(join(imaginiDemo, 'subsol-demo.jpg'), Buffer.from('/9j/4AAQSkZJRgABAQEASABIAAD/2wBDAP//////////////////////////////////////////////////////////////////////////////////////wgALCAABAAEBAREA/8QAFBABAAAAAAAAAAAAAAAAAAAAAP/aAAgBAQABPxA=', 'base64'));
writeFileSync(join(imaginiDemo, 'imagine-invalida.txt'), 'exemplu de fisier care nu este imagine');
const rapoarteSursa = {reguli:JSON.parse(readFileSync(join(radacina,'internal/bd/migrari/0015_pontaj_ore_reguli.sql'),'utf8').match(/'(\{"atribuire"[^']+)\'/)[1]),configuratie:JSON.parse(readFileSync(join(radacina,'internal/rapoarte/testdata/modele-v3.json'),'utf8')),document:JSON.parse(readFileSync(join(radacina,'internal/rapoarte/testdata/raport-pontaj-v3.json'),'utf8')).document,pdf:pdfDemoBuffer.toString('base64')};
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
    if (r.exceptionDetails) throw new Error(r.exceptionDetails.exception?.description || r.exceptionDetails.text);
    return r.result?.value;
}

async function apasa(text) {
    assert(await evalueaza(`(() => { const b = [...document.querySelectorAll('button')].find(b => { const copie = b.cloneNode(true); copie.querySelectorAll('[aria-hidden="true"]').forEach(e => e.remove()); return copie.textContent.trim() === ${JSON.stringify(text)} && !b.disabled; }); if (!b) return false; b.click(); return true; })()`), `Buton lipsă: ${text}`);
}

async function incarcaFisier(selector, cale) {
    const {root}=await cdp('DOM.getDocument');
    const {nodeId}=await cdp('DOM.querySelector',{nodeId:root.nodeId,selector});
    assert(nodeId,`Campul de import lipseste: ${selector}`);
    await cdp('DOM.setFileInputFiles',{nodeId,files:[cale]});
}
async function completeaza(selector, valoare) {
    await evalueaza(`(() => { const e = document.querySelector(${JSON.stringify(selector)}); if (!e) throw new Error('Câmp lipsă');
        Object.getOwnPropertyDescriptor(e.tagName === 'SELECT' ? HTMLSelectElement.prototype : e.tagName === 'TEXTAREA' ? HTMLTextAreaElement.prototype : HTMLInputElement.prototype, 'value').set.call(e, ${JSON.stringify(valoare)});
        e.dispatchEvent(new Event(e.tagName === 'SELECT' ? 'change' : 'input', {bubbles:true})); })()`);
}

// Se injectează înaintea React; toate mutațiile sunt limitate la localStorage
// din profilul temporar al testului, cu datele fictive incluse în repo.
function simulare(demo, cataloage, limba, calendarSursa, versiuneAplicatie, schimburiSursa, rapoarteSursa) {
    let stare = JSON.parse(localStorage.getItem('pontix-test') || 'null') || {angajati: [], sectii: [], norma_intreaga_minute: 420, demo_disponibil: true};
    stare.sarbatori ??= structuredClone(calendarSursa);
    stare.schimburi ??= structuredClone(schimburiSursa);
    stare.ciorne ??= [];
    stare.departamente ??= [];
    stare.tipuriAbsenta ??= [{cod:'CO',denumire:{ro:'Concediu de odihnă',en:'Annual leave'},categorie:'ODIHNA',activ:true}, {cod:'CM',denumire:{ro:'Concediu medical',en:'Medical leave'},categorie:'MEDICAL',activ:true}];
    stare.absente ??= [];
    stare.modeleRapoarte ??= structuredClone(rapoarteSursa.configuratie.modele);
    stare.documenteAbsenta ??= [];
    stare.reguliPontaj ??= {...structuredClone(rapoarteSursa.reguli),revizie:'1'};
    // Date simulate pentru fluxul UI; regulile reale sunt verificate separat în Go.
    function schita(an,luna,modificari=[]) {
        const prefix=`${String(an).padStart(4,'0')}-${String(luna).padStart(2,'0')}-`;
        const raw=new Map(stare.ciorne.map(c=>[c.angajat_id+'/'+c.data,c]));
        for(const c of modificari)raw.set(c.angajat_id+'/'+c.data,c);
        const celule=new Map(),totaluri=new Map();
        for(const c of raw.values()) {
            if(!c.valoare.trim())continue;
            const p=stare.schimburi.find(s=>s.cod===c.valoare),manual=c.valoare.startsWith('@M')?Number(c.valoare.slice(2)):Number(c.valoare)*60;
            let reste=p?p.durata_minute:manual,curent=new Date(c.data+'T00:00:00Z');
            curent.setUTCMinutes(p?.ora_start>=0?p.ora_start:0);
            const cunoscut=Number.isFinite(reste)&&reste>=0&&!(p?.pauza_minute);
            if(!cunoscut)reste=0;
            do {
                const data=curent.toISOString().slice(0,10),credit=stare.reguliPontaj.atribuire==='ZI_INCEPUT'?c.data:data;
                const n=Math.min(reste,1440-curent.getUTCHours()*60-curent.getUTCMinutes());
                if(credit.startsWith(prefix)){
                    const k=c.angajat_id+'/'+credit,z=celule.get(k)||{angajat_id:c.angajat_id,data:credit,minute:0,calculabil:true,surse:[]};
                    z.minute+=n;z.calculabil&&=cunoscut;z.surse.push(c.data);celule.set(k,z);
                    const a=totaluri.get(c.angajat_id)||{total:0,norma:0,znl:0,complet:true};a.total+=n;a.complet&&=cunoscut;
                    if([0,6].includes(curent.getUTCDay())||stare.sarbatori.some(s=>s.data===data))a.znl+=n;else a.norma+=n;
                    totaluri.set(c.angajat_id,a);
                }
                reste-=n;curent=new Date(curent.getTime()+Math.max(n,1)*60000);
            }while(reste>0);
        }
        const absente=stare.absente.filter(z=>z.data.startsWith(prefix));
        for(const a of absente)if(a.activ)celule.set(a.angajat_id+'/'+a.data,{angajat_id:a.angajat_id,data:a.data,minute:0,calculabil:true,surse:[]});
        return {schita:{an,luna,configuratie:structuredClone(stare.reguliPontaj),celule:structuredClone(stare.ciorne.filter(c=>c.data.startsWith(prefix))),absente:structuredClone(absente),calcul:{celule:[...celule.values()],totaluri:[...totaluri].map(([angajat_id,a])=>({angajat_id,valori:Object.fromEntries(stare.reguliPontaj.reguli.map(r=>[r.coloana,{sutimi:({TOTAL_ORE:a.total,ORE_ZNL:a.znl,ORE_NORMA:a.norma}[r.sursa]||0)*r.factor_sutimi,unitate:r.sursa==='ZILE_ABSENTA'?'ZILE':'ORE',definita:r.sursa!=='NEDEFINITA'&&a.complet}]))}))}},eroare:null};
    }
    const configuratiePontaj=()=>({configuratie:structuredClone(stare.reguliPontaj),optiuni:[...stare.schimburi.map(s=>({...s,tip:'SCHIMB'})),...stare.tipuriAbsenta.map(s=>({...s,tip:'ABSENTA'}))],eroare:null});
    const raspuns = () => {
        stare.angajati = stare.angajati.map(a => ({...a, departament_id: a.departament_id || '', norma_minute: a.tip_contract === 'PARTIALA' ? a.norma_partiala_minute : stare.norma_intreaga_minute}));
        stare.demo_disponibil = stare.angajati.length === 0;
        localStorage.setItem('pontix-test', JSON.stringify(stare));
        return {date: structuredClone(stare), eroare: null};
    };
    const calendar = an => ({calendar: {an, sarbatori: structuredClone(stare.sarbatori.filter(s => Number(s.data.slice(0, 4)) === an).sort((a, b) => a.data.localeCompare(b.data))), import_disponibil: an === 2026}, adaugate: 0, eroare: null});
    const absente = an => ({absente:{an,tipuri:structuredClone(stare.tipuriAbsenta),zile:structuredClone(stare.absente.filter(z => Number(z.data.slice(0,4)) === an).sort((a,b) => a.data.localeCompare(b.data)))},eroare:null});
    const modele = () => ({configuratie:{...structuredClone(rapoarteSursa.configuratie),modele:structuredClone(stare.modeleRapoarte)},eroare:null});
    window.go = {main: {App: {
        Limbi: async () => [limba, ...Object.keys(cataloage).filter(l => l !== limba)],
        Catalog: async l => cataloage[l], Informatii: async () => ({versiune: versiuneAplicatie, commit: 'verificare'}),
        DateInitiale: async () => raspuns(),
        CreeazaBackup: async () => { stare.backupuriManuale = (stare.backupuriManuale || 0) + 1; raspuns(); return {cale: 'karios\\backup\\pontix-manual-2026-10-26-110000-v' + versiuneAplicatie + '.db', eroare: null}; },
        CitesteConfiguratiePontaj:async()=>configuratiePontaj(),
        SalveazaConfiguratiePontaj:async c=>{if(c.revizie!==stare.reguliPontaj.revizie)return {eroare:{cod:'pontaj.reguli_conflict'}};stare.reguliPontaj={...structuredClone(c),revizie:String(Number(c.revizie)+1)};raspuns();return configuratiePontaj();},
        PrevizualizeazaSchita:async(an,luna,celule)=>schita(an,luna,celule),
        ListaModeleRapoarte: async () => modele(),
        SalveazaModelRaport: async d => {
            const texte=[d.antet.stanga,d.antet.centru,d.antet.dreapta,d.subsol.stanga,d.subsol.centru,d.subsol.dreapta,d.continut.text].flatMap(v => Object.values(v));
            if (texte.some(v => [...v.matchAll(/\{\{([^{}]+)\}\}/g)].some(m => !rapoarteSursa.configuratie.variabile.includes(m[1])))) return {configuratie:null,eroare:{cod:'rapoarte.variabila_invalida'}};
            const i=stare.modeleRapoarte.findIndex(v => v.id === d.id);
            if(i>=0 && stare.modeleRapoarte[i].revizie!==d.revizie) return {configuratie:null,eroare:{cod:'rapoarte.conflict'}};
            if(i<0){d={...d,id:'model-'+crypto.randomUUID(),revizie:'1'};stare.modeleRapoarte.push(d);}else{d={...d,revizie:String(Number(d.revizie)+1)};stare.modeleRapoarte[i]=d;}
            raspuns();return {...modele(),salvat_id:d.id};
        },
        PregatesteRaport: async c => {
            stare.ultimaCerereRaport=c;stare.ultimTokenRaport=crypto.randomUUID();raspuns();
            const d=stare.modeleRapoarte.find(d => d.id===c.model_id);
            return {document:{...structuredClone(rapoarteSursa.document),titlu:d.denumire[c.limba]||d.denumire.ro,sursa:d.continut.sursa,revizie_model:d.revizie},pdf:rapoarteSursa.pdf,token:stare.ultimTokenRaport,nume_fisier:'raport-demo.pdf',pagini:4,eroare:null};
        },
        SalveazaPDFRaport: async token => { if(token!==stare.ultimTokenRaport) return {eroare:{cod:'rapoarte.expirat'}};stare.exporturiRaport=(stare.exporturiRaport||0)+1;raspuns();return {anulat:false,cale:'raport-demo.pdf',eroare:null}; },
        DeschidePDFPentruTiparire: async token => { if(token!==stare.ultimTokenRaport) return {eroare:{cod:'rapoarte.expirat'}};stare.tipaririRaport=(stare.tipaririRaport||0)+1;raspuns();return {anulat:false,cale:'raport-demo.pdf',eroare:null}; },
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
        AdaugaDepartament: async (sectie_id, denumire) => { stare.departamente.push({id: crypto.randomUUID(), sectie_id, denumire}); return raspuns(); },
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
        CompleteazaSchimburiInitiale: async categorie => {
            const cat = (categorie || 'LAB').toUpperCase();
            if (cat !== 'LAB') return {schimburi: [], adaugate: 0, eroare: {cod: 'schimburi.categorie_invalida', param: {}}};
            const lipsa = schimburiSursa.filter(s => (s.categorie || 'LAB') === cat && !stare.schimburi.some(v => v.cod.toUpperCase() === s.cod));
            stare.schimburi.push(...structuredClone(lipsa)); raspuns(); return {schimburi: structuredClone(stare.schimburi), adaugate: lipsa.length, eroare: null};
        },
        SalveazaSchimb: async s => {
            s.cod = s.cod.trim().toUpperCase();
            if (stare.schimburi.some(v => v.id !== s.id && v.cod === s.cod)) return {schimburi: [], eroare: {cod: 'schimburi.cod_duplicat', param: {}}};
            const i = stare.schimburi.findIndex(v => v.id === s.id);
            if (i < 0) stare.schimburi.push({...s, id: 'schimb-' + Date.now()}); else stare.schimburi[i] = s;
            raspuns(); return {schimburi: structuredClone(stare.schimburi), eroare: null};
        },
        SchimbaActivSchimb: async (id, activ) => { stare.schimburi.find(s => s.id === id).activ = activ; raspuns(); return {schimburi: structuredClone(stare.schimburi), eroare: null}; },
        CitesteAbsente: async an => absente(an),
        SalveazaTipAbsenta: async (initial,tip,an) => {
            tip.cod = tip.cod.trim().toUpperCase();
            if (stare.schimburi.some(s => s.cod === tip.cod) || (!initial && stare.tipuriAbsenta.some(t => t.cod === tip.cod))) return {absente:null,eroare:{cod:'absente.cod_duplicat'}};
            const i = stare.tipuriAbsenta.findIndex(t => t.cod === initial);
            if (i >= 0 && stare.tipuriAbsenta[i].categorie !== tip.categorie && stare.absente.some(z => z.cod === initial)) return {absente:null,eroare:{cod:'absente.categorie_folosita'}};
            if (i < 0) stare.tipuriAbsenta.push(tip); else stare.tipuriAbsenta[i] = tip;
            raspuns(); return absente(an);
        },
        SalveazaAbsente: async (zile,an) => {
            if (zile.some(z => (stare.absente.find(v => v.angajat_id === z.angajat_id && v.data === z.data)?.revizie || '') !== z.revizie)) return {absente:null,eroare:{cod:'absente.conflict'}};
            if (zile.some(z => z.activ && stare.ciorne.some(c => c.angajat_id === z.angajat_id && c.data === z.data && c.valoare.trim()))) return {absente:null,eroare:{cod:'absente.pontaj_existent'}};
            if (zile.some(z => z.inlocuitor_id && z.inlocuitor_id === z.angajat_id)) return {absente:null,eroare:{cod:'absente.inlocuitor_propriu'}};
            for (const z of zile) {
                const i = stare.absente.findIndex(v => v.angajat_id === z.angajat_id && v.data === z.data);
                const existent = i >= 0 ? stare.absente[i] : null;
                let docId = z.document_id || existent?.document_id || '';
                let docNume = z.document_nume || existent?.document_nume || '';
                let docMime = z.document_mime || existent?.document_mime || '';
                if (z.document_continut_b64) {
                    docId = 'doc-' + Date.now() + '-' + Math.random().toString(16).slice(2,6);
                    docNume = z.document_nume || 'cerere.pdf';
                    docMime = z.document_mime || 'application/pdf';
                    stare.documenteAbsenta.push({id: docId, nume_fisier: docNume, tip_mime: docMime, continut_b64: z.document_continut_b64});
                }
                const nou = {
                    ...z,
                    nr_inregistrare: z.nr_inregistrare ?? existent?.nr_inregistrare ?? '',
                    inlocuitor_id: z.inlocuitor_id ?? existent?.inlocuitor_id ?? '',
                    stare_aprobare: z.stare_aprobare || existent?.stare_aprobare || 'IN_ASTEPTARE',
                    aprobat_de: z.aprobat_de ?? existent?.aprobat_de ?? '',
                    document_id: docId,
                    document_nume: docNume,
                    document_mime: docMime,
                    document_continut_b64: '',
                    revizie: String(i < 0 ? 1 : Number(stare.absente[i].revizie)+1)
                };
                if (i < 0) stare.absente.push(nou); else stare.absente[i] = nou;
            }
            raspuns(); return absente(an);
        },
        CitesteDocumentAbsenta: async id => {
            const d = stare.documenteAbsenta.find(x => x.id === id);
            return d ? {document: structuredClone(d), eroare: null} : {document: null, eroare: {cod: 'absente.document_inexistent'}};
        },
        CitesteSchita: async (an,luna)=>schita(an,luna),
        SalveazaSchita: async (an, luna, celule) => {
            if (celule.some(c => stare.absente.some(z => z.activ && z.angajat_id === c.angajat_id && z.data === c.data))) return {schita:null,eroare:{cod:'absente.celula_protejata'}};
            if (celule.some(c => (stare.ciorne.find(v => v.angajat_id === c.angajat_id && v.data === c.data)?.revizie || '') !== c.revizie)) return {schita: null, eroare: {cod: 'pontaj.conflict', param: {}}};
            for (const c of celule) {
                const i = stare.ciorne.findIndex(v => v.angajat_id === c.angajat_id && v.data === c.data);
                const nou = {...c, valoare: c.valoare.trim(), revizie: String(i < 0 ? 1 : Number(stare.ciorne[i].revizie)+1)};
                if (i < 0) stare.ciorne.push(nou); else stare.ciorne[i] = nou;
            }
            raspuns(); return schita(an,luna);
        }
    }}};
}

try {
    let port;
    await asteapta(() => {
        try {
            const linii = readFileSync(join(profil, 'DevToolsActivePort'), 'utf8').split('\n');
            if (!/^\d+$/.test(linii[0]) || !linii[1]?.startsWith('/devtools/browser/')) return false;
            port = linii[0]; return true;
        } catch (err) { if (['ENOENT','EBUSY'].includes(err.code)) return false; throw err; }
    }, 'Edge nu a pregătit portul de verificare.');
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
    const injectat = await cdp('Page.addScriptToEvaluateOnNewDocument', {source: `(${simulare})(${JSON.stringify(demo)},${JSON.stringify(cataloage)},'ro',${JSON.stringify(calendarSursa)},${JSON.stringify(versiuneAplicatie)},${JSON.stringify(schimburiSursa)},${JSON.stringify(rapoarteSursa)})`});
    await cdp('Page.navigate', {url});
    await asteapta(() => evalueaza(`Boolean(document.body?.textContent.includes('Adaugă primii membri'))`), 'Starea goală nu apare.');
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
    await completeaza('.card-departamente select', demo.sectii[0].id);
    await completeaza('.card-departamente input', 'Laborator'); await apasa('Adaugă departament');
    await asteapta(() => evalueaza(`document.querySelectorAll('.arbore-departamente ul li:not(.text-discret)').length === 1`), 'Primul departament nu apare.');
    await completeaza('.card-departamente input', 'Biochimie'); await apasa('Adaugă departament');
    await asteapta(() => evalueaza(`document.querySelectorAll('.arbore-departamente ul li:not(.text-discret)').length === 2`), 'O secție nu poate primi două departamente.');
    await completeaza('.card-departamente select', demo.sectii[1].id);
    await completeaza('.card-departamente input', 'Laborator Sud'); await apasa('Adaugă departament');
    await asteapta(() => evalueaza(`document.querySelectorAll('.arbore-departamente ul li:not(.text-discret)').length === 3`), 'Departamentul celeilalte secții nu apare.');
    const departamentNord = await evalueaza(`JSON.parse(localStorage.getItem('pontix-test')).departamente.find(d => d.denumire.ro === 'Laborator').id`);
    await apasa('Creează copie de siguranță acum');
    await asteapta(() => evalueaza(`document.querySelector('.card-backup [role="status"]')?.textContent.includes('pontix-manual-')`), 'Backupul manual din Setări nu afișează calea.');
    assert.equal(await evalueaza(`JSON.parse(localStorage.getItem('pontix-test')).backupuriManuale`), 1, 'Backupul manual nu a fost înregistrat.');
    await apasa('Angajați'); await completeaza('.cautare input', 'DEMO-001');
    await asteapta(() => evalueaza(`document.querySelector('tbody')?.textContent.includes('8 h')`), 'Norma întreagă nu se actualizează.');
    await apasa('Editează');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog[open]'))`), 'Formular atribuire departament lipsă.');
    assert(!(await evalueaza(`document.querySelector('dialog select[aria-label="Departament"]').textContent.includes('Laborator Sud')`)), 'Formularul arată departamente din altă secție.');
    await completeaza('dialog select[aria-label="Departament"]', departamentNord); await apasa('Salvează');
    await asteapta(() => evalueaza(`!document.querySelector('dialog') && document.querySelector('tbody').textContent.includes('Laborator')`), 'Departamentul angajatului nu se salvează.');
    await completeaza('.cautare input', ''); await completeaza('select[aria-label="Departament"]', departamentNord);
    await asteapta(() => evalueaza(`document.querySelectorAll('tbody tr').length === 1`), 'Filtrul de departament nu funcționează.');
    await completeaza('select[aria-label="Departament"]', '');
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
    assert.equal(await evalueaza(`document.querySelectorAll('.tabel-schimburi tbody tr').length`), 7, 'Programele confirmate nu apar.');
    const interval12 = await evalueaza(`[...document.querySelectorAll('.tabel-schimburi tbody tr')].find(r => r.cells[0].textContent === 'Z12').textContent`);
    assert(interval12.includes('07:00') && interval12.includes('19:00'), 'Z12 nu este afișat 07–19.');
    const intrare = await evalueaza(`[...document.querySelectorAll('.tabel-schimburi tbody tr')].find(r => r.cells[0].textContent === 'IN9').textContent`);
    assert(intrare.includes('22:00') && intrare.includes('07:00') && intrare.includes('9 h') && intrare.includes('+1') && intrare.includes('Intrare specială'), 'IN9 nu este distinct și nominal de 9 ore.');
    await completeaza('.schimburi select[aria-label="Regim de lucru"]','INTRARE_ZI_LIBERA');
    assert.equal(await evalueaza(`document.querySelectorAll('.tabel-schimburi tbody tr').length`),1,'Filtrul intrării speciale nu funcționează.');
    await apasa('Editează'); await asteapta(() => evalueaza(`document.querySelector('dialog select[name="regim"]')?.value === 'INTRARE_ZI_LIBERA' && document.querySelector('dialog select[name="categorie"]')?.value === 'LAB'`),'Regimul sau categoria specială lipsește în formular.');
    await apasa('Salvează'); await asteapta(() => evalueaza(`!document.querySelector('dialog')`),'Regimul nu se păstrează la editare.');
    await completeaza('.schimburi select[aria-label="Regim de lucru"]','LUCRATOARE');
    assert.deepEqual(await evalueaza(`[...document.querySelectorAll('.tabel-schimburi tbody tr')].map(r => r.cells[0].textContent).sort()`),['N9','S7','Z7'],'Filtrul programului obișnuit diferă.');
    const seara = await evalueaza(`[...document.querySelectorAll('.tabel-schimburi tbody tr')].find(r => r.cells[0].textContent === 'S7').textContent`);
    assert(seara.includes('15:00') && seara.includes('22:00'),'S7 nu respectă 15–22.');
    await completeaza('.schimburi select[aria-label="Regim de lucru"]','NELUCRATOARE');
    assert.deepEqual(await evalueaza(`[...document.querySelectorAll('.tabel-schimburi tbody tr')].map(r => r.cells[0].textContent).sort()`),['G24','N12','Z12'],'Filtrul weekend/sărbătoare diferă.');
    await completeaza('.schimburi select[aria-label="Regim de lucru"]','');
    await completeaza('.schimburi select[aria-label="Categorie"]','LAB');
    assert.deepEqual(await evalueaza(`[...document.querySelectorAll('.tabel-schimburi tbody tr')].map(r => r.cells[0].textContent).sort()`),['G24','IN9','N12','N9','S7','Z12','Z7'],'Filtrul categoriei Lab diferă.');
    await completeaza('.schimburi select[aria-label="Categorie"]','CUSTOM');
    assert.equal(await evalueaza(`document.querySelectorAll('.tabel-schimburi tbody tr').length`),0,'Categoria Custom nu este inițial goală în testul izolat.');
    await completeaza('.schimburi select[aria-label="Categorie"]','');
    await apasa('Completează schimburile inițiale');
    await asteapta(() => evalueaza(`document.querySelector('.schimburi [role="status"]')?.textContent.includes('0 schimburi')`), 'Completarea repetată dublează schimburile.');
    await apasa('Adaugă schimb');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog[open]')) && document.querySelector('dialog select[name="categorie"]')?.value === 'CUSTOM'`), 'Formular schimb lipsă sau categorie implicită incorectă.');
    await completeaza('dialog .campuri label:nth-child(1) input', 'test_n12');
    await completeaza('dialog .campuri label:nth-child(2) select', 'NOAPTE');
    await completeaza('dialog .campuri label:nth-child(3) input', 'Tură de test');
    await completeaza('dialog input[type="time"]', '19:00');
    await completeaza('dialog .campuri label:nth-child(6) input', '720');
    await completeaza('dialog .campuri label:nth-child(7) input', '30'); await apasa('Salvează');
    await asteapta(() => evalueaza(`!document.querySelector('dialog') && document.querySelector('.tabel-schimburi').textContent.includes('TEST_N12')`), 'Schimb nesalvat.');
    await completeaza('.schimburi select[aria-label="Categorie"]','CUSTOM');
    assert.deepEqual(await evalueaza(`[...document.querySelectorAll('.tabel-schimburi tbody tr')].map(r => r.cells[0].textContent)`),['TEST_N12'],'Filtrul categoriei Custom nu afișează schimbul nou.');
    await completeaza('.schimburi select[aria-label="Categorie"]','');
    assert(await evalueaza(`document.querySelector('.tabel-schimburi').textContent.includes('19:00') && document.querySelector('.tabel-schimburi').textContent.includes('12 h')`), 'Intervalul schimbului diferă.');
    await evalueaza(`[...document.querySelectorAll('.tabel-schimburi tbody tr')].find(r => r.cells[0].textContent === 'TEST_N12').querySelectorAll('button')[1].click()`);
    await asteapta(() => evalueaza(`[...document.querySelectorAll('.tabel-schimburi tbody tr')].find(r => r.cells[0].textContent === 'TEST_N12').textContent.includes('Inactiv')`), 'Schimbul nu se dezactivează.');
    await apasa('Reactivează');
    await asteapta(() => evalueaza(`document.querySelector('.tabel-schimburi').textContent.includes('Activ')`), 'Schimbul nu se reactivează.');
    await apasa('Adaugă schimb');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog[open]'))`), 'Formular duplicat lipsă.');
    await completeaza('dialog .campuri label:nth-child(1) input', 'TEST_N12');
    await completeaza('dialog .campuri label:nth-child(3) input', 'Duplicat');
    await completeaza('dialog input[type="time"]', '20:00');
    await completeaza('dialog .campuri label:nth-child(6) input', '720'); await apasa('Salvează');
    await asteapta(() => evalueaza(`document.querySelector('dialog [role="alert"]')?.textContent.includes('acest cod')`), 'Codul duplicat nu este semnalat.');
    await apasa('Anulează');
    await apasa('Concedii și absențe');
    await asteapta(() => evalueaza(`document.querySelectorAll('.tabel-tipuri-absenta tbody tr').length === 2`), 'Tipurile CO/CM nu apar.');
    const persoanaAbsenta = await evalueaza(`JSON.parse(localStorage.getItem('pontix-test')).angajati[0].id`);
    await completeaza('.filtre-absente select[aria-label="Angajat"]',persoanaAbsenta);
    await apasa('Adaugă tip de absență');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog input[name="cod"]'))`), 'Formular tip lipsă.');
    await completeaza('dialog input[name="cod"]','N12'); await completeaza('dialog input[name="ro"]','Eveniment'); await apasa('Salvează');
    await asteapta(() => evalueaza(`document.querySelector('dialog [role="alert"]')?.textContent.includes('deja folosit')`), 'Confuzie între cod de schimb și absență.');
    await completeaza('dialog input[name="cod"]','EV'); await completeaza('dialog input[name="en"]','Family event'); await apasa('Salvează');
    await asteapta(() => evalueaza(`!document.querySelector('dialog') && document.querySelectorAll('.tabel-tipuri-absenta tbody tr').length === 3`), 'Tipul special nu s-a salvat.');
    await apasa('Înregistrează zile');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog input[name="de_la"]'))`), 'Formular interval lipsă.');
    await completeaza('dialog input[name="de_la"]','2026-11-27'); await completeaza('dialog input[name="pana_la"]','2026-12-02');
    await asteapta(() => evalueaza(`document.querySelectorAll('.zile-absenta input').length === 6 && [...document.querySelectorAll('dialog button')].some(b => b.textContent.includes('luni–vineri') && !b.disabled)`), 'Calendarul intervalului nu s-a încărcat.');
    await apasa('Selectează luni–vineri, fără sărbători');
    assert.equal(await evalueaza(`document.querySelectorAll('.zile-absenta input:checked').length`),2,'Selecția include weekenduri/30 noiembrie/1 decembrie.');
    await apasa('Salvează');
    await asteapta(() => evalueaza(`!document.querySelector('dialog') && document.querySelectorAll('.tabel-absente tbody tr').length === 2`), 'Intervalul nu s-a salvat.');
    await apasa('Înregistrează zile'); await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog select[name="tip"]'))`),'Formular medical lipsă.');
    await completeaza('dialog select[name="tip"]','CM'); await completeaza('dialog input[name="de_la"]','2026-11-28'); await completeaza('dialog input[name="pana_la"]','2026-11-29');
    await apasa('Salvează');
    await asteapta(() => evalueaza(`!document.querySelector('dialog') && document.querySelectorAll('.tabel-absente tbody tr').length === 4`), 'Medicalul pe date explicite nu apare.');
    await apasa('Înregistrează zile'); await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog select[name="tip"]'))`),'Formular eveniment lipsă.');
    await completeaza('dialog select[name="tip"]','EV'); await completeaza('dialog input[name="de_la"]','2026-11-30'); await completeaza('dialog input[name="pana_la"]','2026-11-30'); await apasa('Salvează');
    await asteapta(() => evalueaza(`!document.querySelector('dialog') && document.querySelectorAll('.tabel-absente tbody tr').length === 5`),'Ziua specială nu apare.');
    assert.deepEqual(await evalueaza(`[...document.querySelectorAll('.sumar-absente tbody tr:first-child td')].slice(2).map(c => c.textContent)`),['2','2','1'],'Sumar incorect pe categorii.');
    await completeaza('.filtre-absente select[aria-label="Categorie"]','SPECIAL');
    assert.equal(await evalueaza(`document.querySelectorAll('.tabel-absente tbody tr').length`),1,'Filtrul categoriei nu funcționează.');
    await evalueaza(`[...document.querySelectorAll('.tabel-tipuri-absenta tbody tr')].find(r => r.cells[0].textContent === 'EV').querySelector('button').click()`);
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog input[name="activ"]'))`),'Editare tip lipsă.');
    await evalueaza(`document.querySelector('dialog input[name="activ"]').click()`); await apasa('Salvează');
    await asteapta(() => evalueaza(`!document.querySelector('dialog') && document.querySelector('.tabel-tipuri-absenta').textContent.includes('Inactiv')`),'Tipul nu se dezactivează.');
    assert.equal(await evalueaza(`document.querySelectorAll('.tabel-absente tbody tr').length`),1,'Dezactivarea a pierdut istoricul.');
    await completeaza('.filtre-absente select[aria-label="Categorie"]','');
    const persoanaInlocuitor = await evalueaza(`JSON.parse(localStorage.getItem('pontix-test')).angajati[1].id`);
    await evalueaza(`[...document.querySelectorAll('.tabel-absente tbody tr')].find(r => r.cells[2].textContent === '2026-11-28').querySelector('button').click()`);
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog select[name="tip"]'))`),'Editarea zilei lipsește.');
    await completeaza('dialog select[name="tip"]','CO');
    await completeaza('dialog input[name="nr_inregistrare"]','REG-104/2026');
    await completeaza('dialog select[name="inlocuitor"]',persoanaInlocuitor);
    await completeaza('dialog select[name="stare_aprobare"]','APROBATA');
    await completeaza('dialog input[name="aprobat_de"]','Dr. Ionescu');
    await incarcaFisier('dialog input[name="cerere_fisier"]',join(imaginiDemo,'cerere-demo.pdf'));
    await asteapta(() => evalueaza(`document.querySelector('.atasament-curent')?.textContent.includes('cerere-demo.pdf')`),'Cererea PDF nu este încărcată în formular.');
    await apasa('Salvează');
    await asteapta(() => evalueaza(`!document.querySelector('dialog') && document.querySelector('.sumar-absente tbody tr').cells[2].textContent === '3'`),'Editarea nu actualizează sumarul.');
    assert(await evalueaza(`document.querySelector('.tabel-absente').textContent.includes('REG-104/2026') && document.querySelector('.tabel-absente').textContent.includes('Dr. Ionescu')`),'Nr. înregistrare sau aprobarea nu apar în tabelul de absențe.');
    await evalueaza(`document.querySelector('.tabel-absente a[href="#previzualizare"]').click()`);
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog iframe.cadru-previzualizare-cerere')?.src.startsWith('blob:'))`),'Previzualizarea modală a cererii de concediu nu se deschide.');
    await apasa('Închide');
    await asteapta(() => evalueaza(`!document.querySelector('dialog')`),'Fereastra de previzualizare a cererii nu se închide.');
    await completeaza('.filtre-absente select[aria-label="Aprobare"]','APROBATA');
    assert.equal(await evalueaza(`document.querySelectorAll('.tabel-absente tbody tr').length`),1,'Filtrul stării de aprobare nu funcționează.');
    await completeaza('.filtre-absente select[aria-label="Aprobare"]','');
    await completeaza('.filtre-absente select[aria-label="Categorie"]','MEDICAL');
    await evalueaza(`document.querySelector('.tabel-absente thead input').click()`); await apasa('Anulează zilele selectate (1)');
    await asteapta(() => evalueaza(`document.querySelectorAll('.tabel-absente tbody tr').length === 0`),'Anularea selecției nu funcționează.');
    await completeaza('.filtre-absente select[aria-label="Categorie"]','');
    assert.deepEqual(await evalueaza(`[...document.querySelectorAll('.sumar-absente tbody tr:first-child td')].slice(2).map(c => c.textContent)`),['3','0','1'],'Anularea nu actualizează sumarul.');
    assert(await evalueaza(`document.querySelector('.sumar-absente tbody tr:first-child td:nth-child(2)').textContent.includes('Rămas:') && document.querySelector('.sumar-absente tbody tr:first-child td:nth-child(2)').textContent.includes('în așteptare: 2')`),'Sumarul nu afișează diferența aritmetică rămasă și zilele în așteptare.');
    await apasa('Înregistrează zile');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog select[name="angajat"]'))`),'Formular verificare înlocuitor lipsă.');
    await completeaza('dialog select[name="angajat"]',persoanaInlocuitor);
    await completeaza('dialog input[name="de_la"]','2026-11-28');
    await completeaza('dialog input[name="pana_la"]','2026-11-28');
    await apasa('Salvează');
    await asteapta(() => evalueaza(`document.querySelector('dialog [role="alert"]')?.textContent.includes('deja desemnat ca înlocuitor')`),'Suprapunerea titular-înlocuitor nu este blocată în formular.');
    await apasa('Anulează');
    const capturaAbsente = await cdp('Page.captureScreenshot',{format:'png'}); writeFileSync(join(radacina,'build/bin/concedii-absente.png'),Buffer.from(capturaAbsente.data,'base64'));
    await apasa('Pontaj');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('.unelte-perioada select'))`),'Pontaj lipsă după concedii.');
    assert(await evalueaza(`[...document.querySelectorAll('select[aria-label="Alege programul de lucru"] option')].some(o => o.value === 'IN9' && o.textContent.includes('Intrare specială'))`),'Sugestiile din pontaj nu disting IN9.');
    await completeaza('.unelte-perioada select','11');
    await asteapta(() => evalueaza(`document.querySelectorAll('.grila-pontaj input.celula-absenta').length === 3`),'Pontajul nu arată absențele din noiembrie.');
    assert(await evalueaza(`[...document.querySelectorAll('.grila-pontaj input.celula-absenta')].some(i => i.title.includes('REG-104/2026') && i.title.includes('Dr. Ionescu'))`),'Tooltip-ul din pontaj nu afișează detaliile cererii de concediu.');
    assert(await evalueaza(`Boolean(document.querySelector('.grila-pontaj td.celula-inlocuitor input')?.title.includes('Înlocuiește:'))`),'Grila de pontaj nu evidențiază ziua pe rândul persoanei care înlocuiește.');
    const protejata = await evalueaza(`(() => { const e = document.querySelector('.grila-pontaj input.celula-absenta'); e.focus(); return {selector:'input[data-rand="'+e.dataset.rand+'"][data-zi="'+e.dataset.zi+'"]',valoare:e.value}; })()`);
    await completeaza('.unelte-grila input','7'); await apasa('Aplică orele manuale');
    await asteapta(() => evalueaza(`document.querySelector('.pontaj-lunar [role="alert"]')?.textContent.includes('Concedii și absențe')`),'Grila suprascrie absența.');
    assert.equal(await evalueaza(`document.querySelector(${JSON.stringify(protejata.selector)}).value`),protejata.valoare,'Absență suprascrisă.');
    await apasa('Concedii și absențe');
    await asteapta(() => evalueaza(`document.querySelectorAll('.tabel-absente tbody tr').length === 4`),'Revenirea a pierdut zilele salvate.');
    await evalueaza(`document.querySelector('.tabel-absente tbody tr button').click()`);
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog select[name="tip"]'))`),'Formular conflict lipsă.');
    await completeaza('dialog select[name="tip"]','CM');
    await evalueaza(`(async () => { const r = await window.go.main.App.CitesteAbsente(2026); await window.go.main.App.SalveazaAbsente([r.absente.zile[0]],2026); })()`);
    await apasa('Salvează');
    await asteapta(() => evalueaza(`document.querySelector('dialog [role="alert"]')?.textContent.includes('modificate între timp')`),'Conflictul absențelor nu este semnalat.');
    await apasa('Anulează'); await apasa('Reîncarcă evidența');
    await asteapta(() => evalueaza(`document.querySelectorAll('.tabel-absente tbody tr').length === 4`),'Reîncărcarea absențelor a eșuat.');
    await completeaza('.modul-absente .unelte-perioada input','2027'); await apasa('Arată anul');
    await asteapta(() => evalueaza(`document.querySelectorAll('.tabel-absente tbody tr').length === 0 && document.querySelectorAll('.tabel-tipuri-absenta tbody tr').length === 3`),'Filtrarea anuală a absențelor a eșuat.');
    assert(await evalueaza(`document.documentElement.scrollWidth <= 1024`),'Modulul de absențe depășește fereastra.');

    await apasa('Configurare pontaj');
    await asteapta(()=>evalueaza(`document.querySelectorAll('.reguli-pontaj tbody tr').length===12`),'Cele 12 reguli nu apar.');
    assert.equal(await evalueaza(`document.querySelector('select[name="sursa_ore_garda"]').value`),'NEDEFINITA','Garda are o regulă inventată.');
    await completeaza('select[name="atribuire"]','ZILE_CALENDARISTICE');
    await completeaza('select[name="sursa_ore_garda"]','ORE_SCHIMBURI');
    await evalueaza(`(()=>{const s=document.querySelector('select[name="coduri_ore_garda"]');[...s.options].forEach(o=>o.selected=o.value==='G24');s.dispatchEvent(new Event('change',{bubbles:true}));})()`);
    await completeaza('input[name="factor_ore_garda"]','2');await apasa('Salvează');
    await asteapta(()=>evalueaza(`JSON.parse(localStorage.getItem('pontix-test')).reguliPontaj.revizie==='2'`),'Regulile nu s-au salvat.');
    await completeaza('select[name="atribuire"]','ZI_INCEPUT');
    await evalueaza(`window.confirm=()=>false`);await apasa('Angajați');
    assert.equal(await evalueaza(`document.querySelector('h1').textContent`),'Configurare pontaj','Configurația nesalvată a fost pierdută.');
    await evalueaza(`(async()=>{const r=await window.go.main.App.CitesteConfiguratiePontaj();await window.go.main.App.SalveazaConfiguratiePontaj(r.configuratie);})()`);
    await apasa('Salvează');await asteapta(()=>evalueaza(`document.querySelector('.reguli-pontaj [role="alert"]')?.textContent.includes('între timp')`),'Conflictul regulilor nu este afișat.');
    await evalueaza(`window.confirm=()=>true`);await apasa('Reîncarcă regulile');
    await asteapta(()=>evalueaza(`document.querySelector('select[name="atribuire"]').value==='ZILE_CALENDARISTICE'`),'Reîncărcarea regulilor nu funcționează.');
    await completeaza('select[name="atribuire"]','ZI_INCEPUT');await apasa('Salvează');
    await asteapta(()=>evalueaza(`!document.querySelector('.reguli-pontaj .buton.primar:not(:disabled)')`),'Salvarea configurației nu s-a încheiat.');
    const capturaReguli=await cdp('Page.captureScreenshot',{format:'png'});writeFileSync(join(radacina,'build/bin/configurare-pontaj.png'),Buffer.from(capturaReguli.data,'base64'));
    await apasa('Definire rapoarte');
    await asteapta(() => evalueaza(`document.querySelectorAll('.tabel-modele-rapoarte tbody tr').length===4`),'Modelele inițiale nu apar.');
    await evalueaza(`[...document.querySelectorAll('.tabel-modele-rapoarte tbody tr')].find(r => r.cells[0].textContent==='Pontaj lunar').querySelector('button').click()`);
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog textarea[name="antet_stanga"]'))`),'Editorul de raport nu apare.');
    assert.deepEqual(await evalueaza(`[...document.querySelectorAll('dialog legend')].slice(0,3).map(e=>e.textContent)`),['1 · Antet','2 · Conținut','3 · Subsol'],'Cele trei componente nu sunt separate.');
    await completeaza('dialog textarea[name="antet_stanga"]','Unitate demo\n{{sectie}}');
    await completeaza('dialog textarea[name="subsol_dreapta"]','Avizat demo');
    await incarcaFisier('dialog input[name="antet_imagine_stanga"]',join(imaginiDemo,'antet-demo.png'));
    await asteapta(()=>evalueaza(`document.querySelector('.zona-model-raport img')?.src.startsWith('data:image/png;base64,')`),'Imaginea PNG nu este importata.');
    await incarcaFisier('dialog input[name="subsol_imagine_dreapta"]',join(imaginiDemo,'subsol-demo.jpg'));
    await asteapta(()=>evalueaza(`document.querySelectorAll('.imagine-model-raport img').length===2`),'Imaginea JPEG nu este importata.');
    await completeaza('dialog input[name="antet_inaltime_stanga"]','22');
    await incarcaFisier('dialog input[name="antet_imagine_stanga"]',join(imaginiDemo,'subsol-demo.jpg'));
    await asteapta(()=>evalueaza(`document.querySelector('input[name="antet_imagine_stanga"]').closest('.zona-model-raport').querySelector('img').alt==='subsol-demo.jpg'`),'Inlocuirea imaginii nu functioneaza.');
    await incarcaFisier('dialog input[name="antet_imagine_stanga"]',join(imaginiDemo,'imagine-invalida.txt'));
    await asteapta(()=>evalueaza(`document.querySelector('dialog [role="alert"]')?.textContent.includes('Imaginea nu este validă')`),'Fisierul care nu este imagine nu este refuzat.');
    assert.equal(await evalueaza(`document.querySelector('input[name="antet_imagine_stanga"]').closest('.zona-model-raport').querySelector('img').alt`),'subsol-demo.jpg','Importul invalid pierde imaginea existenta.');
    await incarcaFisier('dialog input[name="antet_imagine_stanga"]',join(imaginiDemo,'antet-demo.png'));
    await asteapta(()=>evalueaza(`document.querySelector('input[name="antet_imagine_stanga"]').closest('.zona-model-raport').querySelector('img').alt==='antet-demo.png'`),'Restaurarea PNG nu functioneaza.');
    await completeaza('dialog textarea[name="introducere"]','{{executa_sql}}');await apasa('Salvează');
    await asteapta(() => evalueaza(`document.querySelector('dialog [role="alert"]')?.textContent.includes('câmp necunoscut')`),'Câmpul necunoscut nu este refuzat.');
    await completeaza('dialog textarea[name="introducere"]','Perioada {{perioada}}');
    await completeaza('dialog select[name="limba_texte"]','en');
    await completeaza('dialog textarea[name="antet_stanga"]','Demo organisation\n{{sectie}}');
    await completeaza('dialog select[name="limba_texte"]','ro');
    assert.equal(await evalueaza(`document.querySelector('dialog textarea[name="antet_stanga"]').value`),'Unitate demo\n{{sectie}}','Schimbarea limbii a pierdut textul RO.');
    await evalueaza(`document.querySelector('dialog input[name="coloana_functie"]').click()`);
    await completeaza('dialog input[name="margine_mm"]','12');
    await completeaza('dialog input[name="font_pt"]','9');
    const capturaEditor=await cdp('Page.captureScreenshot',{format:'png'});writeFileSync(join(radacina,'build/bin/definire-rapoarte.png'),Buffer.from(capturaEditor.data,'base64'));
    await apasa('Salvează');
    await asteapta(() => evalueaza(`!document.querySelector('dialog') && [...document.querySelectorAll('.tabel-modele-rapoarte tbody tr')].find(r=>r.cells[0].textContent==='Pontaj lunar').cells[2].textContent==='2'`),'Modelul nu se salvează cu revizie nouă.');
    const modeleDupaEditare=await evalueaza(`JSON.parse(localStorage.getItem('pontix-test')).modeleRapoarte`);
    assert.equal(modeleDupaEditare.find(d=>d.id==='raport-pontaj').antet.imagini.stanga.inaltime_mm,22,'Inaltimea imaginii s-a pierdut la inlocuire/salvare.');
    assert.equal(modeleDupaEditare.find(d=>d.id==='raport-pontaj').subsol.imagini.dreapta.nume,'subsol-demo.jpg','Imaginea din subsol nu este salvata.');
    assert.equal(modeleDupaEditare.find(d=>d.id==='raport-fisa').antet.stanga.ro,'','Editarea unui raport a schimbat alt model.');
    assert(modeleDupaEditare.find(d=>d.id==='raport-pontaj').continut.coloane.includes('functie'),'Coloana configurată s-a pierdut.');
    await evalueaza(`[...document.querySelectorAll('.tabel-modele-rapoarte tbody tr')].find(r=>r.cells[0].textContent==='Pontaj lunar').querySelectorAll('button')[1].click()`);
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog input[name="denumire_ro"]'))`),'Duplicarea nu deschide editorul.');
    await completeaza('dialog input[name="denumire_ro"]','Raport demo separat');
    await evalueaza(`document.querySelector('input[name="antet_imagine_stanga"]').closest('.zona-model-raport').querySelector('button').click()`);
    assert(await evalueaza(`!document.querySelector('input[name="antet_imagine_stanga"]').closest('.zona-model-raport').querySelector('img')`),'Imaginea nu poate fi eliminata din copia modelului.');
    await evalueaza(`document.querySelector('dialog input[name="activ_model"]').click()`);await apasa('Salvează');
    await asteapta(() => evalueaza(`!document.querySelector('dialog')&&document.querySelectorAll('.tabel-modele-rapoarte tbody tr').length===5`),'Copia independentă nu a fost creată.');
    await evalueaza(`[...document.querySelectorAll('.tabel-modele-rapoarte tbody tr')].find(r=>r.cells[0].textContent==='Pontaj lunar').querySelector('button').click()`);
    await asteapta(() => evalueaza(`Boolean(document.querySelector('dialog textarea[name="antet_stanga"]'))`),'Formular conflict model lipsă.');
    assert.equal(await evalueaza(`document.querySelector('input[name="antet_imagine_stanga"]').closest('.zona-model-raport').querySelector('img').alt`),'antet-demo.png','Imaginea modelului original s-a pierdut la duplicare sau la redeschidere.');
    await completeaza('dialog textarea[name="antet_stanga"]','Modificare locală nesalvată');
    await evalueaza(`(async()=>{const r=await window.go.main.App.ListaModeleRapoarte();await window.go.main.App.SalveazaModelRaport(r.configuratie.modele.find(d=>d.id==='raport-pontaj'));})()`);
    await apasa('Salvează');
    await asteapta(() => evalueaza(`document.querySelector('dialog [role="alert"]')?.textContent.includes('modificat între timp')`),'Conflictul definiției nu este semnalat.');
    assert.equal(await evalueaza(`document.querySelector('dialog textarea[name="antet_stanga"]').value`),'Modificare locală nesalvată','Conflictul a pierdut editarea locală.');
    await apasa('Anulează');await apasa('Rapoarte');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('select[name="model_raport"]'))`),'Meniul Rapoarte lipsește.');
    assert.equal(await evalueaza(`document.querySelectorAll('select[name="model_raport"] option').length`),4,'Copia inactivă apare la generare.');
    await completeaza('select[name="model_raport"]','raport-pontaj');await completeaza('input[name="an_raport"]','2026');await completeaza('select[name="luna_raport"]','10');
    await completeaza('select[name="sectie_raport"]',demo.sectii[0].id);await completeaza('select[name="limba_raport"]','en');await apasa('Generează și previzualizează');
    await asteapta(() => evalueaza(`document.querySelector('iframe.previzualizare-pdf')?.src.startsWith('blob:')`),'PDF-ul nu apare în previzualizare.');
    assert(await evalueaza(`document.querySelector('.rezultat-raport [role="status"]').textContent.includes('Pagini: 4')`),'Numarul de pagini din PDF nu este afisat.');
    const cerereRaport=await evalueaza(`JSON.parse(localStorage.getItem('pontix-test')).ultimaCerereRaport`);
    assert(cerereRaport.an===2026&&cerereRaport.luna===10&&cerereRaport.sectie_id===demo.sectii[0].id&&cerereRaport.limba==='en','Filtrele raportului nu sunt trimise corect.');
    await apasa('Salvează PDF');await asteapta(() => evalueaza(`document.querySelector('.generare-rapoarte [role="status"]')?.textContent.includes('PDF salvat')`),'Salvarea PDF nu este semnalată.');
    await apasa('Deschide pentru tipărire');await asteapta(() => evalueaza(`document.querySelector('.generare-rapoarte [role="status"]')?.textContent.includes('deschis')`),'Deschiderea pentru tipărire nu este semnalată.');
    assert.equal(await evalueaza(`JSON.parse(localStorage.getItem('pontix-test')).tipaririRaport`),1,'Acțiunea de tipărire nu este apelată o singură dată.');
    const capturaRaport=await cdp('Page.captureScreenshot',{format:'png'});writeFileSync(join(radacina,'build/bin/previzualizare-raport.png'),Buffer.from(capturaRaport.data,'base64'));
    await completeaza('select[name="luna_raport"]','11');
    assert(await evalueaza(`!document.querySelector('iframe.previzualizare-pdf')`),'Schimbarea perioadei păstrează PDF-ul vechi ca și când ar fi nou.');
    await completeaza('select[name="model_raport"]','raport-fisa');
    assert(await evalueaza(`document.querySelector('select[name="angajat_raport"]').required`),'Fișa individuală nu cere angajat.');
    await completeaza('select[name="angajat_raport"]',persoanaAbsenta);await apasa('Generează și previzualizează');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('iframe.previzualizare-pdf'))`),'Fișa individuală nu se previzualizează.');
    assert(await evalueaza(`document.documentElement.scrollWidth<=1024`),'Rapoartele depășesc fereastra.');
    await apasa('Definire rapoarte');
    await asteapta(()=>evalueaza(`Boolean(document.querySelector('.tabel-modele-rapoarte'))`),'Modelele nu reapar.');
    await evalueaza(`[...document.querySelectorAll('.tabel-modele-rapoarte tbody tr')].find(r=>r.cells[0].textContent==='Pontaj detaliat').querySelector('button').click()`);
    await asteapta(()=>evalueaza(`Boolean(document.querySelector('dialog input[name="coloana_totaluri"]'))`),'Grupul final nu apare în editor.');
    assert(await evalueaza(`document.querySelector('dialog input[name="coloana_totaluri"]').disabled && document.querySelector('dialog input[name="coloana_zile"]').disabled`),'Coloanele obligatorii pot fi eliminate.');
    await apasa('Anulează');await apasa('Rapoarte');
    await asteapta(()=>evalueaza(`Boolean(document.querySelector('select[name="model_raport"]'))`),'Generatorul lipsește.');
    await completeaza('select[name="model_raport"]','raport-pontaj-detaliat');await apasa('Generează și previzualizează');
    await asteapta(()=>evalueaza(`document.querySelector('iframe.previzualizare-pdf') && JSON.parse(localStorage.getItem('pontix-test')).ultimaCerereRaport.model_id==='raport-pontaj-detaliat'`),'Raportul detaliat nu este separat.');
    // Încărcăm 50 de persoane strict în profilul simulat, fără date reale.
    await evalueaza(`(() => { const s = JSON.parse(localStorage.getItem('pontix-test')); const bio = s.departamente.find(d => d.denumire.ro === 'Biochimie').id; s.angajati = Array.from({length:50}, (_, i) => ({...s.angajati[0], id:'grila-'+i, nume:'Exemplu '+String(i+1).padStart(2,'0'), prenume:'Test', marca:'TEST-'+String(i+1).padStart(3,'0'), activ:true, departament_id: i % 2 ? bio : s.angajati[0].departament_id})); s.ciorne = []; localStorage.setItem('pontix-test', JSON.stringify(s)); })()`);
    await cdp('Page.reload');
    await asteapta(() => evalueaza(`document.querySelectorAll('tbody tr').length === 50`), 'Fixture 50 persoane lipsă.');
    await apasa('Pontaj');
    await asteapta(() => evalueaza(`Boolean(document.querySelector('.grila-pontaj tbody input'))`), 'Grila lunară lipsește.');
    await completeaza('.unelte-perioada input', '2026'); await apasa('Arată anul');
    await completeaza('.unelte-perioada select', '10');
    await asteapta(() => evalueaza(`document.querySelectorAll('.grila-pontaj tbody input').length === 1550`), 'Grila 50×31 nu are toate celulele.');
    await completeaza('.unelte-grila select[aria-label="Departament"]', departamentNord);
    await asteapta(() => evalueaza(`document.querySelectorAll('.grila-pontaj tbody input').length === 775`), 'Filtrul departamentului în grilă nu funcționează.');
    await completeaza('.unelte-grila select[aria-label="Departament"]', '');
    await asteapta(() => evalueaza(`document.querySelectorAll('.grila-pontaj tbody input').length === 1550`), 'Eliminarea filtrului nu reface grila.');
    await evalueaza(`document.querySelector('input[data-rand="0"][data-zi="0"]').focus()`);
    await completeaza('select[aria-label="Alege programul de lucru"]', 'N12');await apasa('Aplică programul');
    await asteapta(()=>evalueaza(`document.querySelector('input[data-rand="0"][data-zi="0"]').value==='12'`),'Programul nu este afișat în ore.');
    await apasa('Salvează schița');await asteapta(()=>evalueaza(`JSON.parse(localStorage.getItem('pontix-test')).ciorne.some(c=>c.valoare==='N12')`),'Codul programului nu a fost păstrat.');
    await evalueaza(`document.querySelector('input[data-rand="0"][data-zi="0"]').focus()`);
    await evalueaza(`document.querySelector('input[data-rand="0"][data-zi="0"]').dispatchEvent(new KeyboardEvent('keydown', {key:'ArrowRight',shiftKey:true,bubbles:true}))`);
    await completeaza('.unelte-grila input', '7'); await apasa('Aplică orele manuale');
    await asteapta(() => evalueaza(`document.querySelector('input[data-rand="0"][data-zi="1"]').value === '7'`), 'Completarea intervalului nu funcționează.');
    const copiat = await evalueaza(`(() => { const dt = new DataTransfer(); document.querySelector('input[data-rand="0"][data-zi="1"]').dispatchEvent(new ClipboardEvent('copy', {bubbles:true, clipboardData:dt})); return dt.getData('text/plain'); })()`);
    assert.equal(copiat, '7\t7', 'Copiere TSV greșită.');
    const textLipire = '7\t12\r\n5\t7.5\r\n';
    await evalueaza(`document.querySelector('input[data-rand="1"][data-zi="0"]').focus()`);
    await evalueaza(`new Promise(resolve => requestAnimationFrame(() => resolve(true)))`);
    await evalueaza(`(() => { const input = document.querySelector('input[data-rand="1"][data-zi="0"]'); const dt = new DataTransfer(); dt.setData('text/plain',${JSON.stringify(textLipire)}); input.dispatchEvent(new ClipboardEvent('paste',{bubbles:true,clipboardData:dt})); })()`);
    await asteapta(() => evalueaza(`document.querySelector('input[data-rand="2"][data-zi="1"]').value === '7.5'`), 'Lipirea TSV nu funcționează.');
    await apasa('Salvează schița');
    await asteapta(() => evalueaza(`document.querySelector('.pontaj-lunar [role="status"]')?.textContent.includes('salvată')`), 'Schița nu se salvează.');
    const capturaGrila = await cdp('Page.captureScreenshot', {format: 'png'});
    writeFileSync(join(radacina, 'build/bin/pontaj-schita.png'), Buffer.from(capturaGrila.data, 'base64'));
    await apasa('Angajați'); await apasa('Pontaj');
    await asteapta(() => evalueaza(`document.querySelector('input[data-rand="2"][data-zi="1"]')?.value === '7.5'`), 'Schița nu se păstrează la revenire.');
    await completeaza('.unelte-perioada select', '2');
    await asteapta(() => evalueaza(`document.querySelectorAll('.grila-pontaj tbody input').length === 1400`), 'Februarie 2026 nu are 28 de zile.');
    await completeaza('.unelte-perioada input', '2028'); await apasa('Arată anul');
    await asteapta(() => evalueaza(`document.querySelectorAll('.grila-pontaj tbody input').length === 1450`), 'Februarie 2028 nu are 29 de zile.');
    await completeaza('.unelte-perioada input', '2026'); await apasa('Arată anul');
    await completeaza('.unelte-perioada select', '1');
    await asteapta(() => evalueaza(`document.querySelector('.grila-pontaj thead th:nth-child(3)')?.classList.contains('zi-sarbatoare')`), 'Sărbătoarea nu este marcată în grilă.');
    await completeaza('.unelte-perioada select', '10');
    await asteapta(() => evalueaza(`document.querySelector('input[data-rand="0"][data-zi="0"]')?.value === '7'`), 'Revenirea în octombrie a pierdut datele.');
    await evalueaza(`document.querySelector('input[data-rand="0"][data-zi="0"]').focus()`);
    await completeaza('.unelte-grila input','24');await apasa('Aplică orele manuale');
    await asteapta(()=>evalueaza(`document.querySelector('input[data-rand="0"][data-zi="0"]').value==='24'`),'Ore manuale lipsă.');
    await evalueaza(`window.confirm = () => false`); await apasa('Angajați');
    assert(await evalueaza(`document.querySelector('h1').textContent === 'Pontaj'`), 'Navigarea a abandonat schița fără acord.');
    await completeaza('.unelte-perioada input', '2027'); await apasa('Arată anul');
    assert(await evalueaza(`document.querySelector('.unelte-perioada input').value === '2026' && document.querySelector('input[data-rand="0"][data-zi="0"]').value === '24'`), 'Anul afișat diferă de grilă după anulare.');
    await evalueaza(`window.confirm = () => true`);
    await evalueaza(`(async () => { const r = await window.go.main.App.CitesteSchita(2026,10); const c = r.schita.celule.find(c => c.angajat_id === 'grila-0' && c.data === '2026-10-01'); await window.go.main.App.SalveazaSchita(2026,10,[{...c, valoare:'extern'}]); })()`);
    await apasa('Salvează schița');
    await asteapta(() => evalueaza(`document.querySelector('.pontaj-lunar [role="alert"]')?.textContent.includes('altă instanță')`), 'Conflictul nu este semnalat în UI.');
    assert(await evalueaza(`document.querySelector('input[data-rand="0"][data-zi="0"]').value === '24'`), 'Conflictul a pierdut modificarea locală.');
    await apasa('Reîncarcă luna');
    await asteapta(() => evalueaza(`document.querySelector('input[data-rand="0"][data-zi="0"]')?.value === '—'`), 'Reîncărcarea nu arată revizia nouă.');
    await completeaza('.unelte-grila select[aria-label="Departament"]', departamentNord);
    await asteapta(() => evalueaza(`document.querySelectorAll('.grila-pontaj tbody input').length === 775`), 'Filtrul departamentului înainte de precompletare nu funcționează.');
    await completeaza('select[aria-label="Alege programul de lucru"]', 'Z7'); await apasa('Precompletează L–V');
    await asteapta(() => evalueaza(`document.querySelector('input[data-rand="0"][data-zi="4"]')?.value === '7' && document.querySelector('input[data-rand="0"][data-zi="2"]')?.value === ''`), 'Precompletarea L–V nu completează doar zilele lucrătoare.');
    await evalueaza(`document.querySelector('input[data-rand="0"][data-zi="0"]').focus()`);
    await completeaza('select[aria-label="Alege programul de lucru"]', 'Z7'); await apasa('Aplică programul');
    await evalueaza(`document.querySelector('input[data-rand="0"][data-zi="2"]').focus()`);
    await completeaza('select[aria-label="Alege programul de lucru"]', 'G24'); await apasa('Aplică programul');
    await evalueaza(`document.querySelector('input[data-rand="0"][data-zi="3"]').focus()`);
    await completeaza('select[aria-label="Alege programul de lucru"]', 'Z7'); await apasa('Aplică programul');
    await asteapta(() => evalueaza(`document.querySelector('input[data-rand="0"][data-zi="3"]')?.closest('td')?.classList.contains('celula-avertisment') && Boolean(document.querySelector('.grila-pontaj td.total-diferenta'))`), 'Avertismentul după G24 sau diferența de total nu sunt evidențiate.');
    assert(await evalueaza(`document.documentElement.scrollWidth <= 1024`), 'Grila depășește fereastra în afara containerului de derulare.');
    await cdp('Page.removeScriptToEvaluateOnNewDocument', {identifier: injectat.identifier});
    await cdp('Page.addScriptToEvaluateOnNewDocument', {source: `(${simulare})(${JSON.stringify(demo)},${JSON.stringify(cataloage)},'en',${JSON.stringify(calendarSursa)},${JSON.stringify(versiuneAplicatie)},${JSON.stringify(schimburiSursa)},${JSON.stringify(rapoarteSursa)})`});
    await cdp('Page.reload');
    await asteapta(() => evalueaza(`document.querySelector('h1')?.textContent === 'Employees'`), 'Traducerea EN nu apare.');
    await apasa('Leave and absences');
    await asteapta(() => evalueaza(`document.querySelector('h1')?.textContent === 'Leave and absences' && document.querySelector('.tabel-tipuri-absenta')?.textContent.includes('Family event')`),'Modulul de absențe nu este tradus în EN.');
    assert.deepEqual(erori, [], 'Erori JavaScript în interfață.');
    console.log('UI verificată: angajați/departamente, calendar 2026, schimburi și categorii LAB/CUSTOM, concedii/medical/speciale, intervale și sumar, anulare și conflicte, protecție în pontaj, grilă 50×31/TSV; definiții rapoarte, antet/conținut/subsol, RO/EN, duplicare și conflicte, previzualizare PDF, export și tipărire simulate; 1024px și EN.');
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
