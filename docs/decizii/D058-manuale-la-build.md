# D058 · Manualele în livrarea fiecărei versiuni

Stare: propusă. Propusă de: Codex, 2026-10-30. F02-T18. Cerința utilizatorului: manualele se regenerează și se marchează pentru fiecare versiune și se salvează împreună cu executabilul.

Sursa Markdown se păstrează în repo, împreună cu ediția și versiunea revizuită editorial în JSON. Buildul refuză o sursă nerevizuită pentru versiunea curentă și produce documente noi MD/HTML/PDF cu versiune, commit, data buildului și ediție. Manifestul păstrează identitatea și hashurile. Nu se declară livrare completă dacă lipsește un format.

Se reutilizează generatorul standard Python din F02-T17 și Microsoft Edge disponibil pe Windows pentru imprimarea HTML în PDF, cu profil separat și timeout. Nu se adaugă biblioteci sau servicii externe. Lipsa uneltelor produce eroare explicită, fără instalare automată. Generarea se face la construire; aplicația portabilă nu depinde de Python.

Alternative: copierea PDF-ului anterior cu nume nou nu garantează instrucțiuni actuale; generarea numai la cerere poate omite documentele; un serviciu extern introduce rețea și distribuirea documentelor. Aceste variante nu sunt folosite. Actualizarea metadatelor nu înlocuiește revizuirea conținutului. Integrarea contextuală/AI din D057 rămâne separată. Această decizie tehnică rămâne propusă; cerința de livrare a manualelor este explicit autorizată de utilizator.
