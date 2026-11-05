import {useEffect, useId, useRef, useState} from 'react';
import {useTranslation} from 'react-i18next';
import type {Font, Template} from '@pdfme/common';
import type {Designer} from '@pdfme/ui';
import {ListaFonturiSistem} from '../wailsjs/go/main/App';
import {fontSistem, fonturiPentruSablon, pluginuriPdfme} from './raportPdfme';
import {notifica} from './Toast';

// EditorPdfme incadreaza Designer-ul pdfme. Sablonul circula ca text JSON,
// la fel cum se salveaza in model (continut.flexibil.sablon). Fonturile
// Windows se adauga la cerere, unul cate unul (D088).
export function EditorPdfme({sablon, schimba}: {sablon: string; schimba: (s: string) => void}) {
    const {t} = useTranslation();
    const idLista = useId();
    const container = useRef<HTMLDivElement>(null);
    const designer = useRef<Designer | null>(null);
    const fonturi = useRef<Font>({});
    const ultim = useRef(sablon);
    const schimbaRef = useRef(schimba);
    const [eroare, setEroare] = useState('');
    const [disponibile, setDisponibile] = useState<string[]>([]);
    const [ales, setAles] = useState('');
    const [ocupat, setOcupat] = useState(false);
    schimbaRef.current = schimba;

    useEffect(() => {
        let activ = true;
        ListaFonturiSistem().then(r => { if (activ) setDisponibile((r.fonturi || []).map(f => f.nume)); }).catch(() => {});
        const template = JSON.parse(ultim.current) as Template;
        Promise.all([import('@pdfme/ui'), pluginuriPdfme(), fonturiPentruSablon(template)]).then(([{Designer}, plugins, {font, lipsa}]) => {
            if (!activ || !container.current) return;
            fonturi.current = font;
            if (lipsa.length) notifica(t('rapoarte.pdfme_fonturi_lipsa', {fonturi: lipsa.join(', ')}), 'avertisment', 8000);
            const d = new Designer({domContainer: container.current, template, options: {font, lang: 'en'}, plugins});
            d.onChangeTemplate(tpl => {
                const s = JSON.stringify(tpl);
                ultim.current = s;
                schimbaRef.current(s);
            });
            designer.current = d;
        }).catch(() => { if (activ) setEroare(t('rapoarte.pdfme_invalid')); });
        return () => {
            activ = false;
            designer.current?.destroy();
            designer.current = null;
        };
    }, []);

    // Sablonul inlocuit din afara (de exemplu din coloanele SQL) intra in editor.
    useEffect(() => {
        if (designer.current && sablon !== ultim.current) {
            ultim.current = sablon;
            designer.current.updateTemplate(JSON.parse(sablon) as Template);
        }
    }, [sablon]);

    async function adaugaFont() {
        const nume = ales.trim();
        if (!nume || ocupat || !designer.current) return;
        setOcupat(true);
        setEroare('');
        try {
            if (!fonturi.current[nume]) {
                const f = await fontSistem(nume);
                if (!f) {
                    setEroare(t('fonturi.inexistent', {nume}));
                    return;
                }
                fonturi.current = {...fonturi.current, [nume]: f};
                designer.current?.updateOptions({font: fonturi.current});
            }
            notifica(t('rapoarte.pdfme_font_adaugat', {nume}));
            setAles('');
        } catch {
            setEroare(t('fonturi.citire', {nume}));
        } finally {
            setOcupat(false);
        }
    }

    return <>
        <fieldset className="campuri fonturi-pdfme">
            <legend>{t('rapoarte.pdfme_fonturi')}</legend>
            <p>{t('rapoarte.pdfme_fonturi_indiciu')}</p>
            <input name="font_windows" list={idLista} value={ales} onChange={e => setAles(e.target.value)} onKeyDown={e => { if (e.key === 'Enter') { e.preventDefault(); void adaugaFont(); } }} aria-label={t('rapoarte.pdfme_fonturi')}/>
            <datalist id={idLista}>{disponibile.map(n => <option key={n} value={n}/>)}</datalist>
            <button type="button" className="buton secundar" disabled={ocupat || !ales.trim()} onClick={() => void adaugaFont()}>{t('rapoarte.pdfme_adauga_font')}</button>
        </fieldset>
        {eroare && <p role="alert" className="mesaj eroare">{eroare}</p>}
        <div ref={container} className="editor-pdfme"/>
    </>;
}
