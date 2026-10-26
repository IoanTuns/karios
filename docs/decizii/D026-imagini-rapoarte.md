# D026 · Imagini în rapoarte și paginare

Stare: propusă · Data: 2026-10-25 · Task: F09-T08

Ioan Tuns precizează că fiecare raport poate avea una sau mai multe pagini și imagini în antet/subsol. Paginarea existentă continuă automat conținutul și repetă componentele pe fiecare pagină, cu numărul paginii și totalul.

Propunere: fiecare dintre cele trei zone ale antetului și subsolului poate conține o imagine PNG/JPEG, pe lângă textele traductibile. Editorul importă fișierul local, permite înlocuirea/eliminarea și configurarea înălțimii în mm. La randare proporțiile se păstrează, iar imaginea încape în lățimea zonei. Imaginile sunt comune limbilor modelului. Nu se folosesc URL-uri sau căi externe în definiție.

Imaginea se păstrează în definiția JSON din baza criptată, ca base64, împreună cu numele și înălțimea. Migrarea 0014 completează câmpurile opționale ale definițiilor existente, fără schimbarea textelor sau a opțiunilor acestora. Validare Go a conținutului efectiv: PNG/JPEG, maximum 2 MiB per imagine, maximum 4 milioane de pixeli și 4096 px pe axă; înălțime 5–35 mm. Editorul aplică aceleași limite de import. Nu se introduc biblioteci noi.

Cele trei rapoarte cresc la v2 pentru aspectul nou, cu referințe noi păstrând fișierele v1. Se verifică imagini în documente de o pagină și în antetul/subsolul fiecărei pagini a unui document lung, inclusiv transparență PNG și JPEG.

PDF-urile cu imagini sunt normalizate în memorie prin pdfcpu înainte de previzualizare/salvare. Verificarea independentă pypdf a semnalat chei XObject duplicate pentru imaginile repetate produse de maroto/gofpdf; rescrierea deduplică resursele, păstrând documentul și imaginile. Nu se introduc alte biblioteci sau fișiere intermediare.

Compatibilitate cu livrarea de probă 0.20.0: aceasta avea istoricul 0001–0011 și 0013, fără cererile 0012. Goose respinge implicit completarea unei migrări mai vechi. Numai pentru acest istoric exact, cu tabelul definițiilor prezent și ambele surse 0012/0013 în executabil, furnizorul permite aplicarea 0012 înainte de continuarea la 0014. Citirea stării rămâne fără modificarea datelor, iar stocarea păstrează backupul înainte de aplicare. Orice alt gol rămâne refuzat. Test separat de actualizare păstrează modelul și revizia lui; cazuri negative protejează celelalte istorice.

F07-T12 este acum comis (a132a27), deci se integrează în această ramură înaintea livrării comune. Checkoutul utilizatorului rămâne neatins. TEMP/TMP se izolează la generarea Wails, conform constatării I020. Versiunea vizată: 0.21.0, după integrarea tehnică 0.20.1.
