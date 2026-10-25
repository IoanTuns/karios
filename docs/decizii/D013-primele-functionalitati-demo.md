# D013 · Primul flux de angajați și datele demo

Stare: propusă · Data: 2026-10-25 · Taskuri: F06-T03, F07-T01, F07-T07

## Cerință

Ioan Tuns a cerut primele funcționalități implementate și un set de date demo. Construim un flux complet de administrare a angajaților peste schema existentă, cu salvare persistentă.

## Propunere tehnică

- Interfață React cu CSS și elemente HTML native: navigare Angajați / Setări, tabel cu căutare și filtre, formular de editare, dezactivare/reactivare fără ștergere. Nu adăugăm o bibliotecă de componente.
- Accesul SQL se generează cu sqlc, deja ales prin D001. Versiune fixată: [v1.31.1](https://github.com/sqlc-dev/sqlc/releases/tag/v1.31.1), [licență MIT](https://github.com/sqlc-dev/sqlc/blob/v1.31.1/LICENSE). Configurare conform [documentației oficiale](https://docs.sqlc.dev/en/latest/reference/config.html). Rulăm prin `go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1 generate`, fără instalare globală și fără dependență sqlc în executabil.
- Validările și tranzacțiile stau în Go; metodele Wails returnează erori structurate cu cod traductibil. Interfața folosește cataloagele RO/EN existente.
- Setarea normei globale folosește minute întregi în stocare; interfața acceptă ore cu fracții care reprezintă minute întregi. Nu modificăm regulile de plată.
- Datele demo sunt un JSON inclus în aplicație, cu nume fictive și mărci `DEMO-…`, norme întregi/parțiale, concedii și stări active/inactive. Secțiile au denumiri `i18n.Text`.
- Încărcarea demo este explicită, tranzacțională și permisă numai dacă nu există angajați (inclusiv șterși logic) sau pontaje. Nu suprascrie date, nu schimbă norma globală și nu rulează automat. Interfața explică ce va fi adăugat înaintea acțiunii. Datele rămân între porniri.
- Modificările acestui flux sunt înregistrate în jurnal cu actorul provizoriu `local`; identificarea persoanei vine cu F15. F06-T04 rămâne pentru acoperirea celorlalte fluxuri.
- Verificarea UI folosește Edge headless deja instalat și protocolul CDP prin Node, fără biblioteci noi. Doar testul injectează o punte Wails simulată și date fictive într-un profil temporar; testele Go verifică separat backendul real și persistența criptată.

## Variante evaluate

- Date fictive doar în memoria interfeței: nu verifică salvarea și repornirea.
- Demo încărcat automat: ar amesteca exemple cu datele utilizatorului.
- Ștergere definitivă a angajaților: pierde istoricul; folosim activ/inactiv.

## Limite

Pontajul lunar, rapoartele și conturile rămân în taskurile existente. Această tranșă nu pretinde că implementează salarizarea sau întregul sistem de audit.
