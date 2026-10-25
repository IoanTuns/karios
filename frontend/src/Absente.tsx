import {useEffect, useRef, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {CitesteAbsente, CitesteCalendar, SalveazaAbsente, SalveazaTipAbsenta} from '../wailsjs/go/main/App';
import {datebaza, main} from '../wailsjs/go/models';
import {Dialog} from './Formulare';

const categorii = ['ODIHNA', 'MEDICAL', 'SPECIAL'];
const cheieZi = (z: datebaza.ZiAbsenta) => `${z.angajat_id}/${z.data}`;
const numeTip = (tip: datebaza.TipAbsenta, limba: string) => tip.denumire[limba] || tip.denumire.ro || Object.values(tip.denumire).find(Boolean) || tip.cod;

function interval(deLa: string, panaLa: string): string[] {
    if (!/^\d{4}-\d{2}-\d{2}$/.test(deLa) || !/^\d{4}-\d{2}-\d{2}$/.test(panaLa) || deLa < '0001-01-01' || panaLa > '9999-12-31') return [];
    const inceput = Date.parse(`${deLa}T12:00:00Z`), sfarsit = Date.parse(`${panaLa}T12:00:00Z`);
    const numar = (sfarsit - inceput) / 86400000 + 1;
    if (!Number.isInteger(numar) || numar < 1 || numar > 366) return [];
    const zile = Array.from({length: numar}, (_, i) => new Date(inceput + i * 86400000).toISOString().slice(0, 10));
    return zile[0] === deLa && zile[zile.length - 1] === panaLa ? zile : [];
}

export function Absente({date}: {date: datebaza.Date}) {
    const {t, i18n} = useTranslation();
    const [an, setAn] = useState(new Date().getFullYear());
    const [anText, setAnText] = useState(String(an));
    const [dateAbsente, setDateAbsente] = useState<datebaza.AbsenteAn | null>(null);
    const [angajat, setAngajat] = useState('');
    const [categorie, setCategorie] = useState('');
    const [incarcare, setIncarcare] = useState(true);
    const [ocupat, setOcupat] = useState(false);
    const blocat = useRef(false);
    const [eroare, setEroare] = useState('');
    const [mesaj, setMesaj] = useState('');
    const [reincarcare, setReincarcare] = useState(0);
    const [tipEditat, setTipEditat] = useState<datebaza.TipAbsenta | null>(null);
    const [ziEditata, setZiEditata] = useState<datebaza.ZiAbsenta | null | undefined>(undefined);
    const [selectate, setSelectate] = useState<Set<string>>(new Set());

    function aplica(r: main.RaspunsAbsente) {
        if (r.eroare) { setEroare(t(r.eroare.cod, r.eroare.param ?? {})); return false; }
        if (!r.absente) { setEroare(t('date.eroare')); return false; }
        setDateAbsente(r.absente); return true;
    }
    useEffect(() => {
        let activ = true;
        setIncarcare(true); setDateAbsente(null); setEroare(''); setMesaj(''); setSelectate(new Set());
        CitesteAbsente(an).then(r => { if (activ) aplica(r); }).catch(() => { if (activ) setEroare(t('date.eroare')); }).finally(() => { if (activ) setIncarcare(false); });
        return () => { activ = false; };
    }, [an, reincarcare]);

    async function executa(actiune: () => Promise<main.RaspunsAbsente>) {
        if (blocat.current) return;
        blocat.current = true; setOcupat(true); setEroare(''); setMesaj('');
        try { if (aplica(await actiune())) { setTipEditat(null); setZiEditata(undefined); setSelectate(new Set()); setMesaj(t('ui.salvat')); } }
        catch { setEroare(t('date.eroare')); }
        finally { blocat.current = false; setOcupat(false); }
    }
    const tipuri = dateAbsente?.tipuri ?? [];
    const zile = (dateAbsente?.zile ?? []).filter(z => z.activ);
    const tipPentru = (cod: string) => tipuri.find(tip => tip.cod === cod);
    const persoane = date.angajati.filter(a => !angajat || a.id === angajat);
    const vizibile = zile.filter(z => (!angajat || z.angajat_id === angajat) && (!categorie || tipPentru(z.cod)?.categorie === categorie));
    const alese = vizibile.filter(z => selectate.has(cheieZi(z)));
    function schimbaFiltru(actiune: () => void) { actiune(); setSelectate(new Set()); }
    function anuleaza(lista: datebaza.ZiAbsenta[]) {
        if (!lista.length || !window.confirm(t('absente.confirma_anulare', {numar: lista.length}))) return;
        executa(() => SalveazaAbsente(lista.map(z => new datebaza.ZiAbsenta({...z, activ: false})), an));
    }

    return <section className="modul-absente">
        <div className="unelte-calendar unelte-perioada"><form onSubmit={e => { e.preventDefault(); const nou = Number(anText); if (Number.isInteger(nou) && nou >= 1 && nou <= 9999) setAn(nou); else setEroare(t('sarbatori.an_invalid')); }}>
            <label>{t('sarbatori.an')}<input aria-label={t('sarbatori.an')} type="number" min="1" max="9999" required disabled={ocupat} value={anText} onChange={e => setAnText(e.target.value)}/></label><button className="buton secundar" disabled={ocupat}>{t('sarbatori.arata')}</button></form>
            <div className="actiuni-calendar"><button className="buton secundar" disabled={ocupat || incarcare || !dateAbsente} onClick={() => { setEroare(''); setTipEditat(new datebaza.TipAbsenta({cod:'', denumire:{ro:'',en:''}, categorie:'SPECIAL', activ:true})); }}>{t('absente.adauga_tip')}</button>
                <button className="buton primar" disabled={ocupat || incarcare || !dateAbsente || !date.angajati.length || !tipuri.some(tip => tip.activ)} onClick={() => { setEroare(''); setZiEditata(null); }}>{t('absente.adauga')}</button></div></div>
        {!tipEditat && ziEditata === undefined && eroare && <p className="mesaj eroare" role="alert">{eroare}</p>}{mesaj && <p className="mesaj succes" role="status">{mesaj}</p>}
        {incarcare ? <p role="status">{t('ui.incarcare')}</p> : !dateAbsente ? <button className="buton secundar" onClick={() => setReincarcare(x => x + 1)}>{t('ui.reincearca')}</button> : <>
            <div className="unelte filtre-absente"><label>{t('ui.angajat')}<select aria-label={t('ui.angajat')} value={angajat} disabled={ocupat} onChange={e => schimbaFiltru(() => setAngajat(e.target.value))}><option value="">{t('absente.toti_angajatii')}</option>{date.angajati.map(a => <option key={a.id} value={a.id}>{a.nume} {a.prenume} · {a.marca}</option>)}</select></label>
                <label>{t('absente.categorie')}<select aria-label={t('absente.categorie')} disabled={ocupat} value={categorie} onChange={e => schimbaFiltru(() => setCategorie(e.target.value))}><option value="">{t('absente.toate_categoriile')}</option>{categorii.map(c => <option key={c} value={c}>{t(`absente.${c}`)}</option>)}</select></label>
                <button className="legatura" disabled={ocupat} onClick={() => setReincarcare(x => x + 1)}>{t('absente.reincarca')}</button></div>
            <h2>{t('absente.sumar', {an})}</h2><p className="nota-calendar">{t('absente.sumar_indiciu')}</p>
            <div className="panou tabel-derulant"><table className="sumar-absente"><thead><tr><th>{t('ui.angajat')}</th><th>{t('absente.drept_contractual')}</th>{categorii.map(c => <th key={c}>{t(`absente.${c}`)}</th>)}</tr></thead><tbody>
                {persoane.map(a => <tr key={a.id}><td>{a.nume} {a.prenume}<small className="text-discret"> · {a.marca}</small></td><td>{a.concediu_zile_an == null ? t('ui.necompletat') : t('ui.zile', {numar:a.concediu_zile_an})}</td>{categorii.map(c => <td key={c}>{zile.filter(z => z.angajat_id === a.id && tipPentru(z.cod)?.categorie === c).length}</td>)}</tr>)}
            </tbody></table></div>
            <h2>{t('absente.evidenta')}</h2><div className="unelte-calendar"><p className="nota-calendar">{t('absente.zile_inregistrate', {numar:vizibile.length})}</p><button className="buton secundar" disabled={ocupat || !alese.length} onClick={() => anuleaza(alese)}>{t('absente.anuleaza_selectate', {numar:alese.length})}</button></div>
            <div className="panou tabel-derulant"><table className="tabel-absente"><thead><tr><th><input type="checkbox" aria-label={t('absente.selecteaza_toate')} disabled={ocupat || !vizibile.length} checked={vizibile.length > 0 && alese.length === vizibile.length} onChange={e => setSelectate(new Set(e.target.checked ? vizibile.map(cheieZi) : []))}/></th>{['ui.angajat','absente.data','absente.tip','absente.categorie','ui.actiuni'].map(c => <th key={c}>{t(c)}</th>)}</tr></thead><tbody>
                {vizibile.map(z => { const a = date.angajati.find(a => a.id === z.angajat_id), tip = tipPentru(z.cod), k = cheieZi(z); return <tr key={k}>
                    <td><input type="checkbox" aria-label={t('absente.selecteaza_zi', {data:z.data, nume:`${a?.nume ?? ''} ${a?.prenume ?? ''}`})} disabled={ocupat} checked={selectate.has(k)} onChange={e => setSelectate(vechi => { const noi = new Set(vechi); if (e.target.checked) noi.add(k); else noi.delete(k); return noi; })}/></td>
                    <td>{a?.nume} {a?.prenume}<small className="text-discret"> · {a?.marca}</small></td><td>{z.data}</td><td><strong>{z.cod}</strong> · {tip ? numeTip(tip, i18n.language) : z.cod}</td><td>{tip && t(`absente.${tip.categorie}`)}</td><td><div className="actiuni-rand"><button className="legatura" disabled={ocupat} onClick={() => { setEroare(''); setZiEditata(z); }}>{t('ui.editeaza')}</button><button className="legatura discreta" disabled={ocupat} onClick={() => anuleaza([z])}>{t('absente.anuleaza')}</button></div></td></tr>; })}
            </tbody></table>{!vizibile.length && <p className="fara-rezultate">{t('absente.gol')}</p>}</div>
            <h2>{t('absente.tipuri')}</h2><p className="nota-calendar">{t('absente.tipuri_indiciu')}</p><div className="panou tabel-derulant"><table className="tabel-tipuri-absenta"><thead><tr>{['schimburi.cod','schimburi.denumire','absente.categorie','ui.stare','ui.actiuni'].map(c => <th key={c}>{t(c)}</th>)}</tr></thead><tbody>{tipuri.map(tip => <tr key={tip.cod}><td><strong>{tip.cod}</strong></td><td>{numeTip(tip, i18n.language)}</td><td>{t(`absente.${tip.categorie}`)}</td><td>{t(tip.activ ? 'ui.activ' : 'ui.inactiv')}</td><td><button className="legatura" disabled={ocupat} onClick={() => { setEroare(''); setTipEditat(tip); }}>{t('ui.editeaza')}</button></td></tr>)}</tbody></table></div>
        </>}
        {tipEditat && <FormularTip key={tipEditat.cod || 'nou'} tip={tipEditat} ocupat={ocupat} eroare={eroare} inchide={() => { setTipEditat(null); setEroare(''); }} salveaza={tip => executa(() => SalveazaTipAbsenta(tipEditat.cod, tip, an))}/>}
        {ziEditata !== undefined && dateAbsente && <FormularAbsente key={ziEditata ? cheieZi(ziEditata) : 'nou'} zi={ziEditata} date={date} absente={dateAbsente} angajatImplicit={angajat} ocupat={ocupat} eroare={eroare} inchide={() => { setZiEditata(undefined); setEroare(''); }} salveaza={zile => executa(() => SalveazaAbsente(zile, an))}/>}
    </section>;
}

function FormularTip({tip, ocupat, eroare, inchide, salveaza}: {tip: datebaza.TipAbsenta; ocupat: boolean; eroare: string; inchide: () => void; salveaza: (tip: datebaza.TipAbsenta) => Promise<void>}) {
    const {t} = useTranslation();
    const initial = {cod:tip.cod, ro:tip.denumire.ro || '', en:tip.denumire.en || '', categorie:tip.categorie, activ:tip.activ};
    const [campuri, setCampuri] = useState(initial);
    const [validare, setValidare] = useState('');
    function seteaza<K extends keyof typeof campuri>(k: K, v: typeof campuri[K]) { setCampuri(c => ({...c,[k]:v})); }
    function renunta() { if (JSON.stringify(initial) === JSON.stringify(campuri) || window.confirm(t('ui.renunta'))) inchide(); }
    return <Dialog titlu={t(tip.cod ? 'absente.editeaza_tip' : 'absente.adauga_tip')} ocupat={ocupat} inchide={renunta}><form onSubmit={e => { e.preventDefault(); if (!campuri.ro.trim() && !campuri.en.trim()) { setValidare(t('sarbatori.denumire_obligatorie')); return; } setValidare(''); salveaza(new datebaza.TipAbsenta({cod:campuri.cod, denumire:{...tip.denumire,ro:campuri.ro,en:campuri.en}, categorie:campuri.categorie, activ:campuri.activ})); }}>
        <div className="corp-dialog">{(eroare || validare) && <p className="mesaj eroare" role="alert">{validare || eroare}</p>}<fieldset className="campuri" disabled={ocupat}>
            <label>{t('schimburi.cod')} *<input name="cod" required maxLength={16} readOnly={!!tip.cod} value={campuri.cod} onChange={e => seteaza('cod',e.target.value)}/></label>
            <label>{t('absente.categorie')}<select name="categorie" value={campuri.categorie} onChange={e => seteaza('categorie',e.target.value)}>{categorii.map(c => <option key={c} value={c}>{t(`absente.${c}`)}</option>)}</select></label>
            <label className="latime-completa">{t('sarbatori.denumire_ro')}<input name="ro" value={campuri.ro} onChange={e => seteaza('ro',e.target.value)}/></label><label className="latime-completa">{t('sarbatori.denumire_en')}<input name="en" value={campuri.en} onChange={e => seteaza('en',e.target.value)}/></label>
            <label className="bifa"><input name="activ" type="checkbox" checked={campuri.activ} onChange={e => seteaza('activ',e.target.checked)}/>{t('ui.activ_indiciu')}</label>
        </fieldset></div><footer className="subsol-dialog"><button type="button" className="buton secundar" disabled={ocupat} onClick={renunta}>{t('ui.anuleaza')}</button><button className="buton primar" disabled={ocupat}>{t('ui.salveaza')}</button></footer>
    </form></Dialog>;
}

function FormularAbsente({zi, date, absente, angajatImplicit, ocupat, eroare, inchide, salveaza}: {zi: datebaza.ZiAbsenta | null; date: datebaza.Date; absente: datebaza.AbsenteAn; angajatImplicit: string; ocupat: boolean; eroare: string; inchide: () => void; salveaza: (zile: datebaza.ZiAbsenta[]) => Promise<void>}) {
    const {t, i18n} = useTranslation();
    const azi = new Date(), dataImplicita = absente.an === azi.getFullYear() ? `${absente.an}-${String(azi.getMonth()+1).padStart(2,'0')}-${String(azi.getDate()).padStart(2,'0')}` : `${String(absente.an).padStart(4,'0')}-01-01`;
    const initial = {angajat:zi?.angajat_id || angajatImplicit, cod:zi?.cod || absente.tipuri.find(tip => tip.activ && tip.categorie === 'ODIHNA')?.cod || absente.tipuri.find(tip => tip.activ)?.cod || '', deLa:zi?.data || dataImplicita, panaLa:zi?.data || dataImplicita};
    const [campuri, setCampuri] = useState(initial);
    const [excluse, setExcluse] = useState<Set<string>>(new Set());
    const [sarbatori, setSarbatori] = useState<Map<string,string>>(new Map());
    const [calendarGata, setCalendarGata] = useState(false);
    const [calendarEroare, setCalendarEroare] = useState('');
    const [validare, setValidare] = useState('');
    const zile = interval(campuri.deLa,campuri.panaLa);
    const alese = zile.filter(data => !excluse.has(data));
    useEffect(() => {
        let activ = true;
        setCalendarGata(false); setSarbatori(new Map()); setCalendarEroare('');
        const ani = [...new Set(zile.map(data => Number(data.slice(0,4))))];
        if (ani.length) Promise.all(ani.map(CitesteCalendar)).then(lista => {
            if (!activ) return;
            const err = lista.find(r => r.eroare)?.eroare;
            if (err || lista.some(r => !r.calendar)) { setCalendarEroare(t(err?.cod || 'date.eroare', err?.param ?? {})); return; }
            setSarbatori(new Map(lista.flatMap(r => r.calendar!.sarbatori.map(s => [s.data, s.denumire[i18n.language] || s.denumire.ro || s.data] as [string,string])))); setCalendarGata(true);
        }).catch(() => { if (activ) setCalendarEroare(t('date.eroare')); });
        return () => { activ = false; };
    }, [campuri.deLa,campuri.panaLa,i18n.language]);
    function seteaza<K extends keyof typeof campuri>(k: K, v: string) { setCampuri(c => ({...c,[k]:v})); if (k === 'deLa' || k === 'panaLa') setExcluse(new Set()); }
    function renunta() { if ((JSON.stringify(initial) === JSON.stringify(campuri) && !excluse.size) || window.confirm(t('ui.renunta'))) inchide(); }
    return <Dialog titlu={t(zi ? 'absente.editeaza' : 'absente.adauga')} ocupat={ocupat} inchide={renunta}><form onSubmit={e => {
        e.preventDefault(); setValidare('');
        if (!campuri.angajat || !alese.length) { setValidare(t('absente.selectie_invalida')); return; }
        const existente = absente.zile.filter(z => z.angajat_id === campuri.angajat);
        if (!zi && alese.some(data => existente.some(z => z.data === data && z.activ))) { setValidare(t('absente.zi_ocupata')); return; }
        salveaza(alese.map(data => new datebaza.ZiAbsenta({angajat_id:campuri.angajat, data, cod:campuri.cod, activ:true, revizie:existente.find(z => z.data === data)?.revizie || ''})));
    }}><div className="corp-dialog">{(validare || eroare) && <p className="mesaj eroare" role="alert">{validare || eroare}</p>}
        <fieldset className="campuri" disabled={ocupat}><label className="latime-completa">{t('ui.angajat')} *<select name="angajat" required disabled={!!zi} value={campuri.angajat} onChange={e => seteaza('angajat',e.target.value)}><option value="">{t('absente.alege_angajat')}</option>{date.angajati.map(a => <option key={a.id} value={a.id}>{a.nume} {a.prenume} · {a.marca}{!a.activ ? ` · ${t('ui.inactiv')}` : ''}</option>)}</select></label>
            <label className="latime-completa">{t('absente.tip')}<select name="tip" required value={campuri.cod} onChange={e => seteaza('cod',e.target.value)}>{absente.tipuri.filter(tip => tip.activ || tip.cod === zi?.cod).map(tip => <option key={tip.cod} value={tip.cod}>{tip.cod} · {numeTip(tip,i18n.language)}{!tip.activ ? ` · ${t('ui.inactiv')}` : ''}</option>)}</select></label>
            <label>{t('absente.de_la')} *<input name="de_la" type="date" min="0001-01-01" max="9999-12-31" required readOnly={!!zi} value={campuri.deLa} onChange={e => seteaza('deLa',e.target.value)}/></label><label>{t('absente.pana_la')} *<input name="pana_la" type="date" min="0001-01-01" max="9999-12-31" required readOnly={!!zi} value={campuri.panaLa} onChange={e => seteaza('panaLa',e.target.value)}/></label>
        </fieldset><p className="nota-calendar">{t('absente.selectie_indiciu')}</p>{!zile.length && <p className="mesaj eroare">{t('absente.interval_invalid')}</p>}{calendarEroare && <p className="mesaj eroare" role="alert">{calendarEroare}</p>}
        {!zi && <div className="actiuni-calendar"><button type="button" className="legatura" disabled={ocupat || !zile.length} onClick={() => setExcluse(new Set())}>{t('absente.toate_zilele')}</button><button type="button" className="legatura" disabled={ocupat || !calendarGata} onClick={() => setExcluse(new Set(zile.filter(data => [0,6].includes(new Date(`${data}T12:00:00Z`).getUTCDay()) || sarbatori.has(data))))}>{t('absente.luni_vineri')}</button></div>}
        <div className="zile-absenta" role="group" aria-label={t('absente.zile_selectate')}>
            {zile.map(data => { const d = new Date(`${data}T12:00:00Z`), weekend = [0,6].includes(d.getUTCDay()); return <label key={data} className={sarbatori.has(data) ? 'zi-sarbatoare' : weekend ? 'zi-weekend' : ''}><input type="checkbox" disabled={ocupat || !!zi} checked={!excluse.has(data)} onChange={e => setExcluse(vechi => { const noi = new Set(vechi); if (e.target.checked) noi.delete(data); else noi.add(data); return noi; })}/><span>{data} · {new Intl.DateTimeFormat(i18n.language,{weekday:'short',timeZone:'UTC'}).format(d)}{sarbatori.has(data) && <small>{sarbatori.get(data)}</small>}</span></label>; })}
        </div><p className="nota-calendar" role="status">{t('absente.selectate', {numar:alese.length})}</p>
    </div><footer className="subsol-dialog"><button type="button" className="buton secundar" disabled={ocupat} onClick={renunta}>{t('ui.anuleaza')}</button><button className="buton primar" disabled={ocupat || !alese.length}>{t('ui.salveaza')}</button></footer></form></Dialog>;
}
