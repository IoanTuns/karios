# Pontix {{VERSIUNE}} - note pentru utilizatori

Commit livrare: {{COMMIT}}.

## Despre această versiune

Pontix 1.0 este prima versiune pentru folosire curentă. A trecut proba de acceptanță pe un PC Windows curat, de pe un stick USB și la actualizarea de la o versiune veche cu date.

Pontix ține evidența orelor și a zilelor și le predă programului de salarizare; nu calculează bani.

## Ce face

- **Angajații:** contracte, norme și zile de concediu, cu istoric; secții, funcții, grupuri profesionale; mărci separate pentru gardă.
- **Pontajul lunar:**
  - completare în lot, precompletare din program;
  - reguli configurabile pentru coloanele de total;
  - orele suplimentare (OS), ZNL, pauza pe program, gărzile, nopțile cu schimbarea orei.
- **Concediile și absențele:** cereri, aprobare, vedere pe zile și pe cereri.
- **Planificarea turelor:** reguli, publicarea în Pontaj, comparația dintre planificare și pontaj.
- **Rapoartele:**
  - în PDF și Excel, cu editor vizual și rapoarte proprii pe SQL;
  - opțional, o pagină nouă pentru fiecare secție.
- **Exportul pentru salarizare:** CSV sau Excel, cu modele configurabile.
- **Controlul lunii:**
  - **Verifică luna**, înainte de închidere;
  - închiderea și redeschiderea lunii, cu motiv;
  - jurnalul modificărilor.
- **Siguranța datelor:**
  - conturi și roluri, date criptate;
  - backup zilnic și la cerere, restaurare, cheie de recuperare;
  - a doua fereastră deschisă doar în citire.
- **Limba:** română și engleză; traducerile se pot corecta și importa din Setări.
- **Ajutorul:** ghiduri în aplicație și manualul inclus, în MD, HTML și PDF. Asistentul AI este opțional.

## Înainte de actualizare

Urmați ACTUALIZARE.md:
- păstrați copia completă a datelor și cheia de recuperare;
- nu copiați numai baza de date;
- nu folosiți un executabil vechi pe datele migrate de o versiune nouă.

## Limite cunoscute

- **Executabilul nu este încă semnat digital.** La prima pornire, Windows poate afișa un avertisment SmartScreen.
- **Regulile de plată nu vin cu aplicația.** Procentele, sporurile și interpretările legale se configurează de unitate, împreună cu contabilitatea.
- **Exportul pentru salarizare are doar modelele generice** CSV și Excel. Formatul unui program de salarizare anume se adaugă după documentația lui de import, când programul e ales.
- **Unele întrebări despre raportul de referință sunt încă deschise:**
  - ce înseamnă celulele roșii și chenarele;
  - ce coloane lipsesc sau sunt în plus;
  - ce lună se folosește pentru comparația rând cu rând.
- **Valorile „—” nu se pot calcula.** Motivul apare la trecerea mouse-ului peste celulă, iar **Verifică luna** le adună înainte de închidere.
- **Indicatorul zilnic** nu stabilește necesarul de personal sau deficitul.
- **AI generativ** necesită un motor și un model compatibile, distribuite separat. Ghidurile din aplicație merg fără ele.
- **Fișa de acceptanță inclusă** este pentru verificarea pe calculatorul vostru, cu date fictive. Nu este o certificare.
