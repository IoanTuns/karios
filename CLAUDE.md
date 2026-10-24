# Pontix

Aplicație desktop portabilă pentru Windows (un singur `.exe`) pentru pontajul lunar cu ture și tipărirea rapoartelor. Go + Wails + React + SQLite. Designul inițial și starea proiectului sunt în `docs/`.

## Documentare (obligatoriu)

Toate acțiunile și deciziile se documentează, în același commit cu schimbarea. Regulile complete sunt în `docs/README.md`.

- Orice acțiune → o intrare nouă, sus, în `docs/jurnal.md`.
- Un task început, terminat sau blocat → starea lui în `docs/functionalitati/Fxx-*.md`, iar starea funcționalității în `docs/roadmap.md` dacă s-a schimbat.
- Muncă nouă → un task nou (`Fxx-Tyy`) sau o funcționalitate nouă, înainte de a o face.
- O alegere între variante → o decizie nouă `docs/decizii/Dxxx-*.md` (șablon `_sablon.md`). O decizie acceptată nu se rescrie; se înlocuiește.
- ID-urile nu se refolosesc. Documentația e în română.

## Traduceri (obligatoriu)

Niciun text afișat utilizatorului nu se scrie direct în cod (D004).

- Text nou → cheie nouă în `internal/i18n/locales/ro.json` și `en.json`.
- Erori pentru utilizator → `calcul.Eroare{Cod, Param}` sau un tip asemănător, nu `fmt.Errorf` cu text.
- Denumiri introduse de utilizator → `i18n.Text`.

## Cod

- Identificatorii și comentariile sunt în română, fără diacritice în identificatori.
- Logica de calcul stă în `internal/calcul` și nu depinde de bază de date sau interfață.
- Înainte de commit: `gofmt -l -w .`, `go vet ./...`, `go test ./...`.
