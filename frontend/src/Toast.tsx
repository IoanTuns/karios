import {useTranslation} from 'react-i18next';
import {useEffect, useState} from 'react';

export type TipToast = 'succes' | 'info' | 'avertisment' | 'eroare';

export interface ToastItem {
    id: string;
    tip: TipToast;
    mesaj: string;
    durataMs: number;
}

// notifica emite o notificare globală afișată în ToastContainer
export function notifica(mesaj: string, tip: TipToast = 'succes', durataMs = 3500) {
    const ev = new CustomEvent<ToastItem>('pontix-toast', {
        detail: {
            id: `${Date.now()}-${Math.random().toString(36).slice(2, 7)}`,
            tip,
            mesaj,
            durataMs,
        },
    });
    window.dispatchEvent(ev);
}

export function ToastContainer() {
    const {t} = useTranslation();
    const [toaste, setToaste] = useState<ToastItem[]>([]);

    useEffect(() => {
        const handler = (e: Event) => {
            const custom = e as CustomEvent<ToastItem>;
            if (custom.detail) {
                setToaste(prev => [...prev.slice(-4), custom.detail]);
            }
        };
        window.addEventListener('pontix-toast', handler);
        return () => window.removeEventListener('pontix-toast', handler);
    }, []);

    const elimina = (id: string) => {
        setToaste(prev => prev.filter(t => t.id !== id));
    };

    return (
        <div className="container-toasturi" role="region" aria-live="polite" aria-label={t('ui.notificari')}>
            {toaste.map(t => (
                <ToastCard key={t.id} toast={t} onElimina={() => elimina(t.id)} />
            ))}
        </div>
    );
}

function ToastCard({toast, onElimina}: {toast: ToastItem; onElimina: () => void}) {
    const {t} = useTranslation();
    useEffect(() => {
        const timer = setTimeout(onElimina, toast.durataMs);
        return () => clearTimeout(timer);
    }, [toast.durataMs, onElimina]);

    const pictograme: Record<TipToast, string> = {
        succes: '✓',
        info: 'ℹ',
        avertisment: '⚠️',
        eroare: '✕',
    };

    return (
        <div className={`toast-card toast-${toast.tip}`} role="status">
            <span className="simbol-toast" aria-hidden="true">{pictograme[toast.tip]}</span>
            <span className="mesaj-toast">{toast.mesaj}</span>
            <button
                type="button"
                className="buton-inchide-toast"
                onClick={onElimina}
                aria-label={t('ui.inchide_notificarea')}
            >
                ×
            </button>
        </div>
    );
}

