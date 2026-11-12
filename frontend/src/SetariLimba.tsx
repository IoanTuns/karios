import {useCallback, useEffect, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {ImportaFolderTraduceri, Limbi, ScrieFisierTraduceri, StareTraduceri, StergeTraduceriLimba, TraduceriLipsa, VerificaFolderTraduceri} from '../wailsjs/go/main/App';
import {main} from '../wailsjs/go/models';
import {reincarcaTraduceri} from './i18n';
import {alegeLimbaInterfata, alegeLimbaRapoarte, LIMBA_RAPOARTE_INTERFATA, preferintaLimbaRapoarte} from './limba';
import {notifica} from './Toast';

// F03-T05/F03-T08 (D103): limba interfetei, limba implicita a rapoartelor si textele
// netraduse din fiecare limba (cheile din limba implicita fara traducere, i18n.Lipsa).
// F03-T11 (D106): traducerile din folderul karios/traduceri, importate in baza criptata.
export function SetariLimba() {
    const {t, i18n} = useTranslation();
    const [limbi, setLimbi] = useState<string[]>([]);
    const [lipsa, setLipsa] = useState<Record<string, string[]>>({});
    const [rapoarte, setRapoarte] = useState(preferintaLimbaRapoarte());
    const [versiune, setVersiune] = useState(0);

    useEffect(() => {
        let activ = true;
        Limbi().then(async l => {
            const rez: Record<string, string[]> = {};
            for (const limba of l.slice(1)) rez[limba] = await TraduceriLipsa(limba);
            if (activ) { setLimbi(l); setLipsa(rez); }
        }).catch(() => {});
        return () => { activ = false; };
    }, [versiune]);

    // O limba adaugata prin import nu are nume in cataloage; se afiseaza codul.
    const nume = (l: string) => i18n.exists(`rapoarte.${l}`) ? t(`rapoarte.${l}`) : l;
    // Rapoartele au textele lor doar in romana si engleza.
    const limbiRapoarte = limbi.filter(l => l === 'ro' || l === 'en');
    return <section className="panou card-setari setari-limba">
        <h2>{t('limba.titlu')}</h2>
        <p>{t('limba.descriere')}</p>
        <div className="campuri">
            <label>{t('ui.limba_interfata')}<select name="setari_limba_interfata" value={i18n.language} onChange={e => alegeLimbaInterfata(e.target.value)}>
                {limbi.map(l => <option key={l} value={l}>{nume(l)}</option>)}
            </select></label>
            <label>{t('limba.rapoarte')}<select name="setari_limba_rapoarte" value={rapoarte} onChange={e => { alegeLimbaRapoarte(e.target.value); setRapoarte(e.target.value); }}>
                <option value={LIMBA_RAPOARTE_INTERFATA}>{t('limba.rapoarte_interfata')}</option>
                {limbiRapoarte.map(l => <option key={l} value={l}>{nume(l)}</option>)}
            </select><small>{t('limba.rapoarte_indiciu')}</small></label>
        </div>
        <h3>{t('limba.traduceri_lipsa')}</h3>
        {limbi.slice(1).map(l => {
            const chei = lipsa[l] ?? [];
            return <div key={l} className="traduceri-lipsa">
                <p>{chei.length ? t('limba.lipsa_numar', {limba: nume(l), numar: String(chei.length)}) : t('limba.lipsa_niciuna', {limba: nume(l)})}</p>
                {chei.length > 0 && <details><summary>{t('limba.lipsa_lista')}</summary>
                    <ul className="lista-chei">{chei.map(c => <li key={c}><code>{c}</code> — {i18n.getFixedT(limbi[0])(c)}</li>)}</ul>
                </details>}
            </div>;
        })}
        <TraduceriFolder limbi={limbi} nume={nume} schimbat={() => setVersiune(v => v + 1)}/>
    </section>;
}

function TraduceriFolder({limbi, nume, schimbat}: {limbi: string[]; nume: (l: string) => string; schimbat: () => void}) {
    const {t} = useTranslation();
    const [stare, setStare] = useState<main.RaspunsTraduceri | null>(null);
    const [rezumate, setRezumate] = useState<main.RaspunsTraduceri['rezumate']>(undefined);
    const [limbaFisier, setLimbaFisier] = useState('en');
    const [eroare, setEroare] = useState('');
    const [ocupat, setOcupat] = useState(false);

    const citeste = useCallback(() => { StareTraduceri().then(setStare).catch(() => {}); }, []);
    useEffect(citeste, [citeste]);

    async function ruleaza(cerere: () => Promise<main.RaspunsTraduceri>, dupa?: (r: main.RaspunsTraduceri) => Promise<void> | void) {
        setOcupat(true); setEroare('');
        try {
            const r = await cerere();
            setStare(r);
            if (r.eroare) setEroare(t(r.eroare.cod, r.eroare.param ?? {}));
            else if (dupa) await dupa(r);
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setOcupat(false);
        }
    }

    const deImportat = (rezumate ?? []).some(r => r.noi + r.modificate + r.revenite > 0);
    const corecturi = Object.entries(stare?.corecturi ?? {}).sort(([a], [b]) => a.localeCompare(b));
    return <div className="traduceri-folder">
        <h3>{t('traduceri.titlu')}</h3>
        <p className="text-discret">{t('traduceri.descriere')}</p>
        {stare?.folder && <p className="text-discret"><code>{t('traduceri.folder', {folder: stare.folder})}</code></p>}
        <div className="campuri">
            <label>{t('traduceri.limba_fisier')}<input name="traduceri_limba" list="traduceri-limbi" value={limbaFisier} maxLength={5} onChange={e => setLimbaFisier(e.target.value.trim())}/></label>
            <datalist id="traduceri-limbi">{limbi.map(l => <option key={l} value={l}>{nume(l)}</option>)}</datalist>
            <button type="button" className="buton secundar" disabled={ocupat || !limbaFisier} onClick={() => void ruleaza(() => ScrieFisierTraduceri(limbaFisier), r => notifica(t('traduceri.scris', {cale: r.cale ?? ''})))}>{t('traduceri.scrie')}</button>
            <button type="button" className="buton secundar" disabled={ocupat} onClick={() => void ruleaza(VerificaFolderTraduceri, r => setRezumate(r.rezumate ?? []))}>{t('traduceri.verifica')}</button>
            <button type="button" className="buton primar" disabled={ocupat || !deImportat} onClick={() => { if (window.confirm(t('traduceri.confirma_import'))) void ruleaza(ImportaFolderTraduceri, async () => { setRezumate(undefined); await reincarcaTraduceri(); schimbat(); notifica(t('traduceri.importate')); }); }}>{t('traduceri.importa')}</button>
        </div>
        {rezumate && (rezumate.length && deImportat ? null : <p role="status">{t('traduceri.nimic_de_importat')}</p>)}
        {rezumate && rezumate.length > 0 && <table className="tabel-traduceri"><thead><tr>
            {['col_limba', 'col_noi', 'col_modificate', 'col_revenite', 'col_neschimbate', 'col_probleme'].map(c => <th key={c}>{t(`traduceri.${c}`)}</th>)}
        </tr></thead><tbody>{rezumate.map(r => <tr key={r.limba}>
            <td>{nume(r.limba)}</td><td>{r.noi}</td><td>{r.modificate}</td><td>{r.revenite}</td><td>{r.neschimbate}</td>
            <td>{r.probleme?.length ? <details><summary>{r.probleme.length}</summary><ul className="lista-chei">{r.probleme.map(p => <li key={p.cheie}><code>{p.cheie}</code> — {t(`traduceri.motiv_${p.motiv}`)}</li>)}</ul></details> : 0}</td>
        </tr>)}</tbody></table>}
        {corecturi.length ? <ul className="lista-corecturi">{corecturi.map(([l, n]) => <li key={l}>
            {t('traduceri.corecturi', {limba: nume(l), numar: String(n)})}{' '}
            <button type="button" className="legatura" disabled={ocupat} onClick={() => { if (window.confirm(t('traduceri.confirma_stergere', {limba: nume(l)}))) void ruleaza(() => StergeTraduceriLimba(l), async () => { await reincarcaTraduceri(); schimbat(); notifica(t('traduceri.sterse')); }); }}>{t('traduceri.sterge')}</button>
        </li>)}</ul> : <p className="text-discret">{t('traduceri.fara_corecturi')}</p>}
        <p className="text-discret">{t('traduceri.doar_admin')}</p>
        {eroare && <p className="mesaj eroare" role="alert">{eroare}</p>}
    </div>;
}
