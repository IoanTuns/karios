# Ajutor în fiecare meniu și manual în AI

F02-T17 · D057 propusă · 30 octombrie 2026.

Stare: primele două etape sunt implementate în F02-T21 (0.88.13): catalog comun, panou contextual și răspunsuri deterministe din AI. Reformularea cu LLM și tururile vizuale rămân extensii separate.

## Cum ar arăta pentru utilizator

Fiecare ecran are „Cum să faci”, în același loc în antet. Se deschide un panou lateral cu sarcinile ecranului, pașii numerotați, rezultatul așteptat și problemele frecvente. De exemplu, în Concedii: „Înregistrează o cerere”, „Aprobă zilele selectate”, „Atașează documentul”, „Anulează o cerere”. Un link „Citește în manual” deschide capitolul complet. Un link „Întreabă asistentul” trimite întrebarea și ID-ul articolului, nu datele persoanei selectate.

În AI: „Cum înregistrez un concediu?” primește pașii din articol, cu sursa. „Înregistrează concediul pentru EX-001” urmează fluxul existent de cerere și formular, cu validare și confirmare. Ajutorul nu transformă o întrebare într-o scriere de date și nu inventează restricții de rol.

## Acoperirea meniurilor

| Ecran / zonă | Sarcini inițiale | Capitol din manual |
| --- | --- | --- |
| Pornire și contul meu | Acces, recuperare, schimbare parolă | 2 |
| Angajați | Adăugare, import, categorie, dezactivare | 4 |
| Pontaj | Program, ore manuale, precompletare, salvare | 7 |
| Planificare | Tură specială, publicare, comparație | 8 |
| Concedii și absențe | Cerere, aprobare, document, anulare | 6 |
| Rapoarte | Generare, PDF, Excel, tipărire | 10 |
| Asistent AI | Comandă clară, concediu, context, motor | 11 |
| Instrumente AI | Șablon, verificare, export, anulare | 12 |
| Setări generale | Normă, secții, departamente | 5 și 9 |
| Schimburi | Interval, durată, regim, activare | 5 |
| Categorii de program | Permise, implicit, atribuire | 5 |
| Zile libere | An, import, completare listă | 5 |
| Configurare pontaj | Ture peste miezul nopții, coloane | 9 |
| Definire rapoarte | Antet, coloane, subsol, pagină | 10 |
| Jurnalul modificărilor | Autor, diferențe, filtrare, export | 15 |
| Conturi | Roluri, recuperare, dezactivare | 2 și 15 |
| Backup și restaurare | Copie, verificare, repornire | 14 |
| Închiderea lunii | Verificare, motiv, redeschidere | 13 |

## Formatul articolului

Fiecare articol ar avea: `id`, `ecrane`, `intrebari`, `titlu`, `conditii`, `pasi`, `rezultat`, `probleme`, `versiune_minima`, `articole_asociate` și texte pe limbi. Exemplu de ID stabil: `absente.inregistreaza`. Condițiile se verifică din starea aplicației; descrierea lor în document nu este un control de acces.

## Ordinea recomandată

1. Catalog local de articole și panou contextual pe toate ecranele, cu căutare și acces din tastatură.
2. Rutare deterministă a întrebărilor „cum fac” către articole, cu link spre ecran și sursă, inclusiv fără motor AI pornit.
3. Reformulare opțională cu modelul local, folosind numai câteva extrase relevante și un buget limitat de context.
4. Tururi vizuale pentru fluxurile importante, fără apăsarea automată a butoanelor de salvare sau confirmare.

## Criterii de acceptare pentru taskul de integrare

- Fiecare ecran și subsecțiune are cel puțin un articol potrivit.
- Un ID necunoscut sau un articol pentru altă versiune nu generează instrucțiuni inventate.
- Întrebările despre utilizare nu execută comenzi și nu cer inutil motorul LLM.
- Răspunsul AI indică articolul și versiunea, iar întrebările fără acoperire sunt semnalate.
- Căutarea tratează diacriticele și variante uzuale de exprimare; testele acoperă confuzia dintre cerere de concediu și întrebare despre concedii.
- Panoul funcționează cu tastatura, la 1024×640 și în ambele teme; textele noi folosesc i18n RO/EN.
- Documentația nu conține nume reale, parole, chei sau copii de documente ale angajaților.
- Un test verifică existența articolelor pentru toate valorile de ecran din navigare și absența linkurilor rupte.

## Implementare F02-T21

F02-T21 implementează catalogul comun, panoul contextual și rutarea deterministă a întrebărilor de utilizare. Articolele sunt disponibile RO/EN în aplicație și incluse în manualul românesc la build. Nu se implementează în acest task reformularea opțională cu LLM sau tururile care evidențiază controale.
