import assert from 'node:assert/strict';

export function simulareNomenclatoare(){
 const api=window.go.main.App;
 api.ListaBackupuri=async()=>({backupuri:[],programata:''});
 const liste={categorii:[{id:'LAB',denumire:{ro:'Laborator',en:'Laboratory'},activ:true,revizie:'1'},{id:'CUSTOM',denumire:{ro:'Personalizate',en:'Custom'},activ:true,revizie:'1'}],functii:[]};
 const citeste=api.DateInitiale;
 api.DateInitiale=async()=>{const r=await citeste();r.date.categorii_schimburi=structuredClone(liste.categorii);r.date.functii=structuredClone(liste.functii);return r};
 api.SalveazaElementNomenclator=async(tip,e)=>{const lista=liste[tip],index=lista.findIndex(x=>x.id===e.id);const nou={...structuredClone(e),id:e.id||crypto.randomUUID().toUpperCase(),revizie:String(Date.now())};if(index<0)lista.push(nou);else lista[index]=nou;return api.DateInitiale()};
 const salveaza=api.SalveazaAngajat;
 api.SalveazaAngajat=async a=>{a.functie=liste.functii.find(f=>f.id===a.functie_id)?.denumire.ro||a.functie;const r=await salveaza(a);if(r.eroare)return r;return api.DateInitiale()};
}

export async function verificaNomenclatoare({cdp,evalueaza,asteapta,apasa,completeaza,cataloage,limba,captura}){
 const t=k=>cataloage[limba][k];
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.ajutor-antet'))"),'Aplicatie neincarcata');
 async function adauga(tip,ro,en){
  await asteapta(()=>evalueaza(`Boolean(document.querySelector('[data-nomenclator="${tip}"] > button'))`),'Editor lipsa '+tip);
  await evalueaza(`document.querySelector('[data-nomenclator="${tip}"] > button').click()`);
  await completeaza(`input[name="nomenclator_ro"]`,ro);await completeaza(`input[name="nomenclator_en"]`,en);
  await evalueaza(`document.querySelector('[data-nomenclator="${tip}"] form .buton.primar').click()`);
  await asteapta(()=>evalueaza(`!document.querySelector('[data-nomenclator="${tip}"] form')`),'Lista nu s-a salvat');
 }
 await apasa(t('ui.setari'));await adauga('functii','Medic fictiv UI','Test physician UI');
 await captura('functii-'+limba+'.png');
 await apasa(t('ui.angajati'));await apasa(t('ui.adauga_angajat'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('dialog[open] input'))"),'Formular lipsa');
 await completeaza('dialog .campuri label:nth-child(1) input','Exemplu');await completeaza('dialog .campuri label:nth-child(2) input','Fictiv');await completeaza('dialog .campuri label:nth-child(3) input','FUNC-UI');
 const functie=await evalueaza("[...document.querySelectorAll('select[name=angajat_functie] option')].find(o=>o.textContent.includes('UI')).value");
 await completeaza('select[name="angajat_functie"]',functie);await apasa(t('ui.salveaza'));
 await asteapta(()=>evalueaza("!document.querySelector('dialog[open]')"),'Angajat nesalvat');
 await apasa(t('ui.pontaj'));await asteapta(()=>evalueaza("Boolean(document.querySelector('.grila-pontaj tbody'))"),'Pontaj lipsa');
 assert(await evalueaza("[...document.querySelectorAll('.grila-pontaj tbody tr')].some(r=>r.textContent.includes('FUNC-UI')&&r.querySelector('.fix-nume').textContent.includes("+JSON.stringify(limba==='en'?'Test physician UI':'Medic fictiv UI')+"))"),'Functia lipseste din randul angajatului');
 await evalueaza("[...document.querySelectorAll('.grila-pontaj tbody tr')].find(r=>r.textContent.includes('FUNC-UI')).scrollIntoView({block:'center'})");
 await captura('functie-pontaj-'+limba+'.png');
 await apasa(t('ui.setari'));await apasa(t('ui.schimburi'));await adauga('categorii','Ambulatoriu UI','Outpatient UI');
 await apasa(t('schimburi.adauga'));await asteapta(()=>evalueaza("Boolean(document.querySelector('dialog[open]'))"),'Formular program lipsa');
 await completeaza('dialog .campuri label:nth-child(1) input','UI7');await completeaza('dialog .campuri label:nth-child(3) input','Program UI');await completeaza('dialog .campuri label:nth-child(4) input','UI schedule');await completeaza('dialog input[type=time]','07:00');await completeaza('dialog input[type=number]','420');
 const cat=await evalueaza("[...document.querySelectorAll('dialog select[name=categorie] option')].find(o=>o.textContent.includes('UI')).value");await completeaza('dialog select[name="categorie"]',cat);await apasa(t('ui.salveaza'));
 await asteapta(()=>evalueaza("!document.querySelector('dialog[open]')"),'Program nesalvat');
 assert(await evalueaza("[...document.querySelectorAll('.tabel-schimburi tbody tr')].some(r=>r.textContent.includes('UI7')&&r.querySelector('td:nth-child(3)').textContent.includes("+JSON.stringify(limba==='en'?'Outpatient UI':'Ambulatoriu UI')+"))"),'Categoria nu este asociata programului');
 await captura('categorii-schimburi-'+limba+'.png');
 await apasa(t('ui.angajati'));await apasa(t('ui.adauga_angajat'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('select[name=angajat_principal]'))"),'Selector principal lipsa');
 await completeaza('dialog .campuri label:nth-child(1) input','Exemplu');await completeaza('dialog .campuri label:nth-child(2) input','Fictiv');await completeaza('dialog .campuri label:nth-child(3) input','GARDA-UI');
 const principal=await evalueaza("[...document.querySelectorAll('select[name=angajat_principal] option')].find(o=>o.textContent.includes('FUNC-UI')).value");
 await completeaza('select[name=angajat_principal]',principal);await captura('marca-garda-formular-'+limba+'.png');await apasa(t('ui.salveaza'));
 await asteapta(()=>evalueaza("!document.querySelector('dialog[open]')"),'Marca garda nesalvata');
 assert(await evalueaza("[...document.querySelectorAll('tr')].some(r=>r.textContent.includes('GARDA-UI')&&r.textContent.includes('FUNC-UI')&&r.textContent.includes("+JSON.stringify(t('rapoarte.subtip_garda'))+"))"),'Subtipul si marca principala lipsesc din lista');
 await captura('marca-garda-lista-'+limba+'.png');

}
