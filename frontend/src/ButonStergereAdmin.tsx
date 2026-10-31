import {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {StergeConfigurare} from '../wailsjs/go/main/App';
import {datebaza} from '../wailsjs/go/models';
import {useCont} from './Conturi';

export function ButonStergereAdmin({tip,id,revizie,activ,nume,disabled=false,onSters,onOcupat}:{tip:string;id:string;revizie:string;activ:boolean;nume:string;disabled?:boolean;onSters:(date:datebaza.Date)=>void|Promise<void>;onOcupat?:(ocupat:boolean)=>void}){
 const {t}=useTranslation();const cont=useCont();const [ocupat,setOcupat]=useState(false);const [eroare,setEroare]=useState('');
 if(cont?.cont.rol!=='ADMINISTRATOR')return null;
 async function sterge(){
  if(ocupat||disabled||activ||!revizie)return;
  if(!window.confirm(t('stergere.confirma',{nume})))return;
  setOcupat(true);onOcupat?.(true);setEroare('');
  try{
   const r=await StergeConfigurare(tip,id,revizie);
   if(r.eroare){const p={...r.eroare.param};if(p.sursa)p.sursa=t('stergere.sursa.'+p.sursa,{defaultValue:t('stergere.sursa.alte_date')});setEroare(t(r.eroare.cod,p));return;}
   if(r.date)await onSters(r.date);else setEroare(t('date.eroare'));
  }catch{setEroare(t('date.eroare'));}finally{setOcupat(false);onOcupat?.(false);}
 }
 return <span className="actiune-stergere-admin"><button type="button" data-stergere-admin={tip} className="legatura discreta" disabled={disabled||ocupat||activ||!revizie} title={t(activ?'stergere.dezactiveaza_intai':'stergere.indiciu')} onClick={()=>void sterge()}>{t('stergere.buton')}</button>{eroare&&<small role="alert" className="mesaj eroare">{eroare}</small>}</span>;
}
