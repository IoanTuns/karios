import {useEffect, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {ListaSchimburi} from '../wailsjs/go/main/App';
import {datebaza} from '../wailsjs/go/models';
import {Dialog, durata} from './Formulare';
import {oraText} from './Schimburi';

export function CategoriiProgram({date,ocupat,eroare,salveaza,reincarca,curataEroare}: {
    date:datebaza.Date;ocupat:boolean;eroare:string;salveaza:(c:datebaza.CategorieProgram)=>Promise<boolean>;reincarca:()=>Promise<boolean>;curataEroare:()=>void;
}) {
    const {t,i18n}=useTranslation();
    const [schimburi,setSchimburi]=useState<datebaza.Schimb[]>([]);
    const [eroareCitire,setEroareCitire]=useState('');
    const [incarcare,setIncarcare]=useState(true);
    const [editare,setEditare]=useState<datebaza.CategorieProgram|null>(null);
    useEffect(()=>{
        let activ=true;
        ListaSchimburi().then(r=>{if(!activ)return;if(r.eroare)setEroareCitire(t(r.eroare.cod,r.eroare.param??{}));else setSchimburi(r.schimburi??[]);})
            .catch(()=>{if(activ)setEroareCitire(t('date.eroare'));}).finally(()=>{if(activ)setIncarcare(false);});
        return()=>{activ=false;};
    },[]);
    const nume=(c:datebaza.CategorieProgram)=>c.denumire[i18n.language]||c.denumire.ro||Object.values(c.denumire).find(Boolean)||'';
    const program=(id:string)=>schimburi.find(s=>s.id===id)?.cod||t('program.schimb_lipsa');
    const deschide=(c:datebaza.CategorieProgram)=>{curataEroare();setEditare(c);};
    async function reincarcaTot(){
        setIncarcare(true);setEroareCitire('');
        try{const [,r]=await Promise.all([reincarca(),ListaSchimburi()]);if(r.eroare)setEroareCitire(t(r.eroare.cod,r.eroare.param??{}));else setSchimburi(r.schimburi??[]);}
        catch{setEroareCitire(t('date.eroare'));}finally{setIncarcare(false);}
    }
    return <section className="categorii-program">
        <p className="nota-calendar">{t('program.indiciu')}</p>
        {eroareCitire&&<p role="alert" className="mesaj eroare">{eroareCitire}</p>}
        <div className="unelte-calendar"><button type="button" className="buton primar" disabled={ocupat||incarcare||Boolean(eroareCitire)} onClick={()=>deschide(new datebaza.CategorieProgram({id:'',denumire:{ro:'',en:''},schimburi_id:[],schimb_implicit_id:'',activ:true,revizie:''}))}>{t('program.adauga')}</button><button className="buton secundar" disabled={ocupat||incarcare} onClick={()=>void reincarcaTot()}>{t('program.reincarca')}</button></div>
        <div className="panou tabel-derulant"><table className="tabel-categorii-program"><thead><tr>{['denumire','programe','implicit','stare','actiuni'].map(c=><th key={c}>{t(c==='stare'||c==='actiuni'?`ui.${c}`:`program.${c}`)}</th>)}</tr></thead><tbody>
            {(date.categorii_program??[]).map(c=><tr key={c.id}><th scope="row">{nume(c)}</th><td>{c.schimburi_id.map(program).join(', ')}</td><td>{c.schimb_implicit_id?program(c.schimb_implicit_id):t('program.fara_implicit')}</td><td>{t(c.activ?'ui.activ':'ui.inactiv')}</td><td><button type="button" className="legatura" disabled={ocupat||incarcare||Boolean(eroareCitire)} onClick={()=>deschide(c)}>{t('ui.editeaza')}</button></td></tr>)}
        </tbody></table></div>
        {editare&&<FormularCategorie key={editare.id||'nou'} categorie={editare} schimburi={schimburi} ocupat={ocupat} eroare={eroare} inchide={()=>setEditare(null)} salveaza={async c=>{if(await salveaza(c))setEditare(null);}}/>}
    </section>;
}

function FormularCategorie({categorie,schimburi,ocupat,eroare,inchide,salveaza}:{categorie:datebaza.CategorieProgram;schimburi:datebaza.Schimb[];ocupat:boolean;eroare:string;inchide:()=>void;salveaza:(c:datebaza.CategorieProgram)=>Promise<void>}) {
    const {t,i18n}=useTranslation();
    const [c,setC]=useState(()=>new datebaza.CategorieProgram({...categorie,denumire:{...categorie.denumire},schimburi_id:[...categorie.schimburi_id]}));
    const [validare,setValidare]=useState('');
    const renunta=()=>{if(!ocupat&&(JSON.stringify(c)===JSON.stringify(categorie)||window.confirm(t('ui.renunta'))))inchide();};
    const textSchimb=(s:datebaza.Schimb)=>`${s.cod} · ${s.ora_start>=0?oraText(s.ora_start):t('pontaj.fara_ora')} · ${durata(s.durata_minute,t)}${!s.activ?` · ${t('ui.inactiv')}`:''}`;
    return <Dialog titlu={t(c.id?'program.editeaza':'program.adauga')} ocupat={ocupat} inchide={renunta}><form onSubmit={e=>{
        e.preventDefault();setValidare('');
        if(!Object.values(c.denumire).some(v=>v.trim())||!c.schimburi_id.length){setValidare(t('program.categorie_invalida'));return;}
        if(c.schimb_implicit_id&&!c.schimburi_id.includes(c.schimb_implicit_id)){setValidare(t('program.implicit_invalid'));return;}
        void salveaza(c);
    }}><div className="corp-dialog">
        {(validare||eroare)&&<p role="alert" className="mesaj eroare">{validare||eroare}</p>}
        <fieldset disabled={ocupat} className="campuri">
            {['ro','en'].map(l=><label key={l}>{t(`program.denumire_${l}`)}<input name={`categorie_denumire_${l}`} autoFocus={l==='ro'} value={c.denumire[l]??''} onChange={e=>setC(new datebaza.CategorieProgram({...c,denumire:{...c.denumire,[l]:e.target.value}}))}/></label>)}
            <fieldset className="latime-completa selectie-programe"><legend>{t('program.programe')}</legend>
                {schimburi.map(s=><label key={s.id} className="bifa"><input type="checkbox" name={`categorie_schimb_${s.cod}`} checked={c.schimburi_id.includes(s.id)} disabled={!s.activ&&!c.schimburi_id.includes(s.id)} onChange={e=>setC(new datebaza.CategorieProgram({...c,schimburi_id:e.target.checked?[...c.schimburi_id,s.id]:c.schimburi_id.filter(id=>id!==s.id),schimb_implicit_id:!e.target.checked&&c.schimb_implicit_id===s.id?'':c.schimb_implicit_id}))}/><span>{textSchimb(s)}<small>{s.denumire[i18n.language]||s.denumire.ro}</small></span></label>)}
            </fieldset>
            <label className="latime-completa">{t('program.implicit')}<select name="categorie_implicit" value={c.schimb_implicit_id} onChange={e=>setC(new datebaza.CategorieProgram({...c,schimb_implicit_id:e.target.value}))}>
                <option value="">{t('program.fara_implicit')}</option>{schimburi.filter(s=>c.schimburi_id.includes(s.id)).map(s=><option key={s.id} value={s.id}>{textSchimb(s)}</option>)}</select><small>{t('program.implicit_indiciu')}</small></label>
            <label className="bifa latime-completa"><input name="categorie_activa" type="checkbox" checked={c.activ} onChange={e=>setC(new datebaza.CategorieProgram({...c,activ:e.target.checked}))}/>{t('program.activa')}</label>
        </fieldset>
    </div><footer className="subsol-dialog"><button type="button" className="buton secundar" disabled={ocupat} onClick={renunta}>{t('ui.anuleaza')}</button><button className="buton primar" disabled={ocupat}>{t('ui.salveaza')}</button></footer></form></Dialog>;
}
