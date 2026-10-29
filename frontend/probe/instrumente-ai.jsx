// F14-T25: probă a componentelor reale, cu date fictive și punte Wails simulată.
import React from 'react';
import {createRoot} from 'react-dom/client';
import i18next from 'i18next';
import {initReactI18next} from 'react-i18next';
import {AsistentAI} from '../src/AsistentAI';
import {initializeazaSesiuneAI} from '../src/asistentSesiune';
import ro from '../../internal/i18n/locales/ro.json';
import en from '../../internal/i18n/locales/en.json';
import '../src/style.css';

const param = new URLSearchParams(location.search), caz = param.get('caz') || 'plan', limba = param.get('limba') || 'ro';
await i18next.use(initReactI18next).init({lng:limba,resources:{ro:{translation:ro},en:{translation:en}},interpolation:{escapeValue:false}});
const t = (k,p) => i18next.t(`ai.${k}`,p);
const intentie = {an_de_confirmat:true,selectie:{tip:'toti',valoare:'',marci:[],fara_pontaj:false},de_la:'2026-10-30',pana_la:'2026-11-02',cod_schimb:'Z7',filtru_zile:'lucratoare',mod_existente:'pastreaza',mesaj_original:'Toți pe Z7, anul neprecizat',clarificari:['an']};
const oameni = [{id:'a',marca:'FICTIV-1',nume:'Exemplu',prenume:'Ana',activ:true},{id:'b',marca:'FICTIV-2',nume:'Probă',prenume:'Ion',activ:true},{id:'c',marca:'FICTIV-3',nume:'Inactiv',prenume:'Test',activ:false}];
const randuri = [{angajat_id:'a',marca:'FICTIV-1',nume:'Exemplu Ana',data:'2026-10-30',inainte:'6',dupa:'Z7',stare:'ocupata',motiv:''},{angajat_id:'a',marca:'FICTIV-1',nume:'Exemplu Ana',data:'2026-11-02',inainte:'',dupa:'Z7',stare:'noua',motiv:''},{angajat_id:'b',marca:'FICTIV-2',nume:'Probă Ion',data:'2026-10-30',inainte:'',dupa:'Z7',stare:'blocata',motiv:'absenta'},{angajat_id:'b',marca:'FICTIV-2',nume:'Probă Ion',data:'2026-11-02',inainte:'Z7',dupa:'Z7',stare:'identica',motiv:''}];
const sumar = {selectati:2,afectati:0,zile:1,noi:0,inlocuiri:0,pastrate:1,blocate:1,identice:0};
const plan = {intentie,persoane:oameni.filter(a=>a.activ).map(a=>({...a,nume:`${a.nume} ${a.prenume}`})),luni:['2026-10','2026-11'],randuri,token:'token-verificat',noi:1,inlocuiri:0,ocupate:1,blocate:1,identice:1,aplicate:0,selectie_explicata:oameni.map(a=>({...a,nume:`${a.nume} ${a.prenume}`,inclus:a.activ,motiv:a.activ?'toti':'inactiv'})),sumar_luni:[{...sumar,cheie:'2026-10',nume:'2026-10'},{...sumar,cheie:'2026-11',nume:'2026-11',afectati:1,noi:1,pastrate:0,blocate:0,identice:1}],sumar_persoane:[]};
let sabloane = [], exporturi = [], confirmari = 0, anulari = 0;
window.runtime = {EventsOnMultiple:()=>()=>{}};
window.go = {main:{App:{
    StareAI:async()=>({stare:{instalat:true,activ:true,port:5000,model_activ:'fictiv-3b.gguf',context_configurat:32768}}),EsteAutoPornireAI:async()=>true,SarcinaAIActiva:async()=>'',
    ListaSchimburi:async()=>({schimburi:[{id:'z7',cod:'Z7',denumire:{ro:'Zi',en:'Day'},activ:true}]}),
    ListaSabloaneAI:async()=>({sabloane}),SalveazaSablonAI:async(nume,i)=>{sabloane=[{id:'sablon1',nume,intentie:i}];return {sabloane};},StergeSablonAI:async()=>{sabloane=[];return {sabloane};},
    PrevizualizeazaPlanAI:async i=>{if(i.an_de_confirmat)throw Error('An neconfirmat');return {plan:{...plan,intentie:i}};},
    ConfirmaPlanAI:async(i,token)=>{if(token!=='token-verificat')throw Error('Token invalid');confirmari++;return {plan:{...plan,intentie:i,operatie_id:'op1',aplicate:1}};},
    ExportaPlanAICSV:async(i,token,id,l)=>{exporturi.push({i,token,id,l});return {cale:'date-fictive/plan.csv'};},
    VerificaPontajAI:async()=>({verificare:{celule:4,semnale:[{angajat_id:'a',nume:'Exemplu Ana',marca:'FICTIV-1',data:'2026-10-30',valoare:'6',motiv:'ore_manuale'},{angajat_id:'a',nume:'Exemplu Ana',marca:'FICTIV-1',data:'2026-10',valoare:'',motiv:'regula',regula_cod:'TEST',regula_nume:'Limită fictivă',severitate:'AVERTISMENT'}]}}),
    PrevizualizeazaAnulareAI:async()=>({anulare:{operatie_id:'op1',creat_la:'2026-10-29',token:'undo-token',conflicte:caz==='conflict'?1:0,anulata:false,randuri:[{angajat_id:'a',nume:'Exemplu Ana',data:'2026-11-02',acum:'Z7',restaurat:'',plan_acum:['G24'],plan_restaurat:['N12'],motiv:caz==='conflict'?'modificata':''}]}}),
    AnuleazaUltimaOperatieAI:async token=>{if(token!=='undo-token')throw Error('Token undo invalid');anulari++;if(caz==='expirat')return {eroare:{cod:'ai.anulare_conflict'}};return {anulare:{operatie_id:'op1',randuri:[],anulata:true}};},
}}};
sessionStorage.setItem('pontix_asistent_conversatie',JSON.stringify([{id:'plan',rol:'assistant',text:'',intentie}]));initializeazaSesiuneAI();
createRoot(document.getElementById('root')).render(<div className="aplicatie"><aside className="bara-laterala">Pontix</aside><main className="continut"><header className="antet"><h1>{t('titlu')}</h1></header><AsistentAI date={{angajati:oameni,sectii:[],departamente:[],categorii_program:[]}}/></main></div>);
const asteapta = ms=>new Promise(r=>setTimeout(r,ms));
const cere = (c,m)=>{if(!c)throw Error(m);};
const buton = (cheie,rad=document)=>[...rad.querySelectorAll('button')].find(b=>b.textContent.trim()===t(cheie));
async function apasa(cheie,rad=document){const b=buton(cheie,rad);cere(b&&!b.disabled,`Buton indisponibil ${cheie}`);b.click();await asteapta(120);}
async function seteaza(el,v){cere(el,'Câmp lipsă');Object.getOwnPropertyDescriptor(el instanceof HTMLSelectElement?HTMLSelectElement.prototype:HTMLInputElement.prototype,'value').set.call(el,v);el.dispatchEvent(new Event(el instanceof HTMLSelectElement?'change':'input',{bubbles:true}));await asteapta(100);}
const laVedere = e=>{const r=e.getBoundingClientRect();return r.top>=0&&r.bottom<=innerHeight+1&&r.width>0;};
async function verifica(){
    await asteapta(500);
    cere(document.querySelector('[data-actiune="previzualizeaza"]').disabled,'Anul implicit nu cere confirmare');
    const confirmAn=document.querySelector('.corp-plan-ai .optiune-plan-ai input');confirmAn.click();await asteapta(120);
    cere(!document.querySelector('[data-actiune="previzualizeaza"]').disabled,'Anul confirmat nu deblochează planul');
    if(caz==='clarificari'){
        const inceput=document.querySelector('input[type="date"]');await seteaza(inceput,'2027-01-01');cere(document.querySelector('[data-actiune="previzualizeaza"]').disabled,'Perioadă inversată acceptată');await seteaza(inceput,'2026-10-30');
        await apasa('verificare_porneste');cere(document.querySelector('.verificare-pontaj-ai').textContent.includes(t('verificare_ore_manuale')),'Semnalul verificării lipsește');
        cere(document.querySelector('.verificare-pontaj-ai').textContent.includes(t('verificare_regula',{cod:'TEST',nume:'Limită fictivă',severitate:t('severitate_AVERTISMENT')})),'Regula configurată nu este tradusă');
        document.querySelector('.optiuni-motor-ai').open=true;await asteapta(50);cere(document.querySelector('.meniu-motor-ai').textContent.includes('fictiv-3b.gguf')&&document.querySelector('.meniu-motor-ai').textContent.includes('32768'),'Metadate model lipsă');document.querySelector('.optiuni-motor-ai').open=false;
    } else if(caz==='sabloane'){
        const det=document.querySelector('.sabloane-ai');det.open=true;
        await seteaza(det.querySelector('input'),'Program secție');await apasa('sablon_salveaza',det);cere(sabloane.length===1,'Șablon nesalvat');
        await seteaza(det.querySelector('select'),'sablon1');await apasa('sablon_incarca',det);cere(document.querySelector('[data-actiune="previzualizeaza"]').disabled,'Șablonul nu cere revizuirea perioadei');
        window.confirm=()=>true;await apasa('sablon_sterge',det);cere(sabloane.length===0&&det.textContent.includes(t('sablon_sters')),'Ștergerea șablonului nu a fost actualizată');
        const grup=document.querySelector('.campuri-plan-ai select');await seteaza(grup,'persoane');const lista=document.querySelector('.persoane-plan-ai');await seteaza(lista.querySelector('input:not([type="checkbox"])'),'Proba');cere(lista.querySelectorAll('input[type="checkbox"]').length===1,'Căutarea nominală nu filtrează');
        lista.querySelector('button').click();await asteapta(100);cere(lista.querySelector('input[type="checkbox"]').checked,'Selecția rezultatelor nu funcționează');
    } else if(caz==='anulare'||caz==='conflict'||caz==='expirat'){
        await apasa('instrumente');await apasa('anulare_pregateste');const p=document.querySelector('.instrumente-asistent-ai');
        const confirma=buton('anulare_confirma',p);cere(confirma.disabled&&laVedere(confirma),'Anulare fără confirmare sau în afara ferestrei');
        cere(p.textContent.includes(t('anulare_plan_separat',{coduri:'G24'}))&&p.textContent.includes(t('anulare_plan_separat',{coduri:'N12'})),'Planificarea separată nu este prezentată la anulare');
        const check=p.querySelector('input[type="checkbox"]');
        if(caz==='conflict'){cere(check.disabled&&p.textContent.includes(t('anulare_motiv_modificata')),'Conflict neprotejat');cere(anulari===0,'Anulare cu conflict');}
        else{check.click();await asteapta(80);await apasa('anulare_confirma',p);if(caz==='expirat'){cere(anulari===1&&p.textContent.includes(t('anulare_conflict'))&&!buton('anulare_confirma',p),'Anularea expirată încă poate fi confirmată');}else{cere(anulari===1&&p.textContent.includes(t('anulare_reusita')),'Anulare neconfirmată');}}
    } else {
        document.querySelector('[data-actiune="previzualizeaza"]').click();await asteapta(150);
        cere(document.querySelector('.sumar-plan-ai').textContent.includes(t('sumar_detalii',sumar)),'Sumar selectați/afectați greșit');
        const explicatii=document.querySelector('.explicatii-selectie-ai');explicatii.open=true;await seteaza(explicatii.querySelector('select'),'exclus');cere(explicatii.querySelectorAll('li').length===1&&explicatii.textContent.includes(t('selectie_motiv_inactiv')),'Explicație excludere lipsă');
        const filtre=document.querySelector('.previzualizare-plan-ai > .filtre-detalii-ai');await seteaza(filtre.querySelector('input'),'Proba');await seteaza(filtre.querySelector('select'),'blocata');cere(document.querySelectorAll('tbody tr').length===1&&document.querySelector('tbody').textContent.includes(t('plan_motiv_absenta')),'Filtre incorecte');
        await apasa('export_titlu');cere(exporturi.length===1&&exporturi[0].token==='token-verificat'&&!exporturi[0].id,'Export propunere neverificat');
        document.querySelector('.subsol-plan-ai input').click();await asteapta(80);document.querySelector('[data-actiune="confirma"]').click();await asteapta(150);cere(confirmari===1&&document.querySelector('[data-etapa="rezultat"]'),'Aplicare lipsă');
        await apasa('export_titlu');cere(exporturi[1].id==='op1','Export rezultat fără operația persistentă');
        cere(laVedere(document.querySelector('.subsol-plan-ai [role="status"]')),'Rezultat în afara ferestrei');
    }
    cere(document.documentElement.scrollWidth<=innerWidth+1&&document.documentElement.scrollHeight<=innerHeight+1,'Pagina depășește fereastra');
    cere(!document.body.textContent.match(/ai\.(?:sumar|selectie|sablon|anulare|verificare|export)_/),'Chei netraduse în interfață');
    document.body.dataset.rezultat='trecut';
}
verifica().catch(e=>{document.body.dataset.rezultat=e.message;const p=document.createElement('pre');p.textContent=e.stack;document.body.append(p);});
