// Textele interfeței vin din cataloagele Go (internal/i18n), nu din fișiere
// separate, ca interfața, mesajele și rapoartele să folosească aceleași texte (D004).
import i18next from 'i18next';
import {initReactI18next} from 'react-i18next';
import {Catalog, Limbi, SeteazaLimbaInterfata} from '../wailsjs/go/main/App';
import {limbaInterfataSalvata} from './limba';

export async function initI18n(): Promise<void> {
    const limbi = await Limbi();
    const implicita = limbi[0];
    const resources: Record<string, {translation: Record<string, string>}> = {};
    for (const limba of limbi) {
        resources[limba] = {translation: await Catalog(limba)};
    }

    await i18next.use(initReactI18next).init({
        resources,
        lng: limbaInterfataSalvata(limbi) ?? implicita, // limba aleasă anterior (F03-T05, D103)
        fallbackLng: implicita,
        interpolation: {escapeValue: false}, // React escapează deja
    });

    document.documentElement.lang = i18next.language;
    SeteazaLimbaInterfata(i18next.language).catch(() => {});
    document.title = i18next.t('aplicatie.titlu');
}
