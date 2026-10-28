# D037 · Planificare separată și publicare în pontaj

Stare: propusă · Data: 2026-10-26 · Propusă de: Codex · Decisă de: — · Task: F11-T01

Cerință explicită: gărzi și schimburi de după-amiază într-o secțiune separată, din care se publică datele în pontaj. Se păstrează D036: gărzile sunt alese explicit, fără rotații inventate.

Implementarea propusă este un ecran operațional **Planificare**, cu secțiunile **Gărzi** și **După-amiază**, perioadă lunară și înregistrări angajat/dată/program. Planul are persistență, revizii și jurnal proprii. Salvarea planului nu modifică pontajul.

Publicarea are previzualizare și confirmare în aplicație. Implicit completează numai celulele goale. Pentru schițele deja completate se poate solicita explicit înlocuirea, cu afișarea valorilor înainte/după; pontajele efective, absențele și suprapunerile nu se suprascriu. Datele și reviziile sunt reverificate tranzacțional la publicare. Repetarea publicării acelorași valori nu dublează orele. Schimbarea sau retragerea unei poziții din plan nu șterge automat pontajul anterior publicat.

Programele și persoanele se aleg dintre configurațiile existente; categoria activă a angajatului trebuie să permită programul. Se verifică întregul interval al turei, inclusiv zilele următoare și lunile vecine, pentru absențe și suprapuneri. Calendarul evidențiază weekenduri și libere fără a interzice gărzi planificate explicit în acele zile. Câmpurile planificate se publică în schița folosită deja de grila și rapoartele de pontaj; nu sunt o confirmare nouă a muncii efectiv realizate.
