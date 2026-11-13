import {useEffect, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {VerificaLuna} from '../wailsjs/go/main/App';
import {datebaza} from '../wailsjs/go/models';
import {Dialog} from './Formulare';

// F12-T06: „Verifica luna” — constatarile lunii, cele mai grave primele; nu schimba nimic.
export function DialogVerificaLuna({an, luna, nesalvate, inchide}: {an: number; luna: number; nesalvate: number; inchide: () => void}) {
    const {t} = useTranslation();
    const [verificare, setVerificare] = useState<datebaza.VerificareLuna | null>(null);
    const [eroare, setEroare] = useState('');

    useEffect(() => {
        VerificaLuna(an, luna).then(r => {
            if (r.eroare) setEroare(t(r.eroare.cod, r.eroare.param ?? {}));
            else setVerificare(r.verificare ?? null);
        }).catch(() => setEroare(t('date.eroare')));
    }, [an, luna, t]);

    const titlu = (p: datebaza.ProblemaLuna) => p.cod === 'necalculat'
        ? t(`pontaj.motiv_${p.param?.motiv}`, {cod: p.param?.cod ?? ''})
        : t(`verifica_luna.${p.cod}`, {numar: String(p.numar), ...(p.param ?? {})});
    const probleme = verificare?.probleme ?? [];
    return <Dialog titlu={t('verifica_luna.titlu', {luna: `${String(luna).padStart(2, '0')}.${an}`})} ocupat={false} inchide={inchide}>
        <div className="corp-dialog verifica-luna">
            {nesalvate > 0 && <p className="mesaj eroare" role="alert">{t('verifica_luna.nesalvate', {numar: String(nesalvate)})}</p>}
            {eroare && <p className="mesaj eroare" role="alert">{eroare}</p>}
            {!verificare && !eroare && <p role="status">{t('ui.incarcare')}</p>}
            {verificare && (probleme.length === 0
                ? <p className="mesaj succes" role="status">{t('verifica_luna.nimic')}</p>
                : <ul className="lista-verificare-luna">{probleme.map((p, i) => <li key={i} className={`gravitate-${p.gravitate.toLowerCase()}`}>
                    <strong>{t(`verifica_luna.gravitate_${p.gravitate}`)}</strong> {titlu(p)}
                    {p.cod === 'necalculat' && <span> {t('verifica_luna.persoane', {numar: String(p.numar)})}</span>}
                    {p.exemple?.length > 0 && <small className="text-discret"> {p.exemple.join(', ')}{p.numar > p.exemple.length && p.cod !== 'absente_in_asteptare' ? ' …' : ''}</small>}
                </li>)}</ul>)}
            <p className="text-discret">{t('verifica_luna.indiciu')}</p>
        </div>
        <footer className="subsol-dialog"><button type="button" className="buton primar" onClick={inchide}>{t('ui.inchide')}</button></footer>
    </Dialog>;
}
