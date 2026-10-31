# D076 - Prioritatea publicarii in Pontaj

Stare: propusa. Data: 2026-10-31. Task: F11-T20.

## Context

Utilizatorul confirma prioritatea: Planificare, concedii/zile libere, concediu medical. Schitele ocupate nu trebuie sa blocheze programul planificat.

## Propunere tehnica

Publicarea inlocuieste implicit valorile de program/ore din pontaj_ciorne. Se elimina bifa optionala din interfata si se explica prioritatea. Campul Inlocuieste ramane in API pentru compatibilitate cu apelurile vechi, fara a mai limita inlocuirea. Tokenul, reviziile, auditul, protectia absentelor si lunilor inchise se pastreaza. Salvarea planului pregateste programul; publicarea confirma scrierea in Pontaj.

Utilizatorul a confirmat explicit si retragerea pontajelor efective. Inregistrarile raman in tabela cu marcaj logic de retragere, iar auditul pastreaza snapshotul integral inainte/dupa. Publicarea inlocuieste efectivele aceleiasi date/marci; concediile inlocuiesc efectivele intersectate ale marcii, medicalul ale tuturor marcilor asociate. Sunt verificate toate lunile intersectate de o tura. Tokenul include snapshoturile efective, astfel incat modificarea orelor dupa previzualizare sa ceara o previzualizare noua. Editarea manuala si verificarile AI pastreaza protectia suprapunerilor.

## Alternative

Pastrarea bifei pentru celule ocupate ar contrazice prioritatea confirmata. Ignorarea tokenului ar permite suprascrierea unor modificari concurente neprevizualizate.
