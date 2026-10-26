# Curățarea artefactelor · 2026-10-26

F02-T12, cerut de Ioan Tuns. Au fost eliminate 36 de ținte generate, însumând **357420297 octeți (340,86 MiB)**. Inventarul de nume, dimensiuni și momente ale fișierelor din directoarele protejate a rămas identic înainte și după curățare. Nu au existat ținte blocate.

Inventarul Git al tuturor ramurilor nu conținea executabile sau build/bin/dist; istoricul nu a fost rescris. Checkouturile, uneltele, baza veche din build/bin/karios, resursele urmărite și referințele rapoartelor au fost păstrate. Temporarele șterse erau cele create exclusiv de generatorul Wails în sesiunile anterioare. Înaintea ștergerii s-au verificat ținta absolută în directoarele de artefacte, legăturile, fișierele blocate și absența datelor operaționale. Datele nu au fost deschise.

Țintele de mai jos sunt relative la rădăcina repo-ului. Pentru directoarele protejate nu a fost calculat un volum de șters.

| Țintă | Tip | Rezultat | MiB eliminate |
| --- | --- | --- | --- |
| `build/bin/pontix-0.13.0` | director | protejat: conține date | — |
| `build/bin/pontix-0.20.0-rapoarte` | director | protejat: conține date | — |
| `build/bin/pontix-0.21.0-integrare` | director | șters | 25,44 |
| `build/bin/pontix-0.21.0-rapoarte` | director | protejat: conține date | — |
| `build/bin/pontix-0.22.1-rapoarte` | director | protejat: conține date | — |
| `build/bin/calendar-2026.png` | fisier | șters | 0,05 |
| `build/bin/concedii-absente.png` | fisier | șters | 0,05 |
| `build/bin/interfata-demo.png` | fisier | șters | 0,10 |
| `build/bin/pontaj-schita.png` | fisier | șters | 0,06 |
| `build/bin/pontix-0.13.0.exe` | fisier | șters | 18,00 |
| `build/bin/pontix-0.15.0.exe` | fisier | șters | 18,03 |
| `build/bin/pontix-0.16.0.exe` | fisier | șters | 18,50 |
| `build/bin/pontix-0.17.1.exe` | fisier | șters | 18,51 |
| `build/bin/pontix-0.18.0.exe` | fisier | șters | 18,52 |
| `build/bin/pontix-0.19.0.exe` | fisier | șters | 18,55 |
| `build/bin/pontix.exe` | fisier | șters | 18,55 |
| `build/bin/lucru-rapoarte/build/bin/pontix-0.22.2-20261026-094422-1a13ef05517c4123ab42d1dbedbfb2b8` | director | șters | 25,47 |
| `build/bin/lucru-rapoarte/build/bin/rapoarte-demo` | director | șters | 0,97 |
| `build/bin/lucru-rapoarte/build/bin/temporar-wails-f09t08` | director | șters | 1,96 |
| `build/bin/lucru-rapoarte/build/bin/temporar-wails-f09t09` | director | șters | 1,80 |
| `build/bin/lucru-rapoarte/build/bin/temporar-wails-rapoarte` | director | șters | 1,73 |
| `build/bin/lucru-rapoarte/build/bin/temporar-wails-reunire` | director | șters | 1,77 |
| `build/bin/lucru-rapoarte/build/bin/actualizare-f09t09.py` | fisier | șters | 0,01 |
| `build/bin/lucru-rapoarte/build/bin/calendar-2026.png` | fisier | șters | 0,05 |
| `build/bin/lucru-rapoarte/build/bin/concedii-absente.png` | fisier | șters | 0,05 |
| `build/bin/lucru-rapoarte/build/bin/configurare-pontaj.png` | fisier | șters | 0,06 |
| `build/bin/lucru-rapoarte/build/bin/definire-rapoarte.png` | fisier | șters | 0,04 |
| `build/bin/lucru-rapoarte/build/bin/interfata-demo.png` | fisier | șters | 0,11 |
| `build/bin/lucru-rapoarte/build/bin/module-f02t12.json` | fisier | șters | 0,05 |
| `build/bin/lucru-rapoarte/build/bin/pontaj-schita.png` | fisier | șters | 0,06 |
| `build/bin/lucru-rapoarte/build/bin/pontix-0.19.1.exe` | fisier | șters | 25,29 |
| `build/bin/lucru-rapoarte/build/bin/pontix-0.20.0.exe` | fisier | șters | 25,29 |
| `build/bin/lucru-rapoarte/build/bin/pontix-0.21.0-integrare.exe` | fisier | șters | 25,39 |
| `build/bin/lucru-rapoarte/build/bin/pontix-0.21.0.exe` | fisier | șters | 25,39 |
| `build/bin/lucru-rapoarte/build/bin/pontix-0.22.0-rapoarte.exe` | fisier | șters | 25,47 |
| `build/bin/lucru-rapoarte/build/bin/pontix-0.22.1-rapoarte.exe` | fisier | șters | 25,47 |
| `build/bin/lucru-rapoarte/build/bin/previzualizare-raport.png` | fisier | șters | 0,05 |
| `build/bin/lucru-rapoarte/build/bin/proba-build-f02t12.json` | fisier | șters | 0,00 |
| `build/bin/lucru-rapoarte/build/bin/referinte-f09t09.py` | fisier | șters | 0,00 |
| `build/bin/lucru-rapoarte/build/bin/ui-f09t09.py` | fisier | șters | 0,00 |
