import {useEffect, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {ListaSabloaneAI, SalveazaSablonAI, StergeSablonAI} from '../wailsjs/go/main/App';
import {ai, datebaza} from '../wailsjs/go/models';

export function SabloaneAsistentAI({intentie, blocat, onIncarca}: {intentie?: ai.IntentiePlan; blocat: boolean; onIncarca: (i: ai.IntentiePlan) => void}) {
    const {t} = useTranslation();
    const [lista, setLista] = useState<datebaza.SablonAI[]>([]);
    const [nume, setNume] = useState('');
    const [selectat, setSelectat] = useState('');
    const [ocupat, setOcupat] = useState(false);
    const [eroare, setEroare] = useState('');
    const [mesaj, setMesaj] = useState('');
    useEffect(() => {
        let activ = true;
        ListaSabloaneAI().then(r => { if (!activ) return; if (r.eroare) setEroare(t(r.eroare.cod, r.eroare.param ?? {})); else setLista(r.sabloane ?? []); }).catch(() => {if (activ) setEroare(t('date.eroare'));});
        return () => {activ = false;};
    }, []);
    async function modifica(sterge: boolean) {
        if (ocupat || blocat) return;
        setOcupat(true); setEroare(''); setMesaj('');
        try {
            const r = sterge ? await StergeSablonAI(selectat) : await SalveazaSablonAI(nume.trim(), intentie!);
            if (r.eroare) setEroare(t(r.eroare.cod, r.eroare.param ?? {}));
            else {setLista(r.sabloane ?? []); setSelectat(''); setNume(''); setMesaj(t(sterge ? 'ai.sablon_sters' : 'ai.sablon_salvat'));}
        } catch {setEroare(t('date.eroare'));}
        finally {setOcupat(false);}
    }
    return <details className="sabloane-ai"><summary>{t('ai.sabloane')}</summary>
        <p className="text-discret">{t('ai.sabloane_indiciu')}</p>
        {eroare && <p role="alert">{eroare}</p>}{mesaj && <p role="status">{mesaj}</p>}
        <fieldset disabled={blocat || ocupat}>
            <label>{t('ai.sablon_ales')}<select value={selectat} onChange={e => setSelectat(e.target.value)}>
                <option value="">{t('ai.plan_alege')}</option>{lista.map(s => <option key={s.id} value={s.id}>{s.nume}</option>)}
            </select></label>
            <div className="actiuni-plan-ai">
                <button type="button" className="buton secundar mic" disabled={!selectat} onClick={() => {
                    const s = lista.find(s => s.id === selectat); if (s) onIncarca(new ai.IntentiePlan({...s.intentie, an_de_confirmat: true, clarificari: []}));
                }}>{t('ai.sablon_incarca')}</button>
                <button type="button" className="buton secundar mic" disabled={!selectat} onClick={() => {if (window.confirm(t('ai.sablon_confirma_stergere'))) void modifica(true);}}>{t('ai.sablon_sterge')}</button>
            </div>
            {intentie && <><label>{t('ai.sablon_nume')}<input maxLength={100} value={nume} onChange={e => setNume(e.target.value)}/></label>
                <button type="button" className="buton secundar mic" disabled={!nume.trim()} onClick={() => void modifica(false)}>{t('ai.sablon_salveaza')}</button></>}
        </fieldset>
    </details>;
}
