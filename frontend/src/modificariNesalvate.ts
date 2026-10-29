// Ecranele cu modificări nesalvate se raportează aici; aplicația Go întreabă la închiderea
// ferestrei dacă există vreuna (F05-T24). `beforeunload` nu oprește închiderea în Wails.
import i18next from 'i18next';
import {SeteazaModificariNesalvate} from '../wailsjs/go/main/App';

const surse = new Set<string>();

// semnaleazaModificari se folosește într-un useEffect; curățarea retrage sursa.
export function semnaleazaModificari(sursa: string, da: boolean): () => void {
    actualizeaza(sursa, da);
    return () => actualizeaza(sursa, false);
}

function actualizeaza(sursa: string, da: boolean) {
    const inainte = surse.size > 0;
    if (da) surse.add(sursa); else surse.delete(sursa);
    if (inainte !== surse.size > 0 || da) SeteazaModificariNesalvate(surse.size > 0, i18next.language).catch(() => {});
}
