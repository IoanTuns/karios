import assert from 'node:assert/strict';

export async function verificaFocus({cdp,evalueaza,asteapta,apasa,cataloage,limba,captura}) {
 const t=k=>cataloage[limba][k];
 await asteapta(()=>evalueaza("Boolean(document.querySelector('.actiuni-antet'))"),'Antet lipsa');
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
  await captura('zoom-'+ecran+'-'+limba+'.png');
 }
}
