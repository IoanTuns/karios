# F07-T32 - Concedii din Completare Pontaj

Responsabil: Codex. Ramura F07-T32-absente-din-pontaj. Aplicatie 0.92.8, manual 1.19. [D068](decizii/D068-absente-din-pontaj.md) propusa.

## Modificare

Buton Concediu / alt tip de absenta in Completare Pontaj. Formularul comun permite alegerea tipului configurat, aprobarii, inlocuitorului si documentelor pe selectia exacta. Se poate inlocui o zi cu program sau schimba tipul unei absente existente. Celulele afiseaza CO/CM/codul configurat. Salvarea este imediata, comuna cu evidenta absentelor; nu necesita o a doua salvare a schitei pentru acele zile.

Backend: conversie atomica cu verificarea valorilor/reviziilor originale, a lunii deschise si a validarilor de absente. O eroare anuleaza golirea schitei si absentele, inclusiv jurnalul. Fluxul obisnuit SalveazaAbsente pastreaza protectia contra programului existent. Pontajul efectiv nu este sters. Planificarea ramane distincta de realizare. Frontendul pastreaza modificarile nesalvate din afara selectiei, reseteaza istoricul local si respecta incadrarea ierarhica a numelor.

## Verificari

Teste backend: Z7 in CO aprobat, CO in CM, pastrarea celeilalte zile; rollback integral la tip invalid, conflict, luna inchisa, pontaj efectiv, selectie invalida sau absenta respinsa. Edge cu punte simulata RO/EN: program salvat in CO, anulare, eroare simulata, metadate/aprobare, CO in CM, modificari nesalvate in afara selectiei si vizibilitate in evidenta absentelor. Inspectie vizuala a formularului la 1024x640.

Scenariile UI au identificat afisarea orelor zero in locul codului de absenta, corectata. Sincronizarea testelor a fost adaptata la finalizarea salvarii/recalcularii si la selectia prin mouse; verificarile au trecut. Compilarea TypeScript a identificat filtrarea unei stari in afara formularului, limitata apoi la conversie. gofmt, go vet ./... si suita Go completa au trecut. Build oficial complet din 5afeb8f: `build/bin/pontix-0.92.8-20261031-145330-7d34feb3a6b848e38ffef9d30f5353cc/pontix.exe`. Manual 1.19 MD/HTML/PDF alaturi; manifest, hashuri executabil/manuale/arhiva si integritate ZIP verificate.

Nu sunt modificate calcule sau cazuri contabile de referinta. Fara migrare/dependente noi. Datele reale nu au fost deschise; proba interactiva WebView2 ramane la operator. Pentru revenirea la program se anuleaza absenta in evidenta, apoi se aplica programul. Pas urmator: review si proba operator; fara push/merge in main.

TypeScript/Vite, 18 teste frontend, 3 teste manual, 7 teste livrare si control documentatie trecute. Focus/cursor/zoom reverificate in Edge. Dependente verificate: Go 1.27.3 si npm actuale; 62 actualizari tranzitive Go urmarite separat in F14-T21.
