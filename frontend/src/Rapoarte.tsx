import {EditorPrezentare,EditorCopieRaport,stilCSS} from './PrezentareRaport';
import type {TFunction} from 'i18next';
import {EditorFlexibil,esteFlexibil,flexibilInitial} from './RaportFlexibil';
import {ButonStergereAdmin} from './ButonStergereAdmin';
import {FiltruStare,corespundeStare,type StareLista} from './FiltruStare';
import {numeAngajat} from './numeAngajat';
import {FiltreProfesionale,filtruProfesionalGol,corespundeProfesional} from './FiltreProfesionale';
import {useEffect, useRef, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {semnaleazaModificari} from './modificariNesalvate';
import {ListaModeleRapoarte, SalveazaModelRaport, PregatesteRaport, SalveazaPDFRaport, DeschidePDFPentruTiparire, SalveazaExcelRaport} from '../wailsjs/go/main/App';
import {datebaza, main, rapoarte} from '../wailsjs/go/models';
import {Dialog, denumireDepartament, denumireSectie} from './Formulare';
import {completeazaRaport, sablonPePagina} from './raportPdfme';
import {notifica} from './Toast';
import {limbaRapoarte} from './limba';

const copiaza = <T,>(v: T): T => JSON.parse(JSON.stringify(v));
export type DateSimple<T> = T extends Array<infer V> ? DateSimple<V>[] : T extends object ? {[K in keyof T as T[K] extends Function ? never : K]: DateSimple<T[K]>} : T;
const textIn = (v: {[limba: string]: string} | undefined, limba: string) => v?.[limba] || v?.ro || Object.values(v || {}).find(Boolean) || '';
const adresaImagine = (i: DateSimple<rapoarte.Imagine>) => `data:image/${i.date.startsWith('iVBOR') ? 'png' : 'jpeg'};base64,${i.date}`;

// F09-T23: modelele vechi trec pe editorul vizual la primul acces; anuntul spune cate
// si ce s-a aproximat (originalele raman in jurnal).
function anuntaConversia(conv: rapoarte.RezultatConversie[] | undefined, t: TFunction) {
    if (!conv?.length) return;
    const erori = conv.filter(c => c.eroare);
    const note = conv.flatMap(c => (c.note || []).map(n => `${c.denumire.ro || c.id}: ${n}`));
    notifica(t('rapoarte.conversie_anunt', {numar: conv.length - erori.length}) + (note.length ? ' ' + note.join(' ') : ''), note.length ? 'avertisment' : 'info', 12000);
    if (erori.length) notifica(t('rapoarte.conversie_erori', {modele: erori.map(c => c.denumire.ro || c.id).join(', ')}), 'eroare', 12000);
}

function useModele() {
    const {t} = useTranslation();
    const [configuratie, setConfiguratie] = useState<rapoarte.Configuratie | null>(null);
    const [incarcare, setIncarcare] = useState(true);
    const [eroare, setEroare] = useState('');
    const [revizie, setRevizie] = useState(0);
    function aplica(r: main.RaspunsModeleRapoarte) { anuntaConversia(r.conversie, t); if (r.eroare) { setEroare(t(r.eroare.cod,r.eroare.param || {})); return false; } if (!r.configuratie) { setEroare(t('rapoarte.eroare')); return false; } setConfiguratie(r.configuratie); return true; }
    useEffect(() => { let activ = true; setIncarcare(true); setEroare(''); ListaModeleRapoarte().then(r => { if (activ) aplica(r); }).catch(() => { if (activ) setEroare(t('rapoarte.eroare')); }).finally(() => { if (activ) setIncarcare(false); }); return () => { activ = false; }; },[revizie]);
    return {configuratie,incarcare,eroare,setEroare,aplica,reincarca:() => setRevizie(r => r+1)};
}

export function DefinireRapoarte({date}: {date?: datebaza.Date} = {}) {
 const [stareLista,setStareLista]=useState<StareLista>('ACTIVE');
    const {t,i18n} = useTranslation();
    const {configuratie,incarcare,eroare,setEroare,aplica,reincarca} = useModele();
    const [editare,setEditare] = useState<rapoarte.Definitie | null>(null);
    const [ocupat,setOcupat] = useState(false);
    const blocat = useRef(false);
    const [mesaj,setMesaj] = useState('');
    const modele = configuratie?.modele || [];
    function nou(d?: rapoarte.Definitie) {
        setEroare(''); setMesaj('');
        if (d) {
            setEditare(new rapoarte.Definitie({...copiaza(d),id:'',revizie:'',activ:true,continut:{...copiaza(d.continut),ascunde_explicatii:!!d.continut.ascunde_explicatii},denumire:copiaza(d.denumire)}));
            return;
        }
        // F09-T31: un model nou porneste gol (SQL personalizat, fara antet, subsol sau SQL);
        // copierea unui model existent ramane la „Duplică”.
        const gol = () => ({ro:'',en:''});
        const bloc = () => ({stanga:gol(),centru:gol(),dreapta:gol()});
        setEditare(new rapoarte.Definitie({id:'',revizie:'',activ:true,denumire:gol(),antet:bloc(),subsol:bloc(),
            continut:{sursa:'SQL_LIBER',coloane:['rezultat'],text:gol(),ascunde_explicatii:true,flexibil:flexibilInitial()},
            pagina:{format:'A4',orientare:'PEISAJ',margine_mm:10,font_pt:8}}));
    }
    async function salveaza(d: rapoarte.Definitie) {
        if (blocat.current) return;
        blocat.current = true; setOcupat(true); setEroare('');
        try { if (aplica(await SalveazaModelRaport(d))) { setEditare(null); setMesaj(t('ui.salvat')); } }
        catch { setEroare(t('rapoarte.eroare')); }
        finally { blocat.current = false; setOcupat(false); }
    }
    return <section className="definitii-rapoarte"><FiltruStare nume="modele_rapoarte" valoare={stareLista} schimba={setStareLista} disabled={ocupat}/>
        <div className="unelte-calendar"><p className="nota-calendar">{t('rapoarte.modele_indiciu')}</p><div className="actiuni-calendar"><button className="legatura" disabled={ocupat || incarcare} onClick={reincarca}>{t('rapoarte.reincarca')}</button><button className="buton primar" disabled={ocupat || incarcare} onClick={() => nou()}>{t('rapoarte.model_nou')}</button></div></div>
        {mesaj && <p className="mesaj succes" role="status">{mesaj}</p>}{eroare && !editare && <p className="mesaj eroare" role="alert">{eroare}</p>}
        {incarcare ? <p role="status">{t('ui.incarcare')}</p> : <div className="panou tabel-derulant"><table className="tabel-modele-rapoarte"><thead><tr>{['schimburi.denumire','rapoarte.sursa','rapoarte.revizie','rapoarte.pagina','ui.stare','ui.actiuni'].map(k => <th key={k}>{t(k)}</th>)}</tr></thead><tbody>
            {modele.filter(d=>corespundeStare(d.activ,stareLista)).map(d => <tr key={d.id}><td><strong>{textIn(d.denumire,i18n.language)}</strong></td><td>{t(`rapoarte.sursa_${d.continut.sursa}`)}</td><td>{d.revizie}</td><td>{d.pagina.format} · {t(`rapoarte.${d.pagina.orientare}`)}</td><td>{t(d.activ ? 'ui.activ':'ui.inactiv')}</td><td><div className="actiuni-rand"><button className="legatura" disabled={ocupat} onClick={() => { setEroare(''); setEditare(copiaza(d)); }}>{t('ui.editeaza')}</button><button className="legatura discreta" disabled={ocupat} onClick={() => nou(d)}>{t('rapoarte.duplica')}</button><ButonStergereAdmin tip="modele_raport" id={d.id} revizie={d.revizie} activ={d.activ} nume={textIn(d.denumire,i18n.language)} disabled={ocupat||!!editare} onOcupat={setOcupat} onSters={()=>{reincarca();}}/></div></td></tr>)}
        </tbody></table>{!modele.length && <p className="fara-rezultate">{t('rapoarte.lista_goala')}</p>}</div>}
        {!incarcare && !configuratie && <button className="buton secundar" onClick={reincarca}>{t('ui.reincearca')}</button>}
        {editare && configuratie && <EditorModel key={editare.id || 'nou'} initial={editare} configuratie={configuratie} ocupat={ocupat} eroare={eroare} inchide={() => { setEditare(null);setEroare(''); }} salveaza={salveaza} angajati={date?.angajati ?? []}/>}
    </section>;
}

function EditorModel({initial,configuratie,ocupat:salvareInCurs,eroare,inchide,salveaza,angajati}: {initial: rapoarte.Definitie; configuratie: rapoarte.Configuratie; ocupat: boolean; eroare: string; inchide: () => void; salveaza: (d: rapoarte.Definitie) => Promise<void>; angajati: datebaza.Angajat[]}) {
    const {t,i18n} = useTranslation();
    // F09-T32: la deschidere, sablonul vizual se potriveste cu Setări de pagină.
    const [model,setModel] = useState<DateSimple<rapoarte.Definitie>>(() => { const c = copiaza(initial); return cuPagina(c, c.pagina); });
    useEffect(() => { if (model.continut.flexibil?.sablon !== initial.continut.flexibil?.sablon && initial.continut.flexibil?.sablon) notifica(t('rapoarte.pdfme_pagina_potrivita'), 'avertisment', 8000); }, []);
    const [limba,setLimba] = useState(i18n.language === 'en' ? 'en' : 'ro');
    const [validare,setValidare] = useState('');
 const [validatFlex,setValidatFlex]=useState('');
    const [citireImagine,setCitireImagine] = useState(false);
    const importInCurs = useRef(false);
    const ocupat = salvareInCurs || citireImagine;
    const modificate = JSON.stringify(model) !== JSON.stringify(initial);
    // F09-T28: activarea sau dezactivarea nu schimba continutul raportului, deci nu cere previzualizare.
    const doarStare = initial.id !== '' && JSON.stringify({...model, activ: initial.activ}) === JSON.stringify(initial);
    useEffect(() => semnaleazaModificari('rapoarte', modificate), [modificate]);
    function renunta() { if (!modificate || window.confirm(t('ui.renunta'))) inchide(); }
    function textBloc(bloc: 'antet' | 'subsol',zona: 'stanga'|'centru'|'dreapta',v: string) { setModel(m => ({...m,[bloc]:{...m[bloc],[zona]:{...m[bloc][zona],[limba]:v}}})); }
    function seteazaImagine(bloc: 'antet'|'subsol',zona: string,imagine?: DateSimple<rapoarte.Imagine>) {
        setModel(m => { const imagini={...m[bloc].imagini}; if(imagine) imagini[zona]=imagine;else delete imagini[zona];return {...m,[bloc]:{...m[bloc],imagini}}; });
    }
    async function importaImagine(bloc: 'antet'|'subsol',zona: string,fisier?: File) {
        if(!fisier || importInCurs.current) return;
        importInCurs.current=true;setCitireImagine(true);setValidare('');
        try {
            if(fisier.size>2*1024*1024) throw new Error('rapoarte.imagine_prea_mare');
            if(!['image/png','image/jpeg'].includes(fisier.type) || fisier.name.length>200) throw new Error('rapoarte.imagine_invalida');
            const adresa=await new Promise<string>((rezolva,respinge) => { const cititor=new FileReader();cititor.onload=()=>rezolva(String(cititor.result));cititor.onerror=()=>respinge(new Error('rapoarte.imagine_invalida'));cititor.readAsDataURL(fisier); });
            const imagine=await new Promise<HTMLImageElement>((rezolva,respinge) => { const img=new Image();img.onload=()=>rezolva(img);img.onerror=()=>respinge(new Error('rapoarte.imagine_invalida'));img.src=adresa; });
            if(imagine.naturalWidth>4096 || imagine.naturalHeight>4096 || imagine.naturalWidth*imagine.naturalHeight>4_000_000) throw new Error('rapoarte.imagine_prea_mare');
            seteazaImagine(bloc,zona,{nume:fisier.name,date:adresa.slice(adresa.indexOf(',')+1),inaltime_mm:model[bloc].imagini?.[zona]?.inaltime_mm || 15});
        } catch(e) { setValidare(t(e instanceof Error && e.message.startsWith('rapoarte.') ? e.message : 'rapoarte.imagine_invalida')); }
        finally { importInCurs.current=false;setCitireImagine(false); }
    }
    const zone = (bloc: 'antet'|'subsol') => <div className="zone-raport">{(['stanga','centru','dreapta'] as const).map(z => {
        const imagine=model[bloc].imagini?.[z];
        return <div className="zona-model-raport" key={z}>
            <label>{t(`rapoarte.${z}`)}<textarea name={`${bloc}_${z}`} rows={3} maxLength={1000} value={model[bloc][z]?.[limba] || ''} onChange={e => textBloc(bloc,z,e.target.value)}/></label>
            <label>{t('rapoarte.imagine_zona',{zona:t(`rapoarte.${z}`)})}<input name={`${bloc}_imagine_${z}`} type="file" accept="image/png,image/jpeg" onChange={e=>{importaImagine(bloc,z,e.target.files?.[0]);e.target.value='';}}/></label>
            {imagine && <div className="imagine-model-raport"><img src={adresaImagine(imagine)} alt={imagine.nume}/><span>{imagine.nume}</span><label>{t('rapoarte.inaltime_imagine')}<input name={`${bloc}_inaltime_${z}`} type="number" min="5" max="35" required value={imagine.inaltime_mm} onChange={e=>seteazaImagine(bloc,z,{...imagine,inaltime_mm:Number(e.target.value)})}/></label><button type="button" className="legatura" onClick={()=>seteazaImagine(bloc,z)}>{t('rapoarte.elimina_imagine',{zona:t(`rapoarte.${z}`)})}</button></div>}
        </div>;
    })}</div>;
    const previewBloc = (bloc: 'antet'|'subsol') => <div className="zone-previzualizare">{(['stanga','centru','dreapta'] as const).map(z => {const imagine=model[bloc].imagini?.[z];return <div key={z}>{imagine && <img className="imagine-structura" style={{height:`${imagine.inaltime_mm}mm`}} src={adresaImagine(imagine)} alt={imagine.nume}/>}<span style={{display:'block',...stilCSS(model.continut.prezentare?.stiluri?.[bloc+'_'+z])}}>{textIn(model[bloc][z],limba)}</span></div>;})}</div>;
    // In editorul vizual pagina se construieste in pdfme; formularul pastreaza doar ce are efect (UIR-02).
    const vizual = model.continut.flexibil?.mod_aspect === 'PDFME';
    return <Dialog titlu={t(initial.id ? 'rapoarte.editeaza_model':'rapoarte.model_nou')} ocupat={ocupat} inchide={renunta}><form className="formular-model-raport" onSubmit={e => { e.preventDefault(); if(esteFlexibil(model.continut.sursa)&&validatFlex!==JSON.stringify(model)&&!doarStare){setValidare(t('rapoarte.flex_previz_obligatorie'));return;} if (!Object.values(model.denumire).some(v => v.trim())) { setValidare(t('rapoarte.nume_invalid'));return; } if (!model.continut.coloane.length) { setValidare(t('rapoarte.continut_invalid'));return; } setValidare('');salveaza(new rapoarte.Definitie(model)); }}>
        <div className="corp-dialog">{(eroare || validare) && <p className="mesaj eroare" role="alert">{validare || eroare}</p>}
            <fieldset className="campuri" disabled={ocupat}><label>{t('sarbatori.denumire_ro')} *<input name="denumire_ro" maxLength={100} value={model.denumire.ro || ''} onChange={e => setModel(m => ({...m,denumire:{...m.denumire,ro:e.target.value}}))}/></label><label>{t('sarbatori.denumire_en')}<input name="denumire_en" maxLength={100} value={model.denumire.en || ''} onChange={e => setModel(m => ({...m,denumire:{...m.denumire,en:e.target.value}}))}/></label>
                <label>{t('rapoarte.limba_texte')}<select name="limba_texte" value={limba} onChange={e => setLimba(e.target.value)}><option value="ro">{t('rapoarte.ro')}</option><option value="en">{t('rapoarte.en')}</option></select></label><label className="bifa"><input name="activ_model" type="checkbox" checked={model.activ} onChange={e => setModel(m => ({...m,activ:e.target.checked}))}/>{t('rapoarte.model_activ')}</label></fieldset>
            {vizual && <p className="nota-calendar">{t('rapoarte.pdfme_texte_formular')}</p>}
            <p className="variabile-raport"><strong>{t('rapoarte.variabile')}: </strong>{configuratie.variabile.map(v => <code key={v}>{`{{${v}}}`}</code>)}</p>
            <p className="nota-calendar">{t('rapoarte.imagini_indiciu')}</p>{citireImagine && <p role="status">{t('ui.incarcare')}</p>}
            <fieldset className="componenta-raport" disabled={ocupat}><legend>1 · {t('rapoarte.antet')}</legend>{zone('antet')}</fieldset>
            <fieldset className="componenta-raport" disabled={ocupat}><legend>2 · {t('rapoarte.continut')}</legend><label>{t('rapoarte.sursa')}<select name="sursa" value={model.continut.sursa} onChange={e => { const cod=e.target.value,coloane=configuratie.surse.find(s => s.cod===cod)!.coloane;setModel(m => cuPagina({...m,continut:{...m.continut,sursa:cod,coloane:[...coloane],flexibil:esteFlexibil(cod)?m.continut.flexibil||flexibilInitial():undefined}},{...m.pagina,format:cod.startsWith('PONTAJ_')?'A3':'A4',font_pt:cod==='PONTAJ_DETALIAT'?6:m.pagina.font_pt,orientare:cod==='FISA_INDIVIDUALA'?'PORTRET':'PEISAJ'})); }}>{configuratie.surse.filter(s => esteFlexibil(s.cod) || s.cod === model.continut.sursa).map(s => <option key={s.cod} value={s.cod}>{t(`rapoarte.sursa_${s.cod}`)}</option>)}</select></label>{esteFlexibil(model.continut.sursa)&&<EditorFlexibil angajati={angajati} model={model} schimba={flexibil=>setModel(m=>({...m,continut:{...m.continut,flexibil}}))} validat={()=>setValidatFlex(JSON.stringify(model))} invalideaza={()=>setValidatFlex('')}/>}
                <label className="bifa"><input type="checkbox" name="afiseaza_explicatii" checked={!model.continut.ascunde_explicatii} onChange={e=>setModel(m=>({...m,continut:{...m.continut,ascunde_explicatii:!e.target.checked}}))}/>{t('rapoarte.afiseaza_explicatii')}</label>
                <p className="nota-calendar">{t('rapoarte.explicatii_indiciu')}</p>
                <label>{t('rapoarte.introducere')}<textarea name="introducere" rows={2} maxLength={1000} value={model.continut.text?.[limba] || ''} onChange={e => setModel(m => ({...m,continut:{...m.continut,text:{...m.continut.text,[limba]:e.target.value}}}))}/></label>
</fieldset>
            <fieldset className="componenta-raport" disabled={ocupat}><legend>3 · {t('rapoarte.subsol')}</legend>{zone('subsol')}</fieldset>
            <EditorPrezentare disabled={ocupat} doarEtichete={vizual} value={model.continut.prezentare} schimba={prezentare=>setModel(m=>({...m,continut:{...m.continut,prezentare}}))}/>
            <fieldset className="componenta-raport" disabled={ocupat}><legend>{t('rapoarte.pagina')}</legend><div className="campuri"><label>{t('rapoarte.format')}<select name="format" value={model.pagina.format} onChange={e => setModel(m => cuPagina(m,{...m.pagina,format:e.target.value}))}><option value="A4">A4</option><option value="A3">A3</option></select></label><label>{t('rapoarte.orientare')}<select name="orientare" value={model.pagina.orientare} onChange={e => setModel(m => cuPagina(m,{...m.pagina,orientare:e.target.value}))}>{['PORTRET','PEISAJ'].map(o => <option key={o} value={o}>{t(`rapoarte.${o}`)}</option>)}</select></label>
                <label>{t('rapoarte.margine')}<input name="margine_mm" type="number" min="5" max="25" required value={model.pagina.margine_mm} onChange={e => setModel(m => ({...m,pagina:{...m.pagina,margine_mm:Number(e.target.value)}}))}/></label><label>{t('rapoarte.font')}<input name="font_pt" type="number" min="6" max="14" required value={model.pagina.font_pt} onChange={e => setModel(m => ({...m,pagina:{...m.pagina,font_pt:Number(e.target.value)}}))}/></label></div></fieldset>
            {!vizual && <div className="structura-model" aria-label={t('rapoarte.previzualizare')}><h3>{t('rapoarte.antet')}</h3>{previewBloc('antet')}<h3>{t('rapoarte.continut')}</h3><p>{textIn(model.continut.text,limba)}</p><div className="mostra-coloane">{model.continut.coloane.map(c => <span key={c}>{t(`rapoarte.coloana_${c}`)}</span>)}</div><h3>{t('rapoarte.subsol')}</h3>{previewBloc('subsol')}</div>}
        </div><footer className="subsol-dialog"><button type="button" className="buton secundar" disabled={ocupat} onClick={renunta}>{t('ui.anuleaza')}</button><button className="buton primar" disabled={ocupat}>{t('ui.salveaza')}</button></footer>
    </form></Dialog>;
}

// Fisa individuala cere un angajat, si cand e folosita din SQL (tabelul fisa_individuala, D090).
const ceAngajat=(m?:rapoarte.Definitie|null)=>m?.continut.sursa==='FISA_INDIVIDUALA'||/fisa_individuala/i.test(m?.continut.flexibil?.sql||'');

export function GenerareRapoarte({date,initial,onPerioada}: {date: datebaza.Date;initial?:{an:number;luna:number};onPerioada?:(p:{an:number;luna:number})=>void}) {
    const {t,i18n} = useTranslation();
    const {configuratie,incarcare,eroare,setEroare,reincarca} = useModele();
    const [filtruProfesional,setFiltruProfesional]=useState({...filtruProfesionalGol});
    const [selectie,setSelectie] = useState({model:'',an:initial?.an??new Date().getFullYear(),luna:initial?.luna??new Date().getMonth()+1,sectie:'',departament:'',angajat:'',limba:limbaRapoarte()});
    useEffect(()=>onPerioada?.({an:selectie.an,luna:selectie.luna}),[selectie.an,selectie.luna,onPerioada]);
    const [ocupat,setOcupat] = useState(false);
 const [stareCopie,setStareCopie]=useState({modificat:false,ocupat:false});
    const blocat = useRef(false);
    const [rezultat,setRezultat] = useState<main.RaspunsRaport | null>(null);
    const [url,setURL] = useState('');
    const [mesaj,setMesaj] = useState('');
    const modele = (configuratie?.modele || []).filter(d => d.activ);
    const model = modele.find(d => d.id===(selectie.model || modele[0]?.id));
    const persoaneRaport=date.angajati.filter(a=>corespundeProfesional(a,filtruProfesional)&&(!selectie.sectie||a.sectie_id===selectie.sectie)&&(!selectie.departament||a.departament_id===selectie.departament));
    useEffect(()=>{if(selectie.angajat&&!persoaneRaport.some(a=>a.id===selectie.angajat)){setSelectie(s=>({...s,angajat:''}));setRezultat(null);setMesaj(t('ux30.angajat_resetat'));}},[date,filtruProfesional,selectie.angajat,selectie.sectie,selectie.departament]);

    function seteaza<K extends keyof typeof selectie>(k: K,v: typeof selectie[K]) { setSelectie(s => ({...s,[k]:v,...(k==='sectie'?{departament:'',angajat:''}:k==='departament'?{angajat:''}:{})}));setRezultat(null);setEroare('');setMesaj(''); }
    useEffect(() => { if (!rezultat?.pdf) { setURL('');return; } const octeti=Uint8Array.from(atob(rezultat.pdf),c => c.charCodeAt(0)),nou=URL.createObjectURL(new Blob([octeti],{type:'application/pdf'}));setURL(nou);return () => URL.revokeObjectURL(nou); },[rezultat]);
    async function genereaza() {
        if (!model || blocat.current || (selectie.angajat&&!persoaneRaport.some(a=>a.id===selectie.angajat))) return;
        blocat.current = true;setOcupat(true);setEroare('');setMesaj('');setRezultat(null);
        try { const r=await completeazaRaport(await PregatesteRaport(new rapoarte.Cerere({model_id:model.id,revizie:model.revizie,an:selectie.an,luna:selectie.luna,sectie_id:selectie.sectie,departament_id:selectie.departament,...filtruProfesional,angajat_id:selectie.angajat,limba:selectie.limba})));if (r.eroare) setEroare(t(r.eroare.cod,r.eroare.param || {}));else setRezultat(r); }
        catch { setEroare(t('rapoarte.eroare')); } finally { blocat.current=false;setOcupat(false); }
    }
    async function fisier(tipareste: boolean) {
        if (!rezultat || blocat.current) return;
        blocat.current=true;setOcupat(true);setEroare('');setMesaj('');
        try { const r=await (tipareste?DeschidePDFPentruTiparire:SalveazaPDFRaport)(rezultat.token);if(r.eroare)setEroare(t(r.eroare.cod,r.eroare.param||{}));else if(!r.anulat)setMesaj(t(tipareste?'rapoarte.deschis':'rapoarte.fisier_salvat',{cale:r.cale})); }
        catch { setEroare(t('rapoarte.eroare')); } finally { blocat.current=false;setOcupat(false); }
    }
    async function exportaExcelPontaj(){if(!rezultat||blocat.current)return;blocat.current=true;setOcupat(true);setEroare('');try{const r=await SalveazaExcelRaport(rezultat.token);if(r.eroare)setEroare(t(r.eroare.cod,r.eroare.param||{}));else if(!r.anulat)setMesaj(t('rapoarte.fisier_salvat',{cale:r.cale}));}catch{setEroare(t('rapoarte.eroare'));}finally{blocat.current=false;setOcupat(false);}}
    return <section className="generare-rapoarte"><FiltreProfesionale numarRezultate={persoaneRaport.length} disabled={ocupat||stareCopie.modificat||stareCopie.ocupat} date={date} value={filtruProfesional} onChange={f=>{setFiltruProfesional(f);setRezultat(null);setMesaj('');setEroare('');}}/>{eroare && <p className="mesaj eroare" role="alert">{eroare}</p>}{mesaj && <p className="mesaj succes" role="status">{mesaj}</p>}
        <p className="nota-calendar">{t('rapoarte.paginare_indiciu')}</p>
        {incarcare ? <p role="status">{t('ui.incarcare')}</p> : <><form className="panou filtre-raport" onSubmit={e => { e.preventDefault();genereaza(); }}><fieldset disabled={ocupat||stareCopie.modificat||stareCopie.ocupat}><label className="alegere-model">{t('rapoarte.model')}<select name="model_raport" required value={model?.id || ''} onChange={e => seteaza('model',e.target.value)}>{!modele.length && <option value="">{t('rapoarte.selecteaza_model')}</option>}{modele.map(d => <option key={d.id} value={d.id}>{textIn(d.denumire,i18n.language)}</option>)}</select></label>
            <label>{t('sarbatori.an')}<input name="an_raport" type="number" min="1" max="9999" required value={selectie.an} onChange={e => seteaza('an',Number(e.target.value))}/></label>
            {model?.continut.sursa!=='ABSENTE_ANUALE' && <label>{t('pontaj.luna')}<select name="luna_raport" value={selectie.luna} onChange={e => seteaza('luna',Number(e.target.value))}>{Array.from({length:12},(_,i) => <option key={i} value={i+1}>{t(`rapoarte.luna_${String(i+1).padStart(2,'0')}`)}</option>)}</select></label>}
            <label>{t('rapoarte.limba')}<select name="limba_raport" value={selectie.limba} onChange={e => seteaza('limba',e.target.value)}><option value="ro">{t('rapoarte.ro')}</option><option value="en">{t('rapoarte.en')}</option></select></label>
            <label>{t('ui.sectie')}<select name="sectie_raport" value={selectie.sectie} onChange={e => seteaza('sectie',e.target.value)}><option value="">{t('ui.toate_sectiile')}</option>{date.sectii.map(s => <option key={s.id} value={s.id}>{denumireSectie(s,i18n.language)}</option>)}</select></label>
            <label>{t('ui.departament')}<select name="departament_raport" value={selectie.departament} onChange={e => seteaza('departament',e.target.value)}><option value="">{t('ui.toate_departamentele')}</option>{(date.departamente||[]).filter(d => !selectie.sectie || d.sectie_id===selectie.sectie).map(d => <option key={d.id} value={d.id}>{denumireDepartament(d,i18n.language)}</option>)}</select></label>
            <label>{t('ui.angajat')}<select name="angajat_raport" required={ceAngajat(model)} value={selectie.angajat} onChange={e => seteaza('angajat',e.target.value)}><option value="">{t(ceAngajat(model)?'absente.alege_angajat':'absente.toti_angajatii')}</option>{persoaneRaport.map(a => <option key={a.id} value={a.id}>{numeAngajat(a)} · {a.marca}</option>)}</select></label>
            <button className="buton primar" disabled={!model}>{t('rapoarte.genereaza')}</button>
            <button type="button" className="buton secundar" title={!rezultat ? t('motiv.genereaza_intai') : stareCopie.modificat ? t('motiv.copie_nesalvata') : t('pontaj.export_excel_titlu')} disabled={ocupat||stareCopie.modificat||stareCopie.ocupat||!rezultat} onClick={exportaExcelPontaj}>📊 {t('pontaj.export_excel')}</button></fieldset></form>
            {!modele.length && <p className="fara-rezultate">{t('rapoarte.lista_goala')}</p>}<button className="legatura" disabled={ocupat||stareCopie.modificat||stareCopie.ocupat} onClick={() => { setRezultat(null);reincarca(); }}>{t('rapoarte.reincarca')}</button>
            {rezultat?.document && <div className="rezultat-raport"><div className="unelte-calendar"><p role="status">{t('rapoarte.generat',{randuri:rezultat.document.randuri.length,versiune:rezultat.document.versiune_raport,revizie:rezultat.document.revizie_model})} · {t('rapoarte.numar_pagini',{numar:rezultat.pagini})}</p><div className="actiuni-calendar"><button className="buton secundar" disabled={ocupat||stareCopie.modificat||stareCopie.ocupat} onClick={() => fisier(false)}>{t('rapoarte.salveaza_pdf')}</button><button className="buton primar" disabled={ocupat||stareCopie.modificat||stareCopie.ocupat} onClick={() => fisier(true)}>{t('rapoarte.tipareste')}</button></div></div><p className="nota-calendar">{t('rapoarte.tiparire_indiciu')}</p><EditorCopieRaport key={rezultat.token} rezultat={rezultat} disabled={ocupat} onActualizat={setRezultat} onStare={(modificat,ocupat)=>setStareCopie(v=>v.modificat===modificat&&v.ocupat===ocupat?v:{modificat,ocupat})}/>{url && <iframe className="previzualizare-pdf" title={t('rapoarte.previzualizare')} src={`${url}#view=FitH`}/>}</div>}
        </>}
    </section>;
}

// F09-T32: pagina din Setări de pagină este pagina raportului; sablonul pdfme (daca exista) o urmeaza.
function cuPagina<T extends DateSimple<rapoarte.Definitie>>(m: T, pagina: T['pagina']): T {
    const f = m.continut.flexibil;
    if (!f?.sablon || f.mod_aspect !== 'PDFME') return {...m, pagina};
    let sablon = f.sablon;
    try { sablon = sablonPePagina(f.sablon, pagina); } catch { /* sablon invalid: verificarea il semnaleaza */ }
    return sablon === f.sablon ? {...m, pagina} : {...m, pagina, continut: {...m.continut, flexibil: {...f, sablon}}};
}

