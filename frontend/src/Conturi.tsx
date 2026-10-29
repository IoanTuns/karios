// Conturile aplicației (F15, D008 punctul 3, D045): primul administrator, autentificarea,
// recuperarea parolei, „Contul meu” și administrarea conturilor.
import {createContext, useContext, useEffect, useState, type FormEvent, type ReactNode} from 'react';
import {useTranslation} from 'react-i18next';
import {Autentifica, CreeazaCont, CreeazaPrimulAdministrator, IesireDinCont, IntrebariRecuperare, ListaConturi, ModificaCont, RecupereazaCont, ReseteazaParolaCont, SchimbaIntrebarileMele, SchimbaParolaMea, StareCont} from '../wailsjs/go/main/App';
import {conturi, main} from '../wailsjs/go/models';
import {Dialog} from './Formulare';
import {notifica} from './Toast';

interface StareContext {
    cont: conturi.Cont;
    intrebari: string[];
    iesire: () => Promise<void>;
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

function CampParola({eticheta, valoare, schimba, ocupat, nou}: {eticheta: string; valoare: string; schimba: (v: string) => void; ocupat: boolean; nou?: boolean}) {
    return <label>{eticheta}
        <input type="password" value={valoare} disabled={ocupat} autoComplete={nou ? 'new-password' : 'current-password'} onChange={e => schimba(e.target.value)}/>
    </label>;
}

function parolaValida(p: string): boolean {
    return p.length >= 8 && /\p{L}/u.test(p) && /\p{N}/u.test(p);
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

    async function trimite(e: FormEvent) {
        e.preventDefault();
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
        <CampParola eticheta={t('cont.parola')} valoare={parola} schimba={setParola} ocupat={ocupat} nou/>
        <CampParola eticheta={t('cont.confirma_parola')} valoare={confirmare} schimba={setConfirmare} ocupat={ocupat} nou/>
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
        <label>{t('cont.nume')}<input value={nume} disabled={ocupat} autoComplete="username" autoFocus={!ultimulNume} onChange={e => setNume(e.target.value)} required/></label>
        <label>{t('cont.parola')}<input type="password" value={parola} disabled={ocupat} autoComplete="current-password" autoFocus={!!ultimulNume} onChange={e => setParola(e.target.value)} required/></label>
        {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}
        <button className="buton primar" disabled={ocupat || !nume.trim() || !parola}>{t('cont.intra')}</button>
        <button type="button" className="legatura" disabled={ocupat} onClick={incepeRecuperarea}>{t('cont.am_uitat_parola')}</button>
    </form></div>;
}

// PoartaCont arată aplicația doar după autentificare (D045 punctul 4).
export function PoartaCont({children}: {children: ReactNode}) {
    const {t} = useTranslation();
    const [stare, setStare] = useState<main.RaspunsStareCont | null>(null);
    const [eroare, setEroare] = useState('');

    useEffect(() => {
        StareCont().then(setStare).catch(() => setEroare(t('date.eroare')));
    }, [t]);

    if (eroare) return <div className="ecran-cont"><div className="mesaj eroare" role="alert">{eroare}</div></div>;
    if (!stare) return <div className="ecran-cont"><p role="status">{t('ui.incarcare')}</p></div>;
    if (stare.necesita_cont) return <EcranPrimulAdministrator intrebariToate={stare.intrebari} gata={setStare}/>;
    if (!stare.autentificat || !stare.cont) return <EcranAutentificare ultimulNume={stare.ultimul_nume} gata={setStare}/>;

    const valoare: StareContext = {
        cont: stare.cont,
        intrebari: stare.intrebari,
        iesire: async () => { setStare(await IesireDinCont()); },
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
