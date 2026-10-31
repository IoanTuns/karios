import {it,expect} from 'vitest';
import i18next from 'i18next';
import {numeAngajat} from './numeAngajat';
it('marcheaza garda prin contract sau asociere fara a modifica numele',async()=>{await i18next.init({lng:'ro',resources:{ro:{translation:{'ui.nume_subordonat':'{{nume}} ({{tip}})','rapoarte.subtip_garda':'Garda'}}}});const a={nume:'Exemplu',prenume:'Ana',garda:true};expect(numeAngajat(a)).toBe('Exemplu Ana (Garda)');expect(numeAngajat({...a,garda:false,angajat_principal_id:'principal'})).toBe('Exemplu Ana (Garda)');expect(numeAngajat({...a,garda:false})).toBe('Exemplu Ana');expect(a.nume).toBe('Exemplu');});
