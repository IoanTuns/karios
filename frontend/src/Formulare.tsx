import {EligibilitateProgram} from './PanouEligibilitate';
import {numeAngajat} from './numeAngajat';
import {useEffect, useRef, useState, useId, type FormEvent, type ReactNode} from 'react';
import {useTranslation} from 'react-i18next';
import {EditorNomenclator} from './Nomenclatoare';
import {notifica} from './Toast';
import {IstoricContract} from './IstoricContract';
import {RestaurareBackup} from './Restaurare';
import {useCont} from './Conturi';
import type {TFunction} from 'i18next';
import {CitesteUltimeleLiniiJurnal, CreeazaBackup, DeschideFolderJurnale, DetaliiDezvoltare, Informatii} from '../wailsjs/go/main/App';
import {datebaza, main} from '../wailsjs/go/models';
import {SetariLimba} from './SetariLimba';

export function durata(minute: number, t: TFunction): string {
    return t(minute % 60 ? 'ui.ore_minute' : 'ui.ore', {ore: Math.floor(minute / 60), minute: minute % 60});
}

export function denumireSectie(s: datebaza.Sectie, limba: string): string {
    return s.denumire[limba] || s.denumire.ro || Object.values(s.denumire).find(Boolean) || '';
}

export function denumireDepartament(d: datebaza.Departament, limba: string): string {
    return d.denumire[limba] || d.denumire.ro || Object.values(d.denumire).find(Boolean) || '';
}

function minuteDinOre(valoare: string): number | null {
    const minute = Number(valoare) * 60;
    const intreg = Math.round(minute);
    return valoare.trim() && Number.isFinite(minute) && Math.abs(minute - intreg) < 0.00001 && intreg >= 1 && intreg <= 1440 ? intreg : null;
}

export function Dialog({titlu, ocupat, inchide, children, permiteMaximizare = true, clasa, maximizatInitial = false}: {
    titlu: string;
    ocupat: boolean;
    inchide: () => void;
    children: ReactNode;
    permiteMaximizare?: boolean;
    clasa?: string;
    maximizatInitial?: boolean;
}) {
    const ref = useRef<HTMLDialogElement>(null);
    const titluId=useId();
    const [maximizat, setMaximizat] = useState(maximizatInitial);
    const stareCurenta = useRef({ocupat, inchide});
    useEffect(() => { stareCurenta.current = {ocupat, inchide}; }, [ocupat, inchide]);
    const {t} = useTranslation();
    useEffect(() => {
        const dialog = ref.current;
        dialog?.showModal();
        const laTasta = (e: KeyboardEvent) => {
            if (e.key === 'Escape' && !stareCurenta.current.ocupat && document.activeElement?.closest('dialog')===dialog) {
                e.preventDefault();
                stareCurenta.current.inchide();
            }
        };
        window.addEventListener('keydown', laTasta);
        return () => {
            dialog?.close();
            window.removeEventListener('keydown', laTasta);
        };
    }, []);
    return <dialog ref={ref} className={`dialog ${clasa||''} ${maximizat?'dialog-maximizat':''}`} aria-labelledby={titluId}
        onCancel={e => { e.preventDefault(); if (!ocupat) inchide(); }}>
        <div className="antet-dialog">
            <h2 id={titluId}>{titlu}</h2>
            <div className="actiuni-antet-dialog">
                {permiteMaximizare && (
                    <button
                        type="button"
                        className="buton-marire"
                        aria-label={t(maximizat ? 'ui.restaureaza' : 'ui.mareste')}
                        title={t(maximizat ? 'ui.restaureaza' : 'ui.mareste')}
                        disabled={ocupat}
                        onClick={() => setMaximizat(!maximizat)}>
                        {maximizat ? '🗗' : '🗖'}
                    </button>
                )}
                <button type="button" className="buton-inchidere" aria-label={t('ui.inchide')} disabled={ocupat} onClick={inchide}>×</button>
            </div>
        </div>
        {children}
    </dialog>;
}

export function FormularAngajat({angajat, date, ocupat, eroare, inchide, salveaza}: {
    angajat: datebaza.Angajat; date: datebaza.Date; ocupat: boolean; eroare: string;
    inchide: () => void; salveaza: (a: datebaza.Angajat) => Promise<void>;
}) {
    const {t, i18n} = useTranslation();
    const [nume, setNume] = useState(angajat.nume);
    const [prenume, setPrenume] = useState(angajat.prenume);
    const [marca, setMarca] = useState(angajat.marca);
    const [functie, setFunctie] = useState(angajat.functie);
    const [functieId,setFunctieId]=useState(angajat.functie_id||'');
    const [prefixId,setPrefixId]=useState(angajat.prefix_id||'');
    const [principalId,setPrincipalId]=useState(angajat.angajat_principal_id||'');
    const [sectie, setSectie] = useState(angajat.sectie_id);
    const [departament, setDepartament] = useState(angajat.departament_id || '');
    const [categorieProgram,setCategorieProgram]=useState(angajat.categorie_program_id||'');
    const [contract, setContract] = useState(angajat.contract_id||angajat.tip_contract);
    const [ore, setOre] = useState(angajat.norma_partiala_minute == null ? '' : String(angajat.norma_partiala_minute / 60));
    const [concediu, setConcediu] = useState(angajat.concediu_zile_an == null ? '' : String(angajat.concediu_zile_an));
    const [lunaSold,setLunaSold]=useState(angajat.initializare?.luna||'');
    const [oreSold,setOreSold]=useState(String((angajat.initializare?.minute||0)/60));
    const [coInitial,setCoInitial]=useState(String(angajat.initializare?.co_zile||0));
    const [activ, setActiv] = useState(angajat.activ);
    const [validare, setValidare] = useState('');
    const initial = useRef(JSON.stringify([prefixId, nume, prenume, marca, functie, functieId, principalId, sectie, departament, categorieProgram, contract, ore, concediu, activ,lunaSold,oreSold,coInitial]));

    function inchideCuVerificare() {
        const modificat = initial.current !== JSON.stringify([prefixId, nume, prenume, marca, functie, functieId, principalId, sectie, departament, categorieProgram, contract, ore, concediu, activ,lunaSold,oreSold,coInitial]);
        if (!ocupat && (!modificat || window.confirm(t('ui.renunta')))) inchide();
    }

    async function trimite(e: FormEvent) {
        e.preventDefault(); setValidare('');
        if (!nume.trim() || !prenume.trim() || !marca.trim()) { setValidare(t('date.campuri_obligatorii')); return; }
        const modNormare=date.tipuri_contract?.find(c=>c.id===contract)?.normare||(contract==='PARTIALA'?'INDIVIDUALA':'GLOBALA');
        const minute = modNormare === 'INDIVIDUALA' ? minuteDinOre(ore) : null;
        if (modNormare === 'INDIVIDUALA' && minute === null) { setValidare(t('date.norma_invalida')); return; }
        const zile = concediu === '' ? null : Number(concediu);
        if (zile !== null && (!Number.isSafeInteger(zile) || zile < 0)) { setValidare(t('date.concediu_invalid')); return; }
        const minuteSold=Math.round(Number(oreSold.replace(',','.'))*60),zileInitiale=Number(coInitial);
        if((!lunaSold&&(!!angajat.initializare||Number(oreSold.replace(',','.'))!==0||zileInitiale!==0))||(lunaSold&&(!Number.isSafeInteger(minuteSold)||!Number.isSafeInteger(zileInitiale)||zileInitiale<0))){setValidare(t('sold.initializare_invalida'));return;}
        await salveaza(new datebaza.Angajat({...angajat, prefix_id:prefixId, nume, prenume, marca, functie, functie_id:functieId, angajat_principal_id:principalId, sectie_id: sectie, departament_id: departament,
            categorie_program_id:categorieProgram, contract_id:contract, tip_contract: modNormare==='INDIVIDUALA'?'PARTIALA':'INTREAGA', norma_partiala_minute: minute, concediu_zile_an: zile, activ,initializare:lunaSold?{luna:lunaSold,minute:minuteSold,co_zile:zileInitiale,revizie:angajat.initializare?.revizie||0}:undefined}));
    }

    return <Dialog titlu={t(angajat.id ? 'ui.editare_angajat' : 'ui.adauga_angajat')} ocupat={ocupat} inchide={inchideCuVerificare}>
        <form onSubmit={trimite} onInvalidCapture={e=>(e.target as HTMLElement).closest('details')?.setAttribute('open','')} onKeyDown={e => {
            if (e.ctrlKey && e.key === 'Enter') {
                e.preventDefault();
                (e.currentTarget as HTMLFormElement).requestSubmit();
            }
        }}>
            <div className="corp-dialog"><p className="text-discret">{t('ui.formular_indiciu')}</p>
                {(validare || eroare) && <div className="mesaj eroare" role="alert">{validare || eroare}</div>}
                <fieldset disabled={ocupat} className="campuri grup-formular"><legend>{t('ui28.identitate')}</legend>
                    <label>{t('nomenclator.prefix')}<select name="angajat_prefix" value={prefixId} onChange={e=>setPrefixId(e.target.value)}><option value="">{t('nomenclator.fara_prefix')}</option>{(date.prefixe??[]).filter(p=>p.activ||p.id===prefixId).map(p=><option key={p.id} value={p.id}>{p.denumire[i18n.language]||p.denumire.ro||p.denumire.en}{!p.activ?` (${t('ui.inactiv')})`:''}</option>)}</select>{angajat.angajat_principal_id&&!prefixId&&angajat.prefix&&<small>{t('nomenclator.prefix_mostenit')}</small>}</label>
                    <label>{t('ui.nume')} *<input autoFocus required name="angajat_nume" value={nume} onChange={e => setNume(e.target.value)}/></label>
                    <label>{t('ui.prenume')} *<input required name="angajat_prenume" value={prenume} onChange={e => setPrenume(e.target.value)}/></label>
                    <label>{t('ui.marca')} *<input required name="angajat_marca" value={marca} onChange={e => setMarca(e.target.value)}/></label>
                    <label className="latime-completa">{t('ui.garda_principal')}<select name="angajat_principal" value={principalId} onChange={e=>{setPrincipalId(e.target.value);if(e.target.value&&date.tipuri_contract?.some(c=>c.id==='GARDA'&&c.activ))setContract('GARDA');}}><option value="">{t('ui.garda_fara_principal')}</option>{date.angajati.filter(a=>a.id!==angajat.id&&!a.angajat_principal_id&&(a.activ||a.id===principalId)).map(a=><option key={a.id} value={a.id}>{numeAngajat(a)} · {a.marca}</option>)}</select><small>{t('ui.garda_principal_indiciu')}</small>{(principalId||date.tipuri_contract?.find(c=>c.id===contract)?.garda||contract==='GARDA')&&<small>{t(principalId?'ui28.marca_principala':'ui28.garda_independenta')}</small>}</label>
</fieldset><fieldset disabled={ocupat} className="campuri grup-formular"><legend>{t('ui28.incadrare')}</legend><label>{t('ui.functie')}<select name="angajat_functie" value={functieId|| (functie?'__anterioara':'')} onChange={e=>{const id=e.target.value;if(id==='__anterioara')return;setFunctieId(id);setFunctie('');}}><option value="">{t('nomenclator.fara_functie')}</option>{!functieId&&functie&&<option value="__anterioara">{functie}</option>}{(date.functii??[]).filter(f=>f.activ||f.id===functieId).map(f=><option key={f.id} value={f.id}>{f.denumire[i18n.language]||f.denumire.ro||f.denumire.en}{!f.activ?` (${t('ui.inactiv')})`:''}</option>)}</select><small>{t('nomenclator.functii_indiciu')}</small></label>
                    <label>{t('ui.sectie')}<select value={sectie} onChange={e => { setSectie(e.target.value); setDepartament(''); }}>
                        <option value="">{t('ui.fara_sectie')}</option>{date.sectii.map(s => <option value={s.id} key={s.id}>{denumireSectie(s, i18n.language)}</option>)}
                    </select></label>
                    <label>{t('ui.departament')}<select aria-label={t('ui.departament')} value={departament} disabled={!sectie} onChange={e => setDepartament(e.target.value)}>
                        <option value="">{t('ui.fara_departament')}</option>{(date.departamente ?? []).filter(d => d.sectie_id === sectie).map(d => <option value={d.id} key={d.id}>{denumireDepartament(d, i18n.language)}</option>)}
                    </select></label>
</fieldset><fieldset disabled={ocupat} className="campuri grup-formular"><legend>{t('ui28.program_contract')}</legend><label>{t('ui.contract')} *<select value={contract} onChange={e => setContract(e.target.value)}>
                        {(date.tipuri_contract?.length?date.tipuri_contract:[{id:'INTREAGA',denumire:{ro:t('ui.intreaga')},activ:true},{id:'PARTIALA',denumire:{ro:t('ui.partiala')},activ:true}]).filter(c=>c.activ||c.id===contract).map(c=><option key={c.id} value={c.id}>{(c.denumire as Record<string,string>)[i18n.language]||c.denumire.ro}</option>)}</select></label>
                    {(date.tipuri_contract?.find(c=>c.id===contract)?.normare||(contract==='PARTIALA'?'INDIVIDUALA':'GLOBALA')) === 'INDIVIDUALA' ? <label>{t('ui.norma_partiala')}<input type="number" min="0" max="24" step="any" required value={ore} onChange={e => setOre(e.target.value)}/></label>
                        : <div className="norma-preluata"><span>{t('ui.norma')}</span><strong>{(()=>{const c=date.tipuri_contract?.find(c=>c.id===contract);return c?.normare==='FARA'?t('structura.normare_FARA'):c?.normare==='LUNARA'?`${durata(c.norma_minute??0,t)} / ${t('structura.luna')}`:durata(c?.normare==='ZILNICA'?(c.norma_minute??0):date.norma_intreaga_minute,t)})()}</strong><small>{t('ui.norma_globala_indiciu')}</small></div>}
                    <label className="latime-completa">{t('program.categorie')}<select name="angajat_categorie_program" value={categorieProgram} onChange={e=>setCategorieProgram(e.target.value)}><option value="">{t('program.fara_categorie')}</option>{(date.categorii_program??[]).filter(c=>c.activ||c.id===categorieProgram).map(c=><option key={c.id} value={c.id}>{c.denumire[i18n.language]||c.denumire.ro}{!c.activ?` · ${t('ui.inactiv')}`:''}</option>)}</select><small>{t('program.atribuire_indiciu')}</small></label><EligibilitateProgram categorie={date.categorii_program?.find(c=>c.id===categorieProgram)}/>
                    <label className="latime-completa">{t('ui.concediu')}<input type="number" min="0" step="1" value={concediu} onChange={e => setConcediu(e.target.value)} aria-describedby="indiciu-concediu"/>
                        <small id="indiciu-concediu">{t('ui.concediu_indiciu')}</small></label>
</fieldset><details className="detalii-configurare solduri-initiale"><summary>{t('ui28.solduri_initiale')} {lunaSold&&t('ui28.sold_rezumat',{luna:lunaSold,ore:oreSold,co:coInitial})}</summary><fieldset disabled={ocupat} className="campuri"><p className="latime-completa">{t('ui28.formula_sold')}</p><label>{t('sold.luna')}<input type="month" name="initializare_luna" min="2000-01" max="2100-12" value={lunaSold} onChange={e=>setLunaSold(e.target.value)}/></label>
                    <label>{t('sold.ore')}<input name="initializare_ore" inputMode="decimal" value={oreSold} onChange={e=>setOreSold(e.target.value)}/></label>
                    <label>{t('sold.co')}<input name="initializare_co" type="number" min="0" max="366" step="1" value={coInitial} onChange={e=>setCoInitial(e.target.value)}/></label>
                    <p className="latime-completa">{t('sold.indiciu')}</p><button type="button" className="legatura" onClick={()=>window.dispatchEvent(new CustomEvent('pontix-ajutor-contextual',{detail:'angajati'}))}>{t('ajutor.buton')}</button>
</fieldset></details>{angajat.id && <IstoricContract angajatId={angajat.id}/>}<fieldset disabled={ocupat} className="campuri grup-formular"><legend>{t('ui.stare')}</legend>                    <label className="bifa latime-completa"><input type="checkbox" checked={activ} onChange={e => setActiv(e.target.checked)}/>{t('ui.activ_indiciu')}</label>
                </fieldset>
            </div>
            <footer className="subsol-dialog"><button type="button" className="buton secundar" disabled={ocupat} onClick={inchideCuVerificare}>{t('ui.anuleaza')}</button>
                <button type="submit" className="buton primar" disabled={ocupat}>{t(ocupat ? 'ui.salvare' : 'ui.salveaza')}</button></footer>
        </form>
    </Dialog>;
}

export function DialogDemo({ocupat, eroare, inchide, confirma}: {ocupat: boolean; eroare: string; inchide: () => void; confirma: () => Promise<void>}) {
    const {t} = useTranslation();
    return <Dialog titlu={t('ui.demo_titlu')} ocupat={ocupat} inchide={inchide}>
        <div className="corp-dialog"><p>{t('ui.demo_descriere')}</p><p className="text-discret">{t('ui.demo_limita')}</p>
            {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}</div>
        <footer className="subsol-dialog"><button className="buton secundar" disabled={ocupat} onClick={inchide}>{t('ui.anuleaza')}</button>
            <button className="buton primar" disabled={ocupat} onClick={confirma}>{t(ocupat ? 'ui.salvare' : 'ui.demo_confirma')}</button></footer>
    </Dialog>;
}

export function Setari({date, onDate, ocupat, salveazaNorma, adaugaSectie, adaugaDepartament}: {onDate:(date:datebaza.Date)=>void;date: datebaza.Date; ocupat: boolean;
    salveazaNorma: (minute: number) => Promise<boolean>; adaugaSectie: (denumire: string) => Promise<boolean>;
    adaugaDepartament: (sectieID: string, denumire: string) => Promise<boolean>;
}) {
    const {t, i18n} = useTranslation();
    const [ore, setOre] = useState(String(date.norma_intreaga_minute / 60));
    const [sectie, setSectie] = useState('');
    const [sectieDepartament, setSectieDepartament] = useState(date.sectii[0]?.id || '');
    const [numeDepartament, setNumeDepartament] = useState('');
    const [validare, setValidare] = useState('');
    const [ocupatBackup, setOcupatBackup] = useState(false);
    const [mesajBackup, setMesajBackup] = useState('');
    const [backupuriCreate, setBackupuriCreate] = useState(0);
    const esteAdministrator = useCont()?.cont.rol === 'ADMINISTRATOR';
    const [eroareBackup, setEroareBackup] = useState('');
    const [dialogDezvoltare, setDialogDezvoltare] = useState(false);
    const [versiuneCurenta, setVersiuneCurenta] = useState('');
    // R6 (analiza UX 0.92.59): General pe file, ca pagina sa nu mai aiba ~4 ecrane. Sectiunile
    // raman in pagina (hidden), iar fila aleasa se retine local.
    const file = ['structura', 'nomenclatoare', 'backup', 'limba', 'despre'] as const;
    const [fila, setFilaStare] = useState<typeof file[number]>(() => { try { const v = localStorage.getItem('pontix-setari-fila'); return (file as readonly string[]).includes(v || '') ? v as typeof file[number] : 'structura'; } catch { return 'structura'; } });
    const setFila = (f: typeof file[number]) => { setFilaStare(f); try { localStorage.setItem('pontix-setari-fila', f); } catch { /* doar in sesiune */ } };
    useEffect(() => { setOre(String(date.norma_intreaga_minute / 60)); }, [date.norma_intreaga_minute]);
    useEffect(() => {
        let activ = true;
        Informatii().then(info => { if (activ) setVersiuneCurenta(info.versiune); }).catch(() => {});
        return () => { activ = false; };
    }, []);

    async function trimiteNorma(e: FormEvent) {
        e.preventDefault();
        const minute = minuteDinOre(ore);
        if (minute === null) { setValidare(t('date.norma_invalida')); return; }
        setValidare('');
        await salveazaNorma(minute);
    }

    async function creeazaBackup() {
        if (ocupat || ocupatBackup) return;
        setOcupatBackup(true); setMesajBackup(''); setEroareBackup('');
        try {
            const r = await CreeazaBackup();
            if (r.eroare) setEroareBackup(t(r.eroare.cod, r.eroare.param ?? {}));
            else if (r.cale) { setMesajBackup(t('ui.backup_creat', {cale: r.cale})); setBackupuriCreate(n => n + 1); }
            else setEroareBackup(t('date.eroare'));
        } catch {
            setEroareBackup(t('date.eroare'));
        } finally {
            setOcupatBackup(false);
        }
    }

    return <div className="setari">
        <div className="file-setari-generale taburi-sectiuni" role="group" aria-label={t('ui.setari')}>
            {file.map(f => <button key={f} type="button" className={`buton-tab ${fila === f ? 'activ' : ''}`} aria-pressed={fila === f} onClick={() => setFila(f)}>{t(`setari.fila_${f}`)}</button>)}
        </div>
        <div className="sectiune-setari" hidden={fila !== 'structura'}>
        <section className="panou card-setari"><h2>{t('ui.norma_titlu')}</h2><p>{t('ui.norma_descriere')}</p>
            <form onSubmit={trimiteNorma}>
                <label>{t('ui.ore_pe_zi')}<input type="number" min="0" max="24" step="any" required value={ore} disabled={ocupat} onChange={e => setOre(e.target.value)}/></label>
                {validare && <p className="mesaj eroare" role="alert">{validare}</p>}
                <button className="buton primar" disabled={ocupat}>{t(ocupat ? 'ui.salvare' : 'ui.salveaza')}</button>
            </form><p className="nota">{t('ui.norma_istoric')}</p>
        </section>
        <section className="panou card-setari"><h2>{t('ui.sectii_titlu')}</h2><p>{t('ui.sectii_descriere')}</p>
            {date.sectii.length ? <ul className="lista-sectii">{date.sectii.map(s => <li key={s.id}>{denumireSectie(s, i18n.language)}</li>)}</ul> : <p className="text-discret">{t('ui.fara_sectii')}</p>}
            <form onSubmit={async e => { e.preventDefault(); if (sectie.trim() && await adaugaSectie(sectie)) setSectie(''); }}>
                <label>{t('ui.sectie_noua')}<input required value={sectie} disabled={ocupat} onChange={e => setSectie(e.target.value)}/></label>
                <button className="buton secundar" disabled={ocupat || !sectie.trim()}>{t('ui.adauga_sectie')}</button>
            </form>
        </section>
        <section className="panou card-setari card-departamente"><h2>{t('ui.departamente')}</h2><p>{t('ui.departamente_descriere')}</p>
            <div className="organizare-departamente"><div>
                {date.sectii.length ? <ul className="arbore-departamente">{date.sectii.map(s => <li key={s.id}><strong>{denumireSectie(s, i18n.language)}</strong>
                    <ul>{(date.departamente ?? []).filter(d => d.sectie_id === s.id).map(d => <li key={d.id}>{denumireDepartament(d, i18n.language)}</li>)}
                        {!(date.departamente ?? []).some(d => d.sectie_id === s.id) && <li className="text-discret">{t('ui.fara_departamente')}</li>}</ul></li>)}</ul> : <p>{t('ui.fara_sectii')}</p>}
            </div><form onSubmit={async e => { e.preventDefault(); if (numeDepartament.trim() && sectieDepartament && await adaugaDepartament(sectieDepartament, numeDepartament)) setNumeDepartament(''); }}>
                <label>{t('ui.sectie')}<select required aria-label={t('ui.sectie_departament')} disabled={ocupat} value={sectieDepartament} onChange={e => setSectieDepartament(e.target.value)}><option value="">{t('ui.alege_sectia')}</option>{date.sectii.map(s => <option value={s.id} key={s.id}>{denumireSectie(s, i18n.language)}</option>)}</select></label>
                <label>{t('ui.departament_nou')}<input required aria-label={t('ui.departament_nou')} value={numeDepartament} disabled={ocupat} onChange={e => setNumeDepartament(e.target.value)}/></label>
                <button className="buton secundar" disabled={ocupat || !sectieDepartament || !numeDepartament.trim()}>{t('ui.adauga_departament')}</button>
            </form></div>
        </section>
        </div>
        <div className="sectiune-setari" hidden={fila !== 'backup'}>
        <section className="panou card-setari card-backup"><h2>{t('ui.backup_titlu')}</h2><p>{t('ui.backup_descriere')}</p>
            {eroareBackup && <p className="mesaj eroare" role="alert">{eroareBackup}</p>}
            {mesajBackup && <p className="mesaj succes" role="status">{mesajBackup}</p>}
            <button type="button" className="buton secundar" disabled={ocupat || ocupatBackup} onClick={creeazaBackup}>{t(ocupatBackup ? 'ui.salvare' : 'ui.creeaza_backup')}</button>
            {esteAdministrator && <RestaurareBackup reimprospatare={backupuriCreate} ocupat={ocupat || ocupatBackup}/>}
        </section>
        </div>
        <div className="sectiune-setari" hidden={fila !== 'nomenclatoare'}>
        <EditorNomenclator tip="grupuri" lista={date.grupuri_profesionale??[]} onDate={onDate} date={date}/><EditorNomenclator tip="functii" lista={date.functii??[]} onDate={onDate} date={date}/><EditorNomenclator tip="prefixe" lista={date.prefixe??[]} onDate={onDate} date={date}/><EditorNomenclator tip="contracte" lista={date.tipuri_contract??[]} onDate={onDate} date={date}/>
        </div>
        <div className="sectiune-setari" hidden={fila !== 'limba'}><SetariLimba/></div>
        <div className="sectiune-setari" hidden={fila !== 'despre'}>
        <section className="panou card-setari card-despre">
            <h2>{t('ui.despre_titlu')}</h2>
            <p>{t('ui.despre_descriere')}</p>
            <div className="despre-continut">
                <p><strong>{t('aplicatie.titlu')}</strong> — {t('aplicatie.subtitlu')}</p>
                <div className="insigne-sumar-despre">
                    <span className="cip-info">Repo: <strong>karios</strong></span>
                    <span className="cip-info">{t('ui.cip_versiune')} <strong>v{versiuneCurenta || '...'}</strong></span>
                    <span className="cip-info">{t('ui.cip_licenta')} <strong>Apache-2.0</strong></span>
                    <span className="cip-info">{t('ui.cip_offline')}</span>
                    <span className="cip-info">Model: <strong>Qwen2.5 / llama.cpp</strong></span>
                </div>
                <p className="semnatura-dezvoltator-card">{t('aplicatie.dezvoltat_de', {autor: 'Ioan Tuns'})}</p>
                <p className="drepturi-autor">© 2026 Ioan Tuns · {t('aplicatie.licenta_descriere')}</p>
                <div className="actiuni-despre">
                    <button
                        type="button"
                        className="buton primar"
                        onClick={() => setDialogDezvoltare(true)}
                    >
                        ℹ️ {t('ui.buton_detalii_dezvoltare')}
                    </button>
                </div>
            </div>
        </section>
        </div>
        {dialogDezvoltare && <DialogDetaliiDezvoltare inchide={() => setDialogDezvoltare(false)} />}
    </div>;
}

export function DialogDetaliiDezvoltare({inchide}: {inchide: () => void}) {
    const {t, i18n} = useTranslation();
    const [detalii, setDetalii] = useState<main.DetaliiDezvoltare | null>(null);
    const [incarcare, setIncarcare] = useState(true);
    const [copiat, setCopiat] = useState(false);

    useEffect(() => {
        let activ = true;
        DetaliiDezvoltare(i18n.language).then(d => {
            if (activ) {
                setDetalii(d);
                setIncarcare(false);
            }
        }).catch(() => {
            if (activ) setIncarcare(false);
        });
        return () => { activ = false; };
    }, []);

    const copiazaRaport = () => {
        if (!detalii) return;
        const e = (k: string) => t(`dezvoltare.${k}`);
        const sd = detalii.statistici_diagnoza;
        const text = [
            `# ${e('titlu')} — ${detalii.nume_aplicatie}`,
            `- ${e('repo')}: ${detalii.nume_repo}`,
            `- ${e('versiune')}: v${detalii.versiune}`,
            `- ${e('commit')}: ${detalii.commit}`,
            `- ${e('arhitectura')}: ${detalii.arhitectura}`,
            `- ${e('compilator')}: Go ${detalii.versiune_go} | Wails ${detalii.versiune_wails}`,
            `- ${e('cale_binar')}: ${detalii.cale_executabil}`,
            `- ${e('sha256')}: ${detalii.checksum_exe}`,
            `- ${e('autor')}: ${detalii.autor}`,
            `- ${e('licenta')}: ${detalii.licenta}`,
            `- ${e('securitate_date')}: ${detalii.confidentialitate}`,
            '',
            `## ${e('titlu_cerinte')}`,
            `- ${e('sistem_operare')}: ${detalii.cerinte_minime.sistem_operare}`,
            `- ${e('procesor')}: ${detalii.cerinte_minime.procesor}`,
            `- ${e('ram_minim')}: ${detalii.cerinte_minime.ram_minim}`,
            `- ${e('ram_recomandat')}: ${detalii.cerinte_minime.ram_recomandat}`,
            `- ${e('spatiu_disc')}: ${detalii.cerinte_minime.spatiu_disc}`,
            `- ${e('dependinte')}: ${detalii.cerinte_minime.dependinte}`,
            `- ${e('retea')}: ${detalii.cerinte_minime.retea}`,
            '',
            `## ${e('titlu_ai')}`,
            `- ${e('model')}: ${detalii.detalii_llm.nume_model} (${detalii.detalii_llm.dimensiune_model})`,
            `- ${e('format_model')}: ${detalii.detalii_llm.format}`,
            `- ${e('motor')}: ${detalii.detalii_llm.motor_inferenta}`,
            `- ${e('context')}: ${detalii.detalii_llm.fereastra_context}`,
            `- ${e('port')}: ${detalii.detalii_llm.port_local} (127.0.0.1)`,
            `- ${e('stare_motor')}: ${detalii.detalii_llm.stare}`,
            `- ${e('izolare')}: ${detalii.detalii_llm.izolare_retea}`,
            '',
            `## ${e('titlu_diagnoza')}`,
            `- ${e('erori')}: ${sd?.numar_erori || 0}`,
            `- ${e('interogari')}: ${sd?.total_interogari || 0}`,
            `- ${e('timp_mediu')}: ${(sd?.timp_mediu_ms || 0).toFixed(2)} ms`,
            `- ${e('timp_maxim')}: ${(sd?.timp_maxim_ms || 0).toFixed(2)} ms`,
            `- ${e('interogare_lenta')}: ${sd?.lenta_interogare || e('fara_lente')}`,
            `- ${e('fisier_jurnal')}: ${sd?.fisier_jurnal_curent || '—'}`,
            `- ${e('director_jurnale')}: ${sd?.cale_folder_jurnale || '—'}`,
        ].join('\n') + '\n';
        navigator.clipboard.writeText(text).then(() => {
            setCopiat(true);
            setTimeout(() => setCopiat(false), 2500);
        });
    };

    return (
        <Dialog titlu={`ℹ️ ${t('dezvoltare.titlu')}`} permiteMaximizare={true} ocupat={false} inchide={inchide}>
            <div className="corp-dialog detalii-dezvoltare-dialog" style={{ maxHeight: '72vh', overflowY: 'auto' }}>
                {incarcare ? (
                    <p className="incarcare">{t('dezvoltare.incarcare')}</p>
                ) : !detalii ? (
                    <p className="mesaj eroare">{t('dezvoltare.eroare')}</p>
                ) : (
                    <div className="continut-raport-dezvoltare">
                        <section className="sectiune-raport">
                            <h3>📦 {t('dezvoltare.titlu_identificare')}</h3>
                            <div className="tabel-proprietati-tehnice">
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.repo')}:</span>
                                    <strong className="valoare-prop cod-evidential">{detalii.nume_repo}</strong>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.aplicatie')}:</span>
                                    <strong className="valoare-prop">{detalii.nume_aplicatie}</strong>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.versiune')}:</span>
                                    <span className="valoare-prop insigna-versiune">v{detalii.versiune}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.commit')}:</span>
                                    <span className="valoare-prop cod-evidential" title={t('dezvoltare.commit_indiciu')}>{detalii.commit}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.arhitectura')}:</span>
                                    <span className="valoare-prop">{detalii.arhitectura}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.compilator')}:</span>
                                    <span className="valoare-prop">{detalii.versiune_go} · Wails {detalii.versiune_wails}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.licenta')}:</span>
                                    <strong className="valoare-prop">{detalii.licenta}</strong>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.autor')}:</span>
                                    <span className="valoare-prop">{detalii.autor} (© 2026)</span>
                                </div>
                            </div>
                        </section>

                        <section className="sectiune-raport">
                            <h3>🛡️ {t('dezvoltare.titlu_integritate')}</h3>
                            <div className="tabel-proprietati-tehnice">
                                <div className="linie-proprietate coloana-completa">
                                    <span className="eticheta-prop">{t('dezvoltare.cale_binar')}:</span>
                                    <span className="valoare-prop cale-fisier" title={detalii.cale_executabil}>{detalii.cale_executabil}</span>
                                </div>
                                <div className="linie-proprietate coloana-completa">
                                    <span className="eticheta-prop">{t('dezvoltare.sha256')}:</span>
                                    <span className="valoare-prop hash-integritate" title={t('dezvoltare.sha256_indiciu')}>{detalii.checksum_exe}</span>
                                </div>
                                <div className="linie-proprietate coloana-completa">
                                    <span className="eticheta-prop">{t('dezvoltare.securitate_date')}:</span>
                                    <span className="valoare-prop">{detalii.confidentialitate}</span>
                                </div>
                            </div>
                        </section>

                        <section className="sectiune-raport">
                            <h3>💻 {t('dezvoltare.titlu_cerinte')}</h3>
                            <div className="tabel-proprietati-tehnice">
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.sistem_operare')}:</span>
                                    <span className="valoare-prop">{detalii.cerinte_minime.sistem_operare}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.procesor')}:</span>
                                    <span className="valoare-prop">{detalii.cerinte_minime.procesor}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.ram_minim')}:</span>
                                    <span className="valoare-prop">{detalii.cerinte_minime.ram_minim}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.ram_recomandat')}:</span>
                                    <strong className="valoare-prop">{detalii.cerinte_minime.ram_recomandat}</strong>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.spatiu_disc')}:</span>
                                    <span className="valoare-prop">{detalii.cerinte_minime.spatiu_disc}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.dependinte')}:</span>
                                    <span className="valoare-prop">{detalii.cerinte_minime.dependinte}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.retea')}:</span>
                                    <strong className="valoare-prop">{detalii.cerinte_minime.retea}</strong>
                                </div>
                            </div>
                        </section>

                        <section className="sectiune-raport">
                            <h3>🤖 {t('dezvoltare.titlu_ai')}</h3>
                            <div className="tabel-proprietati-tehnice">
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.model')}:</span>
                                    <strong className="valoare-prop">{detalii.detalii_llm.nume_model}</strong>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.dimensiune_model')}:</span>
                                    <span className="valoare-prop">{detalii.detalii_llm.dimensiune_model}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.format_model')}:</span>
                                    <span className="valoare-prop">{detalii.detalii_llm.format}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.motor')}:</span>
                                    <span className="valoare-prop">{detalii.detalii_llm.motor_inferenta}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.context')}:</span>
                                    <span className="valoare-prop">{detalii.detalii_llm.fereastra_context} tokeni</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.port')}:</span>
                                    <span className="valoare-prop">{detalii.detalii_llm.port_local} (127.0.0.1)</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.stare_motor')}:</span>
                                    <span className="valoare-prop insigna-stare-ai">{detalii.detalii_llm.stare}</span>
                                </div>
                                <div className="linie-proprietate coloana-completa">
                                    <span className="eticheta-prop">{t('dezvoltare.izolare')}:</span>
                                    <span className="valoare-prop">{detalii.detalii_llm.izolare_retea}</span>
                                </div>
                            </div>
                        </section>

                        <section className="sectiune-raport">
                            <h3>📊 {t('dezvoltare.titlu_diagnoza')}</h3>
                            <div className="tabel-proprietati-tehnice">
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.erori')}:</span>
                                    <span className="valoare-prop" style={{ fontWeight: 600, color: (detalii.statistici_diagnoza?.numar_erori || 0) > 0 ? '#b02a24' : '#1e7040' }}>
                                        {detalii.statistici_diagnoza?.numar_erori || 0} erori înregistrate
                                    </span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.interogari')}:</span>
                                    <span className="valoare-prop">{detalii.statistici_diagnoza?.total_interogari || 0}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.timp_mediu')}:</span>
                                    <strong className="valoare-prop">{(detalii.statistici_diagnoza?.timp_mediu_ms || 0).toFixed(2)} ms</strong>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.timp_maxim')}:</span>
                                    <span className="valoare-prop">{(detalii.statistici_diagnoza?.timp_maxim_ms || 0).toFixed(2)} ms</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.interogare_lenta')}:</span>
                                    <span className="valoare-prop cod-evidential">{detalii.statistici_diagnoza?.lenta_interogare || t('dezvoltare.fara_lente')}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">{t('dezvoltare.fisier_jurnal')}:</span>
                                    <span className="valoare-prop">{detalii.statistici_diagnoza?.fisier_jurnal_curent || '—'} ({detalii.statistici_diagnoza?.dimensiune_jurnal || '0 B'})</span>
                                </div>
                                <div className="linie-proprietate coloana-completa">
                                    <span className="eticheta-prop">{t('dezvoltare.director_jurnale')}:</span>
                                    <span className="valoare-prop cale-fisier">{detalii.statistici_diagnoza?.cale_folder_jurnale || '—'}</span>
                                </div>
                                {detalii.statistici_diagnoza?.ultima_eroare && (
                                    <div className="linie-proprietate coloana-completa">
                                        <span className="eticheta-prop">{t('dezvoltare.ultima_eroare')}:</span>
                                        <span className="valoare-prop cod-evidential" style={{ color: '#b02a24' }}>
                                            {detalii.statistici_diagnoza.ultima_eroare}
                                        </span>
                                    </div>
                                )}
                            </div>
                            <div style={{ display: 'flex', gap: '10px', marginTop: '10px', flexWrap: 'wrap' }}>
                                <button
                                    type="button"
                                    className="buton secundar"
                                    onClick={() => DeschideFolderJurnale()}
                                >
                                    📁 {t('dezvoltare.deschide_jurnale')}
                                </button>
                                <button
                                    type="button"
                                    className="buton secundar"
                                    onClick={async () => {
                                        const linii = await CitesteUltimeleLiniiJurnal(100);
                                        await navigator.clipboard.writeText(linii);
                                        notifica(t('dezvoltare.linii_copiate'), 'info');
                                    }}
                                >
                                    📋 {t('dezvoltare.copiaza_linii')}
                                </button>
                            </div>
                        </section>
                    </div>
                )}
            </div>
            <footer className="subsol-dialog">
                <button type="button" className="buton secundar" onClick={inchide}>{t('ui.inchide')}</button>
                <button type="button" className="buton primar" disabled={!detalii} onClick={copiazaRaport}>
                    {copiat ? `✓ ${t('dezvoltare.raport_copiat')}` : `📋 ${t('dezvoltare.copiaza_raport')}`}
                </button>
            </footer>
        </Dialog>
    );
}
