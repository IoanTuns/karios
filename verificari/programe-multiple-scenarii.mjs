import assert from 'node:assert/strict';
// Simularea verifica interactiunea si pastrarea pozitiilor; calculul real este testat in Go.
export async function verificaProgrameMultiple({evalueaza,asteapta,apasa,completeaza,cataloage,limba}){
 const t=k=>cataloage[limba][k];
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.actiuni-antet'))"),'Aplicatie lipsa');
 await apasa(t('ui.planificare'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.modul-planificare'))"),'Planificare lipsa');
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.buton-sumar-angajat'))"),'Angajati lipsa');
 await evalueaza("document.querySelector('.buton-sumar-angajat').click()");
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.panou-sumar-angajat'))"),'Sumar lipsa');
 const adauga=async()=>evalueaza(`(()=>{const r=document.querySelector('.tabel-sumar-angajat tbody tr');const b=[...r.querySelectorAll('button')].find(b=>b.textContent.trim()===${JSON.stringify(t('planificare.adauga_program_zi'))});b.click();})()`);
 await adauga();await asteapta(()=>evalueaza("document.querySelectorAll('.programe-zi-item').length===1"),'Primul program lipsa');
 const ids=await evalueaza("JSON.parse(localStorage.getItem('pontix-test')).schimburi.filter(s=>['S7','IN9'].includes(s.cod)).map(s=>({id:s.id,cod:s.cod}))");
 await completeaza('.programe-zi-item select',ids.find(s=>s.cod==='S7').id);
 await adauga();await asteapta(()=>evalueaza("document.querySelectorAll('.programe-zi-item').length===2"),'Al doilea inlocuieste primul');
 const gol=await evalueaza("[...document.querySelectorAll('.programe-zi-item select')].find(e=>!e.value).getAttribute('aria-label')");
 await completeaza('select[aria-label='+JSON.stringify(gol)+']',ids.find(s=>s.cod==='IN9').id);
 assert(await evalueaza("document.querySelectorAll('.programe-zi-item').length===2"),'Pozitii pierdute');
 await evalueaza("document.querySelector('.antet-sumar-angajat button').click()");
 await asteapta(()=>evalueaza("!document.querySelector('dialog[open]')"),'Sumar nesalvat');
 assert(await evalueaza("JSON.parse(localStorage.getItem('pontix-test')).planificare.filter(p=>p.activ).length===2"),'Stocare simulata pierde pozitie');
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.buton-sumar-angajat'))"),'Grila lipsa');
 await evalueaza("document.querySelector('.buton-sumar-angajat').click()");
 await asteapta(()=>evalueaza("document.querySelectorAll('.programe-zi-item').length===2"),'Redeschidere pierde pozitie');
 await evalueaza("document.querySelector('.programe-zi-item button').click()");
 await asteapta(()=>evalueaza("document.querySelectorAll('.programe-zi-item').length===1"),'Stergere neindividuala');
 await evalueaza("document.querySelector('.antet-sumar-angajat button').click()");
 await asteapta(()=>evalueaza("!document.querySelector('dialog[open]')"),'Anulare nesalvata');
 assert(await evalueaza("JSON.parse(localStorage.getItem('pontix-test')).planificare.filter(p=>p.activ).length===1"),'Anularea sterge intregul set');
}
