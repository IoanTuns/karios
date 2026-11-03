# F09-T15 - Explicatii optionale

2026-11-02. Responsabil Codex. Ramura F09-T15-explicatii-optionale. Aplicatie 0.92.27, manual 1.37.

## Comportament

In Setari > Definire rapoarte > Continut: Afiseaza modul de calcul si semnificatia coloanelor. Alegere per model, persistata si pastrata la duplicare. Un model nou porneste debifat. Modelele vechi isi pastreaza afisarea pana la debifare si salvare. Explicatiile, formulele si SQL nu apar in PDF/Excel cand optiunea este debifata; titlul sectiunii de calcul este de asemenea omis. Istoricul tehnic ramane la final, separat. Nu se schimba datele/calculul sau baza de date. La un model flexibil modificarea necesita previzualizare inaintea salvarii.

JSON continut.ascunde_explicatii optional, fara migrare. Rapoarte pontaj14, detaliat15, individual14, absente11, garzi11, flexibil2. Referinte tehnice noi, pastrate cele vechi.

## Verificari

- Probe Go anonime RO/EN pentru toate cele cinci surse existente si sursa flexibila: afiseaza/ascunde/reafiseaza, persistenta, date identice, PDF si Excel. Trecute.
- TypeScript/Vite si 24 teste frontend, trecute.
- Edge RO/EN: model nou debifat, bifare si salvare persistenta, verificare definitie curenta; fara erori JavaScript. Trecute.
- Verificare documentatie 0 erori, 3 teste manuale si 7 teste livrare, trecute.
- gofmt/vet: trecute. Suita Go finala: timeout in conturi, restul pachetelor parcurse inclusiv rapoarte trecute; reluarea separata a trecut in 19.116s, apoi intreaga suita Go a trecut (rezultate in cache pentru pachetele deja verificate). Durata raportata de testul de autentificare peste doua ore, fara modificarea codului/testului pentru incident. Imprimanta fizica si datele reale nu au fost accesate.
- Dependente: Wails stabil2.17.1, actualizat separat in 0.92.26 si verificat cu intreaga suita Go/vet. Excelize2.11.1 ar impune efp prerelease; pastrat2.11.0/efp0.0.2, motiv consemnat. Go1.27.3 si npm directe actuale; actualizarile indirecte raman in F14-T21. CLI Wails2.17.1 se ruleaza temporar cu go run, fara instalare pe sistem.

## Predare

Pas urmator: proba operator si review pentru integrare. Fara push sau merge. Buildul oficial si verificarile manifestului/manualelor se consemneaza dupa compilare.
