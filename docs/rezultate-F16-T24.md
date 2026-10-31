# F16-T24 - Dimensiuni adaptive la pornire

Responsabil Codex. Aplicatie 0.91.1, manual 1.10. Ramura F16-T24-fereastra-adaptiva, peste livrarea 0.91.0.

Zona utila a monitorului principal este citita prin SystemParametersInfoW (SPI_GETWORKAREA), iar GetDpiForSystem furnizeaza scalarea. Dimensiunile logice sunt reduse cu 64 unitati pentru margini si decoratiuni. Limita initiala ramane 1280x800; minimele 1024x640 scad si ele cand ecranul nu le permite. Windows/Wails pastreaza centrarea, maximizarea si redimensionarea manuala. Daca API-ul nu furnizeaza zona utila, se porneste la 800x600 cu minime 640x400.

Nu se schimba rezolutia, scalarea Windows sau dimensiunea fonturilor. Nu sunt afectate baza de date si rapoartele. Nu se gestioneaza automat mutarea ulterioara pe alt monitor; Windows/Wails gestioneaza evenimentele DPI existente.

Verificari: teste pentru 800, 1024, 1280, 1366, 1920 si 3840 pixeli, scalari 100/150/200%, limite minime si rezerva pentru decoratiuni; gofmt si vet trecute. Suita Go completa (`go test -p 1 -timeout 120s ./...`) trecuta. Trei teste de generare manual si control documentatie trecute. API-urile Wails au fost verificate in sursa locala; manifestul executabilului declara permonitorv2. Dependente verificate: Go 1.27.3 si npm actuale; actualizarile tranzitive Go raman urmarite in F14-T21, fara dependente noi.

Build oficial 7084fc3 complet: `build/bin/pontix-0.91.1-20261031-114557-b515fe2547ba49f6a8b80fc0707f434b/pontix.exe`, manual 1.10 MD/HTML/PDF si ZIP. Manifest si hashuri verificate. Pas urmator: proba vizuala de pornire la rezolutia/scalarea utilizatorului. Nu s-au deschis datele reale. Fara push sau merge.
