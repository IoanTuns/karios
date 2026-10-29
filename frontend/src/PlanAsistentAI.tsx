import {useEffect, useRef, useState} from 'react';
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
    const titluRef = useRef<HTMLHeadingElement>(null);
    const corpRef = useRef<HTMLDivElement>(null);
    const text = (denumire: Record<string, string>) => denumire[i18n.language] || denumire.ro || Object.values(denumire)[0] || '';
    const aplicat = !!rezultat;
    const afisat = rezultat ?? plan;
    const etapa = aplicat ? 'rezultat' : plan ? 'verificare' : 'interpretare';
    const selectie = intentie.selectie;
    const grupuri = selectie.tip === 'sectie' ? date.sectii : selectie.tip === 'departament' ? date.departamente : selectie.tip === 'categorie' ? date.categorii_program.filter(c => c.activ) : [];
    const grupuriGasite = grupuri.some(g => g.id === selectie.valoare)
        ? grupuri.filter(g => g.id === selectie.valoare)
        : grupuri.filter(g => Object.values(g.denumire).some(n => n.toLocaleLowerCase() === selectie.valoare.toLocaleLowerCase()));
    const rezolvaValoare = grupuriGasite.length === 1 ? grupuriGasite[0].id : '';
    const complet = !!(selectie.tip && (selectie.tip === 'toti' || (selectie.tip === 'persoane' ? selectie.marci.length : rezolvaValoare)) && intentie.de_la && intentie.pana_la && intentie.cod_schimb && intentie.filtru_zile && intentie.mod_existente);

    useEffect(() => {
        // Schimbarea etapei este explicită; editarea câmpurilor nu mută focusul.
        if (corpRef.current) corpRef.current.scrollTop = 0;
        titluRef.current?.focus({preventScroll: true});
    }, [etapa]);

    const dataAfisata = (valoare: string) => {
        const data = new Date(`${valoare}T12:00:00Z`);
        return Number.isNaN(data.getTime()) ? valoare || t('ai.plan_alege') : new Intl.DateTimeFormat(i18n.language, {dateStyle: 'medium', timeZone: 'UTC'}).format(data);
    };
    const grupAfisat = selectie.tip === 'persoane'
        ? t('ai.plan_selectate_scurt', {numar: selectie.marci.length})
        : selectie.tip === 'toti' ? t('ai.plan_grup_toti')
        : text(grupuriGasite[0]?.denumire ?? {}) || selectie.valoare || t('ai.plan_alege');

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

    return <section className="plan-asistent-ai" data-etapa={etapa} aria-label={t('ai.plan_titlu')}>
        <header className="antet-plan-ai">
            <ol className="etape-plan-ai" aria-label={t('ai.plan_etape')}>
                {['interpretare', 'verificare', 'rezultat'].map((pas, index) => <li key={pas} aria-current={etapa === pas ? 'step' : undefined}>{index + 1}. {t(`ai.plan_etapa_${pas}`)}</li>)}
            </ol>
            <h3 ref={titluRef} tabIndex={-1}>{t(aplicat ? 'ai.plan_rezultat' : plan ? 'ai.plan_previzualizare' : 'ai.plan_titlu')}</h3>
            <p className="context-plan-ai">{grupAfisat} · {intentie.cod_schimb || t('ai.plan_alege')} · {dataAfisata(intentie.de_la)} – {dataAfisata(intentie.pana_la)}
                <span className="reguli-context-plan-ai">{intentie.filtru_zile ? t(`ai.plan_zile_${intentie.filtru_zile}`) : t('ai.plan_intrebare_zile')} · {intentie.mod_existente ? t(`ai.plan_context_${intentie.mod_existente}`) : t('ai.plan_intrebare_existente')}{selectie.fara_pontaj && <> · {t('ai.plan_context_fara_pontaj')}</>}</span>
            </p>
            {eroare && <p className="mesaj eroare" role="alert">{eroare}</p>}
        </header>
        <div className="corp-plan-ai" ref={corpRef}>
        {!afisat && <>
        <p className="text-discret">{t('ai.plan_indiciu')}</p>
        {intentie.mesaj_original && <details className="cerere-originala-ai"><summary>{t('ai.plan_cerere_originala')}</summary><blockquote>{intentie.mesaj_original}</blockquote></details>}
        <fieldset disabled={blocat || ocupat}>
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
        </fieldset>
        </>}
        {afisat && <div className="previzualizare-plan-ai">
            <p>{t('ai.plan_rezumat', {persoane: afisat.persoane.length, luni: afisat.luni.join(', '), noi: afisat.noi, inlocuiri: afisat.inlocuiri, identice: afisat.identice, ocupate: afisat.ocupate, blocate: afisat.blocate})}</p>
            <details><summary>{t('ai.plan_lista_persoane')}</summary><ul>{afisat.persoane.map(a => <li key={a.id}>{a.nume} ({a.marca})</li>)}</ul></details>
            <div className="tabel-plan-ai"><table>
                <caption>{t('ai.plan_detalii')}</caption>
                <thead><tr>{['persoane', 'data', 'inainte', aplicat ? 'valoare_finala' : 'dupa', 'stare'].map(c => <th key={c} scope="col">{t(`ai.plan_${c}`)}</th>)}</tr></thead>
                <tbody>{afisat.randuri.slice(pagina * 50, (pagina + 1) * 50).map(r => {
                    const salvat = r.stare === 'noua' || r.stare === 'inlocuire';
                    const valoare = aplicat && !salvat ? r.inainte : r.dupa;
                    const stareFinala = salvat ? 'salvat' : r.stare === 'ocupata' ? 'pastrat' : r.stare === 'identica' ? 'identic' : 'neaplicat';
                    return <tr key={`${r.angajat_id}/${r.data}`} data-stare={r.stare}>
                        <td>{r.nume} ({r.marca})</td><td>{r.data}</td><td>{r.inainte || '—'}</td><td>{valoare || '—'}</td><td>{t(aplicat ? `ai.plan_final_${stareFinala}` : `ai.plan_stare_${r.stare}`)}{r.motiv && <> · {t(`ai.plan_motiv_${r.motiv}`)}</>}</td>
                    </tr>;
                })}</tbody>
            </table></div>
            {afisat.randuri.length > 50 && <nav className="actiuni-propunere" aria-label={t('ai.plan_pagini')}>
                <button type="button" className="buton secundar mic" disabled={pagina === 0} onClick={() => setPagina(p => p - 1)}>{t('ai.plan_anterior')}</button>
                <span>{pagina + 1} / {Math.ceil(afisat.randuri.length / 50)}</span>
                <button type="button" className="buton secundar mic" disabled={(pagina + 1) * 50 >= afisat.randuri.length} onClick={() => setPagina(p => p + 1)}>{t('ai.plan_urmator')}</button>
            </nav>}
        </div>}
        </div>
        <footer className="subsol-plan-ai">
            {ocupat && <p role="status">{t('ai.plan_in_lucru')}</p>}
            {aplicat && afisat ? <>
                <strong role="status">{t('ai.plan_aplicate', {numar: afisat.aplicate})}</strong>
                <div className="actiuni-plan-ai">{onDeschidePontaj && afisat.luni.map(luna => <button key={luna} type="button" className="buton secundar mic" onClick={() => onDeschidePontaj(Number(luna.slice(0, 4)), Number(luna.slice(5, 7)))}>{t('ai.plan_deschide', {luna})}</button>)}</div>
            </> : plan && afisat ? <>
                <label className="optiune-plan-ai"><input type="checkbox" checked={verificat} disabled={blocat} onChange={e => setVerificat(e.target.checked)}/>{t('ai.plan_verificat')}</label>
                <div className="actiuni-plan-ai">
                    <button type="button" className="buton secundar" disabled={blocat} data-actiune="revizuieste" onClick={() => {setPlan(undefined); setVerificat(false); setPagina(0);}}>{t('ai.plan_revizuieste')}</button>
                    <button type="button" className="buton primar" disabled={blocat || !verificat || afisat.noi + afisat.inlocuiri === 0} data-actiune="confirma" onClick={() => void confirma()}>{t('ai.plan_confirma', {numar: afisat.noi + afisat.inlocuiri})}</button>
                </div>
            </> : <>
                {!complet && <p role="status">{t('ai.plan_clarifica')}</p>}
                <button type="button" className="buton primar" disabled={blocat || !complet} data-actiune="previzualizeaza" onClick={() => void pregateste()}>{t('ai.plan_previzualizeaza')}</button>
            </>}
        </footer>
    </section>;
}
