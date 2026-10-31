import assert from 'node:assert/strict';

export async function verificaAbsentePontaj({evalueaza,asteapta,apasa,completeaza,cataloage,limba,captura}) {
 const t=k=>cataloage[limba][k];
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.actiuni-antet'))"),'Aplicatie lipsa');
 await apasa(t('ui.pontaj'));
 const celula=i=>`input[data-rand="0"][data-zi="${i}"]`;
 await asteapta(()=>evalueaza(`Boolean(document.querySelector('${celula(0)}'))`),'Grila lipsa');
 assert.equal(await evalueaza("document.querySelector('select[name=filtru_stare_angajati_pontaj]').value"),'ACTIVE','Pontaj nu porneste pe activi');
 await completeaza('select[name=filtru_stare_angajati_pontaj]','INACTIVE');
 assert.equal(await evalueaza("document.querySelectorAll('.grila-pontaj tbody tr').length"),await evalueaza("JSON.parse(localStorage.getItem('pontix-test')).angajati.filter(a=>!a.activ).length"),'Filtrul Pontaj amesteca starile');
 await completeaza('select[name=filtru_stare_angajati_pontaj]','ACTIVE');
 await asteapta(()=>evalueaza(`Boolean(document.querySelector('${celula(0)}'))`),'Revenirea la activi nu reface grila');

 await evalueaza(`(()=>{const e=document.querySelector('${celula(0)}');e.dispatchEvent(new MouseEvent('mousedown',{bubbles:true,button:0}));e.focus()})()`);
 await completeaza('select[aria-label="'+t('pontaj.program_completare')+'"]','Z7');await apasa(t('pontaj.aplica_program'));
 await apasa(t('pontaj.salveaza'));
 await asteapta(()=>evalueaza("JSON.parse(localStorage.getItem('pontix-test')).ciorne.some(c=>c.valoare==='Z7')"),'Program nesalvat');
 await asteapta(()=>evalueaza(`!document.querySelector('${celula(0)}').disabled && !document.querySelector('${celula(0)}').closest('td').classList.contains('celula-modificata')`),'Salvare inca ocupata');
 await evalueaza(`(()=>{const e=document.querySelector('${celula(1)}');e.dispatchEvent(new MouseEvent('mousedown',{bubbles:true,button:0}));e.focus()})()`);
 await asteapta(()=>evalueaza(`document.querySelector('${celula(1)}').closest('td').classList.contains('celula-selectata')`),'Selectia zilei a doua lipsa');
 await apasa(t('pontaj.aplica_program')); 
 await asteapta(()=>evalueaza(`document.querySelector('${celula(1)}').closest('td').classList.contains('celula-modificata')`),'Ziua din afara selectiei nu a fost modificata');
 await evalueaza(`(()=>{const e=document.querySelector('${celula(0)}');e.dispatchEvent(new MouseEvent('mousedown',{bubbles:true,button:0}));e.focus()})()`);
 await asteapta(()=>evalueaza(`document.querySelector('${celula(0)}').closest('td').classList.contains('celula-selectata')`),'Selectia primei zile lipsa');
 await apasa(t('pontaj.absenta_titlu'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('dialog select[name=tip]'))"),'Formular absenta lipsa');
 await apasa(t('ui.anuleaza'));
 await asteapta(()=>evalueaza(`document.querySelector('${celula(0)}').value==='7'`),'Anularea a modificat programul');
 await apasa(t('pontaj.absenta_titlu'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('dialog select[name=tip]'))"),'Formular lipsa');
 await completeaza('dialog select[name=tip]','CO');await completeaza('dialog select[name=stare_aprobare]','APROBATA');
 await completeaza('dialog input[name=nr_inregistrare]','TEST-CO');
 await evalueaza('window.__refuzaConversie=true');await apasa(t('ui.salveaza'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('dialog .mesaj.eroare'))"),'Eroare nesemnalata');
 await asteapta(()=>evalueaza(`document.querySelector('${celula(0)}').value==='7'`),'Eroarea a modificat programul');
 await evalueaza('window.__refuzaConversie=false');await captura('concediu-din-pontaj-'+limba+'.png');await apasa(t('ui.salveaza'));
 await asteapta(()=>evalueaza(`document.querySelector('${celula(0)}').value==='CO'`),'CO nu apare in Pontaj');
 await asteapta(()=>evalueaza(`document.querySelector('${celula(1)}').value==='7'`),'Modificarea nesalvata din afara selectiei s-a pierdut');
 assert(await evalueaza("JSON.parse(localStorage.getItem('pontix-test')).absente.some(z=>z.cod==='CO'&&z.stare_aprobare==='APROBATA'&&z.nr_inregistrare==='TEST-CO')"),'Aprobarea/metadatele lipsesc');
 await evalueaza(`(()=>{const e=document.querySelector('${celula(0)}');e.dispatchEvent(new MouseEvent('mousedown',{bubbles:true,button:0}));e.focus()})()`);await apasa(t('pontaj.absenta_titlu'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('dialog select[name=tip]'))"),'Editarea absentei lipsa');
 await completeaza('dialog select[name=tip]','CM');await apasa(t('ui.salveaza'));
 await asteapta(()=>evalueaza(`document.querySelector('${celula(0)}').value==='CM'`),'Schimbarea CO-CM nu s-a aplicat');
 for (const cod of ['CO','CM']) {
  await evalueaza(`(()=>{const e=document.querySelector('${celula(0)}');e.dispatchEvent(new MouseEvent('mousedown',{bubbles:true,button:0}));e.focus()})()`);
  await apasa(t('pontaj.absenta_titlu'));
  await asteapta(()=>evalueaza("Boolean(document.querySelector('dialog select[name=tip]'))"),'Corectie explicita lipsa');
  await completeaza('dialog select[name=tip]',cod);await apasa(t('ui.salveaza'));
  await asteapta(()=>evalueaza(`document.querySelector('${celula(0)}').value==='${cod}'`),'Corectie explicita neaplicata: '+cod);
 }
 await apasa(t('pontaj.salveaza'));
 await asteapta(()=>evalueaza(`!document.querySelector('${celula(1)}').disabled&&!document.querySelector('${celula(1)}').closest('td').classList.contains('celula-modificata')`),'Salvarea celeilalte zile nu s-a terminat');
 await apasa(t('ui.absente'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.tabel-absente, .lista-absente')) || document.body.textContent.includes('TEST-CO')"),'Absenta lipseste din evidenta');
 assert(await evalueaza("document.body.textContent.includes('TEST-CO')"),'Metadatele nu sunt vizibile in Absente');
 await captura('absenta-pontaj-salvata-'+limba+'.png');
 await apasa(t('ui.editeaza'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('dialog input[name=de_la]'))"),'Editare zi lipsa');
 const dataNoua=await evalueaza("(()=>{const s=JSON.parse(localStorage.getItem('pontix-test'));const d=new Date(s.absente.find(z=>z.activ).data+'T12:00:00Z');d.setUTCDate(d.getUTCDate()+7);return d.toISOString().slice(0,10)})()");
 await completeaza('dialog input[name=de_la]',dataNoua);
 assert.equal(await evalueaza("document.querySelector('dialog input[name=pana_la]').value"),dataNoua,'Datele editarii nu sunt sincronizate');
 await apasa(t('ui.salveaza'));
 await asteapta(()=>evalueaza("!document.querySelector('dialog[open]')"),'Mutarea nu s-a salvat');
 assert(await evalueaza(`JSON.parse(localStorage.getItem('pontix-test')).absente.some(z=>z.data==='${dataNoua}'&&z.activ)`),'Ziua nu a fost mutata');
 await evalueaza('window.confirm=()=>true');
 await apasa(t('absente.sterge_zi'));
 await asteapta(()=>evalueaza("JSON.parse(localStorage.getItem('pontix-test')).absente.every(z=>!z.activ)"),'Medicalul nu a fost anulat');
 await completeaza('select[name=filtru_stare_zile_absenta]','INACTIVE');
 await asteapta(()=>evalueaza("document.querySelectorAll('.tabel-absente tbody tr').length>=2"),'Zilele anulate sunt inaccesibile');
 await apasa(t('ui.reactiveaza'));
 await asteapta(()=>evalueaza("JSON.parse(localStorage.getItem('pontix-test')).absente.some(z=>z.activ)"),'Reactivarea nu s-a salvat');
 await completeaza('select[name=filtru_stare_tipuri_absenta]','TOATE');
 for (const activa of [false,true]) {
  await evalueaza("document.querySelector('.tabel-tipuri-absenta tbody tr td:last-child button:nth-child(2)').click()");
  await asteapta(()=>evalueaza(`JSON.parse(localStorage.getItem('pontix-test')).tipuriAbsenta[0].activ===${activa}`),'Starea tipului nu s-a schimbat');
 }



}
