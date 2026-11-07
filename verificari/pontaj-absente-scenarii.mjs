import assert from 'node:assert/strict';

export async function verificaAbsentePontaj({evalueaza,asteapta,apasa,completeaza,cataloage,limba,captura}) {
 const t=k=>cataloage[limba][k];
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.actiuni-antet'))"),'Aplicatie lipsa');
 await apasa(t('ui.pontaj'));
 const celula=i=>`input[data-rand="0"][data-zi="${i}"]`;
 await asteapta(()=>evalueaza(`Boolean(document.querySelector('${celula(0)}'))`),'Grila lipsa');
 await completeaza('.unelte-perioada input[type=number]','2026');
 await evalueaza("document.querySelector('.unelte-perioada form').requestSubmit()");
 await completeaza('select[aria-label="'+t('pontaj.luna')+'"]','10');
 await asteapta(()=>evalueaza("document.querySelector('input[data-rand=\"0\"][data-zi=\"0\"]').getAttribute('aria-label').includes('2026-10-01')"),'Perioada fixa octombrie 2026 lipsa');
 assert.equal(await evalueaza("document.querySelector('select[name=filtru_stare_angajati_pontaj]').value"),'ACTIVE','Pontaj nu porneste pe activi');
 // D092: grila pastreaza doar OL, CO, CM, CFP, CVE, DLG, FMP, cu numele complet in titlu.
 const antetTotaluri=await evalueaza("[...document.querySelectorAll('.grila-pontaj thead tr:first-child th abbr')].map(a=>a.textContent)");
 assert.deepEqual(antetTotaluri,['OL','CO','CM','CFP','CVE','DLG','FMP'],'Coloanele de total din Pontaj: '+JSON.stringify(antetTotaluri));
 assert(await evalueaza("Boolean(document.querySelector('.grila-pontaj thead abbr[title]')?.title.length>5)"),'Titlul complet al coloanei lipseste');
 // R2: la 1024x640 se vad cel putin 5 randuri de angajati fara derularea paginii.
 await evalueaza('window.scrollTo(0,0)');
 const randuriVizibile=await evalueaza("[...document.querySelectorAll('.grila-pontaj tbody tr')].filter(r=>r.querySelector('input')&&r.getBoundingClientRect().bottom<=innerHeight).length");
 assert(randuriVizibile>=5,'Randuri vizibile in Pontaj la '+await evalueaza('innerWidth+"x"+innerHeight')+': '+randuriVizibile);
 // R3: vederea Totaluri arata toate cele 7 coloane in cadrul grilei, fara zile.
 await apasa(t('pontaj.vedere_totaluri'));
 const totaluriVizibile=await evalueaza("(()=>{const c=document.querySelector('.grila-container')||document.querySelector('.grila-pontaj').parentElement;const r=c.getBoundingClientRect();return [...document.querySelectorAll('.grila-pontaj thead tr:first-child th abbr')].filter(a=>{const b=a.closest('th').getBoundingClientRect();return b.width>0&&b.right<=r.right+1&&b.left>=r.left-1;}).length})()");
 assert.equal(totaluriVizibile,7,'Totaluri vizibile in vederea Totaluri: '+totaluriVizibile);
 assert(await evalueaza("[...document.querySelectorAll('.grila-pontaj thead th.col-zi')].every(th=>th.getBoundingClientRect().width===0)"),'Zilele raman vizibile in vederea Totaluri');
 await captura('pontaj-totaluri-'+limba+'.png');
 await apasa(t('pontaj.vedere_zile'));
 // R7: un buton dezactivat spune de ce.
 assert.equal(await evalueaza("[...document.querySelectorAll('button')].find(b=>b.textContent==="+JSON.stringify(t('pontaj.salveaza'))+")?.title"),t('motiv.fara_modificari'),'Salvarea dezactivata fara motiv');
 await completeaza('select[name=filtru_stare_angajati_pontaj]','INACTIVE');
 assert.equal(await evalueaza("document.querySelectorAll('.grila-pontaj tbody tr').length"),await evalueaza("JSON.parse(localStorage.getItem('pontix-test')).angajati.filter(a=>!a.activ).length"),'Filtrul Pontaj amesteca starile');
 await completeaza('select[name=filtru_stare_angajati_pontaj]','ACTIVE');
 await asteapta(()=>evalueaza(`Boolean(document.querySelector('${celula(0)}'))`),'Revenirea la activi nu reface grila');

 await evalueaza(`(()=>{const e=document.querySelector('${celula(0)}');e.dispatchEvent(new MouseEvent('mousedown',{bubbles:true,button:0}));e.focus()})()`);
 await completeaza('select[aria-label="'+t('pontaj.program_completare')+'"]','Z7');await apasa(t('pontaj.aplica_program'));
 await apasa(t('pontaj.salveaza'));
 await asteapta(()=>evalueaza("JSON.parse(localStorage.getItem('pontix-test')).ciorne.some(c=>c.valoare==='Z7')"),'Program nesalvat');
 await asteapta(()=>evalueaza(`!document.querySelector('${celula(0)}').disabled && !document.querySelector('${celula(0)}').closest('td').classList.contains('celula-modificata')`),'Salvare inca ocupata');
 // F05-T20: totalurile necalculate spun de ce, iar sumarul le grupeaza deasupra grilei.
 const motivRegula=t('pontaj.motiv_regula_nedefinita');
 await asteapta(()=>evalueaza(`[...document.querySelectorAll('.grila-pontaj td.total-necalculat')].some(td=>td.title===${JSON.stringify(motivRegula)})`),'Totalul necalculat nu are motivul regulii');
 await evalueaza("document.querySelector('.motive-necalculate summary').click()");
 assert(await evalueaza(`document.querySelector('.motive-necalculate').open && document.querySelector('.motive-necalculate li').textContent.includes(${JSON.stringify(motivRegula.slice(0,30))})`),'Sumarul motivelor lipseste sau e gol');
 await captura('pontaj-motive-'+limba+'.png');
 await evalueaza("document.querySelector('.motive-necalculate summary').click()");
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
 await evalueaza('window.__refuzaConversie=true');await apasa(t('ux30.previzualizeaza'));await asteapta(()=>evalueaza("Boolean(document.querySelector('.ux30-comparatie'))"),'Comparatie absenta lipsa');await apasa(t('ux30.confirma_modificari'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('dialog .mesaj.eroare'))"),'Eroare nesemnalata');
 await asteapta(()=>evalueaza(`document.querySelector('${celula(0)}').value==='7'`),'Eroarea a modificat programul');
 await evalueaza('window.__refuzaConversie=false');await captura('concediu-din-pontaj-'+limba+'.png');await apasa(t('ux30.previzualizeaza'));await asteapta(()=>evalueaza("Boolean(document.querySelector('.ux30-comparatie'))"),'Comparatie absenta lipsa');await apasa(t('ux30.confirma_modificari'));
 await asteapta(()=>evalueaza(`document.querySelector('${celula(0)}').value==='CO'`),'CO nu apare in Pontaj');
 await asteapta(()=>evalueaza(`document.querySelector('${celula(1)}').value==='7'`),'Modificarea nesalvata din afara selectiei s-a pierdut');
 assert(await evalueaza("JSON.parse(localStorage.getItem('pontix-test')).absente.some(z=>z.cod==='CO'&&z.stare_aprobare==='APROBATA'&&z.nr_inregistrare==='TEST-CO')"),'Aprobarea/metadatele lipsesc');
 await evalueaza(`(()=>{const e=document.querySelector('${celula(0)}');e.dispatchEvent(new MouseEvent('mousedown',{bubbles:true,button:0}));e.focus()})()`);await apasa(t('pontaj.absenta_titlu'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('dialog select[name=tip]'))"),'Editarea absentei lipsa');
 await completeaza('dialog select[name=tip]','CM');await apasa(t('ux30.previzualizeaza'));await asteapta(()=>evalueaza("Boolean(document.querySelector('.ux30-comparatie'))"),'Comparatie absenta lipsa');await apasa(t('ux30.confirma_modificari'));
 await asteapta(()=>evalueaza(`document.querySelector('${celula(0)}').value==='CM'`),'Schimbarea CO-CM nu s-a aplicat');
 for (const cod of ['CO','CM']) {
  await evalueaza(`(()=>{const e=document.querySelector('${celula(0)}');e.dispatchEvent(new MouseEvent('mousedown',{bubbles:true,button:0}));e.focus()})()`);
  await apasa(t('pontaj.absenta_titlu'));
  await asteapta(()=>evalueaza("Boolean(document.querySelector('dialog select[name=tip]'))"),'Corectie explicita lipsa');
  await completeaza('dialog select[name=tip]',cod);await apasa(t('ux30.previzualizeaza'));await asteapta(()=>evalueaza("Boolean(document.querySelector('.ux30-comparatie'))"),'Comparatie absenta lipsa');await apasa(t('ux30.confirma_modificari'));
  await asteapta(()=>evalueaza(`document.querySelector('${celula(0)}').value==='${cod}'`),'Corectie explicita neaplicata: '+cod);
 }
 await apasa(t('pontaj.salveaza'));
 await asteapta(()=>evalueaza(`!document.querySelector('${celula(1)}').disabled&&!document.querySelector('${celula(1)}').closest('td').classList.contains('celula-modificata')`),'Salvarea celeilalte zile nu s-a terminat');
 await apasa(t('ui.absente'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.tabel-absente, .lista-absente')) || document.body.textContent.includes('TEST-CO')"),'Absenta lipseste din evidenta');
 assert(await evalueaza("document.body.textContent.includes('TEST-CO')"),'Metadatele nu sunt vizibile in Absente');
 await captura('absenta-pontaj-salvata-'+limba+'.png');
 // F05-T21 (I050): ecranul se deschide pe evidenta zilelor; vederea pe cereri grupeaza si selecteaza zilele.
 assert(await evalueaza("document.querySelector('.taburi-sectiuni [aria-pressed=true]')?.textContent==="+JSON.stringify(t('absente.sectiune_zile'))+""),'Absentele nu se deschid pe zile');
 await apasa(t('absente.vedere_cereri'));
 await asteapta(()=>evalueaza("document.querySelectorAll('.tabel-cereri-absente tbody tr').length>=1"),'Vederea pe cereri lipseste');
 await evalueaza("document.querySelector('.tabel-cereri-absente tbody tr td:last-child button').click()");
 await asteapta(()=>evalueaza("[...document.querySelectorAll('.modul-absente button')].some(b=>b.textContent.startsWith("+JSON.stringify(t('absente.editeaza_selectate').split('(')[0])+")&&!b.textContent.includes('(0)')&&!b.disabled)"),'Selectarea zilelor cererii nu activeaza editarea in lot');
 await captura('absente-cereri-'+limba+'.png');
 await apasa(t('absente.vedere_zile'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.tabel-absente'))"),'Revenirea la zile lipsa');
 await apasa(t('ui.editeaza'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('dialog input[name=de_la]'))"),'Editare zi lipsa');
 const dataNoua=await evalueaza("(()=>{const s=JSON.parse(localStorage.getItem('pontix-test'));const d=new Date(s.absente.find(z=>z.activ).data+'T12:00:00Z');d.setUTCDate(d.getUTCDate()+7);return d.toISOString().slice(0,10)})()");
 await completeaza('dialog input[name=de_la]',dataNoua);
 assert.equal(await evalueaza("document.querySelector('dialog input[name=pana_la]').value"),dataNoua,'Datele editarii nu sunt sincronizate');
 await apasa(t('ux30.previzualizeaza'));await asteapta(()=>evalueaza("Boolean(document.querySelector('.ux30-comparatie'))"),'Comparatie absenta lipsa');await apasa(t('ux30.confirma_modificari'));
 await asteapta(()=>evalueaza("!document.querySelector('dialog[open]')"),'Mutarea nu s-a salvat');
 assert(await evalueaza(`JSON.parse(localStorage.getItem('pontix-test')).absente.some(z=>z.data==='${dataNoua}'&&z.activ)`),'Ziua nu a fost mutata');
 await evalueaza('window.confirm=()=>true');
 await apasa(t('absente.sterge_zi'));
 await asteapta(()=>evalueaza("JSON.parse(localStorage.getItem('pontix-test')).absente.every(z=>!z.activ)"),'Medicalul nu a fost anulat');
 await completeaza('select[name=filtru_stare_zile_absenta]','INACTIVE');
 await asteapta(()=>evalueaza("document.querySelectorAll('.tabel-absente tbody tr').length>=2"),'Zilele anulate sunt inaccesibile');
 await apasa(t('ui.reactiveaza'));await asteapta(()=>evalueaza("Boolean(document.querySelector('dialog select[name=tip]'))"),'Reactivare lipsa');await apasa(t('ux30.previzualizeaza'));await asteapta(()=>evalueaza("Boolean(document.querySelector('.ux30-comparatie'))"),'Comparatie reactivare lipsa');await apasa(t('ux30.confirma_modificari'));
 await asteapta(()=>evalueaza("JSON.parse(localStorage.getItem('pontix-test')).absente.some(z=>z.activ)"),'Reactivarea nu s-a salvat');
 await apasa(t('absente.sectiune_toate'));
 await completeaza('select[name=filtru_stare_tipuri_absenta]','TOATE');
 for (const activa of [false,true]) {
  await evalueaza("document.querySelector('.tabel-tipuri-absenta tbody tr td:last-child button:nth-child(2)').click()");
  await asteapta(()=>evalueaza(`JSON.parse(localStorage.getItem('pontix-test')).tipuriAbsenta[0].activ===${activa}`),'Starea tipului nu s-a schimbat');
 }



}
