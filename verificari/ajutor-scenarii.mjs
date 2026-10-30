// F02-T21: scenarii pe aplicatia reala cu punte si date exclusiv fictive.
import assert from 'node:assert/strict';

export function simulareAjutor(catalog, texte, versiune, limba) {
    const api = window.go.main.App;
    api.Limbi = async () => [limba, ...Object.keys(texte).filter(l => l !== limba)];
    if (limba === 'en') localStorage.setItem('pontix_tema','dark');
    api.StareAudit = async () => ({rezultat:{stare:'valid'}});
    api.StareSeif = async () => ({blocat:false});
    api.StareCont = async () => ({autentificat:true, cont:{id:'test',nume:'Operator fictiv',rol:'ADMINISTRATOR',activ:true}});
    api.StareRecuperare = async () => ({necesara:false});
    api.SeteazaLimbaInterfata = async () => {};
    api.SeteazaModificariNesalvate = async () => {};
    api.SarcinaAIActiva = async () => '';
    api.StareAI = async () => ({stare:{activ:false,instalat:false}});
    api.EsteAutoPornireAI = async () => false;
    window.__eroriAjutor = [];
    api.InregistreazaEroareFrontend = async (...args) => {window.__eroriAjutor.push(args.join(' '));};
    window.runtime ??= {};
    window.runtime.EventsOnMultiple = () => () => {};
    api.ArticoleAjutor = async l => catalog.map(d => {
        const prefix='ajutor.articol.'+d.id+'.';
        return {...d,versiune,titlu:texte[l][prefix+'titlu'],pasi:[1,2,3].map(i=>texte[l][prefix+'pas'+i]),
            rezultat:texte[l][prefix+'rezultat'],probleme:texte[l][prefix+'probleme']};
    });
    api.DeschideManualAjutor = async id => {window.__manualAjutor=id;return null;};
    sessionStorage.setItem('pontix_asistent_conversatie',JSON.stringify([
        {id:'ghid-test',rol:'assistant',text:'Ghid fictiv',articoleAjutor:['absente.inregistreaza'],versiuneAjutor:versiune},
        {id:'ghid-invalid',rol:'assistant',text:'ID fictiv inexistent',articoleAjutor:['necunoscut.fictiv'],versiuneAjutor:versiune},
        {id:'ghid-vechi',rol:'assistant',text:'Versiune fictiva veche',articoleAjutor:['pontaj.completeaza'],versiuneAjutor:'0.0.0'}
    ]));
    void api.IncarcaDemo();
}

export async function verificaAjutor({cdp,evalueaza,asteapta,completeaza,apasa,cataloage,limba,captura}) {
    const t=k=>cataloage[limba][k];
    await asteapta(()=>evalueaza("Boolean(document.querySelector('.ajutor-antet'))"),'Butonul contextual lipseste.');
    await apasa(t('ajutor.buton'));
    await asteapta(()=>evalueaza("Boolean(document.querySelector('dialog[open] .ajutor-articol h3'))"),'Articolul initial lipseste.');
    await completeaza('dialog input[type="search"]','CONCEDIU');
    await asteapta(()=>evalueaza("document.querySelectorAll('.ajutor-lista button').length>0"),'Cautarea globala nu gaseste concediul.');
    await apasa(t('ajutor.articol.absente.inregistreaza.titlu'));
    assert(await evalueaza("document.querySelector('.ajutor-articol').textContent.includes('absente.inregistreaza')"),'Lipseste sursa.');
    await apasa(t('ajutor.manual'));
    await asteapta(()=>evalueaza("window.__manualAjutor==='absente.inregistreaza'"),'ID gresit la deschiderea manualului.');
    await captura('ajutor-'+limba+'.png');
    await completeaza('dialog input[type="search"]','faraacoperire987654');
    await asteapta(()=>evalueaza("document.querySelectorAll('.ajutor-lista button').length===0"),'Cautare fara rezultate incorecta.');
    await completeaza('dialog input[type="search"]','concediu');
    await apasa(t('ajutor.deschide'));
    await asteapta(()=>evalueaza("!document.querySelector('dialog[open]') && document.querySelector('main h1').textContent==="+JSON.stringify(t('ui.absente'))),'Navigarea nu deschide sectiunea.');
    await apasa(t('ui.asistent_ai'));
    await asteapta(()=>evalueaza("document.querySelectorAll('.bula-mesaj .legatura').length>=3"),'Citari AI lipsa.');
    await apasa(t('ajutor.citire').replace('{{id}}','absente.inregistreaza'));
    await asteapta(()=>evalueaza("Boolean(document.querySelector('.ajutor-articol h3'))"),'Citarea valida nu deschide ghidul.');
    await cdp('Input.dispatchKeyEvent',{type:'keyDown',key:'Escape',code:'Escape',windowsVirtualKeyCode:27});
    await asteapta(()=>evalueaza("!document.querySelector('dialog[open]')"),'Escape nu inchide ghidul.');
    assert(await evalueaza("document.activeElement?.textContent.includes('absente.inregistreaza')"),'Focusul nu revine la citare.');
    await apasa(t('ajutor.citire').replace('{{id}}','necunoscut.fictiv'));
    await asteapta(()=>evalueaza("Boolean(document.querySelector('dialog [role=alert]'))"),'ID necunoscut fara mesaj.');
    await cdp('Input.dispatchKeyEvent',{type:'keyDown',key:'Escape',code:'Escape',windowsVirtualKeyCode:27});
    await asteapta(()=>evalueaza("!document.querySelector('dialog[open]')"),'Inchidere ID invalid esuata.');
    await apasa(t('ajutor.citire').replace('{{id}}','pontaj.completeaza'));
    await asteapta(()=>evalueaza("Boolean(document.querySelector('dialog [role=alert]'))"),'Versiune veche acceptata.');
    await captura('ajutor-versiune-'+limba+'.png');
    const geometrie=await evalueaza("(()=>{const d=document.querySelector('dialog'),r=d.getBoundingClientRect();return {latime:r.width,inaltime:r.height,viewport:innerWidth,vertical:innerHeight}})()");
    assert(geometrie.latime<=geometrie.viewport && geometrie.inaltime<=geometrie.vertical,'Dialogul depaseste ecranul.');
    await cdp('Input.dispatchKeyEvent',{type:'keyDown',key:'Escape',code:'Escape',windowsVirtualKeyCode:27});
    await asteapta(()=>evalueaza("!document.querySelector('dialog[open]')"),'Inchidere esuata.');
    await apasa(t('ajutor.buton'));
    await asteapta(()=>evalueaza("Boolean(document.querySelector('.ajutor-articol h3'))"),'Ghid AI lipsa.');
    await apasa(t('ajutor.intreaba'));
    await asteapta(()=>evalueaza("!document.querySelector('dialog[open]') && Boolean(document.querySelector('.formular-intrare-ai input')?.value.trim())"),'Intrebarea nu a ajuns in campul AI.');
    assert(await evalueaza("document.querySelectorAll('.bula-mesaj').length===3"),'Intrebarea a fost trimisa automat.');
    assert.deepEqual(await evalueaza("window.__eroriAjutor"),[],'Erori UI raportate.');
}
