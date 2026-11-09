import {useEffect, useState, type KeyboardEvent} from 'react';
import {useTranslation} from 'react-i18next';
import {AdaugaSchimbareContract, ListaIstoricContract, StergeSchimbareContract} from '../wailsjs/go/main/App';
import {acces, main} from '../wailsjs/go/models';
import {dataAfisata} from './formatare';
import {notifica} from './Toast';

// F07-T43 (I010, D101): schimbarile normei zilnice si ale dreptului la concediu, valabile de la o data.
// Valorile din fisa raman cele de dinaintea primei schimbari. Schimbarile se salveaza imediat.

/** orePeMinute transforma „7,5” sau „7.5” in 450 de minute; gol sau invalid → undefined. */
export function orePeMinute(text: string): number | undefined {
    const t = text.trim().replace(',', '.');
    if (!t) return undefined;
    const v = Number(t);
    return Number.isFinite(v) && v > 0 ? Math.round(v * 60) : NaN;
}

function oreText(minute?: number): string {
    if (minute == null) return '—';
    const h = Math.floor(minute / 60), m = minute % 60;
    return m ? `${h} h ${m} min` : `${h} h`;
}

export function IstoricContract({angajatId}: {angajatId: string}) {
    const {t} = useTranslation();
    const [lista, setLista] = useState<acces.SchimbareContract[]>([]);
    const [deLa, setDeLa] = useState('');
    const [ore, setOre] = useState('');
    const [zile, setZile] = useState('');
    const [eroare, setEroare] = useState('');
    const [ocupat, setOcupat] = useState(false);

    useEffect(() => {
        let activ = true;
        ListaIstoricContract(angajatId).then(r => {
            if (!activ) return;
            if (r.eroare) setEroare(t(r.eroare.cod, r.eroare.param ?? {}));
            else setLista(r.schimbari ?? []);
        }).catch(() => { if (activ) setEroare(t('date.eroare')); });
        return () => { activ = false; };
    }, [angajatId, t]);

    async function aplica(cerere: () => Promise<main.RaspunsIstoricContract>, mesaj: string): Promise<boolean> {
        setOcupat(true);
        setEroare('');
        try {
            const r = await cerere();
            if (r.eroare) { setEroare(t(r.eroare.cod, r.eroare.param ?? {})); return false; }
            setLista(r.schimbari ?? []);
            notifica(t(mesaj));
            return true;
        } catch {
            setEroare(t('date.eroare'));
            return false;
        } finally {
            setOcupat(false);
        }
    }

    async function adauga() {
        const norma = orePeMinute(ore);
        const concediu = zile.trim() === '' ? undefined : Number(zile);
        if (!deLa) { setEroare(t('istoric.data_invalida')); return; }
        if (Number.isNaN(norma)) { setEroare(t('istoric.norma_invalida')); return; }
        if (concediu !== undefined && (!Number.isInteger(concediu) || concediu < 0 || concediu > 366)) { setEroare(t('istoric.concediu_invalid')); return; }
        if (norma === undefined && concediu === undefined) { setEroare(t('istoric.valoare_lipsa')); return; }
        const s = new acces.SchimbareContract({angajat_id: angajatId, de_la: deLa, norma_minute: norma, concediu_zile_an: concediu, creat_la: ''});
        if (await aplica(() => AdaugaSchimbareContract(s), 'istoric.adaugata')) { setDeLa(''); setOre(''); setZile(''); }
    }

    // Enter nu trimite formularul angajatului; adauga schimbarea.
    const laEnter = (e: KeyboardEvent) => { if (e.key === 'Enter') { e.preventDefault(); void adauga(); } };

    return <details className="detalii-configurare istoric-contract">
        <summary>{t('istoric.titlu')}{lista.length ? ' ' + t('istoric.rezumat', {numar: String(lista.length)}) : ''}</summary>
        <p className="text-discret">{t('istoric.indiciu')}</p>
        {lista.length ? <table className="tabel-istoric-contract"><thead><tr>
            <th>{t('istoric.de_la')}</th><th>{t('istoric.norma')}</th><th>{t('istoric.concediu')}</th><th>{t('ui.actiuni')}</th>
        </tr></thead><tbody>{lista.map(s => <tr key={s.de_la}>
            <td>{dataAfisata(s.de_la)}</td>
            <td>{s.norma_minute == null ? t('istoric.neschimbat') : oreText(s.norma_minute)}</td>
            <td>{s.concediu_zile_an == null ? t('istoric.neschimbat') : t('ui.zile', {numar: s.concediu_zile_an})}</td>
            <td><button type="button" className="legatura" disabled={ocupat}
                onClick={() => { if (window.confirm(t('istoric.confirma_stergere', {data: dataAfisata(s.de_la)}))) void aplica(() => StergeSchimbareContract(angajatId, s.de_la), 'istoric.stearsa'); }}>{t('ui.sterge')}</button></td>
        </tr>)}</tbody></table> : <p className="text-discret">{t('istoric.gol')}</p>}
        <fieldset disabled={ocupat} className="campuri">
            <label>{t('istoric.de_la')}<input type="date" name="istoric_de_la" value={deLa} onChange={e => setDeLa(e.target.value)} onKeyDown={laEnter}/></label>
            <label>{t('istoric.norma_noua')}<input name="istoric_norma_ore" inputMode="decimal" placeholder={t('istoric.neschimbat')} value={ore} onChange={e => setOre(e.target.value)} onKeyDown={laEnter}/></label>
            <label>{t('istoric.concediu_nou')}<input name="istoric_concediu" type="number" min="0" max="366" step="1" placeholder={t('istoric.neschimbat')} value={zile} onChange={e => setZile(e.target.value)} onKeyDown={laEnter}/></label>
            <button type="button" className="buton secundar" onClick={() => void adauga()}>{t('istoric.adauga')}</button>
        </fieldset>
        {eroare && <p className="mesaj eroare" role="alert">{eroare}</p>}
    </details>;
}
