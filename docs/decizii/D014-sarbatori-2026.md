# D014 · Sărbătorile din 2026 și interfața calendarului

Stare: propusă · Data: 2026-10-25 · Taskuri: F07-T03

## Cerință

Ioan Tuns a cerut extragerea sărbătorilor publice legale din România pentru 2026, generarea înregistrărilor și interfață pentru administrarea lor.

## Surse verificate

- [Codul muncii, art. 139](https://legislatie.just.ro/Public/DetaliiDocument/128647): lista sărbătorilor și precizarea că sărbătorile mobile ale salariaților creștini urmează cultul acestora.
- [Legea 52/2023](https://legislatie.just.ro/Public/FormaPrintabila/00000G2MKDF31CH3J1Q0MB8DHTA18JV3): 6 și 7 ianuarie.
- [Învierea la Catedrala din Suceava, 12 aprilie 2026](https://www.arhiepiscopiasucevei.ro/bucuria-invierii-domnului-la-catedrala-arhiepiscopala-din-suceava-2/) și [calendarul Rusaliilor, Basilica, 31 mai 2026](https://basilica.ro/calendar-ortodox-rusalii-pogorarea-duhului-sfant/): datele ortodoxe. Pascalia ROEA și paginile lunare ale Patriarhiei nu au putut fi deschise la verificare; datele au fost confirmate în sursele primare ecleziastice de mai sus.

Verificate la 2026-10-25. Nu introducem propuneri legislative drept legi în vigoare, zile-punte recuperabile sau reguli de plată.

## Propunere tehnică

- Reutilizăm tabelul `sarbatori` și anul derivat din data completă. Migrarea de date `0004_sarbatori_romania_2026.sql` creează înregistrările solicitate fără schimbarea structurii tabelului; la actualizarea unei baze existente se aplică backup-ul deja implementat.
- Set JSON verificabil, cu date și denumiri RO/EN, inclus în executabil. Pentru calendarul ortodox sunt 17 sărbători nominale în 16 date distincte: 1 iunie coincide cu a doua zi de Rusalii, deci are o singură înregistrare cu denumire combinată.
- Migrarea adaugă numai date lipsă. În interfață există și completarea explicită din același set, tranzacțională și repetabilă fără duplicate; nu suprascrie înregistrări deja editate și arată calendarul utilizat înainte de încărcare.
- Formular pe fiecare an pentru adăugare și editare, validarea datei reale și a denumirii. Schimbarea datei unei înregistrări este o actualizare, nu ștergere/recreare. Nu adăugăm o acțiune de ștergere.
- Operațiile sunt jurnalizate cu actorul provizoriu `local`, conform fluxurilor existente. Modelul actual este un calendar comun; calendarele individuale pe cult cer modelare separată înainte de a fi folosite la calcule de plată.
- Fără biblioteci noi și fără dependență de internet la rularea aplicației.

## Precizare de domeniu

Ioan Tuns a precizat că dorește toate sărbătorile publice legale din România, nu doar exemplele de Paște și Rusalii. Setul național include toate datele fixe și sărbătorile mobile ortodoxe, etichetate explicit în interfață; acesta nu înlocuiește calendarele individuale pe cult. Nu este introdusă o regulă universală pentru toți salariații.
