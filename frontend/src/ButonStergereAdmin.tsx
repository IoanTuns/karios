import {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {StergeConfigurare} from '../wailsjs/go/main/App';
import {datebaza} from '../wailsjs/go/models';
import {useCont} from './Conturi';

export function ButonStergereAdmin({tip,id,revizie,activ,nume,disabled=false,onSters,onOcupat}:{tip:string;id:string;revizie:string;activ:boolean;nume:string;disabled?:boolean;onSters:(date:datebaza.Date)=>void|Promise<void>;onOcupat?:(ocupat:boolean)=>void}){
 const {t}=useTranslation();const cont=useCont();const [ocupat,setOcupat]=useState(false);const [eroare,setEroare]=useState('');const [sursa,setSursa]=useState('');
 if(cont?.cont.rol!=='ADMINISTRATOR'||activ)return null;
 async function sterge(){
  if(ocupat||disabled||activ||!revizie)return;
  if(!window.confirm(t('stergere.confirma',{nume})))return;
  setOcupat(true);onOcupat?.(true);setEroare('');setSursa('');
  try{
   const r=await StergeConfigurare(tip,id,revizie);
   if(r.eroare){setSursa(r.eroare.param?.sursa||'');const p={...r.eroare.param};if(p.sursa)p.sursa=t('stergere.sursa.'+p.sursa,{defaultValue:t('stergere.sursa.alte_date')});setEroare(t(r.eroare.cod,p));return;}
   if(r.date)await onSters(r.date);else setEroare(t('date.eroare'));
  }catch{setEroare(t('date.eroare'));}finally{setOcupat(false);onOcupat?.(false);}
 }
 const ecranSursa:Record<string,string>={angajati:'angajati',angajati_functii:'angajati',angajati_contracte:'angajati',angajati_prefixe:'angajati',functii:'setari',prefixe:'setari',tipuri_contract:'setari',categorii_program:'categorii_program',categorii_program_schimburi:'categorii_program',clasificari_schimburi:'schimburi',tipuri_schimb:'schimburi',pontaj:'pontaj',pontaj_ciorne:'pontaj',planificare:'planificare',absente_zile:'absente',setari:'setari',sabloane_ai:'asistent_ai',operatii_ai:'asistent_ai'};
 return <details className="actiune-stergere-admin"><summary>{t('ui28.alte_actiuni')}</summary><button type="button" data-stergere-admin={tip} className="legatura discreta" disabled={disabled||ocupat||activ||!revizie} title={t(activ?'stergere.dezactiveaza_intai':'stergere.indiciu')} onClick={()=>void sterge()}>{t('stergere.buton')}</button>{eroare&&<small role="alert" className="mesaj eroare">{eroare}</small>}{ecranSursa[sursa]&&<button type="button" className="legatura" onClick={()=>window.dispatchEvent(new CustomEvent('pontix-configurare',{detail:ecranSursa[sursa]}))}>{t('ui28.deschide_sursa')}</button>}</details>;
}
