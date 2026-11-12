import {useEffect, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {ExportaSalarizare, ListaModeleExportSalarizare, SalveazaModelExportSalarizare} from '../wailsjs/go/main/App';
import {exportsalarizare, main} from '../wailsjs/go/models';
import {Dialog} from './Formulare';
import {notifica} from './Toast';

// F10-T03 (D108): exportul pentru salarizare dupa un model (CSV sau Excel) si editorul modelelor.

const PREFIX_TOTAL = 'total:';

/** numeSursa: eticheta unei surse de coloana (identitate, total din Pontaj sau text fix). */
function useNumeSursa() {
    const {t} = useTranslation();
    return (s: string) => s.startsWith(PREFIX_TOTAL) ? t('export_salarizare.sursa_total', {coloana: t(`pontaj.total_${s.slice(PREFIX_TOTAL.length)}`)}) : t(`export_salarizare.sursa_${s}`);
}

/** DialogExportSalarizare: alege modelul si scrie fisierul pentru luna din Pontaj. */
export function DialogExportSalarizare({an, luna, sectie, numeSectie, inchide}: {an: number; luna: number; sectie: string; numeSectie: string; inchide: () => void}) {
    const {t} = useTranslation();
    const [modele, setModele] = useState<exportsalarizare.Model[]>([]);
    const [model, setModel] = useState('');
    const [doarSectia, setDoarSectia] = useState(Boolean(sectie));
    const [ocupat, setOcupat] = useState(false);
    const [eroare, setEroare] = useState('');

    useEffect(() => {
        ListaModeleExportSalarizare().then(r => {
            const active = (r.modele ?? []).filter(m => m.activ);
            setModele(active);
            setModel(m => m || active[0]?.id || '');
            if (r.eroare) setEroare(t(r.eroare.cod, r.eroare.param ?? {}));
        }).catch(() => setEroare(t('date.eroare')));
    }, [t]);

    async function exporta() {
        setOcupat(true); setEroare('');
        try {
            const r: main.RaspunsExportSalarizare = await ExportaSalarizare(model, an, luna, doarSectia ? sectie : '');
            if (r.eroare) { setEroare(t(r.eroare.cod, r.eroare.param ?? {})); return; }
            if (r.anulat) return;
            notifica(t('export_salarizare.salvat', {cale: r.cale ?? '', randuri: String(r.randuri)}));
            if (r.necalculate) notifica(t('export_salarizare.necalculate', {numar: String(r.necalculate)}), 'avertisment', 8000);
            inchide();
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setOcupat(false);
        }
    }

    return <Dialog titlu={t('export_salarizare.titlu')} ocupat={ocupat} inchide={inchide}>
        <div className="corp-dialog">
            <p className="text-discret">{t('export_salarizare.indiciu', {luna: `${String(luna).padStart(2, '0')}.${an}`})}</p>
            {modele.length ? <label className="latime-completa">{t('export_salarizare.model')}<select name="model_export_salarizare" value={model} onChange={e => setModel(e.target.value)}>
                {modele.map(m => <option key={m.id} value={m.id}>{m.denumire} ({m.format})</option>)}
            </select></label> : <p>{t('export_salarizare.fara_modele')}</p>}
            {sectie && <label className="bifa"><input type="checkbox" checked={doarSectia} onChange={e => setDoarSectia(e.target.checked)}/>{t('export_salarizare.doar_sectia', {sectie: numeSectie})}</label>}
            {eroare && <p className="mesaj eroare" role="alert">{eroare}</p>}
        </div>
        <footer className="subsol-dialog">
            <button type="button" className="buton secundar" disabled={ocupat} onClick={inchide}>{t('ui.anuleaza')}</button>
            <button type="button" className="buton primar" disabled={ocupat || !model} onClick={() => void exporta()}>{t('export_salarizare.exporta')}</button>
        </footer>
    </Dialog>;
}

function modelNou(): exportsalarizare.Model {
    return new exportsalarizare.Model({id: '', denumire: '', format: 'CSV', separator: ';', bom: true, separator_zecimal: ',', zecimale: 2, forma_ore: 'ZECIMAL', antet: true, coloane: [{sursa: 'marca', titlu: 'Marca'}], activ: true, revizie: ''});
}

/** SetariExportSalarizare: lista si editorul modelelor (Setari › General). */
export function SetariExportSalarizare() {
    const {t} = useTranslation();
    const numeSursa = useNumeSursa();
    const [modele, setModele] = useState<exportsalarizare.Model[]>([]);
    const [surse, setSurse] = useState<string[]>([]);
    const [ales, setAles] = useState('');
    const [model, setModel] = useState<exportsalarizare.Model | null>(null);
    const [eroare, setEroare] = useState('');
    const [ocupat, setOcupat] = useState(false);

    function aplica(r: main.RaspunsModeleExport, id?: string) {
        setModele(r.modele ?? []);
        setSurse(r.surse ?? []);
        const tinta = id ?? ales ?? r.modele?.[0]?.id ?? '';
        const gasit = (r.modele ?? []).find(m => m.id === tinta) ?? r.modele?.[0];
        setAles(gasit?.id ?? '');
        setModel(gasit ? new exportsalarizare.Model(structuredClone(gasit)) : null);
        if (r.eroare) setEroare(t(r.eroare.cod, r.eroare.param ?? {}));
    }
    useEffect(() => { ListaModeleExportSalarizare().then(r => aplica(r)).catch(() => setEroare(t('date.eroare'))); }, []); // eslint-disable-line react-hooks/exhaustive-deps

    const schimba = (v: Partial<exportsalarizare.Model>) => setModel(m => m ? new exportsalarizare.Model({...m, ...v}) : m);
    const schimbaColoana = (i: number, v: Partial<exportsalarizare.Coloana>) => schimba({coloane: model!.coloane.map((c, j) => j === i ? {...c, ...v} : c) as exportsalarizare.Coloana[]});
    const muta = (i: number, d: number) => { const c = [...model!.coloane]; const [x] = c.splice(i, 1); c.splice(i + d, 0, x); schimba({coloane: c}); };

    async function salveaza() {
        if (!model) return;
        setOcupat(true); setEroare('');
        try {
            const r = await SalveazaModelExportSalarizare(model);
            if (r.eroare) setEroare(t(r.eroare.cod, r.eroare.param ?? {}));
            else { aplica(r, r.model?.id); notifica(t('export_salarizare.model_salvat')); }
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setOcupat(false);
        }
    }

    return <section className="panou card-setari setari-export-salarizare">
        <h2>{t('export_salarizare.titlu_setari')}</h2>
        <p>{t('export_salarizare.descriere')}</p>
        <div className="campuri">
            <label>{t('export_salarizare.model')}<select name="model_editat" value={ales} onChange={e => aplica({modele, surse} as main.RaspunsModeleExport, e.target.value)}>
                {modele.map(m => <option key={m.id} value={m.id}>{m.denumire}{m.activ ? '' : ` · ${t('ui.inactiv')}`}</option>)}
            </select></label>
            <button type="button" className="buton secundar" disabled={ocupat} onClick={() => { setAles(''); setModel(model ? new exportsalarizare.Model({...structuredClone(model), id: '', revizie: '', denumire: t('export_salarizare.copie', {nume: model.denumire})}) : modelNou()); }}>{t('export_salarizare.model_nou')}</button>
        </div>
        {model && <fieldset className="campuri" disabled={ocupat}>
            <label>{t('export_salarizare.denumire')}<input name="denumire_model_export" value={model.denumire} maxLength={100} onChange={e => schimba({denumire: e.target.value})}/></label>
            <label>{t('export_salarizare.format')}<select name="format_export" value={model.format} onChange={e => schimba({format: e.target.value})}><option value="CSV">CSV</option><option value="XLSX">Excel (.xlsx)</option></select></label>
            {model.format === 'CSV' && <label>{t('export_salarizare.separator')}<select name="separator_export" value={model.separator} onChange={e => schimba({separator: e.target.value})}><option value=";">;</option><option value=",">,</option><option value="TAB">Tab</option></select></label>}
            <label>{t('export_salarizare.separator_zecimal')}<select name="zecimal_export" value={model.separator_zecimal} onChange={e => schimba({separator_zecimal: e.target.value})}><option value=",">,</option><option value=".">.</option></select></label>
            <label>{t('export_salarizare.zecimale')}<input name="zecimale_export" type="number" min="0" max="4" value={model.zecimale} onChange={e => schimba({zecimale: Number(e.target.value)})}/></label>
            <label>{t('export_salarizare.forma_ore')}<select name="forma_ore_export" value={model.forma_ore} onChange={e => schimba({forma_ore: e.target.value})}><option value="ZECIMAL">{t('export_salarizare.ore_zecimal')}</option><option value="ORE_MINUTE">{t('export_salarizare.ore_minute')}</option></select></label>
            <label className="bifa"><input type="checkbox" checked={model.antet} onChange={e => schimba({antet: e.target.checked})}/>{t('export_salarizare.antet')}</label>
            {model.format === 'CSV' && <label className="bifa"><input type="checkbox" checked={model.bom} onChange={e => schimba({bom: e.target.checked})}/>{t('export_salarizare.bom')}</label>}
            <label className="bifa"><input type="checkbox" checked={model.activ} onChange={e => schimba({activ: e.target.checked})}/>{t('export_salarizare.activ')}</label>
        </fieldset>}
        {model && <table className="tabel-coloane-export"><thead><tr><th>#</th><th>{t('export_salarizare.sursa')}</th><th>{t('export_salarizare.titlu_coloana')}</th><th>{t('export_salarizare.valoare_fixa')}</th><th>{t('ui.actiuni')}</th></tr></thead>
            <tbody>{model.coloane.map((c, i) => <tr key={i}>
                <td>{i + 1}</td>
                <td><select name={`sursa_export_${i}`} aria-label={t('export_salarizare.sursa')} disabled={ocupat} value={c.sursa} onChange={e => schimbaColoana(i, {sursa: e.target.value})}>{surse.map(s => <option key={s} value={s}>{numeSursa(s)}</option>)}</select></td>
                <td><input name={`titlu_export_${i}`} aria-label={t('export_salarizare.titlu_coloana')} disabled={ocupat} value={c.titlu} maxLength={100} onChange={e => schimbaColoana(i, {titlu: e.target.value})}/></td>
                <td>{c.sursa === 'CONSTANTA' ? <input name={`valoare_export_${i}`} aria-label={t('export_salarizare.valoare_fixa')} disabled={ocupat} value={c.valoare ?? ''} maxLength={100} onChange={e => schimbaColoana(i, {valoare: e.target.value})}/> : '—'}</td>
                <td><button type="button" className="legatura" disabled={ocupat || i === 0} onClick={() => muta(i, -1)}>↑</button> <button type="button" className="legatura" disabled={ocupat || i === model.coloane.length - 1} onClick={() => muta(i, 1)}>↓</button> <button type="button" className="legatura" disabled={ocupat || model.coloane.length === 1} onClick={() => schimba({coloane: model.coloane.filter((_, j) => j !== i)})}>{t('ui.sterge')}</button></td>
            </tr>)}</tbody>
        </table>}
        {model && <div className="unelte-calendar">
            <button type="button" className="buton secundar" disabled={ocupat} onClick={() => schimba({coloane: [...model.coloane, new exportsalarizare.Coloana({sursa: 'total:OL', titlu: ''})]})}>{t('export_salarizare.adauga_coloana')}</button>
            <button type="button" className="buton primar" disabled={ocupat} onClick={() => void salveaza()}>{t('ui.salveaza')}</button>
        </div>}
        <p className="text-discret">{t('export_salarizare.doar_admin')}</p>
        {eroare && <p className="mesaj eroare" role="alert">{eroare}</p>}
    </section>;
}
