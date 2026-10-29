// Ecranul jurnalului de modificări (F12-T02): cine a schimbat ce și când, cu filtre și export.
import {useEffect, useMemo, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {CitesteJurnal, ExportaJurnalCSV, StareAudit} from '../wailsjs/go/main/App';
import {audit, datebaza} from '../wailsjs/go/models';
import {notifica} from './Toast';

const PE_PAGINA = 100;
const CAMPURI_TEHNICE = new Set(['revizie', 'modificat_la', 'creat_la']);

type Valoare = unknown;

function parseaza(text: string): Valoare {
    if (!text) return undefined;
    try {
        return JSON.parse(text);
    } catch {
        return text;
    }
}

function scurt(v: Valoare): string {
    if (v === undefined || v === null || v === '') return '—';
    const s = typeof v === 'string' ? v : JSON.stringify(v);
    return s.length > 80 ? s.slice(0, 77) + '…' : s;
}

// diferente întoarce câmpurile schimbate între valoarea veche și cea nouă.
function diferente(vechi: string, nou: string): {camp: string; inainte: string; dupa: string}[] {
    const a = parseaza(vechi), b = parseaza(nou);
    const obiect = (x: Valoare): x is Record<string, Valoare> => typeof x === 'object' && x !== null && !Array.isArray(x);
    if (!obiect(a) || !obiect(b)) return [{camp: '', inainte: scurt(a), dupa: scurt(b)}];
    const chei = Array.from(new Set([...Object.keys(a), ...Object.keys(b)])).sort();
    return chei
        .filter(k => !CAMPURI_TEHNICE.has(k) && JSON.stringify(a[k]) !== JSON.stringify(b[k]))
        .map(k => ({camp: k, inainte: scurt(a[k]), dupa: scurt(b[k])}));
}

function rezumatNou(nou: string): string {
    const b = parseaza(nou);
    if (typeof b !== 'object' || b === null || Array.isArray(b)) return scurt(b);
    return Object.entries(b as Record<string, Valoare>)
        .filter(([k, v]) => !CAMPURI_TEHNICE.has(k) && v !== '' && v !== null && v !== undefined)
        .slice(0, 4).map(([k, v]) => `${k}: ${scurt(v)}`).join(' · ');
}

export function Jurnal({date}: {date: datebaza.Date}) {
    const {t, i18n} = useTranslation();
    const [filtru, setFiltru] = useState({tabel: '', angajat_id: '', luna: '', text: ''});
    const [deplasare, setDeplasare] = useState(0);
    const [pagina, setPagina] = useState<datebaza.PaginaJurnal | null>(null);
    const [stareAudit, setStareAudit] = useState<audit.Rezultat | null>(null);
    const [ocupat, setOcupat] = useState(false);
    const [eroare, setEroare] = useState('');

    const angajati = useMemo(() => new Map(date.angajati.map(a => [a.id, a])), [date.angajati]);
    const filtruCerere = useMemo(() => new datebaza.FiltruJurnal({...filtru, limita: PE_PAGINA, deplasare}), [filtru, deplasare]);

    useEffect(() => {
        StareAudit().then(r => setStareAudit(r.rezultat)).catch(() => {});
    }, []);

    useEffect(() => {
        let activ = true;
        const temporizator = window.setTimeout(() => {
            setOcupat(true);
            CitesteJurnal(filtruCerere).then(r => {
                if (!activ) return;
                if (r.eroare) setEroare(t(r.eroare.cod, r.eroare.param ?? {}));
                else { setEroare(''); setPagina(r.pagina); }
            }).catch(() => { if (activ) setEroare(t('date.eroare')); }).finally(() => { if (activ) setOcupat(false); });
        }, 250);
        return () => { activ = false; window.clearTimeout(temporizator); };
    }, [filtruCerere, t]);

    function schimba(camp: keyof typeof filtru, valoare: string) {
        setFiltru(f => ({...f, [camp]: valoare}));
        setDeplasare(0);
    }

    function eticheta(tabel: string): string {
        const k = `jurnal.tabel_${tabel}`;
        const text = t(k);
        return text === k ? tabel : text;
    }

    function rand(i: datebaza.IntrareJurnal): string {
        const [id, zi] = i.rand_id.split('/');
        const a = angajati.get(id);
        if (a) return `${a.nume} ${a.prenume} (${a.marca})${zi ? ` · ${zi}` : ''}`;
        return i.rand_id;
    }

    async function exporta() {
        setOcupat(true);
        try {
            const r = await ExportaJurnalCSV(new datebaza.FiltruJurnal({...filtru}), i18n.language);
            if (r.eroare) setEroare(t(r.eroare.cod, r.eroare.param ?? {}));
            else if (!r.anulat) notifica(t('jurnal.exportat', {numar: String(r.numar), cale: r.cale ?? ''}));
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setOcupat(false);
        }
    }

    const format = new Intl.DateTimeFormat(i18n.language, {dateStyle: 'short', timeStyle: 'medium'});
    const total = pagina?.total ?? 0;
    return <section className="panou ecran-jurnal">
        {stareAudit && stareAudit.stare !== 'inactiv' && <p className={`stare-jurnal ${stareAudit.stare}`}>
            {t(stareAudit.stare === 'ok' ? 'jurnal.sigilat' : 'jurnal.incident', {moment: format.format(new Date(stareAudit.verificat_la))})}
        </p>}
        <div className="filtre-jurnal">
            <label>{t('jurnal.filtru_tabel')}<select value={filtru.tabel} onChange={e => schimba('tabel', e.target.value)}>
                <option value="">{t('jurnal.toate')}</option>
                {(pagina?.tabele ?? []).map(x => <option key={x} value={x}>{eticheta(x)}</option>)}
            </select></label>
            <label>{t('jurnal.filtru_angajat')}<select value={filtru.angajat_id} onChange={e => schimba('angajat_id', e.target.value)}>
                <option value="">{t('jurnal.toti')}</option>
                {[...date.angajati].sort((a, b) => `${a.nume} ${a.prenume}`.localeCompare(`${b.nume} ${b.prenume}`, i18n.language))
                    .map(a => <option key={a.id} value={a.id}>{a.nume} {a.prenume} ({a.marca})</option>)}
            </select></label>
            <label>{t('jurnal.filtru_luna')}<input type="month" value={filtru.luna} onChange={e => schimba('luna', e.target.value)}/></label>
            <label>{t('jurnal.filtru_text')}<input type="search" value={filtru.text} placeholder={t('jurnal.filtru_text_indiciu')} onChange={e => schimba('text', e.target.value)}/></label>
            <button type="button" className="buton secundar" disabled={ocupat || !total} onClick={() => void exporta()}>{t('jurnal.exporta')}</button>
        </div>
        {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}
        <p className="text-discret" role="status">{ocupat ? t('ui.incarcare') : t('jurnal.numar', {numar: String(total)})}</p>
        <div className="tabel-jurnal-derulare">
            <table className="tabel-jurnal">
                <thead><tr>
                    <th scope="col">{t('jurnal.coloana_cand')}</th><th scope="col">{t('jurnal.coloana_cine')}</th>
                    <th scope="col">{t('jurnal.coloana_tabel')}</th><th scope="col">{t('jurnal.coloana_rand')}</th>
                    <th scope="col">{t('jurnal.coloana_modificari')}</th>
                </tr></thead>
                <tbody>
                    {(pagina?.intrari ?? []).map(i => <tr key={i.id}>
                        <td>{format.format(new Date(i.cand))}</td>
                        <td>{i.cine}</td>
                        <td>{eticheta(i.tabel)}</td>
                        <td>{rand(i)}</td>
                        <td>{!i.vechi
                            ? <span className="jurnal-creat">{t('jurnal.creat')}: {rezumatNou(i.nou)}</span>
                            : <ul className="jurnal-diferente">{diferente(i.vechi, i.nou).map((d, k) => <li key={k}>
                                {d.camp && <strong>{d.camp}: </strong>}<span className="vechi">{d.inainte}</span> → <span className="nou">{d.dupa}</span>
                            </li>)}</ul>}
                        </td>
                    </tr>)}
                </tbody>
            </table>
        </div>
        {total > PE_PAGINA && <div className="paginare-jurnal">
            <button type="button" className="buton secundar mic" disabled={ocupat || deplasare === 0} onClick={() => setDeplasare(d => Math.max(0, d - PE_PAGINA))}>◀</button>
            <span>{t('jurnal.pagina', {de_la: String(deplasare + 1), pana_la: String(Math.min(total, deplasare + PE_PAGINA)), total: String(total)})}</span>
            <button type="button" className="buton secundar mic" disabled={ocupat || deplasare + PE_PAGINA >= total} onClick={() => setDeplasare(d => d + PE_PAGINA)}>▶</button>
        </div>}
    </section>;
}
