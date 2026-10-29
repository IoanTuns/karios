import {useEffect, useRef, useState, type FormEvent} from 'react';
import {useTranslation} from 'react-i18next';
import {AplicaPropunereAI, EsteAutoPornireAI, OpresteAI, PornesteAI, SeteazaAutoPornireAI, StareAI} from '../wailsjs/go/main/App';
import {ai, datebaza, main} from '../wailsjs/go/models';
import {marcheazaVizibil, opresteSarcina, pornesteSarcina, reseteazaConversatie as reseteazaSesiunea, seteazaConversatie, seteazaEroare, useSesiuneAI, type MesajUI} from './asistentSesiune';

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
    onDeschidePontaj,
}: {
    date: datebaza.Date;
    onDeschidePontaj?: (an: number, luna: number) => void;
}) {
    const {t} = useTranslation();
    const [stare, setStare] = useState<ai.StareMediu | null>(null);
    const [incarcare, setIncarcare] = useState(true);
    const [ocupat, setOcupat] = useState(false);
    const [mesaj, setMesaj] = useState('');
    const [eroare, setEroare] = useState('');
    const [textIntrare, setTextIntrare] = useState('');
    const sesiune = useSesiuneAI();
    const conversatie = sesiune.conversatie;
    const lucreaza = sesiune.sarcina !== '';
    const blocat = ocupat || lucreaza;
    const [autoPornire, setAutoPornire] = useState(true);

    const anCurent = new Date().getFullYear();
    const lunaCurenta = new Date().getMonth() + 1;

    const capatConversatieRef = useRef<HTMLDivElement>(null);

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
        capatConversatieRef.current?.scrollIntoView({behavior: 'smooth'});
    }, [conversatie, lucreaza]);

    function reseteazaConversatie() {
        reseteazaSesiunea();
        setMesaj('');
        setEroare('');
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

        setEroare('');
        seteazaEroare('');

        // 1. Verificare răspuns afirmativ de confirmare a unei propuneri existente
        const textNorm = text.toLowerCase().replace(/[.!?,]/g, '').trim();
        const esteConfirmare = ['da', 'confirm', 'confirma', 'aplica', 'aplic', 'ok', 'de acord', 'executa', 'fa asa', 'este bine', 'e bine', 'da aplica', 'da te rog', 'sigur', 'salveaza'].includes(textNorm);

        if (esteConfirmare) {
            const mesajCuPropunere = [...conversatie].reverse().find(m =>
                m.rol === 'assistant' &&
                m.propuneri &&
                m.propuneri.some(p => !m.propuneriAplicata?.[p.id])
            );

            if (mesajCuPropunere && mesajCuPropunere.propuneri) {
                const propuneriDeAplicat = mesajCuPropunere.propuneri.filter(p => !mesajCuPropunere.propuneriAplicata?.[p.id]);
                if (propuneriDeAplicat.length > 0) {
                    setTextIntrare('');
                    const mesajUtilizator: MesajUI = { id: `u_${Date.now()}`, rol: 'user', text };
                    seteazaConversatie(prev => [...prev, mesajUtilizator]);
                    for (const p of propuneriDeAplicat) {
                        await aplicaPropunere(mesajCuPropunere.id, p);
                    }
                    return;
                }
            }
        }

        setTextIntrare('');

        const mesajUtilizator: MesajUI = {
            id: `u_${Date.now()}`,
            rol: 'user',
            text,
        };

        const istoricNou = [...conversatie, mesajUtilizator];
        seteazaConversatie(() => istoricNou);

        // Trimite ultimele 10 mesaje pentru context conversațional bogat
        const istoricAPI: ai.MesajConversatie[] = istoricNou
            .filter(m => m.id !== 'bun_venit')
            .slice(-10)
            .map(m => {
                let continut = m.text;
                if (m.rol === 'assistant' && m.propuneri && m.propuneri.length > 0) {
                    const detalii = m.propuneri.map(p => {
                        const stareTxt = m.propuneriAplicata?.[p.id] ? '(deja aplicată)' : '(propusă, neconfirmată)';
                        return `${p.descriere} ${stareTxt}`;
                    }).join('; ');
                    continut = continut && !continut.includes('Acțiuni propuse') ? `${continut}. ${detalii}` : detalii;
                }
                return new ai.MesajConversatie({
                    role: m.rol,
                    content: continut,
                });
            });

        await pornesteSarcina(new main.CerereInterogareAI({
            istoric: istoricAPI,
            an: anCurent,
            luna: lunaCurenta,
        }));
    }

    async function aplicaPropunere(mesajId: string, propunere: ai.PropunereActiune) {
        setOcupat(true);
        setEroare('');
        setMesaj('');

        try {
            const resp = await AplicaPropunereAI(new main.CerereAplicaPropunereAI({
                an: anCurent,
                luna: lunaCurenta,
                propunere,
            }));

            if (resp.eroare) {
                const textEroare = t(resp.eroare.cod, resp.eroare.param ?? {});
                setEroare(textEroare !== resp.eroare.cod ? textEroare : (resp.eroare.param?.detalii ?? resp.eroare.cod));
            } else {
                setMesaj(t('ai.actiune_aplicata'));
                seteazaConversatie(prev => prev.map(m => {
                    if (m.id === mesajId) {
                        return {
                            ...m,
                            propuneriAplicata: {
                                ...(m.propuneriAplicata ?? {}),
                                [propunere.id]: true,
                            },
                        };
                    }
                    return m;
                }));
            }
        } catch (err: any) {
            setEroare(err?.message || t('date.eroare'));
        } finally {
            setOcupat(false);
        }
    }

    return (
        <div className="panou-asistent-ai">
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

            {/* Avertisment ghidat dacă fișierele lipsesc */}
            {stare && !stare.instalat && (
                <div className="card-lipsa-fisiere-ai">
                    <h3>📁 {t('ai.lipsa_fisiere_titlu')}</h3>
                    <p>{t('ai.lipsa_fisiere_detalii')}</p>
                    <div className="cutie-cale-folder">
                        <code>{stare.cale_folder}</code>
                    </div>
                    <ul className="instructiuni-fisiere-ai">
                        <li>
                            <strong>{stare.cale_binar ? '✓' : '1.'} llama-server.exe</strong> (Binarul nativ portabil)
                        </li>
                        <li>
                            <strong>{stare.cale_model ? '✓' : '2.'} model.gguf</strong> (Recomandat: Qwen 2.5 3B Instruct Q4_K_M)
                        </li>
                    </ul>
                </div>
            )}

            {(eroare || sesiune.eroare) && <div className="mesaj eroare" role="alert">{eroare || sesiune.eroare}</div>}
            {mesaj && <div className="mesaj succes" role="status">{mesaj}</div>}

            {/* Zona de conversație */}
            <div className="flux-conversatie-ai">
                {conversatie.map(m => (
                    <div key={m.id} className={`bula-mesaj ${m.rol}`}>
                        <div className="expeditor-mesaj">
                            {m.rol === 'user' ? 'Tu' : 'Pontix AI'}
                        </div>
                        <div className="continut-mesaj">{m.text}</div>

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
                                            <div className="detalii-propunere">
                                                {p.parametri?.cod_schimb && (
                                                    <span className="insigna-detaliu schimb">
                                                        📋 Program: <strong>{String(p.parametri.cod_schimb)}</strong>
                                                    </span>
                                                )}
                                                {p.parametri?.marca && (
                                                    <span className="insigna-detaliu angajat">
                                                        👤 {(() => {
                                                            const ang = gasesteAngajat(String(p.parametri.marca), date.angajati);
                                                            if (ang) {
                                                                return `${ang.nume} ${ang.prenume} (${ang.marca})`;
                                                            }
                                                            return `Marca: ${String(p.parametri.marca)}`;
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
                                                        📅 Data: <strong>{String(p.parametri.data)}</strong>
                                                    </span>
                                                )}
                                                {p.parametri?.data_inceput && p.parametri?.data_sfarsit && (
                                                    <span className="insigna-detaliu data">
                                                        📅 Interval: <strong>{String(p.parametri.data_inceput)} → {String(p.parametri.data_sfarsit)}</strong>
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
                                                ) : (
                                                    <button
                                                        type="button"
                                                        className="buton primar mic"
                                                        disabled={blocat}
                                                        onClick={() => aplicaPropunere(m.id, p)}
                                                    >
                                                        {t('ai.confirma_actiune')}
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
                <div ref={capatConversatieRef} />
            </div>

            {/* Sugestii rapide de comenzi */}
            <div className="sugestii-comenzi-ai">
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
            </div>

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
    );
}

