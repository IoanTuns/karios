import {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {DateInitiale,SalveazaElementNomenclator} from '../wailsjs/go/main/App';
import {datebaza} from '../wailsjs/go/models';

export function EditorNomenclator({tip,lista,onDate}:{tip:'categorii'|'functii';lista:datebaza.ElementNomenclator[];onDate:(date:datebaza.Date)=>void}) {
 const {t,i18n}=useTranslation();
 const [editare,setEditare]=useState<datebaza.ElementNomenclator|null>(null);
 const [ocupat,setOcupat]=useState(false),[eroare,setEroare]=useState('');
 const nume=(e:datebaza.ElementNomenclator)=>e.denumire[i18n.language]||e.denumire.ro||e.denumire.en||'';
 return <section className="panou card-setari nomenclator" data-nomenclator={tip}>
 <h2>{t(tip==='functii'?'nomenclator.functii':'nomenclator.categorii')}</h2>
 <p>{t(tip==='functii'?'nomenclator.functii_indiciu':'nomenclator.categorii_indiciu')}</p>
 {eroare&&<p role="alert" className="mesaj eroare">{eroare}</p>}
 <button type="button" className="buton secundar" disabled={ocupat||!!editare} onClick={()=>{setEroare('');setEditare(new datebaza.ElementNomenclator({id:'',denumire:{ro:'',en:''},activ:true,revizie:''}));}}>{t('nomenclator.adauga')}</button>
 <button type="button" className="buton secundar" disabled={ocupat} onClick={async()=>{if(editare&&!window.confirm(t('ui.renunta')))return;setOcupat(true);setEroare('');try{const r=await DateInitiale();if(r.eroare)setEroare(t(r.eroare.cod,r.eroare.param??{}));else if(r.date){onDate(r.date);setEditare(null);}}catch{setEroare(t('date.eroare'));}finally{setOcupat(false);}}}>{t('program.reincarca')}</button>
 <ul>{lista.map(e=><li key={e.id}>{nume(e)} {!e.activ&&<small>{t('ui.inactiv')}</small>} <button type="button" className="legatura" disabled={ocupat||!!editare} onClick={()=>{setEroare('');setEditare(new datebaza.ElementNomenclator({...e,denumire:{...e.denumire}}));}}>{t('ui.editeaza')}</button></li>)}</ul>
 {editare&&<form onSubmit={async e=>{e.preventDefault();if(ocupat)return;setOcupat(true);setEroare('');try{const r=await SalveazaElementNomenclator(tip,editare);if(r.eroare)setEroare(t(r.eroare.cod,r.eroare.param??{}));else if(r.date){onDate(r.date);setEditare(null);}}catch{setEroare(t('date.eroare'));}finally{setOcupat(false);}}}>
 <fieldset className="campuri" disabled={ocupat}>
 <label>{t('sarbatori.denumire_ro')}<input name="nomenclator_ro" autoFocus value={editare.denumire.ro||''} onChange={e=>setEditare({...editare,denumire:{...editare.denumire,ro:e.target.value}})}/></label>
 <label>{t('sarbatori.denumire_en')}<input name="nomenclator_en" value={editare.denumire.en||''} onChange={e=>setEditare({...editare,denumire:{...editare.denumire,en:e.target.value}})}/></label>
 <label className="bifa"><input type="checkbox" checked={editare.activ} onChange={e=>setEditare({...editare,activ:e.target.checked})}/>{t('ui.activ_indiciu')}</label>
 </fieldset>
 <button type="button" className="buton secundar" disabled={ocupat} onClick={()=>{if(window.confirm(t('ui.renunta')))setEditare(null);}}>{t('ui.anuleaza')}</button>
 <button className="buton primar" disabled={ocupat}>{t('ui.salveaza')}</button>
 </form>}
 </section>;
}
