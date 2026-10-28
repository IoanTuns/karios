import React, {useRef, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {AnalizeazaImportAngajati, AnalizeazaImportSarbatori, DescarcaSablon, ExecutaImportAngajati, ExecutaImportSarbatori} from '../wailsjs/go/main/App';
import {datebaza, main} from '../wailsjs/go/models';

export function descarcaFisierDinBase64(continutBase64: string, numeFisier: string) {
    const bytes = atob(continutBase64);
    const byteNumbers = new Uint8Array(bytes.length);
    for (let i = 0; i < bytes.length; i++) {
        byteNumbers[i] = bytes.charCodeAt(i);
    }
    const tip = numeFisier.endsWith('.xlsx')
        ? 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
        : 'text/csv;charset=utf-8;';
    const blob = new Blob([byteNumbers], {type: tip});
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = numeFisier;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
}

function citesteFisierBase64(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = () => {
            const res = reader.result as string;
            const idx = res.indexOf(',');
            resolve(idx >= 0 ? res.slice(idx + 1) : res);
        };
        reader.onerror = () => reject(reader.error);
        reader.readAsDataURL(file);
    });
}

export function DialogImportAngajati({
    inchide,
    laSucces,
}: {
    inchide: () => void;
    laSucces: (dateNoi: datebaza.Date, mesaj: string) => void;
}) {
    const {t} = useTranslation();
    const inputFisier = useRef<HTMLInputElement>(null);
    const [analiza, setAnaliza] = useState<datebaza.RezultatAnalizaImportAngajati | null>(null);
    const [actualizeazaExistente, setActualizeazaExistente] = useState(true);
    const [ocupat, setOcupat] = useState(false);
    const [eroare, setEroare] = useState('');

    async function descarca(format: 'xlsx' | 'csv') {
        try {
            setOcupat(true);
            const resp = await DescarcaSablon('angajati', format);
            if (resp.eroare) {
                setEroare(t(resp.eroare.cod, resp.eroare.param ?? {}));
            } else if (resp.continut_base64) {
                descarcaFisierDinBase64(resp.continut_base64, resp.nume_fisier);
            }
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setOcupat(false);
        }
    }

    async function schimbaFisier(e: React.ChangeEvent<HTMLInputElement>) {
        const file = e.target.files?.[0];
        if (!file) return;
        setEroare('');
        setOcupat(true);
        try {
            const b64 = await citesteFisierBase64(file);
            const resp = await AnalizeazaImportAngajati(b64, file.name);
            if (resp.eroare) {
                setEroare(t(resp.eroare.cod, resp.eroare.param ?? {}));
                setAnaliza(null);
            } else if (resp.rezultat) {
                setAnaliza(resp.rezultat);
            }
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setOcupat(false);
        }
    }

    async function executaImport() {
        if (!analiza || ocupat) return;
        setEroare('');
        setOcupat(true);
        try {
            const resp = await ExecutaImportAngajati(new datebaza.CerereExecutaImportAngajati({
                randuri: analiza.randuri,
                actualizeaza_existente: actualizeazaExistente,
            }));
            if (resp.eroare) {
                setEroare(t(resp.eroare.cod, resp.eroare.param ?? {}));
            } else if (resp.date) {
                laSucces(resp.date, t('import.succes_angajati', {numar: resp.importati}));
                inchide();
            }
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setOcupat(false);
        }
    }

    const randuriEligibile = analiza?.randuri.filter(r => r.valid && (actualizeazaExistente || !r.exista_deja)) ?? [];

    return (
        <div className="ecran-blocat" role="dialog" aria-modal="true" aria-labelledby="titlu-import-angajati">
            <div className="fereastra-import">
                <header className="antet-import">
                    <div>
                        <h2 id="titlu-import-angajati">{t('import.titlu_angajati')}</h2>
                        <p className="descriere">{t('import.descriere_angajati')}</p>
                    </div>
                    <button type="button" className="buton-inchide" onClick={inchide} aria-label={t('import.buton_anuleaza')}>✕</button>
                </header>

                <div className="corp-import">
                    <div className="sectiune-actiuni-sabloane">
                        <span className="eticheta-sablon">{t('ui.descarca_sablon')}:</span>
                        <button type="button" className="buton secundar mic" disabled={ocupat} onClick={() => descarca('xlsx')}>
                            📊 {t('import.sablon_excel')}
                        </button>
                        <button type="button" className="buton secundar mic" disabled={ocupat} onClick={() => descarca('csv')}>
                            📄 {t('import.sablon_csv')}
                        </button>
                    </div>

                    <div className="sectiune-selectie-fisier">
                        <input
                            ref={inputFisier}
                            type="file"
                            accept=".xlsx,.csv"
                            style={{display: 'none'}}
                            onChange={schimbaFisier}
                        />
                        <button
                            type="button"
                            className="buton primar"
                            disabled={ocupat}
                            onClick={() => inputFisier.current?.click()}
                        >
                            📁 {t('import.alege_fisier')}
                        </button>
                        {analiza && (
                            <span className="info-fisier-ales">
                                {t('import.fisier_selectat', {nume: analiza.nume_fisier, format: analiza.format})}
                            </span>
                        )}
                    </div>

                    {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}

                    {analiza && (
                        <div className="sectiune-analiza">
                            <div className="bara-statistici-import">
                                <p>
                                    {t('import.rezumat', {
                                        total: analiza.total_randuri,
                                        valide: analiza.total_valide,
                                        existente: analiza.total_existente,
                                        invalide: analiza.total_invalide,
                                    })}
                                </p>
                                <label className="optiune-conflict">
                                    <input
                                        type="checkbox"
                                        checked={actualizeazaExistente}
                                        onChange={e => setActualizeazaExistente(e.target.checked)}
                                    />
                                    {t('import.actualizeaza_existente')}
                                </label>
                            </div>

                            <div className="tabel-derulant-import">
                                <table className="tabel-import">
                                    <thead>
                                        <tr>
                                            <th>{t('import.coloana_stare')}</th>
                                            <th>{t('ui.marca')}</th>
                                            <th>{t('ui.angajat')}</th>
                                            <th>{t('ui.functie')}</th>
                                            <th>{t('ui.sectie')}</th>
                                            <th>{t('ui.departament')}</th>
                                            <th>{t('ui.contract')}</th>
                                            <th>{t('ui.norma')}</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {analiza.randuri.map((r, i) => {
                                            const esteExistent = r.exista_deja;
                                            const esteInvalid = !r.valid;
                                            return (
                                                <tr
                                                    key={i}
                                                    className={
                                                        esteInvalid
                                                            ? 'rand-import-eroare'
                                                            : esteExistent
                                                            ? 'rand-import-existent'
                                                            : 'rand-import-valid'
                                                    }
                                                >
                                                    <td>
                                                        {esteInvalid ? (
                                                            <span className="stare-eticheta eroare" title={r.mesaj_eroare}>
                                                                ✕ {r.mesaj_eroare || t('import.stare_eroare')}
                                                            </span>
                                                        ) : esteExistent ? (
                                                            <span className="stare-eticheta avertisment">
                                                                ⚠ {t('import.stare_existent')}
                                                            </span>
                                                        ) : (
                                                            <span className="stare-eticheta succes">
                                                                ✓ {t('import.stare_valid')}
                                                            </span>
                                                        )}
                                                    </td>
                                                    <td><strong>{r.marca || '—'}</strong></td>
                                                    <td>{r.nume} {r.prenume}</td>
                                                    <td>{r.functie || '—'}</td>
                                                    <td>{r.sectie || '—'}</td>
                                                    <td>{r.departament || '—'}</td>
                                                    <td>{r.tip_contract}</td>
                                                    <td>{Math.round(r.norma_minute / 60)}h</td>
                                                </tr>
                                            );
                                        })}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    )}
                </div>

                <footer className="subsol-import">
                    <button type="button" className="buton secundar" disabled={ocupat} onClick={inchide}>
                        {t('import.buton_anuleaza')}
                    </button>
                    <button
                        type="button"
                        className="buton primar"
                        disabled={ocupat || !analiza || randuriEligibile.length === 0}
                        onClick={executaImport}
                    >
                        {t('import.buton_executa', {numar: randuriEligibile.length})}
                    </button>
                </footer>
            </div>
        </div>
    );
}

export function DialogImportSarbatori({
    an,
    inchide,
    laSucces,
}: {
    an: number;
    inchide: () => void;
    laSucces: (calNou: datebaza.CalendarAn, mesaj: string) => void;
}) {
    const {t} = useTranslation();
    const inputFisier = useRef<HTMLInputElement>(null);
    const [analiza, setAnaliza] = useState<datebaza.RezultatAnalizaImportSarbatori | null>(null);
    const [inlocuiesteExistente, setInlocuiesteExistente] = useState(true);
    const [ocupat, setOcupat] = useState(false);
    const [eroare, setEroare] = useState('');

    async function descarca(format: 'xlsx' | 'csv') {
        try {
            setOcupat(true);
            const resp = await DescarcaSablon('sarbatori', format);
            if (resp.eroare) {
                setEroare(t(resp.eroare.cod, resp.eroare.param ?? {}));
            } else if (resp.continut_base64) {
                descarcaFisierDinBase64(resp.continut_base64, resp.nume_fisier);
            }
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setOcupat(false);
        }
    }

    async function schimbaFisier(e: React.ChangeEvent<HTMLInputElement>) {
        const file = e.target.files?.[0];
        if (!file) return;
        setEroare('');
        setOcupat(true);
        try {
            const b64 = await citesteFisierBase64(file);
            const resp = await AnalizeazaImportSarbatori(b64, file.name);
            if (resp.eroare) {
                setEroare(t(resp.eroare.cod, resp.eroare.param ?? {}));
                setAnaliza(null);
            } else if (resp.rezultat) {
                setAnaliza(resp.rezultat);
            }
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setOcupat(false);
        }
    }

    async function executaImport() {
        if (!analiza || ocupat) return;
        setEroare('');
        setOcupat(true);
        try {
            const resp = await ExecutaImportSarbatori(
                new datebaza.CerereExecutaImportSarbatori({
                    randuri: analiza.randuri,
                    inlocuieste_existente: inlocuiesteExistente,
                }),
                an
            );
            if (resp.eroare) {
                setEroare(t(resp.eroare.cod, resp.eroare.param ?? {}));
            } else if (resp.calendar) {
                laSucces(resp.calendar, t('import.succes_sarbatori', {numar: resp.importate}));
                inchide();
            }
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setOcupat(false);
        }
    }

    const randuriEligibile = analiza?.randuri.filter(r => r.valid && (inlocuiesteExistente || !r.exista_deja)) ?? [];

    return (
        <div className="ecran-blocat" role="dialog" aria-modal="true" aria-labelledby="titlu-import-sarbatori">
            <div className="fereastra-import">
                <header className="antet-import">
                    <div>
                        <h2 id="titlu-import-sarbatori">{t('import.titlu_sarbatori')}</h2>
                        <p className="descriere">{t('import.descriere_sarbatori')}</p>
                    </div>
                    <button type="button" className="buton-inchide" onClick={inchide} aria-label={t('import.buton_anuleaza')}>✕</button>
                </header>

                <div className="corp-import">
                    <div className="sectiune-actiuni-sabloane">
                        <span className="eticheta-sablon">{t('ui.descarca_sablon')}:</span>
                        <button type="button" className="buton secundar mic" disabled={ocupat} onClick={() => descarca('xlsx')}>
                            📊 {t('import.sablon_excel')}
                        </button>
                        <button type="button" className="buton secundar mic" disabled={ocupat} onClick={() => descarca('csv')}>
                            📄 {t('import.sablon_csv')}
                        </button>
                    </div>

                    <div className="sectiune-selectie-fisier">
                        <input
                            ref={inputFisier}
                            type="file"
                            accept=".xlsx,.csv"
                            style={{display: 'none'}}
                            onChange={schimbaFisier}
                        />
                        <button
                            type="button"
                            className="buton primar"
                            disabled={ocupat}
                            onClick={() => inputFisier.current?.click()}
                        >
                            📁 {t('import.alege_fisier')}
                        </button>
                        {analiza && (
                            <span className="info-fisier-ales">
                                {t('import.fisier_selectat', {nume: analiza.nume_fisier, format: analiza.format})}
                            </span>
                        )}
                    </div>

                    {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}

                    {analiza && (
                        <div className="sectiune-analiza">
                            <div className="bara-statistici-import">
                                <p>
                                    {t('import.rezumat', {
                                        total: analiza.total_randuri,
                                        valide: analiza.total_valide,
                                        existente: analiza.total_existente,
                                        invalide: analiza.total_invalide,
                                    })}
                                </p>
                                <label className="optiune-conflict">
                                    <input
                                        type="checkbox"
                                        checked={inlocuiesteExistente}
                                        onChange={e => setInlocuiesteExistente(e.target.checked)}
                                    />
                                    {t('import.inlocuieste_existente')}
                                </label>
                            </div>

                            <div className="tabel-derulant-import">
                                <table className="tabel-import">
                                    <thead>
                                        <tr>
                                            <th>{t('import.coloana_stare')}</th>
                                            <th>{t('sarbatori.data')}</th>
                                            <th>{t('sarbatori.denumire')} (RO)</th>
                                            <th>{t('sarbatori.denumire')} (EN)</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {analiza.randuri.map((r, i) => {
                                            const esteExistent = r.exista_deja;
                                            const esteInvalid = !r.valid;
                                            return (
                                                <tr
                                                    key={i}
                                                    className={
                                                        esteInvalid
                                                            ? 'rand-import-eroare'
                                                            : esteExistent
                                                            ? 'rand-import-existent'
                                                            : 'rand-import-valid'
                                                    }
                                                >
                                                    <td>
                                                        {esteInvalid ? (
                                                            <span className="stare-eticheta eroare" title={r.mesaj_eroare}>
                                                                ✕ {r.mesaj_eroare || t('import.stare_eroare')}
                                                            </span>
                                                        ) : esteExistent ? (
                                                            <span className="stare-eticheta avertisment">
                                                                ⚠ {t('import.stare_existent')}
                                                            </span>
                                                        ) : (
                                                            <span className="stare-eticheta succes">
                                                                ✓ {t('import.stare_valid')}
                                                            </span>
                                                        )}
                                                    </td>
                                                    <td><strong>{r.data || '—'}</strong></td>
                                                    <td>{r.denumire_ro}</td>
                                                    <td>{r.denumire_en}</td>
                                                </tr>
                                            );
                                        })}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    )}
                </div>

                <footer className="subsol-import">
                    <button type="button" className="buton secundar" disabled={ocupat} onClick={inchide}>
                        {t('import.buton_anuleaza')}
                    </button>
                    <button
                        type="button"
                        className="buton primar"
                        disabled={ocupat || !analiza || randuriEligibile.length === 0}
                        onClick={executaImport}
                    >
                        {t('import.buton_executa', {numar: randuriEligibile.length})}
                    </button>
                </footer>
            </div>
        </div>
    );
}

