import {useEffect, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {ConfirmaPlanAI, ListaSchimburi, PrevizualizeazaPlanAI} from '../wailsjs/go/main/App';
import {ai, datebaza} from '../wailsjs/go/models';

export function PlanAsistentAI({intentie, date, rezultat, blocat, onSchimba, onAplicat, onOcupat, onDeschidePontaj}: {
    intentie: ai.IntentiePlan;
    date: datebaza.Date;
    rezultat?: datebaza.PlanAI;
    blocat: boolean;
    onSchimba: (i: ai.IntentiePlan) => void;
    onAplicat: (p: datebaza.PlanAI) => void;
    onOcupat: (ocupat: boolean) => void;
    onDeschidePontaj?: (an: number, luna: number) => void;
}) {
    const {t, i18n} = useTranslation();
    const [programe, setPrograme] = useState<datebaza.Schimb[]>([]);
    const [plan, setPlan] = useState<datebaza.PlanAI>();
    const [eroare, setEroare] = useState('');
    const [ocupat, setOcupat] = useState(false);
    const [verificat, setVerificat] = useState(false);
    const [pagina, setPagina] = useState(0);
    const text = (denumire: Record<string, string>) => denumire[i18n.language] || denumire.ro || Object.values(denumire)[0] || '';
    const aplicat = !!rezultat;
    const afisat = rezultat ?? plan;
    const selectie = intentie.selectie;
    const grupuri = selectie.tip === 'sectie' ? date.sectii : selectie.tip === 'departament' ? date.departamente : selectie.tip === 'categorie' ? date.categorii_program.filter(c => c.activ) : [];
    const grupuriGasite = grupuri.some(g => g.id === selectie.valoare)
        ? grupuri.filter(g => g.id === selectie.valoare)
        : grupuri.filter(g => Object.values(g.denumire).some(n => n.toLocaleLowerCase() === selectie.valoare.toLocaleLowerCase()));
    const rezolvaValoare = grupuriGasite.length === 1 ? grupuriGasite[0].id : '';
    const complet = !!(selectie.tip && (selectie.tip === 'toti' || (selectie.tip === 'persoane' ? selectie.marci.length : rezolvaValoare)) && intentie.de_la && intentie.pana_la && intentie.cod_schimb && intentie.filtru_zile && intentie.mod_existente);

    useEffect(() => {
        let activ = true;
        ListaSchimburi().then(r => {
            if (!activ) return;
            if (r.eroare) setEroare(t(r.eroare.cod, r.eroare.param ?? {}));
            else setPrograme(r.schimburi.filter(p => p.activ));
        }).catch(() => { if (activ) setEroare(t('date.eroare')); });
        return () => { activ = false; };
    }, []);

    function schimba(campuri: Partial<ai.IntentiePlan>) {
        setPlan(undefined); setVerificat(false); setEroare(''); setPagina(0);
        onSchimba(new ai.IntentiePlan({...intentie, ...campuri}));
    }

    function schimbaSelectia(campuri: Partial<ai.SelectiePlan>) {
        schimba({selectie: new ai.SelectiePlan({...selectie, ...campuri})});
    }

    async function pregateste() {
        if (!complet || blocat) return;
        setOcupat(true); onOcupat(true); setEroare(''); setPlan(undefined); setVerificat(false); setPagina(0);
        try {
            const r = await PrevizualizeazaPlanAI(new ai.IntentiePlan({...intentie, selectie: {...selectie, valoare: rezolvaValoare || selectie.valoare}}));
            if (r.eroare) setEroare(t(r.eroare.cod, r.eroare.param ?? {}));
            else setPlan(r.plan);
        } catch { setEroare(t('date.eroare')); }
        finally { setOcupat(false); onOcupat(false); }
    }

    async function confirma() {
        if (!plan || !verificat || blocat) return;
        setOcupat(true); onOcupat(true); setEroare('');
        try {
            const r = await ConfirmaPlanAI(plan.intentie, plan.token);
            if (r.eroare) {
                setEroare(t(r.eroare.cod, r.eroare.param ?? {}));
                setPlan(undefined); setVerificat(false);
            } else if (r.plan) { onAplicat(r.plan); }
        } catch { setEroare(t('date.eroare')); setPlan(undefined); setVerificat(false); }
        finally { setOcupat(false); onOcupat(false); }
    }

    return <section className="plan-asistent-ai" aria-label={t('ai.plan_titlu')}>
        <h3>{t('ai.plan_titlu')}</h3>
        <p className="text-discret">{t('ai.plan_indiciu')}</p>
        {intentie.mesaj_original && <blockquote>{intentie.mesaj_original}</blockquote>}
        <fieldset disabled={blocat || ocupat || aplicat}>
            <legend>{t('ai.plan_interpretare')}</legend>
            <div className="campuri-plan-ai">
                <label>{t('ai.plan_persoane')}
                    <select value={selectie.tip} onChange={e => schimbaSelectia({tip: e.target.value, valoare: '', marci: []})}>
                        <option value="">{t('ai.plan_alege')}</option>
                        {['toti', 'sectie', 'departament', 'categorie', 'persoane'].map(tip => <option key={tip} value={tip}>{t(`ai.plan_grup_${tip}`)}</option>)}
                    </select>
                </label>
                {['sectie', 'departament', 'categorie'].includes(selectie.tip) && <label>{t(`ai.plan_grup_${selectie.tip}`)}
                    <select value={rezolvaValoare} onChange={e => schimbaSelectia({valoare: e.target.value})}>
                        <option value="">{t('ai.plan_alege')}</option>
                        {grupuri.map(g => <option key={g.id} value={g.id}>{text(g.denumire)}</option>)}
                    </select>
                </label>}
                <label>{t('ai.plan_program')}
                    <select value={intentie.cod_schimb} onChange={e => schimba({cod_schimb: e.target.value})}>
                        <option value="">{t('ai.plan_alege')}</option>
                        {programe.map(p => <option key={p.id} value={p.cod}>{p.cod} · {text(p.denumire)}</option>)}
                    </select>
                </label>
                <label>{t('ai.plan_de_la')}<input type="date" min="2000-01-01" max="2100-12-31" value={intentie.de_la} onChange={e => schimba({de_la: e.target.value})}/></label>
                <label>{t('ai.plan_pana_la')}<input type="date" min="2000-01-01" max="2100-12-31" value={intentie.pana_la} onChange={e => schimba({pana_la: e.target.value})}/></label>
            </div>
            {selectie.tip === 'persoane' && <details className="persoane-plan-ai" open>
                <summary>{t('ai.plan_selectate', {numar: selectie.marci.length})}</summary>
                <div>{date.angajati.filter(a => a.activ).map(a => <label key={a.id}>
                    <input type="checkbox" checked={selectie.marci.includes(a.marca)} onChange={e => schimbaSelectia({marci: e.target.checked ? [...selectie.marci, a.marca] : selectie.marci.filter(m => m !== a.marca)})}/>
                    {a.nume} {a.prenume} ({a.marca})
                </label>)}</div>
            </details>}
            <label className="optiune-plan-ai"><input type="checkbox" checked={selectie.fara_pontaj} onChange={e => schimbaSelectia({fara_pontaj: e.target.checked})}/>{t('ai.plan_fara_pontaj')}</label>
            <fieldset className="clarificare-plan-ai">
                <legend>{t('ai.plan_intrebare_zile')}</legend>
                {['lucratoare', 'toate', 'weekend'].map(f => <button type="button" key={f} aria-pressed={intentie.filtru_zile === f} className={`buton mic ${intentie.filtru_zile === f ? 'primar' : 'secundar'}`} onClick={() => schimba({filtru_zile: f})}>{t(`ai.plan_zile_${f}`)}</button>)}
            </fieldset>
            <fieldset className="clarificare-plan-ai">
                <legend>{t('ai.plan_intrebare_existente')}</legend>
                {['pastreaza', 'inlocuieste'].map(f => <button type="button" key={f} aria-pressed={intentie.mod_existente === f} className={`buton mic ${intentie.mod_existente === f ? 'primar' : 'secundar'}`} onClick={() => schimba({mod_existente: f})}>{t(`ai.plan_existente_${f}`)}</button>)}
            </fieldset>
            {!complet && <p role="status">{t('ai.plan_clarifica')}</p>}
            {!aplicat && <button type="button" className="buton primar" disabled={!complet} onClick={() => void pregateste()}>{t('ai.plan_previzualizeaza')}</button>}
        </fieldset>
        {eroare && <p className="mesaj eroare" role="alert">{eroare}</p>}
        {ocupat && <p role="status">{t('ai.plan_in_lucru')}</p>}
        {afisat && <div className="previzualizare-plan-ai">
            <h4>{t(aplicat ? 'ai.plan_rezultat' : 'ai.plan_previzualizare')}</h4>
            <p>{t('ai.plan_rezumat', {persoane: afisat.persoane.length, luni: afisat.luni.join(', '), noi: afisat.noi, inlocuiri: afisat.inlocuiri, identice: afisat.identice, ocupate: afisat.ocupate, blocate: afisat.blocate})}</p>
            <details><summary>{t('ai.plan_lista_persoane')}</summary><ul>{afisat.persoane.map(a => <li key={a.id}>{a.nume} ({a.marca})</li>)}</ul></details>
            <div className="tabel-plan-ai"><table>
                <caption>{t('ai.plan_detalii')}</caption>
                <thead><tr>{['persoane', 'data', 'inainte', 'dupa', 'stare'].map(c => <th key={c} scope="col">{t(`ai.plan_${c}`)}</th>)}</tr></thead>
                <tbody>{afisat.randuri.slice(pagina * 50, (pagina + 1) * 50).map(r => <tr key={`${r.angajat_id}/${r.data}`}>
                    <td>{r.nume} ({r.marca})</td><td>{r.data}</td><td>{r.inainte || '—'}</td><td>{r.dupa}</td><td>{t(`ai.plan_stare_${r.stare}`)}{r.motiv && <> · {t(`ai.plan_motiv_${r.motiv}`)}</>}</td>
                </tr>)}</tbody>
            </table></div>
            {afisat.randuri.length > 50 && <nav className="actiuni-propunere" aria-label={t('ai.plan_pagini')}>
                <button type="button" className="buton secundar mic" disabled={pagina === 0} onClick={() => setPagina(p => p - 1)}>{t('ai.plan_anterior')}</button>
                <span>{pagina + 1} / {Math.ceil(afisat.randuri.length / 50)}</span>
                <button type="button" className="buton secundar mic" disabled={(pagina + 1) * 50 >= afisat.randuri.length} onClick={() => setPagina(p => p + 1)}>{t('ai.plan_urmator')}</button>
            </nav>}
            {aplicat ? <div role="status">
                <strong>{t('ai.plan_aplicate', {numar: afisat.aplicate})}</strong>
                {afisat.luni.map(luna => <button key={luna} type="button" className="buton secundar mic" onClick={() => onDeschidePontaj?.(Number(luna.slice(0, 4)), Number(luna.slice(5, 7)))}>{t('ai.plan_deschide', {luna})}</button>)}
            </div> : <>
                <label className="optiune-plan-ai"><input type="checkbox" checked={verificat} disabled={blocat} onChange={e => setVerificat(e.target.checked)}/>{t('ai.plan_verificat')}</label>
                <button type="button" className="buton primar" disabled={blocat || !verificat || afisat.noi + afisat.inlocuiri === 0} onClick={() => void confirma()}>{t('ai.plan_confirma', {numar: afisat.noi + afisat.inlocuiri})}</button>
            </>}
        </div>}
    </section>;
}
