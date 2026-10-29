# F12 · Închiderea lunii și jurnalul de modificări

Etapă: V1.1 · Stare: de făcut · Depinde de: F06, F08, F15

## Scop

O lună închisă nu se mai poate modifica și calculul ei rămâne înghețat. Orice modificare și orice redeschidere se văd în jurnal.

## Taskuri

| ID | Task | Stare | Responsabil | Note |
| --- | --- | --- | --- | --- |
| F12-T01 | Închiderea și redeschiderea lunii; înghețarea pontajului și planificării împotriva modificărilor | gata | Antigravity | Implementat în 0.51.0: migrarea 0019 `stare_luni`, protecție în `SalveazaSchita` și `SalveazaPlanificare` (`pontaj.luna_inchisa`), banner vizual în `Pontaj.tsx` și buton de închidere/redeschidere cu confirmare. |
| F12-T02 | Ecranul jurnalului de modificări | gata | Claude | 0.65.0: ecran în Setări, filtre, diferențe pe câmpuri, export CSV. Vine din F06-T04 (tabelul complet) |
| F12-T03 | „Cine” în jurnal = contul din aplicație | gata | Claude | Vine din F15-T08 |
| F12-T04 | Amprenta lunii închise: calculată la închidere, tipărită în subsolul rapoartelor, verificată în pontaj | în lucru | Claude | [D044](../decizii/D044-jurnal-audit-sigilat.md) punctul 6. 0.60.0: `stare_luni.amprenta`, subsolul rapoartelor lunare (pontaj v4, detaliat v3, fișă v4, gărzi v2), afișare în bannerul lunii închise. Rămas: proba vizuală în PDF și în aplicația pornită |
| F12-T05 | Motivul scris de utilizator la închiderea și redeschiderea lunii, în locul textelor fixe din cod | gata | Claude | 0.66.0: dialog în Pontaj, motiv obligatoriu și în backend, avertizare despre amprentă la redeschidere |
