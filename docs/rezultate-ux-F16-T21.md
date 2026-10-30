# F16-T21: remedieri UX pentru pilot

Responsabil: Codex. Data: 2026-10-30. Aplicatie 0.88.15, manual editia 1.4. Baza: main 1006f87 si analiza F16-T20.

| Constatare | Rezultat |
| --- | --- |
| UX01 | Callbackurile AI verifica ocuparea si confirmarea schitei nesalvate; refuzul pastreaza perioada si celula. Concediul foloseste anul datei cererii. Se pastreaza confirmarea existenta; nu se introduce salvare implicita. |
| UX02 | Indicator neutru pentru persoanele cu inregistrari in lista filtrata; eliminate pragurile procentuale si verdictul optim/deficit. Tooltipul explica scopul si numarul persoanelor afisate este vizibil. |
| UX03 | Ghidurile fara model sunt anuntate explicit; configurarea tehnica este pliabila. |
| UX04 | La latimi sub 1200px drawerul se suprapune peste grila, cu inchidere vizibila si continut derulabil; formularul permite infasurare. |
| UX05 | Cautarea afiseaza un articol din rezultate; la zero rezultate nu mai ramane articolul vechi. Actiunile au pozitie sticky. |
| UX06 | ID-ul tehnic este in detalii; citarea din alta versiune este distincta de articolul inexistent si ofera deschiderea articolului curent. |
| UX07 | Shift+F10 si tasta ContextMenu deschid actiunile celulei, cu focus initial; Tab parcurge butoanele, Escape revine la celula. Rol group corespunde butoanelor native fara a pretinde un meniu ARIA cu sageti. |
| UX08 | Checklistul manualului include backup, cheia separata si diferenta dintre baza si folder complet. Test nou pentru copierea intregului folder fictiv, redeschidere, pastrarea datelor si a backupului, sursa intacta. Nu certifica Windows curat sau alt utilizator real. |

Validare: frontend build si 13 teste Vitest trecute; 3 probe manual trecute; UI Edge RO/light si EN/dark 1024x640 trecute, inclusiv doua refuzuri de navigare AI cu celula nesalvata, indicator neutru si Shift+F10/Escape/focus. Simulatorul foloseste numai propuneri si date fictive. gofmt, go vet si go test -p 1 -timeout 120s ./... trecute, inclusiv copiere completa, backup/restaurare, trecere la D054 si deblocare fara copia Windows.

Limite: fara studiu cu operatori, zoom 200% complet, Windows curat sau imprimare fizica. In Setari exista deja locatia copiei create; nu se adauga API sau stocare pentru o alta evidenta a backupurilor. Actualizarile Go tranzitive raman in F14-T21; Go 1.27.3 si npm sunt curente la verificare. Nu se schimba dependente sau reguli de plata.

Pas urmator: build oficial, apoi review si proba operator pe copie fictiva in Windows separat. Fara push sau merge in main pentru acest task.
