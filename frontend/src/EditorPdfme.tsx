import {useEffect, useRef, useState} from 'react';
import {useTranslation} from 'react-i18next';
import type {Template} from '@pdfme/common';
import type {Designer} from '@pdfme/ui';
import {fonturiPdfme, pluginuriPdfme} from './raportPdfme';

// EditorPdfme incadreaza Designer-ul pdfme. Sablonul circula ca text JSON,
// la fel cum se salveaza in model (continut.flexibil.sablon).
export function EditorPdfme({sablon, schimba}: {sablon: string; schimba: (s: string) => void}) {
    const {t} = useTranslation();
    const container = useRef<HTMLDivElement>(null);
    const designer = useRef<Designer | null>(null);
    const ultim = useRef(sablon);
    const schimbaRef = useRef(schimba);
    const [eroare, setEroare] = useState('');
    schimbaRef.current = schimba;

    useEffect(() => {
        let activ = true;
        Promise.all([import('@pdfme/ui'), pluginuriPdfme(), fonturiPdfme()]).then(([{Designer}, plugins, font]) => {
            if (!activ || !container.current) return;
            const d = new Designer({domContainer: container.current, template: JSON.parse(ultim.current) as Template, options: {font, lang: 'en'}, plugins});
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

    return <>{eroare && <p role="alert" className="mesaj eroare">{eroare}</p>}<div ref={container} className="editor-pdfme"/></>;
}
