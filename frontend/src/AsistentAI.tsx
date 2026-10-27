import React, {useEffect, useRef, useState, type FormEvent} from 'react';
import {useTranslation} from 'react-i18next';
import {EsteAutoPornireAI, InterogheazaAI, OpresteAI, PornesteAI, SalveazaPlanificare, SeteazaAutoPornireAI, StareAI} from '../wailsjs/go/main/App';
import {ai, datebaza, main} from '../wailsjs/go/models';

interface MesajUI {
    id: string;
    rol: 'user' | 'assistant';
    text: string;
    propuneri?: ai.PropunereActiune[];
    propuneriAplicata?: Record<string, boolean>;
}

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

export function AsistentAI({date}: {date: datebaza.Date}) {
    const {t} = useTranslation();
    const [stare, setStare] = useState<ai.StareMediu | null>(null);
    const [incarcare, setIncarcare] = useState(true);
    const [ocupat, setOcupat] = useState(false);
    const [mesaj, setMesaj] = useState('');
    const [eroare, setEroare] = useState('');
    const [textIntrare, setTextIntrare] = useState('');
    const [conversatie, setConversatie] = useState<MesajUI[]>([]);
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
        // Mesaj inițial de întâmpinare
        setConversatie([
            {
                id: 'bun_venit',
                rol: 'assistant',
                text: t('ai.mesaj_bun_venit'),
            },
        ]);
    }, []);

    useEffect(() => {
        capatConversatieRef.current?.scrollIntoView({behavior: 'smooth'});
    }, [conversatie]);

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
        if (!text || ocupat) return;

        setEroare('');
        setTextIntrare('');

        const mesajUtilizator: MesajUI = {
            id: `u_${Date.now()}`,
            rol: 'user',
            text,
        };

        const istoricNou = [...conversatie, mesajUtilizator];
        setConversatie(istoricNou);
        setOcupat(true);

        try {
            // Trimite ultimele 6 mesaje pentru context
            const istoricAPI: ai.MesajConversatie[] = istoricNou
                .filter(m => m.id !== 'bun_venit')
                .slice(-6)
                .map(m => new ai.MesajConversatie({
                    role: m.rol,
                    content: m.text,
                }));

            const resp = await InterogheazaAI(new main.CerereInterogareAI({
                istoric: istoricAPI,
                an: anCurent,
                luna: lunaCurenta,
            }));

            if (resp.eroare) {
                const textEroare = t(resp.eroare.cod, resp.eroare.param ?? {});
                setEroare(textEroare !== resp.eroare.cod ? textEroare : (resp.eroare.param?.detalii ?? resp.eroare.cod));
            } else if (resp.raspuns) {
                const mesajAsistent: MesajUI = {
                    id: `a_${Date.now()}`,
                    rol: 'assistant',
                    text: resp.raspuns.mesaj || (resp.raspuns.propuneri?.length ? t('ai.propuneri_titlu') : ''),
                    propuneri: resp.raspuns.propuneri,
                    propuneriAplicata: {},
                };
                setConversatie([...istoricNou, mesajAsistent]);
            }
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setOcupat(false);
            void actualizeazaStare();
        }
    }

    async function aplicaPropunere(mesajId: string, propunere: ai.PropunereActiune) {
        setOcupat(true);
        setEroare('');
        setMesaj('');

        try {
            if (propunere.tip === 'planifica_interval' || propunere.tip === 'sterge_interval') {
                const marca = String(propunere.parametri['marca'] ?? '').trim();
                const dInceput = String(propunere.parametri['data_inceput'] ?? '').trim();
                const dSfarsit = String(propunere.parametri['data_sfarsit'] ?? '').trim();
                const codSchimb = propunere.tip === 'sterge_interval' ? '' : String(propunere.parametri['cod_schimb'] ?? '').trim();
                const doarLucratoare = Boolean(propunere.parametri['doar_lucratoare']);

                const angajat = gasesteAngajat(marca, date.angajati);

                if (!angajat) {
                    setEroare(`Nu a fost găsit angajatul "${marca}".`);
                    setOcupat(false);
                    return;
                }

                const start = new Date(dInceput);
                const end = new Date(dSfarsit);
                const pozitii: datebaza.PozitiePlanificare[] = [];

                const cur = new Date(start);
                while (cur <= end) {
                    const ziSapt = cur.getDay(); // 0 = duminica, 6 = sambata
                    if (!doarLucratoare || (ziSapt !== 0 && ziSapt !== 6)) {
                        const anStr = cur.getFullYear();
                        const lunaStr = String(cur.getMonth() + 1).padStart(2, '0');
                        const ziStr = String(cur.getDate()).padStart(2, '0');
                        pozitii.push(new datebaza.PozitiePlanificare({
                            angajat_id: angajat.id,
                            data: `${anStr}-${lunaStr}-${ziStr}`,
                            cod_schimb: codSchimb,
                        }));
                    }
                    cur.setDate(cur.getDate() + 1);
                }

                if (pozitii.length === 0) {
                    setEroare("Nu au fost identificate zile eligibile în intervalul solicitat.");
                    setOcupat(false);
                    return;
                }

                const an = start.getFullYear() || anCurent;
                const luna = start.getMonth() + 1 || lunaCurenta;
                const r = await SalveazaPlanificare(an, luna, pozitii);
                if (r.eroare) {
                    setEroare(t(r.eroare.cod, r.eroare.param ?? {}));
                } else {
                    setMesaj(t('ai.actiune_aplicata'));
                    setConversatie(prev => prev.map(m => {
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
            } else if (propunere.tip === 'planifica_tura' || propunere.tip === 'sterge_tura') {
                const marca = String(propunere.parametri['marca'] ?? '').trim();
                const dataStr = String(propunere.parametri['data'] ?? '').trim();
                const codSchimb = propunere.tip === 'sterge_tura' ? '' : String(propunere.parametri['cod_schimb'] ?? '').trim();

                // Găsește ID-ul angajatului după marcă sau nume (tolerant la formatare)
                const angajat = gasesteAngajat(marca, date.angajati);

                if (!angajat) {
                    setEroare(`Nu a fost găsit angajatul cu marca sau numele "${marca}".`);
                    setOcupat(false);
                    return;
                }

                const d = new Date(dataStr);
                const an = d.getFullYear() || anCurent;
                const luna = d.getMonth() + 1 || lunaCurenta;

                const pozitie = new datebaza.PozitiePlanificare({
                    angajat_id: angajat.id,
                    data: dataStr,
                    cod_schimb: codSchimb,
                });

                const r = await SalveazaPlanificare(an, luna, [pozitie]);
                if (r.eroare) {
                    setEroare(t(r.eroare.cod, r.eroare.param ?? {}));
                } else {
                    setMesaj(t('ai.actiune_aplicata'));
                    // Marchează propunerea ca aplicată
                    setConversatie(prev => prev.map(m => {
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
            }
        } catch {
            setEroare(t('date.eroare'));
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

            {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}
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
                                                    <span className="eticheta-aplicat">✓ {t('ai.actiune_aplicata')}</span>
                                                ) : (
                                                    <button
                                                        type="button"
                                                        className="buton primar mic"
                                                        disabled={ocupat}
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
                {ocupat && (
                    <div className="bula-mesaj assistant scriere">
                        <span className="puncte-asteptare">● ● ●</span>
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
                        disabled={ocupat}
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
                    disabled={ocupat}
                    onChange={e => setTextIntrare(e.target.value)}
                    placeholder={t('ai.indiciu_scriere')}
                />
                <button
                    type="submit"
                    className="buton primar"
                    disabled={ocupat || !textIntrare.trim()}
                >
                    {t('ai.trimite')} ➔
                </button>
            </form>
        </div>
    );
}

