import {useEffect, useId, useMemo, useRef, useState} from 'react';
import {useTranslation} from 'react-i18next';
import type {Font, Template} from '@pdfme/common';
import type {Designer} from '@pdfme/ui';
import {ListaFonturiSistem} from '../wailsjs/go/main/App';
import {Dialog} from './Formulare';
import {antetEditabil, antetFix, etichetePdfme, STATIC_PDFME, fontSistem, fonturiDinSablon, fonturiPentruSablon, pluginuriPdfme, TIPURI_PDFME, zoomPotrivit, type StareTabel} from './raportPdfme';
import {notifica} from './Toast';

const FONTURI_GO = new Set(['Go', 'GoBold', 'GoItalic', 'GoBoldItalic']);
const SELECTOR_POPUP_ANTD = '.ant-select-dropdown, .ant-dropdown, .ant-popover, .ant-picker-dropdown, .ant-color-picker, .ant-tooltip, .ant-cascader-dropdown';
// Latimea de la care panoul de proprietati porneste deschis (UIR-09).
const LATIME_PANOU = 1280;

// Variabilele raportului disponibile in textele pdfme, pe langa coloanele SQL.
export const VARIABILE_PDFME = ['titlu', 'perioada', 'an', 'luna', 'data_curenta', 'generat_la', 'antet_stanga', 'antet_centru', 'antet_dreapta', 'subsol_stanga', 'subsol_centru', 'subsol_dreapta', 'introducere', 'currentPage', 'totalPages'];

async function copiaza(text: string): Promise<boolean> {
    try {
        await navigator.clipboard.writeText(text);
        return true;
    } catch {
        // Rezerva pentru WebView fara acces la clipboard prin API.
        const zona = document.createElement('textarea');
        zona.value = text;
        document.body.appendChild(zona);
        zona.select();
        const ok = document.execCommand('copy');
        zona.remove();
        return ok;
    }
}


// EditorPdfme incadreaza Designer-ul pdfme. Sablonul circula ca text JSON,
// la fel cum se salveaza in model (continut.flexibil.sablon).
function EditorPdfme({sablon, schimba, fonturiRef}: {sablon: string; schimba: (s: string) => void; fonturiRef: React.MutableRefObject<((f: Font) => void) | null>}) {
    const {t} = useTranslation();
    const container = useRef<HTMLDivElement>(null);
    const designer = useRef<Designer | null>(null);
    const ultim = useRef(sablon);
    const schimbaRef = useRef(schimba);
    const [eroare, setEroare] = useState('');
    schimbaRef.current = schimba;

    useEffect(() => {
        let activ = true;
        const template = JSON.parse(ultim.current) as Template;
        const etichete = Object.fromEntries(TIPURI_PDFME.map(tip => [tip, t('rapoarte.pdfme_tip_' + tip)]));
        Promise.all([import('@pdfme/ui'), pluginuriPdfme(etichete), fonturiPentruSablon(template)]).then(([{Designer}, plugins, {font, lipsa}]) => {
            if (!activ || !container.current) return;
            if (lipsa.length) notifica(t('rapoarte.pdfme_fonturi_lipsa', {fonturi: lipsa.join(', ')}), 'avertisment', 8000);
            const latime = container.current.clientWidth;
            const panou = latime >= LATIME_PANOU;
            const baza = template.basePdf as {width?: number; height?: number};
            const d = new Designer({domContainer: container.current, template, plugins,
                options: {font, lang: 'en', labels: etichetePdfme(t), sidebarOpen: panou, zoomLevel: zoomPotrivit(latime, container.current.clientHeight, baza.width || 210, baza.height || 297, panou)}});
            d.onChangeTemplate(tpl => {
                const s = JSON.stringify(tpl);
                ultim.current = s;
                schimbaRef.current(s);
            });
            designer.current = d;
            fonturiRef.current = f => d.updateOptions({font: f});
        }).catch(() => { if (activ) setEroare(t('rapoarte.pdfme_invalid')); });
        return () => {
            activ = false;
            fonturiRef.current = null;
            designer.current?.destroy();
            designer.current = null;
        };
    }, []);

    // Sablonul inlocuit din afara (coloane noi, tabel actualizat) intra in editor.
    useEffect(() => {
        if (designer.current && sablon !== ultim.current) {
            ultim.current = sablon;
            designer.current.updateTemplate(JSON.parse(sablon) as Template);
        }
    }, [sablon]);

    return <>{eroare && <p role="alert" className="mesaj eroare">{eroare}</p>}<div ref={container} className="editor-pdfme"/></>;
}

export interface SpatiuPdfmeProps {
    sablon: string;
    schimba: (s: string) => void;
    coloaneSQL: string[];
    stare: StareTabel | null;
    ocupat: boolean;
    eroare: string;
    url: string;
    pagini: number;
    verifica: () => Promise<void>;
    actualizeazaTabel: () => Promise<void>;
    inchide: () => void;
}

// SpatiuPdfme este spatiul de lucru separat al editorului vizual (UIR-01):
// dialog marit, o singura derulare, unelte langa pagina si previzualizare in file.
export function SpatiuPdfme({sablon, schimba, coloaneSQL, stare, ocupat, eroare, url, pagini, verifica, actualizeazaTabel, inchide}: SpatiuPdfmeProps) {
    const {t} = useTranslation();
    const idLista = useId();
    const [fila, setFila] = useState<'pagina' | 'pdf'>('pagina');
    const [disponibile, setDisponibile] = useState<string[]>([]);
    const [ales, setAles] = useState('');
    const [adaugare, setAdaugare] = useState(false);
    const [eroareFont, setEroareFont] = useState('');
    const fonturiAdaugate = useRef<Font>({});
    const aplicaFonturi = useRef<((f: Font) => void) | null>(null);
    const [adaugate, setAdaugate] = useState<string[]>([]);
    // F09-T29: antetul si subsolul (pagina fixa) se editeaza mutate temporar pe pagina.
    const [antet, setAntet] = useState(false);
    const numeInainte = useRef<Set<string>>(new Set());
    const verificaDupaAntet = useRef(false);
    const sablonCurent = useRef(sablon);
    sablonCurent.current = sablon;

    // Listele derulante ale panoului pdfme (antd) se deschid in document.body, sub dialogul
    // modal (top layer), deci nu se vad si nu se pot alege (fonturi, aliniere). Containerele
    // lor se muta in dialog, unde raman deasupra editorului; pozitia o recalculeaza antd.
    useEffect(() => {
        const muta = () => {
            const dialog = document.querySelector('dialog.dialog-pdfme[open]');
            if (!dialog) return;
            for (const el of Array.from(document.body.children)) {
                if (el instanceof HTMLElement && !el.contains(dialog) && el.querySelector(SELECTOR_POPUP_ANTD)) dialog.appendChild(el);
            }
        };
        const observator = new MutationObserver(muta);
        observator.observe(document.body, {childList: true});
        muta();
        return () => observator.disconnect();
    }, []);

    useEffect(() => {
        let activ = true;
        ListaFonturiSistem().then(r => { if (activ) setDisponibile((r.fonturi || []).map(f => f.nume)); }).catch(() => {});
        return () => { activ = false; };
    }, []);

    const folosite = useMemo(() => {
        try {
            return [...fonturiDinSablon(JSON.parse(sablon))].filter(f => !FONTURI_GO.has(f));
        } catch {
            return [];
        }
    }, [sablon]);
    const fonturiAfisate = [...new Set([...folosite, ...adaugate])];

    async function adaugaFont() {
        const nume = ales.trim();
        if (!nume || adaugare) return;
        setAdaugare(true);
        setEroareFont('');
        try {
            const f = await fontSistem(nume);
            if (!f) {
                setEroareFont(t('fonturi.inexistent', {nume}));
                return;
            }
            // Fonturile sablonului sunt deja in editor; se adauga doar cele noi.
            const {font} = await fonturiPentruSablon(JSON.parse(sablon) as Template);
            fonturiAdaugate.current = {...fonturiAdaugate.current, [nume]: f};
            aplicaFonturi.current?.({...font, ...fonturiAdaugate.current});
            setAdaugate(a => (a.includes(nume) ? a : [...a, nume]));
            notifica(t('rapoarte.pdfme_font_adaugat', {nume}));
            setAles('');
        } catch {
            setEroareFont(t('fonturi.citire', {nume}));
        } finally {
            setAdaugare(false);
        }
    }

    async function copiazaCamp(camp: string) {
        const text = '{' + camp + '}';
        notifica(t(await copiaza(text) ? 'rapoarte.pdfme_copiat' : 'rapoarte.pdfme_copiere_esuata', {camp: text}), 'info', 2500);
    }

    async function verificaSiArata() {
        await verifica();
        setFila('pdf');
    }

    function intraInAntet() {
        try {
            const nou = antetEditabil(sablonCurent.current);
            // Doar elementele corpului paginii raman pe pagina la iesire; restul (si cele noi) trec in antet.
            numeInainte.current = new Set((JSON.parse(nou).schemas?.[0] || []).filter((e: Record<string, unknown>) => e[STATIC_PDFME] !== true).map((e: {name?: string}) => String(e.name)));
            setFila('pagina');
            setAntet(true);
            if (nou !== sablonCurent.current) schimba(nou);
        } catch {
            notifica(t('rapoarte.pdfme_invalid'), 'eroare');
        }
    }

    // iesiDinAntet pune elementele inapoi pe pagina fixa; intoarce sablonul rezultat.
    function iesiDinAntet(): string {
        let nou = sablonCurent.current;
        try {
            nou = antetFix(sablonCurent.current, numeInainte.current);
        } catch {
            // Sablon invalid: ramane cum este; verificarea il semnaleaza.
        }
        setAntet(false);
        if (nou !== sablonCurent.current) schimba(nou);
        return nou;
    }

    // Verificarea dupa iesirea din antet asteapta sablonul actualizat in model.
    useEffect(() => {
        if (verificaDupaAntet.current && !antet) {
            verificaDupaAntet.current = false;
            void verificaSiArata();
        }
    }, [sablon, antet]);

    function verificaOricand() {
        if (!antet) {
            void verificaSiArata();
            return;
        }
        // Efectul de mai sus verifica dupa ce modelul primeste sablonul cu antetul la loc.
        verificaDupaAntet.current = true;
        iesiDinAntet();
    }

    function inchideSpatiu() {
        if (antet) iesiDinAntet();
        inchide();
    }

    const nepotrivit = stare && stare.exista && !stare.potrivit;
    return <Dialog titlu={t('rapoarte.pdfme_spatiu_titlu')} ocupat={ocupat} inchide={inchideSpatiu} clasa="dialog-pdfme" maximizatInitial>
        <div className="spatiu-pdfme">
            <div className="bara-pdfme" role="toolbar" aria-label={t('rapoarte.pdfme_spatiu_titlu')}>
                <div className="file-pdfme" role="tablist">
                    <button type="button" role="tab" aria-selected={fila === 'pagina'} className={fila === 'pagina' ? 'activ' : ''} onClick={() => setFila('pagina')}>{t('rapoarte.pdfme_fila_pagina')}</button>
                    <button type="button" role="tab" aria-selected={fila === 'pdf'} className={fila === 'pdf' ? 'activ' : ''} disabled={!url || antet} onClick={() => setFila('pdf')}>{t('rapoarte.pdfme_fila_pdf')}{pagini > 0 ? ` (${t('rapoarte.pdfme_pagini', {numar: pagini})})` : ''}</button>
                </div>
                <details className="campuri-pdfme">
                    <summary>{t('rapoarte.pdfme_campuri')}</summary>
                    <div className="panou-pdfme"><p>{t('rapoarte.pdfme_campuri_indiciu')}</p>
                    <div className="lista-campuri-pdfme">
                        {VARIABILE_PDFME.map(v => <button key={v} type="button" className="eticheta-camp" onClick={() => void copiazaCamp(v)}>{'{' + v + '}'}</button>)}
                        {coloaneSQL.map(v => <button key={'sql-' + v} type="button" className="eticheta-camp sql" title={t('rapoarte.pdfme_coloana_sql')} onClick={() => void copiazaCamp(v)}>{'{' + v + '}'}</button>)}
                    </div></div>
                </details>
                <details className="fonturi-pdfme">
                    <summary>{t('rapoarte.pdfme_fonturi')}{fonturiAfisate.length ? ` (${fonturiAfisate.length})` : ''}</summary>
                    <div className="panou-pdfme"><p>{t('rapoarte.pdfme_fonturi_indiciu')}</p>
                    <div className="campuri">
                        <input name="font_windows" list={idLista} value={ales} onChange={e => setAles(e.target.value)} onKeyDown={e => { if (e.key === 'Enter') { e.preventDefault(); void adaugaFont(); } }} aria-label={t('rapoarte.pdfme_fonturi')}/>
                        <datalist id={idLista}>{disponibile.map(n => <option key={n} value={n}/>)}</datalist>
                        <button type="button" className="buton secundar" disabled={adaugare || !ales.trim()} onClick={() => void adaugaFont()}>{t('rapoarte.pdfme_adauga_font')}</button>
                    </div>
                    {eroareFont && <p role="alert" className="mesaj eroare">{eroareFont}</p>}
                    {fonturiAfisate.length > 0 && <p className="text-discret">{t('rapoarte.pdfme_fonturi_folosite')}: {fonturiAfisate.join(', ')}</p>}</div>
                </details>
                <div className="actiuni-pdfme">
                    <button type="button" className={'buton ' + (antet ? 'primar' : 'secundar')} aria-pressed={antet} disabled={ocupat} onClick={() => (antet ? iesiDinAntet() : intraInAntet())}>{t(antet ? 'rapoarte.pdfme_antet_gata' : 'rapoarte.pdfme_editeaza_antet')}</button>
                    {stare && stare.exista && <button type="button" className={'buton ' + (nepotrivit ? 'primar' : 'secundar')} disabled={ocupat} onClick={() => void actualizeazaTabel()}>{t('rapoarte.pdfme_actualizeaza_tabel')}</button>}
                    <button type="button" className="buton primar" disabled={ocupat} onClick={verificaOricand}>{t('rapoarte.flex_verifica')}</button>
                    <button type="button" className="buton secundar" disabled={ocupat} onClick={inchideSpatiu}>{t('rapoarte.pdfme_gata')}</button>
                </div>
            </div>
            {antet && <p role="status" className="mesaj atentionare indiciu-antet-pdfme">{t('rapoarte.pdfme_antet_indiciu')}</p>}
            {nepotrivit && <p role="status" className="mesaj avertisment">{t('rapoarte.pdfme_tabel_nepotrivit', {adaugate: stare!.adaugate.join(', ') || '—', eliminate: stare!.eliminate.join(', ') || '—'})}</p>}
            {ocupat && <p role="status">{t('ui.incarcare')}</p>}
            {eroare && <p role="alert" className="mesaj eroare">{eroare}</p>}
            <div className="corp-pdfme" hidden={fila !== 'pagina'}><EditorPdfme sablon={sablon} schimba={schimba} fonturiRef={aplicaFonturi}/></div>
            {fila === 'pdf' && url && <iframe title={t('rapoarte.previzualizare')} className="cadru-pdfme" src={url}/>}
        </div>
    </Dialog>;
}
