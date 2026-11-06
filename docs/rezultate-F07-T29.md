# F07-T29 - Compartiment

Responsabil: Codex. Ramura F07-T29-compartiment. Aplicatie 0.92.4, manual 1.15.

La cererea utilizatorului, Departament devine Compartiment in textele romanesti: formulare, filtre, import, rapoarte, legenda, jurnal afisat si ajutor contextual. Manualul romanesc foloseste noua denumire. Denumirile proprii introduse de operator nu sunt schimbate. Traducerea engleza ramane Department.

Sabloanele CSV/XLSX noi au antetul Compartiment, din catalogul i18n. Importul accepta atat Compartiment cat si vechile antete Departament/Department. AI recunoaste compartimentul si departamentul; cautarea ajutorului accepta ambele denumiri. Cheile traducerilor, campurile API si BD raman compatibile; fara migrare. Nu se modifica regulile de calcul. Rapoartele au versiuni noi si referintele anterioare sunt pastrate.

Verificari: gofmt, go vet ./... si suita Go completa trecute; pachetele afectate reverificate dupa modificarea importului/AI/ajutorului. Generarea si reimportul sabloanelor, aliasul AI si PDF-urile validate prin teste. 18 teste frontend, 3 teste manual si control documentatie trecute. Dependente verificate: Go 1.27.3 si npm actuale; 62 actualizari tranzitive Go urmarite separat in F14-T21. Fara dependente noi. Datele reale nu au fost accesate. Proba interactiva desktop nu a fost efectuata.

Build oficial complet din c3f812e: `build/bin/pontix-0.92.4-20261031-140835-e6b0ae3bb60b48bfac4274e31a87dcd9/pontix.exe`. Manual 1.15 MD/HTML/PDF alaturi; manifest, hashuri executabil/manuale/arhiva si integritate ZIP verificate. Fara push sau merge in main.
