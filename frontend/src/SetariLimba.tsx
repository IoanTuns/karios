import {useEffect, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {Limbi, TraduceriLipsa} from '../wailsjs/go/main/App';
import {alegeLimbaInterfata, alegeLimbaRapoarte, LIMBA_RAPOARTE_INTERFATA, preferintaLimbaRapoarte} from './limba';

// F03-T05/F03-T08 (D103): limba interfetei, limba implicita a rapoartelor si textele
// netraduse din fiecare limba (cheile din limba implicita fara traducere, i18n.Lipsa).
export function SetariLimba() {
    const {t, i18n} = useTranslation();
    const [limbi, setLimbi] = useState<string[]>([]);
    const [lipsa, setLipsa] = useState<Record<string, string[]>>({});
    const [rapoarte, setRapoarte] = useState(preferintaLimbaRapoarte());

    useEffect(() => {
        let activ = true;
        Limbi().then(async l => {
            const rez: Record<string, string[]> = {};
            for (const limba of l.slice(1)) rez[limba] = await TraduceriLipsa(limba);
            if (activ) { setLimbi(l); setLipsa(rez); }
        }).catch(() => {});
        return () => { activ = false; };
    }, []);

    const nume = (l: string) => t(`rapoarte.${l}`);
    return <section className="panou card-setari setari-limba">
        <h2>{t('limba.titlu')}</h2>
        <p>{t('limba.descriere')}</p>
        <div className="campuri">
            <label>{t('ui.limba_interfata')}<select name="setari_limba_interfata" value={i18n.language} onChange={e => alegeLimbaInterfata(e.target.value)}>
                {limbi.map(l => <option key={l} value={l}>{nume(l)}</option>)}
            </select></label>
            <label>{t('limba.rapoarte')}<select name="setari_limba_rapoarte" value={rapoarte} onChange={e => { alegeLimbaRapoarte(e.target.value); setRapoarte(e.target.value); }}>
                <option value={LIMBA_RAPOARTE_INTERFATA}>{t('limba.rapoarte_interfata')}</option>
                {limbi.map(l => <option key={l} value={l}>{nume(l)}</option>)}
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
    </section>;
}
