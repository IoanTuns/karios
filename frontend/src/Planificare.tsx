import {useEffect, useMemo, useRef, useState, type KeyboardEvent} from 'react';
import {useTranslation} from 'react-i18next';
import {CitesteCalendar, CitestePlanificare, CitesteSchita, ListaSchimburi, PrevizualizeazaPublicarePlan, PublicaPlanificare, SalveazaPlanificare} from '../wailsjs/go/main/App';
import {datebaza} from '../wailsjs/go/models';
import {Dialog, denumireSectie, denumireDepartament, durata} from './Formulare';
import {oraText} from './Schimburi';

export function programPlanificat(p:datebaza.Schimb) {
    return p.tip==='GARDA'||p.tip==='NOAPTE'||p.regim==='NELUCRATOARE'||p.regim==='INTRARE_ZI_LIBERA'||(p.ora_start>=720&&p.ora_start<1320);
}
export function inSectiune(p:datebaza.Schimb,sectiune:string) {
    if(sectiune==='TOATE'||!sectiune)return programPlanificat(p);
    const garda=p.tip==='GARDA'||p.tip==='NOAPTE'||p.regim==='NELUCRATOARE'||p.regim==='INTRARE_ZI_LIBERA';
    return sectiune==='GARZI'?garda:!garda&&p.ora_start>=720&&p.ora_start<1320;
}
type PublicareCurenta={cerere:datebaza.CererePublicarePlan;rezultat:datebaza.PublicarePlan};
type Pozitie={rand:number;zi:number};
type Selectie={inceput:Pozitie;sfarsit:Pozitie};

export function Planificare({date,initial,onOcupat,onDeschidePontaj}:{date:datebaza.Date;initial?:{an:number;luna:number};onOcupat:(v:boolean)=>void;onDeschidePontaj:(an:number,luna:number)=>void}) {
    const {t,i18n}=useTranslation();
    const [an,setAn]=useState(initial?.an??new Date().getFullYear());
    const [luna,setLuna]=useState(initial?.luna??new Date().getMonth()+1);
    const [anText,setAnText]=useState(String(an));
    const [mod,setMod]=useState<'grila'|'lista'>('grila');
    const [sectiune,setSectiune]=useState('TOATE');
    const [sectie,setSectie]=useState('');
    const [departament,setDepartament]=useState('');
    const [cautare,setCautare]=useState('');
    const [inactive,setInactive]=useState(false);
    const [plan,setPlan]=useState<datebaza.PlanificareLuna|null>(null);
    const [schita,setSchita]=useState<datebaza.SchitaLuna|null>(null);
    const [schimburi,setSchimburi]=useState<datebaza.Schimb[]>([]);
    const [calendar,setCalendar]=useState<datebaza.CalendarAn|null>(null);
    const [selectateLista,setSelectateLista]=useState<Set<string>>(new Set());
    const [editare,setEditare]=useState<datebaza.PozitiePlanificare|null>(null);
    const [publicare,setPublicare]=useState<PublicareCurenta|null>(null);
    const [eroare,setEroare]=useState('');
    const [mesaj,setMesaj]=useState('');
    const [ocupat,setOcupat]=useState(false);
    const [incarcare,setIncarcare]=useState(true);
    const [reincarcare,setReincarcare]=useState(0);
    const [modificari,setModificari]=useState<Map<string,datebaza.PozitiePlanificare>>(new Map());
    const [istoricInapoi,setIstoricInapoi]=useState<Map<string,datebaza.PozitiePlanificare>[]>([]);
    const [istoricInainte,setIstoricInainte]=useState<Map<string,datebaza.PozitiePlanificare>[]>([]);
    const [selectie,setSelectie]=useState<Selectie>({inceput:{rand:0,zi:0},sfarsit:{rand:0,zi:0}});
    const [programAles,setProgramAles]=useState('');
    const grilaRef=useRef<HTMLDivElement>(null);
    const blocat=useRef(false);

    const nume=(id:string)=>{const a=date.angajati.find(a=>a.id===id);return a?`${a.nume} ${a.prenume}`:id;};
    const textSchimb=(id:string)=>{const s=schimburi.find(s=>s.id===id);return s?`${s.cod} · ${s.denumire[i18n.language]||s.denumire.ro} · ${durata(s.durata_minute,t)}`:t('program.schimb_lipsa');};
    const sarbatori=useMemo(()=>new Map((calendar?.sarbatori??[]).map(z=>[z.data,z.denumire[i18n.language]||z.denumire.ro])),[calendar,i18n.language]);

    useEffect(()=>{
        let activ=true;setIncarcare(true);setPlan(null);setEroare('');setMesaj('');setSelectateLista(new Set());setModificari(new Map());setIstoricInapoi([]);setIstoricInainte([]);
        Promise.all([CitestePlanificare(an,luna),ListaSchimburi(),CitesteCalendar(an),CitesteSchita(an,luna)]).then(([p,s,c,k])=>{
            if(!activ)return;const err=p.eroare||s.eroare||c.eroare||k.eroare;
            if(err){setEroare(t(err.cod,err.param??{}));return;}
            if(!p.planificare||!c.calendar){setEroare(t('date.eroare'));return;}
            setPlan(p.planificare);setSchimburi(s.schimburi??[]);setCalendar(c.calendar);setSchita(k.schita??null);
        }).catch(()=>{if(activ)setEroare(t('date.eroare'));}).finally(()=>{if(activ)setIncarcare(false);});
        return()=>{activ=false;};
    },[an,luna,reincarcare]);

    const numarZile=[31,an%400===0||(an%4===0&&an%100!==0)?29:28,31,30,31,30,31,31,30,31,30,31][luna-1];
    const zile=useMemo(()=>Array.from({length:numarZile},(_,i)=>`${String(an).padStart(4,'0')}-${String(luna).padStart(2,'0')}-${String(i+1).padStart(2,'0')}`),[an,luna,numarZile]);

    const pozitiiDB=useMemo(()=>{
        const harta=new Map<string,datebaza.PozitiePlanificare>();
        for(const p of plan?.pozitii??[]){
            harta.set(`${p.angajat_id}/${p.data}`,p);
        }
        return harta;
    },[plan]);

    const absente=useMemo(()=>{
        const harta=new Map<string,datebaza.ZiAbsenta>();
        for(const z of schita?.absente??[]){
            if(z.activ)harta.set(`${z.angajat_id}/${z.data}`,z);
        }
        return harta;
    },[schita]);

    const randuri=useMemo(()=>{
        const cuPlan=new Set((plan?.pozitii??[]).map(p=>p.angajat_id));
        const termen=cautare.trim().toLocaleLowerCase();
        return date.angajati.filter(a=>{
            if(!(a.activ||cuPlan.has(a.id))||(sectie&&a.sectie_id!==sectie)||(departament&&a.departament_id!==departament))return false;
            if(termen&&!`${a.nume} ${a.prenume} ${a.marca}`.toLocaleLowerCase().includes(termen))return false;
            return true;
        });
    },[date.angajati,plan,sectie,departament,cautare]);

    const pozitiiLista=(plan?.pozitii??[]).filter(p=>{
        const a=date.angajati.find(a=>a.id===p.angajat_id),termen=cautare.trim().toLocaleLowerCase();
        return (sectiune==='TOATE'||p.sectiune===sectiune)&&(inactive||p.activ)&&(!sectie||a?.sectie_id===sectie)&&(!departament||a?.departament_id===departament)&&(!termen||`${nume(p.angajat_id)} ${a?.marca??''}`.toLocaleLowerCase().includes(termen));
    });

    const inceputRand=Math.min(selectie.inceput.rand,selectie.sfarsit.rand),sfarsitRand=Math.max(selectie.inceput.rand,selectie.sfarsit.rand);
    const inceputZi=Math.min(selectie.inceput.zi,selectie.sfarsit.zi),sfarsitZi=Math.max(selectie.inceput.zi,selectie.sfarsit.zi);

    function getPozitieCelula(angajatId:string,data:string):datebaza.PozitiePlanificare|null {
        const k=`${angajatId}/${data}`;
        if(modificari.has(k))return modificari.get(k)!;
        const dinDB=pozitiiDB.get(k);
        return dinDB&&dinDB.activ?dinDB:null;
    }

    function schimbaLunaRelativ(delta:number) {
        if(ocupat||(modificari.size>0&&!window.confirm(t('ui.renunta'))))return;
        let nouAn=an,nouaLuna=luna+delta;
        if(nouaLuna<1){nouAn--;nouaLuna=12;}
        else if(nouaLuna>12){nouAn++;nouaLuna=1;}
        if(nouAn<1||nouAn>9999)return;
        setAn(nouAn);setAnText(String(nouAn));setLuna(nouaLuna);
    }

    function aplicaProgramSelectie(schimbCod:string) {
        if(ocupat||!randuri.length)return;
        const sAles=schimburi.find(s=>s.cod===schimbCod);
        if(schimbCod&&!sAles)return;

        const noiModificari=new Map(modificari);
        setIstoricInapoi(stiva=>[...stiva.slice(-49),new Map(modificari)]);
        setIstoricInainte([]);
        let eroriCategorie=0;

        for(let r=inceputRand;r<=sfarsitRand;r++){
            const a=randuri[r];
            if(!a||!a.activ)continue;
            const categorie=date.categorii_program?.find(c=>c.id===a.categorie_program_id);

            for(let z=inceputZi;z<=sfarsitZi;z++){
                const data=zile[z];
                const k=`${a.id}/${data}`;
                const absenta=absente.get(k);
                if(absenta)continue;

                if(!schimbCod){
                    const existent=pozitiiDB.get(k);
                    if(existent&&existent.activ){
                        noiModificari.set(k,new datebaza.PozitiePlanificare({...existent,activ:false}));
                    } else {
                        noiModificari.delete(k);
                    }
                } else {
                    if(categorie&&categorie.activ&&!categorie.schimburi_id.includes(sAles!.id)){
                        eroriCategorie++;
                        continue;
                    }
                    const existent=pozitiiDB.get(k);
                    const sect=inSectiune(sAles!,'DUPA_AMIAZA')?'DUPA_AMIAZA':'GARZI';
                    if(existent){
                        noiModificari.set(k,new datebaza.PozitiePlanificare({...existent,schimb_id:sAles!.id,sectiune:sect,activ:true}));
                    } else {
                        noiModificari.set(k,new datebaza.PozitiePlanificare({id:'',angajat_id:a.id,data,schimb_id:sAles!.id,sectiune:sect,activ:true,revizie:'',stare:'NEPUBLICATA',publicat_la:''}));
                    }
                }
            }
        }
        setModificari(noiModificari);
        if(eroriCategorie>0){
            setEroare(t('program.nepermis',{nume:t('ui.angajati'),cod:schimbCod}));
        } else {
            setEroare('');setMesaj('');
        }
    }

    function anuleazaPas() {
        if(!istoricInapoi.length||ocupat)return;
        const anterioara=istoricInapoi[istoricInapoi.length-1];
        setIstoricInapoi(s=>s.slice(0,-1));
        setIstoricInainte(s=>[...s.slice(-49),new Map(modificari)]);
        setModificari(new Map(anterioara));
        setEroare('');setMesaj('');
    }

    function refaPas() {
        if(!istoricInainte.length||ocupat)return;
        const urmatoare=istoricInainte[istoricInainte.length-1];
        setIstoricInainte(s=>s.slice(0,-1));
        setIstoricInapoi(s=>[...s.slice(-49),new Map(modificari)]);
        setModificari(new Map(urmatoare));
        setEroare('');setMesaj('');
    }

    function focalizeaza(p:Pozitie,extinde:boolean) {
        const urmator={rand:Math.max(0,Math.min(randuri.length-1,p.rand)),zi:Math.max(0,Math.min(numarZile-1,p.zi))};
        setSelectie(s=>extinde?{...s,sfarsit:urmator}:{inceput:urmator,sfarsit:urmator});
        const input=grilaRef.current?.querySelector<HTMLInputElement>(`input[data-rand="${urmator.rand}"][data-zi="${urmator.zi}"]`);
        input?.focus();
        input?.scrollIntoView({block:'nearest',inline:'nearest'});
    }

    function taste(e:KeyboardEvent<HTMLInputElement>,r:number,z:number) {
        if((e.ctrlKey||e.metaKey)&&!e.altKey&&e.key.toLowerCase()==='z'){e.preventDefault();if(e.shiftKey)refaPas();else anuleazaPas();return;}
        if((e.ctrlKey||e.metaKey)&&!e.altKey&&e.key.toLowerCase()==='y'){e.preventDefault();refaPas();return;}
        const directii:Record<string,[number,number]>={ArrowLeft:[0,-1],ArrowRight:[0,1],ArrowUp:[-1,0],ArrowDown:[1,0],Enter:[e.shiftKey?-1:1,0]};
        if(directii[e.key]){e.preventDefault();const [dr,dz]=directii[e.key];focalizeaza({rand:r+dr,zi:z+dz},e.shiftKey&&e.key!=='Enter');}
        else if(e.key==='Tab'){
            let zi=z+(e.shiftKey?-1:1),rand=r;
            if(zi<0){rand--;zi=numarZile-1;}else if(zi>=numarZile){rand++;zi=0;}
            if(rand>=0&&rand<randuri.length){e.preventDefault();focalizeaza({rand,zi},false);}
        }else if(e.key==='Delete'||e.key==='Backspace'){e.preventDefault();aplicaProgramSelectie('');}
        else if(e.key==='Escape'){e.preventDefault();setSelectie({inceput:{rand:r,zi:z},sfarsit:{rand:r,zi:z}});}
    }

    async function salveazaPlanGrila() {
        if(blocat.current||!modificari.size)return;
        blocat.current=true;setOcupat(true);onOcupat(true);setEroare('');setMesaj('');
        try{
            const pozitii=Array.from(modificari.values());
            const r=await SalveazaPlanificare(an,luna,pozitii);
            if(r.eroare)setEroare(t(r.eroare.cod,r.eroare.param??{}));
            else if(r.planificare){
                setPlan(r.planificare);setModificari(new Map());setIstoricInapoi([]);setIstoricInainte([]);setMesaj(t('planificare.salvata'));
            }else setEroare(t('date.eroare'));
        }catch{setEroare(t('date.eroare'));}finally{blocat.current=false;setOcupat(false);onOcupat(false);}
    }

    async function salveazaPozitiiLista(pozitii:datebaza.PozitiePlanificare[]) {
        if(blocat.current)return;blocat.current=true;setOcupat(true);onOcupat(true);setEroare('');setMesaj('');
        try{
            const r=await SalveazaPlanificare(an,luna,pozitii);
            if(r.eroare)setEroare(t(r.eroare.cod,r.eroare.param??{}));
            else if(r.planificare){setPlan(r.planificare);setEditare(null);setMesaj(t('planificare.salvata'));}
            else setEroare(t('date.eroare'));
        }catch{setEroare(t('date.eroare'));}finally{blocat.current=false;setOcupat(false);onOcupat(false);}
    }

    const dePublicatTotal=useMemo(()=>{
        return (plan?.pozitii??[]).filter(p=>p.activ&&(sectiune==='TOATE'||p.sectiune===sectiune));
    },[plan,sectiune]);

    async function previzualizeazaPublicare(inlocuieste=false) {
        if(modificari.size>0){
            if(!window.confirm(t('planificare.salveaza_inainte_publicare')))return;
            await salveazaPlanGrila();
        }
        if(blocat.current)return;
        const pozitiiActive=(plan?.pozitii??[]).filter(p=>p.activ&&(sectiune==='TOATE'||p.sectiune===sectiune));
        if(!pozitiiActive.length){setEroare(t('planificare.fara_pozitii'));return;}
        blocat.current=true;setOcupat(true);onOcupat(true);setEroare('');setMesaj('');
        const cerere=new datebaza.CererePublicarePlan({an,luna,sectiune:sectiune==='TOATE'?'TOATE':sectiune,pozitii_id:pozitiiActive.map(p=>p.id),inlocuieste,token:''});
        try{
            const r=await PrevizualizeazaPublicarePlan(cerere);
            if(r.eroare)setEroare(t(r.eroare.cod,r.eroare.param??{}));
            else if(r.publicare)setPublicare({cerere:new datebaza.CererePublicarePlan({...cerere,token:r.publicare.token}),rezultat:r.publicare});
            else setEroare(t('date.eroare'));
        }catch{setEroare(t('date.eroare'));}finally{blocat.current=false;setOcupat(false);onOcupat(false);}
    }

    async function publica() {
        if(!publicare||blocat.current)return;blocat.current=true;setOcupat(true);onOcupat(true);setEroare('');
        try{
            const r=await PublicaPlanificare(publicare.cerere);
            if(r.eroare){setEroare(t(r.eroare.cod,r.eroare.param??{}));setPublicare(null);return;}
            if(!r.publicare){setEroare(t('date.eroare'));return;}
            setPublicare(null);setMesaj(t('planificare.publicata',{publicate:r.publicare.publicabile,identice:r.publicare.identice,blocate:r.publicare.blocate}));
            const p=await CitestePlanificare(an,luna);if(p.eroare)setEroare(t(p.eroare.cod,p.eroare.param??{}));else if(p.planificare)setPlan(p.planificare);
        }catch{setEroare(t('date.eroare'));}finally{blocat.current=false;setOcupat(false);onOcupat(false);}
    }

    const cipuriSpeciale=['S7','IN9','G24','Z12','N12'];
    const randActiv=randuri[selectie.sfarsit.rand],dataActiva=zile[selectie.sfarsit.zi];
    const cheieActiva=randActiv&&dataActiva?`${randActiv.id}/${dataActiva}`:'';
    const pozitieActiva=randActiv&&dataActiva?getPozitieCelula(randActiv.id,dataActiva):null;
    const schimbActiv=pozitieActiva?schimburi.find(s=>s.id===pozitieActiva.schimb_id):undefined;
    const absentaActiva=cheieActiva?absente.get(cheieActiva):undefined;

    return <section className="modul-planificare">
        <div className="unelte-calendar unelte-perioada">
            <form onSubmit={e=>{e.preventDefault();const a=Number(anText);if(Number.isInteger(a)&&a>=1&&a<=9999)setAn(a);else setEroare(t('sarbatori.an_invalid'));}}>
                <label>{t('sarbatori.an')}<input name="an_planificare" type="number" min="1" max="9999" required value={anText} disabled={ocupat} onChange={e=>setAnText(e.target.value)}/></label>
                <button className="buton secundar" disabled={ocupat}>{t('sarbatori.arata')}</button>
                <button type="button" className="buton secundar buton-luna-nav" aria-label={t('planificare.luna_precedenta')} disabled={ocupat} onClick={()=>schimbaLunaRelativ(-1)}>◀</button>
                <label>{t('pontaj.luna')}<select name="luna_planificare" value={luna} disabled={ocupat} onChange={e=>setLuna(Number(e.target.value))}>
                    {Array.from({length:12},(_,i)=><option key={i} value={i+1}>{new Intl.DateTimeFormat(i18n.language,{month:'long',timeZone:'UTC'}).format(new Date(Date.UTC(2026,i,1)))}</option>)}
                </select></label>
                <button type="button" className="buton secundar buton-luna-nav" aria-label={t('planificare.luna_urmatoare')} disabled={ocupat} onClick={()=>schimbaLunaRelativ(1)}>▶</button>
            </form>
            <div className="comutator-mod" role="group" aria-label={t('planificare.mod_grila')}>
                <button type="button" className={mod==='grila'?'activ':''} onClick={()=>setMod('grila')}>{t('planificare.mod_grila')}</button>
                <button type="button" className={mod==='lista'?'activ':''} onClick={()=>setMod('lista')}>{t('planificare.mod_lista')}</button>
            </div>
            {mod==='grila'&&<button className="buton primar" disabled={ocupat||!modificari.size} onClick={salveazaPlanGrila}>{t('planificare.salveaza_grila')}</button>}
            <button className="buton secundar" disabled={ocupat||!dePublicatTotal.length} onClick={()=>void previzualizeazaPublicare()}>{t('planificare.previzualizeaza')}</button>
            <button className="legatura" disabled={ocupat} onClick={()=>onDeschidePontaj(an,luna)}>{t('planificare.deschide_pontaj')}</button>
            <button className="buton secundar" disabled={ocupat||incarcare} onClick={()=>setReincarcare(x=>x+1)}>{t('planificare.reincarca')}</button>
        </div>

        <div className="unelte-grila filtre-planificare">
            <select aria-label={t('ui.sectie')} value={sectie} disabled={ocupat} onChange={e=>{setSectie(e.target.value);setDepartament('');}}>
                <option value="">{t('ui.toate_sectiile')}</option>
                {date.sectii.map(s=><option key={s.id} value={s.id}>{denumireSectie(s,i18n.language)}</option>)}
            </select>
            <select aria-label={t('ui.departament')} value={departament} disabled={ocupat} onChange={e=>setDepartament(e.target.value)}>
                <option value="">{t('ui.toate_departamentele')}</option>
                {(date.departamente??[]).filter(d=>!sectie||d.sectie_id===sectie).map(d=><option key={d.id} value={d.id}>{!sectie?`${date.sectii.find(s=>s.id===d.sectie_id)?.denumire?denumireSectie(date.sectii.find(s=>s.id===d.sectie_id)!,i18n.language)+' / ':''}`:''}{denumireDepartament(d,i18n.language)}</option>)}
            </select>
            <select aria-label={t('planificare.sectiuni')} value={sectiune} disabled={ocupat} onChange={e=>setSectiune(e.target.value)}>
                <option value="TOATE">{t('planificare.toate_sectiunile')}</option>
                <option value="GARZI">{t('planificare.sectiune_GARZI')}</option>
                <option value="DUPA_AMIAZA">{t('planificare.sectiune_DUPA_AMIAZA')}</option>
            </select>
            <input name="cautare_planificare" aria-label={t('ui.cautare')} placeholder={t('ui.cautare_indiciu')} value={cautare} disabled={ocupat} onChange={e=>setCautare(e.target.value)}/>
            {mod==='lista'&&<label className="bifa"><input type="checkbox" name="planificari_inactive" checked={inactive} disabled={ocupat} onChange={e=>setInactive(e.target.checked)}/>{t('planificare.arata_retrase')}</label>}
            {mod==='grila'&&modificari.size>0&&<span className="stare-modificari" aria-live="polite">{t('planificare.modificari_nesalvate',{numar:modificari.size})}</span>}
        </div>

        {mod==='grila'&&<div className="paleta-planificare" role="toolbar" aria-label={t('planificare.paleta_titlu')}>
            <span className="text-discret">{t('planificare.paleta_titlu')}</span>
            {cipuriSpeciale.map(cod=>{
                const s=schimburi.find(x=>x.cod===cod);
                return <button key={cod} type="button" className="cip-schimb" disabled={ocupat||!randuri.length||!s} title={s?`${s.denumire[i18n.language]||s.denumire.ro} · ${oraText(s.ora_start)}`:cod} onClick={()=>aplicaProgramSelectie(cod)}>{cod}</button>;
            })}
            <select aria-label={t('planificare.alege_program')} value={programAles} disabled={ocupat||!randuri.length} onChange={e=>{const cod=e.target.value;setProgramAles(cod);if(cod){aplicaProgramSelectie(cod);setProgramAles('');}}}>
                <option value="">{t('planificare.alege_program')}</option>
                {schimburi.filter(s=>s.activ&&programPlanificat(s)&&!cipuriSpeciale.includes(s.cod)).map(s=><option key={s.id} value={s.cod}>{s.cod} · {s.denumire[i18n.language]||s.denumire.ro}</option>)}
            </select>
            <button type="button" className="cip-schimb cip-sterge" disabled={ocupat||!randuri.length} onClick={()=>aplicaProgramSelectie('')}>{t('planificare.sterge_selectia')}</button>
            <button type="button" className="buton secundar" disabled={ocupat||!istoricInapoi.length} onClick={anuleazaPas}>{t('planificare.anuleaza_pas')}</button>
            <button type="button" className="buton secundar" disabled={ocupat||!istoricInainte.length} onClick={refaPas}>{t('planificare.refa_pas')}</button>
        </div>}

        {randActiv&&dataActiva&&<div className="inspectie-celula" role="region" aria-label={t('planificare.inspectie_titlu')}>
            <div className="detalii-celula-activa">
                <strong>{randActiv.nume} {randActiv.prenume}</strong>
                <span>{dataActiva}{sarbatori.get(dataActiva)?` · ${sarbatori.get(dataActiva)}`:''}</span>
                {absentaActiva&&<span className="eticheta-contract">{absentaActiva.cod} ({t(`absente.${absentaActiva.cod}`)||absentaActiva.cod})</span>}
                {schimbActiv&&<span className="eticheta-intrare">{schimbActiv.cod} ({schimbActiv.denumire[i18n.language]||schimbActiv.denumire.ro}) · {t(`planificare.stare_${pozitieActiva?.stare||'NEPUBLICATA'}`)}</span>}
                {!schimbActiv&&!absentaActiva&&<span className="text-discret">{t('pontaj.inspectie_gol')}</span>}
            </div>
        </div>}

        {eroare&&!editare&&!publicare&&<p role="alert" className="mesaj eroare">{eroare}</p>}
        {mesaj&&<p role="status" className="mesaj succes">{mesaj}</p>}

        {incarcare?<p role="status">{t('ui.incarcare')}</p>:plan&&mod==='grila'?<>
            <div className="panou grila-container" ref={grilaRef}>
                <table className="grila-planificare" aria-label={t('ui.planificare')}>
                    <thead>
                        <tr>
                            <th className="fix-nume">{t('ui.angajat')}</th>
                            <th className="fix-marca">{t('ui.marca')}</th>
                            {zile.map((data,z)=>{
                                const d=new Date(`${data}T12:00:00Z`),weekend=[0,6].includes(d.getUTCDay());
                                return <th key={data} className={sarbatori.has(data)?'zi-sarbatoare':weekend?'zi-weekend':''} title={sarbatori.get(data)||''} scope="col">
                                    <strong>{z+1}</strong><small>{new Intl.DateTimeFormat(i18n.language,{weekday:'short',timeZone:'UTC'}).format(d)}</small>
                                </th>;
                            })}
                            <th className="fix-total-ture">{t('planificare.total_ture')}</th>
                        </tr>
                    </thead>
                    <tbody>
                        {randuri.map((a,r)=>{
                            let totalTure=0;
                            return <tr key={a.id}>
                                <th scope="row" className="fix-nume">
                                    <strong>{a.nume} {a.prenume}</strong>
                                    {!a.activ&&<small>{t('ui.inactiv')}</small>}
                                    {a.categorie_program_id&&<small className="categorie-angajat">{date.categorii_program?.find(c=>c.id===a.categorie_program_id)?.denumire[i18n.language]||date.categorii_program?.find(c=>c.id===a.categorie_program_id)?.denumire.ro}</small>}
                                </th>
                                <td className="fix-marca">{a.marca}</td>
                                {zile.map((data,z)=>{
                                    const weekend=[0,6].includes(new Date(`${data}T12:00:00Z`).getUTCDay());
                                    const selectata=r>=inceputRand&&r<=sfarsitRand&&z>=inceputZi&&z<=sfarsitZi;
                                    const k=`${a.id}/${data}`;
                                    const absenta=absente.get(k);
                                    const p=getPozitieCelula(a.id,data);
                                    if(p&&p.activ)totalTure++;
                                    const schimb=p?schimburi.find(s=>s.id===p.schimb_id):undefined;
                                    const esteModificata=modificari.has(k);
                                    const valoareAfisata=absenta?absenta.cod:schimb?schimb.cod:'';
                                    const stareClasa=p?`plan-${(p.stare||'nepublicata').toLowerCase()}`:'';

                                    return <td key={data} className={`${sarbatori.has(data)?'zi-sarbatoare':weekend?'zi-weekend':''} ${selectata?'celula-selectata':''} ${esteModificata?'celula-modificata':''}`}>
                                        <input data-rand={r} data-zi={z} aria-label={`${a.nume} ${a.prenume} ${data}`}
                                            readOnly disabled={ocupat}
                                            className={`${absenta?'celula-absenta':''} ${stareClasa}`}
                                            value={valoareAfisata}
                                            title={[absenta?`${t('absente.celula_protejata')}: ${absenta.cod}`:'',schimb?`${schimb.cod} (${schimb.denumire[i18n.language]||schimb.denumire.ro})`:sarbatori.get(data)||data].filter(Boolean).join(' · ')}
                                            onFocus={()=>setSelectie({inceput:{rand:r,zi:z},sfarsit:{rand:r,zi:z}})}
                                            onMouseDown={e=>{if(e.shiftKey){e.preventDefault();focalizeaza({rand:r,zi:z},true);}}}
                                            onMouseEnter={e=>{if(e.buttons===1&&!ocupat)setSelectie(s=>({...s,sfarsit:{rand:r,zi:z}}));}}
                                            onKeyDown={e=>taste(e,r,z)}/>
                                    </td>;
                                })}
                                <td className="fix-total-ture">{totalTure||'—'}</td>
                            </tr>;
                        })}
                    </tbody>
                </table>
                {!randuri.length&&<p className="fara-rezultate">{t('pontaj.fara_angajati')}</p>}
            </div>
            <p className="indiciu-grila">{t('planificare.indiciu')}</p>
        </>:plan&&mod==='lista'?<>
            <div className="unelte-calendar actiuni-planificare">
                <button className="buton primar" disabled={ocupat} onClick={()=>{setEroare('');setEditare(new datebaza.PozitiePlanificare({id:'',angajat_id:'',data:`${an}-${String(luna).padStart(2,'0')}-01`,schimb_id:'',sectiune:sectiune==='TOATE'?'GARZI':sectiune,activ:true,revizie:'',stare:'NEPUBLICATA',publicat_la:''}));}}>{t('planificare.adauga')}</button>
                <p className="nota-calendar">{t('planificare.scop_publicare',{numar:pozitiiLista.filter(p=>p.activ).length})}</p>
            </div>
            <div className="panou tabel-derulant">
                <table className="tabel-planificare">
                    <thead>
                        <tr>
                            <th><input type="checkbox" aria-label={t('planificare.selecteaza_toate')} disabled={ocupat||!pozitiiLista.some(p=>p.activ)} checked={pozitiiLista.some(p=>p.activ)&&pozitiiLista.filter(p=>p.activ).every(p=>selectateLista.has(p.id))} onChange={e=>setSelectateLista(e.target.checked?new Set(pozitiiLista.filter(p=>p.activ).map(p=>p.id)):new Set())}/></th>
                            {['angajat','data','program','stare','actiuni'].map(c=><th key={c}>{t(c==='angajat'||c==='stare'||c==='actiuni'?`ui.${c}`:`planificare.${c}`)}</th>)}
                        </tr>
                    </thead>
                    <tbody>
                        {pozitiiLista.map(p=><tr key={p.id} className={sarbatori.has(p.data)?'plan-zi-libera':[0,6].includes(new Date(`${p.data}T12:00:00Z`).getUTCDay())?'plan-weekend':''}>
                            <td><input type="checkbox" aria-label={t('planificare.selecteaza',{nume:nume(p.angajat_id),data:p.data})} checked={selectateLista.has(p.id)} disabled={ocupat||!p.activ} onChange={e=>setSelectateLista(v=>{const n=new Set(v);if(e.target.checked)n.add(p.id);else n.delete(p.id);return n;})}/></td>
                            <th scope="row">{nume(p.angajat_id)}</th>
                            <td>{p.data}{sarbatori.has(p.data)&&<small>{sarbatori.get(p.data)}</small>}</td>
                            <td>{textSchimb(p.schimb_id)}</td>
                            <td><span className={`stare-plan stare-plan-${p.stare}`}>{t(`planificare.stare_${p.stare}`)}</span>{p.publicat_la&&<small>{new Intl.DateTimeFormat(i18n.language,{dateStyle:'short',timeStyle:'short'}).format(new Date(p.publicat_la))}</small>}</td>
                            <td><div className="actiuni-rand"><button className="legatura" disabled={ocupat||!p.activ} onClick={()=>{setEroare('');setEditare(p);}}>{t('ui.editeaza')}</button><button className="legatura discreta" disabled={ocupat} onClick={()=>{if(window.confirm(t(p.activ?'planificare.confirma_retragerea':'planificare.confirma_reactivarea')))void salveazaPozitiiLista([new datebaza.PozitiePlanificare({...p,activ:!p.activ})]);}}>{t(p.activ?'planificare.retrage':'ui.reactiveaza')}</button></div></td>
                        </tr>)}
                    </tbody>
                </table>
                {!pozitiiLista.length&&<p className="fara-rezultate">{t('planificare.fara_pozitii')}</p>}
            </div>
        </>:null}

        {editare&&<FormularPlanificare key={editare.id||'nou'} pozitie={editare} an={an} luna={luna} date={date} schimburi={schimburi} sarbatori={sarbatori} ocupat={ocupat} eroare={eroare} inchide={()=>{setEditare(null);setEroare('');}} salveaza={salveazaPozitiiLista}/>}

        {publicare&&<Dialog titlu={t('planificare.titlu_publicare',{sectiune:sectiune==='TOATE'?t('planificare.publicare_toate'):t(`planificare.sectiune_${sectiune}`)})} ocupat={ocupat} inchide={()=>{setPublicare(null);setEroare('');}}>
            <div className="corp-dialog previzualizare-publicare">
                <p>{t('planificare.publicare_indiciu')}</p>
                {eroare&&<p role="alert" className="mesaj eroare">{eroare}</p>}
                <label className="bifa"><input type="checkbox" name="inlocuieste_pontaj" checked={publicare.cerere.inlocuieste} disabled={ocupat} onChange={e=>void previzualizeazaPublicare(e.target.checked)}/>{t('planificare.permite_inlocuire')}</label>
                <p><strong>{t('planificare.rezumat_publicare',{publicabile:publicare.rezultat.publicabile,identice:publicare.rezultat.identice,blocate:publicare.rezultat.blocate})}</strong></p>
                <div className="tabel-derulant">
                    <table>
                        <thead>
                            <tr>{['angajat','data','inainte','dupa','rezultat'].map(c=><th key={c}>{t(c==='angajat'?'ui.angajat':`planificare.${c}`)}</th>)}</tr>
                        </thead>
                        <tbody>
                            {publicare.rezultat.randuri.map(r=><tr key={r.pozitie_id}>
                                <th scope="row">{nume(r.angajat_id)}</th>
                                <td>{r.data}</td>
                                <td>{r.inainte||'—'}</td>
                                <td>{textSchimb(r.schimb_id)}</td>
                                <td>{t(`planificare.publicare_${r.stare}`)}{r.motiv&&<small>{t(`planificare.motiv_${r.motiv}`)}</small>}</td>
                            </tr>)}
                        </tbody>
                    </table>
                </div>
            </div>
            <footer className="subsol-dialog">
                <button className="buton secundar" disabled={ocupat} onClick={()=>{setPublicare(null);setEroare('');}}>{t('ui.anuleaza')}</button>
                <button className="buton primar" disabled={ocupat||!(publicare.rezultat.publicabile+publicare.rezultat.identice)} onClick={()=>void publica()}>{t('planificare.publica_eligibile',{numar:publicare.rezultat.publicabile+publicare.rezultat.identice})}</button>
            </footer>
        </Dialog>}
    </section>;
}

function FormularPlanificare({pozitie,an,luna,date,schimburi,sarbatori,ocupat,eroare,inchide,salveaza}:{pozitie:datebaza.PozitiePlanificare;an:number;luna:number;date:datebaza.Date;schimburi:datebaza.Schimb[];sarbatori:Map<string,string>;ocupat:boolean;eroare:string;inchide:()=>void;salveaza:(p:datebaza.PozitiePlanificare[])=>Promise<void>}) {
    const {t,i18n}=useTranslation();
    const [angajat,setAngajat]=useState(pozitie.angajat_id);
    const [schimb,setSchimb]=useState(pozitie.schimb_id);
    const [deLa,setDeLa]=useState(pozitie.data);
    const [panaLa,setPanaLa]=useState(pozitie.data);
    const [excluse,setExcluse]=useState<Set<string>>(new Set());
    const [validare,setValidare]=useState('');
    const initial=useRef(JSON.stringify([angajat,schimb,deLa,panaLa,[]]));
    const renunta=()=>{if(!ocupat&&(initial.current===JSON.stringify([angajat,schimb,deLa,panaLa,[...excluse]])||window.confirm(t('ui.renunta'))))inchide();};
    const prima=`${an}-${String(luna).padStart(2,'0')}-01`,ultima=`${an}-${String(luna).padStart(2,'0')}-${new Date(Date.UTC(an,luna,0)).getUTCDate()}`;
    const programePentru=(id:string)=>{const a=date.angajati.find(a=>a.id===id),c=date.categorii_program?.find(c=>c.id===a?.categorie_program_id);return schimburi.filter(s=>s.activ&&inSectiune(s,pozitie.sectiune)&&(!c||(c.activ&&c.schimburi_id.includes(s.id))));};
    const programe=programePentru(angajat);
    const zile=Array.from({length:new Date(Date.UTC(an,luna,0)).getUTCDate()},(_,i)=>`${an}-${String(luna).padStart(2,'0')}-${String(i+1).padStart(2,'0')}`).filter(d=>d>=deLa&&d<=panaLa);

    return <Dialog titlu={t(pozitie.id?'planificare.editeaza':'planificare.adauga')} ocupat={ocupat} inchide={renunta}>
        <form onSubmit={e=>{
            e.preventDefault();setValidare('');const selectate=zile.filter(d=>!excluse.has(d));
            if(!angajat||!schimb||!selectate.length){setValidare(t('planificare.campuri_obligatorii'));return;}
            const s=schimburi.find(x=>x.id===schimb);
            const sect=s&&inSectiune(s,'DUPA_AMIAZA')?'DUPA_AMIAZA':'GARZI';
            void salveaza(selectate.map(data=>new datebaza.PozitiePlanificare({...pozitie,angajat_id:angajat,schimb_id:schimb,sectiune:sect,data})));
        }}>
            <div className="corp-dialog">
                {(validare||eroare)&&<p role="alert" className="mesaj eroare">{validare||eroare}</p>}
                <fieldset className="campuri" disabled={ocupat}>
                    <label className="latime-completa">{t('ui.angajat')}<select name="plan_angajat" required disabled={Boolean(pozitie.id)} value={angajat} onChange={e=>{const id=e.target.value;setAngajat(id);const lista=programePentru(id);setSchimb(lista.some(s=>s.id===schimb)?schimb:lista.find(s=>s.cod===(pozitie.sectiune==='GARZI'?'G24':'S7'))?.id||lista[0]?.id||'');}}><option value="">{t('planificare.alege_angajat')}</option>{date.angajati.filter(a=>a.activ&&(programePentru(a.id).length||a.id===angajat)).map(a=><option key={a.id} value={a.id}>{a.nume} {a.prenume} · {a.marca}</option>)}</select></label>
                    <label className="latime-completa">{t('planificare.program')}<select name="plan_schimb" required value={schimb} onChange={e=>setSchimb(e.target.value)}><option value="">{t('planificare.alege_program')}</option>{programe.map(s=><option key={s.id} value={s.id}>{s.cod} · {s.denumire[i18n.language]||s.denumire.ro} · {oraText(s.ora_start)} · {durata(s.durata_minute,t)}</option>)}</select><small>{t('planificare.programe_indiciu')}</small></label>
                    <label>{t('absente.de_la')}<input name="plan_de_la" type="date" min={prima} max={ultima} required disabled={Boolean(pozitie.id)} value={deLa} onChange={e=>{setDeLa(e.target.value);if(e.target.value>panaLa)setPanaLa(e.target.value);setExcluse(new Set());}}/></label>
                    <label>{t('absente.pana_la')}<input name="plan_pana_la" type="date" min={deLa} max={ultima} required disabled={Boolean(pozitie.id)} value={panaLa} onChange={e=>{setPanaLa(e.target.value);setExcluse(new Set());}}/></label>
                </fieldset>
                {!pozitie.id&&<><button type="button" className="legatura" disabled={ocupat} onClick={()=>setExcluse(new Set(zile.filter(d=>sarbatori.has(d)||[0,6].includes(new Date(`${d}T12:00:00Z`).getUTCDay()))))}>{t('planificare.selecteaza_lv')}</button><div className="zile-absenta">{zile.map(d=><label key={d} className={sarbatori.has(d)?'zi-sarbatoare':[0,6].includes(new Date(`${d}T12:00:00Z`).getUTCDay())?'zi-weekend':''}><input type="checkbox" disabled={ocupat} checked={!excluse.has(d)} onChange={e=>setExcluse(v=>{const n=new Set(v);if(e.target.checked)n.delete(d);else n.add(d);return n;})}/>{d}{sarbatori.get(d)?` · ${sarbatori.get(d)}`:''}</label>)}</div></>}
            </div>
            <footer className="subsol-dialog">
                <button type="button" className="buton secundar" disabled={ocupat} onClick={renunta}>{t('ui.anuleaza')}</button>
                <button className="buton primar" disabled={ocupat}>{t('planificare.salveaza_plan')}</button>
            </footer>
        </form>
    </Dialog>;
}
