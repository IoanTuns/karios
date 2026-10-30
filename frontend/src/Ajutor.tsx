import {useEffect, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {ArticoleAjutor, DeschideManualAjutor} from '../wailsjs/go/main/App';
import {ajutor} from '../wailsjs/go/models';
import {Dialog} from './Formulare';

const normalizat = (s: string) => s.normalize('NFD').replace(/\p{Diacritic}/gu, '').toLocaleLowerCase();

export function PanouAjutor({ecran, articolId, versiuneCeruta, inchide, navigheaza, intreaba}: {
    ecran: string; articolId?: string; versiuneCeruta?: string; inchide: () => void;
    navigheaza?: (ecran: string) => boolean; intreaba?: (intrebare: string) => boolean;
}) {
    const {t, i18n} = useTranslation();
    const [articole, setArticole] = useState<ajutor.Articol[]>([]);
    const [cautare, setCautare] = useState('');
    const [toate, setToate] = useState(!!articolId);
    const [ales, setAles] = useState(articolId ?? '');
    const [versiuneConfirmata, setVersiuneConfirmata] = useState(false);
    const [eroare, setEroare] = useState('');
    const [incarcare, setIncarcare] = useState(true);
    useEffect(() => {
        let activ = true;
        setIncarcare(true);
        ArticoleAjutor(i18n.language).then(lista => {if (activ) {setArticole(lista); if (!articolId) setAles(lista.find(a => a.ecran === ecran)?.id ?? '');}})
            .catch(() => {if (activ) setEroare(t('date.eroare'));})
            .finally(() => {if (activ) setIncarcare(false);});
        return () => {activ = false;};
    }, [i18n.language]);
    const vizibile = articole.filter(a => (toate || a.ecran === ecran) && normalizat(a.titlu + ' ' + a.termeni.join(' ') + ' ' + a.pasi.join(' ')).includes(normalizat(cautare.trim())));
    const cautareActiva = cautare.trim().length > 0;
    const idVizibil = cautareActiva ? (vizibile.find(a => a.id === ales)?.id ?? vizibile[0]?.id ?? '') : ales;
    const articol = articole.find(a => a.id === idVizibil && (cautareActiva || versiuneConfirmata || !versiuneCeruta || a.versiune === versiuneCeruta));
    const invalid = !cautareActiva && !!ales && !articol && !incarcare;
    function articolesursa() { return articole.some(a => a.id === ales); }
    const titluEcran = t('ui.' + ecran);
    return <Dialog titlu={t('ajutor.titlu')} ocupat={false} inchide={inchide} clasa="dialog-ajutor" permiteMaximizare>
        <div className="corp-dialog ajutor-corp">
            <div className="ajutor-filtre">
                <label>{t('ajutor.cauta')}<input type="search" value={cautare} onChange={e => {setCautare(e.target.value); setToate(true);}}/></label>
                <label className="bifa"><input type="checkbox" checked={toate} onChange={e => setToate(e.target.checked)}/>{t('ajutor.toate')}</label>
            </div>
            {incarcare && <p role="status">{t('ui.incarcare')}</p>}
            {eroare && <p role="alert">{eroare}</p>}
            {invalid && <p role="alert">{t(articolesursa() ? 'ajutor.versiune_diferita' : 'ajutor.articol_lipsa')}{articolesursa() && <button type="button" className="legatura" onClick={() => setVersiuneConfirmata(true)}>{t('ajutor.versiune_curenta')}</button>}</p>}
            <div className="ajutor-coloane">
                <nav aria-label={t('ajutor.titlu')} className="ajutor-lista">
                    {!incarcare && !vizibile.length && <p role="status">{t('ajutor.fara_rezultate')}</p>}
                    {vizibile.map(a => <button type="button" key={a.id} className="buton secundar" aria-pressed={idVizibil === a.id} onClick={() => {setAles(a.id); setVersiuneConfirmata(true); setEroare('');}}>{a.titlu}</button>)}
                </nav>
                <div className="ajutor-articol" aria-live="polite">
                    {articol ? <article key={articol.id}>
                        <h3>{articol.titlu}</h3>
                        <details className="text-discret"><summary>{t('ajutor.detalii_sursa')}</summary>{t('ajutor.sursa', {id: articol.id, versiune: articol.versiune})}</details>
                        <ol>{articol.pasi.map((pas, i) => <li key={i}>{pas}</li>)}</ol>
                        <h4>{t('ajutor.rezultat')}</h4><p>{articol.rezultat}</p>
                        <h4>{t('ajutor.probleme')}</h4><p>{articol.probleme}</p>
                        <div className="ajutor-actiuni">
                            {navigheaza && <button type="button" className="buton secundar" onClick={() => {if (navigheaza(articol.ecran)) inchide();}}>{t('ajutor.deschide')}</button>}
                            {intreaba && <button type="button" className="buton secundar" onClick={() => {if (intreaba(t('ajutor.intrebare', {titlu: articol.titlu}))) inchide();}}>{t('ajutor.intreaba')}</button>}
                            <button type="button" className="legatura" onClick={() => void DeschideManualAjutor(articol.id).then(e => {if (e) setEroare(t(e.cod));}).catch(() => setEroare(t('date.eroare')))}>{t('ajutor.manual')}</button>
                        </div>
                    </article> : !invalid && <p>{cautareActiva && !vizibile.length ? t('ajutor.fara_rezultate') : toate ? t('ajutor.titlu') : titluEcran}</p>}
                </div>
            </div>
        </div>
    </Dialog>;
}

export function ButonAjutor({ecran, clasa = ''}: {ecran: string; clasa?: string}) {
    const {t} = useTranslation();
    const [deschis, setDeschis] = useState(false);
    return <div className={clasa}>
        <button type="button" className="buton secundar" onClick={() => setDeschis(true)}>{t('ajutor.buton')}</button>
        {deschis && <PanouAjutor ecran={ecran} inchide={() => setDeschis(false)}/>}
    </div>;
}
