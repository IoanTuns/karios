# Documentația Pontix

Toate acțiunile și deciziile proiectului se documentează aici, în același commit cu schimbarea pe care o descriu.

| Ce | Unde | Se actualizează când |
| --- | --- | --- |
| Roadmap | [roadmap.md](roadmap.md) | apare o funcționalitate nouă sau își schimbă starea |
| Funcționalități și taskuri | [functionalitati/](functionalitati/) | un task începe, se termină, se blochează sau apare unul nou |
| Decizii | [decizii/](decizii/) | se alege între variante (tehnologie, regulă de calcul, format, proces) |
| Jurnal | [jurnal.md](jurnal.md) | se face orice acțiune: cod, document, instalare, discuție cu contabilitatea |
| Design inițial | [Aplicația de pontaj — tehnologii și decizii.pdf](Aplicația%20de%20pontaj%20—%20tehnologii%20și%20decizii.pdf) | nu se mai modifică; schimbările față de el devin decizii noi |

## Reguli

- **ID-uri.** Funcționalități `F01`, taskuri `F04-T03`, decizii `D001`. Un ID nu se refolosește, nici după anulare.
- **Stări pentru taskuri și funcționalități:** `de făcut`, `în lucru`, `gata`, `blocat`, `amânat`, `anulat`. Un task `blocat` spune de ce.
- **Stări pentru decizii:** `propusă`, `acceptată`, `respinsă`, `înlocuită de Dxxx`. O decizie acceptată nu se rescrie; se scrie una nouă care o înlocuiește.
- **Legături.** Fiecare task trimite la deciziile și fișierele lui; fiecare decizie trimite la taskurile care o aplică.
- **Șabloane:** [funcționalitate](functionalitati/_sablon.md), [decizie](decizii/_sablon.md).
- **Limba documentației:** română.
