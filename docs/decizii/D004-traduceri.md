# D004 · Totul se poate traduce și configura

Stare: acceptată · Data: 2026-10-24 · Taskuri: toate din F03

> Folderul de traduceri de lângă `.exe` devine sursă de import: traducerile corectate se păstrează în baza criptată ([D008](D008-criptarea-datelor.md)).

## Context

Aplicația trebuie să fie disponibilă cel puțin în română și engleză. Cerința s-a extins în aceeași zi: totul să fie configurabil și traductibil, inclusiv rapoartele și datele introduse de utilizator.

## Decizie

- **Niciun text afișat nu e scris direct în cod**, nici în Go, nici în React.
- **Cataloage JSON** plate (cheie → text), câte unul pe limbă, numite după codul limbii (`ro.json`, `en.json`). Cele incluse în `.exe` stau în `internal/i18n/locales/`. Parametrii au forma `{{nume}}`, compatibilă cu react-i18next.
- **Un folder lângă `.exe`** poate corecta texte sau adăuga limbi noi, fără recompilare (`i18n.AdaugaFolder`).
- **Rezerva:** limba cerută → română → cheia însăși. Nu apare niciodată un text gol.
- **Erorile** au un cod stabil și parametri (`calcul.Eroare`), traduse la afișare.
- **Datele utilizatorului** (categorii, coduri de absență, sărbători, coloane de raport) au câte o denumire pe limbă (`i18n.Text`), salvată ca JSON în SQLite.
- **Limba interfeței și limba rapoartelor** sunt setări separate.
- **Codurile interne** (`LUCRATOARE`, `NOAPTE`, `C`, `CM`) rămân aceleași în toate limbile; se traduce doar denumirea afișată.
- Română și engleză trebuie să fie mereu complete. Testul `TestCataloageIncluseComplete` verifică asta.

## Variante respinse

- **Fișiere `.po` / gettext** — mai greu de folosit direct din React.
- **Coduri scurte diferite pe limbă** (de exemplu `CM` / `SL`) — ar complica importul și exportul. Se poate reveni dacă e cerut.

## Consecințe

- Orice text nou înseamnă o cheie nouă în `ro.json` și `en.json`.
- O limbă adăugată de utilizator poate fi incompletă; `i18n.Lipsa` arată ce lipsește.
