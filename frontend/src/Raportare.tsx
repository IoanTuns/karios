// „Raportează o problemă / Părerea ta” (F16-T15, I093): formularul, raportul anonimizat pe care
// utilizatorul îl vede și îl poate modifica, apoi salvarea ca fișier (pentru email) sau deschiderea
// paginii GitHub de issue completate. Aplicația nu trimite nimic singură.
import {useState} from 'react';
import {useTranslation} from 'react-i18next';
import {DeschideFolderRapoarteProbleme, PregatesteRaportProblema, SalveazaRaportProblema, TrimiteRaportProblemaPeGitHub} from '../wailsjs/go/main/App';
import {main, raportare} from '../wailsjs/go/models';
import {Dialog} from './Formulare';
import {notifica} from './Toast';

export function DialogRaportare({ecran, inchide}: {ecran: string; inchide: () => void}) {
    const {t} = useTranslation();
    const [tip, setTip] = useState<'problema' | 'parere'>('problema');
    const [descriere, setDescriere] = useState('');
    const [pasi, setPasi] = useState('');
    const [cuJurnal, setCuJurnal] = useState(true);
    const [raport, setRaport] = useState<{titlu: string; text: string} | null>(null);
    const [salvat, setSalvat] = useState('');
    const [ocupat, setOcupat] = useState(false);
    const [eroare, setEroare] = useState('');

    async function ruleaza(f: () => Promise<main.RaspunsRaportProblema>, dupa: (r: main.RaspunsRaportProblema) => void) {
        setOcupat(true); setEroare('');
        try {
            const r = await f();
            if (r.eroare) setEroare(t(r.eroare.cod, r.eroare.param ?? {})); else dupa(r);
        } catch {
            setEroare(t('date.eroare'));
        } finally {
            setOcupat(false);
        }
    }

    const pregateste = () => ruleaza(
        () => PregatesteRaportProblema(new raportare.Cerere({tip, descriere, pasi, cu_jurnal: cuJurnal, ecran})),
        r => { setRaport({titlu: r.titlu, text: r.text}); setSalvat(''); });

    return <Dialog titlu={t('raportare.titlu_dialog')} ocupat={ocupat} inchide={inchide}>
        <div className="corp-dialog formular-raportare">
            {!raport ? <>
                <fieldset className="tip-raport">
                    <legend>{t('raportare.tip')}</legend>
                    {(['problema', 'parere'] as const).map(x => <label key={x} className="bifa">
                        <input type="radio" name="tip_raport" checked={tip === x} disabled={ocupat} onChange={() => setTip(x)}/>{t(`raportare.tip_${x}`)}
                    </label>)}
                </fieldset>
                <label>{t('raportare.descriere')}
                    <textarea rows={5} value={descriere} disabled={ocupat} maxLength={8000} onChange={e => setDescriere(e.target.value)}/>
                </label>
                {tip === 'problema' && <label>{t('raportare.pasi')}
                    <textarea rows={4} value={pasi} disabled={ocupat} maxLength={8000} onChange={e => setPasi(e.target.value)}/>
                </label>}
                <label className="bifa"><input type="checkbox" checked={cuJurnal} disabled={ocupat} onChange={e => setCuJurnal(e.target.checked)}/>{t('raportare.cu_jurnal')}</label>
                <p className="mesaj atentionare">{t('raportare.avertisment')}</p>
            </> : <>
                <label>{t('raportare.previzualizare')}
                    <textarea className="text-raport" rows={16} value={raport.text} disabled={ocupat} spellCheck={false}
                        onChange={e => { setRaport({...raport, text: e.target.value}); setSalvat(''); }}/>
                </label>
                <p className="text-discret">{t('raportare.github_indiciu')}</p>
                {salvat && <div className="mesaj succes" role="status">
                    {t('raportare.salvat', {cale: salvat})}{' '}
                    <button type="button" className="legatura" onClick={() => void DeschideFolderRapoarteProbleme()}>{t('raportare.deschide_folderul')}</button>
                </div>}
            </>}
            {eroare && <div className="mesaj eroare" role="alert">{eroare}</div>}
        </div>
        <footer className="subsol-dialog">
            {!raport ? <>
                <button type="button" className="buton secundar" disabled={ocupat} onClick={inchide}>{t('ui.anuleaza')}</button>
                <button type="button" className="buton primar" disabled={ocupat || !descriere.trim()} onClick={() => void pregateste()}>{t('raportare.pregateste')}</button>
            </> : <>
                <button type="button" className="buton secundar" disabled={ocupat} onClick={() => setRaport(null)}>{t('raportare.inapoi')}</button>
                <button type="button" className="buton secundar" disabled={ocupat || !raport.text.trim()}
                    onClick={() => void ruleaza(() => SalveazaRaportProblema(raport.text), r => setSalvat(r.cale))}>{t('raportare.salveaza')}</button>
                <button type="button" className="buton primar" disabled={ocupat || !raport.text.trim()}
                    onClick={() => void ruleaza(() => TrimiteRaportProblemaPeGitHub(raport.titlu, raport.text), () => notifica(t('raportare.github_deschis'), 'info'))}>{t('raportare.github')}</button>
            </>}
        </footer>
    </Dialog>;
}
