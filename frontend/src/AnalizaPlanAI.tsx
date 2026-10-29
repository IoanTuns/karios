import {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {datebaza} from '../wailsjs/go/models';

export function AnalizaPlanAI({plan}: {plan: datebaza.PlanAI}) {
    const {t} = useTranslation();
    const [cautare, setCautare] = useState('');
    const [selectie, setSelectie] = useState('');
    const [pagina, setPagina] = useState(0);
    const normal = (s: string) => s.normalize('NFD').replace(/[\u0300-\u036f]/g, '').toLocaleLowerCase();
    const persoane = (plan.selectie_explicata ?? []).filter(a => normal(`${a.nume} ${a.marca}`).includes(normal(cautare)) && (!selectie || (selectie === 'inclus' ? a.inclus : !a.inclus)));
    const sumar = (s: datebaza.SumarPlanAI) => <li key={s.cheie}><strong>{s.nume}{s.nume !== s.cheie && <> ({s.cheie})</>}</strong><p>{t('ai.sumar_detalii', {...s})}</p></li>;
    return <>
        <details className="sumar-plan-ai" open><summary>{t('ai.sumar_luni')}</summary><ul>{(plan.sumar_luni ?? []).map(sumar)}</ul></details>
        <details className="sumar-plan-ai"><summary>{t('ai.sumar_persoane')}</summary><ul>{(plan.sumar_persoane ?? []).map(sumar)}</ul></details>
        <details className="explicatii-selectie-ai"><summary>{t('ai.selectie_explicata')}</summary>
            <p className="text-discret">{t('ai.selectie_indiciu')}</p>
            <div className="filtre-detalii-ai">
                <label>{t('ai.cauta_persoana')}<input value={cautare} onChange={e => {setCautare(e.target.value);setPagina(0);}}/></label>
                <label>{t('ai.selectie_filtru')}<select value={selectie} onChange={e => {setSelectie(e.target.value);setPagina(0);}}>
                    <option value="">{t('ai.filtru_toate')}</option><option value="inclus">{t('ai.selectie_inclus')}</option><option value="exclus">{t('ai.selectie_exclus')}</option>
                </select></label>
            </div>
            <p role="status">{t('ai.rezultate_numar', {numar: persoane.length})}</p>
            <ul>{persoane.slice(pagina * 50, (pagina + 1) * 50).map(a => <li key={a.id}>{a.nume} ({a.marca}) · <strong>{t(a.inclus ? 'ai.selectie_inclus' : 'ai.selectie_exclus')}</strong> · {t(`ai.selectie_motiv_${a.motiv}`)}</li>)}</ul>
            {persoane.length > 50 && <nav className="actiuni-plan-ai" aria-label={t('ai.plan_pagini')}>
                <button type="button" className="buton secundar mic" disabled={pagina === 0} onClick={() => setPagina(p => p - 1)}>{t('ai.plan_anterior')}</button>
                <span>{pagina + 1} / {Math.ceil(persoane.length / 50)}</span>
                <button type="button" className="buton secundar mic" disabled={(pagina + 1) * 50 >= persoane.length} onClick={() => setPagina(p => p + 1)}>{t('ai.plan_urmator')}</button>
            </nav>}
        </details>
    </>;
}
