// Autorul din jurnal e „local”, numele contului sau „<cont>/asistent_ai” (F14-T11, F15-T08).
import type {TFunction} from 'i18next';

export function autorJurnal(t: TFunction, cine: string): string {
    const [baza, sursa] = cine.split('/');
    const nume = baza === 'local' ? t('jurnal.cine_local') : baza;
    return sursa === 'asistent_ai' ? t('jurnal.cine_asistent', {cine: nume}) : nume;
}
