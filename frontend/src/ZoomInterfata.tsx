import {useEffect, useState} from 'react';
import {useTranslation} from 'react-i18next';

const niveluri = [75, 80, 90, 100, 110, 125, 150, 175, 200];
const cheie = 'pontix_zoom_interfata';

export function factorZoomInterfata(): number {
    return Number(document.documentElement.style.getPropertyValue('--zoom-interfata')) || 1;
}

export function ZoomInterfata() {
    const {t} = useTranslation();
    const [nivel, setNivel] = useState(() => {
        try { const v = Number(localStorage.getItem(cheie)); return niveluri.includes(v) ? v : 100; }
        catch { return 100; }
    });
    const schimba = (sens: number) => setNivel(v => niveluri[Math.max(0, Math.min(niveluri.length - 1, niveluri.indexOf(v) + sens))]);
    useEffect(() => {
        const radacina = document.documentElement;
        radacina.style.setProperty('--zoom-interfata', String(nivel / 100));
        try { localStorage.setItem(cheie, String(nivel)); } catch { /* Preferinta ramane valabila in sesiune. */ }
        window.dispatchEvent(new Event('resize'));
    }, [nivel]);
    useEffect(() => {
        const tastatura = (e: KeyboardEvent) => {
            if (!(e.ctrlKey || e.metaKey) || e.altKey) return;
            if (e.key === '+' || e.key === '=') { e.preventDefault(); schimba(1); }
            else if (e.key === '-') { e.preventDefault(); schimba(-1); }
            else if (e.key === '0') { e.preventDefault(); setNivel(100); }
        };
        window.addEventListener('keydown', tastatura);
        return () => window.removeEventListener('keydown', tastatura);
    }, []);
    return <div className="zoom-interfata" role="group" aria-label={t('ui.zoom_interfata')}>
        <button type="button" className="buton secundar" disabled={nivel === niveluri[0]} onClick={() => schimba(-1)} aria-label={t('ui.zoom_micsoreaza')} title={t('ui.zoom_micsoreaza')}>−</button>
        <button type="button" className="buton secundar" onClick={() => setNivel(100)} title={t('ui.zoom_reseteaza')} aria-label={t('ui.zoom_reseteaza')}><span aria-live="polite">{nivel}%</span></button>
        <button type="button" className="buton secundar" disabled={nivel === niveluri[niveluri.length - 1]} onClick={() => schimba(1)} aria-label={t('ui.zoom_mareste')} title={t('ui.zoom_mareste')}>+</button>
    </div>;
}
