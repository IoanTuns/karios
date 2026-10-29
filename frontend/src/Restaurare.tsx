// Restaurarea unui backup din aplicație (F05-T25): verificarea sigiliilor copiei, programarea
// pentru următoarea pornire și închiderea aplicației.
import {useEffect, useState} from 'react';
import {useTranslation} from 'react-i18next';
import {AnuleazaRestaurarea, InchideAplicatia, ListaBackupuri, ProgrameazaRestaurare, VerificaBackup} from '../wailsjs/go/main/App';
import {audit, main, stocare} from '../wailsjs/go/models';
import {Dialog} from './Formulare';

function dimensiune(octeti: number): string {
    return octeti >= 1048576 ? `${(octeti / 1048576).toFixed(1)} MB` : `${Math.max(1, Math.round(octeti / 1024))} KB`;
}

export function RestaurareBackup({reimprospatare, ocupat: ocupatExtern}: {reimprospatare: number; ocupat: boolean}) {
    const {t, i18n} = useTranslation();
    const [lista, setLista] = useState<stocare.Backup[]>([]);
    const [programata, setProgramata] = useState('');
    const [ales, setAles] = useState<{backup: stocare.Backup; rezultat: audit.Rezultat} | null>(null);
    const [ocupat, setOcupat] = useState(false);
    const [eroare, setEroare] = useState('');

    function aplica(r: main.RaspunsBackupuri): boolean {
        if (r.eroare) { setEroare(t(r.eroare.cod, r.eroare.param ?? {})); return false; }
        setEroare(''); setLista(r.backupuri ?? []); setProgramata(r.programata); return true;
    }

    useEffect(() => { ListaBackupuri().then(aplica).catch(() => setEroare(t('date.eroare'))); }, [reimprospatare]);

    async function ruleaza(f: () => Promise<void>) {
        setOcupat(true);
        try { await f(); } catch { setEroare(t('date.eroare')); } finally { setOcupat(false); }
    }

    const format = new Intl.DateTimeFormat(i18n.language, {dateStyle: 'short', timeStyle: 'short'});
    const blocat = ocupat || ocupatExtern;
    return <div className="restaurare-backup">
        <h3>{t('restaurare.titlu')}</h3>
        <p className="text-discret">{t('restaurare.descriere')}</p>
        {programata && <div className="mesaj atentionare" role="status">
            {t('restaurare.programata', {nume: programata})}
            <div className="actiuni-cont">
                <button type="button" className="buton secundar mic" disabled={blocat} onClick={() => void ruleaza(async () => { aplica(await AnuleazaRestaurarea()); })}>{t('restaurare.anuleaza')}</button>
                <button type="button" className="buton primar mic" disabled={blocat} onClick={() => void InchideAplicatia()}>{t('restaurare.inchide_acum')}</button>
            </div>
        </div>}
        {eroare && <p className="mesaj eroare" role="alert">{eroare}</p>}
        {lista.length === 0 ? <p className="text-discret">{t('restaurare.fara_backupuri')}</p> : <ul className="lista-backupuri">
            {lista.map(b => <li key={b.nume}>
                <span><strong>{format.format(new Date(b.modificat))}</strong> · {t(b.manual ? 'restaurare.manual' : 'restaurare.automat')} · {dimensiune(b.dimensiune)}<br/><small className="text-discret">{b.nume}</small></span>
                <button type="button" className="buton secundar mic" disabled={blocat} onClick={() => void ruleaza(async () => {
                    const r = await VerificaBackup(b.nume);
                    if (r.eroare) setEroare(t(r.eroare.cod, r.eroare.param ?? {})); else setAles({backup: b, rezultat: r.rezultat});
                })}>{t('restaurare.restaureaza')}</button>
            </li>)}
        </ul>}
        {ales && <Dialog titlu={t('restaurare.confirma_titlu')} ocupat={blocat} inchide={() => setAles(null)}>
            <div className="corp-dialog">
                <p>{t('restaurare.confirma_descriere', {data: format.format(new Date(ales.backup.modificat))})}</p>
                <p className={ales.rezultat.stare === 'incident' ? 'mesaj eroare' : 'text-discret'}>
                    {t(`restaurare.sigiliu_${ales.rezultat.stare}`, {tabele: (ales.rezultat.tabele ?? []).join(', ')})}
                </p>
                <p className="text-discret">{t('restaurare.siguranta')}</p>
                {eroare && <p className="mesaj eroare" role="alert">{eroare}</p>}
            </div>
            <footer className="subsol-dialog">
                <button className="buton secundar" disabled={blocat} onClick={() => setAles(null)}>{t('ui.anuleaza')}</button>
                <button className="buton primar" disabled={blocat || ales.rezultat.stare === 'incident'} onClick={() => void ruleaza(async () => {
                    if (aplica(await ProgrameazaRestaurare(ales.backup.nume))) setAles(null);
                })}>{t('restaurare.programeaza')}</button>
            </footer>
        </Dialog>}
    </div>;
}
