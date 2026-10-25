import {useEffect, useRef, useState, type FormEvent, type ReactNode} from 'react';
import {useTranslation} from 'react-i18next';
import type {TFunction} from 'i18next';
import {datebaza} from '../wailsjs/go/models';

export function durata(minute: number, t: TFunction): string {
    return t(minute % 60 ? 'ui.ore_minute' : 'ui.ore', {ore: Math.floor(minute / 60), minute: minute % 60});
}

export function denumireSectie(s: datebaza.Sectie, limba: string): string {
    return s.denumire[limba] || s.denumire.ro || Object.values(s.denumire).find(Boolean) || '';
}

function minuteDinOre(valoare: string): number | null {
    const minute = Number(valoare) * 60;
    const intreg = Math.round(minute);
    return valoare.trim() && Number.isFinite(minute) && Math.abs(minute - intreg) < 0.00001 && intreg >= 1 && intreg <= 1440 ? intreg : null;
}

export function Dialog({titlu, ocupat, inchide, children}: {titlu: string; ocupat: boolean; inchide: () => void; children: ReactNode}) {
    const ref = useRef<HTMLDialogElement>(null);
    const {t} = useTranslation();
    useEffect(() => {
        const dialog = ref.current;
        dialog?.showModal();
        return () => { dialog?.close(); };
    }, []);
    return <dialog ref={ref} className="dialog" aria-labelledby="titlu-dialog"
        onCancel={e => { e.preventDefault(); if (!ocupat) inchide(); }}>
        <div className="antet-dialog"><h2 id="titlu-dialog">{titlu}</h2>
            <button type="button" className="buton-inchidere" aria-label={t('ui.inchide')} disabled={ocupat} onClick={inchide}>×</button></div>
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
    const [contract, setContract] = useState(angajat.tip_contract);
    const [ore, setOre] = useState(angajat.norma_partiala_minute == null ? '' : String(angajat.norma_partiala_minute / 60));
    const [concediu, setConcediu] = useState(angajat.concediu_zile_an == null ? '' : String(angajat.concediu_zile_an));
    const [activ, setActiv] = useState(angajat.activ);
    const [validare, setValidare] = useState('');
    const initial = useRef(JSON.stringify([nume, prenume, marca, functie, sectie, contract, ore, concediu, activ]));

    function inchideCuVerificare() {
        const modificat = initial.current !== JSON.stringify([nume, prenume, marca, functie, sectie, contract, ore, concediu, activ]);
        if (!ocupat && (!modificat || window.confirm(t('ui.renunta')))) inchide();
    }

    async function trimite(e: FormEvent) {
        e.preventDefault(); setValidare('');
        if (!nume.trim() || !prenume.trim() || !marca.trim()) { setValidare(t('date.campuri_obligatorii')); return; }
        const minute = contract === 'PARTIALA' ? minuteDinOre(ore) : null;
        if (contract === 'PARTIALA' && minute === null) { setValidare(t('date.norma_invalida')); return; }
        const zile = concediu === '' ? null : Number(concediu);
        if (zile !== null && (!Number.isSafeInteger(zile) || zile < 0)) { setValidare(t('date.concediu_invalid')); return; }
        await salveaza(new datebaza.Angajat({...angajat, nume, prenume, marca, functie, sectie_id: sectie,
            tip_contract: contract, norma_partiala_minute: minute, concediu_zile_an: zile, activ}));
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
                    <label className="latime-completa">{t('ui.sectie')}<select value={sectie} onChange={e => setSectie(e.target.value)}>
                        <option value="">{t('ui.fara_sectie')}</option>{date.sectii.map(s => <option value={s.id} key={s.id}>{denumireSectie(s, i18n.language)}</option>)}
                    </select></label>
                    <label>{t('ui.contract')} *<select value={contract} onChange={e => setContract(e.target.value)}>
                        <option value="INTREAGA">{t('ui.intreaga')}</option><option value="PARTIALA">{t('ui.partiala')}</option></select></label>
                    {contract === 'PARTIALA' ? <label>{t('ui.norma_partiala')}<input type="number" min="0" max="24" step="any" required value={ore} onChange={e => setOre(e.target.value)}/></label>
                        : <div className="norma-preluata"><span>{t('ui.norma')}</span><strong>{durata(date.norma_intreaga_minute, t)}</strong><small>{t('ui.norma_globala_indiciu')}</small></div>}
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

export function Setari({date, ocupat, salveazaNorma, adaugaSectie}: {date: datebaza.Date; ocupat: boolean;
    salveazaNorma: (minute: number) => Promise<boolean>; adaugaSectie: (denumire: string) => Promise<boolean>;
}) {
    const {t, i18n} = useTranslation();
    const [ore, setOre] = useState(String(date.norma_intreaga_minute / 60));
    const [sectie, setSectie] = useState('');
    const [validare, setValidare] = useState('');
    useEffect(() => { setOre(String(date.norma_intreaga_minute / 60)); }, [date.norma_intreaga_minute]);

    async function trimiteNorma(e: FormEvent) {
        e.preventDefault();
        const minute = minuteDinOre(ore);
        if (minute === null) { setValidare(t('date.norma_invalida')); return; }
        setValidare('');
        await salveazaNorma(minute);
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
    </div>;
}
