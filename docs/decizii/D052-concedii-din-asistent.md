# D052 · Cererile de concediu din asistentul AI

Stare: acceptată · Data: 2026-10-29 · Decisă de: Ioan Tuns · Taskuri: F14-T26

## Context

„Maria Daria merge în concediu în octombrie 26 până pe 31” a primit lista concediilor lunii („Nu este înregistrat niciun concediu… Tot personalul este activ și la dispoziție”). Cauzele:

- sinteza din cod răspundea la orice mesaj care conține „în concediu”, fără să deosebească întrebarea de comandă;
- asistentul nu avea nicio acțiune pentru absențe;
- răspunsul trăgea o concluzie falsă („tot personalul e la dispoziție”) din lipsa înregistrărilor.

Cererile de concediu au flux de aprobare ([D024](D024-cereri-concediu-aprobare.md)).

## Decizie (Ioan Tuns, 2026-10-29)

- O cerere de concediu din asistent se creează cu starea **„în așteptare”**.
- **Numărul de înregistrare e obligatoriu** pentru aceste cereri.
- Asistentul **deschide formularul de concediu** pentru persoana respectivă, precompletat cu tipul (CO, sau CM când se pomenește „medical”) și intervalul. Utilizatorul completează numărul de înregistrare și salvează prin calea obișnuită: aceleași validări, suprapuneri, jurnal și aplicare în pontaj. Asistentul nu salvează singur concediul.

### Recunoașterea (propunerea lui Claude)

- Comandă = nu e întrebare (nu începe cu „cine / câți / ce / există…” și nu se termină cu „?”), pomenește un concediu și are un verb de acțiune („merge, pleacă, intră, ia, are…”).
- Angajatul se găsește după marcă sau după numele complet, cu toleranța la greșeli din F14-T15 („Maria Daria” → Marin Daria); doi candidați la fel de apropiați sau niciunul duc la o întrebare de clarificare.
- Intervalul: primele două zile din mesaj, cu luna scrisă lângă zi, luna numită în mesaj sau luna de lucru.
- Cardul arată interpretarea; formularul deschis permite orice corectură înainte de salvare.

## Variante respinse

- **Salvarea directă din asistent** — ar ocoli validările și fluxul formularului; numărul de înregistrare oricum trebuie scris de om.
- **Cerere „aprobată” direct** — respinsă de Ioan Tuns.

## Consecințe

- Numărul de înregistrare e obligatoriu doar pentru cererile pornite din asistent. Dacă trebuie obligatoriu pentru toate cererile, e o schimbare separată (întrebare deschisă).
- Întrebările despre concedii primesc în continuare lista calculată din cod, fără concluzii despre disponibilitate.
