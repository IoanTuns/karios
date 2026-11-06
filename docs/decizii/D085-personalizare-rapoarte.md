# D085 - Prezentare si copii editabile de raport

Stare: acceptată. Responsabil Codex. Data2026-11-03. F09-T19 · Decisă de: Ioan Tuns, 2026-11-06 (F02-T25)

Aspectul (Markdown, poziționare liberă, desenarea gofpdf) este înlocuit de [D087](D087-rapoarte-sql-pdfme.md) și a fost șters în F09-T23; restul rămâne valabil.

Cerere: texte formatabile, inclusiv celule generate (confirmare utilizator), bold optional pentru sambata/duminica/sarbatori, interval de garzi care poate depasi luna pentru grupuri consecutive de zile libere. Exemplu furnizat28-30noiembrie si1decembrie2026; s-a cerut confirmarea aplicarii la ambele margini.

Propunere de prezentare: stiluri pe camp (bold, italic, marime, aliniere) si etichete de coloane, persistate in model. Valorile generate pot fi editate intr-o copie pregatita in server: se accepta doar modificari de texte/stiluri pe pozitii existente, niciodata un Document complet furnizat de client. Exporturile si PDF-ul folosesc aceeasi copie. Modificarile se inregistreaza in jurnal cu valori inainte/dupa, autor si moment; se marcheaza explicit copia personalizata, cu istoric nemodificabil din editor. Recalcularea regenerand modelul creeaza din nou valorile sursa; editarea raportului nu actualizeaza pontajul sau regulile salariale. Patch limitat ca dimensiune, validare atomica si token nou; originalul din cache nu se modifica. Fara biblioteca noua.

Boldul pentru zile speciale foloseste calendarul configurat si o coloana cu data completa, inclusiv daca este ascunsa in tabel. Stilul automat este separat de campurile Weekend/Sarbatori folosite la calculul orelor eligibile.

Utilizatorul confirma explicit grupul complet la ambele margini: noiembrie2026 include1decembrie; decembrie2026 include28-30noiembrie. Se extinde numai prin zile consecutive nelucratoare care ating prima/ultima zi a lunii, folosind weekendul si sarbatorile configurate. Nu se adauga alte sarbatori indepartate si nu se modifica totalurile contractuale ale lunii.

Implementarea limiteaza patch-ul la1000 campuri/2MiB,4000 caractere per camp. Stiluri pe camp/celula,6-40pt, aliniere L/C/R, bold/italic si mostenire explicita. Copia are identificator nou asociat jurnalului si istoricului din PDF; contextul cererii, valorile si stilurile inainte/dupa sunt retinute. Detaliile istoricului raman la final; marcajul scurt de copie personalizata poate aparea in nota. Salvarea exportului pastreaza copia; regenerarea reia sursa, iar editarea ad-hoc nu devine regula de calcul. Fonturi Go Italic/BoldItalic din pachetul existent si aceeasi licenta Go fonts; fara biblioteca noua.
