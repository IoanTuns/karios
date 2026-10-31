# F16-T25 - Zoom global

Responsabil: Codex. Ramura F16-T25-zoom-interfata. Aplicatie 0.92.5, manual 1.16. [D065](decizii/D065-zoom-interfata.md) propusa.

## Comportament

Butoane minus/procent/plus in antetul fiecarui ecran: 75-200%, resetare la 100% prin procent; nivelul se memoreaza local. Ctrl+-/Ctrl++/Ctrl+0 functioneaza si cu dialog deschis. Interfata completa este redimensionata, independent de densitatea S/M/L a grilelor. Lipsa accesului la localStorage nu produce eroare.

Dimensiunile dialogurilor si panourilor raportate la viewport tin cont de factor. Pozitionarea meniului Pontaj si inaltimea panoului AI sunt convertite la coordonatele redimensionate. Meniul lateral si actiunile antetului raman accesibile prin derulare/reasezare. Exporturile PDF/Excel nu sunt modificate.

## Verificari si predare

TypeScript/Vite, 18 teste frontend, 3 teste manual, 7 teste livrare si control documentatie trecute. Edge cu punte simulata RO/EN: limite 75/200%, butoane dezactivate la limite, resetare din tastatura, aplicarea factorului si memorarea locala. Formularele/focusul si densitatea Pontaj/Planificare sunt reverificate. Proba 1024x768 la 200% verifica limitele controalelor si dialogurilor. gofmt, go vet si suita Go completa au trecut; testul pachetului aplicatie a fost repetat cu succes dupa ultimul build frontend. Build oficial complet din 221fdb1: `build/bin/pontix-0.92.5-20261031-141440-6d4168b3b5324eceb9573e6cdcb7d579/pontix.exe`. Manual 1.16 MD/HTML/PDF alaturi; manifest, hashuri executabil/manuale/arhiva si integritate ZIP verificate.

Dependente verificate: Go 1.27.3 si npm actuale; 62 actualizari tranzitive Go urmarite separat in F14-T21. Datele reale nu sunt accesate. Proba interactiva WebView2/desktop ramane la operator. Pas urmator: verificarea confortului de utilizare pe monitorul utilizatorului; fara push/merge in main.
