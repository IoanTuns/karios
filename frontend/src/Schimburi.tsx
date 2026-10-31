import {ButonStergereAdmin} from './ButonStergereAdmin';
import {FiltruStare,corespundeStare,type StareLista} from './FiltruStare';
import {useEffect, useRef, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {CompleteazaSchimburiInitiale, ListaSchimburi, SalveazaSchimb, SchimbaActivSchimb} from '../wailsjs/go/main/App';
import {datebaza, main} from '../wailsjs/go/models';
import {EditorNomenclator} from './Nomenclatoare';
import {Dialog, durata} from './Formulare';

export function oraText(minute: number) { return `${String(Math.floor(minute / 60)).padStart(2, '0')}:${String(minute % 60).padStart(2, '0')}`; }


export function Schimburi({date,onDate}:{date:datebaza.Date;onDate:(date:datebaza.Date)=>void}) {
 const categorii=date.categorii_schimburi??[];
 const tipuri=date.tipuri_program??[], regimuri=date.regimuri_program??[];
    const {t, i18n} = useTranslation();
 const [stareLista,setStareLista]=useState<StareLista>('ACTIVE');
    const [lista, setLista] = useState<datebaza.Schimb[]>([]);
    const [categorie, setCategorie] = useState('');
    const [categoriePrecompletare, setCategoriePrecompletare] = useState('LAB');
    const [regim, setRegim] = useState('');
    const [incarcare, setIncarcare] = useState(true);
    const [ocupat, setOcupat] = useState(false);
    const blocat = useRef(false);
    const [eroare, setEroare] = useState('');
    const [mesaj, setMesaj] = useState('');
    const [editare, setEditare] = useState<datebaza.Schimb | null>(null);

    function aplica(r: main.RaspunsSchimburi) {
        if (r.eroare) { setEroare(t(r.eroare.cod, r.eroare.param ?? {})); return false; }
        setLista(r.schimburi ?? []); return true;
    }
    useEffect(() => {
        let activ = true;
        ListaSchimburi().then(r => { if (activ) aplica(r); }).catch(() => { if (activ) setEroare(t('date.eroare')); }).finally(() => { if (activ) setIncarcare(false); });
        return () => { activ = false; };
    }, []);

    async function executa(actiune: () => Promise<main.RaspunsSchimburi>, completare = false) {
        if (blocat.current) return;
        blocat.current = true; setOcupat(true); setEroare(''); setMesaj('');
        try { const r = await actiune(); if (aplica(r)) { setEditare(null); setMesaj(t(completare ? 'schimburi.completate' : 'ui.salvat', {numar: r.adaugate})); } }
        catch { setEroare(t('date.eroare')); }
        finally { blocat.current = false; setOcupat(false); }
    }

    const numeCategorie=(id:string)=>{const c=categorii.find(c=>c.id===id);return c?.denumire[i18n.language]||c?.denumire.ro||c?.denumire.en||id;};
    const filtrate = lista.filter(s => corespundeStare(s.activ,stareLista) && (!categorie || (s.categorie || 'CUSTOM') === categorie) && (!regim || (s.regim_id || s.regim) === regim));

    return <section className="schimburi"><details className="detalii-configurare"><summary>{t('ui28.program_explica')}</summary><p>{t('ui28.program_indiciu')}</p><button type="button" className="legatura" onClick={()=>window.dispatchEvent(new CustomEvent('pontix-ajutor-contextual',{detail:'schimburi'}))}>{t('ajutor.buton')}</button></details>
 <EditorNomenclator tip="categorii" lista={categorii} onDate={onDate}/>
 <EditorNomenclator tip="tipuri_program" lista={tipuri} onDate={onDate}/>
 <EditorNomenclator tip="regimuri_program" lista={regimuri} onDate={onDate}/>
        <div className="unelte-calendar"><FiltruStare nume="schimburi" valoare={stareLista} schimba={setStareLista} disabled={ocupat}/><p className="nota-calendar">{t('schimburi.indiciu')}</p>
            <div className="actiuni-calendar">
                <select aria-label={t('schimburi.categorie_precompletare')} value={categoriePrecompletare} disabled={ocupat || incarcare} onChange={e => setCategoriePrecompletare(e.target.value)}>
                    <option value="LAB">{numeCategorie('LAB')}</option>
                </select>
                <button className="buton secundar" disabled={ocupat || incarcare || !categorii.some(c=>c.id==='LAB'&&c.activ)} onClick={() => executa(() => CompleteazaSchimburiInitiale(categoriePrecompletare), true)}>{t('schimburi.completeaza_initiale')}</button>
                <button className="buton primar" disabled={ocupat || incarcare} onClick={() => { setEroare(''); setEditare(new datebaza.Schimb({id: '', cod: '', denumire: {ro: '', en: ''}, ora_start: -1, durata_minute: 0, pauza_minute: 0, tip: 'NORMAL', regim: 'LUCRATOARE', tip_id:'NORMAL',regim_id:'LUCRATOARE', categorie: 'CUSTOM', activ: true})); }}>{t('schimburi.adauga')}</button></div></div>
        <p className="nota-calendar">{t('schimburi.sursa_initiale')}</p>
        <p className="nota-calendar">{t('schimburi.intrare_indiciu')}</p>
        <div className="unelte">
            <label>{t('schimburi.categorie')} <select aria-label={t('schimburi.categorie')} value={categorie} disabled={ocupat} onChange={e => setCategorie(e.target.value)}><option value="">{t('schimburi.toate_categoriile')}</option>{categorii.map(c => <option key={c.id} value={c.id}>{numeCategorie(c.id)}{!c.activ?` (${t('ui.inactiv')})`:''}</option>)}</select></label>
            <label>{t('schimburi.regim')} <select aria-label={t('schimburi.regim')} value={regim} disabled={ocupat} onChange={e => setRegim(e.target.value)}><option value="">{t('schimburi.toate_regimurile')}</option>{regimuri.map(r => <option key={r.id} value={r.id}>{r.denumire[i18n.language]||r.denumire.ro}</option>)}</select></label>
        </div>
        {!editare && eroare && <p className="mesaj eroare" role="alert">{eroare}</p>}{mesaj && <p className="mesaj succes" role="status">{mesaj}</p>}
        {incarcare ? <p role="status">{t('ui.incarcare')}</p> : <div className="panou tabel-derulant"><table className="tabel-schimburi"><thead><tr>
            {['cod', 'denumire', 'categorie', 'inceput', 'sfarsit', 'durata', 'pauza', 'tip', 'regim'].map(c => <th key={c}>{t(`schimburi.${c}`)}</th>)}<th>{t('ui.stare')}</th><th>{t('ui.actiuni')}</th></tr></thead><tbody>
            {filtrate.map(s => <tr key={s.id}><td><strong>{s.cod}</strong></td><td>{s.denumire[i18n.language] || s.denumire.ro || Object.values(s.denumire).find(Boolean)}</td>
                <td>{numeCategorie(s.categorie||'CUSTOM')}</td>
                <td>{oraText(s.ora_start)}</td><td>{oraText((s.ora_start + s.durata_minute) % 1440)}{s.ora_start + s.durata_minute >= 1440 && <small className="zi-urmatoare">{t('schimburi.peste_zile', {numar: Math.floor((s.ora_start + s.durata_minute) / 1440)})}</small>}</td><td>{durata(s.durata_minute, t)}</td><td>{t('schimburi.minute', {numar: s.pauza_minute})}</td><td>{tipuri.find(t=>t.id===(s.tip_id||s.tip))?.denumire[i18n.language]||tipuri.find(t=>t.id===(s.tip_id||s.tip))?.denumire.ro||t(`schimburi.tip_${s.tip}`)}</td><td><span className={s.regim === 'INTRARE_ZI_LIBERA' ? 'eticheta-intrare' : ''}>{regimuri.find(r=>r.id===(s.regim_id||s.regim))?.denumire[i18n.language]||regimuri.find(r=>r.id===(s.regim_id||s.regim))?.denumire.ro||t(`schimburi.regim_${s.regim||'LUCRATOARE'}`)}</span></td><td>{t(s.activ ? 'ui.activ' : 'ui.inactiv')}</td>
                <td><div className="actiuni-rand"><button className="legatura" disabled={ocupat} onClick={() => { setEroare(''); setEditare(s); }}>{t('ui.editeaza')}</button>
                    <button className="legatura discreta" disabled={ocupat} onClick={() => executa(() => SchimbaActivSchimb(s.id, !s.activ))}>{t(s.activ ? 'ui.dezactiveaza' : 'ui.reactiveaza')}</button><ButonStergereAdmin tip="schimburi" id={s.id} revizie={s.revizie} activ={s.activ} nume={s.cod} disabled={ocupat||!!editare} onOcupat={setOcupat} onSters={async d=>{onDate(d);await executa(ListaSchimburi);}}/></div></td></tr>)}
        </tbody></table>{filtrate.length === 0 && <p className="fara-rezultate">{t('schimburi.gol')}</p>}</div>}
        {!incarcare && eroare && !editare && <button className="buton secundar" disabled={ocupat} onClick={() => executa(ListaSchimburi)}>{t('ui.reincearca')}</button>}
        {editare && <FormularSchimb tipuri={tipuri} regimuri={regimuri} categorii={categorii} schimb={editare} ocupat={ocupat} eroare={eroare} inchide={() => { setEditare(null); setEroare(''); }} salveaza={s => executa(() => SalveazaSchimb(s))}/>}
    </section>;
}

function FormularSchimb({tipuri,regimuri,categorii, schimb, ocupat, eroare, inchide, salveaza}: {tipuri:datebaza.ElementNomenclator[];regimuri:datebaza.ElementNomenclator[];categorii:datebaza.ElementNomenclator[];schimb: datebaza.Schimb; ocupat: boolean; eroare: string; inchide: () => void; salveaza: (s: datebaza.Schimb) => Promise<void>}) {
    const {t,i18n} = useTranslation();
    const initial = {cod: schimb.cod, ro: schimb.denumire.ro || '', en: schimb.denumire.en || '', ora: schimb.ora_start < 0 ? '' : oraText(schimb.ora_start), durata: schimb.durata_minute > 0 ? String(schimb.durata_minute) : '', pauza: String(schimb.pauza_minute), tip: schimb.tip_id||schimb.tip, regim: schimb.regim_id||schimb.regim||'LUCRATOARE', categorie: schimb.categorie || 'CUSTOM', activ: schimb.activ};
    const [campuri, setCampuri] = useState(initial);
    const [validare, setValidare] = useState('');
    function seteaza<K extends keyof typeof campuri>(cheie: K, valoare: typeof campuri[K]) { setCampuri(c => ({...c, [cheie]: valoare})); }
    function inchideCuVerificare() { if (JSON.stringify(initial) === JSON.stringify(campuri) || window.confirm(t('ui.renunta'))) inchide(); }
    return <Dialog titlu={t(schimb.id ? 'schimburi.editeaza' : 'schimburi.adauga')} ocupat={ocupat} inchide={inchideCuVerificare}>
        <form onSubmit={async e => {
            e.preventDefault(); setValidare('');
            const durata = Number(campuri.durata), pauza = Number(campuri.pauza);
            if (!campuri.ro.trim() && !campuri.en.trim()) { setValidare(t('sarbatori.denumire_obligatorie')); return; }
            if (!Number.isSafeInteger(durata) || !Number.isSafeInteger(pauza) || durata < 1 || durata > 2880 || pauza < 0 || pauza > durata || !/^\d{2}:\d{2}$/.test(campuri.ora)) { setValidare(t('schimburi.interval_invalid')); return; }
            const [ore, minute] = campuri.ora.split(':').map(Number);
            await salveaza(new datebaza.Schimb({id: schimb.id, cod: campuri.cod, denumire: {...schimb.denumire, ro: campuri.ro, en: campuri.en}, ora_start: ore * 60 + minute, durata_minute: durata, pauza_minute: pauza, tip: tipuri.find(t=>t.id===campuri.tip)?.comportament||schimb.tip, regim:regimuri.find(r=>r.id===campuri.regim)?.comportament||schimb.regim, tip_id:campuri.tip,regim_id:campuri.regim, categorie: campuri.categorie, activ: campuri.activ}));
        }}>
            <div className="corp-dialog">{(validare || eroare) && <p className="mesaj eroare" role="alert">{validare || eroare}</p>}
                <fieldset className="campuri" disabled={ocupat}>
                    <label>{t('schimburi.cod')} *<input autoFocus required maxLength={16} value={campuri.cod} onChange={e => seteaza('cod', e.target.value)}/></label>
                    <label>{t('schimburi.tip')}<select value={campuri.tip} onChange={e => seteaza('tip', e.target.value)}>{tipuri.filter(t=>t.activ||t.id===campuri.tip).map(tip => <option key={tip.id} value={tip.id}>{tip.denumire[i18n.language]||tip.denumire.ro}</option>)}</select></label>
                    <label className="latime-completa">{t('sarbatori.denumire_ro')}<input value={campuri.ro} onChange={e => seteaza('ro', e.target.value)}/></label>
                    <label className="latime-completa">{t('sarbatori.denumire_en')}<input value={campuri.en} onChange={e => seteaza('en', e.target.value)}/></label>
                    <label>{t('schimburi.inceput')} *<input type="text" name="ora_inceput" required pattern="([01][0-9]|2[0-3]):[0-5][0-9]" maxLength={5} title={t('schimburi.format_ora')} value={campuri.ora} onChange={e => seteaza('ora', e.target.value)}/></label>
                    <label>{t('schimburi.durata_minute')} *<input type="number" required min="1" max="2880" step="1" value={campuri.durata} onChange={e => seteaza('durata', e.target.value)}/></label>
                    <label>{t('schimburi.pauza_minute')}<input type="number" required min="0" max={campuri.durata || '2880'} step="1" value={campuri.pauza} onChange={e => seteaza('pauza', e.target.value)}/></label>
                    <label className="bifa"><input type="checkbox" checked={campuri.activ} onChange={e => seteaza('activ', e.target.checked)}/>{t('ui.activ_indiciu')}</label>
                    <label className="latime-completa">{t('schimburi.regim')}<select name="regim" value={campuri.regim} onChange={e => seteaza('regim', e.target.value)}>{regimuri.filter(r=>r.activ||r.id===campuri.regim).map(r => <option key={r.id} value={r.id}>{r.denumire[i18n.language]||r.denumire.ro}</option>)}</select></label>
                    <label className="latime-completa">{t('schimburi.categorie')}<select name="categorie" value={campuri.categorie} onChange={e => seteaza('categorie', e.target.value)}>{categorii.filter(c=>c.activ||c.id===campuri.categorie).map(c => <option key={c.id} value={c.id}>{c.denumire[i18n.language]||c.denumire.ro||c.denumire.en}{!c.activ?` (${t('ui.inactiv')})`:''}</option>)}</select></label>
                </fieldset>
            </div><footer className="subsol-dialog"><button type="button" className="buton secundar" disabled={ocupat} onClick={inchideCuVerificare}>{t('ui.anuleaza')}</button><button className="buton primar" disabled={ocupat}>{t(ocupat ? 'ui.salvare' : 'ui.salveaza')}</button></footer>
        </form>
    </Dialog>;
}
