import {useTranslation} from 'react-i18next';

export type StareLista = 'ACTIVE'|'INACTIVE'|'TOATE';
export function corespundeStare(activ:boolean,stare:StareLista):boolean {
 return stare==='TOATE'||(stare==='ACTIVE'?activ:!activ);
}
export function FiltruStare({valoare,schimba,nume,eticheta,disabled=false}:{valoare:StareLista;schimba:(stare:StareLista)=>void;nume:string;eticheta?:string;disabled?:boolean}){
 const {t}=useTranslation();
 return <label className="filtru-stare">{t(eticheta||'ui.filtru_stare')}<select name={'filtru_stare_'+nume} disabled={disabled} value={valoare} onChange={e=>schimba(e.target.value as StareLista)}>{(['ACTIVE','INACTIVE','TOATE'] as const).map(s=><option key={s} value={s}>{t('ui.filtru_stare_'+s)}</option>)}</select></label>;
}
