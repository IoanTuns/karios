import {useEffect, useLayoutEffect, useRef, useState, type FormEvent} from 'react';
import {PlanAsistentAI} from './PlanAsistentAI';
import {InstrumenteAsistentAI} from './InstrumenteAsistentAI';
import {useTranslation} from 'react-i18next';
import type {CerereConcediu} from './cerereConcediu';
import {AplicaPropunereAI, EsteAutoPornireAI, OpresteAI, PornesteAI, SeteazaAutoPornireAI, StareAI} from '../wailsjs/go/main/App';
import {ai, datebaza, main} from '../wailsjs/go/models';
import {marcheazaVizibil, opresteSarcina, pornesteSarcina, reseteazaConversatie as reseteazaSesiunea, seteazaConversatie, seteazaEroare, seteazaPerioada, useSesiuneAI, type MesajUI} from './asistentSesiune';
import {formatData} from './formatare';

function curataText(str: string): string {
    return (str || '')
        .toLowerCase()
        .normalize('NFD')
        .replace(/[\u0300-\u036f]/g, '')
        .trim();
}

function curataAlfanumeric(str: string): string {
    return curataText(str).replace(/[^a-z0-9]/g, '');
}

function gasesteAngajat(marcaSauNume: string, angajati: datebaza.Angajat[]): datebaza.Angajat | undefined {
    const raw = (marcaSauNume || '').trim();
    if (!raw || !angajati || angajati.length === 0) return undefined;

    const norm = curataText(raw);
    const alfa = curataAlfanumeric(raw);

    // 1. Potrivire exactă marcă
    const potrExact = angajati.find(a => curataText(a.marca) === norm);
    if (potrExact) return potrExact;

    // 2. Potrivire alfanumerică marcă (ex: "DEM008" vs "DEMO-008" sau "DEMO008")
    const potrAlfa = angajati.find(a => {
        const aAlfa = curataAlfanumeric(a.marca);
        if (aAlfa === alfa) return true;
        const cifreRaw = alfa.replace(/\D/g, '');
        const cifreA = aAlfa.replace(/\D/g, '');
        if (cifreRaw && cifreA && cifreRaw === cifreA) {
            const litRaw = alfa.replace(/\d/g, '');
            const litA = aAlfa.replace(/\d/g, '');
            if (litA.startsWith(litRaw) || litRaw.startsWith(litA)) {
                return true;
            }
        }
        return false;
    });
    if (potrAlfa) return potrAlfa;

    // 3. Potrivire după nume / prenume complet sau inversat
    const potrNume = angajati.find(a => {
        const np = curataText(`${a.nume} ${a.prenume}`);
        const pn = curataText(`${a.prenume} ${a.nume}`);
        const n = curataText(a.nume);
        const p = curataText(a.prenume);
        return np === norm || pn === norm ||
            np.includes(norm) || pn.includes(norm) ||
            norm.includes(np) || norm.includes(pn) ||
            (n.length >= 3 && norm.includes(n)) ||
            (p.length >= 3 && norm.includes(p));
    });
    if (potrNume) return potrNume;

    return undefined;
}

export function AsistentAI({
    date,
    perioadaEcran,
    onDeschidePontaj,
    onDeschideConcediu,
    onAjutor,
    intrebareAjutor,
}: {
    date: datebaza.Date;
    onAjutor?: (id: string, versiune?: string) => void;
    intrebareAjutor?: {text: string; id: number};
    perioadaEcran?: {an: number; luna: number};
    onDeschidePontaj?: (an: number, luna: number) => void;
    onDeschideConcediu?: (c: CerereConcediu) => void;
}) {
    const {t, i18n} = useTranslation();
    useEffect(() => {if (intrebareAjutor) setTextIntrare(intrebareAjutor.text);}, [intrebareAjutor]);
    const [stare, setStare] = useState<ai.StareMediu | null>(null);
    const [incarcare, setIncarcare] = useState(true);
    const [ocupat, setOcupat] = useState(false);
    const [mesaj, setMesaj] = useState('');
    const [eroare, setEroare] = useState('');
    const [textIntrare, setTextIntrare] = useState('');
    const sesiune = useSesiuneAI();
    const conversatie = sesiune.conversatie;
    const planuri = conversatie.filter(m => m.intentie);
    const [planSelectat, setPlanSelectat] = useState('');
    const planActiv = planuri.find(m => m.id === planSelectat) ?? planuri.at(-1);
    const [vizualizare, setVizualizare] = useState<'conversatie' | 'plan' | 'instrumente'>(planuri.length ? 'plan' : 'conversatie');
    const lucreaza = sesiune.sarcina !== '';
    const blocat = ocupat || lucreaza;
    const [autoPornire, setAutoPornire] = useState(true);

    const {an: anCurent, luna: lunaCurenta} = sesiune.perioada;
    useEffect(() => {
        if (perioadaEcran) seteazaPerioada(perioadaEcran);
    }, [perioadaEcran?.an, perioadaEcran?.luna]);

    const panouRef = useRef<HTMLDivElement>(null);
    const fluxRef = useRef<HTMLDivElement>(null);
    const urmaresteCapatuRef = useRef(true);
    const ultimMesajRef = useRef('');
    const ultimMesajId = conversatie.at(-1)?.id ?? '';

    useLayoutEffect(() => {
        const panou = panouRef.current;
        if (!panou) return;
        const dimensioneaza = () => {
            const parinte = panou.closest('main');
            const subsol = parinte ? parseFloat(getComputedStyle(parinte).paddingBottom) || 0 : 0;
            const inaltime = Math.max(260, window.innerHeight - (panou.getBoundingClientRect().top + window.scrollY) - subsol - 12);
            panou.style.setProperty('--inaltime-asistent', `${inaltime}px`);
        };
        dimensioneaza();
        window.addEventListener('resize', dimensioneaza);
        const observator = new ResizeObserver(dimensioneaza);
        if (panou.parentElement) observator.observe(panou.parentElement);
        return () => { window.removeEventListener('resize', dimensioneaza); observator.disconnect(); };
    }, []);

    async function actualizeazaStare() {
        try {
            const resp = await StareAI();
            if (resp.eroare) {
                setEroare(t(resp.eroare.cod, resp.eroare.param ?? {}));
            } else {
                setStare(resp.stare);
            }
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setIncarcare(false);
        }
    }

    useEffect(() => {
        void actualizeazaStare();
        void EsteAutoPornireAI().then(setAutoPornire).catch(() => {});
    }, []);

    useEffect(() => {
        if (!lucreaza) void actualizeazaStare();
    }, [lucreaza]);

    useEffect(() => {
        marcheazaVizibil(true);
        return () => marcheazaVizibil(false);
    }, []);

    useEffect(() => {
        // Modificarea intenției păstrează ID-ul mesajului și nu derulează istoricul.
        if (!ultimMesajId || ultimMesajRef.current === ultimMesajId) return;
        ultimMesajRef.current = ultimMesajId;
        const ultim = conversatie.at(-1);
        if (ultim?.intentie) {
            setPlanSelectat(ultim.id);
            setVizualizare('plan');
            return;
        }
        if (vizualizare === 'conversatie' && urmaresteCapatuRef.current && fluxRef.current) {
            fluxRef.current.scrollTo({top: fluxRef.current.scrollHeight, behavior: window.matchMedia('(prefers-reduced-motion: reduce)').matches ? 'instant' : 'smooth'});
        }
    }, [ultimMesajId, conversatie, vizualizare]);

    function reseteazaConversatie() {
        reseteazaSesiunea();
        setPlanSelectat(''); setVizualizare('conversatie'); urmaresteCapatuRef.current = true;
        setMesaj('');
        setEroare('');
    }

    function incarcaPlan(intentie?: ai.IntentiePlan) {
        const id = crypto.randomUUID();
        const i = intentie ?? new ai.IntentiePlan({selectie: {tip: '', valoare: '', marci: [], fara_pontaj: false}, de_la: `${anCurent}-${String(lunaCurenta).padStart(2, '0')}-01`, pana_la: `${anCurent}-${String(lunaCurenta).padStart(2, '0')}-${new Date(anCurent, lunaCurenta, 0).getDate()}`, cod_schimb: '', filtru_zile: '', mod_existente: '', mesaj_original: '', clarificari: [], an_de_confirmat: true});
        seteazaConversatie(prev => [...prev, {id, rol: 'assistant', text: t('ai.plan_titlu'), intentie: i}]);
        setPlanSelectat(id); setVizualizare('plan');
    }

    async function toggleAutoPornire(activ: boolean) {
        setAutoPornire(activ);
        try {
            await SeteazaAutoPornireAI(activ);
        } catch {
            // Revert dacă eșuează
            void EsteAutoPornireAI().then(setAutoPornire).catch(() => {});
        }
    }

    async function pornesteMotor() {
        setOcupat(true);
        setEroare('');
        try {
            const resp = await PornesteAI('');
            if (resp.eroare) {
                setEroare(t(resp.eroare.cod, resp.eroare.param ?? {}));
            } else {
                setStare(resp.stare);
            }
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setOcupat(false);
        }
    }

    async function opresteMotor() {
        setOcupat(true);
        try {
            const resp = await OpresteAI();
            setStare(resp.stare);
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setOcupat(false);
        }
    }

    async function trimiteMesaj(textDeTrimis?: string) {
        const text = (textDeTrimis ?? textIntrare).trim();
        if (!text || blocat) return;
        setVizualizare('conversatie'); urmaresteCapatuRef.current = true;

        setEroare('');
        seteazaEroare('');

        // Un „da” sau „ok” nu aplică nimic (F14-T10): propunerile se confirmă doar cu butonul
        // de pe card, unde se vede exact ce se modifică.
        const textNorm = text.toLowerCase().replace(/[.!?,]/g, '').trim();
        const esteConfirmare = ['da', 'confirm', 'confirma', 'confirmă', 'aplica', 'aplică', 'aplic', 'ok', 'de acord', 'executa', 'execută', 'fa asa', 'fă așa', 'este bine', 'e bine', 'da aplica', 'da te rog', 'sigur', 'salveaza', 'salvează'].includes(textNorm);
        const arePropuneriDeschise = conversatie.some(m => m.propuneri?.some(p => !m.propuneriAplicata?.[p.id]));
        if (esteConfirmare && arePropuneriDeschise) {
            setTextIntrare('');
            seteazaConversatie(prev => [...prev,
                {id: `u_${Date.now()}`, rol: 'user', text},
                {id: `a_${Date.now()}`, rol: 'assistant', text: t('ai.confirma_cu_butonul')}]);
            return;
        }

        setTextIntrare('');

        const mesajUtilizator: MesajUI = {
            id: `u_${Date.now()}`,
            rol: 'user',
            text,
        };

        const istoricNou = [...conversatie, mesajUtilizator];
        seteazaConversatie(() => istoricNou);

        // Ultimele 10 mesaje; propunerile nu intră în textul mesajelor asistentului, ci separat,
        // în lista de acțiuni, ca modelul să nu le copieze ca răspuns (F14-T17).
        const istoricAPI: ai.MesajConversatie[] = istoricNou
            .filter(m => m.id !== 'bun_venit')
            .slice(-10)
            .map(m => new ai.MesajConversatie({
                role: m.rol,
                content: m.rol === 'assistant' && m.propuneri?.length && (!m.text || m.text === t('ai.propuneri_titlu'))
                    ? '(propuneri afișate utilizatorului)'
                    : m.text,
            }));
        const actiuni = istoricNou.flatMap(m => (m.propuneri ?? []).map(p =>
            `${m.propuneriAplicata?.[p.id] ? '[aplicată]' : '[neconfirmată]'} ${p.descriere}`)).slice(-10);

        await pornesteSarcina(new main.CerereInterogareAI({
            istoric: istoricAPI,
            actiuni,
            an: anCurent,
            luna: lunaCurenta,
        }));
    }

    function revedePropunere(propunere: ai.PropunereActiune) {
        // Cardurile din sesiuni mai vechi trebuie și ele previzualizate înaintea aplicării.
        const p = propunere.parametri ?? {};
        const marci = Array.isArray(p.marci) ? p.marci.map(String) : p.marca ? [String(p.marca)] : [];
        const intentie = new ai.IntentiePlan({
            selectie: {tip: marci.length ? 'persoane' : '', valoare: '', marci, fara_pontaj: false},
            de_la: String(p.data_inceput ?? p.data ?? ''), pana_la: String(p.data_sfarsit ?? p.data ?? ''),
            cod_schimb: String(p.cod_schimb ?? ''), filtru_zile: '', mod_existente: '',
            mesaj_original: propunere.descriere, clarificari: [],
        });
        seteazaConversatie(prev => [...prev, {id: `plan_${Date.now()}`, rol: 'assistant', text: '', intentie}]);
    }

    async function confirmaStergere(mesajId: string, propunere: ai.PropunereActiune) {
        if (blocat || !propunere.tip.startsWith('sterge_')) return;
        setOcupat(true); setEroare('');
        try {
            const r = await AplicaPropunereAI(new main.CerereAplicaPropunereAI({an: anCurent, luna: lunaCurenta, propunere}));
            if (r.eroare) setEroare(t(r.eroare.cod, r.eroare.param ?? {}));
            else seteazaConversatie(prev => prev.map(m => m.id === mesajId ? {...m, propuneriAplicata: {...m.propuneriAplicata, [propunere.id]: true}} : m));
        } catch { setEroare(t('date.eroare')); }
        finally { setOcupat(false); }
    }

    return (
        <div className="panou-asistent-ai" ref={panouRef}>
            {/* Bara de stare a motorului local */}
            <div className="antet-stare-ai">
                <div className="indicator-stare-motor">
                    <span
                        className={`punct-stare-ai ${
                            stare?.activ ? 'activ' : stare?.instalat ? 'gata' : 'inactiv'
                        }`}
                    />
                    <div>
                        <strong>{t('ai.titlu')}</strong>
                        <span className="detaliu-stare-ai">
                            {stare?.activ
                                ? t('ai.stare_activ', {port: stare.port})
                                : stare?.instalat
                                ? t('ai.stare_gata')
                                : t('ai.stare_inactiv')}
                        </span>
                    </div>
                </div>

                <div className="actiuni-stare-motor">
                    <details className="optiuni-motor-ai">
                    <summary>{t('ai.optiuni_motor')}</summary>
                    <div className="meniu-motor-ai">
                    <p>{t('ai.model_activ')}: <strong>{stare?.model_activ || t('ai.model_nepornit')}</strong></p>
                    <p>{t('ai.context_configurat', {numar: stare?.context_configurat ?? '—'})}</p>
                    <p className="text-discret">{t('ai.context_indiciu')}</p>
                    <p className="text-discret">{t('ai.reset_indiciu')}</p>
                    <label className="comutator-autopornire-ai" title={t('ai.autopornire_descriere')}>
                        <input
                            type="checkbox"
                            checked={autoPornire}
                            onChange={e => void toggleAutoPornire(e.target.checked)}
                        />
                        <span>{t('ai.autopornire')}</span>
                    </label>
                    {stare?.instalat && !stare.activ && (
                        <button
                            type="button"
                            className="buton primar mic"
                            disabled={ocupat || incarcare}
                            onClick={pornesteMotor}
                        >
                            ⚡ {t('ai.porneste')}
                        </button>
                    )}
                    {stare?.activ && (
                        <button
                            type="button"
                            className="buton secundar mic"
                            disabled={ocupat || incarcare}
                            onClick={opresteMotor}
                        >
                            ⏹ {t('ai.opreste')}
                        </button>
                    )}
                    </div>
                    </details>
                    <button
                        type="button"
                        className="buton secundar mic buton-conversatie-noua"
                        disabled={blocat || conversatie.length <= 1}
                        onClick={reseteazaConversatie}
                        title={t('ai.conversatie_noua_descriere')}
                    >
                        🔄 {t('ai.conversatie_noua')}
                    </button>
                </div>
            </div>

            <div className="navigare-lucru-ai" role="group" aria-label={t('ai.zona_lucru')}>
                <button type="button" className={`buton mic ${vizualizare === 'conversatie' ? 'primar' : 'secundar'}`} disabled={blocat} aria-pressed={vizualizare === 'conversatie'} onClick={() => setVizualizare('conversatie')}>{t('ai.zona_conversatie')}</button>
                <button type="button" className={`buton mic ${vizualizare === 'plan' ? 'primar' : 'secundar'}`} disabled={!planActiv || blocat} aria-pressed={vizualizare === 'plan'} onClick={() => setVizualizare('plan')}>{t('ai.zona_plan')}</button>
                <button type="button" className={`buton mic ${vizualizare === 'instrumente' ? 'primar' : 'secundar'}`} disabled={blocat} aria-pressed={vizualizare === 'instrumente'} onClick={() => setVizualizare('instrumente')}>{t('ai.instrumente')}</button>
                <button type="button" className="buton secundar mic" disabled={blocat} onClick={() => incarcaPlan()}>{t('ai.plan_nou')}</button>
                {vizualizare === 'plan' && planuri.length > 1 && <select aria-label={t('ai.alege_plan')} value={planActiv?.id} disabled={blocat} onChange={e => setPlanSelectat(e.target.value)}>
                    {planuri.map((m, index) => <option key={m.id} value={m.id}>{t('ai.plan_numar', {numar: index + 1})} · {m.intentie?.de_la} – {m.intentie?.pana_la}</option>)}
                </select>}
            </div>

            <div className="perioada-asistent" hidden={vizualizare !== 'conversatie'}>
                <label>{t('ai.luna_de_lucru')}
                    <select value={lunaCurenta} disabled={blocat} onChange={e => seteazaPerioada({an: anCurent, luna: Number(e.target.value)})}>
                        {Array.from({length: 12}, (_, i) => <option key={i} value={i + 1}>{formatData(i18n.language, {month: 'long', timeZone: 'UTC'}).format(new Date(Date.UTC(2026, i, 15)))}</option>)}
                    </select>
                </label>
                <input type="number" aria-label={t('sarbatori.an')} min={2000} max={2100} value={anCurent} disabled={blocat}
                    onChange={e => { const an = Number(e.target.value); if (an >= 2000 && an <= 2100) seteazaPerioada({an, luna: lunaCurenta}); }}/>
                <span className="text-discret">{t('ai.luna_de_lucru_indiciu')}</span>
            </div>

            {(eroare || sesiune.eroare) && <div className="mesaj eroare" role="alert">{eroare || sesiune.eroare}</div>}
            {mesaj && <div className="mesaj succes" role="status">{mesaj}</div>}

            <div className="zona-plan-ai" hidden={vizualizare !== 'instrumente'}>
                {vizualizare === 'instrumente' && <InstrumenteAsistentAI blocat={blocat} onOcupat={setOcupat} onIncarca={incarcaPlan} onAnulat={id => seteazaConversatie(prev => prev.map(m => m.rezultatPlan?.operatie_id === id ? {...m, rezultatPlan: new datebaza.PlanAI({...m.rezultatPlan, anulat: true})} : m))}/>}
            </div>

            <div className="zona-plan-ai" hidden={vizualizare !== 'plan'}>
                {planActiv?.intentie && <PlanAsistentAI key={planActiv.id}
                    intentie={planActiv.intentie} date={date} rezultat={planActiv.rezultatPlan} blocat={blocat} onOcupat={setOcupat}
                    onSchimba={intentie => seteazaConversatie(prev => prev.map(m => m.id === planActiv.id ? {...m, intentie} : m))}
                    onAplicat={rezultatPlan => seteazaConversatie(prev => prev.map(m => m.id === planActiv.id ? {...m, rezultatPlan} : m))}
                    onDeschidePontaj={onDeschidePontaj}
                />}
            </div>

            <div className="zona-conversatie-ai" hidden={vizualizare !== 'conversatie'}>

            {/* Zona de conversație */}
            <div className="flux-conversatie-ai" ref={fluxRef} onScroll={() => {
                const flux = fluxRef.current;
                if (flux) urmaresteCapatuRef.current = flux.scrollHeight - flux.scrollTop - flux.clientHeight < 48;
            }}>
                {stare && !stare.instalat && <div className="card-lipsa-fisiere-ai">
                    <h3>📁 {t('ai.lipsa_fisiere_titlu')}</h3>
                    <p>{t('ai.lipsa_fisiere_detalii')}</p>
                    <div className="cutie-cale-folder"><code>{stare.cale_folder}</code></div>
                    <ul className="instructiuni-fisiere-ai">
                        <li><strong>{stare.cale_binar ? '✓' : '1.'} llama-server.exe</strong> ({t('ai.binar_descriere')})</li>
                        <li><strong>{stare.cale_model ? '✓' : '2.'} model.gguf</strong> ({t('ai.model_recomandat')})</li>
                    </ul>
                </div>}
                {conversatie.map(m => (
                    <div key={m.id} className={`bula-mesaj ${m.rol}`}>
                        <div className="expeditor-mesaj">
                            {t(m.rol === 'user' ? 'ai.tu' : 'ai.expeditor_asistent')}
                        </div>
                        <div className="continut-mesaj">{m.text}</div>
                        {m.articoleAjutor?.map(id => <button type="button" key={id} className="legatura" onClick={() => onAjutor?.(id, m.versiuneAjutor)}>{t('ajutor.citire', {id})}</button>)}
                        {m.intentie && <div className="referinta-plan-ai">
                            <span>{m.rezultatPlan ? t('ai.plan_aplicate', {numar: m.rezultatPlan.aplicate}) : t('ai.plan_titlu')}</span>
                            <button type="button" className="buton secundar mic" disabled={blocat} onClick={() => {setPlanSelectat(m.id); setVizualizare('plan');}}>{t('ai.deschide_plan')}</button>
                        </div>}

                        {/* Carduri de acțiune propusă */}
                        {m.propuneri && m.propuneri.length > 0 && (
                            <div className="lista-propuneri-ai">
                                {m.propuneri.map(p => {
                                    const esteAplicata = m.propuneriAplicata?.[p.id];
                                    return (
                                        <div key={p.id} className={`card-propunere ${esteAplicata ? 'aplicata' : ''}`}>
                                            <div className="antet-propunere">
                                                <span className="simbol-actiune">⚡</span>
                                                <strong>{p.descriere}</strong>
                                            </div>
                                            {m.reaplicari?.[p.id] && !esteAplicata && <p className="reaplicare-propunere">↻ {t('ai.deja_aplicata_reaplica')}</p>}
                                            {p.parametri?.interpretare && <p className="interpretare-propunere">
                                                {t('ai.am_inteles')} <strong>{String(p.parametri.interpretare)}</strong>
                                                {Array.isArray(p.parametri?.corecturi) && p.parametri.corecturi.length > 0 &&
                                                    <span className="corecturi-propunere"> · {t('ai.corectat')} {(p.parametri.corecturi as string[]).join(', ')}</span>}
                                            </p>}
                                            <div className="detalii-propunere">
                                                {p.parametri?.cod_schimb && (
                                                    <span className="insigna-detaliu schimb">
                                                        📋 {t('ai.eticheta_program')} <strong>{String(p.parametri.cod_schimb)}</strong>
                                                    </span>
                                                )}
                                                {p.parametri?.marca && (
                                                    <span className="insigna-detaliu angajat">
                                                        👤 {(() => {
                                                            const ang = gasesteAngajat(String(p.parametri.marca), date.angajati);
                                                            if (ang) {
                                                                return `${ang.nume} ${ang.prenume} (${ang.marca})`;
                                                            }
                                                            return t('ai.eticheta_marca', {marca: String(p.parametri.marca)});
                                                        })()}
                                                    </span>
                                                )}
                                                {Array.isArray(p.parametri?.nume) && (
                                                    <details className="insigna-detaliu angajat lista-grup">
                                                        <summary>👥 {t('ai.grup_angajati', {numar: String(p.parametri.nume.length), grup: String(p.parametri.grup ?? '')})}</summary>
                                                        <ul>{(p.parametri.nume as string[]).map(n => <li key={n}>{n}</li>)}</ul>
                                                    </details>
                                                )}
                                                {p.parametri?.data && (
                                                    <span className="insigna-detaliu data">
                                                        📅 {t('ai.eticheta_data')} <strong>{String(p.parametri.data)}</strong>
                                                    </span>
                                                )}
                                                {p.parametri?.data_inceput && p.parametri?.data_sfarsit && (
                                                    <span className="insigna-detaliu data">
                                                        📅 {t('ai.eticheta_interval')} <strong>{String(p.parametri.data_inceput)} → {String(p.parametri.data_sfarsit)}</strong>
                                                        {p.parametri?.doar_lucratoare ? ' (L–V)' : ''}
                                                    </span>
                                                )}
                                            </div>
                                            <div className="actiuni-propunere">
                                                {esteAplicata ? (
                                                    <div className="grup-actiune-aplicata">
                                                        <span className="eticheta-aplicat">✓ {t('ai.actiune_aplicata')}</span>
                                                        {onDeschidePontaj && (
                                                            <button
                                                                type="button"
                                                                className="buton secundar mic buton-deschide-pontaj"
                                                                onClick={() => {
                                                                    const d = String(p.parametri?.data || p.parametri?.data_inceput || '');
                                                                    let a = anCurent;
                                                                    let l = lunaCurenta;
                                                                    if (d) {
                                                                        const parti = d.split('-');
                                                                        if (parti.length >= 2) {
                                                                            a = parseInt(parti[0], 10) || anCurent;
                                                                            l = parseInt(parti[1], 10) || lunaCurenta;
                                                                        }
                                                                    }
                                                                    onDeschidePontaj(a, l);
                                                                }}
                                                            >
                                                                👁 {t('ai.deschide_in_pontaj')}
                                                            </button>
                                                        )}
                                                    </div>
                                                ) : p.tip === 'deschide_concediu' ? (
                                                    <button
                                                        type="button"
                                                        className="buton primar mic"
                                                        disabled={blocat || !onDeschideConcediu}
                                                        onClick={() => onDeschideConcediu?.({
                                                            angajatId: String(p.parametri?.angajat_id ?? ''),
                                                            cod: String(p.parametri?.cod ?? ''),
                                                            deLa: String(p.parametri?.data_inceput ?? ''),
                                                            panaLa: String(p.parametri?.data_sfarsit ?? ''),
                                                        })}
                                                    >
                                                        📝 {t('ai.deschide_formular_concediu')}
                                                    </button>
                                                ) : (
                                                    <button
                                                        type="button"
                                                        className="buton primar mic"
                                                        disabled={blocat}
                                                        onClick={() => p.tip.startsWith('sterge_') ? void confirmaStergere(m.id, p) : revedePropunere(p)}
                                                    >
                                                        {t(p.tip.startsWith('sterge_') ? 'ai.confirma_actiune' : 'ai.plan_revede')}
                                                    </button>
                                                )}
                                            </div>
                                        </div>
                                    );
                                })}
                            </div>
                        )}
                    </div>
                ))}
                {lucreaza && (
                    <div className="bula-mesaj assistant scriere">
                        <span className="puncte-asteptare">● ● ●</span>
                        <span className="text-discret">{t('ai.sarcina_in_lucru')}</span>
                        <button type="button" className="buton secundar mic" onClick={() => void opresteSarcina()}>⏹ {t('ai.opreste_sarcina')}</button>
                    </div>
                )}
            </div>

            {/* Sugestii rapide de comenzi */}
            {conversatie.length <= 1 && <div className="sugestii-comenzi-ai">
                <span className="eticheta-sugestii">{t('ai.sugestii_titlu')}</span>
                {[t('ai.sugestie_1'), t('ai.sugestie_2'), t('ai.sugestie_3')].map((s, idx) => (
                    <button
                        key={idx}
                        type="button"
                        className="cip-sugestie"
                        disabled={blocat}
                        onClick={() => trimiteMesaj(s)}
                    >
                        💬 {s}
                    </button>
                ))}
            </div>}

            {/* Zona de introducere text */}
            <form
                className="formular-intrare-ai"
                onSubmit={(e: FormEvent) => {
                    e.preventDefault();
                    void trimiteMesaj();
                }}
            >
                <input
                    type="text"
                    value={textIntrare}
					maxLength={4000}
                    disabled={blocat}
                    onChange={e => setTextIntrare(e.target.value)}
                    placeholder={t('ai.indiciu_scriere')}
                />
                <button
                    type="submit"
                    className="buton primar"
                    disabled={blocat || !textIntrare.trim()}
                >
                    {t('ai.trimite')} ➔
                </button>
            </form>
            </div>
        </div>
    );
}

