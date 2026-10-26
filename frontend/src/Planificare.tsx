import {useEffect, useMemo, useRef, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {CitesteCalendar, CitestePlanificare, ListaSchimburi, PrevizualizeazaPublicarePlan, PublicaPlanificare, SalveazaPlanificare} from '../wailsjs/go/main/App';
import {datebaza} from '../wailsjs/go/models';
import {Dialog, denumireSectie, durata} from './Formulare';
import {oraText} from './Schimburi';

export function programPlanificat(p:datebaza.Schimb) {
    return p.tip==='GARDA'||p.tip==='NOAPTE'||p.regim==='NELUCRATOARE'||p.regim==='INTRARE_ZI_LIBERA'||(p.ora_start>=720&&p.ora_start<1320);
}
function inSectiune(p:datebaza.Schimb,sectiune:string) {
    const garda=p.tip==='GARDA'||p.tip==='NOAPTE'||p.regim==='NELUCRATOARE'||p.regim==='INTRARE_ZI_LIBERA';
    return sectiune==='GARZI'?garda:!garda&&p.ora_start>=720&&p.ora_start<1320;
}
type PublicareCurenta={cerere:datebaza.CererePublicarePlan;rezultat:datebaza.PublicarePlan};

export function Planificare({date,initial,onOcupat,onDeschidePontaj}:{date:datebaza.Date;initial?:{an:number;luna:number};onOcupat:(v:boolean)=>void;onDeschidePontaj:(an:number,luna:number)=>void}) {
    const {t,i18n}=useTranslation();
    const [an,setAn]=useState(initial?.an??new Date().getFullYear());
    const [luna,setLuna]=useState(initial?.luna??new Date().getMonth()+1);
    const [anText,setAnText]=useState(String(an));
    const [sectiune,setSectiune]=useState('GARZI');
    const [sectie,setSectie]=useState('');const [departament,setDepartament]=useState('');const [cautare,setCautare]=useState('');
    const [inactive,setInactive]=useState(false);
    const [plan,setPlan]=useState<datebaza.PlanificareLuna|null>(null);
    const [schimburi,setSchimburi]=useState<datebaza.Schimb[]>([]);
    const [calendar,setCalendar]=useState<datebaza.CalendarAn|null>(null);
    const [selectate,setSelectate]=useState<Set<string>>(new Set());
    const [editare,setEditare]=useState<datebaza.PozitiePlanificare|null>(null);
    const [publicare,setPublicare]=useState<PublicareCurenta|null>(null);
    const [eroare,setEroare]=useState('');const [mesaj,setMesaj]=useState('');const [ocupat,setOcupat]=useState(false);
    const [incarcare,setIncarcare]=useState(true);const [reincarcare,setReincarcare]=useState(0);
    const blocat=useRef(false);
    const nume=(id:string)=>{const a=date.angajati.find(a=>a.id===id);return a?`${a.nume} ${a.prenume}`:id;};
    const textSchimb=(id:string)=>{const s=schimburi.find(s=>s.id===id);return s?`${s.cod} · ${s.denumire[i18n.language]||s.denumire.ro} · ${durata(s.durata_minute,t)}`:t('program.schimb_lipsa');};
    const sarbatori=useMemo(()=>new Map((calendar?.sarbatori??[]).map(z=>[z.data,z.denumire[i18n.language]||z.denumire.ro])),[calendar,i18n.language]);
    useEffect(()=>{
        let activ=true;setIncarcare(true);setPlan(null);setEroare('');setMesaj('');setSelectate(new Set());
        Promise.all([CitestePlanificare(an,luna),ListaSchimburi(),CitesteCalendar(an)]).then(([p,s,c])=>{
            if(!activ)return;const err=p.eroare||s.eroare||c.eroare;
            if(err){setEroare(t(err.cod,err.param??{}));return;}
            if(!p.planificare||!c.calendar){setEroare(t('date.eroare'));return;}
            setPlan(p.planificare);setSchimburi(s.schimburi??[]);setCalendar(c.calendar);
        }).catch(()=>{if(activ)setEroare(t('date.eroare'));}).finally(()=>{if(activ)setIncarcare(false);});
        return()=>{activ=false;};
    },[an,luna,reincarcare]);
    useEffect(()=>{setSelectate(new Set());},[sectiune,sectie,departament,cautare,inactive]);
    const pozitii=(plan?.pozitii??[]).filter(p=>{
        const a=date.angajati.find(a=>a.id===p.angajat_id),termen=cautare.trim().toLocaleLowerCase();
        return p.sectiune===sectiune&&(inactive||p.activ)&&(!sectie||a?.sectie_id===sectie)&&(!departament||a?.departament_id===departament)&&(!termen||`${nume(p.angajat_id)} ${a?.marca??''}`.toLocaleLowerCase().includes(termen));
    });
    const dePublicat=pozitii.filter(p=>p.activ&&(!selectate.size||selectate.has(p.id)));
    async function salveaza(pozitii:datebaza.PozitiePlanificare[]) {
        if(blocat.current)return;blocat.current=true;setOcupat(true);onOcupat(true);setEroare('');setMesaj('');
        try{const r=await SalveazaPlanificare(an,luna,pozitii);if(r.eroare)setEroare(t(r.eroare.cod,r.eroare.param??{}));else if(r.planificare){setPlan(r.planificare);setEditare(null);setMesaj(t('planificare.salvata'));}else setEroare(t('date.eroare'));}
        catch{setEroare(t('date.eroare'));}finally{blocat.current=false;setOcupat(false);onOcupat(false);}
    }
    async function previzualizeaza(inlocuieste=false) {
        if(blocat.current||!dePublicat.length)return;blocat.current=true;setOcupat(true);onOcupat(true);setEroare('');setMesaj('');
        const cerere=new datebaza.CererePublicarePlan({an,luna,sectiune,pozitii_id:dePublicat.map(p=>p.id),inlocuieste,token:''});
        try{const r=await PrevizualizeazaPublicarePlan(cerere);if(r.eroare)setEroare(t(r.eroare.cod,r.eroare.param??{}));else if(r.publicare)setPublicare({cerere:new datebaza.CererePublicarePlan({...cerere,token:r.publicare.token}),rezultat:r.publicare});else setEroare(t('date.eroare'));}
        catch{setEroare(t('date.eroare'));}finally{blocat.current=false;setOcupat(false);onOcupat(false);}
    }
    async function publica() {
        if(!publicare||blocat.current)return;blocat.current=true;setOcupat(true);onOcupat(true);setEroare('');
        try{
            const r=await PublicaPlanificare(publicare.cerere);
            if(r.eroare){setEroare(t(r.eroare.cod,r.eroare.param??{}));setPublicare(null);return;}
            if(!r.publicare){setEroare(t('date.eroare'));return;}
            setPublicare(null);setMesaj(t('planificare.publicata',{numar:r.publicare.publicabile,identice:r.publicare.identice,blocate:r.publicare.blocate}));setSelectate(new Set());
            const p=await CitestePlanificare(an,luna);if(p.eroare)setEroare(t(p.eroare.cod,p.eroare.param??{}));else if(p.planificare)setPlan(p.planificare);
        }catch{setEroare(t('date.eroare'));}finally{blocat.current=false;setOcupat(false);onOcupat(false);}
    }
    return <section className="modul-planificare">
        <div className="unelte-calendar"><form onSubmit={e=>{e.preventDefault();const a=Number(anText);if(Number.isInteger(a)&&a>=1&&a<=9999)setAn(a);else setEroare(t('sarbatori.an_invalid'));}}>
            <label>{t('sarbatori.an')}<input name="an_planificare" type="number" min="1" max="9999" required value={anText} disabled={ocupat} onChange={e=>setAnText(e.target.value)}/></label><button className="buton secundar" disabled={ocupat}>{t('sarbatori.arata')}</button>
            <label>{t('pontaj.luna')}<select name="luna_planificare" value={luna} disabled={ocupat} onChange={e=>setLuna(Number(e.target.value))}>{Array.from({length:12},(_,i)=><option key={i} value={i+1}>{new Intl.DateTimeFormat(i18n.language,{month:'long',timeZone:'UTC'}).format(new Date(Date.UTC(2026,i,1)))}</option>)}</select></label>
        </form><button className="buton secundar" disabled={ocupat||incarcare} onClick={()=>setReincarcare(x=>x+1)}>{t('planificare.reincarca')}</button></div>
        <div className="taburi-sectiuni" role="group" aria-label={t('planificare.sectiuni')}>{['GARZI','DUPA_AMIAZA'].map(s=><button key={s} className={`buton-tab ${sectiune===s?'activ':''}`} aria-pressed={sectiune===s} disabled={ocupat} onClick={()=>setSectiune(s)}>{t(`planificare.sectiune_${s}`)}</button>)}</div>
        <p className="nota-calendar">{t('planificare.indiciu')}</p>
        {eroare&&!editare&&!publicare&&<p role="alert" className="mesaj eroare">{eroare}</p>}{mesaj&&<p role="status" className="mesaj succes">{mesaj}</p>}
        {incarcare?<p role="status">{t('ui.incarcare')}</p>:plan&&<>
            <div className="unelte-grila filtre-planificare">
                <select aria-label={t('ui.sectie')} value={sectie} disabled={ocupat} onChange={e=>{setSectie(e.target.value);setDepartament('');}}><option value="">{t('ui.toate_sectiile')}</option>{date.sectii.map(s=><option key={s.id} value={s.id}>{denumireSectie(s,i18n.language)}</option>)}</select>
                <select aria-label={t('ui.departament')} value={departament} disabled={ocupat} onChange={e=>setDepartament(e.target.value)}><option value="">{t('ui.toate_departamentele')}</option>{(date.departamente??[]).filter(d=>!sectie||d.sectie_id===sectie).map(d=><option key={d.id} value={d.id}>{d.denumire[i18n.language]||d.denumire.ro}</option>)}</select>
                <input name="cautare_planificare" aria-label={t('ui.cautare')} placeholder={t('ui.cautare_indiciu')} value={cautare} disabled={ocupat} onChange={e=>setCautare(e.target.value)}/>
                <label className="bifa"><input type="checkbox" name="planificari_inactive" checked={inactive} disabled={ocupat} onChange={e=>setInactive(e.target.checked)}/>{t('planificare.arata_retrase')}</label>
            </div>
            <div className="unelte-calendar actiuni-planificare"><button className="buton primar" disabled={ocupat} onClick={()=>{setEroare('');setEditare(new datebaza.PozitiePlanificare({id:'',angajat_id:'',data:`${an}-${String(luna).padStart(2,'0')}-01`,schimb_id:'',sectiune,activ:true,revizie:'',stare:'NEPUBLICATA',publicat_la:''}));}}>{t('planificare.adauga')}</button>
                <button className="buton secundar" disabled={ocupat||!dePublicat.length} onClick={()=>void previzualizeaza()}>{t('planificare.previzualizeaza')}</button>
                <button className="legatura" disabled={ocupat} onClick={()=>onDeschidePontaj(an,luna)}>{t('planificare.deschide_pontaj')}</button></div>
            <p className="nota-calendar">{t('planificare.scop_publicare',{numar:dePublicat.length})}</p>
            <div className="panou tabel-derulant"><table className="tabel-planificare"><thead><tr><th><input type="checkbox" aria-label={t('planificare.selecteaza_toate')} disabled={ocupat||!pozitii.some(p=>p.activ)} checked={pozitii.some(p=>p.activ)&&pozitii.filter(p=>p.activ).every(p=>selectate.has(p.id))} onChange={e=>setSelectate(e.target.checked?new Set(pozitii.filter(p=>p.activ).map(p=>p.id)):new Set())}/></th>{['angajat','data','program','stare','actiuni'].map(c=><th key={c}>{t(c==='angajat'||c==='stare'||c==='actiuni'?`ui.${c}`:`planificare.${c}`)}</th>)}</tr></thead><tbody>
                {pozitii.map(p=><tr key={p.id} className={sarbatori.has(p.data)?'plan-zi-libera':[0,6].includes(new Date(`${p.data}T12:00:00Z`).getUTCDay())?'plan-weekend':''}><td><input type="checkbox" aria-label={t('planificare.selecteaza',{nume:nume(p.angajat_id),data:p.data})} checked={selectate.has(p.id)} disabled={ocupat||!p.activ} onChange={e=>setSelectate(v=>{const n=new Set(v);if(e.target.checked)n.add(p.id);else n.delete(p.id);return n;})}/></td><th scope="row">{nume(p.angajat_id)}</th><td>{p.data}{sarbatori.has(p.data)&&<small>{sarbatori.get(p.data)}</small>}</td><td>{textSchimb(p.schimb_id)}</td><td><span className={`stare-plan stare-plan-${p.stare}`}>{t(`planificare.stare_${p.stare}`)}</span>{p.publicat_la&&<small>{new Intl.DateTimeFormat(i18n.language,{dateStyle:'short',timeStyle:'short'}).format(new Date(p.publicat_la))}</small>}</td>
                    <td><div className="actiuni-rand"><button className="legatura" disabled={ocupat||!p.activ} onClick={()=>{setEroare('');setEditare(p);}}>{t('ui.editeaza')}</button><button className="legatura discreta" disabled={ocupat} onClick={()=>{if(window.confirm(t(p.activ?'planificare.confirma_retragerea':'planificare.confirma_reactivarea')))void salveaza([new datebaza.PozitiePlanificare({...p,activ:!p.activ})]);}}>{t(p.activ?'planificare.retrage':'ui.reactiveaza')}</button></div></td></tr>)}
            </tbody></table>{!pozitii.length&&<p className="fara-rezultate">{t('planificare.fara_pozitii')}</p>}</div>
        </>}
        {editare&&<FormularPlanificare key={editare.id||'nou'} pozitie={editare} an={an} luna={luna} date={date} schimburi={schimburi} sarbatori={sarbatori} ocupat={ocupat} eroare={eroare} inchide={()=>{setEditare(null);setEroare('');}} salveaza={salveaza}/>}
        {publicare&&<Dialog titlu={t('planificare.titlu_publicare')} ocupat={ocupat} inchide={()=>{setPublicare(null);setEroare('');}}><div className="corp-dialog previzualizare-publicare">
            <p>{t('planificare.publicare_indiciu')}</p>{eroare&&<p role="alert" className="mesaj eroare">{eroare}</p>}
            <label className="bifa"><input type="checkbox" name="inlocuieste_pontaj" checked={publicare.cerere.inlocuieste} disabled={ocupat} onChange={e=>void previzualizeaza(e.target.checked)}/>{t('planificare.permite_inlocuire')}</label>
            <p><strong>{t('planificare.rezumat_publicare',{noi:publicare.rezultat.publicabile,identice:publicare.rezultat.identice,blocate:publicare.rezultat.blocate})}</strong></p>
            <div className="tabel-derulant"><table><thead><tr>{['angajat','data','inainte','dupa','rezultat'].map(c=><th key={c}>{t(c==='angajat'?'ui.angajat':`planificare.${c}`)}</th>)}</tr></thead><tbody>{publicare.rezultat.randuri.map(r=><tr key={r.pozitie_id}><th scope="row">{nume(r.angajat_id)}</th><td>{r.data}</td><td>{r.inainte||'—'}</td><td>{textSchimb(r.schimb_id)}</td><td>{t(`planificare.publicare_${r.stare}`)}{r.motiv&&<small>{t(`planificare.motiv_${r.motiv}`)}</small>}</td></tr>)}</tbody></table></div>
        </div><footer className="subsol-dialog"><button className="buton secundar" disabled={ocupat} onClick={()=>{setPublicare(null);setEroare('');}}>{t('ui.anuleaza')}</button><button className="buton primar" disabled={ocupat||!(publicare.rezultat.publicabile+publicare.rezultat.identice)} onClick={()=>void publica()}>{t('planificare.publica_eligibile')}</button></footer></Dialog>}
    </section>;
}

function FormularPlanificare({pozitie,an,luna,date,schimburi,sarbatori,ocupat,eroare,inchide,salveaza}:{pozitie:datebaza.PozitiePlanificare;an:number;luna:number;date:datebaza.Date;schimburi:datebaza.Schimb[];sarbatori:Map<string,string>;ocupat:boolean;eroare:string;inchide:()=>void;salveaza:(p:datebaza.PozitiePlanificare[])=>Promise<void>}) {
    const {t,i18n}=useTranslation();
    const [angajat,setAngajat]=useState(pozitie.angajat_id);const [schimb,setSchimb]=useState(pozitie.schimb_id);
    const [deLa,setDeLa]=useState(pozitie.data);const [panaLa,setPanaLa]=useState(pozitie.data);const [excluse,setExcluse]=useState<Set<string>>(new Set());
    const [validare,setValidare]=useState('');
    const initial=useRef(JSON.stringify([angajat,schimb,deLa,panaLa,[]]));
    const renunta=()=>{if(!ocupat&&(initial.current===JSON.stringify([angajat,schimb,deLa,panaLa,[...excluse]])||window.confirm(t('ui.renunta'))))inchide();};
    const prima=`${an}-${String(luna).padStart(2,'0')}-01`,ultima=`${an}-${String(luna).padStart(2,'0')}-${new Date(Date.UTC(an,luna,0)).getUTCDate()}`;
    const programePentru=(id:string)=>{const a=date.angajati.find(a=>a.id===id),c=date.categorii_program?.find(c=>c.id===a?.categorie_program_id);return schimburi.filter(s=>s.activ&&inSectiune(s,pozitie.sectiune)&&(!c||(c.activ&&c.schimburi_id.includes(s.id))));};
    const programe=programePentru(angajat);
    const zile=Array.from({length:new Date(Date.UTC(an,luna,0)).getUTCDate()},(_,i)=>`${an}-${String(luna).padStart(2,'0')}-${String(i+1).padStart(2,'0')}`).filter(d=>d>=deLa&&d<=panaLa);
    return <Dialog titlu={t(pozitie.id?'planificare.editeaza':'planificare.adauga')} ocupat={ocupat} inchide={renunta}><form onSubmit={e=>{
        e.preventDefault();setValidare('');const selectate=zile.filter(d=>!excluse.has(d));
        if(!angajat||!schimb||!selectate.length){setValidare(t('planificare.campuri_obligatorii'));return;}
        void salveaza(selectate.map(data=>new datebaza.PozitiePlanificare({...pozitie,angajat_id:angajat,schimb_id:schimb,data})));
    }}><div className="corp-dialog">{(validare||eroare)&&<p role="alert" className="mesaj eroare">{validare||eroare}</p>}<fieldset className="campuri" disabled={ocupat}>
        <label className="latime-completa">{t('ui.angajat')}<select name="plan_angajat" required disabled={Boolean(pozitie.id)} value={angajat} onChange={e=>{const id=e.target.value;setAngajat(id);const lista=programePentru(id);setSchimb(lista.some(s=>s.id===schimb)?schimb:lista.find(s=>s.cod===(pozitie.sectiune==='GARZI'?'G24':'S7'))?.id||lista[0]?.id||'');}}><option value="">{t('planificare.alege_angajat')}</option>{date.angajati.filter(a=>a.activ&&(programePentru(a.id).length||a.id===angajat)).map(a=><option key={a.id} value={a.id}>{a.nume} {a.prenume} · {a.marca}</option>)}</select></label>
        <label className="latime-completa">{t('planificare.program')}<select name="plan_schimb" required value={schimb} onChange={e=>setSchimb(e.target.value)}><option value="">{t('planificare.alege_program')}</option>{programe.map(s=><option key={s.id} value={s.id}>{s.cod} · {s.denumire[i18n.language]||s.denumire.ro} · {oraText(s.ora_start)} · {durata(s.durata_minute,t)}</option>)}</select><small>{t('planificare.programe_indiciu')}</small></label>
        <label>{t('absente.de_la')}<input name="plan_de_la" type="date" min={prima} max={ultima} required disabled={Boolean(pozitie.id)} value={deLa} onChange={e=>{setDeLa(e.target.value);if(e.target.value>panaLa)setPanaLa(e.target.value);setExcluse(new Set());}}/></label>
        <label>{t('absente.pana_la')}<input name="plan_pana_la" type="date" min={deLa} max={ultima} required disabled={Boolean(pozitie.id)} value={panaLa} onChange={e=>{setPanaLa(e.target.value);setExcluse(new Set());}}/></label>
    </fieldset>
    {!pozitie.id&&<><button type="button" className="legatura" disabled={ocupat} onClick={()=>setExcluse(new Set(zile.filter(d=>sarbatori.has(d)||[0,6].includes(new Date(`${d}T12:00:00Z`).getUTCDay()))))}>{t('planificare.selecteaza_lv')}</button><div className="zile-absenta">{zile.map(d=><label key={d} className={sarbatori.has(d)?'zi-sarbatoare':[0,6].includes(new Date(`${d}T12:00:00Z`).getUTCDay())?'zi-weekend':''}><input type="checkbox" disabled={ocupat} checked={!excluse.has(d)} onChange={e=>setExcluse(v=>{const n=new Set(v);if(e.target.checked)n.delete(d);else n.add(d);return n;})}/>{d}{sarbatori.get(d)?` · ${sarbatori.get(d)}`:''}</label>)}</div></>}
    </div><footer className="subsol-dialog"><button type="button" className="buton secundar" disabled={ocupat} onClick={renunta}>{t('ui.anuleaza')}</button><button className="buton primar" disabled={ocupat}>{t('planificare.salveaza_plan')}</button></footer></form></Dialog>;
}
