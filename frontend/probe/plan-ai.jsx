import React, {useState} from 'react';
import {createRoot} from 'react-dom/client';
import i18next from 'i18next';
import {initReactI18next} from 'react-i18next';
import {PlanAsistentAI} from '../src/PlanAsistentAI';
import {ai} from '../wailsjs/go/models';
import ro from '../../internal/i18n/locales/ro.json';
import en from '../../internal/i18n/locales/en.json';
import '../src/style.css';

const limba = new URLSearchParams(location.search).get('limba') || 'ro';
const caz = new URLSearchParams(location.search).get('caz') || 'plan';
await i18next.use(initReactI18next).init({lng: limba, resources:{ro:{translation:ro},en:{translation:en}},interpolation:{escapeValue:false}});
const persoane = [{id:'p1',marca:'PROBA-01',nume:'Persoană Alfa'},{id:'p2',marca:'PROBA-02',nume:'Persoană Beta'},{id:'p3',marca:'PROBA-03',nume:'Persoană Gama'}];
const date = {angajati:persoane.map(p=>({...p,prenume:'',activ:true})),sectii:[],departamente:[],categorii_program:[]};
const program = {id:'z7',cod:'Z7',denumire:{ro:'Zi 07:00–14:00',en:'Day 07:00–14:00'},activ:true,ora_start:420,durata_minute:420};
let confirmari = 0;
let expirariRamase = caz === 'ciclu' ? 1 : 0;
window.go={main:{App:{
    ListaSabloaneAI:async()=>({sabloane:[]}),
    ListaSchimburi:async()=>({schimburi:[program]}),
    PrevizualizeazaPlanAI:async intentie=>({plan:{intentie,persoane,luni:['2026-10','2026-11'],noi:127,inlocuiri:caz==='ciclu'?2:0,identice:0,ocupate:1,blocate:1,aplicate:0,token:'proba',randuri:Array.from({length:129},(_,x)=>({angajat_id:persoane[x%3].id,marca:persoane[x%3].marca,nume:persoane[x%3].nume,data:`2026-${x<66?'10':'11'}-${String(x%21+1).padStart(2,'0')}`,inainte:x===0?'6':'',dupa:'Z7',stare:x===0?'ocupata':x===1?'blocata':'noua',motiv:x===1?'absenta':''}))}}),
    ConfirmaPlanAI:async(intentie,token)=>{confirmari++; if(token!=='proba')throw Error('token'); if(expirariRamase>0){expirariRamase--;return {eroare:{cod:'ai.plan_expirat'}};} const r=await window.go.main.App.PrevizualizeazaPlanAI(intentie);r.plan.aplicate=127;return r;},
}}};
function Proba() {
    const [intentie,setIntentie]=useState(new ai.IntentiePlan({selectie:{tip:'toti',valoare:'',marci:[],fara_pontaj:false},de_la:'2026-10-01',pana_la:'2026-11-30',cod_schimb:'Z7',filtru_zile:'',mod_existente:'',mesaj_original:'Toți angajații în octombrie și noiembrie lucrează de la 7 la 14',clarificari:['zile','existente']}));
    const [rezultat,setRezultat]=useState();const [ocupat,setOcupat]=useState(false);const [vizibil,setVizibil]=useState(true);
    window.__comutaPlan=()=>setVizibil(v=>!v);
    return <main style={{maxWidth:1040,margin:'20px auto',padding:16,height:'calc(100dvh - 40px)'}}>{vizibil&&<PlanAsistentAI idPlan="plan-proba" intentie={intentie} date={date} rezultat={rezultat} blocat={ocupat} onOcupat={setOcupat} onSchimba={setIntentie} onAplicat={setRezultat} onDeschidePontaj={()=>{}}/>}</main>;
}
createRoot(document.getElementById('root')).render(<Proba/>);
const asteapta=ms=>new Promise(r=>setTimeout(r,ms));
const t=k=>i18next.t(k);
const cere=(conditie,mesaj)=>{if(!conditie)throw Error(mesaj);};
const ciclu=()=>document.querySelector('.ciclu-plan-ai')?.dataset.ciclu;
// Seteaza valoarea unui input controlat de React si declanseaza evenimentul.
function scrie(el,valoare){Object.getOwnPropertyDescriptor(HTMLInputElement.prototype,'value').set.call(el,valoare);el.dispatchEvent(new Event('input',{bubbles:true}));}
// F14-T35: lipsuri, validare langa camp, rezumat, avertisment, expirare si revenire.
async function probaCiclu(){
    for(let x=0;x<100 && !document.querySelector('.clarificare-plan-ai');x++)await asteapta(50);
    cere(ciclu()==='editare','ciclu initial '+ciclu());
    const lipsuri=[...document.querySelectorAll('.lipsuri-plan-ai button')];
    cere(lipsuri.length===2,'lipsuri '+lipsuri.length);
    lipsuri[0].click();await asteapta(50);
    cere(document.activeElement?.id==='plan-ai-zile','lipsa nu duce la camp: '+document.activeElement?.id);
    scrie(document.getElementById('plan-ai-de-la'),'2026-12-01');await asteapta(100);
    const panaLa=document.getElementById('plan-ai-pana-la');
    cere(panaLa.getAttribute('aria-invalid')==='true'&&document.getElementById('plan-ai-eroare-pana_la')?.textContent===t('ai.validare_ordine'),'ordinea datelor nu e marcata langa camp');
    cere(document.getElementById('plan-ai-de-la').getAttribute('aria-invalid')!=='true','data de inceput marcata gresit');
    scrie(document.getElementById('plan-ai-de-la'),'2026-10-01');await asteapta(100);
    cere(!document.getElementById('plan-ai-eroare-pana_la'),'eroarea ramane dupa corectare');
    document.querySelectorAll('.clarificare-plan-ai')[0].querySelector('button').click(); await asteapta(80);
    document.querySelectorAll('.clarificare-plan-ai')[1].querySelector('button').click(); await asteapta(80);
    cere(!document.querySelector('.lipsuri-plan-ai'),'lipsuri ramase dupa clarificare');
    document.querySelector('[data-actiune="previzualizeaza"]').click();await asteapta(200);
    cere(ciclu()==='pregatit','ciclu dupa previzualizare '+ciclu());
    const luna=new Intl.DateTimeFormat(limba,{month:'long',year:'numeric',timeZone:'UTC'}).format(new Date('2026-10-01T12:00:00Z'));
    cere(document.querySelector('.rezumat-plan-ai')?.textContent.includes(luna),'rezumatul nu numeste luna '+luna);
    cere(document.querySelector('[data-contor="inlocuiri"] dd')?.textContent==='2','contor inlocuiri');
    cere(document.querySelector('.previzualizare-plan-ai .mesaj.avertisment')?.textContent.includes('2'),'avertisment inlocuire lipsa');
    // Revenire: zona parasita si redeschisa nu prezinta tokenul vechi drept confirmabil.
    window.__comutaPlan();await asteapta(80);window.__comutaPlan();await asteapta(150);
    cere(ciclu()==='expirat'&&!document.querySelector('[data-actiune="confirma"]')&&document.body.textContent.includes(t('ai.plan_revenire')),'revenirea nu cere o noua verificare');
    document.querySelector('[data-actiune="actualizeaza"]').click();await asteapta(200);
    cere(ciclu()==='pregatit','actualizarea dupa revenire');
    // Expirare la confirmare: contextul ramane si actualizarea primeste focus.
    const casete=document.querySelectorAll('.subsol-plan-ai .optiune-plan-ai input');casete[casete.length-1].click();await asteapta(80);
    document.querySelector('[data-actiune="confirma"]').click();await asteapta(250);
    cere(ciclu()==='expirat'&&document.activeElement?.dataset.actiune==='actualizeaza','expirarea nu duce la actualizare: '+ciclu()+' '+document.activeElement?.outerHTML.slice(0,80));
    cere(document.querySelector('.context-plan-ai')?.textContent.includes('Z7'),'contextul s-a pierdut la expirare');
    document.activeElement.click();await asteapta(200);
    cere(ciclu()==='pregatit'&&!document.querySelector('[data-actiune="confirma"]').disabled===false,'dupa actualizare confirmarea cere din nou verificarea');
    const c2=document.querySelectorAll('.subsol-plan-ai .optiune-plan-ai input');c2[c2.length-1].click();await asteapta(80);
    document.querySelector('[data-actiune="confirma"]').click();await asteapta(250);
    cere(ciclu()==='aplicat'&&confirmari===2,'aplicare: '+ciclu()+' '+confirmari);
    // Tabelul afiseaza datele formatate, nu ISO.
    cere(!/\b2026-10-\d\d\b/.test(document.querySelector('.tabel-plan-ai tbody')?.textContent||''),'date ISO in tabel');
    cere(document.documentElement.scrollWidth<=innerWidth,'depasire latime');
    document.body.dataset.rezultat='trecut';
}
async function proba(){
    if(caz==='ciclu')return probaCiclu();
    for(let x=0;x<100 && !document.querySelector('.clarificare-plan-ai');x++)await asteapta(50);
    const initial=document.querySelector('[data-actiune="previzualizeaza"]');
    if(!initial?.disabled)throw Error('previzualizare fără clarificări');
    document.querySelectorAll('.clarificare-plan-ai')[0].querySelector('button').click(); await asteapta(100);
    document.querySelectorAll('.clarificare-plan-ai')[1].querySelector('button').click(); await asteapta(100);
    document.querySelector('[data-actiune="previzualizeaza"]').click();await asteapta(200);
    const confirmare=[...document.querySelectorAll('button')].find(b=>b.textContent.includes('127'));
    if(!confirmare?.disabled)throw Error('confirmare fără verificare');
    const casete=document.querySelectorAll('.optiune-plan-ai input');casete[casete.length-1].click();await asteapta(100);
    confirmare.click();await asteapta(200);
    if(confirmari!==1 || !document.querySelector('.subsol-plan-ai [role="status"]'))throw Error('rezultat neverificat');
    if(document.documentElement.scrollWidth>innerWidth)throw Error('depășire lățime');
    document.body.dataset.rezultat='trecut';
}
proba().catch(e=>{document.body.dataset.rezultat=e.message;const p=document.createElement('pre');p.textContent=e.stack;document.body.append(p);});
