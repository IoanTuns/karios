# D024 · Definiții de rapoarte și tipărire

Stare: propusă · Data: 2026-10-25 · Task: F09-T07

Ioan Tuns cere rapoarte cu trei componente — header, content, footer — definite separat pentru fiecare raport într-un meniu propriu. Lucrează simultan la alte modificări; implementarea rapoartelor se face într-un worktree izolat, pornit din 68c0681, fără schimbarea ramurii sale F07-T11.

## Model și interfață propuse

- Meniu **Definire rapoarte** pentru creare/editare/duplicare și dezactivare de modele, cu trei zone explicite Antet, Conținut, Subsol. Meniu **Rapoarte** pentru alegerea modelului, perioadei, limbii și filtrelor, previzualizare, generare PDF și deschidere pentru tipărire.
- Antetul și subsolul au texte traductibile stânga/centru/dreapta. Variabilele sunt limitate la perioadă, secție, departament, angajat și data generării. Conținutul selectează sursa de date și coloanele; fără executare de HTML, SQL sau șabloane arbitrare introduse de utilizator.
- Prima implementare folosește texte/câmpuri și coloane configurabile; întrebarea opțională despre editor vizual a fost trimisă utilizatorului. Setări de pagină per raport: A4/A3, orientare, margini, mărime font.
- Definițiile stau în baza criptată (`0012_definitii_rapoarte.sql`; 0011 aparține taskului paralel F07-T11). Revizia este incrementată în SQL, cu refuzul editărilor concurente și jurnal. Copierea unui model creează o definiție independentă.
- Surse inițiale: pontaj lunar din schițe și absențe, fișă individuală lunară și evidență anuală a absențelor. Documentele indică starea de schiță/reguli neaplicate; nu inventăm totaluri de salarizare sau interpretări ale codurilor.
- PDF-ul este generat local în Go, cu diacritice, antet/subsol repetate și numere de pagină. Aceeași structură de date alimentează previzualizarea. Subsolul tehnic obligatoriu D009 include versiunea raportului, revizia definiției, aplicația, regulile folosite, momentul și actorul local până la F15.
- Tipărirea se face din PDF-ul deschis în aplicația asociată, unde utilizatorul alege imprimanta și opțiunile. Nu trimitem automat documente către imprimantă în verificări. Exportul în afara folderului de date este numai la acțiunea explicită de salvare a utilizatorului.

## Biblioteci și licențe

Versiuni stabile verificate prin Go proxy la 2026-10-25:

- `github.com/johnfercher/maroto/v2` **v2.4.4**, licență **MIT**, verificată din LICENSE al versiunii. Este biblioteca propusă inițial în D001, validată acum prin prototipul A3 și teste de paginare.
- `golang.org/x/image` **v0.47.0**, licență **BSD-3-Clause**, verificată din LICENSE. Fonturile Go Regular/Bold din `font/gofont` sunt incluse în PDF; licența fonturilor se păstrează în notificările distribuției.

Surse: [maroto](https://github.com/johnfercher/maroto), [fonturile Go](https://go.googlesource.com/image/+/refs/tags/v0.47.0/font/gofont/), [runtime Wails v2](https://github.com/wailsapp/wails/tree/v2.17.0/pkg/runtime). Nu instalăm programe pe sistem și nu folosim servicii externe pentru datele rapoartelor.

## Lucru în paralel

Versiunea funcționalității este rezervată 0.19.0 pentru a evita suprapunerea cu 0.18.0 a categoriilor. Integrarea eventualelor commituri F07-T11 se face numai în ramura rapoartelor; checkout-ul utilizatorului rămâne neatins. Deciziile tehnice de aici rămân propuse.
