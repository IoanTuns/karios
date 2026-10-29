import {lazy, Suspense, useEffect, useRef, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {AdaugaDepartament, AdaugaSectie, AtribuieCategorieProgramLot, DateInitiale, IncarcaDemo, Informatii, SalveazaAngajat, SalveazaCategorieProgram, SalveazaNormaIntreaga, SchimbaActivAngajat} from '../wailsjs/go/main/App';
import {datebaza, main} from '../wailsjs/go/models';
import {DialogDemo, FormularAngajat, Setari, durata, denumireSectie, denumireDepartament} from './Formulare';
import {DialogImportAngajati} from './ImportDialogs';
import {ErrorBoundary} from './ErrorBoundary';
import {BannerAudit} from './Audit';
import {initializeazaSesiuneAI, useSesiuneAI} from './asistentSesiune';
import {ToastContainer} from './Toast';
import type {ContextPontaj} from './Pontaj';

const Pontaj = lazy(() => import('./Pontaj').then(m => ({ default: m.Pontaj })));
const Planificare = lazy(() => import('./Planificare').then(m => ({ default: m.Planificare })));
const Absente = lazy(() => import('./Absente').then(m => ({ default: m.Absente })));
const GenerareRapoarte = lazy(() => import('./Rapoarte').then(m => ({ default: m.GenerareRapoarte })));
const DefinireRapoarte = lazy(() => import('./Rapoarte').then(m => ({ default: m.DefinireRapoarte })));
const AsistentAI = lazy(() => import('./AsistentAI').then(m => ({ default: m.AsistentAI })));
const Sarbatori = lazy(() => import('./Sarbatori').then(m => ({ default: m.Sarbatori })));
const Schimburi = lazy(() => import('./Schimburi').then(m => ({ default: m.Schimburi })));
const CategoriiProgram = lazy(() => import('./CategoriiProgram').then(m => ({ default: m.CategoriiProgram })));
const ReguliPontaj = lazy(() => import('./ReguliPontaj').then(m => ({ default: m.ReguliPontaj })));

type Ecran = 'angajati' | 'pontaj' | 'planificare' | 'sarbatori' | 'schimburi' | 'categorii_program' | 'absente' | 'definire_rapoarte' | 'rapoarte' | 'asistent_ai' | 'reguli_pontaj' | 'setari';
const ecranePrincipale: Ecran[] = ['angajati', 'pontaj', 'planificare', 'absente', 'rapoarte', 'asistent_ai', 'setari'];
const sectiuniSetari: Ecran[] = ['setari', 'sarbatori', 'schimburi', 'categorii_program', 'reguli_pontaj', 'definire_rapoarte'];

function normalizat(text: string): string {
    return text.normalize('NFD').replace(/\p{Diacritic}/gu, '').toLocaleLowerCase();
}

function Simbol({tip}: {tip: Ecran}) {
    return <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.7" aria-hidden="true">
        {tip === 'angajati'
            ? <><circle cx="9" cy="8" r="3"/><path d="M3 21v-3a6 6 0 0 1 12 0v3M16 5a3 3 0 0 1 0 6M18 15a5 5 0 0 1 3 5"/></>
            : tip === 'absente' ? <><rect x="4" y="5" width="16" height="16" rx="2"/><path d="M4 10h16M8 3v4M16 3v4M9 15l2 2 4-4"/></>
            : tip === 'sarbatori' ? <><rect x="4" y="5" width="16" height="16" rx="2"/><path d="M4 10h16M8 3v4M16 3v4M8 14h3M8 17h3M14 14h3"/></>
            : tip === 'schimburi' ? <><circle cx="12" cy="12" r="9"/><path d="M12 6v6l4 2"/></>
            : tip === 'pontaj' ? <><rect x="3" y="4" width="18" height="17" rx="2"/><path d="M3 9h18M9 9v12M15 9v12M3 15h18"/></>
            : tip === 'planificare' ? <><rect x="3" y="5" width="18" height="16" rx="2"/><path d="M7 3v4M17 3v4M3 10h18M8 15l3 3 5-5"/></>
            : tip === 'rapoarte' ? <><path d="M14 3H6a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/><path d="M14 3v6h6M8 13h8M8 17h5"/></>
            : tip === 'asistent_ai' ? <><path d="M12 2a4 4 0 0 1 4 4v1a4 4 0 0 1-8 0V6a4 4 0 0 1 4-4zM6 10h12a2 2 0 0 1 2 2v6a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2v-6a2 2 0 0 1 2-2zM9 15h.01M15 15h.01"/></>
            : tip === 'definire_rapoarte' ? <><rect x="4" y="3" width="16" height="18" rx="2"/><path d="M4 8h16M4 16h16M9 8v8"/></>
            : tip === 'setari' ? <><circle cx="12" cy="12" r="3"/><path d="M12 2v3M12 19v3M4.9 4.9l2.1 2.1M17 17l2.1 2.1M2 12h3M19 12h3M4.9 19.1L7 17M17 7l2.1-2.1"/></>
            : <><path d="M4 7h16M4 17h16"/><circle cx="9" cy="7" r="3" fill="currentColor"/><circle cx="15" cy="17" r="3" fill="currentColor"/></>}
    </svg>;
}

function App() {
    const {t, i18n} = useTranslation();
    const [ecran, setEcran] = useState<Ecran>('angajati');
    const [date, setDate] = useState<datebaza.Date | null>(null);
    const [versiune, setVersiune] = useState('');
    const [incarcare, setIncarcare] = useState(true);
    const [ocupat, setOcupat] = useState(false);
    const blocat = useRef(false);
    const [eroare, setEroare] = useState('');
    const [mesaj, setMesaj] = useState('');
    const [cautare, setCautare] = useState('');
    const [stare, setStare] = useState('toti');
    const [sectie, setSectie] = useState('');
    const [departament, setDepartament] = useState('');
    const [editare, setEditare] = useState<datebaza.Angajat | null>(null);
    const [demo, setDemo] = useState(false);
    const [importAngajati, setImportAngajati] = useState(false);
    const [selectati, setSelectati] = useState<string[]>([]);
    const [categorieLot, setCategorieLot] = useState<string>('');
    const [pontajModificat, setPontajModificat] = useState(false);
    const [pontajOcupat, setPontajOcupat] = useState(false);
    const contextPontaj = useRef<ContextPontaj | null>(null);
    const [contextPlanificare,setContextPlanificare]=useState<{an:number;luna:number}|undefined>();
    const [filtruAbsente, setFiltruAbsente] = useState<{an: number; angajatId: string} | null>(null);
    const [sidebarCompact, setSidebarCompact] = useState(() => localStorage.getItem('pontix_sidebar_compact') === 'true');
    const [tema, setTema] = useState<'light' | 'dark'>(() => (localStorage.getItem('pontix_tema') as any) || 'light');

    const sesiuneAI = useSesiuneAI();
    useEffect(() => { initializeazaSesiuneAI(); }, []);

    useEffect(() => {
        document.documentElement.setAttribute('data-theme', tema);
        localStorage.setItem('pontix_tema', tema);
    }, [tema]);

    const comutaTema = () => setTema(t => (t === 'light' ? 'dark' : 'light'));

    const comutaSidebar = () => {
        setSidebarCompact(prev => {
            const nou = !prev;
            localStorage.setItem('pontix_sidebar_compact', String(nou));
            return nou;
        });
    };

    function aplica(raspuns: main.RaspunsDate): boolean {
        if (raspuns.eroare) {
            setEroare(t(raspuns.eroare.cod, raspuns.eroare.param ?? {}));
            return false;
        }
        if (!raspuns.date) { setEroare(t('date.eroare')); return false; }
        setDate(raspuns.date);
        return true;
    }

    useEffect(() => {
        let activ = true;
        Informatii().then(info => { if (activ) setVersiune(info.versiune); }).catch(() => {});
        DateInitiale().then(raspuns => { if (activ) aplica(raspuns); })
            .catch(() => { if (activ) setEroare(t('date.eroare')); })
            .finally(() => { if (activ) setIncarcare(false); });
        return () => { activ = false; };
    }, []);

    async function executa(actiune: () => Promise<main.RaspunsDate>, succes = 'ui.salvat'): Promise<boolean> {
        if (blocat.current) return false;
        blocat.current = true;
        setOcupat(true); setEroare(''); setMesaj('');
        try {
            if (!aplica(await actiune())) return false;
            if (succes) setMesaj(t(succes));
            return true;
        } catch { setEroare(t('date.eroare')); return false; }
        finally { blocat.current = false; setOcupat(false); }
    }

    function deschide(angajat?: datebaza.Angajat) {
        setEroare(''); setMesaj('');
        setEditare(angajat ?? new datebaza.Angajat({id: '', nume: '', prenume: '', marca: '', functie: '',
            sectie_id: '', departament_id: '', tip_contract: 'INTREAGA', norma_partiala_minute: null, norma_minute: 0,
            concediu_zile_an: null, activ: true}));
    }

    const lista = date?.angajati ?? [];
    const vizibili = lista.filter(a => normalizat(`${a.nume} ${a.prenume} ${a.marca}`).includes(normalizat(cautare.trim()))
        && (stare === 'toti' || (stare === 'activi' ? a.activ : !a.activ)) && (!sectie || a.sectie_id === sectie) && (!departament || a.departament_id === departament));
    const sectiePentru = (id: string) => {
        const gasita = date?.sectii.find(s => s.id === id);
        return gasita ? denumireSectie(gasita, i18n.language) : t('ui.fara_sectie');
    };
    const departamentPentru = (id: string) => {
        const gasit = (date?.departamente ?? []).find(d => d.id === id);
        return gasit ? denumireDepartament(gasit, i18n.language) : t('ui.fara_departament');
    };
    const esteInSetari = sectiuniSetari.includes(ecran);
    function schimbaEcran(item: Ecran) {
        if (item !== ecran && pontajModificat && !window.confirm(t('ui.renunta'))) return;
        if (item !== ecran) setPontajModificat(false);
        setEcran(item); setEroare(''); setMesaj('');
    }

    return <div className={`aplicatie ${sidebarCompact ? 'sidebar-compact' : ''}`}>
        <aside className="bara-laterala">
            <div className="marca-aplicatie">
                <span className="sigla" aria-hidden="true">p<span>·</span></span>
                <strong>{t('aplicatie.titlu')}</strong>
                <button
                    type="button"
                    className="buton-pliere-sidebar"
                    onClick={comutaSidebar}
                    title={sidebarCompact ? "Extinde meniul lateral" : "Restrânge meniul lateral (compact)"}
                    aria-label={sidebarCompact ? "Extinde meniul lateral" : "Restrânge meniul lateral"}
                >
                    {sidebarCompact ? "▶" : "◀"}
                </button>
            </div>
            <p className="eticheta-meniu">{t('ui.date_baza')}</p>
            <nav aria-label={t('ui.date_baza')}>
                {ecranePrincipale.map(item => {
                    const activ = item === 'setari' ? esteInSetari : ecran === item;
                    return <div key={item} className="grup-navigare">
                        <button className={`navigare ${item === 'setari' ? 'navigare-setari' : ''} ${activ ? 'selectat' : ''}`}
                            aria-current={activ ? 'page' : undefined} disabled={ocupat || pontajOcupat}
                            title={t(`ui.${item}`)}
                            onClick={() => schimbaEcran(item)}>
                            <Simbol tip={item}/><span className="text-navigare">{t(`ui.${item}`)}</span>
                            {item === 'asistent_ai' && (sesiuneAI.necitit || sesiuneAI.sarcina) && <span
                                className={`semnal-navigare ${sesiuneAI.necitit ? 'necitit' : 'in-lucru'}`}
                                title={t(sesiuneAI.necitit ? 'ai.sarcina_terminata' : 'ai.sarcina_in_lucru')}
                                aria-label={t(sesiuneAI.necitit ? 'ai.sarcina_terminata' : 'ai.sarcina_in_lucru')}/>}</button>
                        {item === 'setari' && esteInSetari && <div className="submeniu-setari" role="group" aria-label={t('ui.setari')}>
                            {sectiuniSetari.map(sub => <button key={sub} type="button" className={`subnavigare ${ecran === sub ? 'selectat' : ''}`}
                                aria-current={ecran === sub ? 'page' : undefined} disabled={ocupat || pontajOcupat}
                                onClick={() => schimbaEcran(sub)}>
                                <Simbol tip={sub}/>{t(sub === 'setari' ? 'ui.setari_general' : `ui.${sub}`)}</button>)}
                        </div>}
                    </div>;
                })}
            </nav>
            <div className="subsol-meniu">
                <div className="semnatura-dezvoltator" title={t('aplicatie.autor', {autor: 'Ioan Tuns'})}>
                    <span className="autor-text">{t('aplicatie.dezvoltat_de', {autor: 'Ioan Tuns'})}</span>
                </div>
                <div className="rand-subsol-info">
                    <span className="punct-local"/>{t('ui.local')}
                    {versiune && <span className="versiune">{t('aplicatie.versiune', {versiune})}</span>}
                    <button
                        type="button"
                        className="buton-comutator-tema"
                        onClick={comutaTema}
                        title={tema === 'dark' ? "Comută la modul luminos" : "Comută la modul nocturn (Dark Mode)"}
                        aria-label="Comută tema"
                    >
                        {tema === 'dark' ? '☀️ Zi' : '🌙 Noapte'}
                    </button>
                    <select name="limba_interfata" className="selector-limba" aria-label={t('ui.limba_interfata')} value={i18n.language === 'en' ? 'en' : 'ro'}
                        onChange={e => { const l = e.target.value; void i18n.changeLanguage(l); document.documentElement.lang = l; document.title = i18n.t('aplicatie.titlu'); }}>
                        <option value="ro">{t('rapoarte.ro')}</option><option value="en">{t('rapoarte.en')}</option>
                    </select>
                </div>
            </div>
        </aside>
        <main className={`continut ${ecran === 'pontaj' ? 'continut-pontaj' : ''} ${esteInSetari ? 'continut-setari' : ''}`}>
            <header className="antet">
                <div><p className="fir-navigare">{t('ui.date_baza')} <span aria-hidden="true">/</span> {esteInSetari && ecran !== 'setari' ? <>{t('ui.setari')} <span aria-hidden="true">/</span> </> : null}{t(`ui.${ecran}`)}</p>
                    <h1>{t(`ui.${ecran}`)}</h1><p className="descriere">{t(`ui.${ecran}_descriere`)}</p></div>
                {ecran === 'angajati' && date && <div className="actiuni-antet">
                    <button type="button" className="buton secundar" disabled={ocupat} onClick={() => { setEroare(''); setImportAngajati(true); }}>
                        <span aria-hidden="true">📁</span>{t('ui.importa')}</button>
                    <button type="button" className="buton primar" disabled={ocupat} onClick={() => deschide()}>
                        <span aria-hidden="true">＋</span>{t('ui.adauga_angajat')}</button>
                </div>}
            </header>
            {esteInSetari && <div className="taburi-sectiuni taburi-setari" role="group" aria-label={t('ui.setari')}>
                {sectiuniSetari.map(sub => <button key={sub} type="button" className={`buton-tab ${ecran === sub ? 'activ' : ''}`}
                    aria-pressed={ecran === sub} disabled={ocupat || pontajOcupat}
                    onClick={() => schimbaEcran(sub)}>{t(sub === 'setari' ? 'ui.setari_general' : `ui.${sub}`)}</button>)}
            </div>}
            <BannerAudit/>
            {!editare && !demo && eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}
            {mesaj && <div className="mesaj succes" role="status">{mesaj}</div>}
            {incarcare ? <div className="asteptare" role="status">{t('ui.incarcare')}</div>
                : !date ? <button className="buton secundar" disabled={ocupat} onClick={() => executa(DateInitiale, '')}>{t('ui.reincearca')}</button>
                : <ErrorBoundary key={ecran} numeModul={ecran}>
                    <Suspense fallback={<div className="asteptare" role="status">{t('ui.incarcare')}</div>}>
                    {ecran === 'pontaj' ? <Pontaj date={date} context={contextPontaj} onModificari={setPontajModificat} onOcupat={setPontajOcupat} onDeschideAbsente={(an, angajatId) => { setFiltruAbsente({an, angajatId}); setPontajModificat(false); setEcran('absente'); setEroare(''); setMesaj(''); }}/>
                    : ecran === 'planificare' ? <Planificare date={date} initial={contextPlanificare} onOcupat={setPontajOcupat} onDeschidePontaj={(an,luna)=>{contextPontaj.current={an,luna,sectie:'',departament:'',categorieProgram:'',cautare:'',filtruRanduri:'',selectie:null,derulareX:0,derulareY:0};setEcran('pontaj');setEroare('');setMesaj('');}}/>
                    : ecran === 'sarbatori' ? <Sarbatori/>
                    : ecran === 'schimburi' ? <Schimburi/>
                    : ecran === 'categorii_program' ? <CategoriiProgram date={date} ocupat={ocupat} eroare={eroare} salveaza={c=>executa(()=>SalveazaCategorieProgram(c))} reincarca={()=>executa(DateInitiale,'')} curataEroare={()=>setEroare('')}/>
                    : ecran === 'absente' ? <Absente date={date} filtruInitial={filtruAbsente}/>
                    : ecran === 'reguli_pontaj' ? <ReguliPontaj onModificari={setPontajModificat} onOcupat={setPontajOcupat}/>
                    : ecran === 'definire_rapoarte' ? <DefinireRapoarte/>
                    : ecran === 'rapoarte' ? <GenerareRapoarte date={date}/>
                    : ecran === 'asistent_ai' ? <AsistentAI date={date} onDeschidePontaj={(an, luna) => { contextPontaj.current = { an, luna, sectie: '', departament: '', categorieProgram: '', cautare: '', filtruRanduri: '', selectie: null, derulareX: 0, derulareY: 0 }; setEcran('pontaj'); setEroare(''); setMesaj(''); }}/>
                    : ecran === 'setari' ? <Setari date={date} ocupat={ocupat}
                        salveazaNorma={minute => executa(() => SalveazaNormaIntreaga(minute))}
                        adaugaSectie={denumire => executa(() => AdaugaSectie({[i18n.language]: denumire}))}
                        adaugaDepartament={(sectieID, denumire) => executa(() => AdaugaDepartament(sectieID, {[i18n.language]: denumire}))}/>
                    : <>
                        <section className="indicatori" aria-label={t('ui.angajati')}>
                            {[['ui.total', lista.length], ['ui.activi', lista.filter(a => a.activ).length],
                                ['ui.partiali', lista.filter(a => a.tip_contract === 'PARTIALA').length]].map(([cheie, numar]) =>
                                <div className="indicator" key={cheie}><span>{t(String(cheie))}</span><strong>{numar}</strong><span className="accent-indicator"/></div>)}
                        </section>
                        <section className="panou lista-angajati" aria-label={t('ui.angajati')}>
                            <div className="unelte">
                                <label className="cautare"><span className="doar-cititor">{t('ui.cautare')}</span>
                                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.7" aria-hidden="true"><circle cx="10" cy="10" r="6"/><path d="m15 15 5 5"/></svg>
                                    <input value={cautare} onChange={e => setCautare(e.target.value)} placeholder={t('ui.cautare_indiciu')}/></label>
                                <select aria-label={t('ui.stare')} value={stare} onChange={e => setStare(e.target.value)}>
                                    <option value="toti">{t('ui.toti')}</option><option value="activi">{t('ui.doar_activi')}</option><option value="inactivi">{t('ui.doar_inactivi')}</option>
                                </select>
                                <select aria-label={t('ui.sectie')} value={sectie} onChange={e => { setSectie(e.target.value); setDepartament(''); }}>
                                    <option value="">{t('ui.toate_sectiile')}</option>
                                    {date.sectii.map(s => <option key={s.id} value={s.id}>{denumireSectie(s, i18n.language)}</option>)}
                                </select>
                                <select aria-label={t('ui.departament')} value={departament} onChange={e => setDepartament(e.target.value)}><option value="">{t('ui.toate_departamentele')}</option>
                                    {(date.departamente ?? []).filter(d => !sectie || d.sectie_id === sectie).map(d => <option key={d.id} value={d.id}>{!sectie ? `${sectiePentru(d.sectie_id)} / ` : ''}{denumireDepartament(d, i18n.language)}</option>)}</select>
                            </div>
                            {selectati.length > 0 && (
                                <div className="bara-actiuni-lot" role="region" aria-label={t('ui.actiuni_lot')}>
                                    <span className="numar-selectati"><strong>{t('ui.selectati_numar', {numar: selectati.length})}</strong></span>
                                    <div className="grup-comenzi-lot">
                                        <span className="eticheta-lot">{t('ui.atribuie_categorie_lot')}:</span>
                                        <select
                                            aria-label={t('ui.atribuie_categorie_lot')}
                                            value={categorieLot}
                                            onChange={e => setCategorieLot(e.target.value)}
                                        >
                                            <option value="">{t('ui.fara_categorie_program')}</option>
                                            {(date.categorii_program ?? []).filter(c => c.activ).map(c => (
                                                <option key={c.id} value={c.id}>
                                                    {c.denumire[i18n.language] || c.denumire.ro}
                                                </option>
                                            ))}
                                        </select>
                                        <button
                                            type="button"
                                            className="buton primar mic"
                                            disabled={ocupat}
                                            onClick={async () => {
                                                const ok = await executa(() => AtribuieCategorieProgramLot(selectati, categorieLot), 'ui.categorie_lot_aplicata');
                                                if (ok) setSelectati([]);
                                            }}
                                        >
                                            {t('ui.aplica_categorie_lot')}
                                        </button>
                                        <button
                                            type="button"
                                            className="buton secundar mic"
                                            disabled={ocupat}
                                            onClick={() => setSelectati([])}
                                        >
                                            {t('ui.deselecteaza_tot')}
                                        </button>
                                    </div>
                                </div>
                            )}
                            {lista.length === 0 ? <div className="stare-goala">
                                <span className="simbol-gol" aria-hidden="true"><Simbol tip="angajati"/></span>
                                <h2>{t('ui.gol_titlu')}</h2><p>{t('ui.gol_descriere')}</p>
                                {date.demo_disponibil && <button className="buton secundar" disabled={ocupat} onClick={() => { setEroare(''); setDemo(true); }}>{t('ui.demo')}</button>}
                            </div> : <div className="tabel-derulant"><table>
                                <thead><tr>
                                    <th scope="col" className="coloana-selectie">
                                        <input
                                            type="checkbox"
                                            aria-label={t('ui.selecteaza_toti')}
                                            checked={vizibili.length > 0 && vizibili.every(a => selectati.includes(a.id))}
                                            onChange={e => {
                                                if (e.target.checked) {
                                                    setSelectati(Array.from(new Set([...selectati, ...vizibili.map(a => a.id)])));
                                                } else {
                                                    setSelectati(selectati.filter(id => !vizibili.some(a => a.id === id)));
                                                }
                                            }}
                                        />
                                    </th>
                                    {['angajat', 'marca', 'sectie', 'departament', 'contract', 'norma', 'concediu', 'stare', 'actiuni'].map(c => <th key={c} scope="col">{t(`ui.${c}`)}</th>)}
                                </tr></thead>
                                <tbody>{vizibili.map(a => <tr key={a.id} className={`${a.activ ? '' : 'rand-inactiv'} ${selectati.includes(a.id) ? 'rand-selectat' : ''}`}>
                                    <td className="coloana-selectie">
                                        <input
                                            type="checkbox"
                                            aria-label={t('ui.selecteaza_angajat', {nume: `${a.nume} ${a.prenume}`})}
                                            checked={selectati.includes(a.id)}
                                            onChange={e => {
                                                if (e.target.checked) {
                                                    setSelectati([...selectati, a.id]);
                                                } else {
                                                    setSelectati(selectati.filter(id => id !== a.id));
                                                }
                                            }}
                                        />
                                    </td>
                                    <td><div className="identitate"><span className="avatar" aria-hidden="true">{a.prenume.charAt(0)}{a.nume.charAt(0)}</span>
                                        <span><strong>{a.nume} {a.prenume}</strong>{a.functie && <small>{a.functie}</small>}</span></div></td>
                                    <td className="marca-angajat">{a.marca}</td><td>{sectiePentru(a.sectie_id)}</td><td>{departamentPentru(a.departament_id)}</td>
                                    <td><span className={`eticheta-contract ${a.tip_contract === 'PARTIALA' ? 'partial' : ''}`}>{t(a.tip_contract === 'PARTIALA' ? 'ui.partiala' : 'ui.intreaga')}</span>{a.categorie_program_id&&<small className="categorie-angajat">{date.categorii_program?.find(c=>c.id===a.categorie_program_id)?.denumire[i18n.language]||date.categorii_program?.find(c=>c.id===a.categorie_program_id)?.denumire.ro}</small>}</td>
                                    <td className="durata">{durata(a.norma_minute, t)}</td><td>{a.concediu_zile_an == null ? <span className="text-discret">{t('ui.necompletat')}</span> : t('ui.zile', {numar: a.concediu_zile_an})}</td>
                                    <td><span className={`stare-angajat ${a.activ ? 'activ' : ''}`}><i/>{t(a.activ ? 'ui.activ' : 'ui.inactiv')}</span></td>
                                    <td><div className="actiuni-rand"><button className="legatura" disabled={ocupat} aria-label={t('ui.editare_pentru', {nume: `${a.nume} ${a.prenume}`})} onClick={() => deschide(a)}>{t('ui.editeaza')}</button>
                                        <button className="legatura discreta" disabled={ocupat} aria-label={`${t(a.activ ? 'ui.dezactiveaza' : 'ui.reactiveaza')} ${a.nume} ${a.prenume}`}
                                            onClick={() => executa(() => SchimbaActivAngajat(a.id, !a.activ))}>{t(a.activ ? 'ui.dezactiveaza' : 'ui.reactiveaza')}</button></div></td>
                                </tr>)}</tbody>
                            </table>{vizibili.length === 0 && <p className="fara-rezultate">{t('ui.fara_rezultate')}</p>}</div>}
                            <footer className="subsol-tabel">{t('ui.rezultate', {vizibile: vizibili.length, total: lista.length})}</footer>
                        </section>
                    </>}
                    </Suspense>
                  </ErrorBoundary>}
        </main>
        {editare && date && <FormularAngajat key={editare.id || 'nou'} angajat={editare} date={date} ocupat={ocupat} eroare={eroare}
            inchide={() => { setEditare(null); setEroare(''); }}
            salveaza={async a => { if (await executa(() => SalveazaAngajat(a))) setEditare(null); }}/>}
        {demo && <DialogDemo ocupat={ocupat} eroare={eroare} inchide={() => { setDemo(false); setEroare(''); }}
            confirma={async () => { if (await executa(IncarcaDemo, 'ui.demo_salvat')) setDemo(false); }}/>}
        {importAngajati && <DialogImportAngajati
            inchide={() => setImportAngajati(false)}
            laSucces={(dateNoi, msg) => {
                setDate(dateNoi);
                setMesaj(msg);
                setEroare('');
            }}/>}
        <ToastContainer />
    </div>;
}

export default App;
