import {useTranslation} from 'react-i18next';
import {numeAngajatFaraPrefix} from './numeAngajat';
import {dataAfisata} from './formatare';

type Persoana = {id: string; marca: string; nume?: string; prenume?: string; angajat_principal_id?: string; garda?: boolean};

// Rezumatul selectiei din grila (R5): numele persoanelor si datele in limba interfetei,
// nu „mărci” si date AAAA-LL-ZZ. Peste 3 persoane se arata numarul celor ramase.
export function SelectieRezumat({persoane, deLa, panaLa}: {persoane: Persoana[]; deLa?: string; panaLa?: string}) {
    const {t} = useTranslation();
    const nume = persoane.slice(0, 3).map(p => `${numeAngajatFaraPrefix(p) || p.marca} (${p.marca})`).join(', ');
    const rest = persoane.length > 3 ? ' ' + t('ux30.selectie_si_altii', {numar: persoane.length - 3}) : '';
    const perioada = !deLa ? '—' : !panaLa || panaLa === deLa ? dataAfisata(deLa) : `${dataAfisata(deLa)} – ${dataAfisata(panaLa)}`;
    return <p className="ux30-selectie" role="status">{t('ux30.selectie_text', {numar: persoane.length, persoane: (nume || '—') + rest, perioada})}</p>;
}
