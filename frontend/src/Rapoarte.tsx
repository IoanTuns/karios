import {useEffect, useRef, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {ListaModeleRapoarte, SalveazaModelRaport, PregatesteRaport, SalveazaPDFRaport, DeschidePDFPentruTiparire, ExportaPontajExcel} from '../wailsjs/go/main/App';
import {datebaza, main, rapoarte} from '../wailsjs/go/models';
import {Dialog, denumireDepartament, denumireSectie} from './Formulare';

const copiaza = <T,>(v: T): T => JSON.parse(JSON.stringify(v));
type DateSimple<T> = T extends Array<infer V> ? DateSimple<V>[] : T extends object ? {[K in keyof T as T[K] extends Function ? never : K]: DateSimple<T[K]>} : T;
const textIn = (v: {[limba: string]: string} | undefined, limba: string) => v?.[limba] || v?.ro || Object.values(v || {}).find(Boolean) || '';
const adresaImagine = (i: DateSimple<rapoarte.Imagine>) => `data:image/${i.date.startsWith('iVBOR') ? 'png' : 'jpeg'};base64,${i.date}`;

function useModele() {
    const {t} = useTranslation();
    const [configuratie, setConfiguratie] = useState<rapoarte.Configuratie | null>(null);
    const [incarcare, setIncarcare] = useState(true);
    const [eroare, setEroare] = useState('');
    const [revizie, setRevizie] = useState(0);
    function aplica(r: main.RaspunsModeleRapoarte) { if (r.eroare) { setEroare(t(r.eroare.cod,r.eroare.param || {})); return false; } if (!r.configuratie) { setEroare(t('rapoarte.eroare')); return false; } setConfiguratie(r.configuratie); return true; }
    useEffect(() => { let activ = true; setIncarcare(true); setEroare(''); ListaModeleRapoarte().then(r => { if (activ) aplica(r); }).catch(() => { if (activ) setEroare(t('rapoarte.eroare')); }).finally(() => { if (activ) setIncarcare(false); }); return () => { activ = false; }; },[revizie]);
    return {configuratie,incarcare,eroare,setEroare,aplica,reincarca:() => setRevizie(r => r+1)};
}

export function DefinireRapoarte() {
    const {t,i18n} = useTranslation();
    const {configuratie,incarcare,eroare,setEroare,aplica,reincarca} = useModele();
    const [editare,setEditare] = useState<rapoarte.Definitie | null>(null);
    const [ocupat,setOcupat] = useState(false);
    const blocat = useRef(false);
    const [mesaj,setMesaj] = useState('');
    const modele = configuratie?.modele || [];
    function nou(d?: rapoarte.Definitie) {
        const baza = d || modele[0]; if (!baza) return;
        setEroare(''); setMesaj('');
        setEditare(new rapoarte.Definitie({...copiaza(baza),id:'',revizie:'',activ:true,denumire:d ? copiaza(d.denumire) : {ro:'',en:''}}));
    }
    async function salveaza(d: rapoarte.Definitie) {
        if (blocat.current) return;
        blocat.current = true; setOcupat(true); setEroare('');
        try { if (aplica(await SalveazaModelRaport(d))) { setEditare(null); setMesaj(t('ui.salvat')); } }
        catch { setEroare(t('rapoarte.eroare')); }
        finally { blocat.current = false; setOcupat(false); }
    }
    return <section className="definitii-rapoarte">
        <div className="unelte-calendar"><p className="nota-calendar">{t('rapoarte.modele_indiciu')}</p><div className="actiuni-calendar"><button className="legatura" disabled={ocupat || incarcare} onClick={reincarca}>{t('rapoarte.reincarca')}</button><button className="buton primar" disabled={ocupat || incarcare || !modele.length} onClick={() => nou()}>{t('rapoarte.model_nou')}</button></div></div>
        {mesaj && <p className="mesaj succes" role="status">{mesaj}</p>}{eroare && !editare && <p className="mesaj eroare" role="alert">{eroare}</p>}
        {incarcare ? <p role="status">{t('ui.incarcare')}</p> : <div className="panou tabel-derulant"><table className="tabel-modele-rapoarte"><thead><tr>{['schimburi.denumire','rapoarte.sursa','rapoarte.revizie','rapoarte.pagina','ui.stare','ui.actiuni'].map(k => <th key={k}>{t(k)}</th>)}</tr></thead><tbody>
            {modele.map(d => <tr key={d.id}><td><strong>{textIn(d.denumire,i18n.language)}</strong></td><td>{t(`rapoarte.sursa_${d.continut.sursa}`)}</td><td>{d.revizie}</td><td>{d.pagina.format} · {t(`rapoarte.${d.pagina.orientare}`)}</td><td>{t(d.activ ? 'ui.activ':'ui.inactiv')}</td><td><div className="actiuni-rand"><button className="legatura" disabled={ocupat} onClick={() => { setEroare(''); setEditare(copiaza(d)); }}>{t('ui.editeaza')}</button><button className="legatura discreta" disabled={ocupat} onClick={() => nou(d)}>{t('rapoarte.duplica')}</button></div></td></tr>)}
        </tbody></table>{!modele.length && <p className="fara-rezultate">{t('rapoarte.lista_goala')}</p>}</div>}
        {!incarcare && !configuratie && <button className="buton secundar" onClick={reincarca}>{t('ui.reincearca')}</button>}
        {editare && configuratie && <EditorModel key={editare.id || 'nou'} initial={editare} configuratie={configuratie} ocupat={ocupat} eroare={eroare} inchide={() => { setEditare(null);setEroare(''); }} salveaza={salveaza}/>}
    </section>;
}

function EditorModel({initial,configuratie,ocupat:salvareInCurs,eroare,inchide,salveaza}: {initial: rapoarte.Definitie; configuratie: rapoarte.Configuratie; ocupat: boolean; eroare: string; inchide: () => void; salveaza: (d: rapoarte.Definitie) => Promise<void>}) {
    const {t,i18n} = useTranslation();
    const [model,setModel] = useState<DateSimple<rapoarte.Definitie>>(copiaza(initial));
    const [limba,setLimba] = useState(i18n.language === 'en' ? 'en' : 'ro');
    const [validare,setValidare] = useState('');
    const [citireImagine,setCitireImagine] = useState(false);
    const importInCurs = useRef(false);
    const ocupat = salvareInCurs || citireImagine;
    const sursa = configuratie.surse.find(s => s.cod === model.continut.sursa)!;
    const modificate = JSON.stringify(model) !== JSON.stringify(initial);
    useEffect(() => { if (!modificate) return; const atentioneaza = (e: BeforeUnloadEvent) => { e.preventDefault();e.returnValue=''; }; window.addEventListener('beforeunload',atentioneaza);return () => window.removeEventListener('beforeunload',atentioneaza); },[modificate]);
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
    function selecteazaColoana(cod: string,bifat: boolean) { setModel(m => ({...m,continut:{...m.continut,coloane:bifat ? [...m.continut.coloane.filter(c=>c!=='totaluri'),cod,...(m.continut.coloane.includes('totaluri')?['totaluri']:[])] : m.continut.coloane.filter(c => c!==cod)}})); }
    function muta(cod: string,sens: number) { setModel(m => { const coloane=[...m.continut.coloane],i=coloane.indexOf(cod),j=i+sens;if (i<0 || j<0 || j>=coloane.length || cod==='totaluri' || coloane[j]==='totaluri') return m;[coloane[i],coloane[j]]=[coloane[j],coloane[i]];return {...m,continut:{...m.continut,coloane}}; }); }
    const zone = (bloc: 'antet'|'subsol') => <div className="zone-raport">{(['stanga','centru','dreapta'] as const).map(z => {
        const imagine=model[bloc].imagini?.[z];
        return <div className="zona-model-raport" key={z}>
            <label>{t(`rapoarte.${z}`)}<textarea name={`${bloc}_${z}`} rows={3} maxLength={1000} value={model[bloc][z]?.[limba] || ''} onChange={e => textBloc(bloc,z,e.target.value)}/></label>
            <label>{t('rapoarte.imagine_zona',{zona:t(`rapoarte.${z}`)})}<input name={`${bloc}_imagine_${z}`} type="file" accept="image/png,image/jpeg" onChange={e=>{importaImagine(bloc,z,e.target.files?.[0]);e.target.value='';}}/></label>
            {imagine && <div className="imagine-model-raport"><img src={adresaImagine(imagine)} alt={imagine.nume}/><span>{imagine.nume}</span><label>{t('rapoarte.inaltime_imagine')}<input name={`${bloc}_inaltime_${z}`} type="number" min="5" max="35" required value={imagine.inaltime_mm} onChange={e=>seteazaImagine(bloc,z,{...imagine,inaltime_mm:Number(e.target.value)})}/></label><button type="button" className="legatura" onClick={()=>seteazaImagine(bloc,z)}>{t('rapoarte.elimina_imagine',{zona:t(`rapoarte.${z}`)})}</button></div>}
        </div>;
    })}</div>;
    const previewBloc = (bloc: 'antet'|'subsol') => <div className="zone-previzualizare">{(['stanga','centru','dreapta'] as const).map(z => {const imagine=model[bloc].imagini?.[z];return <div key={z}>{imagine && <img className="imagine-structura" style={{height:`${imagine.inaltime_mm}mm`}} src={adresaImagine(imagine)} alt={imagine.nume}/>}<span>{textIn(model[bloc][z],limba)}</span></div>;})}</div>;
    return <Dialog titlu={t(initial.id ? 'rapoarte.editeaza_model':'rapoarte.model_nou')} ocupat={ocupat} inchide={renunta}><form className="formular-model-raport" onSubmit={e => { e.preventDefault(); if (!Object.values(model.denumire).some(v => v.trim())) { setValidare(t('rapoarte.nume_invalid'));return; } if (!model.continut.coloane.length) { setValidare(t('rapoarte.continut_invalid'));return; } setValidare('');salveaza(new rapoarte.Definitie(model)); }}>
        <div className="corp-dialog">{(eroare || validare) && <p className="mesaj eroare" role="alert">{validare || eroare}</p>}
            <fieldset className="campuri" disabled={ocupat}><label>{t('sarbatori.denumire_ro')} *<input name="denumire_ro" maxLength={100} value={model.denumire.ro || ''} onChange={e => setModel(m => ({...m,denumire:{...m.denumire,ro:e.target.value}}))}/></label><label>{t('sarbatori.denumire_en')}<input name="denumire_en" maxLength={100} value={model.denumire.en || ''} onChange={e => setModel(m => ({...m,denumire:{...m.denumire,en:e.target.value}}))}/></label>
                <label>{t('rapoarte.limba_texte')}<select name="limba_texte" value={limba} onChange={e => setLimba(e.target.value)}><option value="ro">{t('rapoarte.ro')}</option><option value="en">{t('rapoarte.en')}</option></select></label><label className="bifa"><input name="activ_model" type="checkbox" checked={model.activ} onChange={e => setModel(m => ({...m,activ:e.target.checked}))}/>{t('rapoarte.model_activ')}</label></fieldset>
            <p className="variabile-raport"><strong>{t('rapoarte.variabile')}: </strong>{configuratie.variabile.map(v => <code key={v}>{`{{${v}}}`}</code>)}</p>
            <p className="nota-calendar">{t('rapoarte.imagini_indiciu')}</p>{citireImagine && <p role="status">{t('ui.incarcare')}</p>}
            <fieldset className="componenta-raport" disabled={ocupat}><legend>1 · {t('rapoarte.antet')}</legend>{zone('antet')}</fieldset>
            <fieldset className="componenta-raport" disabled={ocupat}><legend>2 · {t('rapoarte.continut')}</legend><label>{t('rapoarte.sursa')}<select name="sursa" value={model.continut.sursa} onChange={e => { const cod=e.target.value,coloane=configuratie.surse.find(s => s.cod===cod)!.coloane;setModel(m => ({...m,continut:{...m.continut,sursa:cod,coloane:[...coloane]},pagina:{...m.pagina,format:cod.startsWith('PONTAJ_')?'A3':'A4',font_pt:cod==='PONTAJ_DETALIAT'?6:m.pagina.font_pt,orientare:cod==='FISA_INDIVIDUALA'?'PORTRET':'PEISAJ'}})); }}>{configuratie.surse.map(s => <option key={s.cod} value={s.cod}>{t(`rapoarte.sursa_${s.cod}`)}</option>)}</select></label>
                <label>{t('rapoarte.introducere')}<textarea name="introducere" rows={2} maxLength={1000} value={model.continut.text?.[limba] || ''} onChange={e => setModel(m => ({...m,continut:{...m.continut,text:{...m.continut.text,[limba]:e.target.value}}}))}/></label>
                <p>{t('rapoarte.coloane')}</p><div className="coloane-raport">{[...model.continut.coloane,...sursa.coloane.filter(c => !model.continut.coloane.includes(c))].map(c => { const i=model.continut.coloane.indexOf(c);return <div key={c}><label><input type="checkbox" name={`coloana_${c}`} checked={i>=0} disabled={(model.continut.sursa==='PONTAJ_LUNAR' && c==='angajat') || (model.continut.sursa==='PONTAJ_DETALIAT' && ['angajat','zile','totaluri'].includes(c))} onChange={e => selecteazaColoana(c,e.target.checked)}/>{t(`rapoarte.coloana_${c}`)}</label><button type="button" className="legatura" aria-label={t('rapoarte.muta_sus',{coloana:t(`rapoarte.coloana_${c}`)})} disabled={i<=0 || c==='totaluri'} onClick={() => muta(c,-1)}>↑</button><button type="button" className="legatura" aria-label={t('rapoarte.muta_jos',{coloana:t(`rapoarte.coloana_${c}`)})} disabled={i<0 || i===model.continut.coloane.length-1 || model.continut.coloane[i+1]==='totaluri'} onClick={() => muta(c,1)}>↓</button></div>; })}</div></fieldset>
            <fieldset className="componenta-raport" disabled={ocupat}><legend>3 · {t('rapoarte.subsol')}</legend>{zone('subsol')}</fieldset>
            <fieldset className="componenta-raport" disabled={ocupat}><legend>{t('rapoarte.pagina')}</legend><div className="campuri"><label>{t('rapoarte.format')}<select name="format" value={model.pagina.format} onChange={e => setModel(m => ({...m,pagina:{...m.pagina,format:e.target.value}}))}><option value="A4">A4</option><option value="A3">A3</option></select></label><label>{t('rapoarte.orientare')}<select name="orientare" value={model.pagina.orientare} onChange={e => setModel(m => ({...m,pagina:{...m.pagina,orientare:e.target.value}}))}>{['PORTRET','PEISAJ'].map(o => <option key={o} value={o}>{t(`rapoarte.${o}`)}</option>)}</select></label>
                <label>{t('rapoarte.margine')}<input name="margine_mm" type="number" min="5" max="25" required value={model.pagina.margine_mm} onChange={e => setModel(m => ({...m,pagina:{...m.pagina,margine_mm:Number(e.target.value)}}))}/></label><label>{t('rapoarte.font')}<input name="font_pt" type="number" min="6" max="14" required value={model.pagina.font_pt} onChange={e => setModel(m => ({...m,pagina:{...m.pagina,font_pt:Number(e.target.value)}}))}/></label></div></fieldset>
            <div className="structura-model" aria-label={t('rapoarte.previzualizare')}><h3>{t('rapoarte.antet')}</h3>{previewBloc('antet')}<h3>{t('rapoarte.continut')}</h3><p>{textIn(model.continut.text,limba)}</p><div className="mostra-coloane">{model.continut.coloane.map(c => <span key={c}>{t(`rapoarte.coloana_${c}`)}</span>)}</div><h3>{t('rapoarte.subsol')}</h3>{previewBloc('subsol')}</div>
        </div><footer className="subsol-dialog"><button type="button" className="buton secundar" disabled={ocupat} onClick={renunta}>{t('ui.anuleaza')}</button><button className="buton primar" disabled={ocupat}>{t('ui.salveaza')}</button></footer>
    </form></Dialog>;
}

export function GenerareRapoarte({date}: {date: datebaza.Date}) {
    const {t,i18n} = useTranslation();
    const {configuratie,incarcare,eroare,setEroare,reincarca} = useModele();
    const [selectie,setSelectie] = useState({model:'',an:new Date().getFullYear(),luna:new Date().getMonth()+1,sectie:'',departament:'',angajat:'',limba:i18n.language==='en'?'en':'ro'});
    const [ocupat,setOcupat] = useState(false);
    const blocat = useRef(false);
    const [rezultat,setRezultat] = useState<main.RaspunsRaport | null>(null);
    const [url,setURL] = useState('');
    const [mesaj,setMesaj] = useState('');
    const modele = (configuratie?.modele || []).filter(d => d.activ);
    const model = modele.find(d => d.id===(selectie.model || modele[0]?.id));
    function seteaza<K extends keyof typeof selectie>(k: K,v: typeof selectie[K]) { setSelectie(s => ({...s,[k]:v,...(k==='sectie'?{departament:'',angajat:''}:k==='departament'?{angajat:''}:{})}));setRezultat(null);setEroare('');setMesaj(''); }
    useEffect(() => { if (!rezultat?.pdf) { setURL('');return; } const octeti=Uint8Array.from(atob(rezultat.pdf),c => c.charCodeAt(0)),nou=URL.createObjectURL(new Blob([octeti],{type:'application/pdf'}));setURL(nou);return () => URL.revokeObjectURL(nou); },[rezultat]);
    async function genereaza() {
        if (!model || blocat.current) return;
        blocat.current = true;setOcupat(true);setEroare('');setMesaj('');setRezultat(null);
        try { const r=await PregatesteRaport(new rapoarte.Cerere({model_id:model.id,revizie:model.revizie,an:selectie.an,luna:selectie.luna,sectie_id:selectie.sectie,departament_id:selectie.departament,angajat_id:selectie.angajat,limba:selectie.limba}));if (r.eroare) setEroare(t(r.eroare.cod,r.eroare.param || {}));else setRezultat(r); }
        catch { setEroare(t('rapoarte.eroare')); } finally { blocat.current=false;setOcupat(false); }
    }
    async function fisier(tipareste: boolean) {
        if (!rezultat || blocat.current) return;
        blocat.current=true;setOcupat(true);setEroare('');setMesaj('');
        try { const r=await (tipareste?DeschidePDFPentruTiparire:SalveazaPDFRaport)(rezultat.token);if(r.eroare)setEroare(t(r.eroare.cod,r.eroare.param||{}));else if(!r.anulat)setMesaj(t(tipareste?'rapoarte.deschis':'rapoarte.fisier_salvat',{cale:r.cale})); }
        catch { setEroare(t('rapoarte.eroare')); } finally { blocat.current=false;setOcupat(false); }
    }
    async function exportaExcelPontaj() {
        if (blocat.current) return;
        blocat.current=true;setOcupat(true);setEroare('');setMesaj('');
        try {
            const r = await ExportaPontajExcel(selectie.an, selectie.luna);
            if (r.eroare) setEroare(t(r.eroare.cod, r.eroare.param||{}));
            else if (r.cale) setMesaj(t('rapoarte.fisier_salvat', {cale: r.cale}));
            else if (!r.anulat && r.continut) {
                const octeti = Uint8Array.from(atob(r.continut), c => c.charCodeAt(0));
                const blob = new Blob([octeti], {type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'});
                const urlFisier = URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = urlFisier;
                a.download = r.nume_fisier || `Pontix_Pontaj_${selectie.an}_${String(selectie.luna).padStart(2, '0')}.xlsx`;
                document.body.appendChild(a);
                a.click();
                document.body.removeChild(a);
                URL.revokeObjectURL(urlFisier);
                setMesaj(t('pontaj.export_excel_succes'));
            }
        } catch { setEroare(t('rapoarte.eroare')); }
        finally { blocat.current=false;setOcupat(false); }
    }
    return <section className="generare-rapoarte">{eroare && <p className="mesaj eroare" role="alert">{eroare}</p>}{mesaj && <p className="mesaj succes" role="status">{mesaj}</p>}
        <p className="nota-calendar">{t('rapoarte.paginare_indiciu')}</p>
        {incarcare ? <p role="status">{t('ui.incarcare')}</p> : <><form className="panou filtre-raport" onSubmit={e => { e.preventDefault();genereaza(); }}><fieldset disabled={ocupat}><label className="alegere-model">{t('rapoarte.model')}<select name="model_raport" required value={model?.id || ''} onChange={e => seteaza('model',e.target.value)}>{!modele.length && <option value="">{t('rapoarte.selecteaza_model')}</option>}{modele.map(d => <option key={d.id} value={d.id}>{textIn(d.denumire,i18n.language)}</option>)}</select></label>
            <label>{t('sarbatori.an')}<input name="an_raport" type="number" min="1" max="9999" required value={selectie.an} onChange={e => seteaza('an',Number(e.target.value))}/></label>
            {model?.continut.sursa!=='ABSENTE_ANUALE' && <label>{t('pontaj.luna')}<select name="luna_raport" value={selectie.luna} onChange={e => seteaza('luna',Number(e.target.value))}>{Array.from({length:12},(_,i) => <option key={i} value={i+1}>{t(`rapoarte.luna_${String(i+1).padStart(2,'0')}`)}</option>)}</select></label>}
            <label>{t('rapoarte.limba')}<select name="limba_raport" value={selectie.limba} onChange={e => seteaza('limba',e.target.value)}><option value="ro">{t('rapoarte.ro')}</option><option value="en">{t('rapoarte.en')}</option></select></label>
            <label>{t('ui.sectie')}<select name="sectie_raport" value={selectie.sectie} onChange={e => seteaza('sectie',e.target.value)}><option value="">{t('ui.toate_sectiile')}</option>{date.sectii.map(s => <option key={s.id} value={s.id}>{denumireSectie(s,i18n.language)}</option>)}</select></label>
            <label>{t('ui.departament')}<select name="departament_raport" value={selectie.departament} onChange={e => seteaza('departament',e.target.value)}><option value="">{t('ui.toate_departamentele')}</option>{(date.departamente||[]).filter(d => !selectie.sectie || d.sectie_id===selectie.sectie).map(d => <option key={d.id} value={d.id}>{denumireDepartament(d,i18n.language)}</option>)}</select></label>
            <label>{t('ui.angajat')}<select name="angajat_raport" required={model?.continut.sursa==='FISA_INDIVIDUALA'} value={selectie.angajat} onChange={e => seteaza('angajat',e.target.value)}><option value="">{t(model?.continut.sursa==='FISA_INDIVIDUALA'?'absente.alege_angajat':'absente.toti_angajatii')}</option>{date.angajati.filter(a => (!selectie.sectie||a.sectie_id===selectie.sectie)&&(!selectie.departament||a.departament_id===selectie.departament)).map(a => <option key={a.id} value={a.id}>{a.nume} {a.prenume} · {a.marca}</option>)}</select></label>
            <button className="buton primar" disabled={!model}>{t('rapoarte.genereaza')}</button>
            <button type="button" className="buton secundar" title="Export Excel (.xlsx)" disabled={ocupat} onClick={exportaExcelPontaj}>📊 {t('pontaj.export_excel')}</button></fieldset></form>
            {!modele.length && <p className="fara-rezultate">{t('rapoarte.lista_goala')}</p>}<button className="legatura" disabled={ocupat} onClick={() => { setRezultat(null);reincarca(); }}>{t('rapoarte.reincarca')}</button>
            {rezultat?.document && <div className="rezultat-raport"><div className="unelte-calendar"><p role="status">{t('rapoarte.generat',{randuri:rezultat.document.randuri.length,versiune:rezultat.document.versiune_raport,revizie:rezultat.document.revizie_model})} · {t('rapoarte.numar_pagini',{numar:rezultat.pagini})}</p><div className="actiuni-calendar"><button className="buton secundar" disabled={ocupat} onClick={() => fisier(false)}>{t('rapoarte.salveaza_pdf')}</button><button className="buton primar" disabled={ocupat} onClick={() => fisier(true)}>{t('rapoarte.tipareste')}</button></div></div><p className="nota-calendar">{t('rapoarte.tiparire_indiciu')}</p>{url && <iframe className="previzualizare-pdf" title={t('rapoarte.previzualizare')} src={`${url}#view=FitH`}/>}</div>}
        </>}
    </section>;
}
