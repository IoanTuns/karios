import assert from 'node:assert/strict';

export async function verificaFocus({cdp,evalueaza,asteapta,apasa,cataloage,limba,captura}) {
 const t=k=>cataloage[limba][k];
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.actiuni-antet'))"),'Antet lipsa');
 // F16-T29: identificarea si antetul raman vizibile si in tabelele obisnuite.
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.tabel-angajati[data-tabel-accesibil]'))"),'Tabel angajati lipsa');
 await evalueaza("(()=>{const g=document.querySelector('.tabel-angajati');g.scrollLeft=100000;g.scrollTop=100;g.scrollIntoView({block:'center'})})()");
 await new Promise(r=>setTimeout(r,80));
 const identitate=await evalueaza("(()=>{const g=document.querySelector('.tabel-angajati'),r=g.getBoundingClientRect(),c=[...g.querySelectorAll('thead th')].slice(0,3).map(e=>{const b=e.getBoundingClientRect();return {left:b.left,right:b.right,top:b.top}});return {left:r.left,right:r.right,top:r.top,c}})()");
 assert(identitate.c[1].left>=identitate.c[0].right-2 && identitate.c[2].left>=identitate.c[1].right-2 && identitate.c[2].right<=identitate.right,'Identitatea sau marca se suprapun in lista '+JSON.stringify(identitate));
 assert(identitate.c.every(c=>c.top>=identitate.top-2 && c.top<=identitate.top+4),'Antetul listei nu ramane la derulare');
 await captura('tabel-angajati-derulat-'+limba+'.png');
 await evalueaza("(()=>{const g=document.querySelector('.tabel-angajati');g.scrollLeft=0;g.scrollTop=0})()");
 // F16-T25: butoane globale, limite, tastatura si memorare.
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.zoom-interfata'))"),'Zoom global lipsa');
 for (let i=0;i<12;i++) await evalueaza("document.querySelector('.zoom-interfata button:last-child').click()");
 await asteapta(()=>evalueaza("document.querySelector('.zoom-interfata').textContent.includes('200%')"),'Limita maxima zoom');
 assert(await evalueaza("document.querySelector('.zoom-interfata button:last-child').disabled"));
 assert.equal(await evalueaza("getComputedStyle(document.documentElement).zoom"),'2');
 await cdp('Emulation.setDeviceMetricsOverride',{width:1024,height:768,deviceScaleFactor:1,mobile:false});
 await evalueaza("document.querySelector('.zoom-interfata').scrollIntoView({block:'nearest'})");
 assert(await evalueaza("(()=>{const r=document.querySelector('.zoom-interfata').getBoundingClientRect();return r.left>=0 && r.right<=innerWidth})()"),'Controalele globale ies din ecran la 200%');
 await apasa(t('ui.adauga_angajat'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('dialog[open]'))"),'Dialog la 200% lipsa');
 assert(await evalueaza("(()=>{const r=document.querySelector('dialog[open]').getBoundingClientRect();return r.left>=0&&r.right<=innerWidth+1&&r.top>=0&&r.bottom<=innerHeight+1})()"),'Dialogul depaseste ecranul la zoom 200%');
 await captura('zoom-global-dialog-200-'+limba+'.png');
 await cdp('Input.dispatchKeyEvent',{type:'keyDown',key:'Escape',code:'Escape',windowsVirtualKeyCode:27});
 await asteapta(()=>evalueaza("!document.querySelector('dialog[open]')"),'Dialog la 200% ramas deschis');
 await captura('zoom-global-200-'+limba+'.png');
 await cdp('Emulation.setDeviceMetricsOverride',{width:1440,height:960,deviceScaleFactor:1,mobile:false});

 for (let i=0;i<12;i++) await evalueaza("document.querySelector('.zoom-interfata button:first-child').click()");
 await asteapta(()=>evalueaza("document.querySelector('.zoom-interfata').textContent.includes('75%')"),'Limita minima zoom');
 assert(await evalueaza("document.querySelector('.zoom-interfata button:first-child').disabled"));
 await cdp('Input.dispatchKeyEvent',{type:'keyDown',key:'0',code:'Digit0',modifiers:2,windowsVirtualKeyCode:48});
 await asteapta(()=>evalueaza("document.querySelector('.zoom-interfata').textContent.includes('100%')"),'Resetare tastatura zoom');
 await cdp('Input.dispatchKeyEvent',{type:'keyDown',key:'+',code:'Equal',modifiers:2,windowsVirtualKeyCode:187});
 await asteapta(()=>evalueaza("localStorage.getItem('pontix_zoom_interfata')==='110'"),'Preferinta zoom nesalvata');
 await evalueaza("document.querySelector('.zoom-interfata button:nth-child(2)').click()");
 await apasa(t('ui.adauga_angajat'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('dialog[open] input'))"),'Formular lipsa');
 assert(await evalueaza("Boolean(document.querySelector('dialog .buton-marire'))"),'Butonul de marire lipseste.');
 await evalueaza("document.querySelector('dialog .buton-marire').click()");
 await evalueaza("window.__campuri=[...document.querySelectorAll('dialog input')];window.__campuri[1].focus()");
 for(const litera of 'Ana'){
  await cdp('Input.insertText',{text:litera});
  await new Promise(r=>setTimeout(r,40));
  assert(await evalueaza("document.activeElement===window.__campuri[1]"),'Tastarea pierde focusul din Prenume.');
 }
 assert.equal(await evalueaza("window.__campuri[1].value"),'Ana');
 await evalueaza("window.__campuri[1].setSelectionRange(1,1)");
 await cdp('Input.insertText',{text:'X'});
 await new Promise(r=>setTimeout(r,40));
 assert.deepEqual(await evalueaza("[window.__campuri[1].value,window.__campuri[1].selectionStart]"),['AXna',2],'Cursorul nu ramane in mijlocul textului.');
 assert(await evalueaza("Boolean(document.querySelector('dialog.dialog-maximizat'))"),'Marirea se pierde dupa tastare.');
 await evalueaza("document.querySelector('dialog .buton-marire').click()");
 for(const index of [0,2]){
  await evalueaza(`window.__campuri[${index}].focus();window.__campuri[${index}].setSelectionRange(0,0)`);
  await cdp('Input.insertText',{text:'Q'});
  await new Promise(r=>setTimeout(r,40));
  assert(await evalueaza(`document.activeElement===window.__campuri[${index}] && window.__campuri[${index}].selectionStart===1`),'Focus pierdut in campul '+index);
 }
 await evalueaza("document.querySelector('details.solduri-initiale').open=true");
 // F07-T27: initializarea pastreaza solduri negative si CO in formular.
 assert(await evalueaza("Boolean(document.querySelector('input[name=initializare_luna][type=month]'))"),'Luna de initializare lipsa');
 await evalueaza("document.querySelector('input[name=initializare_ore]').focus();document.querySelector('input[name=initializare_ore]').select()");
 await cdp('Input.insertText',{text:'-2.5'});
 await asteapta(()=>evalueaza("document.querySelector('input[name=initializare_ore]').value==='-2.5'"),'Soldul negativ nu se pastreaza');
 await evalueaza("document.querySelector('input[name=initializare_co]').focus();document.querySelector('input[name=initializare_co]').select()");
 await cdp('Input.insertText',{text:'8'});
 await asteapta(()=>evalueaza("document.querySelector('input[name=initializare_co]').value==='8' && document.querySelector('input[name=initializare_ore]').value==='-2.5'"),'CO initial pierde soldul de ore');
 await captura('initializare-solduri-'+limba+'.png');
 await evalueaza("window.__confirmInitial=window.confirm;window.__confirmRefuz=0;window.confirm=()=>{window.__confirmRefuz++;return false}");
 await cdp('Input.dispatchKeyEvent',{type:'keyDown',key:'Escape',code:'Escape',windowsVirtualKeyCode:27});
 await new Promise(r=>setTimeout(r,40));
 assert(await evalueaza("window.__confirmRefuz===1 && Boolean(document.querySelector('dialog[open]'))"),'Escape foloseste o stare veche a formularului.');
 await evalueaza("window.confirm=window.__confirmInitial");
 await captura('formular-focus-'+limba+'.png');
 await cdp('Input.dispatchKeyEvent',{type:'keyDown',key:'Escape',code:'Escape',windowsVirtualKeyCode:27});
 await asteapta(()=>evalueaza("!document.querySelector('dialog[open]')"),'Escape nu inchide formularul.');
 await evalueaza("window.go.main.App.AnalizeazaImportAngajati=async()=>({rezultat:{nume_fisier:'fictiv.csv',format:'csv',total_randuri:1,total_valide:1,total_existente:0,total_invalide:0,randuri:[{nume:'Fictiv',prenume:'Test',marca:'TEST-IMPORT',functie:'Test',sectie:'',departament:'',tip_contract:'PARTIALA',norma_minute:270,concediu_zile_an:22,valid:true,exista_deja:false}]}})");
 await apasa(t('ui.importa'));
 await asteapta(()=>evalueaza("Boolean(document.querySelector('input[type=file]'))"),'Dialog import lipsa.');
 await evalueaza("(()=>{const e=document.querySelector('input[type=file]'),d=new DataTransfer();d.items.add(new File(['Nume;Zile_Concediu_An\\nFictiv;22'],'fictiv.csv',{type:'text/csv'}));e.files=d.files;e.dispatchEvent(new Event('change',{bubbles:true}))})()");
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.sectiune-analiza tbody tr'))"),'Previzualizare import lipsa.');
 assert.equal(await evalueaza("document.querySelector('.sectiune-analiza tbody tr td:last-child').textContent"),'22','Concediul nu apare in previzualizare.');
 assert(await evalueaza("document.querySelector('.sectiune-analiza tbody tr td:nth-last-child(2)').textContent.includes('30')"),'Norma fractionara este rotunjita in previzualizare.');
 await captura('import-concediu-'+limba+'.png');
 await evalueaza("document.querySelector('.fereastra-import .buton-inchidere, .fereastra-import button[aria-label]').click()");
 await asteapta(()=>evalueaza("!document.querySelector('.fereastra-import')"),'Dialog import ramas deschis.');
 for(const ecran of ['pontaj','planificare']){
  await apasa(t('ui.'+ecran));
  await asteapta(()=>evalueaza("Boolean(document.querySelector('.grup-zoom-grila'))"),'Zoom lipsa '+ecran);
  for(const width of [1024,1366]){
   await cdp('Emulation.setDeviceMetricsOverride',{width,height:768,deviceScaleFactor:1,mobile:false});
   await evalueaza("document.querySelector('.grup-zoom-grila').scrollIntoView({block:'nearest'})");
   const g=await evalueaza("(()=>{const e=document.querySelector('.grup-zoom-grila'),r=e.getBoundingClientRect();return {w:r.width,x:r.left,right:r.right,buttons:[...e.querySelectorAll('button')].map(b=>b.getBoundingClientRect().width)}})()");
   assert(g.x>=0&&g.right<=width&&g.buttons.length===3&&g.buttons.every(n=>n>=25),'Zoom comprimat sau ascuns '+ecran);
  }
  for(const [index,mod] of [[0,'compact'],[1,'normal'],[2,'extins']]){
   await evalueaza(`document.querySelectorAll('.grup-zoom-grila button')[${index}].click()`);
   await asteapta(()=>evalueaza(`Boolean(document.querySelector('.zoom-${mod}'))`),'Zoom neaplicat '+ecran+' '+mod);
  }
  for(const tema of ['light','dark']){
   await evalueaza(`document.documentElement.setAttribute('data-theme','${tema}')`);
   await evalueaza("document.querySelector('.grila-container').scrollLeft=100000");
   await new Promise(r=>setTimeout(r,80));
   const poz=await evalueaza("(()=>{const g=document.querySelector('.grila-container'),r=g.getBoundingClientRect(),n=g.querySelector('thead .fix-nume').getBoundingClientRect(),m=g.querySelector('thead .fix-marca').getBoundingClientRect();return {left:r.left,right:r.right,nLeft:n.left,nRight:n.right,mLeft:m.left,mRight:m.right,scroll:g.scrollLeft,latime:g.clientWidth,role:g.getAttribute('role')}})()");
   assert(poz.scroll>0,'Grila nu se deruleaza '+ecran);
   assert(poz.nLeft>=poz.left-2&&poz.nRight<poz.right&&poz.mLeft>=poz.nRight-2&&poz.mRight<=poz.right,'Identificarea nu ramane vizibila sau se suprapune '+ecran+' '+JSON.stringify(poz));
   assert.equal(poz.role,'region','Tabelul nu are navigare accesibila');
   await captura('derulare-dreapta-'+ecran+'-'+tema+'-'+limba+'.png');
   await evalueaza("document.querySelector('.grila-container').focus()");
   await cdp('Input.dispatchKeyEvent',{type:'keyDown',key:'Home',code:'Home',modifiers:1,windowsVirtualKeyCode:36});
   await asteapta(()=>evalueaza("document.querySelector('.grila-container').scrollLeft===0"),'Alt+Home nu revine la inceput');
  }
  for(const nivel of [150,200]){
   await evalueaza(`(()=>{const n=Number(localStorage.getItem('pontix_zoom_interfata')||100);const lista=[75,80,90,100,110,125,150,175,200];for(let i=lista.indexOf(n);i<lista.indexOf(${nivel});i++)document.querySelector('.zoom-interfata button:last-child').click();})()`);
   await asteapta(()=>evalueaza(`localStorage.getItem('pontix_zoom_interfata')==='${nivel}'`),'Zoom global neaplicat');
   await evalueaza("document.querySelector('.grila-container').scrollLeft=100000;document.querySelector('.grila-container').scrollTop=100");
   await new Promise(r=>setTimeout(r,100));
   const p=await evalueaza("(()=>{const g=document.querySelector('.grila-container'),r=g.getBoundingClientRect(),n=g.querySelector('thead .fix-nume').getBoundingClientRect(),m=g.querySelector('thead .fix-marca').getBoundingClientRect();return {left:r.left,right:r.right,nLeft:n.left,nRight:n.right,mLeft:m.left,mRight:m.right}})()");
   assert(p.nLeft>=p.left-3&&p.mLeft>=p.nRight-3&&p.mRight<=p.right,'Repere suprapuse la zoom '+nivel+' '+ecran+' '+JSON.stringify(p));
   await captura('tabel-zoom-'+nivel+'-'+ecran+'-'+limba+'.png');
  }
  await evalueaza("document.querySelector('.zoom-interfata button:nth-child(2)').click()");
  await captura('zoom-'+ecran+'-'+limba+'.png');
 }
}
