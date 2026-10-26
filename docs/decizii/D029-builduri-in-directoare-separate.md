# D029 · Director separat pentru fiecare build și curățarea artefactelor

Stare: propusă · Data: 2026-10-26 · Task: F02-T12

Ioan Tuns cere curățarea fișierelor generate și câte un director pentru fiecare build nou. Inventarul Git nu conține executabile, build/bin sau dist; nu este necesară rescrierea istoricului. Legăturile Wails, resursele Windows și referințele rapoartelor rămân urmărite, fiind necesare compilării și verificării.

Propunere de implementare: build.ps1 creează un director unic în build/bin, cu versiune, dată/oră și identificator aleator pentru a evita suprascrierile inclusiv la două builduri simultane. În el se livrează executabilul verificat și un manifest cu versiune, commit și SHA256. Temporarele Wails sunt izolate în același director pentru a nu accesa datele din TEMP-ul comun. Scriptul nu copiază baze de date între builduri și nu curăță buildurile anterioare.

Curățarea cerută acum elimină numai artefacte locale generate și identificate în inventar. Se protejează directoarele cu date operaționale, checkouturile/worktree-urile, uneltele și fișierele urmărite de Git. Înainte de orice ștergere recursivă se verifică ținta absolută, apartenența la build/bin și absența datelor, legăturilor și checkouturilor. Fișierele blocate nu se forțează și aplicațiile nu se opresc. Inventarul și rezultatele curățării se păstrează în documentația taskului.

Înlocuiește mecanismul de ieșire comună din D017; păstrează interdicția opțiunii -clean și protejarea datelor portabile. Nu adaugă biblioteci și nu schimbă rapoarte, calcule sau schema bazei.
