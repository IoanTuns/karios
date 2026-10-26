import {useEffect, useRef, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {AdaugaDepartament, AdaugaSectie, DateInitiale, IncarcaDemo, Informatii, SalveazaAngajat, SalveazaNormaIntreaga, SchimbaActivAngajat} from '../wailsjs/go/main/App';
import {datebaza, main} from '../wailsjs/go/models';
import {DialogDemo, FormularAngajat, Setari, durata, denumireSectie, denumireDepartament} from './Formulare';
import {Sarbatori} from './Sarbatori';
import {Schimburi} from './Schimburi';
import {Pontaj} from './Pontaj';
import {Absente} from './Absente';
import {ReguliPontaj} from './ReguliPontaj';
import {DefinireRapoarte, GenerareRapoarte} from './Rapoarte';

type Ecran = 'angajati' | 'pontaj' | 'sarbatori' | 'schimburi' | 'absente' | 'definire_rapoarte' | 'rapoarte' | 'reguli_pontaj' | 'setari';
const ecranePrincipale: Ecran[] = ['angajati', 'pontaj', 'absente', 'rapoarte', 'setari'];
const sectiuniSetari: Ecran[] = ['setari', 'sarbatori', 'schimburi', 'reguli_pontaj', 'definire_rapoarte'];

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
            : tip === 'rapoarte' ? <><path d="M14 3H6a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/><path d="M14 3v6h6M8 13h8M8 17h5"/></>
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
    const [pontajModificat, setPontajModificat] = useState(false);
    const [pontajOcupat, setPontajOcupat] = useState(false);
    const [filtruAbsente, setFiltruAbsente] = useState<{an: number; angajatId: string} | null>(null);

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

    return <div className="aplicatie">
        <aside className="bara-laterala">
            <div className="marca-aplicatie"><span className="sigla" aria-hidden="true">p<span>·</span></span><strong>{t('aplicatie.titlu')}</strong></div>
            <p className="eticheta-meniu">{t('ui.date_baza')}</p>
            <nav aria-label={t('ui.date_baza')}>
                {ecranePrincipale.map(item => {
                    const activ = item === 'setari' ? esteInSetari : ecran === item;
                    return <div key={item} className="grup-navigare">
                        <button className={`navigare ${item === 'setari' ? 'navigare-setari' : ''} ${activ ? 'selectat' : ''}`}
                            aria-current={activ ? 'page' : undefined} disabled={ocupat || pontajOcupat}
                            onClick={() => schimbaEcran(item)}>
                            <Simbol tip={item}/>{t(`ui.${item}`)}</button>
                        {item === 'setari' && esteInSetari && <div className="submeniu-setari" role="group" aria-label={t('ui.setari')}>
                            {sectiuniSetari.map(sub => <button key={sub} type="button" className={`subnavigare ${ecran === sub ? 'selectat' : ''}`}
                                aria-current={ecran === sub ? 'page' : undefined} disabled={ocupat || pontajOcupat}
                                onClick={() => schimbaEcran(sub)}>
                                <Simbol tip={sub}/>{t(sub === 'setari' ? 'ui.setari_general' : `ui.${sub}`)}</button>)}
                        </div>}
                    </div>;
                })}
            </nav>
            <div className="subsol-meniu"><span className="punct-local"/>{t('ui.local')}
                {versiune && <span className="versiune">{t('aplicatie.versiune', {versiune})}</span>}
                <select name="limba_interfata" className="selector-limba" aria-label={t('ui.limba_interfata')} value={i18n.language === 'en' ? 'en' : 'ro'}
                    onChange={e => { const l = e.target.value; void i18n.changeLanguage(l); document.documentElement.lang = l; document.title = i18n.t('aplicatie.titlu'); }}>
                    <option value="ro">{t('rapoarte.ro')}</option><option value="en">{t('rapoarte.en')}</option>
                </select>
            </div>
        </aside>
        <main className={`continut ${ecran === 'pontaj' ? 'continut-pontaj' : ''} ${esteInSetari ? 'continut-setari' : ''}`}>
            <header className="antet">
                <div><p className="fir-navigare">{t('ui.date_baza')} <span aria-hidden="true">/</span> {esteInSetari && ecran !== 'setari' ? <>{t('ui.setari')} <span aria-hidden="true">/</span> </> : null}{t(`ui.${ecran}`)}</p>
                    <h1>{t(`ui.${ecran}`)}</h1><p className="descriere">{t(`ui.${ecran}_descriere`)}</p></div>
                {ecran === 'angajati' && date && <button className="buton primar" disabled={ocupat} onClick={() => deschide()}>
                    <span aria-hidden="true">＋</span>{t('ui.adauga_angajat')}</button>}
            </header>
            {esteInSetari && <div className="taburi-sectiuni taburi-setari" role="group" aria-label={t('ui.setari')}>
                {sectiuniSetari.map(sub => <button key={sub} type="button" className={`buton-tab ${ecran === sub ? 'activ' : ''}`}
                    aria-pressed={ecran === sub} disabled={ocupat || pontajOcupat}
                    onClick={() => schimbaEcran(sub)}>{t(sub === 'setari' ? 'ui.setari_general' : `ui.${sub}`)}</button>)}
            </div>}
            {!editare && !demo && eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}
            {mesaj && <div className="mesaj succes" role="status">{mesaj}</div>}
            {incarcare ? <div className="asteptare" role="status">{t('ui.incarcare')}</div>
                : !date ? <button className="buton secundar" disabled={ocupat} onClick={() => executa(DateInitiale, '')}>{t('ui.reincearca')}</button>
                : ecran === 'pontaj' ? <Pontaj date={date} onModificari={setPontajModificat} onOcupat={setPontajOcupat} onDeschideAbsente={(an, angajatId) => { setFiltruAbsente({an, angajatId}); setPontajModificat(false); setEcran('absente'); setEroare(''); setMesaj(''); }}/>
                : ecran === 'sarbatori' ? <Sarbatori/>
                : ecran === 'schimburi' ? <Schimburi/>
                : ecran === 'absente' ? <Absente date={date} filtruInitial={filtruAbsente}/>
                : ecran === 'reguli_pontaj' ? <ReguliPontaj onModificari={setPontajModificat} onOcupat={setPontajOcupat}/>
                : ecran === 'definire_rapoarte' ? <DefinireRapoarte/>
                : ecran === 'rapoarte' ? <GenerareRapoarte date={date}/>
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
                        {lista.length === 0 ? <div className="stare-goala">
                            <span className="simbol-gol" aria-hidden="true"><Simbol tip="angajati"/></span>
                            <h2>{t('ui.gol_titlu')}</h2><p>{t('ui.gol_descriere')}</p>
                            {date.demo_disponibil && <button className="buton secundar" disabled={ocupat} onClick={() => { setEroare(''); setDemo(true); }}>{t('ui.demo')}</button>}
                        </div> : <div className="tabel-derulant"><table>
                            <thead><tr>{['angajat', 'marca', 'sectie', 'departament', 'contract', 'norma', 'concediu', 'stare', 'actiuni'].map(c => <th key={c} scope="col">{t(`ui.${c}`)}</th>)}</tr></thead>
                            <tbody>{vizibili.map(a => <tr key={a.id} className={a.activ ? '' : 'rand-inactiv'}>
                                <td><div className="identitate"><span className="avatar" aria-hidden="true">{a.prenume.charAt(0)}{a.nume.charAt(0)}</span>
                                    <span><strong>{a.nume} {a.prenume}</strong>{a.functie && <small>{a.functie}</small>}</span></div></td>
                                <td className="marca-angajat">{a.marca}</td><td>{sectiePentru(a.sectie_id)}</td><td>{departamentPentru(a.departament_id)}</td>
                                <td><span className={`eticheta-contract ${a.tip_contract === 'PARTIALA' ? 'partial' : ''}`}>{t(a.tip_contract === 'PARTIALA' ? 'ui.partiala' : 'ui.intreaga')}</span></td>
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
        </main>
        {editare && date && <FormularAngajat key={editare.id || 'nou'} angajat={editare} date={date} ocupat={ocupat} eroare={eroare}
            inchide={() => { setEditare(null); setEroare(''); }}
            salveaza={async a => { if (await executa(() => SalveazaAngajat(a))) setEditare(null); }}/>}
        {demo && <DialogDemo ocupat={ocupat} eroare={eroare} inchide={() => { setDemo(false); setEroare(''); }}
            confirma={async () => { if (await executa(IncarcaDemo, 'ui.demo_salvat')) setDemo(false); }}/>}
    </div>;
}

export default App;
