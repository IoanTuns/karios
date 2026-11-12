import {ReguliConditionate} from './ReguliConditionate';
import {useEffect, useRef, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {semnaleazaModificari} from './modificariNesalvate';
import {CitesteConfiguratiePontaj, SalveazaConfiguratiePontaj} from '../wailsjs/go/main/App';
import {calcul, datebaza} from '../wailsjs/go/models';

// D092: surse noi pentru OL, CVE, TZLN si TOCN. Unitatea sursei trebuie sa fie a coloanei.
const surseOre = ['TOTAL_ORE','ORE_ZNL','ORE_NORMA','ORE_SCHIMBURI','ORE_SCHIMBURI_ZNL','ORE_LUCRATE_CODURI','ORE_NORMA_LUNA'];
const surseZile = ['ZILE_ABSENTA','ZILE_ABSENTA_ALTELE','ZILE_LUCRATOARE_NETE'];
const surse = ['NEDEFINITA',...surseOre,...surseZile];
const zile = ['CO','CEV','CFS','CM','deleg','CFP','CVE','DLG','FMP','TZLN'];
const ore = ['ore_lucrate','T','ZNL','ore_norma','ore_garda','OL','TOCN'];
const sursaPermisa = (coloana:string,s:string) => s==='NEDEFINITA' || (zile.includes(coloana) ? surseZile.includes(s) : ore.includes(coloana) ? surseOre.includes(s) : true);
// Sursele cu lista de coduri: obligatorie la cele dintai, optionala la OL (incluse) si CVE (excluse).
const cuCoduriObligatorii = (s:string) => s==='ORE_SCHIMBURI'||s==='ORE_SCHIMBURI_ZNL'||s==='ZILE_ABSENTA';
const cuCoduri = (s:string) => cuCoduriObligatorii(s)||s==='ORE_LUCRATE_CODURI'||s==='ZILE_ABSENTA_ALTELE';
const tipCoduri = (s:string) => s==='ZILE_ABSENTA'||s==='ZILE_ABSENTA_ALTELE'||s==='ORE_LUCRATE_CODURI' ? 'ABSENTA' : 'SCHIMB';
const cuFactor = (s:string) => surseOre.includes(s);
const cuAprobare = (s:string) => s==='ZILE_ABSENTA'||s==='ZILE_ABSENTA_ALTELE';

export function ReguliPontaj({onModificari,onOcupat,date}: {date:datebaza.Date;onModificari:(v:boolean)=>void;onOcupat:(v:boolean)=>void}) {
    const {t,i18n} = useTranslation();
    const [config,setConfig] = useState<calcul.ConfiguratiePontaj | null>(null);
    const [original,setOriginal] = useState('');
    const [optiuni,setOptiuni] = useState<datebaza.OptiuneRegula[]>([]);
    const [eroare,setEroare] = useState('');
    // F05-T20 (I049): eroarea de la server, legata de coloana si campul greseli.
    const [eroareCamp,setEroareCamp] = useState<{coloana:string;camp:string;cod:string}|null>(null);
    const [mesaj,setMesaj] = useState('');
    const [ocupat,setOcupat] = useState(false);
    const blocat = useRef(false);
    const modificat = Boolean(config && JSON.stringify(config)!==original);
    useEffect(()=>{onModificari(modificat);return ()=>onModificari(false);},[modificat,onModificari]);
    useEffect(()=>semnaleazaModificari('reguli', modificat),[modificat]);
    async function incarca() {
        if(blocat.current)return;blocat.current=true;setOcupat(true);onOcupat(true);setEroare('');
        try {const r=await CitesteConfiguratiePontaj();if(r.eroare)setEroare(t(r.eroare.cod,r.eroare.param||{}));else if(r.configuratie){setConfig(r.configuratie);setOriginal(JSON.stringify(r.configuratie));setOptiuni(r.optiuni||[]);}}
        catch {setEroare(t('date.eroare'));}finally{blocat.current=false;setOcupat(false);onOcupat(false);}
    }
    useEffect(()=>{incarca();},[]);
    function schimba(i:number,v:Partial<calcul.RegulaColoana>){setConfig(c=>c ? new calcul.ConfiguratiePontaj({...c,reguli:c.reguli.map((r,j)=>j===i?{...r,...v}:r)}):c);setMesaj('');setEroareCamp(null);}
    async function salveaza(){
        if(!config||blocat.current)return;blocat.current=true;setOcupat(true);onOcupat(true);setEroare('');setMesaj('');
        const invalide=config.reguli.filter(r=>eroareRand(r));
        if(invalide.length){setEroare(t('pontaj.reguli_corectati',{numar:String(invalide.length)}));blocat.current=false;setOcupat(false);onOcupat(false);
            const r=invalide[0],e=eroareRand(r)!;window.setTimeout(()=>(document.querySelector(`[name="${e.camp}_${r.coloana}"]`) as HTMLElement|null)?.focus(),0);return;}
        try {const r=await SalveazaConfiguratiePontaj(config);if(r.eroare){const p=r.eroare.param||{};setEroare(t(r.eroare.cod,{...p,coloana:p.coloana?t(`pontaj.total_${p.coloana}`):''}));if(p.coloana&&p.camp)setEroareCamp({coloana:p.coloana,camp:p.camp,cod:r.eroare.cod});}else if(r.configuratie){setConfig(r.configuratie);setOriginal(JSON.stringify(r.configuratie));setOptiuni(r.optiuni||[]);setMesaj(t('ui.salvat'));}}
        catch {setEroare(t('date.eroare'));}finally{blocat.current=false;setOcupat(false);onOcupat(false);}
    }
    // eroareRand: verificarea de dinaintea salvarii, aceleasi limite ca in calcul.
    const eroareRand=(r:calcul.RegulaColoana):{camp:string;cod:string}|null=>{
        const selectie=cuCoduriObligatorii(r.sursa);
        if(selectie&&!(r.coduri||[]).length)return {camp:'coduri',cod:'pontaj.regula_fara_coduri'};
        if(cuFactor(r.sursa)&&!(Number.isFinite(r.factor_sutimi)&&r.factor_sutimi>=1&&r.factor_sutimi<=1000))return {camp:'factor',cod:'pontaj.regula_factor'};
        if(eroareCamp&&eroareCamp.coloana===r.coloana)return {camp:eroareCamp.camp,cod:eroareCamp.cod};
        return null;
    };
    const mesajCamp=(r:calcul.RegulaColoana,camp:string)=>{const e=eroareRand(r);return e&&e.camp===camp?<p id={`eroare_${camp}_${r.coloana}`} className="eroare-camp" role="alert">{t(e.cod,{coloana:t(`pontaj.total_${r.coloana}`)})}</p>:null;};
    const invalid=(r:calcul.RegulaColoana,camp:string)=>{const e=eroareRand(r);return e&&e.camp===camp?{'aria-invalid':true,'aria-describedby':`eroare_${camp}_${r.coloana}`}:{};};
    return <section className="reguli-pontaj">
        {eroare&&<p role="alert" className="mesaj eroare">{eroare}</p>}{mesaj&&<p role="status" className="mesaj succes">{mesaj}</p>}
        <p className="nota-calendar">{t('pontaj.reguli_indiciu')}</p>
        <div className="unelte-calendar"><button className="buton secundar" disabled={ocupat} onClick={()=>{if(!modificat||window.confirm(t('ui.renunta')))incarca();}}>{t('pontaj.reincarca_reguli')}</button>
        <button className="buton primar" disabled={ocupat||!modificat} onClick={salveaza}>{t('ui.salveaza')}</button></div>
        {config&&<><label className="camp-atribuire">{t('pontaj.atribuire')}<select name="atribuire" disabled={ocupat} value={config.atribuire} onChange={e=>setConfig(new calcul.ConfiguratiePontaj({...config,atribuire:e.target.value}))}>
            <option value="ZI_INCEPUT">{t('pontaj.atribuire_ZI_INCEPUT')}</option><option value="ZILE_CALENDARISTICE">{t('pontaj.atribuire_ZILE_CALENDARISTICE')}</option></select></label>
        {/* F04-T08 (D102): turele din noptile cu schimbarea orei, durata din program sau orele efective. */}
        <label className="camp-atribuire">{t('pontaj.schimbare_ora')}<select name="schimbare_ora" disabled={ocupat} value={config.schimbare_ora||'NOMINALA'} onChange={e=>setConfig(new calcul.ConfiguratiePontaj({...config,schimbare_ora:e.target.value}))}>
            <option value="NOMINALA">{t('pontaj.schimbare_ora_NOMINALA')}</option><option value="EFECTIVA">{t('pontaj.schimbare_ora_EFECTIVA')}</option></select></label>
        <p className="nota-calendar">{t('pontaj.schimbare_ora_indiciu')}</p>
        <p className="nota-calendar">{t('pontaj.calendar_indiciu')}</p><ReguliConditionate date={date} lista={config.conditionate??[]} ocupat={ocupat} onChange={r=>setConfig(new calcul.ConfiguratiePontaj({...config,conditionate:r}))}/>
        <div className="panou tabel-derulant"><table className="tabel-reguli-pontaj"><thead><tr>{['coloana_regula','sursa_regula','selectie_regula','factor_regula','aprobare_regula'].map(c=><th key={c}>{t(`pontaj.${c}`)}</th>)}</tr></thead><tbody>
            {config.reguli.map((r,i)=>{const selectie=cuCoduri(r.sursa);return <tr key={r.coloana}><th>{t(`pontaj.total_${r.coloana}`)}</th><td><select name={`sursa_${r.coloana}`} aria-label={t('pontaj.sursa_pentru',{coloana:t(`pontaj.total_${r.coloana}`)})} disabled={ocupat} value={r.sursa} {...invalid(r,'sursa')} onChange={e=>schimba(i,{sursa:e.target.value,coduri:e.target.value==='ZILE_ABSENTA_ALTELE'?['CO','CM','CFP','DLG','FMP'].filter(c=>optiuni.some(o=>o.tip==='ABSENTA'&&o.cod===c)):[],factor_sutimi:100,numai_aprobate:false})}>
                {surse.filter(s=>sursaPermisa(r.coloana,s)).map(s=><option key={s} value={s}>{t(`pontaj.sursa_${s}`)}</option>)}</select>{mesajCamp(r,'sursa')}</td>
                <td>{selectie?<>{/* F05-T23 (I052): o singura selectie, cu butoane folosibile de la tastatura; numele complet e in eticheta si tooltip. */}
                    <div className="optiuni-rapide-regula" role="group" aria-label={t('pontaj.selectie_pentru',{coloana:t(`pontaj.total_${r.coloana}`)})} {...invalid(r,'coduri')}>{optiuni.filter(o=>o.tip===tipCoduri(r.sursa)).map((o,index)=>{const bifat=r.coduri.includes(o.cod);const nume=o.denumire[i18n.language]||o.denumire.ro||o.cod;return <button key={o.cod} name={index===0?`coduri_${r.coloana}`:undefined} type="button" className={`cip-schimb ${bifat?'activ':''}`} aria-pressed={bifat} aria-label={`${o.cod} · ${nume}`} title={nume} disabled={ocupat} onClick={()=>schimba(i,{coduri:bifat?r.coduri.filter(c=>c!==o.cod):[...r.coduri,o.cod]})}>{o.cod}</button>;})}</div>{mesajCamp(r,'coduri')}</>:t('pontaj.fara_selectie')}</td>
                <td>{cuFactor(r.sursa)?<input name={`factor_${r.coloana}`} aria-label={t('pontaj.factor_pentru',{coloana:t(`pontaj.total_${r.coloana}`)})} type="number" min="0.01" max="10" step="0.01" disabled={ocupat} value={r.factor_sutimi/100} {...invalid(r,'factor')} onChange={e=>schimba(i,{factor_sutimi:Math.round(Number(e.target.value)*100)})}/>:t('pontaj.fara_selectie')}{mesajCamp(r,'factor')}</td>
                <td>{cuAprobare(r.sursa)&&<input type="checkbox" aria-label={t('pontaj.aprobare_pentru',{coloana:t(`pontaj.total_${r.coloana}`)})} disabled={ocupat} checked={r.numai_aprobate} onChange={e=>schimba(i,{numai_aprobate:e.target.checked})}/>}</td></tr>;})}
        </tbody></table></div><p>{t('pontaj.revizie_reguli',{revizie:config.revizie})}</p></>}
    </section>;
}
