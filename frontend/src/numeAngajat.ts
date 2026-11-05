import i18next from 'i18next';

type IdentitateAngajat = {nume?: string; prenume?: string; angajat_principal_id?: string; garda?: boolean; prefix?: Record<string, string> | null};

// Marcajul apartine afisarii; numele si prenumele salvate raman nemodificate.
// Grilele Pontaj si Planificare folosesc aceasta varianta, fara prefix (D089).
export function numeAngajatFaraPrefix(a?: IdentitateAngajat | null): string {
    if (!a) return '';
    const nume = [a.nume, a.prenume].filter(Boolean).join(' ').trim();
    return a.garda || a.angajat_principal_id
        ? i18next.t('ui.nume_subordonat', {nume, tip: i18next.t('rapoarte.subtip_garda')})
        : nume;
}

// prefixAngajat da titulatura (Dr., Biol. …) in limba interfetei, cu romana ca rezerva.
export function prefixAngajat(a?: IdentitateAngajat | null): string {
    const p = a?.prefix;
    if (!p) return '';
    return (p[i18next.language] || p.ro || p.en || '').trim();
}

// numeAngajat: numele afisat in liste si selectoare, cu prefixul inainte.
export function numeAngajat(a?: IdentitateAngajat | null): string {
    const nume = numeAngajatFaraPrefix(a);
    const prefix = prefixAngajat(a);
    return prefix && nume ? `${prefix} ${nume}` : nume;
}
