import type {datebaza} from '../wailsjs/go/models';
export type VerificareConfigurare=''|'functii'|'categorii'|'co'|'garda';
export function corespundeVerificare(a:datebaza.Angajat,v:VerificareConfigurare){switch(v){case 'functii':return !a.functie_id&&!a.functie;case 'categorii':return !a.categorie_program_id;case 'co':return !a.garda&&a.concediu_zile_an==null;case 'garda':return a.garda&&!a.angajat_principal_id;default:return true;}}
