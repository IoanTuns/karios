# D030 · Verificarea suprapunerilor pentru înlocuitor, evidențierea în pontaj, diferența aritmetică a concediului și previzualizarea cererii

Stare: propusă · Data: 2026-10-26 · Propusă de: Antigravity · Decisă de: —

## Context

După introducerea cererilor de concediu semnate, a persoanei înlocuitoare și a stării de aprobare (F07-T12, [D024](D024-cereri-concediu-aprobare.md)), în utilizarea curentă apar patru nevoi directe:

1. Prevenirea situației în care persoana desemnată ca înlocuitor este ea însăși în concediu/absență în aceeași zi (sau în care o persoană deja desemnată ca înlocuitor își înregistrează concediu fără a se observa suprapunerea), precum și semnalarea cazului în care aceeași persoană înlocuiește mai mulți colegi în aceeași zi.
2. Vizibilitatea în grila lunară de pontaj nu doar pe rândul angajatului absent, ci și pe rândul persoanei care îl înlocuiește, pentru a facilita planificarea turelor.
3. Afișarea diferenței aritmetice dintre dreptul anual contractual curent și zilele de concediu de odihnă înregistrate (plus câte sunt încă în așteptare), fără a inventa reguli legale de reportare sau indemnizație.
4. Previzualizarea rapidă în fereastră modală a cererii semnate atașate (PDF, PNG, JPEG), fără a obliga utilizatorul să salveze fișierul pe disc pentru o simplă verificare vizuală.

## Variante

1. **Verificare doar în interfață vs. invariant în serviciul Go și interfață.**
   - Doar în interfață: loturile trimise direct sau concurent pot lăsa un înlocuitor simultan absent în aceeași zi.
   - În serviciul Go + interfață: serviciul refuză tranzacțional suprapunerea dintre o absență activă (nerespinsă) și rolul de înlocuitor în aceeași zi, iar interfața afișează imediat atât blocajul, cât și avertismentul neblocant când aceeași persoană înlocuiește doi colegi în aceeași zi.
2. **Previzualizare externă vs. previzualizare modală în aplicație.**
   - Doar descărcare: lasă copii temporare ale cererilor în folderul Downloads al sistemului.
   - Previzualizare modală cu `Blob URL` în memorie (eliberat la închiderea ferestrei) + opțiune explicită de descărcare: păstrează documentul în memorie pentru verificare și oferă descărcarea doar când utilizatorul o cere.

## Decizie

1. **Suprapuneri înlocuitor (F07-T13, I025):**
   - O zi de absență activă cu `stare_aprobare != 'RESPINSA'` nu poate desemna ca `inlocuitor_id` un angajat care are el însuși o absență activă cu `stare_aprobare != 'RESPINSA'` în aceeași dată (în baza de date sau în același lot). Eroare: `absente.inlocuitor_absent`.
   - Simetric, un angajat nu poate înregistra o absență activă (`stare_aprobare != 'RESPINSA'`) într-o zi în care este deja desemnat ca `inlocuitor_id` pe o absență activă (`stare_aprobare != 'RESPINSA'`) a altui coleg, fără schimbarea prealabilă a înlocuitorului. Eroare: `absente.titular_este_inlocuitor`.
   - Dacă persoana aleasă ca înlocuitor înlocuiește deja un alt coleg în aceeași zi, salvarea rămâne permisă, dar formularul și tabelul afișează un avertisment vizibil.
2. **Evidențiere pe rândul înlocuitorului în Pontaj (I026):**
   - În grila lunară de pontaj, celulele zilelor în care un angajat înlocuiește unul sau mai mulți colegi (pentru absențe active nerespinse) primesc clasa `celula-inlocuitor` și afișează în explicația celulei (`title`) lista colegilor înlocuiți și codul absenței.
3. **Diferență aritmetică în sumarul anual (I027):**
   - În tabelul de sumar din Concedii și absențe, coloana dreptului anual afișează, lângă numărul de zile din contract, diferența aritmetică `Rămas: X zile` față de zilele de odihnă active (și numărul de zile `În așteptare`), cu evidențiere când diferența este negativă. Nu se deduc reportări din anii anteriori (acestea rămân în I010 / F01).
4. **Previzualizare modală a cererii semnate (I028):**
   - Tabelul de absențe oferă atât **Previzualizează** (fereastră modală cu `<iframe>` pentru PDF sau `<img>` pentru PNG/JPEG din `Blob URL` revocat la închidere), cât și **Descarcă**.

## Consecințe

- Nu este necesară o migrare nouă a bazei de date (schema din migrarea `0012_cereri_concediu` și interogarea existentă `ListaAbsenteInterval` acoperă toate verificările).
- Se aplică în F07-T13.
