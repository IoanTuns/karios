import {useEffect, useRef, useState, type FormEvent, type ReactNode} from 'react';
import {useTranslation} from 'react-i18next';
import type {TFunction} from 'i18next';
import {CreeazaBackup, DetaliiDezvoltare, Informatii} from '../wailsjs/go/main/App';
import {datebaza, main} from '../wailsjs/go/models';

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

export function Dialog({titlu, ocupat, inchide, children, permiteMaximizare, clasa}: {
    titlu: string;
    ocupat: boolean;
    inchide: () => void;
    children: ReactNode;
    permiteMaximizare?: boolean;
    clasa?: string;
}) {
    const ref = useRef<HTMLDialogElement>(null);
    const [maximizat, setMaximizat] = useState(false);
    const {t} = useTranslation();
    useEffect(() => {
        const dialog = ref.current;
        dialog?.showModal();
        return () => { dialog?.close(); };
    }, []);
    return <dialog ref={ref} className={`dialog ${clasa||''} ${maximizat?'dialog-maximizat':''}`} aria-labelledby="titlu-dialog"
        onCancel={e => { e.preventDefault(); if (!ocupat) inchide(); }}>
        <div className="antet-dialog">
            <h2 id="titlu-dialog">{titlu}</h2>
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
    const [sectie, setSectie] = useState(angajat.sectie_id);
    const [departament, setDepartament] = useState(angajat.departament_id || '');
    const [categorieProgram,setCategorieProgram]=useState(angajat.categorie_program_id||'');
    const [contract, setContract] = useState(angajat.tip_contract);
    const [ore, setOre] = useState(angajat.norma_partiala_minute == null ? '' : String(angajat.norma_partiala_minute / 60));
    const [concediu, setConcediu] = useState(angajat.concediu_zile_an == null ? '' : String(angajat.concediu_zile_an));
    const [activ, setActiv] = useState(angajat.activ);
    const [validare, setValidare] = useState('');
    const initial = useRef(JSON.stringify([nume, prenume, marca, functie, sectie, departament, categorieProgram, contract, ore, concediu, activ]));

    function inchideCuVerificare() {
        const modificat = initial.current !== JSON.stringify([nume, prenume, marca, functie, sectie, departament, categorieProgram, contract, ore, concediu, activ]);
        if (!ocupat && (!modificat || window.confirm(t('ui.renunta')))) inchide();
    }

    async function trimite(e: FormEvent) {
        e.preventDefault(); setValidare('');
        if (!nume.trim() || !prenume.trim() || !marca.trim()) { setValidare(t('date.campuri_obligatorii')); return; }
        const minute = contract === 'PARTIALA' ? minuteDinOre(ore) : null;
        if (contract === 'PARTIALA' && minute === null) { setValidare(t('date.norma_invalida')); return; }
        const zile = concediu === '' ? null : Number(concediu);
        if (zile !== null && (!Number.isSafeInteger(zile) || zile < 0)) { setValidare(t('date.concediu_invalid')); return; }
        await salveaza(new datebaza.Angajat({...angajat, nume, prenume, marca, functie, sectie_id: sectie, departament_id: departament,
            categorie_program_id:categorieProgram, tip_contract: contract, norma_partiala_minute: minute, concediu_zile_an: zile, activ}));
    }

    return <Dialog titlu={t(angajat.id ? 'ui.editare_angajat' : 'ui.adauga_angajat')} ocupat={ocupat} inchide={inchideCuVerificare}>
        <form onSubmit={trimite}>
            <div className="corp-dialog"><p className="text-discret">{t('ui.formular_indiciu')}</p>
                {(validare || eroare) && <div className="mesaj eroare" role="alert">{validare || eroare}</div>}
                <fieldset disabled={ocupat} className="campuri">
                    <label>{t('ui.nume')} *<input autoFocus required value={nume} onChange={e => setNume(e.target.value)}/></label>
                    <label>{t('ui.prenume')} *<input required value={prenume} onChange={e => setPrenume(e.target.value)}/></label>
                    <label>{t('ui.marca')} *<input required value={marca} onChange={e => setMarca(e.target.value)}/></label>
                    <label>{t('ui.functie')}<input value={functie} onChange={e => setFunctie(e.target.value)}/></label>
                    <label>{t('ui.sectie')}<select value={sectie} onChange={e => { setSectie(e.target.value); setDepartament(''); }}>
                        <option value="">{t('ui.fara_sectie')}</option>{date.sectii.map(s => <option value={s.id} key={s.id}>{denumireSectie(s, i18n.language)}</option>)}
                    </select></label>
                    <label>{t('ui.departament')}<select aria-label={t('ui.departament')} value={departament} disabled={!sectie} onChange={e => setDepartament(e.target.value)}>
                        <option value="">{t('ui.fara_departament')}</option>{(date.departamente ?? []).filter(d => d.sectie_id === sectie).map(d => <option value={d.id} key={d.id}>{denumireDepartament(d, i18n.language)}</option>)}
                    </select></label>
                    <label>{t('ui.contract')} *<select value={contract} onChange={e => setContract(e.target.value)}>
                        <option value="INTREAGA">{t('ui.intreaga')}</option><option value="PARTIALA">{t('ui.partiala')}</option></select></label>
                    {contract === 'PARTIALA' ? <label>{t('ui.norma_partiala')}<input type="number" min="0" max="24" step="any" required value={ore} onChange={e => setOre(e.target.value)}/></label>
                        : <div className="norma-preluata"><span>{t('ui.norma')}</span><strong>{durata(date.norma_intreaga_minute, t)}</strong><small>{t('ui.norma_globala_indiciu')}</small></div>}
                    <label className="latime-completa">{t('program.categorie')}<select name="angajat_categorie_program" value={categorieProgram} onChange={e=>setCategorieProgram(e.target.value)}><option value="">{t('program.fara_categorie')}</option>{(date.categorii_program??[]).filter(c=>c.activ||c.id===categorieProgram).map(c=><option key={c.id} value={c.id}>{c.denumire[i18n.language]||c.denumire.ro}{!c.activ?` · ${t('ui.inactiv')}`:''}</option>)}</select><small>{t('program.atribuire_indiciu')}</small></label>
                    <label className="latime-completa">{t('ui.concediu')}<input type="number" min="0" step="1" value={concediu} onChange={e => setConcediu(e.target.value)} aria-describedby="indiciu-concediu"/>
                        <small id="indiciu-concediu">{t('ui.concediu_indiciu')}</small></label>
                    <label className="bifa latime-completa"><input type="checkbox" checked={activ} onChange={e => setActiv(e.target.checked)}/>{t('ui.activ_indiciu')}</label>
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

export function Setari({date, ocupat, salveazaNorma, adaugaSectie, adaugaDepartament}: {date: datebaza.Date; ocupat: boolean;
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
    const [eroareBackup, setEroareBackup] = useState('');
    const [dialogDezvoltare, setDialogDezvoltare] = useState(false);
    const [versiuneCurenta, setVersiuneCurenta] = useState('');
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
            else if (r.cale) setMesajBackup(t('ui.backup_creat', {cale: r.cale}));
            else setEroareBackup(t('date.eroare'));
        } catch {
            setEroareBackup(t('date.eroare'));
        } finally {
            setOcupatBackup(false);
        }
    }

    return <div className="setari">
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
        <section className="panou card-setari card-backup"><h2>{t('ui.backup_titlu')}</h2><p>{t('ui.backup_descriere')}</p>
            {eroareBackup && <p className="mesaj eroare" role="alert">{eroareBackup}</p>}
            {mesajBackup && <p className="mesaj succes" role="status">{mesajBackup}</p>}
            <button type="button" className="buton secundar" disabled={ocupat || ocupatBackup} onClick={creeazaBackup}>{t(ocupatBackup ? 'ui.salvare' : 'ui.creeaza_backup')}</button>
        </section>
        <section className="panou card-setari card-despre">
            <h2>{t('ui.despre_titlu')}</h2>
            <p>{t('ui.despre_descriere')}</p>
            <div className="despre-continut">
                <p><strong>{t('aplicatie.titlu')}</strong> — {t('aplicatie.subtitlu')}</p>
                <div className="insigne-sumar-despre">
                    <span className="cip-info">Repo: <strong>karios</strong></span>
                    <span className="cip-info">Versiune: <strong>v{versiuneCurenta || '...'}</strong></span>
                    <span className="cip-info">Licență: <strong>Apache-2.0</strong></span>
                    <span className="cip-info">100% Offline & Local</span>
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
                        ℹ️ Detalii dezvoltare, integritate și cerințe sistem
                    </button>
                </div>
            </div>
        </section>
        {dialogDezvoltare && <DialogDetaliiDezvoltare inchide={() => setDialogDezvoltare(false)} />}
    </div>;
}

export function DialogDetaliiDezvoltare({inchide}: {inchide: () => void}) {
    const [detalii, setDetalii] = useState<main.DetaliiDezvoltare | null>(null);
    const [incarcare, setIncarcare] = useState(true);
    const [copiat, setCopiat] = useState(false);

    useEffect(() => {
        let activ = true;
        DetaliiDezvoltare().then(d => {
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
        const text = `# Raport Tehnic & Dezvoltare — ${detalii.nume_aplicatie}
- Proiect: ${detalii.nume_aplicatie} (Repository: ${detalii.nume_repo})
- Versiune compilare: v${detalii.versiune}
- Git Commit: ${detalii.commit}
- Arhitectură: ${detalii.arhitectura}
- Mediu build: Go ${detalii.versiune_go} | Wails ${detalii.versiune_wails}
- Executabil: ${detalii.cale_executabil}
- SHA-256 Checksum: ${detalii.checksum_exe}
- Autor: ${detalii.autor}
- Licență: ${detalii.licenta}
- Confidențialitate: ${detalii.confidentialitate}

## Cerințe Minime de Sistem
- Sistem de operare: ${detalii.cerinte_minime.sistem_operare}
- Procesor: ${detalii.cerinte_minime.procesor}
- RAM: ${detalii.cerinte_minime.ram_minim} (Recomandat: ${detalii.cerinte_minime.ram_recomandat})
- Spațiu disc: ${detalii.cerinte_minime.spatiu_disc}
- Dependințe: ${detalii.cerinte_minime.dependinte}
- Rețea: ${detalii.cerinte_minime.retea}

## Modul Local AI & LLM
- Model utilizat: ${detalii.detalii_llm.nume_model} (${detalii.detalii_llm.dimensiune_model})
- Format: ${detalii.detalii_llm.format}
- Motor inferență: ${detalii.detalii_llm.motor_inferenta}
- Fereastră context (n_ctx): ${detalii.detalii_llm.fereastra_context} tokeni
- Port local: ${detalii.detalii_llm.port_local} (127.0.0.1)
- Stare: ${detalii.detalii_llm.stare}
- Izolare rețea: ${detalii.detalii_llm.izolare_retea}
`;
        navigator.clipboard.writeText(text).then(() => {
            setCopiat(true);
            setTimeout(() => setCopiat(false), 2500);
        });
    };

    return (
        <Dialog titlu="ℹ️ Detalii Tehnice de Dezvoltare, Integritate & Sistem" permiteMaximizare={true} ocupat={false} inchide={inchide}>
            <div className="corp-dialog detalii-dezvoltare-dialog" style={{ maxHeight: '72vh', overflowY: 'auto' }}>
                {incarcare ? (
                    <p className="incarcare">Se colectează datele de sistem...</p>
                ) : !detalii ? (
                    <p className="mesaj eroare">Nu s-au putut încărca detaliile de dezvoltare.</p>
                ) : (
                    <div className="continut-raport-dezvoltare">
                        <section className="sectiune-raport">
                            <h3>📦 Identificare & Versiune de Compilare</h3>
                            <div className="tabel-proprietati-tehnice">
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Nume Repository:</span>
                                    <strong className="valoare-prop cod-evidential">{detalii.nume_repo}</strong>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Aplicație:</span>
                                    <strong className="valoare-prop">{detalii.nume_aplicatie}</strong>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Versiune compilată:</span>
                                    <span className="valoare-prop insigna-versiune">v{detalii.versiune}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Git Commit (SHA):</span>
                                    <span className="valoare-prop cod-evidential" title="Hash-ul commit-ului din care s-a generat binarul">{detalii.commit}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Arhitectură & OS:</span>
                                    <span className="valoare-prop">{detalii.arhitectura}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Compilator & Framework:</span>
                                    <span className="valoare-prop">{detalii.versiune_go} · Wails {detalii.versiune_wails}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Licență Repository:</span>
                                    <strong className="valoare-prop">{detalii.licenta}</strong>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Autor & Copyright:</span>
                                    <span className="valoare-prop">{detalii.autor} (© 2026)</span>
                                </div>
                            </div>
                        </section>

                        <section className="sectiune-raport">
                            <h3>🛡️ Validare Integritate & Securitate Executabil</h3>
                            <div className="tabel-proprietati-tehnice">
                                <div className="linie-proprietate coloana-completa">
                                    <span className="eticheta-prop">Cale binar:</span>
                                    <span className="valoare-prop cale-fisier" title={detalii.cale_executabil}>{detalii.cale_executabil}</span>
                                </div>
                                <div className="linie-proprietate coloana-completa">
                                    <span className="eticheta-prop">SHA-256 Validare Integritate:</span>
                                    <span className="valoare-prop hash-integritate" title="Amprenta SHA-256 a fișierului executabil">{detalii.checksum_exe}</span>
                                </div>
                                <div className="linie-proprietate coloana-completa">
                                    <span className="eticheta-prop">Securitate Date:</span>
                                    <span className="valoare-prop">{detalii.confidentialitate}</span>
                                </div>
                            </div>
                        </section>

                        <section className="sectiune-raport">
                            <h3>💻 Cerințe Minime de Sistem</h3>
                            <div className="tabel-proprietati-tehnice">
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Sistem de operare:</span>
                                    <span className="valoare-prop">{detalii.cerinte_minime.sistem_operare}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Procesor (CPU):</span>
                                    <span className="valoare-prop">{detalii.cerinte_minime.procesor}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Memorie RAM minimă:</span>
                                    <span className="valoare-prop">{detalii.cerinte_minime.ram_minim}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Memorie RAM recomandată:</span>
                                    <strong className="valoare-prop">{detalii.cerinte_minime.ram_recomandat}</strong>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Spațiu pe disc:</span>
                                    <span className="valoare-prop">{detalii.cerinte_minime.spatiu_disc}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Dependințe Windows:</span>
                                    <span className="valoare-prop">{detalii.cerinte_minime.dependinte}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Cerințe rețea:</span>
                                    <strong className="valoare-prop">{detalii.cerinte_minime.retea}</strong>
                                </div>
                            </div>
                        </section>

                        <section className="sectiune-raport">
                            <h3>🤖 Modul Local AI & Model LLM</h3>
                            <div className="tabel-proprietati-tehnice">
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Model utilizat:</span>
                                    <strong className="valoare-prop">{detalii.detalii_llm.nume_model}</strong>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Dimensiune model pe disc:</span>
                                    <span className="valoare-prop">{detalii.detalii_llm.dimensiune_model}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Format model:</span>
                                    <span className="valoare-prop">{detalii.detalii_llm.format}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Motor inferență:</span>
                                    <span className="valoare-prop">{detalii.detalii_llm.motor_inferenta}</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Fereastră context:</span>
                                    <span className="valoare-prop">{detalii.detalii_llm.fereastra_context} tokeni</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Port & Host local:</span>
                                    <span className="valoare-prop">{detalii.detalii_llm.port_local} (127.0.0.1)</span>
                                </div>
                                <div className="linie-proprietate">
                                    <span className="eticheta-prop">Stare motor local:</span>
                                    <span className="valoare-prop insigna-stare-ai">{detalii.detalii_llm.stare}</span>
                                </div>
                                <div className="linie-proprietate coloana-completa">
                                    <span className="eticheta-prop">Securitate & Izolare:</span>
                                    <span className="valoare-prop">{detalii.detalii_llm.izolare_retea}</span>
                                </div>
                            </div>
                        </section>
                    </div>
                )}
            </div>
            <footer className="subsol-dialog">
                <button type="button" className="buton secundar" onClick={inchide}>Închide</button>
                <button type="button" className="buton primar" disabled={!detalii} onClick={copiazaRaport}>
                    {copiat ? '✓ Raport copiat!' : '📋 Copiază raportul tehnic'}
                </button>
            </footer>
        </Dialog>
    );
}
