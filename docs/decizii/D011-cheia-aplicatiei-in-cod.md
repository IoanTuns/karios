# D011 · Cheia aplicației în cod (deocamdată)

Stare: acceptată · Data: 2026-10-25 · Decisă de: Ioan Tuns · Taskuri: F05-T12

Înlocuiește parțial [D008](D008-criptarea-datelor.md) punctul 2 („cheia aplicației nu intră în repo”) și face o excepție de la [AGENTS.md](../../AGENTS.md) secțiunea 4 („niciun secret în repo”).

## Context

Seiful (F05-T10) are nevoie de o cheie a aplicației, aceeași la toate versiunile. D008 cerea ca ea să fie injectată la build dintr-un fișier păstrat în afara repo-ului, de un om anume. Rămânea de decis cine o păstrează.

## Decizie

- **O singură cheie, în cod** (`seif.CheieAplicatie`), generată aleator pe 2026-10-25.
- **Scopul, în cuvintele lui Ioan Tuns:** datele să nu poată fi văzute fără executabil.
- Cheia nu se schimbă. Altă cheie ar face ilizibile toate folderele `karios/` existente.
- Variabila poate fi încă înlocuită la build (`-ldflags`), deci mutarea cheii în afara codului rămâne posibilă ([I008](../idei.md)).

## Ce protejează și ce nu

| Situație | Protejat? |
| --- | --- |
| Folderul `karios/` deschis cu alt program sau copiat fără `.exe` | **da** |
| Cineva care are `.exe` | **nu**, la fel ca în D008 |
| Cineva care are codul sursă (repo-ul) | **nu**: vede cheia direct |

## Variante respinse

- **Cheie în afara repo-ului, injectată la build** (D008) — amânată: cere o persoană care păstrează cheia și un pas în plus la build.

## Consecințe

- Repo-ul trebuie tratat ca privat: cine îl primește poate decripta orice folder `karios/`.
- Fără gestionarea unui fișier-cheie, orice build (inclusiv al altui agent) poate deschide datele existente.
- Mutarea cheii mai târziu: versiunea nouă citește seiful cu cheia veche și reîmpachetează copia „A” cu cheia nouă, fără să recripteze baza.
