# Roadmap

Prima versiune utilizabilă (V1) înseamnă pontajul lunar complet, calculat corect, cu centralizatorul tipărit și exportul pentru salarizare. Restul se adaugă peste ea, fără rescriere.

Actualizat: 2026-10-25

| ID | Funcționalitate | Etapă | Stare | Depinde de |
| --- | --- | --- | --- | --- |
| [F01](functionalitati/F01-cerinte-reguli-calcul.md) | Cerințe și reguli de calcul (cu contabilitatea) | V1 | în lucru | — |
| [F02](functionalitati/F02-fundatia-proiectului.md) | Fundația proiectului | V1 | în lucru | — |
| [F03](functionalitati/F03-traduceri.md) | Traduceri (minim română și engleză) | V1 | în lucru | F02 |
| [F04](functionalitati/F04-motor-calcul.md) | Motorul de calcul al orelor și tarifelor | V1 | în lucru | F01, F02 |
| [F05](functionalitati/F05-aplicatie-desktop-portabila.md) | Aplicația desktop portabilă | V1 | în lucru | F02 |
| [F06](functionalitati/F06-baza-de-date.md) | Baza de date și migrările | V1 | în lucru | F05 |
| [F07](functionalitati/F07-date-de-baza.md) | Date de bază | V1 | în lucru | F03, F06 |
| [F08](functionalitati/F08-grila-pontaj.md) | Grila lunară de pontaj | V1 | de făcut | F04, F07 |
| [F09](functionalitati/F09-rapoarte-pdf.md) | Rapoarte PDF | V1 | de făcut | F04, F08 |
| [F10](functionalitati/F10-export-excel.md) | Export Excel pentru salarizare | V1 | de făcut | F01, F04 |
| [F11](functionalitati/F11-planificare.md) | Planificare și planificat vs. realizat | V1.1 | de făcut | F08 |
| [F12](functionalitati/F12-inchidere-luna.md) | Închiderea lunii și jurnalul de modificări | V1.1 | de făcut | F06, F08, F15 |
| [F13](functionalitati/F13-autentificare-sincronizare.md) | Autentificare și sincronizare | Ulterior | amânat | F06 |
| [F14](functionalitati/F14-asistent-llm.md) | Asistent LLM local | Ulterior | amânat | F07, F08 |
| [F15](functionalitati/F15-conturi-acces.md) | Conturi și acces (auto-login, roluri, recuperare) | V1 | de făcut | F05, F06 |

## Ordinea de lucru

1. **Riscul întâi:** F01 și F04 merg în paralel. Regulile de plată sunt partea care poate ieși greșit; restul e muncă previzibilă.
2. **Scheletul portabil (F05)** se face devreme, ca să verificăm pe un PC curat lucrurile care țin de Windows.
3. **Prototipurile de grilă și PDF** (F08-T01, F09-T01) se fac înainte de a construi pe ele, fiindcă decid biblioteca.
4. Apoi F06 → F07 → F08 → F09 și F10, în ordinea din documentul de design.
