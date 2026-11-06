# D078 - Remedii UI si salvare verificata a exporturilor

Stare: acceptată. Data: 2026-10-31. F16-T29. Utilizatorul a autorizat remedierea auditului UI28 si apoi a cerut amanarea investigatiei blocajului Excel · Decisă de: Ioan Tuns, 2026-11-06 (F02-T25)

Identitatea de garda este afisata consecvent din atributul contractului sau asocierea explicita; nu se modifica numele salvat si nu se creeaza asocieri dupa nume. Se pastreaza criteriile existente de eligibilitate a programelor, dar se explica motivele excluderii. Se grupeaza formularele si filtrele, se afiseaza stari textuale si se pastreaza repere fixe la derulare. Panourile native details pastreaza navigarea cu tastatura, fara biblioteca noua. Legaturile de configurare respecta protectia modificarilor nesalvate si drepturile.

Salvarea XLSX verifica arhiva, CRC-urile si redeschiderea foilor; scrie si sincronizeaza un fisier temporar in acelasi folder, verifica bytes/structura, apoi redenumeste peste destinatie. La eroare destinatia precedenta nu se elimina. Esuarea dialogului/salvarii nu mai declanseaza in tacere descarcarea din browser. La testarea variantei de browser, URL-ul temporar nu este revocat imediat dupa click.

S-au deschis normal in Microsoft Excel registre anonime inclusiv cu 300 persoane; nu s-a reprodus blocajul utilizatorului. Investigatia specifica este amanata explicit, iar protectiile deja implementate se pastreaza. Nicio instalare, dependenta noua sau modificare de calcul salarial.
