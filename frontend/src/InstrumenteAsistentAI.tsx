import {useRef, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {AnuleazaUltimaOperatieAI, PrevizualizeazaAnulareAI} from '../wailsjs/go/main/App';
import {ai, datebaza} from '../wailsjs/go/models';
import {SabloaneAsistentAI} from './SabloaneAsistentAI';

export function InstrumenteAsistentAI({blocat, onOcupat, onIncarca, onAnulat}: {
    blocat: boolean; onOcupat: (v: boolean) => void; onIncarca: (i: ai.IntentiePlan) => void; onAnulat: (id: string) => void;
}) {
    const {t} = useTranslation();
    const [anulare, setAnulare] = useState<datebaza.AnulareAI>();
    const [eroare, setEroare] = useState('');
    const [verificat, setVerificat] = useState(false);
    const [pagina, setPagina] = useState(0);
    const ocupatRef = useRef(false);
    async function executa(confirma: boolean) {
        if (ocupatRef.current || blocat || (confirma && (!anulare || !verificat))) return;
        ocupatRef.current = true; onOcupat(true); setEroare(''); setVerificat(false);
        try {
            const r = confirma ? await AnuleazaUltimaOperatieAI(anulare!.token) : await PrevizualizeazaAnulareAI();
            if (r.eroare) {setEroare(t(r.eroare.cod, r.eroare.param ?? {})); setAnulare(undefined);}
            else {setAnulare(r.anulare); setPagina(0); if (r.anulare?.anulata) onAnulat(r.anulare.operatie_id);}
        } catch {setEroare(t('date.eroare')); setAnulare(undefined);}
        finally {ocupatRef.current = false; onOcupat(false);}
    }
    return <section className="plan-asistent-ai instrumente-asistent-ai" aria-label={t('ai.instrumente')}>
        <header className="antet-plan-ai"><h3>{t('ai.instrumente')}</h3>{eroare && <p role="alert" className="mesaj eroare">{eroare}</p>}</header>
        <div className="corp-plan-ai">
            <SabloaneAsistentAI blocat={blocat} onIncarca={onIncarca}/>
            <h4>{t('ai.anulare_titlu')}</h4><p>{t('ai.anulare_indiciu')}</p>
            {!anulare?.anulata && <button type="button" className="buton secundar" disabled={blocat} onClick={() => void executa(false)}>{t('ai.anulare_pregateste')}</button>}
            {anulare && !anulare.operatie_id && <p role="status">{t('ai.anulare_nimic')}</p>}
            {anulare?.operatie_id && <>
                <p>{t('ai.anulare_rezumat', {numar: anulare.randuri.length, conflicte: anulare.conflicte, data: anulare.creat_la})}</p>
                <div className="tabel-plan-ai"><table><caption>{t('ai.anulare_detalii')}</caption>
                    <thead><tr>{['plan_persoane','plan_data','anulare_acum','anulare_restaurat','motiv'].map(c => <th key={c}>{t(`ai.${c}`)}</th>)}</tr></thead>
                    <tbody>{anulare.randuri.slice(pagina * 50, (pagina + 1) * 50).map(r => <tr key={`${r.angajat_id}/${r.data}`}><td>{r.nume}</td><td>{r.data}</td><td>{r.acum || '—'}<small className="plan-separat-anulare-ai">{t('ai.anulare_plan_separat', {coduri:r.plan_acum?.join(', ') || '—'})}</small></td><td>{r.restaurat || '—'}<small className="plan-separat-anulare-ai">{t('ai.anulare_plan_separat', {coduri:r.plan_restaurat?.join(', ') || '—'})}</small></td><td>{r.motiv ? t(`ai.anulare_motiv_${r.motiv}`) : '—'}</td></tr>)}</tbody>
                </table></div>
                {anulare.randuri.length > 50 && <nav className="actiuni-plan-ai" aria-label={t('ai.plan_pagini')}>
                    <button type="button" className="buton secundar mic" disabled={pagina === 0} onClick={() => setPagina(p => p - 1)}>{t('ai.plan_anterior')}</button>
                    <span>{pagina + 1} / {Math.ceil(anulare.randuri.length / 50)}</span>
                    <button type="button" className="buton secundar mic" disabled={(pagina + 1) * 50 >= anulare.randuri.length} onClick={() => setPagina(p => p + 1)}>{t('ai.plan_urmator')}</button>
                </nav>}
            </>}
        </div>
        <footer className="subsol-plan-ai">
            {anulare?.anulata ? <strong role="status">{t('ai.anulare_reusita')}</strong> : anulare?.operatie_id ? <>
                <label className="optiune-plan-ai"><input type="checkbox" disabled={blocat || anulare.conflicte > 0} checked={verificat} onChange={e => setVerificat(e.target.checked)}/>{t('ai.anulare_verificat')}</label>
                <button type="button" className="buton primar" disabled={blocat || !verificat || anulare.conflicte > 0} onClick={() => void executa(true)}>{t('ai.anulare_confirma')}</button>
            </> : <p>{t('ai.instrumente_indiciu')}</p>}
        </footer>
    </section>;
}
