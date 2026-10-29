// F14-T24: regresii în componenta completă, numai date fictive și Wails simulat.
import React from 'react';
import {createRoot} from 'react-dom/client';
import i18next from 'i18next';
import {initReactI18next} from 'react-i18next';
import {AsistentAI} from '../src/AsistentAI';
import {initializeazaSesiuneAI, seteazaConversatie} from '../src/asistentSesiune';
import ro from '../../internal/i18n/locales/ro.json';
import en from '../../internal/i18n/locales/en.json';
import '../src/style.css';

const parametri = new URLSearchParams(location.search);
const limba = parametri.get('limba') || 'ro', caz = parametri.get('caz') || 'rezultat';
const faraMotor = caz === 'fara-motor';
await i18next.use(initReactI18next).init({lng: limba, resources: {ro: {translation: ro}, en: {translation: en}}, interpolation: {escapeValue: false}});
const t = (cheie, param) => i18next.t(cheie, param);
const intentie = {selectie: {tip: 'toti', valoare: '', marci: [], fara_pontaj: false}, de_la: '2026-10-01', pana_la: '2026-11-30', cod_schimb: 'Z7', filtru_zile: 'lucratoare', mod_existente: 'pastreaza', mesaj_original: 'Toți angajații în octombrie și noiembrie lucrează de la 7 la 14', clarificari: []};
const persoane = Array.from({length: 3}, (_, x) => ({id: `p${x}`, marca: `PROBA-${x + 1}`, nume: `Persoană ${x + 1}`, prenume: 'Fictivă', activ: true}));
const date = {angajati: persoane, sectii: [], departamente: [], categorii_program: []};
const randuri = [];
for (const a of persoane) for (let d = new Date('2026-10-01T00:00:00Z'); d <= new Date('2026-11-30T00:00:00Z'); d.setUTCDate(d.getUTCDate() + 1)) {
    if (![0, 6].includes(d.getUTCDay())) randuri.push({angajat_id: a.id, marca: a.marca, nume: a.nume, data: d.toISOString().slice(0, 10), inainte: '', dupa: 'Z7', stare: 'noua', motiv: ''});
}
Object.assign(randuri[0], {inainte: '6', stare: 'ocupata'});
Object.assign(randuri[1], {stare: 'blocata', motiv: 'absenta'});
Object.assign(randuri[2], {inainte: 'Z7', stare: 'identica'});
Object.assign(randuri[3], {inainte: '4', stare: 'inlocuire'});
const plan = {intentie, persoane, luni: ['2026-10', '2026-11'], randuri, token: 'proba', noi: 125, inlocuiri: 1, identice: 1, ocupate: 1, blocate: 1, aplicate: 126};
let confirmari = 0, intentiePrevizualizata, lunaDeschisa;
window.runtime = {EventsOnMultiple: () => () => {}};
window.go = {main: {App: {
    StareAI: async () => ({stare: {instalat: !faraMotor, activ: !faraMotor, port: 5000, cale_folder: 'date-fictive/ai'}}), EsteAutoPornireAI: async () => true, SarcinaAIActiva: async () => '',
    ListaSabloaneAI: async () => ({sabloane: []}),
    ListaSchimburi: async () => ({schimburi: [{id: 'z7', cod: 'Z7', denumire: {ro: 'Zi 07:00–14:00', en: 'Day 07:00–14:00'}, activ: true}]}),
    PrevizualizeazaPlanAI: async i => {intentiePrevizualizata = i; return {plan: {...plan, intentie: i, aplicate: 0}};},
    ConfirmaPlanAI: async (i, token) => {
        confirmari++;
        if (token !== 'proba' || JSON.stringify(i) !== JSON.stringify(intentiePrevizualizata)) throw Error('Planul confirmat diferă de previzualizare');
        return caz === 'eroare' ? {eroare: {cod: 'ai.plan_expirat'}} : {plan: {...plan, intentie: i}};
    },
}}};
const mesaje = Array.from({length: 16}, (_, x) => ({id: `istoric-${x}`, rol: x % 2 ? 'assistant' : 'user', text: `Mesaj fictiv ${x + 1}. ${'Conținut de probă. '.repeat(8)}`}));
mesaje.splice(2, 0, {id: 'plan-vechi', rol: 'assistant', text: '', intentie});
mesaje.push({id: 'plan-nou', rol: 'assistant', text: '', intentie: {...intentie, de_la: '2026-11-01'}});
sessionStorage.setItem('pontix_asistent_conversatie', JSON.stringify(faraMotor ? [{id: 'bun_venit', rol: 'assistant', text: t('ai.mesaj_bun_venit')}] : mesaje));
initializeazaSesiuneAI();
let apeluriLaCapat = 0;
const scrollOriginal = Element.prototype.scrollTo;
Element.prototype.scrollTo = function (...args) {if (this.classList.contains('flux-conversatie-ai')) apeluriLaCapat++; return scrollOriginal.apply(this, args);};
createRoot(document.getElementById('root')).render(<div className="aplicatie"><aside className="bara-laterala"><strong>Karios</strong><p>{t('ui.asistent_ai')}</p></aside><main className="continut"><header className="antet"><div><p className="fir-navigare">{t('ui.date_baza')} / {t('ui.asistent_ai')}</p><h1>{t('ui.asistent_ai')}</h1><p className="descriere">{t('ui.asistent_ai_descriere')}</p></div></header><AsistentAI date={date} onDeschidePontaj={(an, luna) => {lunaDeschisa = {an, luna};}}/></main></div>);
const asteapta = ms => new Promise(r => setTimeout(r, ms));
const cere = (conditie, mesaj) => {if (!conditie) throw Error(mesaj);};
const rect = e => {const r = e.getBoundingClientRect(); return {top: Math.round(r.top), bottom: Math.round(r.bottom), width: Math.round(r.width), height: Math.round(r.height)};};
const vizibil = e => {const r = e.getBoundingClientRect(); return r.top >= 0 && r.bottom <= innerHeight + 1 && r.width > 0;};
const metrici = {caz, limba, latime: innerWidth, inaltime: innerHeight};
async function verifica() {
    if (faraMotor) {
        for (let n = 0; n < 100 && !document.querySelector('.card-lipsa-fisiere-ai'); n++) await asteapta(50);
        cere(vizibil(document.querySelector('.formular-intrare-ai input')), 'Intrarea nu este vizibilă fără motor');
        cere(document.querySelectorAll('.navigare-lucru-ai button')[1].disabled, 'Plan gol accesibil fără intenție');
        cere(document.documentElement.scrollHeight <= innerHeight + 1, 'Pagina depășește fereastra fără motor');
        document.body.dataset.rezultat = 'trecut';
        return;
    }
    for (let n = 0; n < 100 && !document.querySelector('.plan-asistent-ai'); n++) await asteapta(50);
    await asteapta(300);
    const selector = document.querySelector('.navigare-lucru-ai select');
    selector.value = 'plan-vechi'; selector.dispatchEvent(new Event('change', {bubbles: true})); await asteapta(150);
    const corp = document.querySelector('.corp-plan-ai');
    corp.scrollTop = 120;
    const optiune = document.querySelector('.clarificare-plan-ai button');
    optiune.focus({preventScroll: true});
    const scroll = corp.scrollTop, pagina = window.scrollY, apeluri = apeluriLaCapat;
    optiune.click(); await asteapta(200);
    cere(document.activeElement === optiune, 'Focus pierdut la editare');
    cere(Math.abs(corp.scrollTop - scroll) <= 1 && window.scrollY === pagina, 'Poziție pierdută la editare');
    cere(apeluriLaCapat === apeluri && selector.value === 'plan-vechi', 'Editarea a activat derularea sau alt plan');
    const previzualizeaza = document.querySelector('[data-actiune="previzualizeaza"]');
    cere(vizibil(previzualizeaza), 'Previzualizarea nu este la vedere');
    metrici.editareStabila = true; metrici.formular = rect(document.querySelector('.plan-asistent-ai')); metrici.subsolInterpretare = rect(document.querySelector('.subsol-plan-ai'));
    if (caz === 'istoric') {
        document.querySelector('.navigare-lucru-ai button').click(); await asteapta(100);
        const flux = document.querySelector('.flux-conversatie-ai');
        flux.scrollTop = 0; flux.dispatchEvent(new Event('scroll')); await asteapta(50);
        const inainte = apeluriLaCapat;
        seteazaConversatie(prev => [...prev, {id: 'mesaj-la-inceput', rol: 'assistant', text: 'Răspuns fictiv nou.'}]); await asteapta(100);
        cere(apeluriLaCapat === inainte && flux.scrollTop === 0, 'Mesajul nou a întrerupt citirea istoricului');
        flux.scrollTop = flux.scrollHeight; flux.dispatchEvent(new Event('scroll')); await asteapta(50);
        seteazaConversatie(prev => [...prev, {id: 'mesaj-la-capat', rol: 'assistant', text: 'Alt răspuns fictiv nou.'}]); await asteapta(100);
        cere(apeluriLaCapat === inainte + 1, 'Mesajul nou nu urmărește capătul');
        metrici.urmarireConditionata = true;
    } else if (caz !== 'editare') {
        previzualizeaza.click(); await asteapta(200);
        const confirma = document.querySelector('[data-actiune="confirma"]');
        cere(document.querySelector('[data-etapa="verificare"]') && vizibil(confirma), 'Confirmarea nu este vizibilă în verificare');
        cere(document.activeElement === document.querySelector('.antet-plan-ai h3'), 'Previzualizarea nu primește focus');
        cere(confirma.disabled, 'Confirmare fără bifarea verificării');
        const corpVerificare = document.querySelector('.corp-plan-ai');
        corpVerificare.scrollTop = corpVerificare.scrollHeight; await asteapta(50);
        cere(vizibil(confirma), 'Confirmarea dispare când se derulează tabelul');
        metrici.confirmare = rect(confirma); metrici.corpVerificare = rect(corpVerificare);
        if (caz !== 'previzualizare') {
            document.querySelector('.subsol-plan-ai input[type="checkbox"]').click(); await asteapta(50);
            confirma.click(); await asteapta(200);
            cere(confirmari === 1, 'Număr greșit de confirmări');
            if (caz === 'eroare') {
                cere(document.querySelector('[data-etapa="interpretare"]') && document.querySelector('[role="alert"]')?.textContent.includes(t('ai.plan_expirat')), 'Plan expirat neafișat');
                cere(vizibil(document.querySelector('[data-actiune="previzualizeaza"]')), 'Reluarea nu este accesibilă');
                cere(!document.querySelector('[data-actiune="confirma"]'), 'Tokenul expirat încă poate fi confirmat');
            } else {
                cere(document.querySelector('[data-etapa="rezultat"]'), 'Lipsește rezultatul');
                const pastrat = document.querySelector('tr[data-stare="ocupata"]');
                const blocat = document.querySelector('tr[data-stare="blocata"]');
                const salvat = document.querySelector('tr[data-stare="inlocuire"]');
                cere(pastrat.cells[3].textContent === '6' && pastrat.cells[4].textContent === t('ai.plan_final_pastrat'), 'Celula păstrată este prezentată ca modificată');
                cere(blocat.cells[3].textContent === '—' && blocat.cells[4].textContent.includes(t('ai.plan_final_neaplicat')), 'Absența este prezentată ca aplicată');
                cere(salvat.cells[3].textContent === 'Z7' && salvat.cells[4].textContent === t('ai.plan_final_salvat'), 'Lipsește valoarea finală salvată');
                cere(vizibil(document.querySelector('.subsol-plan-ai [role="status"]')), 'Rezultatul nu este la vedere');
                document.querySelector('.subsol-plan-ai .actiuni-plan-ai button').click();
                cere(lunaDeschisa?.an === 2026 && lunaDeschisa.luna === 10, 'Lună greșită la deschiderea pontajului');
                metrici.valoriFinaleCorecte = true;
            }
        }
    }
    cere(document.documentElement.scrollWidth <= innerWidth + 1 && document.documentElement.scrollHeight <= innerHeight + 1, 'Derulare la nivelul paginii');
    document.body.dataset.rezultat = 'trecut';
    const rezultat = document.createElement('script'); rezultat.type = 'application/json'; rezultat.id = 'metrici-proba'; rezultat.textContent = JSON.stringify(metrici); document.body.append(rezultat);
}
verifica().catch(e => {document.body.dataset.rezultat = e.message; const p = document.createElement('pre'); p.textContent = e.stack; document.body.append(p);});
