import {useEffect, useRef, useState, type ChangeEvent} from 'react';
import {useTranslation} from 'react-i18next';
import {CitesteAbsente, CitesteCalendar, CitesteDocumentAbsenta, SalveazaAbsente, SalveazaTipAbsenta} from '../wailsjs/go/main/App';
import {datebaza, main} from '../wailsjs/go/models';
import {Dialog} from './Formulare';

const categorii = ['ODIHNA', 'MEDICAL', 'SPECIAL'];
const stariAprobare = ['IN_ASTEPTARE', 'APROBATA', 'RESPINSA'];
const cheieZi = (z: datebaza.ZiAbsenta) => `${z.angajat_id}/${z.data}`;
const numeTip = (tip: datebaza.TipAbsenta, limba: string) => tip.denumire[limba] || tip.denumire.ro || Object.values(tip.denumire).find(Boolean) || tip.cod;
const normalizat = (text: string) => text.normalize('NFD').replace(/\p{Diacritic}/gu, '').toLocaleLowerCase();

function interval(deLa: string, panaLa: string): string[] {
    if (!/^\d{4}-\d{2}-\d{2}$/.test(deLa) || !/^\d{4}-\d{2}-\d{2}$/.test(panaLa) || deLa < '0001-01-01' || panaLa > '9999-12-31') return [];
    const inceput = Date.parse(`${deLa}T12:00:00Z`), sfarsit = Date.parse(`${panaLa}T12:00:00Z`);
    const numar = (sfarsit - inceput) / 86400000 + 1;
    if (!Number.isInteger(numar) || numar < 1 || numar > 366) return [];
    const zile = Array.from({length: numar}, (_, i) => new Date(inceput + i * 86400000).toISOString().slice(0, 10));
    return zile[0] === deLa && zile[zile.length - 1] === panaLa ? zile : [];
}

export function Absente({date, filtruInitial}: {date: datebaza.Date; filtruInitial?: {an: number; angajatId: string} | null}) {
    const {t, i18n} = useTranslation();
    const [an, setAn] = useState(filtruInitial?.an ?? new Date().getFullYear());
    const [anText, setAnText] = useState(String(an));
    const [dateAbsente, setDateAbsente] = useState<datebaza.AbsenteAn | null>(null);
    const [angajat, setAngajat] = useState(filtruInitial?.angajatId ?? '');
    const [categorie, setCategorie] = useState('');
    const [filtruAprobare, setFiltruAprobare] = useState('');
    const [cautareAbsente, setCautareAbsente] = useState('');
    const [sectiuneActiva, setSectiuneActiva] = useState<'toate' | 'zile' | 'sumar' | 'tipuri'>('toate');
    const [incarcare, setIncarcare] = useState(true);
    const [ocupat, setOcupat] = useState(false);
    const blocat = useRef(false);
    const [eroare, setEroare] = useState('');
    const [mesaj, setMesaj] = useState('');
    const [reincarcare, setReincarcare] = useState(0);
    const [tipEditat, setTipEditat] = useState<datebaza.TipAbsenta | null>(null);
    const [ziEditata, setZiEditata] = useState<datebaza.ZiAbsenta | null | undefined>(undefined);
    const [lotEditat, setLotEditat] = useState<datebaza.ZiAbsenta[] | null>(null);
    const [docPrevizualizat, setDocPrevizualizat] = useState<datebaza.DocumentAbsenta | null>(null);
    const [selectate, setSelectate] = useState<Set<string>>(new Set());

    useEffect(() => {
        if (!filtruInitial) return;
        setAn(filtruInitial.an);
        setAnText(String(filtruInitial.an));
        setAngajat(filtruInitial.angajatId);
    }, [filtruInitial]);

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
        try { if (aplica(await actiune())) { setTipEditat(null); setZiEditata(undefined); setLotEditat(null); setSelectate(new Set()); setMesaj(t('ui.salvat')); } }
        catch { setEroare(t('date.eroare')); }
        finally { blocat.current = false; setOcupat(false); }
    }

    function blobDinDocument(doc: datebaza.DocumentAbsenta): Blob {
        const bin = atob(doc.continut_b64);
        const octeti = new Uint8Array(bin.length);
        for (let i = 0; i < bin.length; i++) octeti[i] = bin.charCodeAt(i);
        return new Blob([octeti], {type: doc.tip_mime});
    }

    function descarcaObiectDocument(doc: datebaza.DocumentAbsenta) {
        const url = URL.createObjectURL(blobDinDocument(doc));
        const a = document.createElement('a');
        a.href = url; a.download = doc.nume_fisier;
        document.body.appendChild(a); a.click(); a.remove();
        setTimeout(() => URL.revokeObjectURL(url), 1000);
    }

    async function incarcaDocument(id: string): Promise<datebaza.DocumentAbsenta | null> {
        if (blocat.current || !id) return null;
        blocat.current = true; setOcupat(true); setEroare('');
        try {
            const r = await CitesteDocumentAbsenta(id);
            if (r.eroare) { setEroare(t(r.eroare.cod, r.eroare.param ?? {})); return null; }
            if (!r.document) { setEroare(t('date.eroare')); return null; }
            return r.document;
        } catch { setEroare(t('date.eroare')); return null; }
        finally { blocat.current = false; setOcupat(false); }
    }

    async function descarcaDocument(id: string) {
        const doc = await incarcaDocument(id);
        if (doc) descarcaObiectDocument(doc);
    }

    async function previzualizeazaDocument(id: string) {
        const doc = await incarcaDocument(id);
        if (doc) setDocPrevizualizat(doc);
    }

    const tipuri = dateAbsente?.tipuri ?? [];
    const zile = (dateAbsente?.zile ?? []).filter(z => z.activ);
    const frecventaInlocuiri = new Map<string, number>();
    for (const z of zile) {
        if (z.inlocuitor_id && (z.stare_aprobare || 'IN_ASTEPTARE') !== 'RESPINSA') {
            const k = `${z.inlocuitor_id}/${z.data}`;
            frecventaInlocuiri.set(k, (frecventaInlocuiri.get(k) ?? 0) + 1);
        }
    }
    const tipPentru = (cod: string) => tipuri.find(tip => tip.cod === cod);
    const persoane = date.angajati.filter(a => !angajat || a.id === angajat);
    const textCautat = normalizat(cautareAbsente.trim());
    const vizibile = zile.filter(z => {
        if (angajat && z.angajat_id !== angajat) return false;
        const tip = tipPentru(z.cod);
        if (categorie && tip?.categorie !== categorie) return false;
        if (filtruAprobare && (z.stare_aprobare || 'IN_ASTEPTARE') !== filtruAprobare) return false;
        if (textCautat) {
            const a = date.angajati.find(x => x.id === z.angajat_id);
            const inl = z.inlocuitor_id ? date.angajati.find(x => x.id === z.inlocuitor_id) : undefined;
            const textRand = normalizat(`${a?.nume ?? ''} ${a?.prenume ?? ''} ${a?.marca ?? ''} ${z.data} ${z.cod} ${tip ? numeTip(tip, i18n.language) : ''} ${z.nr_inregistrare ?? ''} ${inl?.nume ?? ''} ${inl?.prenume ?? ''} ${inl?.marca ?? ''} ${z.aprobat_de ?? ''} ${z.document_nume ?? ''}`);
            if (!textRand.includes(textCautat)) return false;
        }
        return true;
    });
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
                <button className="buton primar" disabled={ocupat || incarcare || !dateAbsente || !date.angajati.length || !tipuri.some(tip => tip.activ)} onClick={() => { setEroare(''); setLotEditat(null); setZiEditata(null); }}>{t('absente.adauga')}</button></div></div>
        {!tipEditat && ziEditata === undefined && !lotEditat && eroare && <p className="mesaj eroare" role="alert">{eroare}</p>}{mesaj && <p className="mesaj succes" role="status">{mesaj}</p>}
        {incarcare ? <p role="status">{t('ui.incarcare')}</p> : !dateAbsente ? <button className="buton secundar" onClick={() => setReincarcare(x => x + 1)}>{t('ui.reincearca')}</button> : <>
            <div className="unelte filtre-absente"><label>{t('ui.angajat')}<select aria-label={t('ui.angajat')} value={angajat} disabled={ocupat} onChange={e => schimbaFiltru(() => setAngajat(e.target.value))}><option value="">{t('absente.toti_angajatii')}</option>{date.angajati.map(a => <option key={a.id} value={a.id}>{a.nume} {a.prenume} · {a.marca}</option>)}</select></label>
                <label>{t('absente.categorie')}<select aria-label={t('absente.categorie')} disabled={ocupat} value={categorie} onChange={e => schimbaFiltru(() => setCategorie(e.target.value))}><option value="">{t('absente.toate_categoriile')}</option>{categorii.map(c => <option key={c} value={c}>{t(`absente.${c}`)}</option>)}</select></label>
                <label>{t('absente.aprobare')}<select aria-label={t('absente.aprobare')} disabled={ocupat} value={filtruAprobare} onChange={e => schimbaFiltru(() => setFiltruAprobare(e.target.value))}><option value="">{t('absente.toate_starile')}</option>{stariAprobare.map(s => <option key={s} value={s}>{t(`absente.aprobare_${s}`)}</option>)}</select></label>
                <label>{t('absente.cauta')}<input name="cautare_absente" aria-label={t('absente.cauta')} placeholder={t('absente.cauta_indiciu')} disabled={ocupat} value={cautareAbsente} onChange={e => schimbaFiltru(() => setCautareAbsente(e.target.value))}/></label>
                <button className="legatura" disabled={ocupat} onClick={() => setReincarcare(x => x + 1)}>{t('absente.reincarca')}</button></div>
            <div className="taburi-sectiuni" role="group" aria-label={t('ui.absente')}>
                {(['toate', 'zile', 'sumar', 'tipuri'] as const).map(s => <button key={s} type="button" className={`buton-tab ${sectiuneActiva === s ? 'activ' : ''}`} aria-pressed={sectiuneActiva === s} onClick={() => setSectiuneActiva(s)}>{t(`absente.sectiune_${s}`)}</button>)}
            </div>
            {(sectiuneActiva === 'toate' || sectiuneActiva === 'sumar') && <>
            <h2>{t('absente.sumar', {an})}</h2><p className="nota-calendar">{t('absente.sumar_indiciu')}</p>
            <div className="panou tabel-derulant"><table className="sumar-absente"><thead><tr><th>{t('ui.angajat')}</th><th>{t('absente.drept_contractual')}</th>{categorii.map(c => <th key={c}>{t(`absente.${c}`)}</th>)}</tr></thead><tbody>
                {persoane.map(a => {
                    const odihnaNerespinsa = zile.filter(z => z.angajat_id === a.id && tipPentru(z.cod)?.categorie === 'ODIHNA' && (z.stare_aprobare || 'IN_ASTEPTARE') !== 'RESPINSA');
                    const inAsteptare = odihnaNerespinsa.filter(z => (z.stare_aprobare || 'IN_ASTEPTARE') === 'IN_ASTEPTARE').length;
                    const ramas = a.concediu_zile_an == null ? null : a.concediu_zile_an - odihnaNerespinsa.length;
                    return <tr key={a.id}><td>{a.nume} {a.prenume}<small className="text-discret"> · {a.marca}</small></td><td>{a.concediu_zile_an == null ? t('ui.necompletat') : <><div>{t('ui.zile', {numar:a.concediu_zile_an})}</div><small className={ramas !== null && ramas < 0 ? 'sold-concediu depasit' : 'sold-concediu'}>{t('absente.ramas_contractual', {ramas: t('ui.zile', {numar: ramas ?? 0})})}{inAsteptare > 0 ? ` · ${t('absente.in_asteptare_scurt', {numar: inAsteptare})}` : ''}</small></>}</td>{categorii.map(c => <td key={c}>{zile.filter(z => z.angajat_id === a.id && tipPentru(z.cod)?.categorie === c).length}</td>)}</tr>;
                })}
            </tbody></table></div>
            </>}
            {(sectiuneActiva === 'toate' || sectiuneActiva === 'zile') && <>
            <h2>{t('absente.evidenta')}</h2>
            <div className="unelte-calendar">
                <p className="nota-calendar">{t('absente.zile_inregistrate', {numar:vizibile.length})}</p>
                <div className="actiuni-calendar">
                    <button className="buton secundar" disabled={ocupat || !alese.length} onClick={() => { setEroare(''); setZiEditata(undefined); setLotEditat(alese); }}>{t('absente.editeaza_selectate', {numar:alese.length})}</button>
                    <button className="buton secundar" disabled={ocupat || !alese.length} onClick={() => anuleaza(alese)}>{t('absente.anuleaza_selectate', {numar:alese.length})}</button>
                </div>
            </div>
            <div className="panou tabel-derulant"><table className="tabel-absente"><thead><tr><th><input type="checkbox" aria-label={t('absente.selecteaza_toate')} disabled={ocupat || !vizibile.length} checked={vizibile.length > 0 && alese.length === vizibile.length} onChange={e => setSelectate(new Set(e.target.checked ? vizibile.map(cheieZi) : []))}/></th>{['ui.angajat','absente.data','absente.tip','absente.categorie','absente.nr_inregistrare','absente.inlocuitor','absente.aprobare','absente.cerere_atasata','ui.actiuni'].map(c => <th key={c}>{t(c)}</th>)}</tr></thead><tbody>
                {vizibile.map(z => {
                    const a = date.angajati.find(a => a.id === z.angajat_id);
                    const inlocuitor = z.inlocuitor_id ? date.angajati.find(x => x.id === z.inlocuitor_id) : undefined;
                    const tip = tipPentru(z.cod);
                    const k = cheieZi(z);
                    const stare = z.stare_aprobare || 'IN_ASTEPTARE';
                    const dublaInlocuire = Boolean(z.inlocuitor_id && stare !== 'RESPINSA' && (frecventaInlocuiri.get(`${z.inlocuitor_id}/${z.data}`) ?? 0) > 1);
                    const zileCerere = (z.document_id || z.nr_inregistrare)
                        ? vizibile.filter(v => z.document_id ? v.document_id === z.document_id : (v.nr_inregistrare === z.nr_inregistrare && v.angajat_id === z.angajat_id))
                        : [];
                    return <tr key={k}>
                        <td><input type="checkbox" aria-label={t('absente.selecteaza_zi', {data:z.data, nume:`${a?.nume ?? ''} ${a?.prenume ?? ''}`})} disabled={ocupat} checked={selectate.has(k)} onChange={e => setSelectate(vechi => { const noi = new Set(vechi); if (e.target.checked) noi.add(k); else noi.delete(k); return noi; })}/></td>
                        <td>{a?.nume} {a?.prenume}<small className="text-discret"> · {a?.marca}</small></td>
                        <td>{z.data}</td>
                        <td><strong>{z.cod}</strong> · {tip ? numeTip(tip, i18n.language) : z.cod}</td>
                        <td>{tip && t(`absente.${tip.categorie}`)}</td>
                        <td>{z.nr_inregistrare || '—'}</td>
                        <td>{inlocuitor ? <>{inlocuitor.nume} {inlocuitor.prenume}<small className="text-discret"> · {inlocuitor.marca}</small>{dublaInlocuire && <small className="zi-urmatoare avertisment-text">{t('absente.dubla_inlocuire_eticheta')}</small>}</> : '—'}</td>
                        <td><span className={`eticheta-aprobare aprobare-${stare}`}>{t(`absente.aprobare_${stare}`)}</span>{z.aprobat_de && <small className="zi-urmatoare">{z.aprobat_de}</small>}</td>
                        <td>{z.document_id ? <div className="actiuni-rand"><a href="#previzualizare" className="legatura" aria-label={t('absente.previzualizeaza_cerere', {nume: z.document_nume || z.data})} onClick={e => { e.preventDefault(); previzualizeazaDocument(z.document_id); }}>{z.document_nume}</a><a href="#cerere" className="legatura discreta" aria-label={t('absente.descarca_cerere', {nume: z.document_nume || z.data})} onClick={e => { e.preventDefault(); descarcaDocument(z.document_id); }}>{t('absente.descarca_scurt')}</a></div> : <span className="text-discret">{t('absente.fara_atasament')}</span>}</td>
                        <td><div className="actiuni-rand"><button className="legatura" disabled={ocupat} onClick={() => { setEroare(''); setLotEditat(null); setZiEditata(z); }}>{t('ui.editeaza')}</button><button className="legatura discreta" disabled={ocupat} onClick={() => anuleaza([z])}>{t('absente.anuleaza')}</button>{zileCerere.length > 1 && <button type="button" className="legatura discreta buton-selecteaza-cererea" disabled={ocupat} onClick={() => setSelectate(new Set(zileCerere.map(cheieZi)))}>{t('absente.selecteaza_cererea')}</button>}</div></td>
                    </tr>;
                })}
            </tbody></table>{!vizibile.length && <p className="fara-rezultate">{t('absente.gol')}</p>}</div>
            </>}
            {(sectiuneActiva === 'toate' || sectiuneActiva === 'tipuri') && <>
            <h2>{t('absente.tipuri')}</h2><p className="nota-calendar">{t('absente.tipuri_indiciu')}</p><div className="panou tabel-derulant"><table className="tabel-tipuri-absenta"><thead><tr>{['schimburi.cod','schimburi.denumire','absente.categorie','ui.stare','ui.actiuni'].map(c => <th key={c}>{t(c)}</th>)}</tr></thead><tbody>{tipuri.map(tip => <tr key={tip.cod}><td><strong>{tip.cod}</strong></td><td>{numeTip(tip, i18n.language)}</td><td>{t(`absente.${tip.categorie}`)}</td><td>{t(tip.activ ? 'ui.activ' : 'ui.inactiv')}</td><td><button className="legatura" disabled={ocupat} onClick={() => { setEroare(''); setTipEditat(tip); }}>{t('ui.editeaza')}</button></td></tr>)}</tbody></table></div>
            </>}
        </>}
        {tipEditat && <FormularTip key={tipEditat.cod || 'nou'} tip={tipEditat} ocupat={ocupat} eroare={eroare} inchide={() => { setTipEditat(null); setEroare(''); }} salveaza={tip => executa(() => SalveazaTipAbsenta(tipEditat.cod, tip, an))}/>}
        {(ziEditata !== undefined || lotEditat) && dateAbsente && <FormularAbsente key={lotEditat ? `lot-${lotEditat.map(cheieZi).join(',')}` : ziEditata ? cheieZi(ziEditata) : 'nou'} zi={ziEditata ?? null} lot={lotEditat} date={date} absente={dateAbsente} angajatImplicit={angajat} ocupat={ocupat} eroare={eroare} inchide={() => { setZiEditata(undefined); setLotEditat(null); setEroare(''); }} salveaza={zile => executa(() => SalveazaAbsente(zile, an))}/>}
        {docPrevizualizat && <PrevizualizareCerere doc={docPrevizualizat} inchide={() => setDocPrevizualizat(null)} descarca={() => descarcaObiectDocument(docPrevizualizat)}/>}
    </section>;
}

function PrevizualizareCerere({doc, inchide, descarca}: {doc: datebaza.DocumentAbsenta; inchide: () => void; descarca: () => void}) {
    const {t} = useTranslation();
    const [url, setUrl] = useState('');
    useEffect(() => {
        const bin = atob(doc.continut_b64);
        const octeti = new Uint8Array(bin.length);
        for (let i = 0; i < bin.length; i++) octeti[i] = bin.charCodeAt(i);
        const u = URL.createObjectURL(new Blob([octeti], {type: doc.tip_mime}));
        setUrl(u);
        return () => URL.revokeObjectURL(u);
    }, [doc]);
    return <Dialog titlu={t('absente.previzualizare_titlu', {nume: doc.nume_fisier})} ocupat={false} inchide={inchide}>
        <div className="corp-dialog previzualizare-cerere-corp">
            {url && (doc.tip_mime === 'application/pdf' ? <iframe className="cadru-previzualizare-cerere" src={url} title={doc.nume_fisier}/> : <img className="imagine-previzualizare-cerere" src={url} alt={doc.nume_fisier}/>)}
        </div>
        <footer className="subsol-dialog">
            <button type="button" className="buton secundar" onClick={descarca}>{t('absente.descarca_scurt')}</button>
            <button type="button" className="buton primar" onClick={inchide}>{t('absente.inchide_previzualizare')}</button>
        </footer>
    </Dialog>;
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

function FormularAbsente({zi, lot, date, absente, angajatImplicit, ocupat, eroare, inchide, salveaza}: {zi: datebaza.ZiAbsenta | null; lot: datebaza.ZiAbsenta[] | null; date: datebaza.Date; absente: datebaza.AbsenteAn; angajatImplicit: string; ocupat: boolean; eroare: string; inchide: () => void; salveaza: (zile: datebaza.ZiAbsenta[]) => Promise<void>}) {
    const {t, i18n} = useTranslation();
    const sursa = zi ?? (lot?.[0] ?? null);
    const modEditare = Boolean(zi || lot);
    const acelasiAngajatLot = !lot || new Set(lot.map(z => z.angajat_id)).size === 1;
    const azi = new Date(), dataImplicita = absente.an === azi.getFullYear() ? `${absente.an}-${String(azi.getMonth()+1).padStart(2,'0')}-${String(azi.getDate()).padStart(2,'0')}` : `${String(absente.an).padStart(4,'0')}-01-01`;
    const initial = {
        angajat: sursa?.angajat_id || angajatImplicit,
        cod: sursa?.cod || absente.tipuri.find(tip => tip.activ && tip.categorie === 'ODIHNA')?.cod || absente.tipuri.find(tip => tip.activ)?.cod || '',
        deLa: zi?.data || (lot ? lot[0].data : dataImplicita),
        panaLa: zi?.data || (lot ? lot[lot.length - 1].data : dataImplicita),
        nrInregistrare: sursa?.nr_inregistrare || '',
        inlocuitor: sursa?.inlocuitor_id || '',
        stareAprobare: sursa?.stare_aprobare || 'IN_ASTEPTARE',
        aprobatDe: sursa?.aprobat_de || '',
        documentId: sursa?.document_id || '',
        documentNume: sursa?.document_nume || '',
        documentMime: sursa?.document_mime || '',
        documentContinutB64: ''
    };
    const [campuri, setCampuri] = useState(initial);
    const [excluse, setExcluse] = useState<Set<string>>(new Set());
    const [sarbatori, setSarbatori] = useState<Map<string,string>>(new Map());
    const [calendarGata, setCalendarGata] = useState(false);
    const [calendarEroare, setCalendarEroare] = useState('');
    const [validare, setValidare] = useState('');
    const zile = lot ? lot.map(z => z.data) : interval(campuri.deLa,campuri.panaLa);
    const alese = lot ? zile : zile.filter(data => !excluse.has(data));
    useEffect(() => {
        if (lot) { setCalendarGata(true); return; }
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
    }, [campuri.deLa,campuri.panaLa,i18n.language,lot]);
    function seteaza<K extends keyof typeof campuri>(k: K, v: string) {
        setCampuri(c => {
            const nou = {...c, [k]: v};
            if (k === 'angajat' && nou.inlocuitor === v) nou.inlocuitor = '';
            return nou;
        });
        if (k === 'deLa' || k === 'panaLa') setExcluse(new Set());
    }
    function incarcaFisier(e: ChangeEvent<HTMLInputElement>) {
        const fisier = e.target.files?.[0];
        if (!fisier) return;
        const mime = fisier.type || (fisier.name.toLowerCase().endsWith('.pdf') ? 'application/pdf' : fisier.name.toLowerCase().endsWith('.png') ? 'image/png' : fisier.name.toLowerCase().match(/\.jpe?g$/) ? 'image/jpeg' : '');
        if (!['application/pdf', 'image/png', 'image/jpeg'].includes(mime) || fisier.size < 1 || fisier.size > 10 * 1024 * 1024) {
            setValidare(t('absente.document_invalid'));
            e.target.value = '';
            return;
        }
        const cititor = new FileReader();
        cititor.onload = () => {
            const rezultat = typeof cititor.result === 'string' ? cititor.result : '';
            const virgula = rezultat.indexOf(',');
            const b64 = virgula >= 0 ? rezultat.slice(virgula + 1) : '';
            if (!b64) { setValidare(t('absente.document_invalid')); return; }
            setValidare('');
            setCampuri(c => ({...c, documentId: '', documentNume: fisier.name, documentMime: mime, documentContinutB64: b64}));
        };
        cititor.onerror = () => setValidare(t('absente.document_invalid'));
        cititor.readAsDataURL(fisier);
    }
    function renunta() { if ((JSON.stringify(initial) === JSON.stringify(campuri) && !excluse.size) || window.confirm(t('ui.renunta'))) inchide(); }
    const inlocuitoriDisponibili = date.angajati.filter(a => (!acelasiAngajatLot || a.id !== campuri.angajat) && (a.activ || a.id === campuri.inlocuitor));
    const alteZileActive = absente.zile.filter(z => z.activ && (z.stare_aprobare || 'IN_ASTEPTARE') !== 'RESPINSA' && !(lot ? lot.some(l => l.angajat_id === z.angajat_id && l.data === z.data) : zi && zi.angajat_id === z.angajat_id && zi.data === z.data));
    const zileInlocuitorAbsent = campuri.inlocuitor && campuri.stareAprobare !== 'RESPINSA' ? alese.filter(data => alteZileActive.some(z => z.angajat_id === campuri.inlocuitor && z.data === data)) : [];
    const zileTitularInlocuieste = campuri.stareAprobare !== 'RESPINSA' ? (lot ? lot.filter(l => alteZileActive.some(z => z.inlocuitor_id === l.angajat_id && z.data === l.data)).map(l => l.data) : campuri.angajat ? alese.filter(data => alteZileActive.some(z => z.inlocuitor_id === campuri.angajat && z.data === data)) : []) : [];
    const zileDublaInlocuire = campuri.inlocuitor && campuri.stareAprobare !== 'RESPINSA' ? [...new Set(alese.filter(data => alteZileActive.some(z => z.inlocuitor_id === campuri.inlocuitor && z.data === data)))] : [];
    return <Dialog titlu={t(modEditare ? 'absente.editeaza' : 'absente.adauga')} ocupat={ocupat} inchide={renunta}><form onSubmit={e => {
        e.preventDefault(); setValidare('');
        if ((!lot && !campuri.angajat) || !alese.length) { setValidare(t('absente.selectie_invalida')); return; }
        if (zileInlocuitorAbsent.length) { setValidare(t('absente.inlocuitor_absent')); return; }
        if (zileTitularInlocuieste.length) { setValidare(t('absente.titular_este_inlocuitor')); return; }
        if (lot) {
            if (lot.some(z => campuri.inlocuitor && campuri.inlocuitor === z.angajat_id)) { setValidare(t('absente.inlocuitor_invalid')); return; }
            salveaza(lot.map(z => new datebaza.ZiAbsenta({
                ...z,
                cod: campuri.cod || z.cod,
                activ: true,
                nr_inregistrare: campuri.nrInregistrare,
                inlocuitor_id: campuri.inlocuitor,
                stare_aprobare: campuri.stareAprobare,
                aprobat_de: campuri.aprobatDe,
                document_id: campuri.documentId,
                document_nume: campuri.documentNume,
                document_mime: campuri.documentMime,
                document_continut_b64: campuri.documentContinutB64 || undefined
            })));
            return;
        }
        const existente = absente.zile.filter(z => z.angajat_id === campuri.angajat);
        if (!zi && alese.some(data => existente.some(z => z.data === data && z.activ))) { setValidare(t('absente.zi_ocupata')); return; }
        salveaza(alese.map(data => new datebaza.ZiAbsenta({
            angajat_id: campuri.angajat,
            data,
            cod: campuri.cod,
            activ: true,
            revizie: existente.find(z => z.data === data)?.revizie || '',
            nr_inregistrare: campuri.nrInregistrare,
            inlocuitor_id: campuri.inlocuitor,
            stare_aprobare: campuri.stareAprobare,
            aprobat_de: campuri.aprobatDe,
            document_id: campuri.documentId,
            document_nume: campuri.documentNume,
            document_mime: campuri.documentMime,
            document_continut_b64: campuri.documentContinutB64 || undefined
        })));
    }}><div className="corp-dialog">{(validare || eroare) && <p className="mesaj eroare" role="alert">{validare || eroare}</p>}
        {zileDublaInlocuire.length > 0 && <p className="mesaj avertisment" role="status">{t('absente.avertisment_dubla_inlocuire', {zile: zileDublaInlocuire.join(', ')})}</p>}
        <fieldset className="campuri" disabled={ocupat}>
            {acelasiAngajatLot && <label className="latime-completa">{t('ui.angajat')} *<select name="angajat" required disabled={modEditare} value={campuri.angajat} onChange={e => seteaza('angajat',e.target.value)}><option value="">{t('absente.alege_angajat')}</option>{date.angajati.map(a => <option key={a.id} value={a.id}>{a.nume} {a.prenume} · {a.marca}{!a.activ ? ` · ${t('ui.inactiv')}` : ''}</option>)}</select></label>}
            <label className="latime-completa">{t('absente.tip')}<select name="tip" required value={campuri.cod} onChange={e => seteaza('cod',e.target.value)}>{absente.tipuri.filter(tip => tip.activ || tip.cod === sursa?.cod).map(tip => <option key={tip.cod} value={tip.cod}>{tip.cod} · {numeTip(tip,i18n.language)}{!tip.activ ? ` · ${t('ui.inactiv')}` : ''}</option>)}</select></label>
            {!lot && <><label>{t('absente.de_la')} *<input name="de_la" type="date" min="0001-01-01" max="9999-12-31" required readOnly={!!zi} value={campuri.deLa} onChange={e => seteaza('deLa',e.target.value)}/></label><label>{t('absente.pana_la')} *<input name="pana_la" type="date" min="0001-01-01" max="9999-12-31" required readOnly={!!zi} value={campuri.panaLa} onChange={e => seteaza('panaLa',e.target.value)}/></label></>}
            <label>{t('absente.nr_inregistrare')}<input name="nr_inregistrare" maxLength={64} value={campuri.nrInregistrare} onChange={e => seteaza('nrInregistrare',e.target.value)}/></label>
            <label>{t('absente.inlocuitor')}<select name="inlocuitor" value={campuri.inlocuitor} onChange={e => seteaza('inlocuitor',e.target.value)}><option value="">{t('absente.fara_inlocuitor')}</option>{inlocuitoriDisponibili.map(a => <option key={a.id} value={a.id}>{a.nume} {a.prenume} · {a.marca}{!a.activ ? ` · ${t('ui.inactiv')}` : ''}</option>)}</select></label>
            <label>{t('absente.aprobare')}<select name="stare_aprobare" value={campuri.stareAprobare} onChange={e => seteaza('stareAprobare',e.target.value)}>{stariAprobare.map(s => <option key={s} value={s}>{t(`absente.aprobare_${s}`)}</option>)}</select></label>
            <label>{t('absente.aprobat_de')}<input name="aprobat_de" maxLength={120} value={campuri.aprobatDe} onChange={e => seteaza('aprobatDe',e.target.value)}/></label>
            <label className="latime-completa">{t('absente.cerere_atasata')}<input name="cerere_fisier" type="file" accept=".pdf,application/pdf,image/png,image/jpeg" onChange={incarcaFisier}/><small className="text-discret">{t('absente.fisier_indiciu')}</small></label>
            {(campuri.documentNume || campuri.documentId) && <div className="latime-completa atasament-curent"><span>{campuri.documentNume}</span><button type="button" className="legatura discreta" onClick={() => setCampuri(c => ({...c, documentId: '', documentNume: '', documentMime: '', documentContinutB64: ''}))}>{t('absente.elimina_atasament')}</button></div>}
        </fieldset>
        {!lot && <><p className="nota-calendar">{t('absente.selectie_indiciu')}</p>{!zile.length && <p className="mesaj eroare">{t('absente.interval_invalid')}</p>}{calendarEroare && <p className="mesaj eroare" role="alert">{calendarEroare}</p>}
        {!zi && <div className="actiuni-calendar"><button type="button" className="legatura" disabled={ocupat || !zile.length} onClick={() => setExcluse(new Set())}>{t('absente.toate_zilele')}</button><button type="button" className="legatura" disabled={ocupat || !calendarGata} onClick={() => setExcluse(new Set(zile.filter(data => [0,6].includes(new Date(`${data}T12:00:00Z`).getUTCDay()) || sarbatori.has(data))))}>{t('absente.luni_vineri')}</button></div>}
        <div className="zile-absenta" role="group" aria-label={t('absente.zile_selectate')}>
            {zile.map(data => { const d = new Date(`${data}T12:00:00Z`), weekend = [0,6].includes(d.getUTCDay()); return <label key={data} className={sarbatori.has(data) ? 'zi-sarbatoare' : weekend ? 'zi-weekend' : ''}><input type="checkbox" disabled={ocupat || !!zi} checked={!excluse.has(data)} onChange={e => setExcluse(vechi => { const noi = new Set(vechi); if (e.target.checked) noi.delete(data); else noi.add(data); return noi; })}/><span>{data} · {new Intl.DateTimeFormat(i18n.language,{weekday:'short',timeZone:'UTC'}).format(d)}{sarbatori.has(data) && <small>{sarbatori.get(data)}</small>}</span></label>; })}
        </div></>}
        <p className="nota-calendar" role="status">{t('absente.selectate', {numar:alese.length})}</p>
    </div><footer className="subsol-dialog"><button type="button" className="buton secundar" disabled={ocupat} onClick={renunta}>{t('ui.anuleaza')}</button><button className="buton primar" disabled={ocupat || !alese.length}>{t('ui.salveaza')}</button></footer></form></Dialog>;
}
