# F07-T39 - Categorii editabile de absenta

Responsabil: Codex. Stare: gata. Ramura F07-T39-categorii-absente.

Utilizatorul precizeaza ca lista Categorii trebuie administrata, cu regulile de ore confirmate in F07-T38. Se adauga nomenclator cu cod stabil, denumire RO/EN, activ/inactiv si efect FARA_ORE/NORMA/PROGRAM; tipurile asociate mostenesc regula categoriei. Natura tehnica odihna/medical/alta ramane necesara prioritatii medicalului si concediului contractual, distincta de denumirea editabila.

Compatibilitate: migratie noua 0030, categorii initiale peste datele existente, fara schimbarea efectelor deja inregistrate. Editarea denumirii/starii este permisa; schimbarea naturii/efectului folosit in istoric este protejata, pentru alta regula se creeaza o categorie noua. Previzualizarea si accesul din Planificare se pastreaza. Fara reguli monetare, fara date reale, fara push/merge.

Etape: persistenta/serviciu si teste, administrare UI si selectie dinamica, traduceri/manual/versionare, verificari si build. Pas urmator: implementare.


## Implementare 0.92.24

Fila separata Categorii in Concedii si absente: creare, editare denumiri RO/EN, natura, efect, activ/inactiv; filtru Active implicit/Inactive/Toate. Formularul tipului selecteaza categoria configurata si afiseaza regula mostenita, fara un al doilea selector contradictoriu de efect.

Serviciul valideaza categoria si ignora natura/efectul trimise separat de client cand exista categorie_id. Editarea regulii nefolosite se propaga atomic; redenumirea si schimbarea starii sunt permise si pentru categoriile folosite. Regula/natura folosite in istoric sunt protejate. Categoriile inactive nu permit tipuri sau absente noi, dar istoricul existent ramane editabil si calculabil. Revizii si jurnal pentru configurari.

Migrarea 0030 nu schimba datele existente: separa categoriile pentru efectele diferite deja configurate, pastreaza natura si asociaza tipurile. Testul de upgrade din 0029 verifica explicit codul, efectul, inactivul, data modificarii si 420 minute creditate. Tipurile create prin apeluri vechi primesc asociere compatibila. Nu se elimina cele trei naturi tehnice, necesare concediului anual si prioritatii medicalului.

Filtrul, sumarul si raportul de absente folosesc denumirile configurate; rapoarte versionate pontaj13, detaliat14, individual13, absente10, garzi10. Referintele tehnice noi sunt separate; cele vechi nu se modifica. Manual 1.34 revizuit pentru configurarea categoriei inaintea tipului si accesul din Planificare.

Verificari anonime: categorie creata, regula mostenita chiar daca clientul transmite alte valori, redenumire si conflict de revizie, propagare inainte de utilizare, protectie dupa utilizare, dezactivare/reactivare, blocare absente noi pe categorie inactiva, prioritate medicala personalizata; raport cu numele personalizat RO/EN. Edge UX30 RO/EN: creare categorie, redenumire, alegere din tip, regula mostenita si filtre active/inactive. Regresia Pontaj-absente trecuta dupa fixarea explicita a perioadei octombrie 2026: proba veche depindea de luna curenta si esua pe duminica de 1 noiembrie, fara defect de produs. Nicio data reala accesata.

TypeScript/Vite, 24 teste frontend, 3 teste manuale, 7 teste livrare si documentatie 0 erori. Go: gofmt/vet si suita completa. Dependente directe npm si Go actuale; indirectele Go raman in F14-T21 pentru compatibilitate, fara instalari noi. D081 ramane propusa pentru reprezentarea tehnica; cerinta functionala vine explicit de la utilizator.

Pas urmator: proba operator si eventuala integrare aprobata. Fara push/merge. Buildul si integritatea pachetului se consemneaza in jurnal dupa verificare.
