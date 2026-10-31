import {FiltruStare,corespundeStare,type StareLista} from './FiltruStare';
import {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {DateInitiale,SalveazaElementNomenclator} from '../wailsjs/go/main/App';
import {datebaza} from '../wailsjs/go/models';

export function EditorNomenclator({tip,lista,onDate,date}:{tip:'categorii'|'functii'|'grupuri'|'contracte'|'tipuri_program'|'regimuri_program';lista:datebaza.ElementNomenclator[];onDate:(date:datebaza.Date)=>void;date?:datebaza.Date}) {
 const {t,i18n}=useTranslation();
 const [stareLista,setStareLista]=useState<StareLista>('ACTIVE');
 const [editare,setEditare]=useState<datebaza.ElementNomenclator|null>(null);
 const [ocupat,setOcupat]=useState(false),[eroare,setEroare]=useState('');
 const nume=(e:datebaza.ElementNomenclator)=>e.denumire[i18n.language]||e.denumire.ro||e.denumire.en||'';
 return <section className="panou card-setari nomenclator" data-nomenclator={tip}>
 <h2>{t('nomenclator.'+tip)}</h2>
 <p>{t('nomenclator.'+tip+'_indiciu')}</p>
 {eroare&&<p role="alert" className="mesaj eroare">{eroare}</p>}
 <button type="button" className="buton secundar" disabled={ocupat||!!editare} onClick={()=>{setEroare('');setEditare(new datebaza.ElementNomenclator({id:'',denumire:{ro:'',en:''},activ:true,revizie:'',ordine:0,grup_id:'',conducere:false,normare:'FARA',norma_minute:undefined,categorie_program_id:'',garda:false,comportament:tip==='regimuri_program'?'LUCRATOARE':'NORMAL',utilizari:0}));}}>{t('nomenclator.adauga')}</button>
 <button type="button" className="buton secundar" disabled={ocupat} onClick={async()=>{if(editare&&!window.confirm(t('ui.renunta')))return;setOcupat(true);setEroare('');try{const r=await DateInitiale();if(r.eroare)setEroare(t(r.eroare.cod,r.eroare.param??{}));else if(r.date){onDate(r.date);setEditare(null);}}catch{setEroare(t('date.eroare'));}finally{setOcupat(false);}}}>{t('program.reincarca')}</button>
 <FiltruStare nume={tip} valoare={stareLista} schimba={setStareLista} disabled={ocupat}/>
 <ul>{lista.filter(e=>corespundeStare(e.activ,stareLista)).map(e=><li key={e.id}>{nume(e)} {!e.activ&&<small>{t('ui.inactiv')}</small>} <button type="button" className="legatura" disabled={ocupat||!!editare} onClick={()=>{setEroare('');setEditare(new datebaza.ElementNomenclator({...e,denumire:{...e.denumire}}));}}>{t('ui.editeaza')}</button></li>)}</ul>
 {editare&&<form onSubmit={async e=>{e.preventDefault();if(ocupat)return;setOcupat(true);setEroare('');try{const r=await SalveazaElementNomenclator(tip,editare);if(r.eroare)setEroare(t(r.eroare.cod,r.eroare.param??{}));else if(r.date){onDate(r.date);setEditare(null);}}catch{setEroare(t('date.eroare'));}finally{setOcupat(false);}}}>
 <fieldset className="campuri" disabled={ocupat}>
 <label>{t('sarbatori.denumire_ro')}<input name="nomenclator_ro" autoFocus value={editare.denumire.ro||''} onChange={e=>setEditare({...editare,denumire:{...editare.denumire,ro:e.target.value}})}/></label>
 <label>{t('sarbatori.denumire_en')}<input name="nomenclator_en" value={editare.denumire.en||''} onChange={e=>setEditare({...editare,denumire:{...editare.denumire,en:e.target.value}})}/></label>
 {(tip==='tipuri_program'||tip==='regimuri_program')&&<label>{t('clasificari.comportament')}<select name="comportament" disabled={editare.utilizari>0} value={editare.comportament} onChange={e=>setEditare({...editare,comportament:e.target.value})}>{(tip==='tipuri_program'?['NORMAL','TURA','NOAPTE','GARDA']:['LUCRATOARE','NELUCRATOARE','INTRARE_ZI_LIBERA']).map(c=><option key={c} value={c}>{t((tip==='tipuri_program'?'schimburi.tip_':'schimburi.regim_')+c)}</option>)}</select>{editare.utilizari>0&&<small>{t('clasificari.utilizat')}</small>}</label>}
 {(tip==='grupuri')&&<label>{t('structura.ordine')}<input name="ordine" type="number" min="0" max="100000" value={editare.ordine??0} onChange={e=>setEditare({...editare,ordine:Number(e.target.value)})}/></label>}
 {tip==='functii'&&<><label>{t('structura.grup')}<select name="grup_functie" value={editare.grup_id||''} onChange={e=>setEditare({...editare,grup_id:e.target.value})}><option value="">{t('structura.fara_grup')}</option>{(date?.grupuri_profesionale??[]).filter(g=>g.activ||g.id===editare.grup_id).map(g=><option key={g.id} value={g.id}>{nume(g)}</option>)}</select></label><label className="bifa"><input name="conducere" type="checkbox" checked={editare.conducere??false} onChange={e=>setEditare({...editare,conducere:e.target.checked})}/>{t('structura.conducere')}</label></>}
 {tip==='contracte'&&<><label>{t('structura.normare')}<select name="normare" value={editare.normare||'FARA'} onChange={e=>setEditare({...editare,normare:e.target.value,norma_minute:undefined})}>{['GLOBALA','INDIVIDUALA','ZILNICA','LUNARA','FARA'].map(m=><option key={m} value={m}>{t('structura.normare_'+m)}</option>)}</select></label>{['ZILNICA','LUNARA'].includes(editare.normare)&&<label>{t('structura.minute')}<input name="norma_contract" required type="number" min="1" max={editare.normare==='ZILNICA'?1440:44640} value={editare.norma_minute??''} onChange={e=>setEditare({...editare,norma_minute:e.target.value===''?undefined:Number(e.target.value)})}/></label>}<label>{t('program.categorie')}<select value={editare.categorie_program_id||''} onChange={e=>setEditare({...editare,categorie_program_id:e.target.value})}><option value="">{t('program.fara_categorie')}</option>{(date?.categorii_program??[]).map(c=><option key={c.id} value={c.id}>{c.denumire[i18n.language]||c.denumire.ro}</option>)}</select></label><label className="bifa"><input type="checkbox" checked={editare.garda??false} onChange={e=>setEditare({...editare,garda:e.target.checked})}/>{t('structura.garda')}</label></>}
 <label className="bifa"><input type="checkbox" checked={editare.activ} onChange={e=>setEditare({...editare,activ:e.target.checked})}/>{t('ui.activ_indiciu')}</label>
 </fieldset>
 <button type="button" className="buton secundar" disabled={ocupat} onClick={()=>{if(window.confirm(t('ui.renunta')))setEditare(null);}}>{t('ui.anuleaza')}</button>
 <button className="buton primar" disabled={ocupat}>{t('ui.salveaza')}</button>
 </form>}
 </section>;
}
