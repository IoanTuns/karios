# D027 · Reunirea ramurilor prin PR și build din main

Stare: propusă · Data: 2026-10-26 · Task: F02-T11

Ioan Tuns cere explicit reunirea tuturor ramurilor în main printr-un PR și construirea ultimei versiuni. Această cerere autorizează pushul ramurii de integrare, crearea PR-ului și merge-ul în main, fără reconfirmare.

Propunere de proces: ramură de integrare din F09-T08, care conține deja toate vârfurile locale cu excepția notei finale F07-T12. După fetch se verifică și ramurile remote; fiecare vârf trebuie să fie strămoș al integrării. Merge commit pentru PR, păstrând istoricul și identificarea tuturor taskurilor; fără squash/rebase sau ștergerea ramurilor. Modificările nesalvate nu se includ și nu se suprascriu.

Înainte de push se verifică faptul că origin este repo-ul privat corect. Nu se publică executabile sau date prin PR. Documentația și codul fac parte din cererea explicită de integrare; cheia aplicației rămâne în locația existentă, conform D011.

Verificările se fac în worktree izolat, cu temporare Wails proprii, fără curățarea bin și fără acces la datele utilizatorului. După testare se creează PR-ul, se face merge și se verifică SHA-ul main local/remote. Buildul final rulează prin build.ps1 din sursele main; versiunea crește numai dacă sunt necesare schimbări de cod sau de dependențe, conform D009. Executabilul se livrează sub nume nou.
