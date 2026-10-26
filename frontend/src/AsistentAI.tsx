import React, {useEffect, useRef, useState, type FormEvent} from 'react';
import {useTranslation} from 'react-i18next';
import {InterogheazaAI, OpresteAI, PornesteAI, SalveazaPlanificare, StareAI} from '../wailsjs/go/main/App';
import {ai, datebaza, main} from '../wailsjs/go/models';

interface MesajUI {
    id: string;
    rol: 'user' | 'assistant';
    text: string;
    propuneri?: ai.PropunereActiune[];
    propuneriAplicata?: Record<string, boolean>;
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
                setEroare(t(resp.eroare.cod, resp.eroare.param ?? {}));
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
            if (propunere.tip === 'planifica_tura' || propunere.tip === 'sterge_tura') {
                const marca = String(propunere.parametri['marca'] ?? '').trim();
                const dataStr = String(propunere.parametri['data'] ?? '').trim();
                const codSchimb = propunere.tip === 'sterge_tura' ? '' : String(propunere.parametri['cod_schimb'] ?? '').trim();

                // Găsește ID-ul angajatului după marcă sau nume
                const angajat = date.angajati.find(a => 
                    a.marca.toLowerCase() === marca.toLowerCase() ||
                    `${a.nume} ${a.prenume}`.toLowerCase().includes(marca.toLowerCase())
                );

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

