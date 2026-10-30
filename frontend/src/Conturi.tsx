// Conturile aplicației (F15, D008 punctul 3, D045): primul administrator, autentificarea,
// recuperarea parolei, „Contul meu” și administrarea conturilor.
import {createContext, useContext, useEffect, useLayoutEffect, useRef, useState, type FormEvent, type ReactNode} from 'react';
import {useTranslation} from 'react-i18next';
import {parolaValida} from './parola';
import {Autentifica, ConfirmaCheiaRecuperare, CreeazaCont, CreeazaPrimulAdministrator, DeblocheazaCuParola, DeblocheazaCuRecuperare, GenereazaCheieRecuperareNoua, IesireDinCont, IntrebariRecuperare, ListaConturi, ModificaCont, RecupereazaCont, ReseteazaParolaCont, SchimbaIntrebarileMele, SchimbaParolaMea, StareCont, StareRecuperare, StareSeif} from '../wailsjs/go/main/App';
import {conturi, main} from '../wailsjs/go/models';
import {Dialog} from './Formulare';
import {notifica} from './Toast';

interface StareContext {
    cont: conturi.Cont;
    intrebari: string[];
    iesire: () => Promise<void>;
    // cheieNoua: administratorul înlocuiește cheia de recuperare; poarta cere apoi confirmarea ei.
    cheieNoua: () => Promise<void>;
}

const ContextCont = createContext<StareContext | null>(null);

export function useCont(): StareContext | null {
    return useContext(ContextCont);
}

function textEroare(t: (k: string, p?: Record<string, string>) => string, e?: main.MesajEroare): string {
    if (!e) return '';
    const p = {...(e.param ?? {})};
    if (p.pana_la) p.pana_la = new Date(p.pana_la).toLocaleTimeString();
    return t(e.cod, p);
}

// FormularIntrebari alege 3 întrebări diferite din listă și răspunsurile lor (D045 punctul 1).
function FormularIntrebari({toate, intrebari, raspunsuri, schimba, ocupat}: {
    toate: string[]; intrebari: string[]; raspunsuri: string[]; ocupat: boolean;
    schimba: (intrebari: string[], raspunsuri: string[]) => void;
}) {
    const {t} = useTranslation();
    return <fieldset className="intrebari-cont">
        <legend>{t('cont.intrebari_titlu')}</legend>
        {[0, 1, 2].map(i => <div key={i} className="rand-intrebare">
            <select aria-label={t('cont.intrebarea', {numar: String(i + 1)})} value={intrebari[i] ?? ''} disabled={ocupat}
                onChange={e => { const q = [...intrebari]; q[i] = e.target.value; schimba(q, raspunsuri); }}>
                <option value="">{t('cont.alege_intrebarea', {numar: String(i + 1)})}</option>
                {toate.map(q => <option key={q} value={q} disabled={intrebari.includes(q) && intrebari[i] !== q}>{t(`cont.${q}`)}</option>)}
            </select>
            <input type="text" aria-label={t('cont.raspunsul', {numar: String(i + 1)})} placeholder={t('cont.raspunsul', {numar: String(i + 1)})}
                value={raspunsuri[i] ?? ''} disabled={ocupat} autoComplete="off"
                onChange={e => { const r = [...raspunsuri]; r[i] = e.target.value; schimba(intrebari, r); }}/>
        </div>)}
    </fieldset>;
}

// CampParola: onInput prinde și valorile scrise de completarea automată a WebView2, care nu
// declanșează mereu onChange; butonul arată parola, ca utilizatorul să verifice ce a scris.
function CampParola({eticheta, valoare, schimba, ocupat, nou, nume}: {eticheta: string; valoare: string; schimba: (v: string) => void; ocupat: boolean; nou?: boolean; nume?: string}) {
    const {t} = useTranslation();
    const [vizibila, setVizibila] = useState(false);
    return <label>{eticheta}
        <span className="camp-parola">
            <input type={vizibila ? 'text' : 'password'} name={nume} value={valoare} disabled={ocupat} autoComplete={nou ? 'new-password' : 'current-password'}
                onChange={e => schimba(e.target.value)} onInput={e => schimba((e.target as HTMLInputElement).value)}/>
            <button type="button" className="buton-text" aria-pressed={vizibila} disabled={ocupat} onClick={() => setVizibila(v => !v)}>
                {t(vizibila ? 'cont.ascunde_parola' : 'cont.arata_parola')}
            </button>
        </span>
    </label>;
}

function EcranPrimulAdministrator({intrebariToate, gata}: {intrebariToate: string[]; gata: (r: main.RaspunsStareCont) => void}) {
    const {t} = useTranslation();
    const [nume, setNume] = useState('');
    const [numeAfisat, setNumeAfisat] = useState('');
    const [parola, setParola] = useState('');
    const [confirmare, setConfirmare] = useState('');
    const [intrebari, setIntrebari] = useState<string[]>(['', '', '']);
    const [raspunsuri, setRaspunsuri] = useState<string[]>(['', '', '']);
    const [ocupat, setOcupat] = useState(false);
    const [eroare, setEroare] = useState('');

    async function trimite(e: FormEvent<HTMLFormElement>) {
        e.preventDefault();
        // Valorile se citesc din formular: completarea automată poate ocoli starea React.
        const f = new FormData(e.currentTarget);
        const parola = String(f.get('parola') ?? ''), confirmare = String(f.get('confirmare_parola') ?? '');
        setParola(parola); setConfirmare(confirmare);
        if (!parolaValida(parola)) { setEroare(t('cont.parola_slaba')); return; }
        if (parola !== confirmare) { setEroare(t('cont.parole_diferite')); return; }
        setOcupat(true); setEroare('');
        try {
            const r = await CreeazaPrimulAdministrator(new conturi.CerereCont({nume, nume_afisat: numeAfisat, rol: 'ADMINISTRATOR', parola, intrebari, raspunsuri}));
            if (r.eroare) setEroare(textEroare(t, r.eroare)); else gata(r);
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setOcupat(false);
        }
    }

    return <div className="ecran-cont"><form className="card-cont" onSubmit={e => void trimite(e)}>
        <h1>{t('cont.primul_titlu')}</h1>
        <p className="text-discret">{t('cont.primul_descriere')}</p>
        <label>{t('cont.nume')}<input value={nume} disabled={ocupat} autoComplete="username" onChange={e => setNume(e.target.value)} required/></label>
        <label>{t('cont.nume_afisat')}<input value={numeAfisat} disabled={ocupat} onChange={e => setNumeAfisat(e.target.value)}/></label>
        <CampParola eticheta={t('cont.parola')} nume="parola" valoare={parola} schimba={setParola} ocupat={ocupat} nou/>
        <CampParola eticheta={t('cont.confirma_parola')} nume="confirmare_parola" valoare={confirmare} schimba={setConfirmare} ocupat={ocupat} nou/>
        <p className="text-discret">{t('cont.regula_parola')}</p>
        <FormularIntrebari toate={intrebariToate} intrebari={intrebari} raspunsuri={raspunsuri} ocupat={ocupat} schimba={(q, r) => { setIntrebari(q); setRaspunsuri(r); }}/>
        {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}
        <button className="buton primar" disabled={ocupat || !nume.trim()}>{t(ocupat ? 'ui.salvare' : 'cont.creeaza_administrator')}</button>
    </form></div>;
}

function EcranAutentificare({ultimulNume, gata}: {ultimulNume: string; gata: (r: main.RaspunsStareCont) => void}) {
    const {t} = useTranslation();
    const [nume, setNume] = useState(ultimulNume);
    const [parola, setParola] = useState('');
    const [recuperare, setRecuperare] = useState<string[] | null>(null);
    const [raspunsuri, setRaspunsuri] = useState<string[]>(['', '', '']);
    const [parolaNoua, setParolaNoua] = useState('');
    const [ocupat, setOcupat] = useState(false);
    const [eroare, setEroare] = useState('');

    async function ruleaza(f: () => Promise<void>) {
        setOcupat(true); setEroare('');
        try { await f(); } catch { setEroare(t('date.eroare')); } finally { setOcupat(false); }
    }

    function autentifica(e: FormEvent) {
        e.preventDefault();
        void ruleaza(async () => {
            const r = await Autentifica(nume, parola);
            if (r.eroare) { setEroare(textEroare(t, r.eroare)); setParola(''); } else gata(r);
        });
    }

    function incepeRecuperarea() {
        if (!nume.trim()) { setEroare(t('cont.scrie_numele')); return; }
        void ruleaza(async () => {
            const r = await IntrebariRecuperare(nume);
            if (r.eroare) setEroare(textEroare(t, r.eroare)); else { setRecuperare(r.intrebari); setRaspunsuri(['', '', '']); setParolaNoua(''); }
        });
    }

    function recupereaza(e: FormEvent) {
        e.preventDefault();
        if (!parolaValida(parolaNoua)) { setEroare(t('cont.parola_slaba')); return; }
        void ruleaza(async () => {
            const r = await RecupereazaCont(nume, raspunsuri, parolaNoua);
            if (r.eroare) setEroare(textEroare(t, r.eroare)); else gata(r);
        });
    }

    if (recuperare) {
        return <div className="ecran-cont"><form className="card-cont" onSubmit={recupereaza}>
            <h1>{t('cont.recuperare_titlu')}</h1>
            <p className="text-discret">{t('cont.recuperare_descriere', {nume})}</p>
            {recuperare.map((q, i) => <label key={q}>{t(`cont.${q}`)}
                <input type="text" value={raspunsuri[i]} disabled={ocupat} autoComplete="off" onChange={e => { const r = [...raspunsuri]; r[i] = e.target.value; setRaspunsuri(r); }}/>
            </label>)}
            <CampParola eticheta={t('cont.parola_noua')} valoare={parolaNoua} schimba={setParolaNoua} ocupat={ocupat} nou/>
            <p className="text-discret">{t('cont.regula_parola')}</p>
            {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}
            <div className="actiuni-cont">
                <button type="button" className="buton secundar" disabled={ocupat} onClick={() => { setRecuperare(null); setEroare(''); }}>{t('ui.anuleaza')}</button>
                <button className="buton primar" disabled={ocupat}>{t('cont.schimba_parola')}</button>
            </div>
            <p className="text-discret">{t('cont.recuperare_administrator')}</p>
        </form></div>;
    }

    return <div className="ecran-cont"><form className="card-cont" onSubmit={autentifica}>
        <h1>{t('cont.autentificare_titlu')}</h1>
        <label>{t('cont.nume')}<input value={nume} disabled={ocupat} autoComplete="username" autoFocus={!ultimulNume} onChange={e => setNume(e.target.value)} onInput={e => setNume((e.target as HTMLInputElement).value)} required/></label>
        <label>{t('cont.parola')}<input type="password" value={parola} disabled={ocupat} autoComplete="current-password" autoFocus={!!ultimulNume} onChange={e => setParola(e.target.value)} onInput={e => setParola((e.target as HTMLInputElement).value)} required/></label>
        {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}
        <button className="buton primar" disabled={ocupat || !nume.trim() || !parola}>{t('cont.intra')}</button>
        <button type="button" className="legatura" disabled={ocupat} onClick={incepeRecuperarea}>{t('cont.am_uitat_parola')}</button>
    </form></div>;
}

function CampCheieRecuperare({valoare, schimba, ocupat, confirma = false}: {
    valoare: string; schimba: (valoare: string) => void; ocupat: boolean; confirma?: boolean;
}) {
    const camp = useRef<HTMLInputElement>(null);
    const selectie = useRef<{inceput: number; sfarsit: number} | null>(null);
    const [, actualizeazaSelectia] = useState(0);
    const formateaza = (text: string) => text.toUpperCase().replace(/[^A-Z0-9]/g, '').match(/.{1,5}/g)?.join('-') ?? '';
    useLayoutEffect(() => {
        if (selectie.current && camp.current === document.activeElement) {
            camp.current?.setSelectionRange(selectie.current.inceput, selectie.current.sfarsit);
        }
        selectie.current = null;
    });
    return <input ref={camp} type="text" value={valoare} disabled={ocupat} autoComplete="off" spellCheck={false}
        onPaste={confirma ? e => e.preventDefault() : undefined}
        onChange={e => {
            const text = e.currentTarget.value;
            const pozitie = (index: number) => formateaza(text.slice(0, index)).length;
            selectie.current = {inceput: pozitie(e.currentTarget.selectionStart ?? text.length), sfarsit: pozitie(e.currentTarget.selectionEnd ?? text.length)};
            schimba(formateaza(text));
            // Refacem selecția și când o cratimă introdusă manual nu schimbă valoarea.
            actualizeazaSelectia(n => n + 1);
        }}/>;
}

// EcranDeblocare: seiful nu are copia acestui utilizator Windows (alt calculator, alt cont
// Windows sau Windows reinstalat). Datele se deschid cu parola unui cont sau cu cheia de
// recuperare (D054); după aceea, pornirile următoare nu mai cer nimic.
function EcranDeblocare({gata}: {gata: (r: main.RaspunsStareCont) => void}) {
    const {t} = useTranslation();
    const [nume, setNume] = useState('');
    const [parola, setParola] = useState('');
    const [cuCheie, setCuCheie] = useState(false);
    const [cheie, setCheie] = useState('');
    const [ocupat, setOcupat] = useState(false);
    const [eroare, setEroare] = useState('');

    function trimite(e: FormEvent) {
        e.preventDefault();
        setOcupat(true); setEroare('');
        (cuCheie ? DeblocheazaCuRecuperare(cheie) : DeblocheazaCuParola(nume, parola))
            .then(r => { if (r.eroare) { setEroare(textEroare(t, r.eroare)); setParola(''); } else gata(r); })
            .catch(() => setEroare(t('date.eroare')))
            .finally(() => setOcupat(false));
    }

    return <div className="ecran-cont"><form className="card-cont" onSubmit={trimite}>
        <h1>{t('seif.deblocare_titlu')}</h1>
        <p className="text-discret">{t('seif.deblocare_descriere')}</p>
        {cuCheie
            ? <label>{t('seif.cheie_recuperare')}
                <CampCheieRecuperare valoare={cheie} schimba={setCheie} ocupat={ocupat}/>
            </label>
            : <>
                <label>{t('cont.nume')}<input type="text" value={nume} disabled={ocupat} autoComplete="username" onChange={e => setNume(e.target.value)}/></label>
                <CampParola eticheta={t('cont.parola')} valoare={parola} schimba={setParola} ocupat={ocupat}/>
            </>}
        {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}
        <div className="actiuni-cont">
            <button type="button" className="buton secundar" disabled={ocupat} onClick={() => { setCuCheie(v => !v); setEroare(''); }}>
                {t(cuCheie ? 'seif.cu_parola' : 'seif.cu_cheia')}
            </button>
            <button className="buton primar" disabled={ocupat || (cuCheie ? !cheie.trim() : !nume.trim() || !parola)}>{t('seif.deblocheaza')}</button>
        </div>
    </form></div>;
}

// EcranCheieRecuperare: cheia de recuperare e obligatorie (D054 punctul 3). Administratorul o
// notează și o retastează; ceilalți așteaptă ca un administrator să o confirme.
function EcranCheieRecuperare({stare, gata, iesire}: {stare: main.RaspunsRecuperare; gata: (r: main.RaspunsRecuperare) => void; iesire: () => void}) {
    const {t} = useTranslation();
    const [cheie, setCheie] = useState('');
    const [ocupat, setOcupat] = useState(false);
    const [eroare, setEroare] = useState('');

    if (!stare.cod) {
        return <div className="ecran-cont"><div className="card-cont">
            <h1>{t('seif.recuperare_titlu')}</h1>
            <p>{t('seif.recuperare_asteapta_admin')}</p>
            <div className="actiuni-cont"><button type="button" className="buton secundar" onClick={iesire}>{t('cont.iesire')}</button></div>
        </div></div>;
    }

    function confirma(e: FormEvent) {
        e.preventDefault();
        setOcupat(true); setEroare('');
        ConfirmaCheiaRecuperare(cheie)
            .then(r => { if (r.eroare) setEroare(textEroare(t, r.eroare)); else gata(r); })
            .catch(() => setEroare(t('date.eroare')))
            .finally(() => setOcupat(false));
    }

    // Lipirea e oprită în câmpul de confirmare: cheia trebuie copiată de pe hârtie, nu din ecran.
    return <div className="ecran-cont"><form className="card-cont" onSubmit={confirma}>
        <h1>{t('seif.recuperare_titlu')}</h1>
        <p>{t('seif.recuperare_descriere')}</p>
        <p className="cheie-recuperare" aria-label={t('seif.cheie_recuperare')}><code>{stare.cod}</code></p>
        <p className="text-discret">{t('seif.recuperare_pastrare')}</p>
        <label>{t('seif.recuperare_retasteaza')}
            <CampCheieRecuperare valoare={cheie} schimba={setCheie} ocupat={ocupat} confirma/>
        </label>
        {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}
        <div className="actiuni-cont">
            <button className="buton primar" disabled={ocupat || !cheie.trim()}>{t('seif.recuperare_confirma')}</button>
        </div>
    </form></div>;
}

// PoartaCont arată aplicația doar după deblocarea datelor (D054), autentificare (D045 punctul 4)
// și confirmarea cheii de recuperare.
export function PoartaCont({children}: {children: ReactNode}) {
    const {t} = useTranslation();
    const [blocat, setBlocat] = useState<boolean | null>(null);
    const [stare, setStare] = useState<main.RaspunsStareCont | null>(null);
    const [recuperare, setRecuperare] = useState<main.RaspunsRecuperare | null>(null);
    const [eroare, setEroare] = useState('');

    useEffect(() => {
        StareSeif().then(s => {
            setBlocat(s.blocat);
            if (!s.blocat) return StareCont().then(setStare);
        }).catch(() => setEroare(t('date.eroare')));
    }, [t]);

    const idCont = stare?.autentificat ? stare.cont?.id : undefined;
    useEffect(() => {
        setRecuperare(null);
        if (!idCont) return;
        StareRecuperare().then(setRecuperare).catch(() => setEroare(t('date.eroare')));
    }, [idCont, t]);

    if (eroare) return <div className="ecran-cont"><div className="mesaj eroare" role="alert">{eroare}</div></div>;
    if (blocat) return <EcranDeblocare gata={r => { setBlocat(false); setStare(r); }}/>;
    if (!stare) return <div className="ecran-cont"><p role="status">{t('ui.incarcare')}</p></div>;
    if (stare.necesita_cont) return <EcranPrimulAdministrator intrebariToate={stare.intrebari} gata={setStare}/>;
    if (!stare.autentificat || !stare.cont) return <EcranAutentificare ultimulNume={stare.ultimul_nume} gata={setStare}/>;
    if (!recuperare) return <div className="ecran-cont"><p role="status">{t('ui.incarcare')}</p></div>;
    if (recuperare.necesara) return <EcranCheieRecuperare stare={recuperare} gata={setRecuperare} iesire={() => void IesireDinCont().then(setStare)}/>;

    const valoare: StareContext = {
        cont: stare.cont,
        intrebari: stare.intrebari,
        iesire: async () => { setStare(await IesireDinCont()); },
        cheieNoua: async () => {
            const r = await GenereazaCheieRecuperareNoua();
            if (r.eroare) notifica(textEroare(t, r.eroare)); else setRecuperare(r);
        },
    };
    // Cheia remontează aplicația la schimbarea contului, ca datele să se reîncarce.
    return <ContextCont.Provider value={valoare}><div key={stare.cont.id} className="radacina-cont">{children}</div></ContextCont.Provider>;
}

export function DialogContulMeu({inchide}: {inchide: () => void}) {
    const {t} = useTranslation();
    const ctx = useCont();
    const [veche, setVeche] = useState('');
    const [noua, setNoua] = useState('');
    const [parolaIntrebari, setParolaIntrebari] = useState('');
    const [intrebari, setIntrebari] = useState<string[]>(ctx?.cont.intrebari ?? ['', '', '']);
    const [raspunsuri, setRaspunsuri] = useState<string[]>(['', '', '']);
    const [ocupat, setOcupat] = useState(false);
    const [eroare, setEroare] = useState('');
    if (!ctx) return null;

    async function ruleaza(f: () => Promise<main.RaspunsOperatieCont>, succes: string) {
        setOcupat(true); setEroare('');
        try {
            const r = await f();
            if (r.eroare) setEroare(textEroare(t, r.eroare)); else { notifica(t(succes)); inchide(); }
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setOcupat(false);
        }
    }

    return <Dialog titlu={t('cont.contul_meu')} ocupat={ocupat} inchide={inchide}>
        <div className="corp-dialog formular-cont">
            <p><strong>{ctx.cont.nume_afisat || ctx.cont.nume}</strong> · {ctx.cont.nume} · {t(`cont.rol_${ctx.cont.rol}`)}</p>
            <h3>{t('cont.schimba_parola')}</h3>
            <CampParola eticheta={t('cont.parola_veche')} valoare={veche} schimba={setVeche} ocupat={ocupat}/>
            <CampParola eticheta={t('cont.parola_noua')} valoare={noua} schimba={setNoua} ocupat={ocupat} nou/>
            <p className="text-discret">{t('cont.regula_parola')}</p>
            <button type="button" className="buton primar" disabled={ocupat || !veche || !parolaValida(noua)} onClick={() => void ruleaza(() => SchimbaParolaMea(veche, noua), 'cont.parola_schimbata')}>{t('cont.schimba_parola')}</button>
            <h3>{t('cont.intrebari_titlu')}</h3>
            <FormularIntrebari toate={ctx.intrebari} intrebari={intrebari} raspunsuri={raspunsuri} ocupat={ocupat} schimba={(q, r) => { setIntrebari(q); setRaspunsuri(r); }}/>
            <CampParola eticheta={t('cont.parola_confirmare')} valoare={parolaIntrebari} schimba={setParolaIntrebari} ocupat={ocupat}/>
            <button type="button" className="buton secundar" disabled={ocupat || !parolaIntrebari} onClick={() => void ruleaza(() => SchimbaIntrebarileMele(parolaIntrebari, intrebari, raspunsuri), 'cont.intrebari_schimbate')}>{t('cont.salveaza_intrebarile')}</button>
            {ctx.cont.rol === 'ADMINISTRATOR' && <>
                <h3>{t('seif.recuperare_titlu')}</h3>
                <p className="text-discret">{t('seif.cheie_noua_descriere')}</p>
                <button type="button" className="buton secundar" disabled={ocupat} onClick={() => { if (window.confirm(t('seif.cheie_noua_confirmare'))) { inchide(); void ctx.cheieNoua(); } }}>{t('seif.cheie_noua')}</button>
            </>}
            {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}
        </div>
        <footer className="subsol-dialog"><button className="buton secundar" disabled={ocupat} onClick={inchide}>{t('ui.inchide')}</button></footer>
    </Dialog>;
}

// EcranConturi: administrarea conturilor (F15-T07), doar pentru administrator.
export function EcranConturi() {
    const {t} = useTranslation();
    const ctx = useCont();
    const [lista, setLista] = useState<conturi.Cont[]>([]);
    const [ocupat, setOcupat] = useState(false);
    const [eroare, setEroare] = useState('');
    const [nou, setNou] = useState<{nume: string; nume_afisat: string; rol: string; parola: string; intrebari: string[]; raspunsuri: string[]} | null>(null);
    const [resetare, setResetare] = useState<{cont: conturi.Cont; parola: string} | null>(null);

    function aplica(r: main.RaspunsConturi): boolean {
        if (r.eroare) { setEroare(textEroare(t, r.eroare)); return false; }
        setEroare(''); setLista(r.conturi ?? []); return true;
    }

    async function ruleaza(f: () => Promise<main.RaspunsConturi>, succes?: string): Promise<boolean> {
        setOcupat(true);
        try {
            const ok = aplica(await f());
            if (ok && succes) notifica(t(succes));
            return ok;
        } catch {
            setEroare(t('date.eroare'));
            return false;
        } finally {
            setOcupat(false);
        }
    }

    useEffect(() => { void ruleaza(ListaConturi); }, []);

    if (!ctx?.cont || ctx.cont.rol !== 'ADMINISTRATOR') return <section className="panou"><p>{t('cont.doar_administrator')}</p></section>;

    return <section className="panou ecran-conturi">
        <div className="actiuni-antet"><button type="button" className="buton primar" disabled={ocupat}
            onClick={() => setNou({nume: '', nume_afisat: '', rol: 'OPERATOR', parola: '', intrebari: ['', '', ''], raspunsuri: ['', '', '']})}>＋ {t('cont.adauga')}</button></div>
        {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}
        <table>
            <thead><tr><th scope="col">{t('cont.nume')}</th><th scope="col">{t('cont.nume_afisat')}</th><th scope="col">{t('cont.rol')}</th><th scope="col">{t('cont.stare')}</th><th scope="col"/></tr></thead>
            <tbody>{lista.map(c => <tr key={c.id}>
                <td>{c.nume}{c.id === ctx.cont.id ? ` (${t('cont.tu')})` : ''}</td>
                <td>{c.nume_afisat}</td>
                <td><select aria-label={t('cont.rol')} value={c.rol} disabled={ocupat}
                    onChange={e => void ruleaza(() => ModificaCont(new conturi.Cont({...c, rol: e.target.value})), 'ui.salvat')}>
                    <option value="ADMINISTRATOR">{t('cont.rol_ADMINISTRATOR')}</option><option value="OPERATOR">{t('cont.rol_OPERATOR')}</option>
                </select></td>
                <td>{!c.activ ? t('cont.inactiv') : c.blocat_pana_la && new Date(c.blocat_pana_la) > new Date() ? t('cont.blocat_pana', {ora: new Date(c.blocat_pana_la).toLocaleTimeString()}) : t('cont.activ')}</td>
                <td className="actiuni-rand">
                    <button type="button" className="buton secundar mic" disabled={ocupat} onClick={() => setResetare({cont: c, parola: ''})}>{t('cont.reseteaza_parola')}</button>
                    <button type="button" className="buton secundar mic" disabled={ocupat}
                        onClick={() => void ruleaza(() => ModificaCont(new conturi.Cont({...c, activ: !c.activ})), 'ui.salvat')}>{t(c.activ ? 'cont.dezactiveaza' : 'cont.activeaza')}</button>
                </td>
            </tr>)}</tbody>
        </table>
        {nou && <Dialog titlu={t('cont.adauga')} ocupat={ocupat} inchide={() => setNou(null)}>
            <div className="corp-dialog formular-cont">
                <label>{t('cont.nume')}<input value={nou.nume} disabled={ocupat} onChange={e => setNou({...nou, nume: e.target.value})}/></label>
                <label>{t('cont.nume_afisat')}<input value={nou.nume_afisat} disabled={ocupat} onChange={e => setNou({...nou, nume_afisat: e.target.value})}/></label>
                <label>{t('cont.rol')}<select value={nou.rol} disabled={ocupat} onChange={e => setNou({...nou, rol: e.target.value})}>
                    <option value="OPERATOR">{t('cont.rol_OPERATOR')}</option><option value="ADMINISTRATOR">{t('cont.rol_ADMINISTRATOR')}</option>
                </select></label>
                <CampParola eticheta={t('cont.parola_initiala')} valoare={nou.parola} schimba={p => setNou({...nou, parola: p})} ocupat={ocupat} nou/>
                <p className="text-discret">{t('cont.regula_parola')}</p>
                <FormularIntrebari toate={ctx.intrebari} intrebari={nou.intrebari} raspunsuri={nou.raspunsuri} ocupat={ocupat} schimba={(q, r) => setNou({...nou, intrebari: q, raspunsuri: r})}/>
                <p className="text-discret">{t('cont.intrebari_impreuna')}</p>
                {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}
            </div>
            <footer className="subsol-dialog">
                <button className="buton secundar" disabled={ocupat} onClick={() => setNou(null)}>{t('ui.anuleaza')}</button>
                <button className="buton primar" disabled={ocupat || !nou.nume.trim() || !parolaValida(nou.parola)}
                    onClick={() => void ruleaza(() => CreeazaCont(new conturi.CerereCont(nou)), 'cont.creat').then(ok => { if (ok) setNou(null); })}>{t('cont.adauga')}</button>
            </footer>
        </Dialog>}
        {resetare && <Dialog titlu={t('cont.reseteaza_parola_pentru', {nume: resetare.cont.nume})} ocupat={ocupat} inchide={() => setResetare(null)}>
            <div className="corp-dialog formular-cont">
                <CampParola eticheta={t('cont.parola_noua')} valoare={resetare.parola} schimba={p => setResetare({...resetare, parola: p})} ocupat={ocupat} nou/>
                <p className="text-discret">{t('cont.regula_parola')}</p>
                {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}
            </div>
            <footer className="subsol-dialog">
                <button className="buton secundar" disabled={ocupat} onClick={() => setResetare(null)}>{t('ui.anuleaza')}</button>
                <button className="buton primar" disabled={ocupat || !parolaValida(resetare.parola)}
                    onClick={() => void ruleaza(() => ReseteazaParolaCont(resetare.cont.id, resetare.parola), 'cont.parola_resetata').then(ok => { if (ok) setResetare(null); })}>{t('cont.reseteaza_parola')}</button>
            </footer>
        </Dialog>}
    </section>;
}
