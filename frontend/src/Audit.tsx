// Avertizarea la incidentul de audit și confirmarea lui cu motiv (D044).
import {useEffect, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {ConfirmaIncidentAudit, StareAudit} from '../wailsjs/go/main/App';
import {audit} from '../wailsjs/go/models';
import {Dialog} from './Formulare';
import {notifica} from './Toast';

export function BannerAudit() {
    const {t} = useTranslation();
    const [rezultat, setRezultat] = useState<audit.Rezultat | null>(null);
    const [dialog, setDialog] = useState(false);
    const [motiv, setMotiv] = useState('');
    const [ocupat, setOcupat] = useState(false);
    const [eroare, setEroare] = useState('');

    useEffect(() => {
        let activ = true;
        StareAudit().then(r => { if (activ) setRezultat(r.rezultat); }).catch(() => {});
        return () => { activ = false; };
    }, []);

    if (rezultat?.stare !== 'incident') return null;

    async function confirma() {
        if (!motiv.trim()) { setEroare(t('audit.motiv_obligatoriu')); return; }
        setOcupat(true);
        setEroare('');
        try {
            const r = await ConfirmaIncidentAudit(motiv);
            if (r.eroare) { setEroare(t(r.eroare.cod, r.eroare.param ?? {})); return; }
            setRezultat(r.rezultat);
            setDialog(false);
            notifica(t('audit.confirmat'));
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setOcupat(false);
        }
    }

    const tabele = (rezultat.tabele ?? []).join(', ');
    return <>
        <div className="banner-audit" role="alert">
            <span className="simbol-audit" aria-hidden="true">⚠</span>
            <div>
                <strong>{t('audit.titlu_incident')}</strong>
                <p>{t('audit.descriere_incident')}</p>
                {tabele && <p>{t('audit.tabele', {tabele})}</p>}
                {rezultat.rand_rupt > 0 && <p>{t('audit.rand_rupt', {rand: String(rezultat.rand_rupt)})}</p>}
                {rezultat.copie_veche && <p>{t('audit.copie_veche')}</p>}
                <p className="text-discret">{t('audit.verificat_la', {moment: new Date(rezultat.verificat_la).toLocaleString()})}</p>
            </div>
            <button type="button" className="buton secundar mic" onClick={() => { setMotiv(''); setEroare(''); setDialog(true); }}>{t('audit.confirma')}</button>
        </div>
        {dialog && <Dialog titlu={t('audit.confirma')} ocupat={ocupat} inchide={() => setDialog(false)}>
            <div className="corp-dialog">
                <p>{t('audit.confirma_descriere')}</p>
                <label>{t('audit.motiv')}
                    <textarea name="motiv_audit" rows={3} value={motiv} disabled={ocupat} onChange={e => setMotiv(e.target.value)} required/>
                </label>
                {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}
            </div>
            <footer className="subsol-dialog">
                <button className="buton secundar" disabled={ocupat} onClick={() => setDialog(false)}>{t('ui.anuleaza')}</button>
                <button className="buton primar" disabled={ocupat || !motiv.trim()} onClick={() => void confirma()}>{t(ocupat ? 'ui.salvare' : 'audit.confirma')}</button>
            </footer>
        </Dialog>}
    </>;
}

// AmprentaLuna arată amprenta unei luni închise și dacă datele curente îi mai corespund.
export function AmprentaLuna({amprenta, corespunde}: {amprenta: string; corespunde: boolean}) {
    const {t} = useTranslation();
    if (!amprenta) return null;
    return <p className={corespunde ? 'amprenta-luna' : 'amprenta-luna diferita'}>
        {t('pontaj.amprenta', {amprenta})} · {t(corespunde ? 'pontaj.amprenta_corespunde' : 'pontaj.amprenta_diferita')}
    </p>;
}
