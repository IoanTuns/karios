import {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {datebaza} from '../wailsjs/go/models';
import {Dialog} from './Formulare';
export function FormularCategorieAbsenta({categorie,ocupat,eroare,inchide,salveaza}:{categorie:datebaza.CategorieAbsenta;ocupat:boolean;eroare:string;inchide:()=>void;salveaza:(c:datebaza.CategorieAbsenta)=>Promise<void>}){
 const {t}=useTranslation();const initial={cod:categorie.cod,ro:categorie.denumire.ro||'',en:categorie.denumire.en||'',natura:categorie.natura,efect:categorie.efect_ore,activ:categorie.activ};const [campuri,setCampuri]=useState(initial);const [validare,setValidare]=useState('');
 function renunta(){if(JSON.stringify(initial)===JSON.stringify(campuri)||window.confirm(t('ui.renunta')))inchide();}
 return <Dialog titlu={t(categorie.cod?'absente.editeaza_categorie':'absente.adauga_categorie')} ocupat={ocupat} inchide={renunta}><form onSubmit={e=>{e.preventDefault();if(!campuri.ro.trim()&&!campuri.en.trim()){setValidare(t('sarbatori.denumire_obligatorie'));return;}void salveaza(new datebaza.CategorieAbsenta({cod:campuri.cod,denumire:{ro:campuri.ro,en:campuri.en},natura:campuri.natura,efect_ore:campuri.efect,activ:campuri.activ,revizie:categorie.revizie}));}}><div className="corp-dialog">{(eroare||validare)&&<p role="alert" className="mesaj eroare">{eroare||validare}</p>}<fieldset className="campuri" disabled={ocupat}>
 <label>{t('schimburi.cod')} *<input name="cod_categorie_absenta" required maxLength={48} readOnly={!!categorie.cod} value={campuri.cod} onChange={e=>setCampuri(c=>({...c,cod:e.target.value}))}/></label>
 <label>{t('sarbatori.denumire_ro')}<input name="categorie_absenta_ro" value={campuri.ro} onChange={e=>setCampuri(c=>({...c,ro:e.target.value}))}/></label>
 <label>{t('sarbatori.denumire_en')}<input name="categorie_absenta_en" value={campuri.en} onChange={e=>setCampuri(c=>({...c,en:e.target.value}))}/></label>
 <label>{t('absente.natura')}<select name="natura_absenta" value={campuri.natura} onChange={e=>{const natura=e.target.value;setCampuri(c=>({...c,natura,efect:natura!=='SPECIAL'&&['NORMA','PROGRAM'].includes(c.efect)?'FARA_ORE':c.efect}));}}>{['ODIHNA','MEDICAL','SPECIAL'].map(v=><option key={v} value={v}>{t(v==='SPECIAL'?'absente.natura_SPECIAL':'absente.'+v)}</option>)}</select></label>
 <label>{t('absente.efect_ore')}<select name="efect_ore" value={campuri.efect} onChange={e=>setCampuri(c=>({...c,efect:e.target.value}))}>{[...(categorie.efect_ore==='ANTERIOR'?['ANTERIOR']:[]),'FARA_ORE',...(campuri.natura==='SPECIAL'?['NORMA','PROGRAM']:[])].map(v=><option key={v} value={v}>{t('absente.efect_'+v)}</option>)}</select></label>
 <label className="bifa"><input type="checkbox" name="categorie_absenta_activa" checked={campuri.activ} onChange={e=>setCampuri(c=>({...c,activ:e.target.checked}))}/>{t('ui.activ')}</label>
 </fieldset><p className="nota-calendar">{t('absente.categorii_indiciu')}</p></div><footer className="subsol-dialog"><button type="button" className="buton secundar" disabled={ocupat} onClick={renunta}>{t('ui.anuleaza')}</button><button className="buton primar" disabled={ocupat}>{t('ui.salveaza')}</button></footer></form></Dialog>;
}
