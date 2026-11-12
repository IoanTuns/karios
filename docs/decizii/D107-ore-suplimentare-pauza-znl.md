# D107 · Orele suplimentare, pauza și zilele din ZNL

Stare: acceptată (regulile) · propusă (detaliile tehnice) · Data: 2026-11-12 · Taskuri: F04-T04, F04-T05, F04-T11 · Decisă de: Ioan Tuns, 2026-11-12 (regulile, prin răspunsurile la chestionarul F01-T02 și la întrebările de clarificare)

## Context

Orele suplimentare (F04-T04) și scăderea pauzei (F04-T05) așteptau regulile de la contabilitate. Ioan Tuns a răspuns la chestionarul de ore (F01, 2026-11-12) și la patru întrebări de clarificare. Pontix nu calculează bani (D100): se numără doar ore.

## Decizie

### Regulile (alese de Ioan Tuns)

1. **Orele suplimentare** sunt orele lucrate într-o zi peste norma zilnică a acelei zile. Apar într-o coloană separată.
2. **În weekend și în sărbători** persoana nu are normă, deci toate orele lucrate sunt suplimentare. Ele apar și în ZNL; cele două coloane nu se exclud.
3. **Rândurile de gardă nu au ore suplimentare.** La gardă rămâne doar împărțirea pe zile lucrătoare / weekend și sărbători (Ore gardă 200%, Ore lucrătoare, D093).
4. **Pauza** se configurează pe fiecare program: dacă se scade și când are loc.
5. **ZNL:** programele care intră se aleg în Configurare pontaj, iar intrarea zilelor libere adăugate de utilizator (de exemplu o punte) e o opțiune.

### Detaliile tehnice (propuse)

**Ore suplimentare (coloana OS)**
- Pentru fiecare persoană și zi se iau minutele lucrate în acea zi: programe și ore scrise de mână, după repartizarea configurată (ziua începerii sau zile calendaristice). Creditele absențelor (FMP, absențele cu ore) nu intră.
- Norma zilei este norma zilnică valabilă în acea zi (istoricul D101). Într-o zi de weekend sau într-o zi liberă din calendar, norma e 0.
- OS = suma, pe zile, a lui max(0, minute lucrate − norma zilei).
- Valoarea e „—” (necalculată, cu motiv) când:
  - persoana nu are normă zilnică (contract fără normă sau cu normă lunară);
  - rândul e de gardă;
  - ziua are o valoare necalculabilă.
- Coloana OS se configurează în Configurare pontaj, ca celelalte (sursa „Ore suplimentare”, cu factor).
  - Apare în grila Pontaj după OL, în tabelul de totaluri din rapoartele SQL (`os`) și în Pontaj detaliat, unde se poate alege ca orice coloană.
  - Exportul Excel nu se schimbă acum.

**Pauza (pe program, în Setări › Programe)**
- Fiecare program cu pauză are un mod:
  - **Nu se scade:** pauza e doar informativă; durata rămâne întreagă.
  - **Din total:** pauza se scade din orele zilei în care începe tura.
  - **La o oră fixă:** de exemplu 12:00–12:30. Se scade din intervalul care o conține, deci contează la împărțirea pe zile și la ZNL.
  - **La mijlocul turei:** intervalul pauzei e centrat pe mijlocul turei.
- Programele care au deja pauză primesc la migrare modul **nestabilit** și rămân necalculate, ca până acum, până când utilizatorul alege modul. Datele vechi nu se schimbă pe ascuns.
- Pauza poate fi cel mult jumătate din durata programului.

**Zilele din ZNL**
- Fiecare zi din calendar are un tip: **sărbătoare legală** sau **zi liberă adăugată**.
- **Clasificarea zilelor existente** se face la migrare:
  - pentru anii 2025–2027, din calendarul legal inclus în aplicație;
  - pentru ceilalți ani, toate zilele rămân „sărbătoare legală”.
- **În Sărbători legale** se alege tipul pentru fiecare zi.
- **Opțiunea din Configurare pontaj** (ZNL): **zilele libere adăugate intră în ZNL** sau **nu intră**.
  - Implicit intră, ca până acum, deci totalurile existente nu se schimbă.
  - Ca la schimbarea orei (D102), opțiunea nu se poate schimba cât timp e închisă o lună care conține o zi liberă adăugată.
  - Weekendul și sărbătorile legale intră mereu.
- **Programele din ZNL:** se folosește sursa existentă „Ore din programele selectate, în ZNL”, cu programele alese (de exemplu cele de tip Tură și Noapte).
  - Configurația salvată nu se schimbă automat. Utilizatorul alege o dată programele în Configurare pontaj.

## Variante respinse

- **Normă aplicată și în weekend.** Ioan Tuns a ales „toate orele din weekend sunt suplimentare”.
- **Ore suplimentare la gărzi.** Respinsă: gărzile au doar împărțirea D093.
- **Doar programele de tip Tură în ZNL, fix.** Ioan Tuns a ales alegerea din Configurare pontaj.
- **Schimbarea automată a configurației ZNL existente.** Ar schimba pe ascuns totalurile deja raportate.

## Consecințe

- Coloana nouă OS intră în calcul, în grilă și în rapoarte. Rapoartele care folosesc calculul primesc versiune nouă (D009).
- Migrări noi:
  - modul pauzei în programe;
  - tipul zilei în calendar.
- Exportul pentru salarizare (F10-T01/T03) rămâne deschis: programul de salarizare nu e ales (E1).
