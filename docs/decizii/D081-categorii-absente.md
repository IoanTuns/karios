# D081 - Categorii administrabile ale absentelor

Stare: propusa. Data: 2026-11-02. Task F07-T39.

Utilizatorul precizeaza explicit ca lista de categorii trebuie editabila si sa detina regula de ore din F07-T38. Implementarea anterioara facuse editabile numai tipurile.

Propunere tehnica: nomenclator separat categorii_absente cu cod stabil, denumire RO/EN, natura tehnica, efect_ore, activ si revizie. Tipul refera categoria si preia natura/efectul, validate exclusiv in backend. Denumirile si starile sunt editabile; natura/efectul utilizate in pontaj ori absente istorice sunt protejate. Editarea regulii nefolosite se propaga atomic tipurilor asociate. O categorie inactiva nu permite tipuri noi sau absente noi; istoricul ramane editabil si calculabil.

Natura odihna/medical/alta este distincta de numele categoriei: pastreaza dreptul anual si prioritatea medicalului. Nu folosim denumirea libera pentru calcul. Nu eliminam prioritatea medicalului prin redenumire. FARA_ORE reduce obligatia, NORMA si PROGRAM crediteaza orele aprobate fara reducerea obligatiei, conform confirmarilor utilizatorului. Fara reguli de plata noi.

Migrarea 0030 asociaza fiecare tip existent unei categorii cu acelasi efect si aceeasi natura, fara schimbarea instantaneelor, reviziilor sau starilor tipului. Varianta ANTERIOR este numai de compatibilitate pentru istoricul existent. API-urile vechi sunt adaptate la aceeasi asociere; nu se pastreaza doua reguli contradictorii.

Alternative: lista fixa de categorii respinsa explicit de utilizator; natura dedusa din numele editabil respinsa pentru prioritate medicala; regula duplicata editabila si pe tip si pe categorie respinsa pentru inconsistenta; modificarea retroactiva a unei reguli folosite respinsa pentru protejarea istoricului.
