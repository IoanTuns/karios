# D105 · Ștergerea tabelelor și a codului de tarif din D003

Stare: acceptată · Data: 2026-11-12 · Taskuri: F04-T10 · Idee: I131 · Decisă de: Ioan Tuns, 2026-11-12 („Poți face I132 și I131”, după ce ștergerea a fost prezentată ca cerând aprobare)

## Context

D003 a pregătit împărțirea orelor pe categorii de tarif (zi, noapte, weekend, sărbătoare). Pregătirea cuprindea:

- tabelele `categorii_ore`, `seturi_reguli` și `reguli_interval`;
- funcțiile `calcul.Segmenteaza` și `SetReguli.Valideaza`.

Între timp aplicația a mers pe altă cale:

- Coloanele centralizatorului se definesc în Configurare pontaj (D092, D098).
- Zilele libere și ZNL folosesc calendarul.
- Pontix nu calculează sporuri și salarii (D100).

Analiza F07-T05 a arătat că nici tabelele, nici funcțiile nu sunt folosite. Le folosea doar un test al schemei.

## Decizie

- **Migrarea `0038_stergere_tarife.sql`** șterge tabelele `reguli_interval`, `seturi_reguli` și `categorii_ore`, cu indexul lor. Migrarea inversă le recreează goale.
- **Codul șters:**
  - fișierele `internal/calcul/segmenteaza.go` și `segmenteaza_test.go`;
  - din `internal/calcul/reguli.go`: tipurile `Interval` și `SetReguli`, cu validarea lor;
  - cele patru mesaje `reguli.*` ale validării.
- **Ce rămâne:** `calcul.TipZi`, `TipuriZi` și `Calendar`, folosite de asistentul AI, de absențe și de planificare.
- **Auditul (D044)** își reface sigiliile după migrare, ca după orice migrare, doar pentru tabelele existente.
- **D003 rămâne valabilă pentru principii.** Momentele se socotesc în fusul Europe/Bucharest, iar tipul zilei vine din calendar. Partea despre intervalele de tarif e înlocuită de această decizie.

## Variante respinse

- **Păstrarea tabelelor pentru o eventuală versiune cu tarife.** Nu există o astfel de cerință, iar D100 o exclude. O versiune viitoare ar porni oricum de la regulile primite de la contabilitate, nu de la schema din D003.
- **Păstrarea doar a codului.** Fără tabele și fără apelanți, codul ar rămâne o capcană pentru cine îl găsește.

## Consecințe

- Baza are trei tabele mai puțin. O bază veche se actualizează automat, cu backup înainte de migrare, ca de obicei.
- Testele de migrare numără o migrare în plus.
- Testele de traduceri folosesc ca exemplu de mesaj cu parametri `pontaj.schimbare_ora_luna_inchisa`, nu `reguli.gol`.
