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
await i18next.use(initReactI18next).init({lng: limba, resources:{ro:{translation:ro},en:{translation:en}},interpolation:{escapeValue:false}});
const persoane = [{id:'p1',marca:'PROBA-01',nume:'Persoană Alfa'},{id:'p2',marca:'PROBA-02',nume:'Persoană Beta'},{id:'p3',marca:'PROBA-03',nume:'Persoană Gama'}];
const date = {angajati:persoane.map(p=>({...p,prenume:'',activ:true})),sectii:[],departamente:[],categorii_program:[]};
const program = {id:'z7',cod:'Z7',denumire:{ro:'Zi 07:00–14:00',en:'Day 07:00–14:00'},activ:true,ora_start:420,durata_minute:420};
let confirmari = 0;
window.go={main:{App:{
    ListaSchimburi:async()=>({schimburi:[program]}),
    PrevizualizeazaPlanAI:async intentie=>({plan:{intentie,persoane,luni:['2026-10','2026-11'],noi:127,inlocuiri:0,identice:0,ocupate:1,blocate:1,aplicate:0,token:'proba',randuri:Array.from({length:129},(_,x)=>({angajat_id:persoane[x%3].id,marca:persoane[x%3].marca,nume:persoane[x%3].nume,data:`2026-${x<66?'10':'11'}-${String(x%21+1).padStart(2,'0')}`,inainte:x===0?'6':'',dupa:'Z7',stare:x===0?'ocupata':x===1?'blocata':'noua',motiv:x===1?'absenta':''}))}}),
    ConfirmaPlanAI:async(intentie,token)=>{confirmari++; if(token!=='proba')throw Error('token'); const r=await window.go.main.App.PrevizualizeazaPlanAI(intentie);r.plan.aplicate=127;return r;},
}}};
function Proba() {
    const [intentie,setIntentie]=useState(new ai.IntentiePlan({selectie:{tip:'toti',valoare:'',marci:[],fara_pontaj:false},de_la:'2026-10-01',pana_la:'2026-11-30',cod_schimb:'Z7',filtru_zile:'',mod_existente:'',mesaj_original:'Toți angajații în octombrie și noiembrie lucrează de la 7 la 14',clarificari:['zile','existente']}));
    const [rezultat,setRezultat]=useState();const [ocupat,setOcupat]=useState(false);
    return <main style={{maxWidth:1040,margin:'20px auto',padding:16}}><PlanAsistentAI intentie={intentie} date={date} rezultat={rezultat} blocat={ocupat} onOcupat={setOcupat} onSchimba={setIntentie} onAplicat={setRezultat} onDeschidePontaj={()=>{}}/></main>;
}
createRoot(document.getElementById('root')).render(<Proba/>);
const asteapta=ms=>new Promise(r=>setTimeout(r,ms));
async function proba(){
    for(let x=0;x<100 && !document.querySelector('.clarificare-plan-ai');x++)await asteapta(50);
    const initial=document.querySelector('.plan-asistent-ai > fieldset > button');
    if(!initial?.disabled)throw Error('previzualizare fără clarificări');
    document.querySelectorAll('.clarificare-plan-ai')[0].querySelector('button').click(); await asteapta(100);
    document.querySelectorAll('.clarificare-plan-ai')[1].querySelector('button').click(); await asteapta(100);
    document.querySelector('.plan-asistent-ai > fieldset > button').click();await asteapta(200);
    const confirmare=[...document.querySelectorAll('button')].find(b=>b.textContent.includes('127'));
    if(!confirmare?.disabled)throw Error('confirmare fără verificare');
    const casete=document.querySelectorAll('.optiune-plan-ai input');casete[casete.length-1].click();await asteapta(100);
    confirmare.click();await asteapta(200);
    if(confirmari!==1 || !document.querySelector('.previzualizare-plan-ai [role="status"] strong'))throw Error('rezultat neverificat');
    if(document.documentElement.scrollWidth>innerWidth)throw Error('depășire lățime');
    document.body.dataset.rezultat='trecut';
}
proba().catch(e=>{document.body.dataset.rezultat=e.message;const p=document.createElement('pre');p.textContent=e.stack;document.body.append(p);});
