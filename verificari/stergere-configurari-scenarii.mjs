import assert from 'node:assert/strict';

export function simulareStergereConfigurari(){
 const api=window.go.main.App,citeste=api.DateInitiale,stareCont=api.StareCont;
 api.ListaBackupuri=async()=>({backupuri:[],programata:''});
 const grupuri=[];
 api.StareCont=async()=>{const r=await stareCont();r.cont.rol=localStorage.getItem('rol-stergere-test')||'ADMINISTRATOR';return r;};
 api.DateInitiale=async()=>{const r=await citeste();r.date.grupuri_profesionale=structuredClone(grupuri);return r;};
 api.SalveazaElementNomenclator=async(tip,e)=>{const i=grupuri.findIndex(g=>g.id===e.id);const nou={...structuredClone(e),id:e.id||'grup-de-sters',revizie:String(Date.now())};if(i>=0)grupuri[i]=nou;else grupuri.push(nou);return api.DateInitiale();};
 window.__stergeri=0;
 api.StergeConfigurare=async(tip,id,rev)=>{
  window.__stergeri++;
  if(localStorage.getItem('rol-stergere-test')==='OPERATOR')return {eroare:{cod:'cont.doar_administrator'}};
  if(window.__refuzaStergere)return {eroare:{cod:'stergere.folosit',param:{sursa:'functii'}}};
  const i=grupuri.findIndex(g=>g.id===id&&g.revizie===rev&&!g.activ);
  if(i<0)return {eroare:{cod:'stergere.conflict'}};
  grupuri.splice(i,1);return api.DateInitiale();
 };
}

export async function verificaStergereConfigurari({cdp,evalueaza,asteapta,apasa,completeaza,cataloage,limba,captura}){
 const t=k=>cataloage[limba][k];
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.ajutor-antet'))"),'Aplicatie lipsa');
 await apasa(t('ui.setari'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('[data-nomenclator=grupuri]'))"),'Editor lipsa');
 await evalueaza("document.querySelector('[data-nomenclator=grupuri] > button').click()");
 await completeaza('input[name=nomenclator_ro]','Configurare temporara');
 await completeaza('input[name=nomenclator_en]','Temporary configuration');
 await evalueaza("document.querySelector('[data-nomenclator=grupuri] form input[type=checkbox]').click()");
 await evalueaza("document.querySelector('[data-nomenclator=grupuri] form .buton.primar').click()");
 await asteapta(()=>evalueaza("!document.querySelector('[data-nomenclator=grupuri] form')"),'Salvare lipsa');
 await completeaza('select[name=filtru_stare_grupuri]','INACTIVE');
 await asteapta(()=>evalueaza("Boolean(document.querySelector('[data-nomenclator=grupuri] [data-stergere-admin]'))"),'Buton admin lipsa');
 await evalueaza("document.querySelector('[data-nomenclator=grupuri] details.actiune-stergere-admin').open=true");
 await evalueaza('window.confirm=()=>false');
 await evalueaza("document.querySelector('[data-nomenclator=grupuri] [data-stergere-admin]').click()");
 assert.equal(await evalueaza('window.__stergeri'),0,'Renuntarea a trimis stergere');
 await evalueaza('window.confirm=()=>true;window.__refuzaStergere=true');
 await evalueaza("document.querySelector('[data-nomenclator=grupuri] [data-stergere-admin]').click()");
 await asteapta(()=>evalueaza("Boolean(document.querySelector('[data-nomenclator=grupuri] [role=alert]'))"),'Utilizarea nu este explicata');
 await captura('stergere-refuzata-'+limba+'.png');
 await evalueaza('window.__refuzaStergere=false');
 await evalueaza("document.querySelector('[data-nomenclator=grupuri] [data-stergere-admin]').click()");
 await asteapta(()=>evalueaza("document.querySelectorAll('[data-nomenclator=grupuri] li').length===0"),'Inregistrarea nu a disparut');
 await evalueaza("localStorage.setItem('rol-stergere-test','OPERATOR')");
 await cdp('Page.reload');
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.ajutor-antet'))"),'Reincarcare lipsa');
 await apasa(t('ui.setari'));
 assert.equal(await evalueaza("document.querySelectorAll('[data-stergere-admin]').length"),0,'Operatorul vede stergerea');
 assert.equal(await evalueaza("(async()=> (await window.go.main.App.StergeConfigurare('grupuri','x','x')).eroare.cod)()"),'cont.doar_administrator','Apel direct operator permis');
}
