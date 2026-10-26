# Documentația Pontix

Orice idee, decizie și acțiune a proiectului se documentează aici, în același commit cu schimbarea pe care o descriu. Regulile de lucru pentru agenți și oameni sunt în [AGENTS.md](../AGENTS.md).

| Ce | Unde | Se actualizează când |
| --- | --- | --- |
| Idei | [idei.md](idei.md) | apare orice idee, chiar dacă nu se face acum |
| Roadmap | [roadmap.md](roadmap.md) | apare o funcționalitate nouă sau își schimbă starea |
| Funcționalități și taskuri | [functionalitati/](functionalitati/) | un task începe, se termină, se blochează sau apare unul nou |
| Decizii | [decizii/](decizii/) | se alege între variante (tehnologie, regulă de calcul, format, proces) |
| Jurnal | [jurnal.md](jurnal.md) | se face orice acțiune: cod, document, instalare, discuție cu contabilitatea |
| Versiuni | [versiuni.md](versiuni.md) | se schimbă codul, un raport sau baza de date (D009) |
| Ghid de utilizare | [primele funcționalități și demo](ghid-primele-functionalitati.md) | se schimbă fluxurile disponibile în aplicație |
| Ghid rapoarte | [definire, PDF și tipărire](ghid-rapoarte.md) | se schimbă modelele, generarea sau tipărirea rapoartelor |
| Configurare pontaj | [ore, ture peste miezul nopții și totaluri](ghid-configurare-pontaj.md) | se schimbă regulile orelor sau completarea grilei |
| Construire și livrare | [directoare separate pentru fiecare build](ghid-build.md) | se schimbă scriptul de build sau organizarea artefactelor |
| Design inițial | [Aplicația de pontaj — tehnologii și decizii.pdf](Aplicația%20de%20pontaj%20—%20tehnologii%20și%20decizii.pdf) | nu se mai modifică; schimbările față de el devin decizii noi |

## Reguli

- **ID-uri.** Funcționalități `F01`, taskuri `F04-T03`, decizii `D001`, idei `I001`. Un ID nu se refolosește, nici după anulare.
- **Stări pentru idei:** `nouă`, `în analiză`, `transformată`, `amânată`, `respinsă`. O idee nu se șterge.
- **Cine.** Tabelele de taskuri au coloana `Responsabil` (un singur om sau agent per task); jurnalul și ideile au coloana `Cine`.
- **Stări pentru taskuri și funcționalități:** `de făcut`, `în lucru`, `gata`, `blocat`, `amânat`, `anulat`. Un task `blocat` spune de ce.
- **Stări pentru decizii:** `propusă`, `acceptată`, `respinsă`, `înlocuită de Dxxx`. O decizie acceptată nu se rescrie; se scrie una nouă care o înlocuiește.
- **Legături.** Fiecare task trimite la deciziile și fișierele lui; fiecare decizie trimite la taskurile care o aplică.
- **Șabloane:** [funcționalitate](functionalitati/_sablon.md), [decizie](decizii/_sablon.md).
- **Limba documentației:** română.
