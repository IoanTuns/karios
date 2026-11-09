import assert from 'node:assert/strict';
// Simularea verifica interactiunea si pastrarea pozitiilor; calculul real este testat in Go.
export async function verificaProgrameMultiple({cdp,evalueaza,asteapta,apasa,completeaza,cataloage,limba,captura}){
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
 // F07-T42 (D097): celula planificata are fundalul programului, iar legenda arata programul.
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.grila-planificare td.celula-program'))"),'Celula planificata fara culoarea programului');
 const culoare=await evalueaza(`(()=>{const p=JSON.parse(localStorage.getItem('pontix-test'));const poz=p.planificare.find(x=>x.activ);const s=p.schimburi.find(x=>x.id===poz.schimb_id);const td=document.querySelector('.grila-planificare td.celula-program');return {asteptat:s.culoare_afisata,var:td.style.getPropertyValue('--culoare-program'),fundal:getComputedStyle(td).backgroundColor,legenda:Boolean(document.querySelector('.legenda-programe .mostra-program'))};})()`);
 assert.equal(culoare.var,culoare.asteptat,'Culoarea celulei nu e cea a programului');
 assert(culoare.legenda,'Legenda culorilor lipseste');
 assert.notEqual(culoare.fundal,'rgba(0, 0, 0, 0)','Fundalul celulei nu se aplica');
 await evalueaza("document.querySelector('.grila-planificare td.celula-program').scrollIntoView({block:'center'})");
 await captura?.('t42-01-planificare-culori.png');
 // Grila Pontaj: aceeasi culoare pe celula cu programul publicat sau introdus.
 await apasa(t('ui.pontaj'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.grila-pontaj'))"),'Pontaj lipsa');
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.grila-pontaj td.col-zi input'))"),'Celule Pontaj lipsa');
 // F08-T16 (I124): programul S7 intr-o zi lucratoare a fiecarui angajat; celula ia culoarea lui, iar legenda il arata.
 const zi=await evalueaza(`(()=>{const td=[...document.querySelectorAll('.grila-pontaj td.col-zi')].find(td=>!td.className.match(/zi-weekend|zi-sarbatoare/));return td.querySelector('input').getAttribute('aria-label').match(/[0-9]{4}-[0-9]{2}-[0-9]{2}/)[0];})()`);
 await evalueaza(`(()=>{const s=JSON.parse(localStorage.getItem('pontix-test'));s.ciorne=(s.ciorne||[]).filter(c=>c.data!==${JSON.stringify(zi)});for(const a of s.angajati.filter(a=>a.activ!==false))s.ciorne.push({angajat_id:a.id,data:${JSON.stringify(zi)},valoare:'S7',revizie:'1'});localStorage.setItem('pontix-test',JSON.stringify(s));})()`);
 await cdp('Page.reload');
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.actiuni-antet'))"),'Aplicatie lipsa dupa reincarcare');
 await apasa(t('ui.pontaj'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.grila-pontaj td.celula-program'))"),'Celula Pontaj fara culoarea programului');
 const pontaj=await evalueaza(`(()=>{const p=JSON.parse(localStorage.getItem('pontix-test'));const s=p.schimburi.find(x=>x.cod==='S7');const td=document.querySelector('.grila-pontaj td.celula-program');const l=[...document.querySelectorAll('.legenda-culori-pontaj span')].map(e=>({cod:e.textContent.trim(),fundal:e.querySelector('i')?.style.background}));return {asteptat:s.culoare_afisata,var:td.style.getPropertyValue('--culoare-program'),fundal:getComputedStyle(td).backgroundColor,legenda:l};})()`);
 assert.equal(pontaj.var,pontaj.asteptat,'Culoarea celulei Pontaj nu e cea a programului');
 assert.notEqual(pontaj.fundal,'rgba(0, 0, 0, 0)','Fundalul celulei Pontaj nu se aplica');
 assert(pontaj.legenda.some(x=>x.cod==='S7'&&x.fundal),'Legenda Pontajului nu arata S7: '+JSON.stringify(pontaj.legenda));
 await evalueaza("document.querySelector('.grila-pontaj td.celula-program').scrollIntoView({block:'center',inline:'center'})");
 await captura?.('t42-02-pontaj.png');
 await evalueaza("document.querySelector('.legenda-grila').scrollIntoView({block:'start'})");
 await captura?.('t16-01-pontaj-legenda.png');
}
