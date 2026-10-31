import i18next from 'i18next';

type IdentitateAngajat = {nume?: string; prenume?: string; angajat_principal_id?: string; garda?: boolean};

// Marcajul apartine afisarii; numele si prenumele salvate raman nemodificate.
export function numeAngajat(a?: IdentitateAngajat | null): string {
    if (!a) return '';
    const nume = [a.nume, a.prenume].filter(Boolean).join(' ').trim();
    return a.garda || a.angajat_principal_id
        ? i18next.t('ui.nume_subordonat', {nume, tip: i18next.t('rapoarte.subtip_garda')})
        : nume;
}
