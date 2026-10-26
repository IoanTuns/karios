import {useEffect, useRef, useState, type FormEvent} from 'react';
import {useTranslation} from 'react-i18next';
import {CitesteCalendar, CompleteazaSarbatori2026, CompleteazaSarbatoriAn, SalveazaSarbatoare} from '../wailsjs/go/main/App';
import {datebaza, main} from '../wailsjs/go/models';
import {Dialog} from './Formulare';

function denumire(s: datebaza.Sarbatoare, limba: string) {
    return s.denumire[limba] || s.denumire.ro || Object.values(s.denumire).find(Boolean) || '';
}

export function Sarbatori() {
    const {t, i18n} = useTranslation();
    const [an, setAn] = useState(2026);
    const [anText, setAnText] = useState('2026');
    const [calendar, setCalendar] = useState<datebaza.CalendarAn | null>(null);
    const [incarcare, setIncarcare] = useState(true);
    const [ocupat, setOcupat] = useState(false);
    const blocat = useRef(false);
    const [eroare, setEroare] = useState('');
    const [mesaj, setMesaj] = useState('');
    const [editare, setEditare] = useState<datebaza.Sarbatoare | null>(null);
    const [importa, setImporta] = useState(false);

    function aplica(r: main.RaspunsCalendar) {
        if (r.eroare) { setEroare(t(r.eroare.cod, r.eroare.param ?? {})); return false; }
        if (!r.calendar) { setEroare(t('date.eroare')); return false; }
        setCalendar(r.calendar);
        return true;
    }

    useEffect(() => {
        let activ = true;
        setIncarcare(true); setCalendar(null); setEroare(''); setMesaj('');
        CitesteCalendar(an).then(r => { if (activ) aplica(r); })
            .catch(() => { if (activ) setEroare(t('date.eroare')); })
            .finally(() => { if (activ) setIncarcare(false); });
        return () => { activ = false; };
    }, [an]);

    async function executa(actiune: () => Promise<main.RaspunsCalendar>, importare = false) {
        if (blocat.current) return;
        blocat.current = true; setOcupat(true); setEroare(''); setMesaj('');
        try {
            const r = await actiune();
            if (aplica(r)) {
                setEditare(null); setImporta(false);
                setMesaj(t(importare ? 'sarbatori.completate' : 'ui.salvat', {numar: r.adaugate}));
            }
        } catch { setEroare(t('date.eroare')); }
        finally { setOcupat(false); blocat.current = false; }
    }

    function schimbaAn(e: FormEvent) {
        e.preventDefault();
        const valoare = Number(anText);
        if (!Number.isInteger(valoare) || valoare < 1 || valoare > 9999) { setEroare(t('sarbatori.an_invalid')); return; }
        setAn(valoare);
    }

    const zile = calendar?.sarbatori ?? [];
    const nelucratoare = zile.filter(s => [0, 6].includes(new Date(`${s.data}T12:00:00Z`).getUTCDay())).length;
    return <section className="calendar-sarbatori">
        <div className="unelte-calendar">
            <form onSubmit={schimbaAn}><label>{t('sarbatori.an')}<input aria-label={t('sarbatori.an')} type="number" min="1" max="9999" required value={anText} disabled={ocupat} onChange={e => setAnText(e.target.value)}/></label>
                <button className="buton secundar" disabled={ocupat}>{t('sarbatori.arata')}</button></form>
            <div className="actiuni-calendar">
                {calendar?.import_disponibil && <button className="buton secundar" disabled={ocupat || incarcare || !calendar} onClick={() => { setEroare(''); setImporta(true); }}>{t('sarbatori.completeaza')}</button>}
                <button className="buton primar" disabled={ocupat || incarcare || !calendar} onClick={() => { setEroare(''); setEditare(new datebaza.Sarbatoare({data: `${String(an).padStart(4, '0')}-01-01`, denumire: {ro: '', en: ''}})); }}>{t('sarbatori.adauga')}</button>
            </div>
        </div>
        {!editare && !importa && eroare && <div role="alert" className="mesaj eroare">{eroare}</div>}
        {mesaj && <div className="mesaj succes" role="status">{mesaj}</div>}
        {incarcare ? <p role="status">{t('ui.incarcare')}</p> : !calendar ? <button className="buton secundar" disabled={ocupat} onClick={() => executa(() => CitesteCalendar(an))}>{t('ui.reincearca')}</button> : <>
            <div className="indicatori">
                <div className="indicator"><span>{t('sarbatori.date_distincte')}</span><strong>{zile.length}</strong></div>
                <div className="indicator"><span>{t('sarbatori.luni_vineri')}</span><strong>{zile.length - nelucratoare}</strong></div>
                <div className="indicator"><span>{t('sarbatori.weekend')}</span><strong>{nelucratoare}</strong></div>
            </div>
            {calendar?.import_disponibil && <p className="nota-calendar">{t('sarbatori.precizare')}</p>}
            <div className="panou tabel-derulant"><table className="tabel-sarbatori"><thead><tr><th>{t('sarbatori.data')}</th><th>{t('sarbatori.zi')}</th><th>{t('sarbatori.denumire')}</th><th>{t('ui.actiuni')}</th></tr></thead>
                <tbody>{zile.map(s => <tr key={s.data}><td>{new Intl.DateTimeFormat(i18n.language).format(new Date(`${s.data}T12:00:00Z`))}</td>
                    <td>{new Intl.DateTimeFormat(i18n.language, {weekday: 'long', timeZone: 'UTC'}).format(new Date(`${s.data}T12:00:00Z`))}</td>
                    <td>{denumire(s, i18n.language)}</td><td><button className="legatura" disabled={ocupat} aria-label={t('sarbatori.editeaza_data', {data: s.data})} onClick={() => { setEroare(''); setEditare(s); }}>{t('ui.editeaza')}</button></td></tr>)}</tbody>
            </table>{zile.length === 0 && <p className="fara-rezultate">{t('sarbatori.gol')}</p>}</div>
        </>}
        {editare && <FormularSarbatoare sarbatoare={editare} exista={zile.some(s => s === editare)} ocupat={ocupat} eroare={eroare}
            inchide={() => { setEditare(null); setEroare(''); }}
            salveaza={s => executa(() => SalveazaSarbatoare(zile.some(z => z === editare) ? editare.data : '', s, an))}/>}
        {importa && <Dialog titlu={t('sarbatori.completeaza')} ocupat={ocupat} inchide={() => { setImporta(false); setEroare(''); }}>
            <div className="corp-dialog"><p>{t('sarbatori.import_indiciu')}</p><p className="text-discret">{t('sarbatori.precizare')}</p>{eroare && <p role="alert" className="mesaj eroare">{eroare}</p>}</div>
            <footer className="subsol-dialog"><button className="buton secundar" disabled={ocupat} onClick={() => setImporta(false)}>{t('ui.anuleaza')}</button>
                <button className="buton primar" disabled={ocupat} onClick={() => executa(() => CompleteazaSarbatoriAn(an), true)}>{t(ocupat ? 'ui.salvare' : 'sarbatori.confirma_import')}</button></footer>
        </Dialog>}
    </section>;
}

function FormularSarbatoare({sarbatoare, exista, ocupat, eroare, inchide, salveaza}: {sarbatoare: datebaza.Sarbatoare; exista: boolean; ocupat: boolean; eroare: string; inchide: () => void; salveaza: (s: datebaza.Sarbatoare) => Promise<void>}) {
    const {t} = useTranslation();
    const [data, setData] = useState(sarbatoare.data);
    const [ro, setRo] = useState(sarbatoare.denumire.ro || '');
    const [en, setEn] = useState(sarbatoare.denumire.en || '');
    const [validare, setValidare] = useState('');
    function inchideCuVerificare() {
        const modificat = data !== sarbatoare.data || ro !== (sarbatoare.denumire.ro || '') || en !== (sarbatoare.denumire.en || '');
        if (!modificat || window.confirm(t('ui.renunta'))) inchide();
    }
    return <Dialog titlu={t(exista ? 'sarbatori.editeaza' : 'sarbatori.adauga')} ocupat={ocupat} inchide={inchideCuVerificare}>
        <form onSubmit={async e => { e.preventDefault(); setValidare(''); if (!ro.trim() && !en.trim()) { setValidare(t('sarbatori.denumire_obligatorie')); return; }
            await salveaza(new datebaza.Sarbatoare({data, denumire: {...sarbatoare.denumire, ro, en}})); }}>
            <div className="corp-dialog">{(validare || eroare) && <p className="mesaj eroare" role="alert">{validare || eroare}</p>}
                <fieldset disabled={ocupat} className="campuri"><label className="latime-completa">{t('sarbatori.data')} *<input autoFocus required type="date" min="0001-01-01" max="9999-12-31" value={data} onChange={e => setData(e.target.value)}/></label>
                    <label className="latime-completa">{t('sarbatori.denumire_ro')}<input value={ro} onChange={e => setRo(e.target.value)}/></label>
                    <label className="latime-completa">{t('sarbatori.denumire_en')}<input value={en} onChange={e => setEn(e.target.value)}/></label></fieldset>
            </div><footer className="subsol-dialog"><button type="button" className="buton secundar" disabled={ocupat} onClick={inchideCuVerificare}>{t('ui.anuleaza')}</button><button className="buton primar" disabled={ocupat}>{t(ocupat ? 'ui.salvare' : 'ui.salveaza')}</button></footer>
        </form>
    </Dialog>;
}
