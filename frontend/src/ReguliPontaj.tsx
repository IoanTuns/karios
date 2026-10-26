import {useEffect, useRef, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {CitesteConfiguratiePontaj, SalveazaConfiguratiePontaj} from '../wailsjs/go/main/App';
import {calcul, datebaza} from '../wailsjs/go/models';

const surse = ['NEDEFINITA','TOTAL_ORE','ORE_ZNL','ORE_NORMA','ORE_SCHIMBURI','ZILE_ABSENTA'];
const zile = ['CO','CEV','CFS','CM','deleg'];
const ore = ['ore_lucrate','T','ZNL','ore_norma','ore_garda'];

export function ReguliPontaj({onModificari,onOcupat}: {onModificari:(v:boolean)=>void;onOcupat:(v:boolean)=>void}) {
    const {t,i18n} = useTranslation();
    const [config,setConfig] = useState<calcul.ConfiguratiePontaj | null>(null);
    const [original,setOriginal] = useState('');
    const [optiuni,setOptiuni] = useState<datebaza.OptiuneRegula[]>([]);
    const [eroare,setEroare] = useState('');
    const [mesaj,setMesaj] = useState('');
    const [ocupat,setOcupat] = useState(false);
    const blocat = useRef(false);
    const modificat = Boolean(config && JSON.stringify(config)!==original);
    useEffect(()=>{onModificari(modificat);return ()=>onModificari(false);},[modificat,onModificari]);
    useEffect(()=>{if(!modificat)return;const avertizeaza=(e:BeforeUnloadEvent)=>{e.preventDefault();e.returnValue='';};window.addEventListener('beforeunload',avertizeaza);return()=>window.removeEventListener('beforeunload',avertizeaza);},[modificat]);
    async function incarca() {
        if(blocat.current)return;blocat.current=true;setOcupat(true);onOcupat(true);setEroare('');
        try {const r=await CitesteConfiguratiePontaj();if(r.eroare)setEroare(t(r.eroare.cod,r.eroare.param||{}));else if(r.configuratie){setConfig(r.configuratie);setOriginal(JSON.stringify(r.configuratie));setOptiuni(r.optiuni||[]);}}
        catch {setEroare(t('date.eroare'));}finally{blocat.current=false;setOcupat(false);onOcupat(false);}
    }
    useEffect(()=>{incarca();},[]);
    function schimba(i:number,v:Partial<calcul.RegulaColoana>){setConfig(c=>c ? new calcul.ConfiguratiePontaj({...c,reguli:c.reguli.map((r,j)=>j===i?{...r,...v}:r)}):c);setMesaj('');}
    async function salveaza(){
        if(!config||blocat.current)return;blocat.current=true;setOcupat(true);onOcupat(true);setEroare('');setMesaj('');
        try {const r=await SalveazaConfiguratiePontaj(config);if(r.eroare)setEroare(t(r.eroare.cod,r.eroare.param||{}));else if(r.configuratie){setConfig(r.configuratie);setOriginal(JSON.stringify(r.configuratie));setOptiuni(r.optiuni||[]);setMesaj(t('ui.salvat'));}}
        catch {setEroare(t('date.eroare'));}finally{blocat.current=false;setOcupat(false);onOcupat(false);}
    }
    return <section className="reguli-pontaj">
        {eroare&&<p role="alert" className="mesaj eroare">{eroare}</p>}{mesaj&&<p role="status" className="mesaj succes">{mesaj}</p>}
        <p className="nota-calendar">{t('pontaj.reguli_indiciu')}</p>
        <div className="unelte-calendar"><button className="buton secundar" disabled={ocupat} onClick={()=>{if(!modificat||window.confirm(t('ui.renunta')))incarca();}}>{t('pontaj.reincarca_reguli')}</button>
        <button className="buton primar" disabled={ocupat||!modificat} onClick={salveaza}>{t('ui.salveaza')}</button></div>
        {config&&<><label className="camp-atribuire">{t('pontaj.atribuire')}<select name="atribuire" disabled={ocupat} value={config.atribuire} onChange={e=>setConfig(new calcul.ConfiguratiePontaj({...config,atribuire:e.target.value}))}>
            <option value="ZI_INCEPUT">{t('pontaj.atribuire_ZI_INCEPUT')}</option><option value="ZILE_CALENDARISTICE">{t('pontaj.atribuire_ZILE_CALENDARISTICE')}</option></select></label>
        <p className="nota-calendar">{t('pontaj.calendar_indiciu')}</p>
        <div className="panou tabel-derulant"><table><thead><tr>{['coloana_regula','sursa_regula','selectie_regula','factor_regula','aprobare_regula'].map(c=><th key={c}>{t(`pontaj.${c}`)}</th>)}</tr></thead><tbody>
            {config.reguli.map((r,i)=>{const selectie=r.sursa==='ORE_SCHIMBURI'||r.sursa==='ZILE_ABSENTA';return <tr key={r.coloana}><th>{t(`pontaj.total_${r.coloana}`)}</th><td><select name={`sursa_${r.coloana}`} aria-label={t('pontaj.sursa_pentru',{coloana:t(`pontaj.total_${r.coloana}`)})} disabled={ocupat} value={r.sursa} onChange={e=>schimba(i,{sursa:e.target.value,coduri:[],factor_sutimi:100,numai_aprobate:false})}>
                {surse.filter(s=>s==='NEDEFINITA'||(!zile.includes(r.coloana)||s==='ZILE_ABSENTA')&&(!ore.includes(r.coloana)||s!=='ZILE_ABSENTA')).map(s=><option key={s} value={s}>{t(`pontaj.sursa_${s}`)}</option>)}</select></td>
                <td>{selectie?<select name={`coduri_${r.coloana}`} aria-label={t('pontaj.selectie_pentru',{coloana:t(`pontaj.total_${r.coloana}`)})} multiple required disabled={ocupat} value={r.coduri} onChange={e=>schimba(i,{coduri:Array.from(e.target.selectedOptions,o=>o.value)})}>
                    {optiuni.filter(o=>o.tip===(r.sursa==='ZILE_ABSENTA'?'ABSENTA':'SCHIMB')).map(o=><option key={o.cod} value={o.cod}>{o.denumire[i18n.language]||o.denumire.ro||o.cod} ({o.cod})</option>)}</select>:t('pontaj.fara_selectie')}</td>
                <td>{r.sursa!=='NEDEFINITA'&&r.sursa!=='ZILE_ABSENTA'?<input name={`factor_${r.coloana}`} aria-label={t('pontaj.factor_pentru',{coloana:t(`pontaj.total_${r.coloana}`)})} type="number" min="0.01" max="10" step="0.01" disabled={ocupat} value={r.factor_sutimi/100} onChange={e=>schimba(i,{factor_sutimi:Math.round(Number(e.target.value)*100)})}/>:t('pontaj.fara_selectie')}</td>
                <td>{r.sursa==='ZILE_ABSENTA'&&<input type="checkbox" aria-label={t('pontaj.aprobare_pentru',{coloana:t(`pontaj.total_${r.coloana}`)})} disabled={ocupat} checked={r.numai_aprobate} onChange={e=>schimba(i,{numai_aprobate:e.target.checked})}/>}</td></tr>;})}
        </tbody></table></div><p>{t('pontaj.revizie_reguli',{revizie:config.revizie})}</p></>}
    </section>;
}
