# D101 · Istoricul normei zilnice și al dreptului la concediu

Stare: acceptată (regulile) · propusă (detaliile tehnice) · Data: 2026-11-09 · Taskuri: F07-T43 · Idee: I010 · Decisă de: Ioan Tuns, 2026-11-09 (cele trei reguli de mai jos)

Implementare: implementată din 0.92.92; toate punctele verificate în cod pe 2026-11-12, pe 0.92.102 (Claude, la cererea lui Ioan Tuns). Starea deciziei nu se schimbă prin această notă.

## Context

Norma zilnică și dreptul anual la concediu erau o singură valoare pe angajat. O schimbare, de exemplu trecerea la jumătate de normă sau mai multe zile de concediu, rescria și lunile trecute la recalculare. I010 cere perioade de valabilitate.

## Decizie

### Regulile (alese de Ioan Tuns)

1. **Norma se aplică pe zile.** Fiecare zi lucrătoare folosește norma valabilă în ziua aceea.
   - Exemplu: 8 h până pe 14 martie și 4 h de la 15 martie dau în martie 8 h × zilele lucrătoare 1–14 + 4 h × zilele lucrătoare 15–31.
2. **Dreptul la concediu e proporțional** cu zilele calendaristice în care a fost valabilă fiecare valoare. Rezultatul se rotunjește la cel mai apropiat întreg (0,5 în sus).
   - Exemplu: 21 de zile până la 30 iunie și 25 de la 1 iulie dau 21 × 181/365 + 25 × 184/365 ≈ 23.
3. **Lunile închise nu se modifică.** O schimbare care începe într-o lună închisă sau înaintea ei e refuzată, fie că se adaugă, fie că se șterge. Luna se redeschide întâi, cu motiv, ca până acum.

### Detaliile tehnice (propuse)

- **Tabelul.** Schimbările stau în `istoric_contract(angajat_id, de_la, norma_minute, concediu_zile_an)`.
  - Un câmp gol înseamnă „neschimbat la această dată”; cel puțin unul trebuie completat.
  - Două schimbări ale aceluiași angajat nu pot avea aceeași dată.
- **Valorile din fișă** sunt cele de dinaintea primei schimbări. Fără schimbări, totul rămâne ca înainte.
- **Limitele normei din istoric.** Norma din istoric înlocuiește norma zilnică doar unde angajatul are normă zilnică. Contractele fără normă zilnică (gardă, lunară, fără normă) nu primesc normă din istoric.
- **Drept necompletat.** Dacă o parte a anului nu are drept de concediu completat, dreptul pe acel an rămâne „necompletat”. Lipsa datelor nu se tratează ca zero zile.
- **Unde se aplică istoricul:**
  - Ore lucrătoare, L FIN și L Început folosesc norma pe zile.
  - Lista angajaților afișează norma valabilă azi.
  - Precompletarea compară durata programului cu norma zilei.
  - Sumarul concediilor arată dreptul pe anul afișat.
  - În rapoartele SQL, `angajati.concediu_zile_an` este dreptul pe anul raportului, iar `angajati.norma_zilnica` e norma de la începutul lunii raportului. Se adaugă tabelul `istoric_contract`.
- **Salvarea.** Schimbările se salvează imediat din fișa angajatului. Intră în jurnalul de modificări și sub sigiliul de audit (D044).

## Variante respinse

- **Norma de la 1 sau de la sfârșitul lunii pentru toată luna.** Ioan Tuns a ales norma pe zile.
- **Dreptul de la 1 ianuarie sau cel mai nou, pe tot anul.** Ioan Tuns a ales calculul proporțional.
- **Acceptarea cu avertisment în lunile închise.** Respinsă: datele unei luni închise nu se schimbă fără redeschiderea ei.
- **Istoric complet al fișei**, cu dată pe toate câmpurile. E prea mult pentru nevoia de acum: I010 privește doar norma și concediul.

## Consecințe

- Recalcularea unei luni vechi folosește valorile de atunci, nu pe cele de azi.
- Trei rapoarte au versiune nouă (D009): Pontaj detaliat (25), Fișa individuală (22) și rapoartele flexibile (20). Pentru angajații fără istoric, valorile nu se schimbă.
- Un SQL cu `SELECT * FROM angajati` vede în `concediu_zile_an` dreptul pe anul raportului, nu valoarea din fișă.
